package gateway

import (
	"encoding/json"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lukabrx/uber-clone/internal/models"
)

var (
	drivers = make(map[string]models.Driver)
	trips   = make(map[string]models.Trip)
	mu      sync.Mutex
)

func RegisterDriver(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var driver models.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	driver.ID = uuid.New().String()
	driver.IsAvailable = true
	drivers[driver.ID] = driver

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(driver)
}

func UnregisterDriver(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	delete(drivers, id)
	w.WriteHeader(http.StatusNoContent)
}

func CheckDriverAvailability(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	driver, exists := drivers[id]
	if !exists {
		http.Error(w, "Driver not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"is_available": driver.IsAvailable})
}

func BookTrip(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var trip models.Trip
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For simplicity, find the first available driver
	var assignedDriverID string
	for _, driver := range drivers {
		if driver.IsAvailable {
			assignedDriverID = driver.ID
			break
		}
	}

	if assignedDriverID == "" {
		http.Error(w, "No available drivers", http.StatusServiceUnavailable)
		return
	}

	trip.ID = uuid.New().String()
	trip.DriverID = assignedDriverID
	trip.Status = "in_progress"
	trip.RequestTime = time.Now()

	// Simulate price calculation using a simple distance formula
	trip.Price = calculatePrice(trip.StartLat, trip.StartLon, trip.EndLat, trip.EndLon)

	trips[trip.ID] = trip

	// Make the assigned driver unavailable
	driver := drivers[assignedDriverID]
	driver.IsAvailable = false
	drivers[assignedDriverID] = driver

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(trip)
}

// calculatePrice simulates price calculation based on distance.
// This is where you would integrate with a service like OSRM.
func calculatePrice(lat1, lon1, lat2, lon2 float64) float64 {
	// A very simple distance calculation (not geographically accurate)
	distance := math.Sqrt(math.Pow(lat2-lat1, 2) + math.Pow(lon2-lon1, 2))
	baseFare := 2.50
	perKmRate := 1.50
	return baseFare + (distance * 100 * perKmRate)
}
