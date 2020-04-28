package config

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type (
	// Config configuration schema
	Config struct {
		Server serverConfig `json:"server" mapstructure:"server"`
		Video  videoConfig  `json:"video" mapstructure:"video"`
		GCS    gcsConfig    `json:"gcs" mapstructure:"gcs"`
		S3     s3Config     `json:"s3" mapstructure:"s3"`
	}

	serverConfig struct {
		HTTPPort int `mapstructure:"http_port" json:"http_port"`
	}

	gcsConfig struct {
		Enable bool   `mapstructure:"enable" json:"enable"`
		Bucket string `mapstructure:"bucket" json:"bucket"`
	}

	s3Config struct {
		Enable bool   `mapstructure:"enable" json:"enable"`
		Bucket string `mapstructure:"bucket" json:"bucket"`
	}

	videoConfig struct {
		MaxSize            int64           `json:"max_size"`
		MaxSizeMB          int64           `mapstructure:"max_size_mb" json:"max_size_mb"`
		AllowedTypes       map[string]bool `mapstructure:"allowed_types" json:"allowed_types"`
		AllowedTypesString string          `json:"allowed_types_string"`
	}
)

// Default returns Default configurations
func Default() *Config {
	return &Config{}
}

// Load configuration from path
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := Default()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if cfg.Video.MaxSizeMB > 0 {
		cfg.Video.MaxSize = cfg.Video.MaxSizeMB * 1024 * 1024
	}

	if cfg.Video.AllowedTypes != nil {
		tmp := ""
		for t := range cfg.Video.AllowedTypes {
			if tmp != "" {
				tmp += ","
			}
			tmp += t
		}
		cfg.Video.AllowedTypesString = tmp
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Server.HTTPPort == 0 {
		return fmt.Errorf("HTTP port cannot be empty")
	}

	if c.GCS.Enable && c.S3.Enable {
		return fmt.Errorf("Cannot enable both storages")
	}

	if c.GCS.Enable && c.GCS.Bucket == "" {
		return fmt.Errorf("GCS bucket cannot be empty")
	}

	if c.S3.Enable && c.S3.Bucket == "" {
		return fmt.Errorf("S3 bucket cannot be empty")
	}

	return nil
}

// ServeHTTP returns all the live configurations, FOR DEBUG only
func (c *Config) ServeHTTP(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c)
}
