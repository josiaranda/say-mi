package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration structure
type Config struct {
	Fallback   string      `yaml:"fallback"`
	Categories interface{} `yaml:"categories"` // Can be nested map or flat map
}

// LoadConfig reads and parses the YAML config file from the given path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// FindCategory traverses the nested config using dot notation
// e.g., "female.sayako.email" navigates through categories["female"]["sayako"]["email"]
func (c *Config) FindCategory(path string) ([]string, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, false
	}

	current := c.Categories

	for i, part := range parts {
		// Try to navigate as a map
		asMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}

		value, exists := asMap[part]
		if !exists {
			return nil, false
		}

		// Check if this is the final part (should be a list of strings)
		if i == len(parts)-1 {
			return extractStringSlice(value)
		}

		// Continue navigating
		current = value
	}

	return nil, false
}

// extractStringSlice converts an interface{} to []string if possible
func extractStringSlice(value interface{}) ([]string, bool) {
	// Try as []interface{} (how YAML unmarshals arrays)
	if slice, ok := value.([]interface{}); ok {
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		if len(result) > 0 {
			return result, true
		}
	}

	// Try as []string directly
	if slice, ok := value.([]string); ok {
		return slice, true
	}

	return nil, false
}
