package trip

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/lukabrx/uber-clone/internal/models"
)

type MemoryRepository struct {
	trips map[string]*models.Trip
	mu    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		trips: make(map[string]*models.Trip),
	}
}

func (r *MemoryRepository) CreateTrip(trip models.Trip) (models.Trip, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	trip.ID = uuid.New().String()
	r.trips[trip.ID] = &trip

	return trip, nil
}

func (r *MemoryRepository) GetTripByID(id string) (models.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	trip, ok := r.trips[id]
	if !ok {
		return models.Trip{}, errors.New("trip not found")
	}
	return *trip, nil
}

func (r *MemoryRepository) UpdateTrip(trip models.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.trips[trip.ID]
	if !ok {
		return errors.New("trip not found")
	}
	*existing = trip
	return nil
}

func (r *MemoryRepository) GetInProgressTrips() ([]models.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var inProgressTrips []models.Trip
	for _, trip := range r.trips {
		if trip.Status == "in_progress" {
			inProgressTrips = append(inProgressTrips, *trip)
		}
	}
	return inProgressTrips, nil
}
