package executor

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"interactive-feedback-mcp/internal/types"
)

func TestNewCommandExecutor(t *testing.T) {
	executor := NewCommandExecutor()
	assert.NotNil(t, executor)
	assert.NotNil(t, executor.processes)
}

func TestCommandExecutor_ExecuteCommand(t *testing.T) {
	executor := NewCommandExecutor()

	// Test successful command execution
	var command string
	if runtime.GOOS == "windows" {
		command = "echo Hello World"
	} else {
		command = "echo 'Hello World'"
	}

	handle, err := executor.ExecuteCommand(command, ".")
	require.NoError(t, err)
	assert.NotNil(t, handle)
	assert.True(t, handle.IsRunning)
	assert.NotNil(t, handle.Output)
	assert.NotNil(t, handle.Done)

	// Read output
	var output strings.Builder
	timeout := time.After(5 * time.Second)

	for {
		select {
		case line := <-handle.Output:
			output.WriteString(line)
			if strings.Contains(line, "Hello World") {
				goto done
			}
		case <-timeout:
			t.Fatal("Timeout waiting for command output")
		}
	}

done:
	assert.Contains(t, output.String(), "Hello World")

	// Wait for completion
	select {
	case err := <-handle.Done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for command completion")
	}

	assert.False(t, handle.IsRunning)
}

func TestCommandExecutor_ExecuteCommand_InvalidCommand(t *testing.T) {
	executor := NewCommandExecutor()

	// Test invalid command
	handle, err := executor.ExecuteCommand("nonexistentcommand12345", ".")
	if err != nil {
		// Some systems may return error immediately
		assert.Error(t, err)
		return
	}

	// If command starts, it should fail
	select {
	case err := <-handle.Done:
		assert.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for command failure")
	}
}

func TestCommandExecutor_KillProcessTree(t *testing.T) {
	executor := NewCommandExecutor()

	// Start a long-running command
	var command string
	if runtime.GOOS == "windows" {
		command = "ping 127.0.0.1 -n 10"
	} else {
		command = "sleep 10"
	}

	handle, err := executor.ExecuteCommand(command, ".")
	require.NoError(t, err)

	// Kill the process
	err = executor.KillProcessTree(handle.PID)
	assert.NoError(t, err)

	// Wait for process to be killed
	select {
	case err := <-handle.Done:
		assert.Error(t, err) // Should be killed
	case <-time.After(5 * time.Second):
		t.Fatal("Process was not killed")
	}

	assert.False(t, handle.IsRunning)
}

func TestCommandExecutor_ConcurrentCommands(t *testing.T) {
	executor := NewCommandExecutor()

	// Start multiple commands concurrently
	var command string
	if runtime.GOOS == "windows" {
		command = "echo test"
	} else {
		command = "echo 'test'"
	}

	handles := make([]*types.CommandHandle, 3)
	for i := 0; i < 3; i++ {
		handle, err := executor.ExecuteCommand(command, ".")
		require.NoError(t, err)
		handles[i] = handle
	}

	// Wait for all commands to complete
	for _, handle := range handles {
		select {
		case err := <-handle.Done:
			assert.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for command completion")
		}
	}
}
