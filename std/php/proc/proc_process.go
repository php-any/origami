package proc

import (
	"context"
	"io"
	"os/exec"
	"sync"
)

// ProcessInfo 存储进程信息
type ProcessInfo struct {
	Cmd        *exec.Cmd
	Command    string
	Pid        int
	Running    bool
	ExitCode   int
	mutex      sync.RWMutex
	done       chan struct{} // 进程结束后关闭，用于 proc_close 阻塞等待
	stopCancel func() bool
	closed     bool
	pipes      []io.Closer
	doneOnce   sync.Once
}

// NewProcessInfo 创建进程信息
func NewProcessInfo(cmd *exec.Cmd, command string) *ProcessInfo {
	return &ProcessInfo{
		Cmd:      cmd,
		Command:  command,
		Pid:      cmd.Process.Pid,
		Running:  true,
		ExitCode: -1,
		done:     make(chan struct{}),
	}
}

// WaitDone 阻塞等待进程结束（proc_close 使用）
func (p *ProcessInfo) WaitDone() {
	<-p.done
}

func (p *ProcessInfo) WaitContext(ctx context.Context) bool {
	if ctx.Err() != nil {
		_ = p.Close()
		return false
	}
	select {
	case <-p.done:
		return ctx.Err() == nil
	case <-ctx.Done():
		_ = p.Close()
		return false
	}
}

func (p *ProcessInfo) Close() error {
	p.mutex.Lock()
	if p.closed {
		p.mutex.Unlock()
		return nil
	}
	p.closed = true
	stop, cmd, running := p.stopCancel, p.Cmd, p.Running
	pipes := p.pipes
	p.pipes = nil
	p.stopCancel = nil
	p.mutex.Unlock()
	if stop != nil {
		stop()
	}
	for _, pipe := range pipes {
		_ = pipe.Close()
	}
	if running && cmd != nil && cmd.Process != nil {
		if cmd.Cancel != nil {
			return cmd.Cancel()
		}
		return cmd.Process.Kill()
	}
	return nil
}

func (p *ProcessInfo) IsClosed() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.closed
}

func (p *ProcessInfo) AddPipe(pipe io.Closer) {
	p.mutex.Lock()
	if p.closed {
		p.mutex.Unlock()
		_ = pipe.Close()
		return
	}
	p.pipes = append(p.pipes, pipe)
	p.mutex.Unlock()
}

// ClosePipes precedes proc_close's wait so a child waiting for EOF can exit.
func (p *ProcessInfo) ClosePipes() {
	p.mutex.Lock()
	pipes := p.pipes
	p.pipes = nil
	p.mutex.Unlock()
	for _, pipe := range pipes {
		_ = pipe.Close()
	}
}

func (p *ProcessInfo) BindRequestCancel(stop func() bool) {
	p.mutex.Lock()
	closed := p.closed || !p.Running
	if !closed {
		p.stopCancel = stop
	}
	p.mutex.Unlock()
	if closed {
		stop()
	}
}

// markDone 标记进程已结束，并关闭 done channel
func (p *ProcessInfo) markDone() {
	p.doneOnce.Do(func() {
		p.mutex.Lock()
		stop := p.stopCancel
		p.stopCancel = nil
		p.mutex.Unlock()
		if stop != nil {
			stop()
		}
		close(p.done)
	})
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
