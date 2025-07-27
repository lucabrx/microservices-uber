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

func (ts *TripSimulator) Start() {
	log.Println("Starting trip completion simulator...")
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
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
				_, err := ts.service.CompleteTrip(trip.ID)
				if err != nil {
					log.Printf("Simulator failed to complete trip %s: %v", trip.ID, err)
					continue
				}

				log.Printf("Simulator completed trip %s for driver %s", trip.ID, trip.DriverID)

			}
		}
	}()
}
