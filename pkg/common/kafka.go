package common

import (
	"context"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

// Kafka bağlantısını başlatır (singleton)
func InitKafka() error {
	addr := os.Getenv("KAFKA_ADDR")
	if addr == "" {
		addr = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "crm-events"
	}
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return nil
}

// Kafka'ya event publish eder
func PublishEvent(ctx context.Context, key, value string) error {
	return kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	})
}
