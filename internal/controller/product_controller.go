package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProductController handles product-related HTTP requests
type ProductController struct {
	logger *zap.Logger
}

// NewProductController creates a new product controller
func NewProductController(logger *zap.Logger) *ProductController {
	return &ProductController{
		logger: logger,
	}
}

// GetProducts handles GET /products request
func (c *ProductController) GetProducts(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"products": []string{"product1", "product2"},
	})
}

// CreateProduct handles POST /products request
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	ctx.JSON(201, gin.H{
		"message": "Product created successfully",
	})
}

// GetProductByID handles GET /products/:id request
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"id":    id,
		"name":  "Sample Product",
		"price": 99.99,
	})
}

// UpdateProduct handles PUT /products/:id request
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": "Product updated successfully",
		"id":      id,
	})
}

// DeleteProduct handles DELETE /products/:id request
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(200, gin.H{
		"message": "Product deleted successfully",
		"id":      id,
	})
} 