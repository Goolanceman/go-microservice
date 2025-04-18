package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/goolanceman/go-microservice/pkg/logger"
	"go.uber.org/zap"
)

// NATS implements the Messaging interface
type NATS struct {
	conn          *nats.Conn
	subscriptions map[string]*nats.Subscription
	mu            sync.RWMutex
	logger        *logger.Logger
}

// NewNATS creates a new NATS client
func NewNATS(url string, logger *logger.Logger) (*NATS, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &NATS{
		conn:          conn,
		subscriptions: make(map[string]*nats.Subscription),
		logger:        logger,
	}, nil
}

// SubscribeMultiple subscribes to multiple topics with the same handler
func (n *NATS) SubscribeMultiple(ctx context.Context, topics []string, handler Handler) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.conn == nil {
		return ErrNotConnected
	}

	// Create a wait group to track all subscriptions
	var wg sync.WaitGroup
	errChan := make(chan error, len(topics))

	for _, topic := range topics {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()

			// Create subscription
			sub, err := n.conn.Subscribe(t, func(msg *nats.Msg) {
				// Create a new context for each message
				msgCtx := context.Background()
				
				// Create message wrapper
				message := &Message{
					Topic:   t,
					Payload: msg.Data,
					Headers: make(map[string]string),
				}

				// Copy headers from NATS message
				for key, values := range msg.Header {
					if len(values) > 0 {
						message.Headers[key] = values[0]
					}
				}

				// Call handler
				if err := handler(msgCtx, message); err != nil {
					n.logger.Error("Error handling message",
						zap.String("topic", t),
						zap.Error(err))
				}
			})

			if err != nil {
				errChan <- fmt.Errorf("failed to subscribe to topic %s: %w", t, err)
				return
			}

			// Store subscription
			n.subscriptions[t] = sub
		}(topic)
	}

	// Wait for all subscriptions to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple subscription errors: %v", errs)
	}

	return nil
}

// Subscribe subscribes to a single topic
func (n *NATS) Subscribe(ctx context.Context, topic string, handler Handler) error {
	return n.SubscribeMultiple(ctx, []string{topic}, handler)
}

// Close closes the NATS connection and all subscriptions
func (n *NATS) Close() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Unsubscribe from all topics
	for _, sub := range n.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			n.logger.Error("Failed to unsubscribe", err)
		}
	}

	// Close connection
	n.conn.Close()
	return nil
}

// ... existing code ... 