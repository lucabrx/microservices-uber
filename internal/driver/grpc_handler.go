package driver

import (
	"context"

	pb "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	"github.com/lukabrx/uber-clone/internal/models"
)

type GrpcHandler struct {
	service *Service
	pb.UnimplementedDriverServiceServer
}

func NewGrpcHandler(service *Service) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driver, err := h.service.RegisterDriver(models.Driver{Name: req.Name, Lat: req.Lat, Lon: req.Lon})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterDriverResponse{Driver: &pb.Driver{Id: driver.ID, Name: driver.Name, Lat: driver.Lat, Lon: driver.Lon}}, nil
}

func (h *GrpcHandler) FindAvailableDrivers(ctx context.Context, req *pb.FindAvailableDriversRequest) (*pb.FindAvailableDriversResponse, error) {
	drivers := h.service.FindClosestAvailableDrivers(req.Lat, req.Lon)
	var pbDrivers []*pb.Driver
	for _, d := range drivers {
		pbDrivers = append(pbDrivers, &pb.Driver{Id: d.ID, Name: d.Name, Lat: d.Lat, Lon: d.Lon})
	}
	return &pb.FindAvailableDriversResponse{Drivers: pbDrivers}, nil
}

func (h *GrpcHandler) UpdateDriverStatus(ctx context.Context, req *pb.UpdateDriverStatusRequest) (*pb.UpdateDriverStatusResponse, error) {
	err := h.service.UpdateDriverStatus(req.Id, req.IsAvailable)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateDriverStatusResponse{}, nil
}
