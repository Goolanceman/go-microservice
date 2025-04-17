package v1

import (
	"go-microservice/internal/types"
)

// RegisterRoutes registers all v1 routes
func RegisterRoutes() []types.Route {
	var routes []types.Route

	// Register user routes
	routes = append(routes, RegisterUserRoutes()...)

	// Register product routes
	routes = append(routes, RegisterProductRoutes()...)

	return routes
}
