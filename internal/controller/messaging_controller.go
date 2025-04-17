package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"go-microservice/internal/service"
	"go.uber.org/zap"
)

// MessagingController handles HTTP requests for messaging operations
type MessagingController struct {
	service *service.MessagingService
	logger  *zap.Logger
}

// NewMessagingController creates a new messaging controller
func NewMessagingController(service *service.MessagingService, logger *zap.Logger) *MessagingController {
	return &MessagingController{
		service: service,
		logger:  logger,
	}
}

// PublishMessage handles HTTP requests to publish messages
func (c *MessagingController) PublishMessage(ctx *gin.Context) {
	var req struct {
		Topic   string          `json:"topic" binding:"required"`
		Payload json.RawMessage `json:"payload" binding:"required"`
	}

	// Bind and validate request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("Failed to bind request",
			zap.Error(err))
		ctx.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Publish message
	if err := c.service.PublishMessage(ctx.Request.Context(), req.Topic, req.Payload); err != nil {
		c.logger.Error("Failed to publish message",
			zap.Error(err),
			zap.String("topic", req.Topic))
		ctx.JSON(500, gin.H{
			"error": "Failed to publish message",
		})
		return
	}

	// Return success response
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Message published successfully",
	})
}

// GetTopics handles HTTP requests to get available topics
func (c *MessagingController) GetTopics(ctx *gin.Context) {
	// In a real implementation, you might want to fetch this from configuration or Kafka
	topics := []string{"topic1", "topic2", "topic3"}

	ctx.JSON(200, gin.H{
		"topics": topics,
	})
} 