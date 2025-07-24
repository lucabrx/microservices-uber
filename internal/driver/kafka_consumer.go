package driver

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const TripEventsTopic = "trip_events"

type TripEvent struct {
	EventType string `json:"event_type"`
	TripID    string `json:"trip_id"`
	DriverID  string `json:"driver_id"`
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
	service  *Service
}

func NewKafkaConsumer(bootstrapServers, groupID string, service *Service) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c, service: service}, nil
}

func (kc *KafkaConsumer) SubscribeAndListen() {
	err := kc.consumer.SubscribeTopics([]string{TripEventsTopic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	log.Println("Subscribed to topic and listening for trip events...")
	go func() {
		for {
			msg, err := kc.consumer.ReadMessage(-1)
			if err != nil {
				// handle error
				continue
			}

			log.Printf("Received message from topic %s", *msg.TopicPartition.Topic)

			var event TripEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Could not unmarshal event: %v", err)
				continue
			}

			switch event.EventType {
			case "TRIP_CREATED":
				log.Printf("Processing trip creation for driver %s", event.DriverID)
				err := kc.service.UpdateDriverStatus(event.DriverID, false)
				if err != nil {
					log.Printf("Error updating driver status for trip creation: %v", err)
				}

			case "TRIP_COMPLETED":
				log.Printf("Processing trip completion for driver %s", event.DriverID)
				err := kc.service.UpdateDriverStatus(event.DriverID, true)
				if err != nil {
					log.Printf("Error updating driver status for trip completion: %v", err)
				}

			default:
				log.Printf("Unknown event type received: %s", event.EventType)
			}
		}
	}()
}
