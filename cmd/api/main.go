package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-microservice/internal/config"
	"go-microservice/pkg/http"
	"go-microservice/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Read environment to pick config
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Load config file based on env
	cfg, err := config.LoadConfig("config/" + env + ".json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger from config
	if err := logger.Init(cfg.Server.LogLevel, cfg.Server.LogFile); err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Create HTTP server
	httpConfig := &http.Config{
		Port:         cfg.Server.Port,
		Environment:  cfg.Server.Environment,
		AllowOrigins: cfg.Server.AllowOrigins,
	}
	server := http.NewServer(httpConfig)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.Start(ctx); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	logger.Info("Shutting down gracefully...")
	cancel()
}
