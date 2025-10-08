# Clarifying Questions for Display Modes Extension

Please answer the following questions inline to help me create a detailed PRD:

## 1. Dynamic Menu Structure
Currently the script has hardcoded categories (gamma-values, input-source, display-modes). How should the dynamic menu work?

a) Single-level menu: Flatten all options into one menu (no categories)
b) Two-level menu: Keep the current category → option structure
c) Multi-level menu: Support unlimited nesting depth based on YAML structure
d) Other (please describe):

**Your answer:** C

## 2. YAML File Location & Discovery
How should the script find and use YAML configuration files?

a) Always use `~/.philips-display.yaml` (current behavior)
b) Accept a command-line argument for custom YAML file path
c) Search multiple locations (e.g., `~/.config/display-modes/`, current directory)
d) Support multiple YAML files and let user choose which one to use
e) Other (please describe):

**Your answer:** B

## 3. YAML File Validation
What should happen if the YAML file is malformed or missing required fields?

a) Exit with error message
b) Show warning but continue with available options
c) Fall back to creating default config (current behavior)
d) Interactive validation with helpful error messages
e) Other (please describe):

**Your answer:** If malformed, exit with error

## 4. Command Execution
The current script uses `ddcutil` commands. Should this be:

a) Kept as-is (only support ddcutil commands in YAML)
b) Support any shell command (more flexible)
c) Support command templates with variables/parameters
d) Add command validation/safety checks before execution
e) Other (please describe):

**Your answer:** ddcutil is necessary to send the commands to the monitor, so it is a strong dependency

## 5. User Feedback
After executing a command, what feedback should the user receive?

a) Silent execution (current behavior)
b) Show the command being executed (current: echo only)
c) Show command + success/failure status
d) Show command + full output from ddcutil/executed command
e) Configurable verbosity level
f) Other (please describe):

**Your answer:** C

## 6. Additional Features
Which of these features should be included? (Select all that apply)

a) Command history/favorites (remember recently used options)
b) Dry-run mode (show command without executing)
c) Support for command aliases/shortcuts
d) Interactive YAML editor to add new options
e) Support for multiple monitors (bus selection for ddcutil)
f) None of the above
g) Other (please describe):

**Your answer:** B, E only in a next iteration

## 7. Error Handling
What should happen if a command fails to execute?

a) Silent failure (exit)
b) Show error and exit
c) Show error and return to menu
d) Retry mechanism
e) Other (please describe):

**Your answer:** C

## 8. Backwards Compatibility
Should the extended script:

a) Maintain full backwards compatibility with existing YAML format
b) Require migration to new YAML format
c) Support both old and new formats
d) Don't care about backwards compatibility

**Your answer:** D

## 9. Dependencies
The script currently uses: bash, yq, awk, fzf, ddcutil. Should we:

a) Keep all current dependencies
b) Reduce dependencies (specify which ones to remove):
c) Add new dependencies if needed (e.g., for validation)
d) Make some dependencies optional with fallbacks

**Your answer:** A

## 10. Non-Goals / Out of Scope
What should this feature explicitly NOT do? (Select all that apply)

a) Modify/write to the YAML file from the script
b) Support non-Philips monitors
c) Provide GUI interface
d) Support remote monitor control
e) Auto-detect monitor capabilities
f) Other (please describe):

**Your answer:** A, E, moreover I want a TUI interface

---

Please fill in your answers above and let me know when you're ready for me to generate the PRD!
