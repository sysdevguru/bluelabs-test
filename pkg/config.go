package pkg

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	envPrefix = "settings"
)

// Server contains the configuration for the HTTP server.
type Server struct {
	// IdleTimeout is the maximum amount of time to wait for an open connection
	// when processing no requests and keep-alives are enabled. If this value is
	// 0, ReadTimeout value be used.
	IdleTimeout time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`

	// Port is the HTTP server port.
	Port int `envconfig:"PORT" default:"8080"`

	// ReadTimeout is the maximum duration for reading the entire request,
	// including the body.
	ReadTimeout time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"1s"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"2s"`
}

// Database contains configuration for the Postgres Database.
type Database struct {
	URL                string `envconfig:"DATABASE_URL" required:"true"`
	LogLevel           string `envconfig:"DATABASE_LOG_LEVEL" default:"warn"`
	MaxOpenConnections int    `envconfig:"DATABASE_MAX_OPEN_CONNECTIONS" default:"10"`
}

// Config is the global config struct.
type Config struct {
	Database Database
	Server   Server
}

// Load configuration from environment.
func Load() (Config, error) {
	config := Config{}

	if err := envconfig.Process(envPrefix, &config); err != nil {
		return config, err
	}

	return config, nil
}
