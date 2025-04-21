package v2

import (
	"go-microservice/internal/controller"
	"go-microservice/internal/types"
)

// Automatically register when this file is imported.
func init() {
	routeRegistry = append(routeRegistry, RegisterTaskRoutes)
}

// RegisterTaskRoutes registers all task-related routes
func RegisterTaskRoutes() []types.Route {
	taskController := controller.NewTaskController()
	return []types.Route{
		{
			Method:  "GET",
			Path:    "/tasks",
			Handler: taskController.GetTasks,
		},
		{
			Method:  "POST",
			Path:    "/tasks",
			Handler: taskController.CreateTask,
		},
		{
			Method:  "GET",
			Path:    "/tasks/:id",
			Handler: taskController.GetTask,
		},
		{
			Method:  "PUT",
			Path:    "/tasks/:id",
			Handler: taskController.UpdateTask,
		},
		{
			Method:  "DELETE",
			Path:    "/tasks/:id",
			Handler: taskController.DeleteTask,
		},
	}
}
