package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/goolanceman/go-microservice/pkg/http"
	"github.com/goolanceman/go-microservice/internal/config"
	"github.com/goolanceman/go-microservice/pkg/logger"
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
	server := http.NewServer(cfg.Server)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	if err := server.Start(ctx); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	// Wait for interrupt signal
	<-sigChan
	logger.Info("Shutting down gracefully...")
	cancel()
} 