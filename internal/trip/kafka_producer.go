package trip

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var TripEventsTopic = "trip_events"

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(bootstrapServers string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: p}, nil
}

type TripCreatedEvent struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
}

func (kp *KafkaProducer) ProduceTripCreated(event TripCreatedEvent) {
	value, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal TripCreatedEvent: %v", err)
		return
	}

	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &TripEventsTopic, Partition: kafka.PartitionAny},
		Value:          value,
		Key:            []byte(event.TripID),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce TripCreatedEvent: %v", err)
		return
	}
}

func (kp *KafkaProducer) Close() {
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}
