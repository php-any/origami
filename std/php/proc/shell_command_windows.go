//go:build windows

package proc

import (
	"context"
	"os/exec"
	"syscall"
)

func shellCommand(ctx context.Context, command string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "cmd.exe")
	// cmd.exe parses its /C tail itself; Go's normal argv quoting follows
	// CommandLineToArgvW and corrupts nested quotes and redirections here.
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: syscall.EscapeArg(cmd.Path) + " /D /S /C \"" + command + "\"", HideWindow: true}
	return cmd
}
