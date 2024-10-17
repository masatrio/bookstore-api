package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set up environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("SERVICE_NAME", "BookstoreAPI")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("TRACING", "true")
	os.Setenv("SERVER_READ_TIMEOUT", "10")
	os.Setenv("SERVER_WRITE_TIMEOUT", "20")
	os.Setenv("SERVER_IDLE_TIMEOUT", "30")
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("JWT_EXPIRY", "3600")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
	os.Setenv("DB_MAX_IDLE_CONNECTION", "5")
	os.Setenv("DB_MAX_ACTIVE_CONNECTION", "50")
	os.Setenv("DB_MAX_IDLE_TIME", "300")
	os.Setenv("DB_TIMEOUT", "100")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("REDIS_PASSWORD", "redispass")
	os.Setenv("REDIS_DB", "1")

	cfg := LoadConfig()

	// Assertions to validate loaded config values
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "BookstoreAPI", cfg.Server.ServiceName)
	assert.Equal(t, "debug", cfg.Server.LogLevel)
	assert.True(t, cfg.Server.Tracing)
	assert.Equal(t, 10, cfg.Server.ReadTimeout)
	assert.Equal(t, 20, cfg.Server.WriteTimeout)
	assert.Equal(t, 30, cfg.Server.IdleTimeout)

	assert.Equal(t, "testsecret", cfg.JWT.Secret)
	assert.Equal(t, 3600, cfg.JWT.Expiry)

	assert.Equal(t, "postgres://user:pass@localhost/db", cfg.Database.URL)
	assert.Equal(t, 5, cfg.Database.MaxIdleConnection)
	assert.Equal(t, 50, cfg.Database.MaxActiveConnection)
	assert.Equal(t, 300, cfg.Database.MaxIdleTime)
	assert.Equal(t, 100, cfg.Database.Timeout)

	assert.Equal(t, "redis://localhost:6379", cfg.Redis.URL)
	assert.Equal(t, "redispass", cfg.Redis.Password)
	assert.Equal(t, 1, cfg.Redis.DB)
}
