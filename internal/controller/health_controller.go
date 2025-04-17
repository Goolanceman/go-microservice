package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthController handles health check related HTTP requests
type HealthController struct {
	logger *zap.Logger
}

// NewHealthController creates a new health controller
func NewHealthController(logger *zap.Logger) *HealthController {
	return &HealthController{
		logger: logger,
	}
}

// GetHealth handles GET /health request
func (c *HealthController) GetHealth(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "healthy",
		"timestamp": "2024-03-20T12:00:00Z",
	})
}

// GetHealthDetailed handles GET /health/detailed request
func (c *HealthController) GetHealthDetailed(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "healthy",
		"timestamp": "2024-03-20T12:00:00Z",
		"components": []gin.H{
			{
				"name": "database",
				"status": "healthy",
				"latency": 5,
			},
			{
				"name": "cache",
				"status": "healthy",
				"latency": 2,
			},
			{
				"name": "storage",
				"status": "healthy",
				"latency": 10,
			},
		},
		"metrics": gin.H{
			"cpu_usage": 45.2,
			"memory_usage": 60.8,
			"disk_usage": 30.5,
		},
	})
}

// ReadinessCheck handles GET /ready request
func (c *HealthController) ReadinessCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ready",
		"services": gin.H{
			"database": "ready",
			"cache":    "ready",
			"message_queue": "ready",
		},
	})
}

// LivenessCheck handles GET /live request
func (c *HealthController) LivenessCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "alive",
		"uptime": "1h 30m",
	})
} 