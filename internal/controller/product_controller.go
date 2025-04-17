package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-microservice/internal/service"
)

// ProductController handles product-related HTTP requests
type ProductController struct {
	service service.Service
}

// NewProductController creates a new product controller
func NewProductController() *ProductController {
	return &ProductController{}
}

// GetProducts handles GET /products request
func (c *ProductController) GetProducts(ctx *gin.Context) {
	// TODO: Implement get products logic using service
	ctx.JSON(http.StatusOK, gin.H{
		"products": []string{"product1", "product2"},
	})
}

// CreateProduct handles POST /products request
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	// TODO: Implement create product logic using service
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
	})
}

// GetProduct handles GET /products/:id request
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO: Implement get product logic using service
	ctx.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  "Sample Product",
		"price": 99.99,
	})
}

// UpdateProduct handles PUT /products/:id request
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO: Implement update product logic using service
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"id":      id,
	})
}

// DeleteProduct handles DELETE /products/:id request
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO: Implement delete product logic using service
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"id":      id,
	})
}
