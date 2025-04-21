package v1

import (
	"go-microservice/internal/controller"
	"go-microservice/internal/types"
)

// Automatically register when this file is imported.
func init() {
	routeRegistry = append(routeRegistry, RegisterUserRoutes)
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes() []types.Route {
	userController := controller.NewUserController()

	return []types.Route{
		{
			Method:  "GET",
			Path:    "/users",
			Handler: userController.GetUsers,
		},
		{
			Method:  "POST",
			Path:    "/users",
			Handler: userController.CreateUser,
		},
		{
			Method:  "GET",
			Path:    "/users/:id",
			Handler: userController.GetUser,
		},
		{
			Method:  "PUT",
			Path:    "/users/:id",
			Handler: userController.UpdateUser,
		},
		{
			Method:  "DELETE",
			Path:    "/users/:id",
			Handler: userController.DeleteUser,
		},
	}
}
