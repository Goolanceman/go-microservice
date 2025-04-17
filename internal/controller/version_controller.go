package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VersionController handles version related HTTP requests
type VersionController struct {
	logger *zap.Logger
}

// NewVersionController creates a new version controller
func NewVersionController(logger *zap.Logger) *VersionController {
	return &VersionController{
		logger: logger,
	}
}

// GetVersion handles GET /version request
func (c *VersionController) GetVersion(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": "1.0.0",
		"build": gin.H{
			"commit": "abc123",
			"date": "2024-03-20T12:00:00Z",
		},
	})
}

// GetVersionDetailed handles GET /version/detailed request
func (c *VersionController) GetVersionDetailed(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": "1.0.0",
		"build": gin.H{
			"commit": "abc123",
			"date": "2024-03-20T12:00:00Z",
			"branch": "main",
			"go_version": "1.21.0",
		},
		"dependencies": []gin.H{
			{
				"name": "gin",
				"version": "v1.9.1",
			},
			{
				"name": "zap",
				"version": "v1.26.0",
			},
			{
				"name": "gorm",
				"version": "v1.25.5",
			},
		},
	})
}

// GetVersionInfo handles GET /version/info request
func (c *VersionController) GetVersionInfo(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"service": "go-microservice",
		"version": "1.0.0",
		"environment": "production",
		"uptime": "24h",
		"status": "healthy",
	})
} 