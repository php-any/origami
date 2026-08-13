package proc

import (
	"os/exec"
	"sync"
)

// ProcessInfo 存储进程信息
type ProcessInfo struct {
	Cmd      *exec.Cmd
	Command  string
	Pid      int
	Running  bool
	ExitCode int
	mutex    sync.RWMutex
}

// NewProcessInfo 创建进程信息
func NewProcessInfo(cmd *exec.Cmd, command string) *ProcessInfo {
	return &ProcessInfo{
		Cmd:      cmd,
		Command:  command,
		Pid:      cmd.Process.Pid,
		Running:  true,
		ExitCode: -1,
	}
}

// SetRunning 设置运行状态
func (p *ProcessInfo) SetRunning(running bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Running = running
}

// SetExitCode 设置退出码
func (p *ProcessInfo) SetExitCode(exitCode int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ExitCode = exitCode
}

// GetRunning 获取运行状态
func (p *ProcessInfo) GetRunning() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Running
}

// GetExitCode 获取退出码
func (p *ProcessInfo) GetExitCode() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.ExitCode
}
