package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-microservice/internal/models"
	"go-microservice/internal/service"
	"go-microservice/pkg/logger"
)

// UserController handles user-related HTTP requests
type UserController struct {
	service service.Service
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{}
}

// GetUsers handles GET /users request
func (c *UserController) GetUsers(ctx *gin.Context) {
	// TODO: Implement get users logic using service
	ctx.JSON(http.StatusOK, gin.H{
		"users": []string{"user1", "user2"},
	})
}

// CreateUser handles POST /users request
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateUser(ctx, &user); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// GetUser handles GET /users/:id request
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := c.service.GetUser(ctx, id)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err))
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
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := c.service.UpdateUser(ctx, &user); err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id request
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if err := c.service.DeleteUser(ctx, id); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
