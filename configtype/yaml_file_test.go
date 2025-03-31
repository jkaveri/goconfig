package configtype

import (
	"os"
	"path/filepath"
	"testing"
)

type TestYAMLConfig struct {
	Name    string `yaml:"name"`
	Version int    `yaml:"version"`
}

func TestYAMLFileImplementation(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Test case 1: Valid YAML file
	t.Run("valid yaml file", func(t *testing.T) {
		// Create a test YAML file
		yamlContent := `name: test
version: 1`
		filePath := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(filePath, []byte(yamlContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test loading the file
		config := &YAMLFile[TestYAMLConfig]{
			FilePath: filePath,
		}
		if err := config.parseYAMLFile(); err != nil {
			t.Errorf("Failed to parse YAML file: %v", err)
		}

		// Verify the data
		if config.Data.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", config.Data.Name)
		}
		if config.Data.Version != 1 {
			t.Errorf("Expected version 1, got %d", config.Data.Version)
		}
	})

	// Test case 2: Environment variable expansion
	t.Run("environment variable expansion", func(t *testing.T) {
		// Set test environment variable
		os.Setenv("TEST_NAME", "env_test")
		defer os.Unsetenv("TEST_NAME")

		// Create a test YAML file with environment variable
		yamlContent := `name: $TEST_NAME
version: 1`
		filePath := filepath.Join(tmpDir, "env_config.yaml")
		if err := os.WriteFile(filePath, []byte(yamlContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test loading the file
		config := &YAMLFile[TestYAMLConfig]{
			FilePath: filePath,
		}
		if err := config.parseYAMLFile(); err != nil {
			t.Errorf("Failed to parse YAML file: %v", err)
		}

		// Verify the data
		if config.Data.Name != "env_test" {
			t.Errorf("Expected name 'env_test', got '%s'", config.Data.Name)
		}
	})

	// Test case 3: UnmarshalText
	t.Run("unmarshal text", func(t *testing.T) {
		// Create a test YAML file
		yamlContent := `name: test
version: 1`
		filePath := filepath.Join(tmpDir, "unmarshal_config.yaml")
		if err := os.WriteFile(filePath, []byte(yamlContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test UnmarshalText
		config := &YAMLFile[TestYAMLConfig]{}
		if err := config.UnmarshalText([]byte(filePath)); err != nil {
			t.Errorf("Failed to unmarshal text: %v", err)
		}

		// Verify the data
		if config.Data.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", config.Data.Name)
		}
	})

	// Test case 4: Reload
	t.Run("reload", func(t *testing.T) {
		// Create initial test YAML file
		yamlContent := `name: test
version: 1`
		filePath := filepath.Join(tmpDir, "reload_config.yaml")
		if err := os.WriteFile(filePath, []byte(yamlContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Load initial config
		config := &YAMLFile[TestYAMLConfig]{
			FilePath: filePath,
		}
		if err := config.parseYAMLFile(); err != nil {
			t.Fatalf("Failed to parse initial YAML file: %v", err)
		}

		// Update the file
		newContent := `name: reloaded
version: 2`
		if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
			t.Fatalf("Failed to update test file: %v", err)
		}

		// Reload the config
		if err := config.Reload(); err != nil {
			t.Errorf("Failed to reload config: %v", err)
		}

		// Verify the reloaded data
		if config.Data.Name != "reloaded" {
			t.Errorf("Expected name 'reloaded', got '%s'", config.Data.Name)
		}
		if config.Data.Version != 2 {
			t.Errorf("Expected version 2, got %d", config.Data.Version)
		}
	})

	// Test case 5: Error cases
	t.Run("error cases", func(t *testing.T) {
		// Test non-existent file
		config := &YAMLFile[TestYAMLConfig]{
			FilePath: "non_existent.yaml",
		}
		if err := config.parseYAMLFile(); err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}

		// Test invalid YAML
		invalidYamlContent := `name: test
version: invalid`
		filePath := filepath.Join(tmpDir, "invalid_config.yaml")
		if err := os.WriteFile(filePath, []byte(invalidYamlContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		config = &YAMLFile[TestYAMLConfig]{
			FilePath: filePath,
		}
		if err := config.parseYAMLFile(); err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}

		// Test empty file path in Reload
		emptyConfig := &YAMLFile[TestYAMLConfig]{}
		if err := emptyConfig.Reload(); err != nil {
			t.Errorf("Expected nil error for empty file path, got %v", err)
		}
	})
}
