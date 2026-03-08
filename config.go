package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration structure
type Config struct {
	Fallback string      `yaml:"fallback"`
	Voices   interface{} `yaml:"voices"` // Can be nested map or flat map
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
// e.g., "kyoko.greeting" navigates through voices["kyoko"]["greeting"]
func (c *Config) FindCategory(path string) ([]string, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, false
	}

	current := c.Voices
	DebugLog("FindCategory: starting with path '%s', parts=%v", path, parts)

	for i, part := range parts {
		// Try to navigate as a map
		asMap, ok := current.(map[string]interface{})
		if !ok {
			DebugLog("FindCategory: part '%s' is not a map, type=%T", part, current)
			return nil, false
		}

		value, exists := asMap[part]
		if !exists {
			DebugLog("FindCategory: part '%s' not found in map, available keys: %v", part, getMapKeys(asMap))
			return nil, false
		}

		DebugLog("FindCategory: part '%s' found, value type=%T", part, value)

		// Check if this is the final part (should be a list of strings)
		if i == len(parts)-1 {
			result, ok := extractStringSlice(value)
			DebugLog("FindCategory: extractStringSlice returned %d items, ok=%v", len(result), ok)
			return result, ok
		}

		// Continue navigating
		current = value
	}

	return nil, false
}

// getMapKeys returns the keys of a map for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// extractStringSlice converts an interface{} to []string if possible
// Handles: []string, []interface{} of strings, or []interface{} of objects with "audio" field
func extractStringSlice(value interface{}) ([]string, bool) {
	// Try as []interface{} (how YAML unmarshals arrays)
	if slice, ok := value.([]interface{}); ok {
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			// Try as string
			if str, ok := item.(string); ok {
				result = append(result, str)
				continue
			}
			// Try as map with "audio" field (for ai_voices.yaml structure)
			if m, ok := item.(map[string]interface{}); ok {
				if audio, exists := m["audio"]; exists {
					if audioStr, ok := audio.(string); ok && audioStr != "" {
						result = append(result, audioStr)
					}
				}
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
