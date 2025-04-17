package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/routes"
	"github.com/mansoor/go-microservice/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	engine *gin.Engine
	config *Config
	logger *zap.Logger
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
	engine.Use(loggingMiddleware(logger))

	// Add CORS middleware if origins are specified
	if cfg.AllowOrigins != "" {
		engine.Use(corsMiddleware(cfg.AllowOrigins))
	}

	return &Server{
		engine: engine,
		config: cfg,
		logger: logger,
	}
}

// RegisterRoutes registers routes from the specified path
func (s *Server) RegisterRoutes() error {
	// Load route groups
	routeGroups := routes.RegisterRoutes(s.logger)

	// Register each route group
	for _, group := range routeGroups {
		// Create route group
		routerGroup := s.engine.Group(group.Prefix)

		// Add group-level middleware
		for _, middleware := range group.Handlers {
			routerGroup.Use(middleware)
		}

		// Register routes
		for _, route := range group.Routes {
			// Add route-level middleware
			handlers := append(route.Middlewares, route.Handler)

			// Register route based on method
			switch route.Method {
			case "GET":
				routerGroup.GET(route.Path, handlers...)
			case "POST":
				routerGroup.POST(route.Path, handlers...)
			case "PUT":
				routerGroup.PUT(route.Path, handlers...)
			case "DELETE":
				routerGroup.DELETE(route.Path, handlers...)
			case "PATCH":
				routerGroup.PATCH(route.Path, handlers...)
			default:
				return fmt.Errorf("unsupported HTTP method: %s", route.Method)
			}
		}
	}

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
		s.logger.Info("Starting HTTP server",
			zap.String("port", s.config.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server",
				zap.Error(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown server
	s.logger.Info("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}

// loggingMiddleware adds logging to all requests
func loggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
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