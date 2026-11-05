// Package config loads the service configuration from an optional YAML file
// and environment variables. Environment variables take precedence and use the
// upper-cased key path joined by underscores, e.g. DATABASE_HOST or HTTP_PORT.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config is the root configuration of the service.
type Config struct {
	App      App      `mapstructure:"app"`
	HTTP     HTTP     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
	Database Database `mapstructure:"database"`
}

// App holds metadata about the running service.
type App struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// HTTP holds the HTTP server settings.
type HTTP struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

// Addr returns the address the HTTP server listens on.
func (h HTTP) Addr() string {
	return net.JoinHostPort(h.Host, strconv.Itoa(h.Port))
}

// Log holds the logger settings.
type Log struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json or text
}

// Database holds the PostgreSQL connection settings.
type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxConns int32  `mapstructure:"max_conns"`
}

// URL returns the connection string in postgres:// URL form.
func (d Database) URL() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     net.JoinHostPort(d.Host, strconv.Itoa(d.Port)),
		Path:     d.Name,
		RawQuery: url.Values{"sslmode": {d.SSLMode}}.Encode(),
	}
	return u.String()
}

// Load reads the configuration from the YAML file at path (if it exists),
// applies environment variable overrides and validates the result.
func Load(path string) (Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return Config{}, fmt.Errorf("read config file %q: %w", path, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate reports whether the configuration is usable.
func (c Config) Validate() error {
	var errs []error

	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		errs = append(errs, fmt.Errorf("http.port must be between 1 and 65535, got %d", c.HTTP.Port))
	}
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		errs = append(errs, fmt.Errorf("log.level must be one of debug, info, warn, error, got %q", c.Log.Level))
	}
	switch c.Log.Format {
	case "json", "text":
	default:
		errs = append(errs, fmt.Errorf("log.format must be json or text, got %q", c.Log.Format))
	}
	if c.Database.Host == "" {
		errs = append(errs, errors.New("database.host is required"))
	}
	if c.Database.Name == "" {
		errs = append(errs, errors.New("database.name is required"))
	}

	return errors.Join(errs...)
}

// setDefaults registers every key so that AutomaticEnv can override it even
// when it is absent from the configuration file.
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "api-microservice")
	v.SetDefault("app.version", "dev")
	v.SetDefault("app.environment", "local")

	v.SetDefault("http.host", "")
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.read_header_timeout", 5*time.Second)
	v.SetDefault("http.read_timeout", 15*time.Second)
	v.SetDefault("http.write_timeout", 15*time.Second)
	v.SetDefault("http.idle_timeout", 60*time.Second)
	v.SetDefault("http.shutdown_timeout", 10*time.Second)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "postgres")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_conns", 10)
}
