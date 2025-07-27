package driver

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lukabrx/uber-clone/internal/models"
	"github.com/lukabrx/uber-clone/internal/types"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(bootstrapServers string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer}, nil
}

func (kp *KafkaProducer) ProduceAvailableDriverUpdate(driver models.Driver) {
	value, err := json.Marshal(driver)
	if err != nil {
		log.Printf("Failed to marshal driver location: %v", err)
		return
	}

	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &types.DriverLocationTopic, Partition: kafka.PartitionAny},
		Value:          value,
		Key:            []byte(driver.ID),
	}, nil)

	if err != nil {
		log.Printf("Failed to produce message: %v", err)
	}
}

func (kp *KafkaProducer) Close() {
	// Wait up to 15 seconds for all messages to be sent then close the producer
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}
