package v1

import (
	"go-microservice/internal/controller"
	"go-microservice/internal/types"
)

// RegisterProductRoutes registers all product-related routes
func RegisterProductRoutes() []types.Route {
	productController := controller.NewProductController()

	return []types.Route{
		{
			Method:  "GET",
			Path:    "/products",
			Handler: productController.GetProducts,
		},
		{
			Method:  "POST",
			Path:    "/products",
			Handler: productController.CreateProduct,
		},
		{
			Method:  "GET",
			Path:    "/products/:id",
			Handler: productController.GetProduct,
		},
		{
			Method:  "PUT",
			Path:    "/products/:id",
			Handler: productController.UpdateProduct,
		},
		{
			Method:  "DELETE",
			Path:    "/products/:id",
			Handler: productController.DeleteProduct,
		},
	}
}
