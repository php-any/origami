package proc

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// ProcCloseFunction 实现 proc_close 函数
// 关闭由 proc_open 打开的进程
type ProcCloseFunction struct{}

func NewProcCloseFunction() data.FuncStmt {
	return &ProcCloseFunction{}
}

func (f *ProcCloseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取进程资源
	processValue, _ := ctx.GetIndexValue(0)
	if processValue == nil {
		return data.NewIntValue(-1), nil
	}

	// 从资源对象中获取 ProcessInfo
	var procInfo *ProcessInfo
	if res, ok := processValue.(*core.ResourceValue); ok {
		resource := res.GetResource()
		if info, ok := resource.(*ProcessInfo); ok {
			procInfo = info
		} else {
			return data.NewIntValue(-1), nil
		}
	} else {
		return data.NewIntValue(-1), nil
	}

	// 获取退出码
	exitCode := procInfo.GetExitCode()

	// 如果进程还在运行，终止它
	if procInfo.GetRunning() && procInfo.Cmd != nil && procInfo.Cmd.Process != nil {
		procInfo.Cmd.Process.Kill()
		procInfo.SetRunning(false)
	}

	return data.NewIntValue(exitCode), nil
}

func (f *ProcCloseFunction) GetName() string {
	return "proc_close"
}

func (f *ProcCloseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "process", 0, nil, nil),
	}
}

func (f *ProcCloseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "process", 0, data.NewBaseType("resource")),
	}
}
