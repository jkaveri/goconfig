package configtype

import (
	"encoding"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

var _ encoding.TextUnmarshaler = (*JSONFile[any])(nil)

// JSONFile represents a configuration file in JSON format.
// It implements encoding.TextUnmarshaler to allow loading from environment variables.
// The generic type T specifies the type of the configuration data.
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
//		DB configtype.JSONFile[DBConfig] `env:"DB_CONFIG"`
//	}
//
//	// Set environment variable to point to JSON config file
//	// export DB_CONFIG=/path/to/db_config.json
//
//	// The JSON file at /path/to/db_config.json should contain:
//	// {
//	//   "host": "localhost",
//	//   "port": 5432,
//	//   "database": "mydb"
//	// }
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
type JSONFile[T any] struct {
	// FilePath is the path to the JSON configuration file
	FilePath string
	// Data contains the parsed configuration data
	Data T
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It reads the JSON file path from the provided text and loads the configuration.
// The file path can contain environment variables that will be expanded.
func (f *JSONFile[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	f.FilePath = string(data)

	return f.parseJSONFile()
}

// parseJSONFile reads and parses the JSON configuration file.
// It expands environment variables in the file content before parsing.
func (f *JSONFile[T]) parseJSONFile() error {
	jsonData, err := os.ReadFile(f.FilePath)
	if err != nil {
		return errors.Wrapf(err, "cannot load json file: %s", f.FilePath)
	}

	jsonStr := os.ExpandEnv(string(jsonData))

	if err := json.Unmarshal([]byte(jsonStr), &f.Data); err != nil {
		return errors.Wrapf(err, "failed to unmarshal json config: %s, data: %s", f.FilePath, jsonStr)
	}

	return nil
}

// Reload reloads the JSON configuration file.
// It is useful when the configuration file has been modified and needs to be reloaded.
// If no file path is set, it returns nil without doing anything.
func (f *JSONFile[T]) Reload() error {
	if f.FilePath == "" {
		return nil
	}

	return f.parseJSONFile()
}
