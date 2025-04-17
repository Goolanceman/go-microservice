package v2

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/controller"
)

// RegisterTaskRoutes registers all task-related routes
func RegisterTaskRoutes(logger *zap.Logger) []Route {
	taskController := controller.NewTaskController(logger, nil) // TODO: Pass actual service
	return []Route{
		{
			Method:  "GET",
			Path:    "/tasks",
			Handler: taskController.handleGetTasks,
		},
		{
			Method:  "POST",
			Path:    "/tasks",
			Handler: taskController.handlePostTask,
		},
		{
			Method:  "GET",
			Path:    "/tasks/:id",
			Handler: taskController.handleGetTaskByID,
		},
		{
			Method:  "PUT",
			Path:    "/tasks/:id",
			Handler: taskController.handlePutTask,
		},
		{
			Method:  "DELETE",
			Path:    "/tasks/:id",
			Handler: taskController.handleDeleteTask,
		},
	}
}