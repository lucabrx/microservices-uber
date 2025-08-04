package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	pb "github.com/lukabrx/uber-clone/api/proto/auth/v1"
	"github.com/lukabrx/uber-clone/internal/auth"
	"github.com/lukabrx/uber-clone/internal/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found or error loading .env file:", err)
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID environment variable not set")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET environment variable not set")
	}

	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		log.Fatal("GOOGLE_REDIRECT_URL environment variable not set")
	}

	pasetoSymmetricKey := os.Getenv("PASETO_SYMMETRIC_KEY")
	if pasetoSymmetricKey == "" {
		log.Fatal("PASETO_SYMMETRIC_KEY environment variable not set")
	}

	googleOauthConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	pasetoMaker, err := auth.NewPasetoMaker(pasetoSymmetricKey)
	if err != nil {
		log.Fatalf("failed to create paseto maker: %v", err)
	}

	userRepo := user.NewMemoryRepository()
	refreshTokenRepo := auth.NewRefreshTokenRepository()

	service := auth.NewService(pasetoMaker, googleOauthConfig, userRepo, refreshTokenRepo)
	handler := auth.NewGrpcHandler(service)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen on port 50053: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, handler)

	log.Println("Auth gRPC server listening at :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve auth gRPC server: %v", err)
	}
}
