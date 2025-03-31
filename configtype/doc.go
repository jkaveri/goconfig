// Package configtype provides a flexible and type-safe way to load and manage configuration
// from various sources. It supports multiple configuration formats including JSON, YAML, and TOML,
// with built-in environment variable expansion support.
//
// The package implements the encoding.TextUnmarshaler interface to allow configuration loading
// from environment variables, making it easy to integrate with various configuration management systems.
//
// Key features:
//   - Support for multiple configuration formats (JSON, YAML, TOML)
//   - Environment variable expansion in both file paths and configuration content
//   - Generic type support for type-safe configuration loading
//   - Base64 encoding support for sensitive data
//   - Hot reloading capability for configuration files
//
// Example usage:
//
//	type DBConfig struct {
//		Host     string `json:"host"`
//		Port     int    `json:"port"`
//		Database string `json:"database"`
//	}
//
//	type AppConfig struct {
//		DB     configtype.JSONFile[DBConfig] `env:"DB_CONFIG"`
//		Secret configtype.Base64             `env:"API_SECRET"`
//	}
//
//	// Set environment variables
//	// export DB_CONFIG=/path/to/db_config.json
//	// export API_SECRET=dGVzdC1zZWNyZXQ=
//
//	// Load configuration
//	var config AppConfig
//	if err := env.Parse(&config); err != nil {
//		log.Fatal(err)
//	}
//
//	// Access configuration
//	fmt.Printf("Database: %s:%d/%s\n",
//		config.DB.Data.Host,
//		config.DB.Data.Port,
//		config.DB.Data.Database)
//	fmt.Printf("Secret: %s\n", config.Secret)
//
// The package provides the following types:
//   - JSONFile[T]: For loading JSON configuration files
//   - YAMLFile[T]: For loading YAML configuration files
//   - TOMLFile[T]: For loading TOML configuration files
//   - Base64: For handling base64-encoded configuration values
//
// Each file-based configuration type supports:
//   - Environment variable expansion in file paths
//   - Environment variable expansion in configuration content
//   - Hot reloading via the Reload() method
//   - Type-safe configuration loading through generics
package configtype
