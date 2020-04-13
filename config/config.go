package config

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Config struct {
	Server serverConfig `mapstructure:"server_config"`
}

type serverConfig struct {
	HTTPPort int `mapstructure:"http_port"`
}

func Default() *Config {
	return &Config{}
}

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

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	return nil
}

// ServeHTTP returns all the live configurations, FOR DEBUG only
func (c *Config) ServeHTTP(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c)
}
