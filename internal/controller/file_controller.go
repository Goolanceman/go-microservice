package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mansoor/go-microservice/internal/service"
)

// FileController handles file-related HTTP requests
type FileController struct {
	logger  *zap.Logger
	service service.Service
}

// NewFileController creates a new file controller
func NewFileController(logger *zap.Logger, service service.Service) *FileController {
	return &FileController{
		logger:  logger,
		service: service,
	}
}

// UploadFile handles POST /files request
func (c *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.logger.Error("Failed to open uploaded file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	defer src.Close()

	// Read file content
	content := make([]byte, file.Size)
	if _, err := src.Read(content); err != nil {
		c.logger.Error("Failed to read file content", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Upload file
	url, err := c.service.UploadFile(ctx, file.Filename, content)
	if err != nil {
		c.logger.Error("Failed to upload file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"url": url})
}

// DownloadFile handles GET /files/:id request
func (c *FileController) DownloadFile(ctx *gin.Context) {
	id := ctx.Param("id")
	content, err := c.service.DownloadFile(ctx, id)
	if err != nil {
		c.logger.Error("Failed to download file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	ctx.Data(http.StatusOK, "application/octet-stream", content)
}

// DeleteFile handles DELETE /files/:id request
func (c *FileController) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteFile(ctx, id); err != nil {
		c.logger.Error("Failed to delete file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	ctx.Status(http.StatusNoContent)
} 