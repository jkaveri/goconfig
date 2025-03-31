package goconfig

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type BasicConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Timeout  time.Duration
	Debug    bool
	Numbers  []int
	Settings map[string]string
}

type NestedConfig struct {
	Server BasicConfig
	DB     struct {
		Host     string
		Port     int
		Password string `env:"DB_PASSWORD"`
	}
}

type PointerConfig struct {
	Server *BasicConfig
	DB     *struct {
		Host string
		Port int
	}
}

type CustomTransformerConfig struct {
	DBConnection string
	APIKey       string
	UserID       string
}

func TestBasicConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("TIMEOUT", "5s")
	os.Setenv("DEBUG", "true")
	os.Setenv("NUMBERS", "1,2,3,4")
	os.Setenv("SETTINGS", `{"key":"value"}`)

	var cfg BasicConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.True(t, cfg.Debug)
	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Numbers)
	assert.Equal(t, map[string]string{"key": "value"}, cfg.Settings)
}

func TestNestedConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_PASSWORD", "secret")

	var cfg NestedConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "db.example.com", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "secret", cfg.DB.Password)
}

func TestPointerConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")

	var cfg PointerConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.NotNil(t, cfg.Server)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.NotNil(t, cfg.DB)
	assert.Equal(t, "db.example.com", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
}

func TestCustomTransformer(t *testing.T) {
	// Set up test environment variables
	os.Setenv("dbconnection", "postgres://localhost:5432/db")
	os.Setenv("apikey", "secret-key")
	os.Setenv("userid", "123")

	var cfg CustomTransformerConfig
	loader := New(WithKeyTransformer(strings.ToLower))
	err := loader.Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, "postgres://localhost:5432/db", cfg.DBConnection)
	assert.Equal(t, "secret-key", cfg.APIKey)
	assert.Equal(t, "123", cfg.UserID)
}

func TestPrefixAndSeparator(t *testing.T) {
	// Set up test environment variables
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("APP.DB.HOST", "db.example.com")

	var cfg NestedConfig
	loader := New(
		WithPrefix("APP"),
		WithSeparator("."),
	)
	err := loader.Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "db.example.com", cfg.DB.Host)
}

func TestArraySeparator(t *testing.T) {
	// Set up test environment variables
	os.Setenv("NUMBERS", "1;2;3;4")

	var cfg BasicConfig
	loader := New(WithArraySeparator(";"))
	err := loader.Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Numbers)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		config  interface{}
		env     map[string]string
		wantErr bool
	}{
		{
			name:    "invalid pointer",
			config:  BasicConfig{}, // Not a pointer
			wantErr: true,
		},
		{
			name:   "invalid port number",
			config: &BasicConfig{},
			env: map[string]string{
				"PORT": "invalid",
			},
			wantErr: true,
		},
		{
			name:   "invalid duration",
			config: &BasicConfig{},
			env: map[string]string{
				"TIMEOUT": "invalid",
			},
			wantErr: true,
		},
		{
			name:   "invalid json map",
			config: &BasicConfig{},
			env: map[string]string{
				"SETTINGS": "invalid json",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			err := Load(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
