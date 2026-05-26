package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer wraps Kafka writer for publishing events
type Producer struct {
	writer *kafka.Writer
}

// Config holds Kafka producer configuration
type Config struct {
	Brokers []string
	Topic   string
}

// Event represents a generic event structure for Kafka
type Event struct {
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// NewProducer creates a new Kafka producer
func NewProducer(config Config) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	// Test connection
	conn, err := kafka.Dial("tcp", config.Brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	return &Producer{
		writer: writer,
	}, nil
}

// PublishEvent publishes a JSON event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, event Event) error {
	// Marshal event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Write to Kafka
	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(event.Type),
			Value: eventBytes,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Event published: type=%s, data=%v\n", event.Type, event.Data)
	return nil
}

// Close closes the Kafka writer
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
