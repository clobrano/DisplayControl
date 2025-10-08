package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clobrano/DisplayControl/internal/config"
	"github.com/clobrano/DisplayControl/internal/ddcutil"
)

func TestNewModel(t *testing.T) {
	cfg := &config.Config{
		Features: map[string]config.VCPFeature{
			"brightness": {
				Code:        "0x10",
				Description: "Brightness",
				Values: map[string]string{
					"0":   "0x00",
					"100": "0x64",
				},
			},
		},
	}

	executor := ddcutil.NewExecutor("")
	model := NewModel(cfg, executor, true)

	if model.config == nil {
		t.Error("Config should not be nil")
	}

	if len(model.items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(model.items))
	}

	if model.cursor != 0 {
		t.Errorf("Cursor should start at 0, got %d", model.cursor)
	}
}

func TestModel_Navigation(t *testing.T) {
	cfg := &config.Config{
		Features: map[string]config.VCPFeature{
			"brightness": {
				Code:        "0x10",
				Description: "Brightness",
				Values: map[string]string{
					"0":   "0x00",
					"50":  "0x32",
					"100": "0x64",
				},
			},
			"contrast": {
				Code:        "0x12",
				Description: "Contrast",
				Values: map[string]string{
					"0":   "0x00",
					"100": "0x64",
				},
			},
		},
	}

	executor := ddcutil.NewExecutor("")
	model := NewModel(cfg, executor, true)

	// Test down navigation
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(Model)
	if model.cursor != 1 {
		t.Errorf("Cursor should be at 1 after down, got %d", model.cursor)
	}

	// Test up navigation
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(Model)
	if model.cursor != 0 {
		t.Errorf("Cursor should be at 0 after up, got %d", model.cursor)
	}
}

func TestModel_SelectionAndBack(t *testing.T) {
	cfg := &config.Config{
		Features: map[string]config.VCPFeature{
			"brightness": {
				Code:        "0x10",
				Description: "Brightness",
				Values: map[string]string{
					"0":   "0x00",
					"100": "0x64",
				},
			},
		},
	}

	executor := ddcutil.NewExecutor("")
	model := NewModel(cfg, executor, true)

	// Select brightness
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)

	if len(model.currentPath) != 1 {
		t.Errorf("Expected path length 1, got %d", len(model.currentPath))
	}

	if model.currentPath[0] != "brightness" {
		t.Errorf("Expected path 'brightness', got %s", model.currentPath[0])
	}

	// Go back
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	model = updated.(Model)

	if len(model.currentPath) != 0 {
		t.Errorf("Expected path length 0 after back, got %d", len(model.currentPath))
	}
}

func TestModel_View(t *testing.T) {
	cfg := &config.Config{
		Features: map[string]config.VCPFeature{
			"brightness": {
				Code:        "0x10",
				Description: "Brightness",
				Values: map[string]string{
					"0": "0x00",
				},
			},
		},
	}

	executor := ddcutil.NewExecutor("")
	model := NewModel(cfg, executor, true)

	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !contains(view, "Main Menu") {
		t.Error("View should contain 'Main Menu'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
