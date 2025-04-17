package types

import "github.com/gin-gonic/gin"

// Route represents a single API route
type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
} 