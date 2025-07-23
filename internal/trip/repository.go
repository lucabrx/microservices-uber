package trip

import (
	"sync"

	"github.com/google/uuid"
	"github.com/lukabrx/uber-clone/internal/models"
)

type MemoryRepository struct {
	trips map[string]models.Trip
	mu    sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		trips: make(map[string]models.Trip),
	}
}

func (r *MemoryRepository) CreateTrip(trip models.Trip) (models.Trip, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	trip.ID = uuid.New().String()
	r.trips[trip.ID] = trip

	return trip, nil
}
