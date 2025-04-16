package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goolanceman/go-microservice/internal/config"
	"github.com/goolanceman/go-microservice/internal/controller"
	"github.com/goolanceman/go-microservice/internal/router"
	"github.com/goolanceman/go-microservice/pkg/messaging"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig("config/dev.json")
	if err != nil {
		logger.Fatal("Failed to load configuration",
			zap.Error(err))
	}

	// Initialize Kafka client
	kafkaClient, err := messaging.NewKafka(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka client",
			zap.Error(err))
	}
	defer kafkaClient.Close()

	// Initialize messaging service
	messagingService := service.NewMessagingService(kafkaClient, logger, cfg)

	// Initialize messaging controller
	messagingController := controller.NewMessagingController(messagingService, logger)

	// Initialize router
	router := router.NewRouter(messagingController, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.GetHandler(),
	}

	// Start Kafka consumers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := messagingService.StartConsumers(ctx); err != nil {
		logger.Fatal("Failed to start Kafka consumers",
			zap.Error(err))
	}

	// Start HTTP server in a goroutine
	go func() {
		logger.Info("Starting server",
			zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server",
				zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown",
			zap.Error(err))
	}

	logger.Info("Server exiting")
} 