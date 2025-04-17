package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-microservice/internal/models"
	"go-microservice/internal/service"
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
		c.logger.Error("Failed to open uploaded file", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	defer src.Close()

	// Read file content
	content := make([]byte, file.Size)
	if _, err := src.Read(content); err != nil {
		c.logger.Error("Failed to read file content", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Create file model
	fileModel := &models.File{
		Name:        file.Filename,
		Size:        file.Size,
		ContentType: file.Header.Get("Content-Type"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Upload file using service
	if err := c.service.UploadFile(ctx, fileModel); err != nil {
		c.logger.Error("Failed to upload file", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded successfully",
		"file":    fileModel,
	})
}

// DownloadFile handles GET /files/:id request
func (c *FileController) DownloadFile(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	file, err := c.service.DownloadFile(ctx, id)
	if err != nil {
		c.logger.Error("Failed to download file", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DeleteFile handles DELETE /files/:id request
func (c *FileController) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	if err := c.service.DeleteFile(ctx, id); err != nil {
		c.logger.Error("Failed to delete file", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
} 