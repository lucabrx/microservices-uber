package driver

import (
	"log"
	"math"
	"sort"

	"github.com/lukabrx/uber-clone/internal/models"
)

type Service struct {
	repo     *MemoryRepository
	producer *KafkaProducer
}

func NewService(repo *MemoryRepository, producer *KafkaProducer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) RegisterDriver(d models.Driver) (*models.Driver, error) {
	driver, err := s.repo.RegisterDriver(d)
	if err != nil {
		return nil, err
	}

	log.Printf("New driver registered %s, publishing availability update", driver.ID)
	s.producer.ProduceAvailableDriverUpdate(driver)

	return &driver, nil
}

func (s *Service) UpdateDriverStatus(id string, isAvailable bool) error {
	err := s.repo.UpdateDriverStatus(id, isAvailable)
	if err != nil {
		return err
	}

	driver, err := s.repo.GetDriverByID(id)
	if err != nil {
		return err
	}

	log.Printf("Driver %s status updated to available: %v, publishing update", id, isAvailable)
	s.producer.ProduceAvailableDriverUpdate(*driver)

	return nil
}

func (s *Service) FindClosestAvailableDrivers(lat, lon float64) []models.Driver {
	drivers := s.repo.GetAvailableDrivers()

	sort.Slice(drivers, func(i, j int) bool {
		distI := calculateDistance(lat, lon, drivers[i].Lat, drivers[i].Lon)
		distJ := calculateDistance(lat, lon, drivers[j].Lat, drivers[j].Lon)
		return distI < distJ
	})

	return drivers
}

// calculateDistance calculates the distance between two points on Earth.
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	lat1_rad := lat1 * (math.Pi / 180.0)
	lat2_rad := lat2 * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1_rad)*math.Cos(lat2_rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (s *Service) IsDriverAvailable(id string) (bool, error) {
	return s.repo.IsDriverAvailable(id)
}
