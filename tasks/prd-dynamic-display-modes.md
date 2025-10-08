# PRD: DisplayControl - Dynamic Display Control Menu System (Go Implementation)

## Introduction/Overview

DisplayControl reimplements the existing `display-modes.sh` script as a Go program (`dispctl`) that supports dynamic, multi-level menu generation based on YAML configuration files derived from display capabilities. Currently, the bash script has hardcoded menu categories and structure. The new Go implementation will parse YAML structures generated from `ddcutil capabilities` output and generate interactive TUI menus that allow users to navigate through VCP features to control their display via DDC/CI commands.

**Problem it solves:** Users need a flexible, automated way to control their monitors using DDC/CI without manually crafting commands or modifying scripts. By parsing display capabilities directly, the tool eliminates manual configuration and ensures compatibility with the monitor's actual features. Go provides better error handling, easier deployment (single binary), and more robust YAML parsing compared to the bash implementation.

## Goals

1. Reimplement display-modes.sh as a standalone Go binary named `dispctl`
2. Parse `ddcutil capabilities` output and generate YAML configuration automatically
3. Enable arbitrary nesting depth in YAML configuration files
4. Accept custom YAML file paths via command-line arguments
5. Provide clear error messages for malformed YAML files
6. Display command execution status (success/failure) to users
7. Implement dry-run mode for previewing commands without execution
8. Support error recovery by returning to menu after failed commands
9. Maintain TUI interface for better user experience
10. Provide a single, easily distributable binary with no runtime dependencies except ddcutil

## User Stories

1. **As a new user**, I want to automatically generate a configuration from my display's capabilities so that I don't need to manually research VCP codes.

2. **As a monitor power user**, I want to navigate through my display's VCP features in an organized menu so that I can quickly access and change settings.

3. **As a user with multiple displays**, I want to specify which display to use via command-line argument so that I can control different monitors.

4. **As a cautious user**, I want to preview what command will execute before it runs so that I can verify I'm making the right change.

5. **As a user troubleshooting display issues**, I want to see whether commands succeed or fail so that I can identify problems with my configuration.

6. **As a user creating a new YAML config**, I want clear error messages when my YAML is malformed so that I can fix it quickly.

## Functional Requirements

1. The program **must** be written in Go and compile to a single binary named `dispctl`.

2. The program **must** support a `--generate` or `-g` flag that:
   - Executes `ddcutil capabilities` (with optional `--display` flag)
   - Parses the output to extract VCP features, codes, and available values
   - Generates a YAML file in the format expected by the menu system
   - Saves to `$XDG_CONFIG_HOME/display-control/config.yaml` (or `~/.config/display-control/config.yaml` if `XDG_CONFIG_HOME` is not set) or a user-specified path

3. The program **must** support unlimited nesting depth in YAML structure, dynamically generating menus for each level.

4. The program **must** accept an optional command-line argument to specify the YAML file path (e.g., `./dispctl -f ~/.config/my-display.yaml`).

5. The program **must** support a `-f` or `--file` flag for specifying the config file path.

6. If no argument is provided, the program **must** default to `$XDG_CONFIG_HOME/display-control/config.yaml` (or `~/.config/display-control/config.yaml` if `XDG_CONFIG_HOME` is not set).

7. The program **must** validate YAML syntax and exit with a descriptive error message if the file is malformed.

8. The program **must** use an interactive TUI library to provide menu selection at each nesting level (replacing fzf functionality).

9. The program **must** parse YAML configuration derived from `ddcutil capabilities` output, where:
   - Branch nodes represent VCP features with their hex codes and metadata
   - Leaf nodes contain the hex values to send via `ddcutil setvcp`
   - When a leaf is selected, the program constructs and executes: `ddcutil setvcp <feature_code> <selected_value>`

10. The program **must** traverse the YAML tree until reaching a leaf node (terminal value containing the VCP value to send).

11. After command execution, the program **must** display:
    - The command that was executed (e.g., `ddcutil setvcp 0x10 0x64`)
    - Success or failure status
    - Command output/error output

12. The program **must** support a `--dry-run` flag that shows the selected command without executing it.

13. When a command fails, the program **must**:
    - Display the error message
    - Return to the top-level menu (allowing the user to try again)

14. The program **must** have only one external runtime dependency: `ddcutil` (system command).

15. The generated YAML **must** preserve:
    - VCP feature codes (e.g., `0x10` for Brightness)
    - Feature names/descriptions (e.g., "Brightness", "Display Mode")
    - Available values with their descriptions (e.g., "01: sRGB", "11: HDMI-1")

16. The YAML structure **must** use nested key-value pairs where:
    - Branch nodes are objects/maps containing VCP feature metadata and nested structures
    - Leaf nodes are hex values to be passed to `ddcutil setvcp <feature_code> <value>`

## Non-Goals (Out of Scope)

1. **Modifying the YAML file from within the program** - The program is read-only for configuration after generation.

2. **Supporting non-DDC/CI monitors** - Only monitors supported by ddcutil will work.

3. **Remote monitor control** - Only local monitor control via ddcutil.

4. **Multiple monitor support in a single invocation** - Users must use `--display` flag or generate separate configs. Full multi-monitor support is planned for a future iteration.

