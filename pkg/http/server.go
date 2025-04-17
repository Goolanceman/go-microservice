package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-microservice/internal/routes"
	"go-microservice/internal/service"
	"go-microservice/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	engine  *gin.Engine
	config  *Config
	logger  *zap.Logger
	service service.Service
}

// Config holds HTTP server configuration
type Config struct {
	Port         string
	Environment  string
	AllowOrigins string
	RoutesPath   string // Path to routes directory
}

// NewServer creates a new HTTP server instance
func NewServer(cfg *Config) *Server {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create new Gin engine
	engine := gin.New()

	// Add default middleware
	engine.Use(gin.Recovery())
	engine.Use(loggingMiddleware())

	// Add CORS middleware if origins are specified
	if cfg.AllowOrigins != "" {
		engine.Use(corsMiddleware(cfg.AllowOrigins))
	}

	return &Server{
		engine: engine,
		config: cfg,
	}
}

// RegisterRoutes registers routes from the specified path
func (s *Server) RegisterRoutes() error {
	// Register routes using the engine, logger, and service
	routes.RegisterRoutes(s.engine)
	return nil
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	// Register routes
	if err := s.RegisterRoutes(); err != nil {
		return fmt.Errorf("failed to register routes: %w", err)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.engine,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server",
			zap.String("port", s.config.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server",
				zap.Error(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown server
	logger.Info("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}

// loggingMiddleware adds logging to all requests
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(allowOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
