//go:build windows

package utils

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"
)

var resumeProcess = windows.NewLazySystemDLL("ntdll.dll").NewProc("NtResumeProcess")

// StartCommandTree contains descendants before the child executes any code.
// Starting suspended closes the race between Start and job assignment.
func StartCommandTree(cmd *exec.Cmd) (func(), error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, err
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
		windows.CloseHandle(job)
		return nil, err
	}
	var mutex sync.Mutex
	closed := false
	finish := func() {
		mutex.Lock()
		defer mutex.Unlock()
		if !closed {
			closed = true
			windows.CloseHandle(job)
		}
	}
	cmd.Cancel = func() error {
		mutex.Lock()
		defer mutex.Unlock()
		if closed {
			return os.ErrProcessDone
		}
		if err := windows.TerminateJobObject(job, 1); err != nil {
			return err
		}
		// Cancellation may precede assignment while Start is returning.
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return nil
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_SUSPENDED | windows.CREATE_NO_WINDOW
	if err := cmd.Start(); err != nil {
		finish()
		return nil, err
	}
	process, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE|windows.PROCESS_SUSPEND_RESUME, false, uint32(cmd.Process.Pid))
	if err == nil {
		err = windows.AssignProcessToJobObject(job, process)
		if err == nil {
			status, _, _ := resumeProcess.Call(uintptr(process))
			if int32(status) < 0 {
				err = fmt.Errorf("NtResumeProcess: NTSTATUS %#x", status)
			}
		}
		windows.CloseHandle(process)
	}
	if err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		finish()
		return nil, err
	}
	return finish, nil
}
