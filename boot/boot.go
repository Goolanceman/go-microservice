package boot

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

type App struct {
	Config *config.Config
	Server *http.Server
}

func Start(ctx context.Context) {
	app := &App{}

	// Initialize configuration
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	app.Config = cfg

	// Initialize logger
	if err := initLogger(cfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer logger.Sync()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	logger.Info("Configuration loaded successfully",
		zap.String("env", env),
	)

	logger.Info("Logger initialized successfully")

	// Initialize HTTP server
	app.Server = initHTTPServer()
	logger.Info("HTTP server initialized successfully")

	// Register routes before starting the server
	if err := app.Server.RegisterRoutes(); err != nil {
		logger.Fatal("Failed to register routes", zap.Error(err))
	}

	logger.Info("Application routes registered successfully")
	logger.Info("Application boot sequence completed. Server is running.")

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := app.Server.Start(ctx); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	<-sigChan
	logger.Info("Shutdown signal received. Shutting down gracefully...")
}

func initConfig() (*config.Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	return config.LoadConfig("config/" + env + ".json")
}

func initLogger(cfg *config.Config) error {
	return logger.Init(cfg.Server.LogLevel, cfg.Server.LogFile)
}

func initHTTPServer() *http.Server {
	return http.NewServer()
}
