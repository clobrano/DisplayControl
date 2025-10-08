package generator

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-config.yaml")

	// Mock ddcutil capabilities output
	mockCapabilities := `MCCS version: 2.1
VCP Features:
   Feature: 10 (Brightness)
      Values: 00..64
   Feature: 60 (Input Source)
      Values:
         0F: DisplayPort-1
         11: HDMI-1
`

	// Create generator with mock executor
	mockExec := func(displayID string) (string, error) {
		return mockCapabilities, nil
	}

	err := GenerateWithExecutor(mockExec, "", outputPath)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Read and parse the generated YAML
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated config: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse generated YAML: %v", err)
	}

	// Verify brightness feature exists
	brightness, ok := config["brightness"].(map[string]interface{})
	if !ok {
		t.Fatal("brightness feature not found in config")
	}

	if code := brightness["_code"]; code != "0x10" {
		t.Errorf("brightness _code = %v, want 0x10", code)
	}

	// Verify input_source feature exists
	inputSource, ok := config["input_source"].(map[string]interface{})
	if !ok {
		t.Fatal("input_source feature not found in config")
	}

	if code := inputSource["_code"]; code != "0x60" {
		t.Errorf("input_source _code = %v, want 0x60", code)
	}
}

func TestGenerate_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")

	mockCapabilities := `VCP Features:
   Feature: 10 (Brightness)
      Values: 00..64
`

	mockExec := func(displayID string) (string, error) {
		return mockCapabilities, nil
	}

	err := GenerateWithExecutor(mockExec, "", outputPath)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify file and directories were created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created in nested directory")
	}
}

func TestExecuteCapabilities(t *testing.T) {
	tests := []struct {
		name      string
		displayID string
		wantArgs  []string
	}{
		{
			name:      "no display specified",
			displayID: "",
			wantArgs:  []string{"capabilities"},
		},
		{
			name:      "display specified",
			displayID: "1",
			wantArgs:  []string{"capabilities", "--display", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildCapabilitiesArgs(tt.displayID)
			if len(args) != len(tt.wantArgs) {
				t.Errorf("args length = %v, want %v", len(args), len(tt.wantArgs))
				return
			}
			for i, arg := range args {
				if arg != tt.wantArgs[i] {
					t.Errorf("args[%d] = %v, want %v", i, arg, tt.wantArgs[i])
				}
			}
		})
	}
}
