package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port        int
	ServiceName string
	LogLevel    string
	Tracing     bool
}

type JWTConfig struct {
	Secret string
	Expiry int
}

type DatabaseConfig struct {
	URL                 string
	MaxIdleConnection   int
	MaxActiveConnection int
	MaxIdleTime         int // in seconds
	Timeout             int // in seconds
}

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
}

var cfg *Config
var once sync.Once

// LoadConfig loads the configuration from the .env file and ensures it's only loaded once.
func LoadConfig() *Config {
	once.Do(func() {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		// Load server config
		portStr := os.Getenv("PORT")
		if portStr == "" {
			panic("PORT environment variable is not set")
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			panic("Invalid PORT environment variable: must be an integer")
		}

		serviceName := os.Getenv("SERVICE_NAME")
		if serviceName == "" {
			panic("SERVICE_NAME environment variable is not set")
		}

		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			panic("LOG_LEVEL environment variable is not set")
		}

		tracingStr := os.Getenv("TRACING")
		if tracingStr == "" {
			panic("TRACING environment variable is not set")
		}

		tracing, err := strconv.ParseBool(tracingStr)
		if err != nil {
			tracing = false // Default to false if there's an error
		}

		// Load JWT config
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			panic("JWT_SECRET environment variable is not set")
		}

		jwtExpiryStr := os.Getenv("JWT_EXPIRY")
		if jwtExpiryStr == "" {
			panic("JWT_EXPIRY environment variable is not set")
		}

		jwtExpiry, err := strconv.Atoi(jwtExpiryStr)
		if err != nil {
			panic("Invalid JWT_EXPIRY environment variable: must be an integer")
		}

		// Load database config
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			panic("DATABASE_URL environment variable is not set")
		}

		maxIdleConnStr := os.Getenv("DB_MAX_IDLE_CONNECTION")
		maxActiveConnStr := os.Getenv("DB_MAX_ACTIVE_CONNECTION")
		maxIdleTimeStr := os.Getenv("DB_MAX_IDLE_TIME")
		timeoutStr := os.Getenv("DB_TIMEOUT")

		maxIdleConnection := 10    // Default value
		maxActiveConnection := 100 // Default value
		maxIdleTime := 60          // Default value in seconds
		timeout := 30              // Default value in seconds

		if maxIdleConnStr != "" {
			var err error
			maxIdleConnection, err = strconv.Atoi(maxIdleConnStr)
			if err != nil {
				panic("Invalid DB_MAX_IDLE_CONNECTION environment variable: must be an integer")
			}
		}

		if maxActiveConnStr != "" {
			var err error
			maxActiveConnection, err = strconv.Atoi(maxActiveConnStr)
			if err != nil {
				panic("Invalid DB_MAX_ACTIVE_CONNECTION environment variable: must be an integer")
			}
		}

		if maxIdleTimeStr != "" {
			var err error
			maxIdleTime, err = strconv.Atoi(maxIdleTimeStr)
			if err != nil {
				panic("Invalid DB_MAX_IDLE_TIME environment variable: must be an integer")
			}
		}

		if timeoutStr != "" {
			var err error
			timeout, err = strconv.Atoi(timeoutStr)
			if err != nil {
				panic("Invalid DB_TIMEOUT environment variable: must be an integer")
			}
		}

		cfg = &Config{
			Server: ServerConfig{
				Port:        port,
				ServiceName: serviceName,
				LogLevel:    logLevel,
				Tracing:     tracing,
			},
			JWT: JWTConfig{
				Secret: jwtSecret,
				Expiry: jwtExpiry,
			},
			Database: DatabaseConfig{
				URL:                 databaseURL,
				MaxIdleConnection:   maxIdleConnection,
				MaxActiveConnection: maxActiveConnection,
				MaxIdleTime:         maxIdleTime,
				Timeout:             timeout,
			},
		}
	})

	return cfg
}
