## Relevant Files

- `DisplayControl/cmd/dispctl/main.go` - Entry point for the dispctl CLI application
- `DisplayControl/internal/config/config.go` - YAML configuration parsing and validation logic
- `DisplayControl/internal/config/config_test.go` - Unit tests for config package
- `DisplayControl/internal/parser/capabilities.go` - Parser for ddcutil capabilities output
- `DisplayControl/internal/parser/capabilities_test.go` - Unit tests for capabilities parser
- `DisplayControl/internal/ui/menu.go` - Bubbletea TUI implementation for menu navigation
- `DisplayControl/internal/ui/menu_test.go` - Unit tests for TUI menu logic
- `DisplayControl/internal/ddcutil/executor.go` - DDC/CI command execution wrapper
- `DisplayControl/internal/ddcutil/executor_test.go` - Unit tests for command executor
- `DisplayControl/internal/generator/generator.go` - YAML configuration file generator from capabilities
- `DisplayControl/internal/generator/generator_test.go` - Unit tests for generator
- `DisplayControl/go.mod` - Go module dependencies
- `DisplayControl/go.sum` - Go module checksums
- `DisplayControl/Makefile` - Build automation scripts
- `DisplayControl/README.md` - Project documentation and usage instructions

### Notes

- This is a new Go project that will be created in the `DisplayControl/` directory within the current workspace.
- The Go module path will be `github.com/clobrano/DisplayControl`.
- The current bash implementation (`display-modes.sh`) uses `yq` and `fzf` for YAML parsing and menu selection.
- The Go implementation will be self-contained and replace the bash script.
- Unit tests should be placed alongside code files (e.g., `config.go` and `config_test.go` in the same directory).
- Use `go test ./...` to run all tests from within the `DisplayControl/` directory.
- Use `go build -o dispctl ./cmd/dispctl` to build the binary.

## Tasks

- [x] 1.0 Project Setup and Infrastructure
  - [x] 1.1 Create `DisplayControl/` directory in current workspace
  - [x] 1.2 Initialize Go module with `go mod init github.com/clobrano/DisplayControl` inside `DisplayControl/` directory
  - [x] 1.3 Create directory structure: `DisplayControl/cmd/dispctl/`, `DisplayControl/internal/config/`, `DisplayControl/internal/parser/`, `DisplayControl/internal/ui/`, `DisplayControl/internal/ddcutil/`, `DisplayControl/internal/generator/`
  - [x] 1.4 Add initial dependencies: `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/bubbles`, `gopkg.in/yaml.v3`, `github.com/spf13/cobra` or `flag` package
  - [x] 1.5 Create `.gitignore` file in `DisplayControl/` (ignore `dispctl` binary, test coverage files, IDE files)
  - [x] 1.6 Create basic `Makefile` in `DisplayControl/` with targets: `build`, `test`, `clean`, `install`
  - [x] 1.7 Create version constant or version file for `--version` flag support

- [x] 2.0 YAML Configuration Parsing and Validation
  - [x] 2.1 Define Go structs for YAML structure (`VCPFeature`, `ConfigRoot`) in `internal/config/config.go`
  - [x] 2.2 Implement struct tags for YAML unmarshaling with `_code`, `_description`, and `values` fields
  - [x] 2.3 Implement `LoadConfig(filePath string) (*ConfigRoot, error)` function to read and parse YAML files
  - [x] 2.4 Implement YAML validation logic to check for required fields (`_code` for features)
  - [x] 2.5 Implement helper function to determine if a node is a leaf (terminal value) or branch (nested feature)
  - [x] 2.6 Add error handling for file not found, invalid YAML syntax, and malformed structure
  - [x] 2.7 Write unit tests for valid YAML files, malformed YAML, missing required fields, deeply nested structures
  - [x] 2.8 Implement config file path resolution: check XDG_CONFIG_HOME, fallback to ~/.config/display-control/config.yaml

- [x] 3.0 DDC/CI Capabilities Parser
  - [x] 3.1 Create `ParseCapabilities(output string) (map[string]VCPFeature, error)` function in `internal/parser/capabilities.go`
  - [x] 3.2 Implement regex or text parsing to extract VCP feature codes (e.g., `Feature: 10 (Brightness)`)
  - [x] 3.3 Implement parsing for feature names/descriptions from capabilities output
  - [x] 3.4 Implement parsing for discrete values (e.g., `Values: 01: sRGB, 05: Games`)
  - [x] 3.5 Handle continuous features (brightness 0-100) by detecting value ranges and generating appropriate value sets
  - [x] 3.6 Convert hex feature codes to proper format (e.g., `10` to `0x10`)
  - [x] 3.7 Handle edge cases: features without values, unparseable output, empty capabilities
  - [x] 3.8 Write unit tests with sample ddcutil capabilities output, edge cases, and malformed input

- [x] 4.0 YAML Configuration Generator
  - [x] 4.1 Create `Generate(displayID string, outputPath string) error` function in `internal/generator/generator.go`
  - [x] 4.2 Implement `ddcutil capabilities` command execution with optional `--display` flag support
  - [x] 4.3 Integrate capabilities parser to convert ddcutil output to internal data structures
  - [x] 4.4 Implement YAML serialization of parsed VCP features using `gopkg.in/yaml.v3`
  - [x] 4.5 Ensure generated YAML matches expected structure with `_code`, `_description`, and `values` fields
  - [x] 4.6 Create output directory if it doesn't exist (e.g., `~/.config/display-control/`)
  - [x] 4.7 Handle file write permissions and errors gracefully
  - [x] 4.8 Write unit tests for generator with mocked ddcutil output and file system operations

