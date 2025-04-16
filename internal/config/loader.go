package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// LoadConfig reads the configuration from the specified JSON file and environment variables
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName(env)
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// Enable environment variable overrides
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set defaults if not specified
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.LogLevel == "" {
		config.Server.LogLevel = "info"
	}

	return &config, nil
}

// GetConfigPath returns the path of the loaded configuration file
func GetConfigPath() string {
	return viper.ConfigFileUsed()
} 