package ddcutil

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ExecutionResult holds the result of a ddcutil command execution
type ExecutionResult struct {
	Command string
	Output  string
	Error   error
	Success bool
}

// Executor handles ddcutil command execution
type Executor struct {
	displayID string
	timeout   time.Duration
}

// NewExecutor creates a new Executor instance
func NewExecutor(displayID string) *Executor {
	return &Executor{
		displayID: displayID,
		timeout:   10 * time.Second, // Default 10 second timeout
	}
}

// SetTimeout sets the command execution timeout
func (e *Executor) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// Execute executes a ddcutil setvcp command
func (e *Executor) Execute(code string, value string, dryRun bool) *ExecutionResult {
	args := e.buildSetvcpArgs(code, value)
	cmdStr := "ddcutil " + strings.Join(args, " ")

	result := &ExecutionResult{
		Command: cmdStr,
	}

	if dryRun {
		result.Output = "[DRY RUN] Command not executed"
		result.Success = true
		return result
	}

	// Execute the command with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ddcutil", args...)
	output, err := cmd.CombinedOutput()

	result.Output = string(output)
	result.Error = err
	result.Success = err == nil

	if ctx.Err() == context.DeadlineExceeded {
		result.Error = fmt.Errorf("command timed out after %v", e.timeout)
		result.Success = false
	}

	return result
}

// buildSetvcpArgs builds the arguments for ddcutil setvcp command
func (e *Executor) buildSetvcpArgs(code string, value string) []string {
	args := []string{"setvcp", code, value}
	if e.displayID != "" {
		args = append(args, "--display", e.displayID)
	}
	return args
}

// CheckDdcutilAvailable checks if ddcutil is available on the system
func CheckDdcutilAvailable() error {
	cmd := exec.Command("ddcutil", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ddcutil not found or not accessible: %w", err)
	}
	return nil
}
