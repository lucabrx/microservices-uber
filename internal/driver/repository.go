package driver

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/lukabrx/uber-clone/internal/models"
)

type MemoryRepository struct {
	drivers map[string]models.Driver
	mu      sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		drivers: make(map[string]models.Driver),
	}
}

func (r *MemoryRepository) RegisterDriver(driver models.Driver) (models.Driver, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	driver.ID = uuid.New().String()
	driver.IsAvailable = true
	r.drivers[driver.ID] = driver

	return driver, nil
}

func (r *MemoryRepository) UpdateDriverStatus(id string, isAvailable bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	driver, ok := r.drivers[id]
	if !ok {
		return errors.New("driver not found")
	}
	driver.IsAvailable = isAvailable
	r.drivers[id] = driver
	return nil
}

func (r *MemoryRepository) GetAvailableDrivers() []models.Driver {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []models.Driver
	for _, driver := range r.drivers {
		if driver.IsAvailable {
			available = append(available, driver)
		}
	}
	return available
}

func (r *MemoryRepository) GetAllDrivers() []models.Driver {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []models.Driver
	for _, driver := range r.drivers {
		all = append(all, driver)
	}
	return all
}
