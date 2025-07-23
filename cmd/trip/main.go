package main

import (
	"log"
	"net"

	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	pb_trip "github.com/lukabrx/uber-clone/api/proto/trip/v1"
	"github.com/lukabrx/uber-clone/internal/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to Driver service
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to driver service: %v", err)
	}
	defer conn.Close()
	driverClient := pb_driver.NewDriverServiceClient(conn)

	repo := trip.NewMemoryRepository()
	service := trip.NewService(repo, driverClient)
	handler := trip.NewGrpcHandler(service)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb_trip.RegisterTripServiceServer(s, handler)

	log.Println("Trip gRPC server listening at :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
