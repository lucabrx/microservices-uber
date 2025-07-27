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

	kafkaProducer, err := driver.NewKafkaProducer("localhost:29092")
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Wiring: Repository -> Service -> Handler
	repo := driver.NewMemoryRepository()
	service := driver.NewService(repo, kafkaProducer)
	handler := driver.NewGrpcHandler(service)

	s := grpc.NewServer()
	pb.RegisterDriverServiceServer(s, handler)

	kafkaConsumer, err := driver.NewKafkaConsumer("localhost:29092", "driver_service_group", service)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	kafkaConsumer.SubscribeAndListen()

	log.Println("Driver gRPC server listening at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