5. **Creating default config on first run** - Users must explicitly run `--generate` to create their configuration.

## Design Considerations

### TUI Interface
- Use a Go TUI library such as:
  - `github.com/charmbracelet/bubbletea` + `github.com/charmbracelet/bubbles` (recommended, modern approach)
  - `github.com/manifoldco/promptui` (simpler, fzf-like experience)
  - `github.com/AlecAivazis/survey/v2` (feature-rich prompts)
- Support keyboard navigation (arrow keys, Enter to select, Esc/q to quit)
- Display breadcrumb trail showing current menu path (optional enhancement)
- Use clear, human-readable option labels from YAML keys

### YAML Structure Example

Generated from `ddcutil capabilities`:

```yaml
brightness:
  _code: "0x10"
  _description: "Brightness"
  values:
    0: "0x00"
    50: "0x32"
    75: "0x4b"
    100: "0x64"

display_mode:
  _code: "0xdc"
  _description: "Display Mode"
  values:
    standard: "0x00"
    productivity: "0x01"
    games: "0x05"
    sports: "0x06"

input_source:
  _code: "0x60"
  _description: "Input Source"
  values:
    hdmi-1: "0x11"
    displayport-1: "0x0f"
    displayport-2: "0x10"

gamma:
  _code: "0x72"
  _description: "Gamma"
  values:
    1.80: "0x50"
    2.00: "0x64"
    2.20: "0x78"
    2.40: "0x8c"
```

Navigation flow: User selects `brightness` → `75` → program executes `ddcutil setvcp 0x10 0x4b`

### Error Messages
- **Malformed YAML**: "Error: Invalid YAML syntax in [file_path]. Please check the file format."
- **File not found**: "Error: Configuration file not found at [file_path]."
- **Command execution failure**: "Error: Command failed with exit code [N]: [error_output]"

## Technical Considerations

1. **Capabilities Parsing**:
   - Parse `ddcutil capabilities` text output using regex or text parsing
   - Extract VCP feature codes (e.g., `Feature: 10 (Brightness)`)
   - Extract available values for features with discrete values (e.g., `Values: 01: sRGB`)
   - Handle features without predefined values (continuous features like brightness 0-100)
   - Convert parsed data to YAML structure with `_code`, `_description`, and `values` fields

2. **YAML Parsing**:
   - Use `gopkg.in/yaml.v3` or `github.com/goccy/go-yaml` for YAML parsing
   - Parse YAML into custom structs that represent VCP features
   - Validate YAML syntax during unmarshaling
   - Extract `_code` metadata to construct ddcutil commands
   - Recursively determine if a node is a leaf (value) or branch (feature)

3. **Menu Generation**:
   - Recursive or iterative approach to handle arbitrary depth
   - Track current VCP feature being navigated
   - Extract value options from `values` map
   - Present human-readable labels via TUI library
   - Store the VCP code from `_code` field for command construction

4. **Command Construction**:
   - When leaf node is selected, extract VCP code from parent feature's `_code` field
   - Extract hex value from selected leaf node
   - Construct command: `ddcutil setvcp <_code> <selected_value>`
   - Type assert interface{} values to determine structure

5. **Command Execution**:
   - Use `os/exec` package to run constructed ddcutil commands
   - Capture stdout, stderr, and exit code
   - Handle ddcutil-specific errors and timeouts

6. **Error Handling**:
   - Return errors up the call stack with context
   - Display errors in TUI
   - Allow retry by returning to menu root
   - Validate ddcutil is available on startup

7. **Dry-run Implementation**:
   - Use flag package or `github.com/spf13/cobra` for CLI argument parsing
   - Pass dry-run flag through execution flow
   - Display command instead of executing when flag is set

8. **Build and Distribution**:
   - Use Go modules for dependency management
   - Provide Makefile or build script for easy compilation
   - Support static linking for portability

## Success Metrics

1. **Flexibility**: Users can create YAML configs with 3+ nesting levels without script modifications.

2. **Reliability**: Malformed YAML files produce helpful error messages 100% of the time.

3. **Usability**: Command execution feedback reduces user confusion about whether settings were applied.

4. **Safety**: Dry-run mode prevents accidental monitor changes during config testing.

## Open Questions

1. Should the program support command sequences (multiple ddcutil commands for one option)? Yes, once one command is issued and succeeds, the program should let the user select another command. The UI should provide a selection to move back to previous menu or exit

2. How should the menu return behavior work? Always return to root menu, or allow navigating back one level? Provide a way to move to the root and to move back 1 level

3. What should the breadcrumb trail look like if implemented? (e.g., `Display Mode > Games`) the example is OK

4. Should the program check for ddcutil availability on startup and provide a helpful error if missing? YES

5. Which TUI library should be used? (`bubbletea` for flexibility, `promptui` for simplicity, or other?) bubbletea

6. Should the program support configuration via environment variables (e.g., `DISPLAY_CONTROL_CONFIG`)? no

7. Should there be a `--version` flag to show the program version? Yes

8. How should continuous VCP features (brightness 0-100) be represented in the YAML? Generate all values, or support ranges? ranges

9. Should the `--generate` flag support a `--display` argument to specify which display to query? yes

10. Should generated YAML files include comments explaining the structure and VCP codes? no
