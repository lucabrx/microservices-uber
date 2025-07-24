package trip

import (
	"context"

	pb "github.com/lukabrx/uber-clone/api/proto/trip/v1"
	"github.com/lukabrx/uber-clone/internal/models"
)

type GrpcHandler struct {
	service *Service
	pb.UnimplementedTripServiceServer
}

func NewGrpcHandler(service *Service) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	trip, err := h.service.CreateTrip(models.Trip{
		RiderID:  req.RiderId,
		DriverID: req.DriverId,
		StartLat: req.StartLat,
		StartLon: req.StartLon,
		EndLat:   req.EndLat,
		EndLon:   req.EndLon,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTripResponse{Trip: &pb.Trip{Id: trip.ID, RiderId: trip.RiderID, DriverId: trip.DriverID, Status: trip.Status, Price: trip.Price}}, nil
}

func (h *GrpcHandler) CompleteTrip(ctx context.Context, req *pb.CompleteTripRequest) (*pb.CompleteTripResponse, error) {
	trip, err := h.service.CompleteTrip(req.GetTripId())
	if err != nil {
		return nil, err
	}

	pbTrip := &pb.Trip{
		Id:       trip.ID,
		RiderId:  trip.RiderID,
		DriverId: trip.DriverID,
		Status:   trip.Status,
		Price:    trip.Price,
	}

	return &pb.CompleteTripResponse{Trip: pbTrip}, nil
}
