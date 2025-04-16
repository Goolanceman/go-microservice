package messaging

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/your-project/logger"
	"github.com/your-project/models"
	"go.uber.org/zap"
)

type Kafka struct {
	consumer *sarama.ConsumerGroup
	mu       *sync.Mutex
	logger   *logger.Logger
}

// SubscribeMultiple subscribes to multiple topics with the same handler
func (k *Kafka) SubscribeMultiple(ctx context.Context, topics []string, handler Handler) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.consumer == nil {
		return ErrNotConnected
	}

	// Subscribe to all topics
	if err := k.consumer.SubscribeTopics(topics, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	// Start consuming messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := k.consumer.ReadMessage(-1)
				if err != nil {
					k.logger.Error("Error reading message",
						zap.Error(err))
					continue
				}

				// Create message wrapper
				message := &Message{
					Topic:   *msg.TopicPartition.Topic,
					Payload: msg.Value,
					Headers: make(map[string]string),
				}

				// Copy headers from Kafka message
				for _, h := range msg.Headers {
					message.Headers[string(h.Key)] = string(h.Value)
				}

				// Call handler
				if err := handler(ctx, message); err != nil {
					k.logger.Error("Error handling message",
						zap.String("topic", *msg.TopicPartition.Topic),
						zap.Error(err))
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