package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// Error sends a standardized error JSON response
func Error(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	})
}
