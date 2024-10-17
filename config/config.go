package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port         int
	ServiceName  string
	LogLevel     string
	Tracing      bool
	ReadTimeout  int // in seconds
	WriteTimeout int // in seconds
	IdleTimeout  int // in seconds
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

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

var cfg *Config
var once sync.Once

// LoadConfig loads the configuration from the .env file and ensures it's only loaded once.
func LoadConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

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
			tracing = false
		}

		// Load read, write, and idle timeouts from environment variables
		readTimeout := getEnvAsInt("SERVER_READ_TIMEOUT", 5)
		writeTimeout := getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)
		idleTimeout := getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)

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

		maxIdleConnection := 10
		maxActiveConnection := 100
		maxIdleTime := 60
		dbTimeout := 30

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
			dbTimeout, err = strconv.Atoi(timeoutStr)
			if err != nil {
				panic("Invalid DB_TIMEOUT environment variable: must be an integer")
			}
		}

		// Load Redis config
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			panic("REDIS_URL environment variable is not set")
		}

		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDBStr := os.Getenv("REDIS_DB")
		redisDB := 0

		if redisDBStr != "" {
			var err error
			redisDB, err = strconv.Atoi(redisDBStr)
			if err != nil {
				panic("Invalid REDIS_DB environment variable: must be an integer")
			}
		}

		cfg = &Config{
			Server: ServerConfig{
				Port:         port,
				ServiceName:  serviceName,
				LogLevel:     logLevel,
				Tracing:      tracing,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
				IdleTimeout:  idleTimeout,
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
				Timeout:             dbTimeout,
			},
			Redis: RedisConfig{
				URL:      redisURL,
				Password: redisPassword,
				DB:       redisDB,
			},
		}
	})

	return cfg
}

// Helper function to read environment variable as integer with a fallback default
func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default %d", name, defaultValue)
		return defaultValue
	}
	return value
}
