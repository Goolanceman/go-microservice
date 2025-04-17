package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/routes/v1"
	"github.com/mansoor/go-microservice/internal/routes/v2"
)

// RouteGroup represents a group of routes
type RouteGroup struct {
	Prefix    string
	Handlers  []gin.HandlerFunc
	Routes    []Route
}

// Route represents a single route
type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

// RegisterRoutes registers all route groups
func RegisterRoutes(logger *zap.Logger) []RouteGroup {
	// Initialize route groups
	routeGroups := make([]RouteGroup, 0)

	// Register v1 routes
	v1Routes := make([]Route, 0)
	v1Routes = append(v1Routes, v1.RegisterUserRoutes(logger)...)
	v1Routes = append(v1Routes, v1.RegisterProductRoutes(logger)...)

	routeGroups = append(routeGroups, RouteGroup{
		Prefix: "/api/v1",
		Routes: v1Routes,
	})

	// Register v2 routes
	v2Routes := make([]Route, 0)
	v2Routes = append(v2Routes, v2.RegisterTaskRoutes(logger)...)

	routeGroups = append(routeGroups, RouteGroup{
		Prefix: "/api/v2",
		Routes: v2Routes,
	})

	return routeGroups
}