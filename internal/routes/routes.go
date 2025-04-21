package routes

import (
	"github.com/gin-gonic/gin"

	v1 "go-microservice/internal/routes/v1"
	v2 "go-microservice/internal/routes/v2"
	"go-microservice/internal/types"
)

// RouteGroup represents a group of routes
type RouteGroup struct {
	Prefix   string
	Handlers []gin.HandlerFunc
	Routes   []types.Route
}

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.Engine) {
	// Register v1 routes
	v1Group := router.Group("/api/v1")
	for _, route := range v1.RegisterRoutes() {
		v1Group.Handle(route.Method, route.Path, route.Handler)
	}

	// Register v2 routes
	v2Group := router.Group("/api/v2")
	for _, route := range v2.RegisterRoutes() {
		v2Group.Handle(route.Method, route.Path, route.Handler)
	}
}
