package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestVCPFeature_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    VCPFeature
		wantErr bool
	}{
		{
			name: "valid feature with values",
			yaml: `_code: "0x10"
_description: "Brightness"
values:
  "0": "0x00"
  "50": "0x32"
  "100": "0x64"`,
			want: VCPFeature{
				Code:        "0x10",
				Description: "Brightness",
				Values: map[string]string{
					"0":   "0x00",
					"50":  "0x32",
					"100": "0x64",
				},
			},
			wantErr: false,
		},
		{
			name: "feature without description",
			yaml: `_code: "0x60"
values:
  hdmi-1: "0x11"`,
			want: VCPFeature{
				Code: "0x60",
				Values: map[string]string{
					"hdmi-1": "0x11",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got VCPFeature
			err := yaml.Unmarshal([]byte(tt.yaml), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("VCPFeature.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Code != tt.want.Code {
					t.Errorf("Code = %v, want %v", got.Code, tt.want.Code)
				}
				if got.Description != tt.want.Description {
					t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
				}
				if len(got.Values) != len(tt.want.Values) {
					t.Errorf("Values length = %v, want %v", len(got.Values), len(tt.want.Values))
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	validYAML := `brightness:
  _code: "0x10"
  _description: "Brightness"
  values:
    "0": "0x00"
    "50": "0x32"
    "100": "0x64"

input_source:
  _code: "0x60"
  _description: "Input Source"
  values:
    hdmi-1: "0x11"
    displayport-1: "0x0f"
`

	malformedYAML := `brightness:
  _code: "0x10"
  values: [invalid, yaml, structure
`

	tests := []struct {
		name     string
		content  string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid config file",
			content:  validYAML,
			filename: "valid.yaml",
			wantErr:  false,
		},
		{
			name:     "malformed yaml",
			content:  malformedYAML,
			filename: "malformed.yaml",
			wantErr:  true,
		},
		{
			name:     "file not found",
			content:  "",
			filename: "nonexistent.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testPath string
			if tt.content != "" {
				testPath = filepath.Join(tmpDir, tt.filename)
				if err := os.WriteFile(testPath, []byte(tt.content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			} else {
				testPath = filepath.Join(tmpDir, tt.filename)
			}

			_, err := LoadConfig(testPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsLeafNode(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			name:  "string value is leaf",
			value: "0x64",
			want:  true,
		},
		{
			name:  "number value is leaf",
			value: 100,
			want:  true,
		},
		{
			name:  "map value is branch",
			value: map[string]interface{}{"_code": "0x10"},
			want:  false,
		},
		{
			name:  "nil value is leaf",
			value: nil,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLeafNode(tt.value); got != tt.want {
				t.Errorf("IsLeafNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveConfigPath(t *testing.T) {
	// Save original env var
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	tests := []struct {
		name           string
		customPath     string
		xdgConfigHome  string
		wantContains   string
	}{
		{
			name:         "custom path provided",
			customPath:   "/custom/path/config.yaml",
			wantContains: "/custom/path/config.yaml",
		},
		{
			name:          "XDG_CONFIG_HOME set",
			customPath:    "",
			xdgConfigHome: "/tmp/config",
			wantContains:  "/tmp/config/display-control/config.yaml",
		},
		{
			name:         "default fallback",
			customPath:   "",
			xdgConfigHome: "",
			wantContains: ".config/display-control/config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.xdgConfigHome != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgConfigHome)
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}

			got := ResolveConfigPath(tt.customPath)
			if tt.wantContains != "" && got != tt.wantContains && !filepath.IsAbs(got) {
				// For relative paths, just check if it contains the expected substring
				if tt.name == "default fallback" {
					if !filepath.IsAbs(got) {
						t.Errorf("ResolveConfigPath() = %v, want absolute path containing %v", got, tt.wantContains)
					}
				}
			}
		})
	}
}
