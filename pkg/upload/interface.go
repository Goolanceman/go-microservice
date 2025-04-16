package upload

import (
	"context"
	"io"
)

// Uploader defines the interface for file upload operations
type Uploader interface {
	// Upload uploads a file to the storage backend
	Upload(ctx context.Context, path string, content io.Reader, contentType string) (string, error)
	
	// Download retrieves a file from the storage backend
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	
	// Delete removes a file from the storage backend
	Delete(ctx context.Context, path string) error
	
	// GetURL returns the public URL for a file (if supported by the backend)
	GetURL(ctx context.Context, path string) (string, error)
}

// Factory function type for creating new uploaders
type UploaderFactory func(cfg interface{}) (Uploader, error)

// Registry of available uploader implementations
var uploaderFactories = make(map[string]UploaderFactory)

// RegisterUploader registers a new uploader implementation
func RegisterUploader(name string, factory UploaderFactory) {
	uploaderFactories[name] = factory
}

// NewUploader creates a new uploader instance based on the specified backend
func NewUploader(backend string, cfg interface{}) (Uploader, error) {
	factory, exists := uploaderFactories[backend]
	if !exists {
		return nil, ErrUnsupportedBackend
	}
	return factory(cfg)
}

// Error types
var (
	ErrUnsupportedBackend = NewUploadError("unsupported upload backend")
	ErrInvalidConfig     = NewUploadError("invalid configuration")
	ErrUploadFailed     = NewUploadError("upload failed")
	ErrDownloadFailed   = NewUploadError("download failed")
	ErrDeleteFailed     = NewUploadError("delete failed")
)

// UploadError represents an error that occurred during upload operations
type UploadError struct {
	Message string
}

func (e *UploadError) Error() string {
	return e.Message
}

// NewUploadError creates a new UploadError
func NewUploadError(message string) *UploadError {
	return &UploadError{Message: message}
} 