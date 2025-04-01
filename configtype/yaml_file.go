package configtype

import (
	"encoding"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var _ encoding.TextUnmarshaler = (*YAMLFile[any])(nil)

// YAMLFile represents a configuration file in YAML format.
// It implements encoding.TextUnmarshaler to allow loading from environment variables.
// The generic type T specifies the type of the configuration data.
//
// Example usage:
//
//	type DBConfig struct {
//		Host     string `yaml:"host"`
//		Port     int    `yaml:"port"`
//		Database string `yaml:"database"`
//	}
//
//	type AppConfig struct {
//		DB configtype.YAMLFile[DBConfig] `env:"DB_CONFIG"`
//	}
//
//	// Set environment variable to point to YAML config file
//	// export DB_CONFIG=/path/to/db_config.yaml
//	// export DB_HOST=localhost
//	// export DB_PORT=5432
//
//	// The YAML file at /path/to/db_config.yaml can contain environment variables:
//	// host: $DB_HOST
//	// port: $DB_PORT
//	// database: mydb
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
type YAMLFile[T any] struct {
	// FilePath is the path to the YAML configuration file
	FilePath string
	// Data contains the parsed configuration data
	Data T
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It reads the YAML file path from the provided text and loads the configuration.
// The file path can contain environment variables that will be expanded.
func (f *YAMLFile[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	f.FilePath = string(data)
	return f.parseYAMLFile()
}

// parseYAMLFile reads and parses the YAML configuration file.
// It expands any environment variables in the file path and file content.
func (f *YAMLFile[T]) parseYAMLFile() error {
	if f.FilePath == "" {
		return nil
	}

	// Expand environment variables in the file path
	expandedPath := os.ExpandEnv(f.FilePath)

	// Read the file
	content, err := os.ReadFile(expandedPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read YAML file: %s", expandedPath)
	}

	// Expand environment variables in the content
	expandedContent := os.ExpandEnv(string(content))

	// Parse YAML content
	if err := yaml.Unmarshal([]byte(expandedContent), &f.Data); err != nil {
		return errors.Wrapf(err, "failed to parse YAML file: %s", expandedPath)
	}

	return nil
}

// Reload reloads the YAML configuration file.
// This is useful when the configuration file has been updated and you want to load the new values.
func (f *YAMLFile[T]) Reload() error {
	return f.parseYAMLFile()
}
