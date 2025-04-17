package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TaskController handles task-related HTTP requests
type TaskController struct {
	logger *zap.Logger
}

// NewTaskController creates a new task controller
func NewTaskController(logger *zap.Logger) *TaskController {
	return &TaskController{
		logger: logger,
	}
}

// GetTasks handles GET /tasks request
func (c *TaskController) GetTasks(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"tasks": []string{"task1", "task2"},
	})
}

// CreateTask handles POST /tasks request
func (c *TaskController) CreateTask(ctx *gin.Context) {
	ctx.JSON(201, gin.H{
		"message": "Task created successfully",
	})
}

// GetTaskByID handles GET /tasks/:id request
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"id":          id,
		"title":       "Sample Task",
		"description": "This is a sample task",
		"status":      "pending",
	})
}

// UpdateTask handles PUT /tasks/:id request
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": "Task updated successfully",
		"id":      id,
	})
}

// DeleteTask handles DELETE /tasks/:id request
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": "Task deleted successfully",
		"id":      id,
	})
} 