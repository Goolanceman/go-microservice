package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MetricsController handles metrics related HTTP requests
type MetricsController struct {
	logger *zap.Logger
}

// NewMetricsController creates a new metrics controller
func NewMetricsController(logger *zap.Logger) *MetricsController {
	return &MetricsController{
		logger: logger,
	}
}

// GetMetrics handles GET /metrics request
func (c *MetricsController) GetMetrics(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"metrics": gin.H{
			"requests": gin.H{
				"total": 1000,
				"successful": 980,
				"failed": 20,
				"rate": 50.5,
			},
			"latency": gin.H{
				"p50": 100,
				"p90": 200,
				"p99": 500,
			},
			"system": gin.H{
				"cpu_usage": 45.2,
				"memory_usage": 60.8,
				"disk_usage": 30.5,
			},
		},
	})
}

// GetMetricsDetailed handles GET /metrics/detailed request
func (c *MetricsController) GetMetricsDetailed(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"timestamp": "2024-03-20T12:00:00Z",
		"metrics": []gin.H{
			{
				"name": "http_requests_total",
				"type": "counter",
				"value": 1000,
				"labels": gin.H{
					"method": "GET",
					"path": "/api/v1/users",
				},
			},
			{
				"name": "http_request_duration_seconds",
				"type": "histogram",
				"buckets": []float64{0.1, 0.5, 1.0, 2.5, 5.0},
				"count": 1000,
				"sum": 500.0,
			},
			{
				"name": "system_cpu_usage",
				"type": "gauge",
				"value": 45.2,
				"unit": "percent",
			},
		},
	})
}

// GetMetricsByType handles GET /metrics/:type request
func (c *MetricsController) GetMetricsByType(ctx *gin.Context) {
	metricType := ctx.Param("type")
	
	// Example metrics by type
	metrics := map[string]gin.H{
		"requests": {
			"total": 1000,
			"success": 950,
			"error": 50,
		},
		"response_time": {
			"avg": 150,
			"p95": 200,
			"p99": 300,
		},
		"system": {
			"cpu_usage": 45.5,
			"memory_usage": 60.2,
			"disk_usage": 30.8,
		},
	}

	if data, exists := metrics[metricType]; exists {
		ctx.JSON(200, gin.H{
			"type": metricType,
			"data": data,
		})
	} else {
		ctx.JSON(404, gin.H{
			"error": "Metric type not found",
		})
	}
} 