package ddcutil

import (
	"strings"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		value     string
		displayID string
		dryRun    bool
		wantCmd   string
	}{
		{
			name:    "basic setvcp command",
			code:    "0x10",
			value:   "0x64",
			dryRun:  false,
			wantCmd: "ddcutil setvcp 0x10 0x64",
		},
		{
			name:      "setvcp with display ID",
			code:      "0x60",
			value:     "0x11",
			displayID: "1",
			dryRun:    false,
			wantCmd:   "ddcutil setvcp 0x60 0x11 --display 1",
		},
		{
			name:    "dry-run mode",
			code:    "0x10",
			value:   "0x32",
			dryRun:  true,
			wantCmd: "ddcutil setvcp 0x10 0x32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor(tt.displayID)
			result := executor.Execute(tt.code, tt.value, tt.dryRun)

			if !strings.Contains(result.Command, "setvcp") {
				t.Errorf("Command doesn't contain 'setvcp': %v", result.Command)
			}

			if !strings.Contains(result.Command, tt.code) {
				t.Errorf("Command doesn't contain code %v: %v", tt.code, result.Command)
			}

			if !strings.Contains(result.Command, tt.value) {
				t.Errorf("Command doesn't contain value %v: %v", tt.value, result.Command)
			}

			if tt.displayID != "" && !strings.Contains(result.Command, tt.displayID) {
				t.Errorf("Command doesn't contain display ID %v: %v", tt.displayID, result.Command)
			}

			if tt.dryRun && result.Output != "[DRY RUN] Command not executed" {
				t.Errorf("Dry-run should not execute command, got output: %v", result.Output)
			}
		})
	}
}

func TestCheckDdcutilAvailable(t *testing.T) {
	// This test will pass if ddcutil is installed, otherwise skip
	err := CheckDdcutilAvailable()
	if err != nil {
		t.Skipf("ddcutil not available: %v", err)
	}
}

func TestBuildSetvcpArgs(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		value     string
		displayID string
		want      []string
	}{
		{
			name:  "basic command",
			code:  "0x10",
			value: "0x64",
			want:  []string{"setvcp", "0x10", "0x64"},
		},
		{
			name:      "with display ID",
			code:      "0x60",
			value:     "0x11",
			displayID: "1",
			want:      []string{"setvcp", "0x60", "0x11", "--display", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor(tt.displayID)
			args := executor.buildSetvcpArgs(tt.code, tt.value)

			if len(args) != len(tt.want) {
				t.Errorf("args length = %v, want %v", len(args), len(tt.want))
				return
			}

			for i, arg := range args {
				if arg != tt.want[i] {
					t.Errorf("args[%d] = %v, want %v", i, arg, tt.want[i])
				}
			}
		})
	}
}

func TestExecutor_Timeout(t *testing.T) {
	executor := NewExecutor("")
	executor.SetTimeout(1 * time.Millisecond) // Very short timeout

	// This should timeout since ddcutil operations take longer
	// Only run if ddcutil is available
	if CheckDdcutilAvailable() != nil {
		t.Skip("ddcutil not available")
	}

	// Note: Actual timeout test would require a real ddcutil command
	// which we can't reliably test without hardware
	if executor.timeout != 1*time.Millisecond {
		t.Errorf("Timeout = %v, want %v", executor.timeout, 1*time.Millisecond)
	}
}
