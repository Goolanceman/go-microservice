package v2

import (
	"go-microservice/internal/types"
)

// RegisterRoutes registers all v2 routes
func RegisterRoutes() []types.Route {
	var routes []types.Route

	// Register task routes
	routes = append(routes, RegisterTaskRoutes()...)

	return routes
} 