- [x] 5.0 DDC/CI Command Executor
  - [x] 5.1 Create `Execute(code string, value string, dryRun bool) (string, error)` function in `internal/ddcutil/executor.go`
  - [x] 5.2 Implement command construction: `ddcutil setvcp <code> <value>` with optional `--display` flag
  - [x] 5.3 Use `os/exec` to execute ddcutil commands and capture stdout, stderr, and exit code
  - [x] 5.4 Implement dry-run mode that returns the command string without execution
  - [x] 5.5 Add command timeout handling (e.g., 10 second timeout for ddcutil operations)
  - [x] 5.6 Parse and return command output and error messages in structured format
  - [x] 5.7 Implement `CheckDdcutilAvailable() error` function to verify ddcutil is installed and accessible
  - [x] 5.8 Write unit tests with mocked command execution, dry-run validation, timeout scenarios, and error cases

- [x] 6.0 Bubbletea TUI Menu System
  - [x] 6.1 Create bubbletea model struct in `internal/ui/menu.go` with state: current menu items, breadcrumb trail, selected VCP feature
  - [x] 6.2 Implement `Init()` method to initialize the bubbletea model
  - [x] 6.3 Implement `Update(msg tea.Msg)` method to handle keyboard input (arrow keys, Enter, Esc, 'q', 'b' for back, 'r' for root)
  - [x] 6.4 Implement `View()` method to render current menu items with breadcrumb trail
  - [x] 6.5 Implement menu navigation logic: track current depth, parent nodes, and navigation history
  - [x] 6.6 Implement breadcrumb trail display showing current path (e.g., "Display Mode > Games")
  - [x] 6.7 Add "← Back" and "↑ Root Menu" options at the top of each submenu
  - [x] 6.8 Implement leaf node detection and command execution trigger when terminal value is selected
  - [x] 6.9 Display command execution results (success/failure, command output) in TUI after execution
  - [x] 6.10 Implement error display and automatic return to menu after failed commands
  - [x] 6.11 Add quit confirmation or direct quit on 'q' or Esc key
  - [x] 6.12 Write unit tests for state transitions, navigation logic, and edge cases

- [x] 7.0 CLI Application and Command-Line Interface
  - [x] 7.1 Create `main.go` in `cmd/dispctl/` with basic CLI structure
  - [x] 7.2 Implement `--version` / `-v` flag to display program version
  - [x] 7.3 Implement `--generate` / `-g` flag to trigger config generation mode
  - [x] 7.4 Implement `--file` / `-f` flag to specify custom YAML config file path
  - [x] 7.5 Implement `--dry-run` / `-d` flag for command preview mode
  - [x] 7.6 Implement `--display` flag to specify target display for both generation and execution
  - [x] 7.7 Add startup check for ddcutil availability using `CheckDdcutilAvailable()`
  - [x] 7.8 Implement config file path resolution logic (XDG_CONFIG_HOME, default fallback)
  - [x] 7.9 Wire together: config loading → TUI initialization → command execution → result display
  - [x] 7.10 Implement graceful error handling with user-friendly error messages at top level
  - [x] 7.11 Add help text and usage examples via `--help` flag

- [x] 8.0 Integration and End-to-End Testing
  - [x] 8.1 Create integration test directory and sample YAML configurations for testing
  - [x] 8.2 Write end-to-end test for config generation workflow (mocked ddcutil capabilities)
  - [x] 8.3 Write end-to-end test for menu navigation and command execution (mocked ddcutil setvcp)
  - [x] 8.4 Test error scenarios: missing config file, malformed YAML, ddcutil not found
  - [x] 8.5 Test dry-run mode end-to-end to ensure no commands execute
  - [x] 8.6 Test navigation: deep nesting, back navigation, root menu return
  - [x] 8.7 Test with multi-level nested YAML structures (3+ levels deep)
  - [x] 8.8 Perform manual testing with actual ddcutil on real hardware (if available)

- [x] 9.0 Documentation and Distribution
  - [x] 9.1 Write comprehensive README.md with project description, features, and prerequisites
  - [x] 9.2 Document installation instructions (build from source, binary releases)
  - [x] 9.3 Document usage examples: generating config, running menu, dry-run mode, multi-display setup
  - [x] 9.4 Document YAML structure format with examples and explanation of `_code`, `_description`, `values`
  - [x] 9.5 Document keyboard shortcuts and TUI navigation in README
  - [x] 9.6 Add troubleshooting section for common issues (ddcutil permissions, I2C device access)
  - [x] 9.7 Create example YAML configuration files in `examples/` directory
  - [x] 9.8 Add build instructions for static linking (for portability)
  - [x] 9.9 Consider adding GitHub Actions or CI configuration for automated builds and releases
  - [x] 9.10 Update Makefile with `install` target to copy binary to `/usr/local/bin` or `~/.local/bin`
