package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goolanceman/go-microservice/internal/config"
)

// LocalConfig holds local storage settings
type LocalConfig struct {
	BaseDir    string `mapstructure:"base_dir"`
	BaseURL    string `mapstructure:"base_url"`
	CreateDirs bool   `mapstructure:"create_dirs"`
}

// LocalUploader implements the Uploader interface for local storage
type LocalUploader struct {
	config *LocalConfig
}

func init() {
	RegisterUploader("local", NewLocalUploader)
}

// NewLocalUploader creates a new local storage uploader instance
func NewLocalUploader(cfg interface{}) (Uploader, error) {
	localCfg, ok := cfg.(*LocalConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// Ensure base directory exists
	if localCfg.CreateDirs {
		if err := os.MkdirAll(localCfg.BaseDir, 0755); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
		}
	}

	return &LocalUploader{
		config: localCfg,
	}, nil
}

// Upload implements the Uploader interface
func (u *LocalUploader) Upload(ctx context.Context, filepath string, content io.Reader, contentType string) (string, error) {
	// Create full path
	fullPath := filepath.Join(u.config.BaseDir, filepath)

	// Create directories if needed
	if u.config.CreateDirs {
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Return URL
	return u.GetURL(ctx, filepath)
}

// Download implements the Uploader interface
func (u *LocalUploader) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(u.config.BaseDir, filepath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return file, nil
}

// Delete implements the Uploader interface
func (u *LocalUploader) Delete(ctx context.Context, filepath string) error {
	fullPath := filepath.Join(u.config.BaseDir, filepath)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}

// GetURL implements the Uploader interface
func (u *LocalUploader) GetURL(ctx context.Context, filepath string) (string, error) {
	if u.config.BaseURL == "" {
		return "", fmt.Errorf("base URL not configured")
	}

	// Ensure base URL ends with a slash
	baseURL := u.config.BaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return baseURL + filepath, nil
}

// Close implements the Uploader interface
func (u *LocalUploader) Close() error {
	return nil // No cleanup needed for local storage
} 