package trip

import (
	"log"
	"time"
)

type TripSimulator struct {
	service  *Service
	producer *KafkaProducer
}

func NewTripSimulator(service *Service, producer *KafkaProducer) *TripSimulator {
	return &TripSimulator{
		service:  service,
		producer: producer,
	}
}

// Start begins the simulation loop.
func (ts *TripSimulator) Start() {
	log.Println("Starting trip completion simulator...")
	// Run this in a separate goroutine so it doesn't block the main thread.
	go func() {
		// Set the ticker to run every 2 minutes.
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Simulator checking for in-progress trips...")
			inProgressTrips, err := ts.service.GetInProgressTrips()
			if err != nil {
				log.Printf("Simulator failed to get in-progress trips: %v", err)
				continue
			}

			if len(inProgressTrips) == 0 {
				log.Println("No in-progress trips to complete.")
				continue
			}

			for _, trip := range inProgressTrips {
				// 1. Mark the trip as completed in the trip service's database.
				_, err := ts.service.CompleteTrip(trip.ID)
				if err != nil {
					log.Printf("Simulator failed to complete trip %s: %v", trip.ID, err)
					continue
				}

				log.Printf("Simulator completed trip %s for driver %s", trip.ID, trip.DriverID)

				// 2. Produce an event to Kafka to notify other services.
				ts.producer.ProduceTripCompleted(trip.ID, trip.DriverID)
			}
		}
	}()
}
