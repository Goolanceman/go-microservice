package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go-microservice/internal/config"
	"go-microservice/pkg/logger"
	"go.uber.org/zap"
)

// Kafka implements the Messaging interface
type Kafka struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	mu       sync.Mutex
}

// NewKafka creates a new Kafka client
func NewKafka(cfg *config.Config) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create producer
	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Create consumer group
	group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, config)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Kafka{
		producer: producer,
		consumer: group,
	}, nil
}

// Publish publishes a message to a topic
func (k *Kafka) Publish(ctx context.Context, msg *Message) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.producer == nil {
		return ErrNotConnected
	}

	// Create Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(msg.Payload),
	}

	// Add headers
	if len(msg.Headers) > 0 {
		var headers []sarama.RecordHeader
		for key, value := range msg.Headers {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(key),
				Value: []byte(value),
			})
		}
		kafkaMsg.Headers = headers
	}

	// Send message
	partition, offset, err := k.producer.SendMessage(kafkaMsg)
	if err != nil {
		logger.Error("Failed to send message", zap.Error(err))
		return fmt.Errorf("failed to send message: %w", err)
	}

	logger.Debug("Message published successfully",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))

	return nil
}

// SubscribeMultiple subscribes to multiple topics with the same handler
func (k *Kafka) SubscribeMultiple(ctx context.Context, topics []string, handler Handler) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.consumer == nil {
		return ErrNotConnected
	}

	// Create consumer group handler
	h := &consumerGroupHandler{
		handler: handler,
	}

	// Start consuming in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := k.consumer.Consume(ctx, topics, h); err != nil {
					logger.Error("Error from consumer", zap.Error(err))
				}
			}
		}
	}()

	return nil
}

// Subscribe subscribes to a single topic
func (k *Kafka) Subscribe(ctx context.Context, topic string, handler Handler) error {
	return k.SubscribeMultiple(ctx, []string{topic}, handler)
}

// Close closes the Kafka connections
func (k *Kafka) Close() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	var errs []error

	if k.producer != nil {
		if err := k.producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
		}
	}

	if k.consumer != nil {
		if err := k.consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple close errors: %v", errs)
	}

	return nil
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler Handler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			if msg == nil {
				return nil
			}

			logger.Debug("Received message",
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset))

			message := &Message{
				Topic:   msg.Topic,
				Payload: msg.Value,
				Headers: make(map[string]string),
			}

			// Copy headers
			for _, header := range msg.Headers {
				message.Headers[string(header.Key)] = string(header.Value)
			}

			// Call handler
			if err := h.handler(session.Context(), message); err != nil {
				logger.Error("Failed to handle message",
					zap.String("topic", msg.Topic),
					zap.Error(err))
				continue
			}

			// Mark message as processed
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
} 