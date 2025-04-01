package configtype

import (
	"encoding"
	"encoding/base64"

	"github.com/pkg/errors"
)

var _ encoding.TextUnmarshaler = (*Base64)(nil)

// Base64 represents a base64-encoded string value.
// It implements encoding.TextUnmarshaler to allow loading from environment variables.
//
// Example usage:
//
//	type AppConfig struct {
//		Secret configtype.Base64 `env:"SECRET_CONFIG"`
//	}
//
//	// Set environment variable with base64-encoded string
//	// export SECRET_CONFIG=dGVzdC1zZWNyZXQ=
//
//	// The base64-encoded string contains: "test-secret"
//
//	// Load configuration
//	var config AppConfig
//	if err := env.Parse(&config); err != nil {
//		log.Fatal(err)
//	}
//
//	// Access configuration
//	fmt.Printf("Secret: %s\n", config.Secret)
type Base64 string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It decodes the base64-encoded text into a string.
func (b *Base64) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return errors.Wrapf(err, "failed to decode base64 string")
	}

	*b = Base64(decoded)
	return nil
}
