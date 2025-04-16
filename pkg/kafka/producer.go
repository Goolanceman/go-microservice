package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/goolanceman/go-microservice/internal/config"
	"github.com/goolanceman/go-microservice/pkg/logger"
)

// Producer handles Kafka message publishing
type Producer struct {
	writer *kafka.Writer
	config *config.KafkaConfig
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.KafkaConfig) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.ProducerTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        true,
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
	}); err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}

	logger.Info("Successfully connected to Kafka producer",
		logger.Strings("brokers", cfg.Brokers),
		logger.String("topic", cfg.ProducerTopic),
	)

	return &Producer{
		writer: writer,
		config: cfg,
	}, nil
}

// Publish sends a message to Kafka
func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// HealthCheck performs a health check on the Kafka producer
func (p *Producer) HealthCheck(ctx context.Context) error {
	return p.Publish(ctx, []byte("healthcheck"), []byte("ping"))
} 