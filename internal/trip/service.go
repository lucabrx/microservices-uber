package trip

import (
	"context"
	"errors"

	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	"github.com/lukabrx/uber-clone/internal/models"
	pricecalculator "github.com/lukabrx/uber-clone/internal/price_calculator"
)

type Service struct {
	repo         *MemoryRepository
	driverClient pb_driver.DriverServiceClient
}

func NewService(repo *MemoryRepository, driverClient pb_driver.DriverServiceClient) *Service {
	return &Service{repo: repo, driverClient: driverClient}
}

func (s *Service) CreateTrip(req models.Trip) (models.Trip, error) {
	findReq := &pb_driver.FindAvailableDriversRequest{Lat: req.StartLat, Lon: req.StartLon}
	findResp, err := s.driverClient.FindAvailableDrivers(context.Background(), findReq)
	if err != nil {
		return models.Trip{}, err
	}
	if len(findResp.Drivers) == 0 {
		return models.Trip{}, errors.New("no available drivers")
	}
	closestDriver := findResp.Drivers[0]

	price, err := pricecalculator.CalculatePrice(req.StartLat, req.StartLon, req.EndLat, req.EndLon)
	if err != nil {
		return models.Trip{}, err
	}

	updateReq := &pb_driver.UpdateDriverStatusRequest{Id: closestDriver.Id, IsAvailable: false}
	_, err = s.driverClient.UpdateDriverStatus(context.Background(), updateReq)
	if err != nil {
		return models.Trip{}, err
	}

	trip := models.Trip{
		RiderID:  req.RiderID,
		DriverID: closestDriver.Id,
		Status:   "in_progress",
		Price:    price,
	}
	return s.repo.CreateTrip(trip)
}
