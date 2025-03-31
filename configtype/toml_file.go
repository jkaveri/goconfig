package configtype

import (
	"encoding"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var _ encoding.TextUnmarshaler = (*TOMLFile[interface{}])(nil)

// TOMLFile represents a configuration file in TOML format.
// It implements encoding.TextUnmarshaler to allow loading from environment variables.
// The generic type T specifies the type of the configuration data.
//
// Example usage:
//
//	type DBConfig struct {
//		Host     string `toml:"host"`
//		Port     int    `toml:"port"`
//		Database string `toml:"database"`
//	}
//
//	type AppConfig struct {
//		DB configtype.TOMLFile[DBConfig] `env:"DB_CONFIG"`
//	}
//
//	// Set environment variable to point to TOML config file
//	// export DB_CONFIG=/path/to/db_config.toml
//	// export DB_HOST=localhost
//	// export DB_PORT=5432
//
//	// The TOML file at /path/to/db_config.toml can contain environment variables:
//	// host = "$DB_HOST"
//	// port = $DB_PORT
//	// database = "mydb"
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
type TOMLFile[T any] struct {
	// FilePath is the path to the TOML configuration file
	FilePath string
	// Data contains the parsed configuration data
	Data T
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It reads the TOML file path from the provided text and loads the configuration.
// The file path can contain environment variables that will be expanded.
func (f *TOMLFile[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	f.FilePath = string(data)
	return f.parseTOMLFile()
}

// parseTOMLFile reads and parses the TOML configuration file.
// It expands any environment variables in the file path and file content.
func (f *TOMLFile[T]) parseTOMLFile() error {
	if f.FilePath == "" {
		return nil
	}

	// Expand environment variables in the file path
	expandedPath := os.ExpandEnv(f.FilePath)

	// Read the file
	content, err := os.ReadFile(expandedPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read TOML file: %s", expandedPath)
	}

	// Expand environment variables in the content
	expandedContent := os.ExpandEnv(string(content))

	// Parse TOML content
	if _, err := toml.Decode(expandedContent, &f.Data); err != nil {
		return errors.Wrapf(err, "failed to parse TOML file: %s", expandedPath)
	}

	return nil
}

// Reload reloads the TOML configuration file.
// This is useful when the configuration file has been updated and you want to load the new values.
func (f *TOMLFile[T]) Reload() error {
	return f.parseTOMLFile()
}
