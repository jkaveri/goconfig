package configtype

import (
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

func TestJSONFileImplementation(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Test case 1: Valid JSON file
	t.Run("valid json file", func(t *testing.T) {
		// Create a test JSON file
		jsonContent := `{"name": "test", "version": 1}`
		filePath := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(filePath, []byte(jsonContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test loading the file
		config := &JSONFile[TestConfig]{
			FilePath: filePath,
		}
		if err := config.parseJSONFile(); err != nil {
			t.Errorf("Failed to parse JSON file: %v", err)
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

		// Create a test JSON file with environment variable
		jsonContent := `{"name": "$TEST_NAME", "version": 1}`
		filePath := filepath.Join(tmpDir, "env_config.json")
		if err := os.WriteFile(filePath, []byte(jsonContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test loading the file
		config := &JSONFile[TestConfig]{
			FilePath: filePath,
		}
		if err := config.parseJSONFile(); err != nil {
			t.Errorf("Failed to parse JSON file: %v", err)
		}

		// Verify the data
		if config.Data.Name != "env_test" {
			t.Errorf("Expected name 'env_test', got '%s'", config.Data.Name)
		}
	})

	// Test case 3: UnmarshalText
	t.Run("unmarshal text", func(t *testing.T) {
		// Create a test JSON file
		jsonContent := `{"name": "test", "version": 1}`
		filePath := filepath.Join(tmpDir, "unmarshal_config.json")
		if err := os.WriteFile(filePath, []byte(jsonContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Test UnmarshalText
		config := &JSONFile[TestConfig]{}
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
		// Create initial test JSON file
		jsonContent := `{"name": "test", "version": 1}`
		filePath := filepath.Join(tmpDir, "reload_config.json")
		if err := os.WriteFile(filePath, []byte(jsonContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Load initial config
		config := &JSONFile[TestConfig]{
			FilePath: filePath,
		}
		if err := config.parseJSONFile(); err != nil {
			t.Fatalf("Failed to parse initial JSON file: %v", err)
		}

		// Update the file
		newContent := `{"name": "reloaded", "version": 2}`
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
		config := &JSONFile[TestConfig]{
			FilePath: "non_existent.json",
		}
		if err := config.parseJSONFile(); err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}

		// Test invalid JSON
		invalidJSONContent := `{"name": "test", "version": invalid}`
		filePath := filepath.Join(tmpDir, "invalid_config.json")
		if err := os.WriteFile(filePath, []byte(invalidJSONContent), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		config = &JSONFile[TestConfig]{
			FilePath: filePath,
		}
		if err := config.parseJSONFile(); err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}

		// Test empty file path in Reload
		emptyConfig := &JSONFile[TestConfig]{}
		if err := emptyConfig.Reload(); err != nil {
			t.Errorf("Expected nil error for empty file path, got %v", err)
		}
	})
}
