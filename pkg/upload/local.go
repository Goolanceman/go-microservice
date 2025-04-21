package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
func (u *LocalUploader) Upload(ctx context.Context, relativePath string, content io.Reader, contentType string) (string, error) {
	fullPath := filepath.Join(u.config.BaseDir, relativePath)

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
	return u.GetURL(ctx, relativePath)
}

// Download implements the Uploader interface
func (u *LocalUploader) Download(ctx context.Context, relativePath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(u.config.BaseDir, relativePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return file, nil
}

// Delete implements the Uploader interface
func (u *LocalUploader) Delete(ctx context.Context, relativePath string) error {
	fullPath := filepath.Join(u.config.BaseDir, relativePath)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}

// GetURL implements the Uploader interface
func (u *LocalUploader) GetURL(ctx context.Context, relativePath string) (string, error) {
	if u.config.BaseURL == "" {
		return "", fmt.Errorf("base URL not configured")
	}

	baseURL := u.config.BaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// Normalize slashes for URLs
	urlPath := strings.ReplaceAll(relativePath, string(os.PathSeparator), "/")
	return baseURL + urlPath, nil
}

// Close implements the Uploader interface
func (u *LocalUploader) Close() error {
	return nil // No cleanup needed for local storage
}
