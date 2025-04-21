package middleware

import (
	"time"

	"go-microservice/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggingMiddleware adds logging to all requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		c.Next()

		latencyMs := time.Since(start).Milliseconds()

		var errors []string
		for _, e := range c.Errors {
			errors = append(errors, e.Error())
		}

		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		logger.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
			zap.Any("headers", headers),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("latency_ms", latencyMs),
			zap.Strings("errors", errors),
		)
	}
}
