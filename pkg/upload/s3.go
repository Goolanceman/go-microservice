package upload

import (
	"context"
	"fmt"
	"io"

	"go-microservice/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Uploader implements the Uploader interface for AWS S3
type S3Uploader struct {
	client *s3.Client
	bucket string
	region string
}

func init() {
	RegisterUploader("s3", NewS3Uploader)
}

// NewS3Uploader creates a new S3 uploader instance
func NewS3Uploader(cfg interface{}) (Uploader, error) {
	s3Cfg, ok := cfg.(*config.S3Config)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// Create S3 client
	client := s3.New(s3.Options{
		Region: s3Cfg.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			s3Cfg.AccessKeyID,
			s3Cfg.SecretAccessKey,
			"",
		)),
	})

	return &S3Uploader{
		client: client,
		bucket: s3Cfg.Bucket,
		region: s3Cfg.Region,
	}, nil
}

// Upload implements the Uploader interface
func (u *S3Uploader) Upload(ctx context.Context, filepath string, content io.Reader, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      &u.bucket,
		Key:         &filepath,
		Body:        content,
		ContentType: &contentType,
	}

	_, err := u.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, filepath)
	return url, nil
}

// Download implements the Uploader interface
func (u *S3Uploader) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: &u.bucket,
		Key:    &filepath,
	}

	result, err := u.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}

	return result.Body, nil
}

// Delete implements the Uploader interface
func (u *S3Uploader) Delete(ctx context.Context, filepath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &u.bucket,
		Key:    &filepath,
	}

	_, err := u.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	return nil
}

// GetURL implements the Uploader interface
func (u *S3Uploader) GetURL(ctx context.Context, filepath string) (string, error) {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, filepath), nil
}
