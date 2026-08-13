package proc

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// ProcTerminateFunction 实现 proc_terminate 函数
// 杀死由 proc_open 打开的进程
type ProcTerminateFunction struct{}

func NewProcTerminateFunction() data.FuncStmt {
	return &ProcTerminateFunction{}
}

func (f *ProcTerminateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 终止进程（跨平台兼容，Windows 不支持 syscall.Signal）
	if procInfo.Cmd != nil && procInfo.Cmd.Process != nil {
		err := procInfo.Cmd.Process.Kill()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		procInfo.SetRunning(false)
		return data.NewBoolValue(true), nil
	}

	return data.NewBoolValue(false), nil
}

func (f *ProcTerminateFunction) GetName() string {
	return "proc_terminate"
}

func (f *ProcTerminateFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "process", 0, nil, nil),
		node.NewParameter(nil, "signal", 1, node.NewIntLiteral(nil, "15"), nil), // 默认 SIGTERM
	}
}

func (f *ProcTerminateFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "process", 0, data.NewBaseType("resource")),
		node.NewVariable(nil, "signal", 1, data.NewBaseType("int")),
	}
}
