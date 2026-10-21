package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
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
	Host            string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port            int           `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User            string        `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Name            string        `yaml:"name" env:"DB_NAME" env-default:"warehouse_db"`
	SSLMode         string        `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	Password        string        `yaml:"-" env:"DB_PASSWORD"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"5m"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"-" env:"JWT_SECRET"`
	RefreshCleanupInterval time.Duration `yaml: "refresh_cleanup_interval" env-default: "24h"`
}

type RateLimitConfig struct {
	AuthBurst       int           `yaml:"auth_burst" env-default:"5"`
	AuthInterval    time.Duration `yaml:"auth_interval" env-default:"12s"`
	RefreshBurst    int           `yaml:"refresh_burst" env-default: "10"`
	RefreshInterval time.Duration `yaml:"refresh_interval" env-default: "6s"`
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
	if cfg.DB.Password, err = secretValue("DB_PASSWORD", cfg.DB.Password); err != nil {
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

func (c DBConfig) DSN() string {
	q := url.Values{}
	q.Set("sslmode", c.SSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     "/" + c.Name,
		RawQuery: q.Encode(),
	}
	return u.String()
}
