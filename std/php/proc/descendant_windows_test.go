//go:build windows

package proc

import (
	"golang.org/x/sys/windows"
	"testing"
)

func watchDescendant(t *testing.T, pid int) func() bool {
	process, err := windows.OpenProcess(windows.SYNCHRONIZE, false, uint32(pid))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { windows.CloseHandle(process) })
	return func() bool {
		status, err := windows.WaitForSingleObject(process, 3000)
		return err == nil && status == windows.WAIT_OBJECT_0
	}
}
