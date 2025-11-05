package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clobrano/DisplayControl/internal/config"
	"github.com/clobrano/DisplayControl/internal/ddcutil"
)

// smartSort sorts items intelligently: numerically if all items are numbers, alphabetically otherwise
func smartSort(items []string) {
	// Check if all items can be parsed as numbers
	allNumbers := true
	numbers := make([]float64, len(items))

	for i, item := range items {
		num, err := strconv.ParseFloat(item, 64)
		if err != nil {
			allNumbers = false
			break
		}
		numbers[i] = num
	}

	if allNumbers {
		// Sort numerically
		sort.Slice(items, func(i, j int) bool {
			numI, _ := strconv.ParseFloat(items[i], 64)
			numJ, _ := strconv.ParseFloat(items[j], 64)
			return numI < numJ
		})
	} else {
		// Sort alphabetically
		sort.Strings(items)
	}
}

// Model represents the application state for the TUI
type Model struct {
	config      *config.Config
	executor    *ddcutil.Executor
	dryRun      bool
	breadcrumbs []string
	currentPath []string
	cursor      int
	items       []string
	result      *ddcutil.ExecutionResult
	err         error
	quitting    bool
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config, executor *ddcutil.Executor, dryRun bool) Model {
	items := make([]string, 0, len(cfg.Features))
	for name := range cfg.Features {
		items = append(items, name)
	}
	smartSort(items)

	return Model{
		config:      cfg,
		executor:    executor,
		dryRun:      dryRun,
		breadcrumbs: []string{"Main Menu"},
		currentPath: []string{},
		items:       items,
		cursor:      0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter":
			return m.handleSelection()

		case "b", "esc":
			return m.goBack()

		case "r":
			return m.goToRoot()
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var s strings.Builder

	// Breadcrumb trail
	s.WriteString(fmt.Sprintf("📍 %s\n\n", strings.Join(m.breadcrumbs, " > ")))

	// Show result if there is one
	if m.result != nil {
		s.WriteString("┌─ Command Execution Result ─────────────\n")
		s.WriteString(fmt.Sprintf("│ Command: %s\n", m.result.Command))
		if m.result.Success {
			s.WriteString("│ Status: ✓ SUCCESS\n")
		} else {
			s.WriteString("│ Status: ✗ FAILED\n")
		}
		if m.result.Output != "" {
			s.WriteString(fmt.Sprintf("│ Output: %s\n", m.result.Output))
		}
		if m.result.Error != nil {
			s.WriteString(fmt.Sprintf("│ Error: %s\n", m.result.Error))
		}
		s.WriteString("└─────────────────────────────────────────\n\n")
	}

	// Show error if there is one
	if m.err != nil {
		s.WriteString(fmt.Sprintf("❌ Error: %s\n\n", m.err))
	}

	// Navigation options
	if len(m.currentPath) > 0 {
		s.WriteString("  ← [b] Back  ↑ [r] Root Menu\n\n")
	}

	// Menu items
	for i, item := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	// Help text
	s.WriteString("\n")
	s.WriteString("Navigation: ↑/↓ or j/k • Enter to select • q to quit\n")

	return s.String()
}

// handleSelection processes the user's selection
func (m Model) handleSelection() (Model, tea.Cmd) {
	if len(m.items) == 0 {
		return m, nil
	}

	selected := m.items[m.cursor]

	// Check if we're at the feature level or value level
	if len(m.currentPath) == 0 {
		// Selected a feature, navigate into it
		feature, exists := m.config.Features[selected]
		if !exists {
			m.err = fmt.Errorf("feature not found: %s", selected)
			return m, nil
		}

		// Build value items
		values := make([]string, 0, len(feature.Values))
		for name := range feature.Values {
			values = append(values, name)
		}
		smartSort(values)

		m.currentPath = append(m.currentPath, selected)
		m.breadcrumbs = append(m.breadcrumbs, feature.Description)
		m.items = values
		m.cursor = 0
		m.result = nil
		m.err = nil

	} else {
		// Selected a value, execute the command
		featureName := m.currentPath[0]
		feature := m.config.Features[featureName]
		hexValue := feature.Values[selected]

		// Execute the command
		result := m.executor.Execute(feature.Code, hexValue, m.dryRun)
		m.result = result
		m.err = nil
	}

	return m, nil
}

// goBack navigates back one level
func (m Model) goBack() (Model, tea.Cmd) {
	if len(m.currentPath) == 0 {
		return m, nil
	}

	// Remove last path element
	m.currentPath = m.currentPath[:len(m.currentPath)-1]
	m.breadcrumbs = m.breadcrumbs[:len(m.breadcrumbs)-1]

	// Rebuild items for the current level
	if len(m.currentPath) == 0 {
		// Back to root - show features
		items := make([]string, 0, len(m.config.Features))
		for name := range m.config.Features {
			items = append(items, name)
		}
		smartSort(items)
		m.items = items
	}

	m.cursor = 0
	m.result = nil
	m.err = nil

	return m, nil
}

// goToRoot navigates back to the root menu
func (m Model) goToRoot() (Model, tea.Cmd) {
	m.currentPath = []string{}
	m.breadcrumbs = []string{"Main Menu"}

	// Rebuild items for root level
	items := make([]string, 0, len(m.config.Features))
	for name := range m.config.Features {
		items = append(items, name)
	}
	smartSort(items)
	m.items = items
	m.cursor = 0
	m.result = nil
	m.err = nil

	return m, nil
}
