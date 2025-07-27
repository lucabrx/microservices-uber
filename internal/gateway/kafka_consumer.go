package gateway

import (
	"context"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	"github.com/lukabrx/uber-clone/internal/models"
	"github.com/lukabrx/uber-clone/internal/types"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	hub      *Hub
}

func NewKafkaConsumer(bootstrapServers, groupID string, hub *Hub) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c, hub: hub}, nil
}

func (kc *KafkaConsumer) SubscribeAndListen(ctx context.Context) {
	err := kc.consumer.SubscribeTopics([]string{types.DriverLocationTopic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	log.Println("Gateway consumer subscribed and listening for driver location updates...")
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping gateway kafka consumer.")
				kc.consumer.Close()
				return
			default:
				msg, err := kc.consumer.ReadMessage(100) // 100ms timeout for reading messages
				if err != nil {
					// Ignore timeout errors
					if err.(kafka.Error).Code() == kafka.ErrTimedOut {
						continue
					}
					log.Printf("Consumer error: %v (%v)\n", err, msg)
					continue
				}

				log.Printf("Received driver location update from topic %s", *msg.TopicPartition.Topic)

				var driver models.Driver
				if err := json.Unmarshal(msg.Value, &driver); err != nil {
					log.Printf("Could not unmarshal driver data: %v", err)
					continue
				}

				res, err := kc.hub.driverClient.FindAvailableDrivers(context.Background(), &pb_driver.FindAvailableDriversRequest{})
				if err != nil {
					log.Printf("Failed to find available drivers after update: %v", err)
					continue
				}
				kc.hub.Broadcast(res.Drivers)
			}
		}
	}()
}
