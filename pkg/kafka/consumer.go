package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go-microservice/internal/config"
	"go.uber.org/zap"
)

// MessageHandler defines the interface for processing Kafka messages
type MessageHandler interface {
	HandleMessage(ctx context.Context, key, value []byte) error
}

// Consumer handles Kafka message consumption
type Consumer struct {
	reader  *kafka.Reader
	config  *config.KafkaConfig
	handler MessageHandler
	logger  *zap.Logger
	wg      sync.WaitGroup
	done    chan struct{}
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.KafkaConfig, handler MessageHandler, logger *zap.Logger, topic string) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		GroupID:     cfg.ConsumerGroup,
		Topic:       topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
	})

	logger.Info("Successfully created Kafka consumer",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
		zap.String("group", cfg.ConsumerGroup),
	)

	return &Consumer{
		reader:  reader,
		config:  cfg,
		handler: handler,
		logger:  logger,
		done:    make(chan struct{}),
	}, nil
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) {
	c.wg.Add(1)
	go c.consumeLoop(ctx)
}

// Stop gracefully stops the consumer
func (c *Consumer) Stop() {
	close(c.done)
	c.wg.Wait()
	if err := c.reader.Close(); err != nil {
		c.logger.Error("Failed to close Kafka reader", zap.Error(err))
	}
}

// consumeLoop continuously reads and processes messages
func (c *Consumer) consumeLoop(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-c.done:
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("Failed to read message", zap.Error(err))
				time.Sleep(time.Second) // Prevent tight loop on error
				continue
			}

			if err := c.handler.HandleMessage(ctx, msg.Key, msg.Value); err != nil {
				c.logger.Error("Failed to handle message",
					zap.Error(err),
					zap.String("key", string(msg.Key)),
				)
				continue
			}

			c.logger.Debug("Successfully processed message",
				zap.String("key", string(msg.Key)),
				zap.Int("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
			)
		}
	}
}

// HealthCheck performs a health check on the Kafka consumer
func (c *Consumer) HealthCheck(ctx context.Context) error {
	// Check if we can read from the topic
	_, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to read from Kafka: %w", err)
	}
	return nil
} 