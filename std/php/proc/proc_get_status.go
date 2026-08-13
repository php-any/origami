package proc

import (
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// ProcGetStatusFunction 实现 proc_get_status 函数
// 获取由 proc_open 打开的进程的信息
type ProcGetStatusFunction struct{}

func NewProcGetStatusFunction() data.FuncStmt {
	return &ProcGetStatusFunction{}
}

func (f *ProcGetStatusFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取进程资源
	processValue, _ := ctx.GetIndexValue(0)
	if processValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 从资源对象中获取 ProcessInfo
	var procInfo *ProcessInfo
	if res, ok := processValue.(*core.ResourceValue); ok {
		resource := res.GetResource()
		if info, ok := resource.(*ProcessInfo); ok {
			procInfo = info
		} else {
			return data.NewBoolValue(false), nil
		}
	} else {
		return data.NewBoolValue(false), nil
	}

	// 检查进程是否还在运行
	running := procInfo.GetRunning()
	if procInfo.Cmd != nil && procInfo.Cmd.Process != nil {
		// 尝试发送信号 0 来检查进程是否还在运行
		err := procInfo.Cmd.Process.Signal(syscall.Signal(0))
		if err != nil {
			running = false
			procInfo.SetRunning(false)
		}
	}

	// 创建状态数组
	status := data.NewObjectValue()
	status.SetProperty("command", data.NewStringValue(procInfo.Command))
	status.SetProperty("pid", data.NewIntValue(procInfo.Pid))
	status.SetProperty("running", data.NewBoolValue(running))
	status.SetProperty("signaled", data.NewBoolValue(false))
	status.SetProperty("stopped", data.NewBoolValue(false))
	status.SetProperty("exitcode", data.NewIntValue(procInfo.GetExitCode()))
	status.SetProperty("termsig", data.NewIntValue(0))
	status.SetProperty("stopsig", data.NewIntValue(0))

	return status, nil
}

func (f *ProcGetStatusFunction) GetName() string {
	return "proc_get_status"
}

func (f *ProcGetStatusFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "process", 0, nil, nil),
	}
}

func (f *ProcGetStatusFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "process", 0, data.NewBaseType("resource")),
	}
}
