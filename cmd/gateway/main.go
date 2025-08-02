package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	pb_auth "github.com/lukabrx/uber-clone/api/proto/auth/v1"
	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	pb_trip "github.com/lukabrx/uber-clone/api/proto/trip/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/lukabrx/uber-clone/internal/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found or error loading .env file:", err)
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID environment variable not set")
	}

	fmt.Println("Using Google Client ID:", googleClientID)

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET environment variable not set")
	}

	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		log.Fatal("GOOGLE_REDIRECT_URL environment variable not set")
	}

	driverConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to driver service: %v", err)
	}
	defer driverConn.Close()

	tripConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to trip service: %v", err)
	}
	defer tripConn.Close()

	authConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to auth service: %v", err)
	}
	defer authConn.Close()

	tripClient := pb_trip.NewTripServiceClient(tripConn)
	driverClient := pb_driver.NewDriverServiceClient(driverConn)
	authClient := pb_auth.NewAuthServiceClient(authConn)

	hub := gateway.NewHub(driverClient)
	googleOauthConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	httpHandler := gateway.NewHttpHandler(driverClient, tripClient, authClient, hub, googleOauthConfig)

	kafkaConsumer, err := gateway.NewKafkaConsumer("localhost:29092", "gateway_group", hub)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer for gateway: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	kafkaConsumer.SubscribeAndListen(ctx)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/auth/google/login", httpHandler.HandleGoogleLogin)
	r.Get("/auth/google/callback", httpHandler.HandleGoogleCallback)
	r.Get("/ws/drivers/available", httpHandler.StreamAvailableDrivers)

	r.Group(func(r chi.Router) {
		r.Use(httpHandler.AuthMiddleware)

		r.Post("/drivers", httpHandler.RegisterDriver)
		r.Get("/drivers/available", httpHandler.FindAvailableDrivers)
		r.Post("/trips", httpHandler.CreateTrip)
		r.Patch("/trips/{id}/complete", httpHandler.CompleteTrip)
	})

	log.Println("Gateway server starting on :8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
