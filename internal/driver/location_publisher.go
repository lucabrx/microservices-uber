package driver

import (
	"log"
	"time"
)

type LocationUpdater struct {
	service  *Service
	producer *KafkaProducer
}

func NewLocationUpdater(service *Service, producer *KafkaProducer) *LocationUpdater {
	return &LocationUpdater{service: service, producer: producer}
}

func (lu *LocationUpdater) Start() {
	log.Println("Starting location updater...")
	// Run this in a separate goroutine so it doesn't block the main thread
	go func() {
		ticker := time.NewTicker(60 * time.Second) // Publish updates every 5 seconds
		defer ticker.Stop()

		for range ticker.C {
			drivers := lu.service.repo.GetAllDrivers()
			for _, driver := range drivers {
				// In a real app, you'd update the driver's Lat/Lon here
				log.Printf("Publishing location for driver %s", driver.ID)
				lu.producer.ProduceLocationUpdate(driver)
			}
		}
	}()
}
