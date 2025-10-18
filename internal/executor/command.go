package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"interactive-feedback-mcp/internal/types"
)

type CommandExecutor struct {
	processes map[int]*types.CommandHandle
	mutex     sync.RWMutex
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		processes: make(map[int]*types.CommandHandle),
	}
}

func (ce *CommandExecutor) ExecuteCommand(command, workingDir string) (*types.CommandHandle, error) {
	ce.mutex.Lock()
	defer ce.mutex.Unlock()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	cmd.Dir = workingDir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Create process group for easier cleanup
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	handle := &types.CommandHandle{
		PID:       cmd.Process.Pid,
		StartTime: time.Now(),
		IsRunning: true,
		Output:    make(chan string, 100), // Buffered channel for performance
		Done:      make(chan error, 1),
	}

	ce.processes[handle.PID] = handle

	// Start goroutines for reading output
	go ce.readOutput(stdout.(*os.File), handle.Output, false)
	go ce.readOutput(stderr.(*os.File), handle.Output, true)
	go ce.waitForCompletion(cmd, handle)

	return handle, nil
}

func (ce *CommandExecutor) readOutput(pipe *os.File, output chan<- string, isError bool) {
	scanner := bufio.NewScanner(pipe)
	prefix := ""
	if isError {
		prefix = "[ERROR] "
	}

	for scanner.Scan() {
		line := prefix + scanner.Text() + "\n"
		select {
		case output <- line:
		default:
			// Channel is full, skip this line to prevent blocking
		}
	}
}

func (ce *CommandExecutor) waitForCompletion(cmd *exec.Cmd, handle *types.CommandHandle) {
	err := cmd.Wait()
	handle.IsRunning = false
	handle.Done <- err
	close(handle.Output)
	close(handle.Done)

	ce.mutex.Lock()
	delete(ce.processes, handle.PID)
	ce.mutex.Unlock()
}

func (ce *CommandExecutor) KillProcessTree(pid int) error {
	ce.mutex.Lock()
	defer ce.mutex.Unlock()

	handle, exists := ce.processes[pid]
	if !exists {
		return fmt.Errorf("process %d not found", pid)
	}

	// Kill the process group (includes all child processes)
	if runtime.GOOS != "windows" {
		if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
			// If SIGTERM fails, try SIGKILL
			syscall.Kill(-pid, syscall.SIGKILL)
		}
	} else {
		// Windows: Use taskkill to kill process tree
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
		cmd.Run()
	}

	// Also try to kill using gopsutil for additional cleanup
	if proc, err := process.NewProcess(int32(pid)); err == nil {
		children, _ := proc.Children()
		for _, child := range children {
			child.Kill()
		}
		proc.Kill()
	}

	handle.IsRunning = false
	return nil
}
