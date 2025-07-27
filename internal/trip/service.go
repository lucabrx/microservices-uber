package trip

import (
	"context"

	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	"github.com/lukabrx/uber-clone/internal/models"
	pricecalculator "github.com/lukabrx/uber-clone/internal/price_calculator"
)

type Service struct {
	repo          *MemoryRepository
	driverClient  pb_driver.DriverServiceClient
	kafkaProducer *KafkaProducer
}

func NewService(repo *MemoryRepository, driverClient pb_driver.DriverServiceClient, kafkaProducer *KafkaProducer) *Service {
	return &Service{repo: repo, driverClient: driverClient, kafkaProducer: kafkaProducer}
}

func (s *Service) CreateTrip(req models.Trip) (models.Trip, error) {
	price, err := pricecalculator.CalculatePrice(req.StartLat, req.StartLon, req.EndLat, req.EndLon)
	if err != nil {
		return models.Trip{}, err
	}

	trip := models.Trip{
		RiderID:  req.RiderID,
		DriverID: req.DriverID,
		Status:   "in_progress",
		Price:    price,
	}
	createdTrip, err := s.repo.CreateTrip(trip)
	if err != nil {
		return models.Trip{}, err
	}

	s.kafkaProducer.ProduceTripCreated(createdTrip.ID, createdTrip.DriverID)

	return createdTrip, nil
}

func (s *Service) CompleteTrip(tripID string) (models.Trip, error) {
	trip, err := s.repo.GetTripByID(tripID)
	if err != nil {
		return models.Trip{}, err
	}

	trip.Status = models.TripStatusCompleted
	if err := s.repo.UpdateTrip(trip); err != nil {
		return models.Trip{}, err
	}

	updateReq := &pb_driver.UpdateDriverStatusRequest{
		Id:          trip.DriverID,
		IsAvailable: true, // The driver is now free
	}
	_, err = s.driverClient.UpdateDriverStatus(context.Background(), updateReq)
	if err != nil {
		return models.Trip{}, err
	}

	s.kafkaProducer.ProduceTripCompleted(trip.ID, trip.DriverID)

	return trip, nil
}
func (s *Service) GetInProgressTrips() ([]models.Trip, error) {
	return s.repo.GetInProgressTrips()
}
