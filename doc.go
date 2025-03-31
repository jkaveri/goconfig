// Package goconfig provides a flexible and powerful configuration utility for Go applications.
// It allows loading environment variables into structs with support for various data types
// and customization options.
//
// Key features:
//   - Load environment variables into struct fields automatically
//   - Support for nested structs and pointers
//   - Customizable field name transformation
//   - Support for various data types (string, bool, int, uint, float, duration, slices, maps)
//   - Configurable prefix and separators
//   - Tag-based field mapping with env and alias tags
//   - Anonymous struct embedding support
//   - Panic recovery for safe error handling
//
// Example usage:
//
//	type Config struct {
//	    Host     string        `env:"HOST"`
//	    Port     int
//	    Timeout  time.Duration
//	    Debug    bool
//	    Numbers  []int
//	    Settings map[string]string
//	}
//
//	var cfg Config
//	if err := goconfig.Load(&cfg); err != nil {
//	    log.Fatal(err)
//	}
package goconfig
