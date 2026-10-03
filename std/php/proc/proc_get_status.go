package proc

import (
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

	// proc_open 的 Wait 协程负责更新生命周期；查询不能发送平台相关信号或改写状态。
	running := procInfo.GetRunning()

	// 创建状态数组
	status := data.NewArrayValue(nil).(*data.ArrayValue)
	status.SetStringKey("command", data.NewStringValue(procInfo.Command))
	status.SetStringKey("pid", data.NewIntValue(procInfo.Pid))
	status.SetStringKey("running", data.NewBoolValue(running))
	status.SetStringKey("signaled", data.NewBoolValue(false))
	status.SetStringKey("stopped", data.NewBoolValue(false))
	status.SetStringKey("exitcode", data.NewIntValue(procInfo.GetExitCode()))
	status.SetStringKey("termsig", data.NewIntValue(0))
	status.SetStringKey("stopsig", data.NewIntValue(0))

	return status, nil
}

func (f *ProcGetStatusFunction) GetName() string {
	return "proc_get_status"
}

var procGetStatusFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "process", 0, nil, nil),
}

func (f *ProcGetStatusFunction) GetParams() []data.GetValue {
	return procGetStatusFunctionGetParams
}

var procGetStatusFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "process", 0, data.NewBaseType("resource")),
}

func (f *ProcGetStatusFunction) GetVariables() []data.Variable {
	return procGetStatusFunctionGetVariables
}
