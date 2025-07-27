package trip

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lukabrx/uber-clone/internal/types"
)

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

func (kp *KafkaProducer) ProduceTripCreated(tripID, driverID string) {
	event := types.TripEvent{
		EventType: types.TripCreatedEvent,
		TripID:    tripID,
		DriverID:  driverID,
	}
	value, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal TripCreatedEvent: %v", err)
		return
	}

	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &types.TripEventsTopic, Partition: kafka.PartitionAny},
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

func (kp *KafkaProducer) ProduceTripCompleted(tripID, driverID string) {
	event := types.TripEvent{
		EventType: types.TripCompletedEvent,
		TripID:    tripID,
		DriverID:  driverID,
	}
	value, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal TripCompletedEvent: %v", err)
		return
	}
	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &types.TripEventsTopic, Partition: kafka.PartitionAny},
		Value:          value,
		Key:            []byte(event.TripID),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce TripCompletedEvent: %v", err)
		return
	}
}
