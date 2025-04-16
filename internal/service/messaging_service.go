package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goolanceman/go-microservice/internal/config"
	"github.com/goolanceman/go-microservice/pkg/messaging"
	"go.uber.org/zap"
)

// Message represents a Kafka message
type Message struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// MessagingService handles Kafka operations
type MessagingService struct {
	kafkaClient messaging.Messaging
	logger      *zap.Logger
	config      *config.Config
}

// NewMessagingService creates a new messaging service
func NewMessagingService(kafkaClient messaging.Messaging, logger *zap.Logger, cfg *config.Config) *MessagingService {
	return &MessagingService{
		kafkaClient: kafkaClient,
		logger:      logger,
		config:      cfg,
	}
}

// StartConsumers starts Kafka consumers for all configured topics
func (s *MessagingService) StartConsumers(ctx context.Context) error {
	// Subscribe to all configured topics
	for _, topic := range s.config.Kafka.Topics {
		err := s.kafkaClient.Subscribe(ctx, topic, func(ctx context.Context, msg *messaging.Message) error {
			s.logger.Info("Received message",
				zap.String("topic", msg.Topic),
				zap.String("payload", string(msg.Payload)))

			// Process message based on topic
			switch msg.Topic {
			case "topic1":
				return s.handleTopic1(msg)
			case "topic2":
				return s.handleTopic2(msg)
			case "topic3":
				return s.handleTopic3(msg)
			default:
				s.logger.Warn("Unknown topic",
					zap.String("topic", msg.Topic))
				return nil
			}
		})

		if err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
		}
	}

	return nil
}

// PublishMessage publishes a message to a Kafka topic
func (s *MessagingService) PublishMessage(ctx context.Context, topic string, payload interface{}) error {
	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create message
	msg := &messaging.Message{
		Topic:   topic,
		Payload: jsonPayload,
	}

	// Publish message
	if err := s.kafkaClient.Publish(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Topic-specific handlers
func (s *MessagingService) handleTopic1(msg *messaging.Message) error {
	s.logger.Info("Processing topic1 message",
		zap.String("payload", string(msg.Payload)))
	// Add your topic1 specific logic here
	return nil
}

func (s *MessagingService) handleTopic2(msg *messaging.Message) error {
	s.logger.Info("Processing topic2 message",
		zap.String("payload", string(msg.Payload)))
	// Add your topic2 specific logic here
	return nil
}

func (s *MessagingService) handleTopic3(msg *messaging.Message) error {
	s.logger.Info("Processing topic3 message",
		zap.String("payload", string(msg.Payload)))
	// Add your topic3 specific logic here
	return nil
} 