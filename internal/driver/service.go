package driver

import (
	"math"
	"sort"

	"github.com/lukabrx/uber-clone/internal/models"
)

type Service struct {
	repo *MemoryRepository
}

func NewService(repo *MemoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterDriver(driver models.Driver) (models.Driver, error) {
	return s.repo.RegisterDriver(driver)
}

func (s *Service) UpdateDriverStatus(id string, isAvailable bool) error {
	return s.repo.UpdateDriverStatus(id, isAvailable)
}

func (s *Service) FindClosestAvailableDrivers(lat, lon float64) []models.Driver {
	drivers := s.repo.GetAvailableDrivers()

	sort.Slice(drivers, func(i, j int) bool {
		distI := haversine(lat, lon, drivers[i].Lat, drivers[i].Lon)
		distJ := haversine(lat, lon, drivers[j].Lat, drivers[j].Lon)
		return distI < distJ
	})

	return drivers
}

// haversine calculates the distance between two points on Earth.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
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
