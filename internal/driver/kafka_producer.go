package driver

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lukabrx/uber-clone/internal/models"
)

var DriverLocationTopic = "driver_locations"

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

func (kp *KafkaProducer) ProduceLocationUpdate(driver models.Driver) {
	value, err := json.Marshal(driver)
	if err != nil {
		log.Printf("Failed to marshal driver location: %v", err)
		return
	}

	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &DriverLocationTopic, Partition: kafka.PartitionAny},
		Value:          value,
		Key:            []byte(driver.ID),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
}

func (kp *KafkaProducer) Close() {
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}
