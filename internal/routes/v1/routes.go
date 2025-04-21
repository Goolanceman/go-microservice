package v1

import (
	"go-microservice/internal/types"
)

// Declare the shared registry slice.
var routeRegistry []func() []types.Route

// RegisterRoutes collects all registered v1 routes.
func RegisterRoutes() []types.Route {
	var routes []types.Route
	for _, regFunc := range routeRegistry {
		routes = append(routes, regFunc()...)
	}
	return routes
}
