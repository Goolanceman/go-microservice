// middleware/gzip.go
package middleware

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func GzipMiddleware() gin.HandlerFunc {
	// Best balance for APIs: compression level 5
	return gzip.Gzip(gzip.DefaultCompression)
}
