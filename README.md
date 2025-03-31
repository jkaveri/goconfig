# goconfig

A flexible and powerful Go configuration utility that loads environment variables into structs with support for various data types and customization options.

## Features

- Load environment variables into struct fields automatically
- Support for nested structs and pointers
- Customizable field name transformation
- Support for various data types:
  - Basic types (string, bool, int, uint, float)
  - Duration
  - Slices (with custom separator)
  - Maps (via JSON)
  - Custom types implementing `encoding.TextUnmarshaler`
- Configurable prefix and separators
- Tag-based field mapping with `env` and `alias` tags
- Anonymous struct embedding support
- Panic recovery for safe error handling


## Installation

```bash
go get github.com/jkaveri/goconfig
```

## Usage

### Basic Example

```go
import "github.com/jkaveri/goconfig"

type Config struct {
    Host     string `env:"HOST"`
    Port     int    `env:"PORT"`
    Timeout  time.Duration
    Debug    bool
    Numbers  []int
    Settings map[string]string
}

func main() {
    var cfg Config
    if err := goconfig.Load(&cfg); err != nil {
        log.Fatal(err)
    }
}
```


## ConfigType Package

The `configtype` package provides additional functionality for loading configuration from various file formats and handling special data types:

- Support for multiple configuration formats:
  - JSON files
  - YAML files
  - TOML files
- Base64 encoding support for sensitive data
- Environment variable expansion in both file paths and configuration content
- Hot reloading capability for configuration files
- Generic type support for type-safe configuration loading

Example usage with configtype:

```go
import (
    "github.com/jkaveri/goconfig"
    "github.com/jkaveri/goconfig/configtype"
)

type DBConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
}

type AppConfig struct {
    DB     configtype.JSONFile[DBConfig] `env:"DB_CONFIG"`
    Secret configtype.Base64             `env:"API_SECRET"`
}

func main() {
    var cfg AppConfig
    if err := goconfig.Load(&cfg); err != nil {
        log.Fatal(err)
    }

    // Access configuration
    fmt.Printf("Database: %s:%d/%s\n",
        cfg.DB.Data.Host,
        cfg.DB.Data.Port,
        cfg.DB.Data.Database)
    fmt.Printf("Secret: %s\n", cfg.Secret)
}
```

### With Options

```go
import "github.com/jkaveri/goconfig"

type Config struct {
    Host string `env:"HOST"`
    Port int    `env:"PORT"`
}

func main() {
    var cfg Config
    loader := goconfig.New(
        goconfig.WithPrefix("APP"),
        goconfig.WithSeparator("."),
        goconfig.WithArraySeparator(";"),
    )
    if err := loader.LoadToStruct(&cfg); err != nil {
        log.Fatal(err)
    }
}
```

### Field Tags

- `env`: Exact environment variable name
- `alias`: Alternative name for the field (will be combined with prefix)

```go
type Config struct {
    Host     string `env:"SERVER_HOST"`     // Uses exact name
    Port     int    `alias:"server_port"`   // Uses prefix + name
    Timeout  time.Duration                  // Uses field name transformed
}
```

## Environment Variables

Given the following struct:

```go
type Config struct {
    Host     string        `env:"HOST"`
    Port     int
    Timeout  time.Duration
    Debug    bool
    Numbers  []int
    Settings map[string]string
}
```

The following environment variables would be loaded:

```bash
HOST=localhost
PORT=8080
TIMEOUT=5s
DEBUG=true
NUMBERS=1,2,3,4
SETTINGS={"key":"value"}
```

## Options

- `WithPrefix(prefix string)`: Set prefix for all environment variables
- `WithSeparator(sep string)`: Set separator for nested fields (default: "_")
- `WithArraySeparator(sep string)`: Set separator for array values (default: ",")
- `WithFieldNameTransformer(fn func(string) string)`: Set custom field name transformation

## License

MIT License
