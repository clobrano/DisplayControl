package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clobrano/DisplayControl/internal/config"
	"github.com/clobrano/DisplayControl/internal/ddcutil"
	"github.com/clobrano/DisplayControl/internal/generator"
	"github.com/clobrano/DisplayControl/internal/ui"
	"github.com/clobrano/DisplayControl/internal/version"
)

func main() {
	// Define flags
	generateFlag := flag.Bool("generate", false, "Generate configuration from ddcutil capabilities")
	generateShort := flag.Bool("g", false, "Generate configuration (shorthand)")
	versionFlag := flag.Bool("version", false, "Display version information")
	versionShort := flag.Bool("v", false, "Display version (shorthand)")
	fileFlag := flag.String("file", "", "Path to configuration file")
	fileShort := flag.String("f", "", "Path to configuration file (shorthand)")
	dryRunFlag := flag.Bool("dry-run", false, "Preview commands without executing")
	dryRunShort := flag.Bool("d", false, "Dry-run mode (shorthand)")
	displayFlag := flag.String("display", "", "Display ID for ddcutil commands")

	flag.Parse()

	// Handle version flag
	if *versionFlag || *versionShort {
		fmt.Printf("dispctl version %s\n", version.Version)
		os.Exit(0)
	}

	// Check if ddcutil is available
	if err := ddcutil.CheckDdcutilAvailable(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please ensure ddcutil is installed and accessible.\n")
		os.Exit(1)
	}

	// Handle generate flag
	if *generateFlag || *generateShort {
		configPath := *fileFlag
		if configPath == "" {
			configPath = *fileShort
		}
		if configPath == "" {
			configPath = config.ResolveConfigPath("")
		}

		fmt.Printf("Generating configuration from ddcutil capabilities...\n")
		if err := generator.Generate(*displayFlag, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Configuration generated successfully: %s\n", configPath)
		os.Exit(0)
	}

	// Resolve config file path
	configPath := *fileFlag
	if configPath == "" {
		configPath = *fileShort
	}
	configPath = config.ResolveConfigPath(configPath)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nHint: Generate a configuration with: dispctl --generate\n")
		os.Exit(1)
	}

	// Check if config has features
	if len(cfg.Features) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No features found in configuration file\n")
		fmt.Fprintf(os.Stderr, "Configuration file: %s\n", configPath)
		os.Exit(1)
	}

	// Determine dry-run mode
	dryRun := *dryRunFlag || *dryRunShort

	// Create executor
	executor := ddcutil.NewExecutor(*displayFlag)

	// Create TUI model
	model := ui.NewModel(cfg, executor, dryRun)

	// Run the TUI
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
