package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go-microservice/internal/config"
	"go.uber.org/zap"
)

// Producer handles Kafka message publishing
type Producer struct {
	writer *kafka.Writer
	config *config.KafkaConfig
	logger *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.KafkaConfig, logger *zap.Logger, topic string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        topic,
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
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
	)

	return &Producer{
		writer: writer,
		config: cfg,
		logger: logger,
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