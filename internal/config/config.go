package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP      HTTPConfig      `yaml:"http"`
	DB        DBConfig        `yaml:"db"`
	Auth      AuthConfig      `yaml:"auth"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type HTTPConfig struct {
	Addr              string        `yaml:"addr" env:"HTTP_ADDR" env-default:":8080"`
	RequestTimeout    time.Duration `yaml:"request_timeout" env-default:"5s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"5s"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout" env-default:"15s"`
}

type DBConfig struct {
	URL             string        `yaml:"-" env:"DATABASE_URL"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"5m"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"-" env:"JWT_SECRET"`
}

type RateLimitConfig struct {
	AuthBurst    int           `yaml:"auth_burst" env-default:"5"`
	AuthInterval time.Duration `yaml:"auth_interval" env-default:"12s"`
}

func Load() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yaml"
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var err error
	if cfg.DB.URL, err = secretValue("DATABASE_URL", cfg.DB.URL); err != nil {
		return nil, err
	}
	if cfg.Auth.JWTSecret, err = secretValue("JWT_SECRET", cfg.Auth.JWTSecret); err != nil {
		return nil, err
	}

	if len(cfg.Auth.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	return &cfg, nil
}

func secretValue(name, current string) (string, error) {
	if path := os.Getenv(name + "_FILE"); path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read %s_FILE: %w", name, err)
		}
		return strings.TrimSpace(string(b)), nil
	}
	if current == "" {
		return "", fmt.Errorf("%s (or %s_FILE) is required", name, name)
	}
	return current, nil
}