package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"github.com/goolanceman/go-microservice/internal/config"
)

// SFTPConfig holds SFTP connection settings
type SFTPConfig struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	PrivateKey string `mapstructure:"private_key"`
	BaseDir    string `mapstructure:"base_dir"`
}

// SFTPUploader implements the Uploader interface for SFTP
type SFTPUploader struct {
	client  *sftp.Client
	config  *SFTPConfig
	baseURL string
}

func init() {
	RegisterUploader("sftp", NewSFTPUploader)
}

// NewSFTPUploader creates a new SFTP uploader instance
func NewSFTPUploader(cfg interface{}) (Uploader, error) {
	sftpCfg, ok := cfg.(*SFTPConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// Configure SSH client
	config := &ssh.ClientConfig{
		User:            sftpCfg.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Set authentication method
	if sftpCfg.PrivateKey != "" {
		key, err := os.ReadFile(sftpCfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.Password(sftpCfg.Password)}
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%s", sftpCfg.Host, sftpCfg.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	// Create SFTP client
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	// Create base directory if it doesn't exist
	if sftpCfg.BaseDir != "" {
		if err := client.MkdirAll(sftpCfg.BaseDir); err != nil {
			// Ignore error if directory already exists
		}
	}

	return &SFTPUploader{
		client:  client,
		config:  sftpCfg,
		baseURL: fmt.Sprintf("sftp://%s@%s:%s", sftpCfg.Username, sftpCfg.Host, sftpCfg.Port),
	}, nil
}

// Upload implements the Uploader interface
func (u *SFTPUploader) Upload(ctx context.Context, filepath string, content io.Reader, contentType string) (string, error) {
	// Create full path
	fullPath := path.Join(u.config.BaseDir, filepath)

	// Create directories if they don't exist
	dir := path.Dir(fullPath)
	if err := u.client.MkdirAll(dir); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := u.client.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Return SFTP URL
	return fmt.Sprintf("%s/%s", u.baseURL, fullPath), nil
}

// Download implements the Uploader interface
func (u *SFTPUploader) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	fullPath := path.Join(u.config.BaseDir, filepath)
	file, err := u.client.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return file, nil
}

// Delete implements the Uploader interface
func (u *SFTPUploader) Delete(ctx context.Context, filepath string) error {
	fullPath := path.Join(u.config.BaseDir, filepath)
	if err := u.client.Remove(fullPath); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}

// GetURL implements the Uploader interface
func (u *SFTPUploader) GetURL(ctx context.Context, filepath string) (string, error) {
	fullPath := path.Join(u.config.BaseDir, filepath)
	return fmt.Sprintf("%s/%s", u.baseURL, fullPath), nil
}

// Close closes the SFTP connection
func (u *SFTPUploader) Close() error {
	return u.client.Close()
} 