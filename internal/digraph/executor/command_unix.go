//go:build !windows
// +build !windows

package executor

import (
	"os"
	"os/exec"
	"syscall"
)

// SetupCommand configures Unix-specific command attributes
func SetupCommand(cmd *exec.Cmd) {
	setupCommand(cmd)
}

// setupCommand configures Unix-specific command attributes
func setupCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
}

// killProcessGroup kills the process group on Unix systems
func killProcessGroup(cmd *exec.Cmd, sig os.Signal) error {
	if cmd != nil && cmd.Process != nil {
		return syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
	}
	return nil
}

// killMultipleProcessGroups kills multiple process groups on Unix systems
func killMultipleProcessGroups(cmds map[string]*exec.Cmd, sig os.Signal) error {
	var lastErr error
	for _, cmd := range cmds {
		if cmd != nil && cmd.Process != nil {
			if err := syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal)); err != nil {
				lastErr = err
			}
		}
	}
	return lastErr
}
