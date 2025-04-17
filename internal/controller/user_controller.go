package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/model"
	"github.com/mansoor/go-microservice/internal/service"
)

// UserController handles user-related HTTP requests
type UserController struct {
	logger  *zap.Logger
	service service.Service
}

// NewUserController creates a new user controller
func NewUserController(logger *zap.Logger, service service.Service) *UserController {
	return &UserController{
		logger:  logger,
		service: service,
	}
}

// CreateUser handles POST /users request
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := c.service.CreateUser(ctx, user)
	if err != nil {
		c.logger.Error("Failed to create user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

// GetUser handles GET /users/:id request
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetUser(ctx, id)
	if err != nil {
		c.logger.Error("Failed to get user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id request
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := c.service.UpdateUser(ctx, id, user)
	if err != nil {
		c.logger.Error("Failed to update user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE /users/:id request
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteUser(ctx, id); err != nil {
		c.logger.Error("Failed to delete user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.Status(http.StatusNoContent)
} 