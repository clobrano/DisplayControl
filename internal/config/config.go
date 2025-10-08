package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// VCPFeature represents a VCP feature with its code, description, and possible values
type VCPFeature struct {
	Code        string            `yaml:"_code"`
	Description string            `yaml:"_description"`
	Values      map[string]string `yaml:"values"`
}

// Config represents the root configuration structure
type Config struct {
	Features map[string]VCPFeature
}

// LoadConfig reads and parses a YAML configuration file
func LoadConfig(filePath string) (*Config, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML into a generic map first
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Convert to Config structure
	config := &Config{
		Features: make(map[string]VCPFeature),
	}

	for name, value := range rawConfig {
		featureMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		feature := VCPFeature{
			Values: make(map[string]string),
		}

		// Extract _code
		if code, ok := featureMap["_code"].(string); ok {
			feature.Code = code
		}

		// Extract _description
		if desc, ok := featureMap["_description"].(string); ok {
			feature.Description = desc
		}

		// Extract values
		if values, ok := featureMap["values"].(map[string]interface{}); ok {
			for k, v := range values {
				if strVal, ok := v.(string); ok {
					feature.Values[k] = strVal
				}
			}
		}

		// Validate that feature has a code
		if feature.Code == "" {
			return nil, fmt.Errorf("feature '%s' missing required '_code' field", name)
		}

		config.Features[name] = feature
	}

	return config, nil
}

// IsLeafNode determines if a value is a terminal/leaf node (not a nested structure)
func IsLeafNode(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}:
		return false
	case map[interface{}]interface{}:
		return false
	default:
		return true
	}
}

// ResolveConfigPath resolves the configuration file path
// If customPath is provided, it returns that
// Otherwise, it checks XDG_CONFIG_HOME or falls back to ~/.config/display-control/config.yaml
func ResolveConfigPath(customPath string) string {
	if customPath != "" {
		return customPath
	}

	// Check XDG_CONFIG_HOME
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "display-control", "config.yaml")
	}

	// Fallback to ~/.config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home directory, use current directory
		return filepath.Join(".config", "display-control", "config.yaml")
	}

	return filepath.Join(homeDir, ".config", "display-control", "config.yaml")
}
