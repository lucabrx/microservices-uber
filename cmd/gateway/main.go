package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	pb_trip "github.com/lukabrx/uber-clone/api/proto/trip/v1"
	"github.com/lukabrx/uber-clone/internal/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	driverConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to driver service: %v", err)
	}
	defer driverConn.Close()
	driverClient := pb_driver.NewDriverServiceClient(driverConn)

	tripConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to trip service: %v", err)
	}
	defer tripConn.Close()
	tripClient := pb_trip.NewTripServiceClient(tripConn)

	httpHandler := gateway.NewHttpHandler(driverClient, tripClient)

	r := chi.NewRouter()
	r.Post("/drivers", httpHandler.RegisterDriver)
	r.Get("/drivers/available", httpHandler.FindAvailableDrivers)
	r.Post("/trips", httpHandler.CreateTrip)
	r.Patch("/trips/{id}/complete", httpHandler.CompleteTrip)

	log.Println("Gateway server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
