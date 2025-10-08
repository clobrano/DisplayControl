package parser

import (
	"testing"
)

func TestParseCapabilities(t *testing.T) {
	sampleOutput := `MCCS version: 2.1
VCP Features:
   Feature: 10 (Brightness)
      Values: 00..64
   Feature: 12 (Contrast)
      Values: 00..64
   Feature: 14 (Select color preset)
      Values:
         01: sRGB
         04: 5000K
         05: 6500K
         06: 7500K
         08: 9300K
         0B: User 1
   Feature: 60 (Input Source)
      Values:
         0F: DisplayPort-1
         10: DisplayPort-2
         11: HDMI-1
   Feature: DC (Display Mode)
      Values:
         00: Standard/Default mode
         03: Movie
         0C: User
   Feature: E0 (Manufacturer Specific)
   Feature: E1 (Power control)
      Values: 00..04
`

	features, err := ParseCapabilities(sampleOutput)
	if err != nil {
		t.Fatalf("ParseCapabilities() error = %v", err)
	}

	tests := []struct {
		name         string
		featureName  string
		wantCode     string
		wantDesc     string
		wantValues   int
		checkValue   string
		checkValueHex string
	}{
		{
			name:         "brightness continuous range",
			featureName:  "brightness",
			wantCode:     "0x10",
			wantDesc:     "Brightness",
			wantValues:   5, // 0, 25, 50, 75, 100
			checkValue:   "0",
			checkValueHex: "0x00",
		},
		{
			name:         "color preset discrete values",
			featureName:  "color_preset",
			wantCode:     "0x14", // Feature: 14 in ddcutil output (hex 14)
			wantDesc:     "Select color preset",
			wantValues:   6,
			checkValue:   "sRGB",
			checkValueHex: "0x01",
		},
		{
			name:         "input source discrete values",
			featureName:  "input_source",
			wantCode:     "0x60",
			wantDesc:     "Input Source",
			wantValues:   3,
			checkValue:   "DisplayPort-1",
			checkValueHex: "0x0f",
		},
		{
			name:         "display mode hex code",
			featureName:  "display_mode",
			wantCode:     "0xdc",
			wantDesc:     "Display Mode",
			wantValues:   3,
			checkValue:   "Standard/Default mode",
			checkValueHex: "0x00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feature, ok := features[tt.featureName]
			if !ok {
				t.Errorf("Feature %s not found", tt.featureName)
				return
			}

			if feature.Code != tt.wantCode {
				t.Errorf("Code = %v, want %v", feature.Code, tt.wantCode)
			}

			if feature.Description != tt.wantDesc {
				t.Errorf("Description = %v, want %v", feature.Description, tt.wantDesc)
			}

			if len(feature.Values) != tt.wantValues {
				t.Errorf("Values length = %v, want %v", len(feature.Values), tt.wantValues)
			}

			if tt.checkValue != "" {
				hexVal, ok := feature.Values[tt.checkValue]
				if !ok {
					t.Errorf("Value %s not found in feature values", tt.checkValue)
					return
				}
				if hexVal != tt.checkValueHex {
					t.Errorf("Value %s = %v, want %v", tt.checkValue, hexVal, tt.checkValueHex)
				}
			}
		})
	}
}

func TestParseCapabilities_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantLen int
	}{
		{
			name:    "empty output",
			input:   "",
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "no features",
			input:   "MCCS version: 2.1\nVCP Features:\n",
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "feature without values",
			input: `VCP Features:
   Feature: E0 (Manufacturer Specific)
`,
			wantErr: false,
			wantLen: 0, // Features without values are skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			features, err := ParseCapabilities(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCapabilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(features) != tt.wantLen {
				t.Errorf("Features length = %v, want %v", len(features), tt.wantLen)
			}
		})
	}
}

func TestDecimalToHex(t *testing.T) {
	tests := []struct {
		hexInput string
		want     string
	}{
		{"10", "0x10"},    // hex 10 = decimal 16
		{"16", "0x16"},    // hex 16 = decimal 22
		{"DC", "0xdc"},    // hex DC = decimal 220
		{"E0", "0xe0"},    // hex E0 = decimal 224
		{"60", "0x60"},    // hex 60 = decimal 96
		{"0A", "0x0a"},    // hex 0A = decimal 10
	}

	for _, tt := range tests {
		t.Run(tt.hexInput, func(t *testing.T) {
			got := decimalToHex(tt.hexInput)
			if got != tt.want {
				t.Errorf("decimalToHex(%s) = %v, want %v", tt.hexInput, got, tt.want)
			}
		})
	}
}

func TestSanitizeFeatureName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Brightness", "brightness"},
		{"Select color preset", "color_preset"},
		{"Input Source", "input_source"},
		{"Display Mode", "display_mode"},
		{"Manufacturer Specific", "manufacturer_specific"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeFeatureName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeFeatureName(%s) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
