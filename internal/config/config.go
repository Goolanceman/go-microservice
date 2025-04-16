package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Features FeaturesConfig `mapstructure:"features"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Environment  string `mapstructure:"environment"`
	LogLevel     string `mapstructure:"log_level"`
	AllowOrigins string `mapstructure:"allow_origins"`
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// KafkaConfig holds Kafka connection and consumer settings
type KafkaConfig struct {
	Brokers           []string       `mapstructure:"brokers"`
	GroupID           string         `mapstructure:"group_id"`
	Topics            []string       `mapstructure:"topics"`
	AutoOffsetReset   string         `mapstructure:"auto_offset_reset"`
	SessionTimeout    time.Duration  `mapstructure:"session_timeout"`
	HeartbeatInterval time.Duration  `mapstructure:"heartbeat_interval"`
	MaxPollInterval   time.Duration  `mapstructure:"max_poll_interval"`
	MaxPollRecords    int            `mapstructure:"max_poll_records"`
	Security          SecurityConfig `mapstructure:"security"`
}

// SecurityConfig holds Kafka security settings
type SecurityConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	CAFile     string `mapstructure:"ca_file"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
	SkipVerify bool   `mapstructure:"skip_verify"`
}

// UploadConfig holds upload service configuration
type UploadConfig struct {
	Backend    string      `mapstructure:"backend"` // "s3", "minio", or "gcs"
	S3Config   S3Config      `mapstructure:"s3"`
	MinioConfig MinioConfig   `mapstructure:"minio"`
	GCSConfig  GCSConfig     `mapstructure:"gcs"`
}

type S3Config struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

type GCSConfig struct {
	ProjectID      string `mapstructure:"project_id"`
	Bucket         string `mapstructure:"bucket"`
	CredentialsFile string `mapstructure:"credentials_file"`
}

// FeaturesConfig holds feature flags
type FeaturesConfig struct {
	EnableRedis bool `mapstructure:"enable_redis"`
	EnableKafka bool `mapstructure:"enable_kafka"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
} 