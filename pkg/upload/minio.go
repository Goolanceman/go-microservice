package upload

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/goolanceman/go-microservice/internal/config"
)

// MinioUploader implements the Uploader interface for MinIO
type MinioUploader struct {
	client *minio.Client
	bucket string
}

func init() {
	RegisterUploader("minio", NewMinioUploader)
}

// NewMinioUploader creates a new MinIO uploader instance
func NewMinioUploader(cfg interface{}) (Uploader, error) {
	minioCfg, ok := cfg.(*config.MinioConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// Initialize minio client
	client, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKeyID, minioCfg.SecretAccessKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	// Check if bucket exists, create if not
	exists, err := client.BucketExists(context.Background(), minioCfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), minioCfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioUploader{
		client: client,
		bucket: minioCfg.Bucket,
	}, nil
}

// Upload implements the Uploader interface
func (u *MinioUploader) Upload(ctx context.Context, filepath string, content io.Reader, contentType string) (string, error) {
	_, err := u.client.PutObject(ctx, u.bucket, filepath, content, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Generate presigned URL
	url, err := u.client.PresignedGetObject(ctx, u.bucket, filepath, time.Hour*24*7, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// Download implements the Uploader interface
func (u *MinioUploader) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	object, err := u.client.GetObject(ctx, u.bucket, filepath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return object, nil
}

// Delete implements the Uploader interface
func (u *MinioUploader) Delete(ctx context.Context, filepath string) error {
	if err := u.client.RemoveObject(ctx, u.bucket, filepath, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}

// GetURL implements the Uploader interface
func (u *MinioUploader) GetURL(ctx context.Context, filepath string) (string, error) {
	url, err := u.client.PresignedGetObject(ctx, u.bucket, filepath, time.Hour*24*7, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
} 