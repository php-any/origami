//go:build !windows

package utils

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

func StartCommandTree(cmd *exec.Cmd) (func(), error) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return func() { _ = cmd.Cancel() }, nil
}
