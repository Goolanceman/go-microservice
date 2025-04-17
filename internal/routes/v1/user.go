package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/controller"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(logger *zap.Logger) []Route {
	userController := controller.NewUserController(logger, nil) // TODO: Pass actual service
	
	return []Route{
		{
			Method:  "GET",
			Path:    "/users",
			Handler: userController.handleGetUsers,
		},
		{
			Method:  "POST",
			Path:    "/users",
			Handler: userController.handlePostUsers,
		},
		{
			Method:  "GET",
			Path:    "/users/:id",
			Handler: userController.handleGetUserByID,
		},
		{
			Method:  "PUT",
			Path:    "/users/:id",
			Handler: userController.handlePutUser,
		},
		{
			Method:  "DELETE",
			Path:    "/users/:id",
			Handler: userController.handleDeleteUser,
		},
	}
}