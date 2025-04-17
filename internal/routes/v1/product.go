package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/controller"
)

// RegisterProductRoutes registers all product-related routes
func RegisterProductRoutes(logger *zap.Logger) []Route {
	productController := controller.NewProductController(logger, nil) // TODO: Pass actual service

	return []Route{
		{
			Method:  "GET",
			Path:    "/products",
			Handler: productController.handleGetProducts,
		},
		{
			Method:  "POST",
			Path:    "/products",
			Handler: productController.handlePostProduct,
		},
		{
			Method:  "GET",
			Path:    "/products/:id",
			Handler: productController.handleGetProductByID,
		},
		{
			Method:  "PUT",
			Path:    "/products/:id",
			Handler: productController.handlePutProduct,
		},
		{
			Method:  "DELETE",
			Path:    "/products/:id",
			Handler: productController.handleDeleteProduct,
		},
	}
}