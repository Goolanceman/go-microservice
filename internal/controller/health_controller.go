package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goolanceman/go-microservice/internal/service"
	"github.com/goolanceman/go-microservice/pkg/logger"
)

// HealthController handles health check endpoints
type HealthController struct {
	healthService *service.HealthService
}

// NewHealthController creates a new health controller
func NewHealthController(healthService *service.HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}

// RegisterRoutes registers the health check routes
func (c *HealthController) RegisterRoutes(router *gin.Engine) {
	router.GET("/healthz", c.Liveness)
	router.GET("/readyz", c.Readiness)
}

// Liveness handles the liveness probe
// @Summary Liveness probe
// @Description Check if the service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (c *HealthController) Liveness(ctx *gin.Context) {
	status := c.healthService.CheckLiveness()
	ctx.JSON(http.StatusOK, status)
}

// Readiness handles the readiness probe
// @Summary Readiness probe
// @Description Check if the service is ready to handle requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /readyz [get]
func (c *HealthController) Readiness(ctx *gin.Context) {
	status := c.healthService.CheckReadiness()
	if status["status"] == "ok" {
		ctx.JSON(http.StatusOK, status)
	} else {
		ctx.JSON(http.StatusServiceUnavailable, status)
	}
} 