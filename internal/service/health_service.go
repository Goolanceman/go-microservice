package service

import (
	"context"
	"sync"
	"time"

	"github.com/goolanceman/go-microservice/internal/config"
	"github.com/goolanceman/go-microservice/pkg/database"
	"github.com/goolanceman/go-microservice/pkg/kafka"
)

// HealthService handles health checks for the service
type HealthService struct {
	config        *config.Config
	redisClient   *database.RedisClient
	kafkaProducer *kafka.Producer
	kafkaConsumer *kafka.Consumer
	mu            sync.RWMutex
	ready         bool
}

// NewHealthService creates a new health service
func NewHealthService(
	cfg *config.Config,
	redisClient *database.RedisClient,
	kafkaProducer *kafka.Producer,
	kafkaConsumer *kafka.Consumer,
) *HealthService {
	return &HealthService{
		config:        cfg,
		redisClient:   redisClient,
		kafkaProducer: kafkaProducer,
		kafkaConsumer: kafkaConsumer,
		ready:         false,
	}
}

// SetReady sets the service readiness status
func (s *HealthService) SetReady(ready bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = ready
}

// CheckLiveness performs a basic liveness check
func (s *HealthService) CheckLiveness() map[string]string {
	return map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
}

// CheckReadiness performs a comprehensive readiness check
func (s *HealthService) CheckReadiness() map[string]string {
	status := map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}

	// Check Redis if enabled
	if s.config.Features.EnableRedis {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.redisClient.HealthCheck(ctx); err != nil {
			status["status"] = "error"
			status["redis"] = "unhealthy"
		} else {
			status["redis"] = "healthy"
		}
	}

	// Check Kafka if enabled
	if s.config.Features.EnableKafka {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.kafkaProducer.HealthCheck(ctx); err != nil {
			status["status"] = "error"
			status["kafka_producer"] = "unhealthy"
		} else {
			status["kafka_producer"] = "healthy"
		}

		if err := s.kafkaConsumer.HealthCheck(ctx); err != nil {
			status["status"] = "error"
			status["kafka_consumer"] = "unhealthy"
		} else {
			status["kafka_consumer"] = "healthy"
		}
	}

	// Check service readiness flag
	s.mu.RLock()
	if !s.ready {
		status["status"] = "error"
		status["service"] = "not ready"
	} else {
		status["service"] = "ready"
	}
	s.mu.RUnlock()

	return status
} 