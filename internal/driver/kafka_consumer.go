package driver

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const TripEventsTopic = "trip_events"

type TripCreatedEvent struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
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
			if err == nil {
				log.Printf("Received message from topic %s: %s\n", *msg.TopicPartition.Topic, string(msg.Value))

				var event TripCreatedEvent
				if err := json.Unmarshal(msg.Value, &event); err == nil {
					kc.service.UpdateDriverStatus(event.DriverID, false) // Mark as unavailable
				}
			}
		}
	}()
}
