package upload

import (
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/jlaffaye/ftp"
	"go-microservice/internal/config"
)

// FTPConfig holds FTP connection settings
type FTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	BaseDir  string `mapstructure:"base_dir"`
}

// FTPUploader implements the Uploader interface for FTP
type FTPUploader struct {
	client  *ftp.ServerConn
	config  *FTPConfig
	baseURL string
}

func init() {
	RegisterUploader("ftp", NewFTPUploader)
}

// NewFTPUploader creates a new FTP uploader instance
func NewFTPUploader(cfg interface{}) (Uploader, error) {
	ftpCfg, ok := cfg.(*FTPConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// Connect to FTP server
	addr := fmt.Sprintf("%s:%s", ftpCfg.Host, ftpCfg.Port)
	client, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	// Login
	if err := client.Login(ftpCfg.Username, ftpCfg.Password); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	// Create base directory if it doesn't exist
	if ftpCfg.BaseDir != "" {
		if err := client.MakeDir(ftpCfg.BaseDir); err != nil {
			// Ignore error if directory already exists
		}
	}

	return &FTPUploader{
		client:  client,
		config:  ftpCfg,
		baseURL: fmt.Sprintf("ftp://%s:%s@%s:%s", ftpCfg.Username, ftpCfg.Password, ftpCfg.Host, ftpCfg.Port),
	}, nil
}

// Upload implements the Uploader interface
func (u *FTPUploader) Upload(ctx context.Context, filepath string, content io.Reader, contentType string) (string, error) {
	// Create full path
	fullPath := path.Join(u.config.BaseDir, filepath)

	// Create directories if they don't exist
	dir := path.Dir(fullPath)
	if err := u.client.MakeDir(dir); err != nil {
		// Ignore error if directory already exists
	}

	// Upload file
	if err := u.client.Stor(fullPath, content); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Return FTP URL
	return fmt.Sprintf("%s/%s", u.baseURL, fullPath), nil
}

// Download implements the Uploader interface
func (u *FTPUploader) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	fullPath := path.Join(u.config.BaseDir, filepath)
	reader, err := u.client.Retr(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return reader, nil
}

// Delete implements the Uploader interface
func (u *FTPUploader) Delete(ctx context.Context, filepath string) error {
	fullPath := path.Join(u.config.BaseDir, filepath)
	if err := u.client.Delete(fullPath); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}

// GetURL implements the Uploader interface
func (u *FTPUploader) GetURL(ctx context.Context, filepath string) (string, error) {
	fullPath := path.Join(u.config.BaseDir, filepath)
	return fmt.Sprintf("%s/%s", u.baseURL, fullPath), nil
}

// Close closes the FTP connection
func (u *FTPUploader) Close() error {
	return u.client.Quit()
} 