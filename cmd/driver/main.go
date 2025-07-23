package main

import (
	"log"
	"net"

	pb "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	"github.com/lukabrx/uber-clone/internal/driver"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Wiring: Repository -> Service -> Handler
	repo := driver.NewMemoryRepository()
	service := driver.NewService(repo)
	handler := driver.NewGrpcHandler(service)

	s := grpc.NewServer()
	pb.RegisterDriverServiceServer(s, handler)

	log.Println("Driver gRPC server listening at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
