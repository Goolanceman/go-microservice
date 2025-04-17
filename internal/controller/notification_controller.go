package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NotificationController handles notification-related HTTP requests
type NotificationController struct {
	logger *zap.Logger
}

// NewNotificationController creates a new notification controller
func NewNotificationController(logger *zap.Logger) *NotificationController {
	return &NotificationController{
		logger: logger,
	}
}

// GetNotifications handles GET /notifications request
func (c *NotificationController) GetNotifications(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"notifications": []string{"notification1", "notification2"},
	})
}

// CreateNotification handles POST /notifications request
func (c *NotificationController) CreateNotification(ctx *gin.Context) {
	ctx.JSON(201, gin.H{
		"message": "Notification created successfully",
	})
}

// GetNotificationByID handles GET /notifications/:id request
func (c *NotificationController) GetNotificationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"id":      id,
		"type":    "email",
		"content": "Your order has been shipped",
		"status":  "sent",
	})
}

// UpdateNotification handles PUT /notifications/:id request
func (c *NotificationController) UpdateNotification(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": "Notification updated successfully",
		"id":      id,
	})
}

// DeleteNotification handles DELETE /notifications/:id request
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": "Notification deleted successfully",
		"id":      id,
	})
} 