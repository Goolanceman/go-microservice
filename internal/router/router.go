package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goolanceman/go-microservice/internal/controller"
	"go.uber.org/zap"
	"time"
)

// Router handles HTTP routing
type Router struct {
	router *gin.Engine
	logger *zap.Logger
}

// NewRouter creates a new router instance
func NewRouter(messagingController *controller.MessagingController, logger *zap.Logger) *Router {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create new Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(logger))

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Messaging routes
		messaging := v1.Group("/messaging")
		{
			messaging.POST("/publish", messagingController.PublishMessage)
			messaging.GET("/topics", messagingController.GetTopics)
		}

		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.String(200, "OK")
		})
	}

	return &Router{
		router: r,
		logger: logger,
	}
}

// GetHandler returns the HTTP handler
func (r *Router) GetHandler() *gin.Engine {
	return r.router
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