package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"cloud.google.com/go/storage"
	"go-microservice/internal/config"
	"google.golang.org/api/option"
)

// GCSUploader implements the Uploader interface for Google Cloud Storage
type GCSUploader struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func init() {
	RegisterUploader("gcs", NewGCSUploader)
}

// NewGCSUploader creates a new GCS uploader instance
func NewGCSUploader(cfg interface{}) (Uploader, error) {
	gcsCfg, ok := cfg.(*config.GCSConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// Initialize GCS client
	ctx := context.Background()
	var client *storage.Client
	var err error

	if gcsCfg.CredentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(gcsCfg.CredentialsFile))
	} else {
		client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	bucket := client.Bucket(gcsCfg.Bucket)

	// Check if bucket exists
	if _, err := bucket.Attrs(ctx); err != nil {
		return nil, fmt.Errorf("bucket does not exist or access denied: %w", err)
	}

	return &GCSUploader{
		client: client,
		bucket: bucket,
	}, nil
}

// Upload implements the Uploader interface
func (u *GCSUploader) Upload(ctx context.Context, filepath string, content io.Reader, contentType string) (string, error) {
	obj := u.bucket.Object(filepath)
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, content); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Generate signed URL
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(24 * 7 * time.Hour),
	}

	url, err := obj.SignedURL(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// Download implements the Uploader interface
func (u *GCSUploader) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	obj := u.bucket.Object(filepath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return reader, nil
}

// Delete implements the Uploader interface
func (u *GCSUploader) Delete(ctx context.Context, filepath string) error {
	obj := u.bucket.Object(filepath)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}

// GetURL implements the Uploader interface
func (u *GCSUploader) GetURL(ctx context.Context, filepath string) (string, error) {
	obj := u.bucket.Object(filepath)
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(24 * 7 * time.Hour),
	}

	url, err := obj.SignedURL(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return url, nil
} 