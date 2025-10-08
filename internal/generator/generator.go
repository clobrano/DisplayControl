package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/clobrano/DisplayControl/internal/parser"
	"gopkg.in/yaml.v3"
)

// CapabilitiesExecutor is a function type for executing ddcutil capabilities
type CapabilitiesExecutor func(displayID string) (string, error)

// Generate generates a YAML configuration file from ddcutil capabilities output
func Generate(displayID string, outputPath string) error {
	return GenerateWithExecutor(executeCapabilities, displayID, outputPath)
}

// GenerateWithExecutor generates a YAML config using a custom executor (for testing)
func GenerateWithExecutor(executor CapabilitiesExecutor, displayID string, outputPath string) error {
	// Execute ddcutil capabilities
	output, err := executor(displayID)
	if err != nil {
		return fmt.Errorf("failed to execute ddcutil capabilities: %w", err)
	}

	// Parse capabilities
	features, err := parser.ParseCapabilities(output)
	if err != nil {
		return fmt.Errorf("failed to parse capabilities: %w", err)
	}

	// Convert to YAML-friendly structure
	yamlData := make(map[string]interface{})
	for name, feature := range features {
		yamlData[name] = map[string]interface{}{
			"_code":        feature.Code,
			"_description": feature.Description,
			"values":       feature.Values,
		}
	}

	// Marshal to YAML
	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// executeCapabilities executes ddcutil capabilities command
func executeCapabilities(displayID string) (string, error) {
	args := buildCapabilitiesArgs(displayID)
	cmd := exec.Command("ddcutil", args...)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ddcutil command failed: %w", err)
	}

	return string(output), nil
}

// buildCapabilitiesArgs builds the ddcutil command arguments
func buildCapabilitiesArgs(displayID string) []string {
	args := []string{"capabilities"}
	if displayID != "" {
		args = append(args, "--display", displayID)
	}
	return args
}
