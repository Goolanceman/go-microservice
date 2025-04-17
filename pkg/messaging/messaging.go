package messaging

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrNotConnected = errors.New("not connected to messaging system")
)

// Message represents a message in the messaging system
type Message struct {
	Topic   string
	Payload []byte
	Headers map[string]string
}

// Handler is a function that handles incoming messages
type Handler func(context.Context, *Message) error

// Messaging defines the interface for messaging systems
type Messaging interface {
	// Subscribe subscribes to a single topic
	Subscribe(ctx context.Context, topic string, handler Handler) error
	
	// SubscribeMultiple subscribes to multiple topics with the same handler
	SubscribeMultiple(ctx context.Context, topics []string, handler Handler) error
	
	// Publish publishes a message to a topic
	Publish(ctx context.Context, msg *Message) error
	
	// Close closes the messaging connection
	Close() error
} 