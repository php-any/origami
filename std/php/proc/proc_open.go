package proc

import (
	"io"
	"os/exec"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// ProcOpenFunction 实现 proc_open 函数
// 执行命令并打开文件指针
type ProcOpenFunction struct{}

func NewProcOpenFunction() data.FuncStmt {
	return &ProcOpenFunction{}
}

func (f *ProcOpenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取命令参数
	cmdValue, _ := ctx.GetIndexValue(0)
	if cmdValue == nil {
		return data.NewBoolValue(false), nil
	}

	var cmd string
	if s, ok := cmdValue.(data.AsString); ok {
		cmd = s.AsString()
	} else {
		cmd = cmdValue.AsString()
	}

	if cmd == "" {
		return data.NewBoolValue(false), nil
	}

	// 获取描述符数组（可选）
	descriptorspecValue, _ := ctx.GetIndexValue(1)
	var descriptorspec map[string]interface{}
	if descriptorspecValue != nil {
		if obj, ok := descriptorspecValue.(*data.ObjectValue); ok {
			descriptorspec = make(map[string]interface{})
			obj.RangeProperties(func(key string, value data.Value) bool {
				descriptorspec[key] = value
				return true
			})
		}
	}

	// 获取管道数组（可选，用于返回文件指针）
	pipesValue, _ := ctx.GetIndexValue(2)
	var pipes *data.ObjectValue
	if pipesValue != nil {
		if obj, ok := pipesValue.(*data.ObjectValue); ok {
			pipes = obj
		} else {
			pipes = data.NewObjectValue()
		}
	} else {
		pipes = data.NewObjectValue()
	}

	// 创建命令
	cmdObj := exec.Command("sh", "-c", cmd)

	// 处理描述符
	var stdinPipe io.WriteCloser
	var stdoutPipe, stderrPipe io.ReadCloser
	if len(descriptorspec) == 0 {
		// 默认描述符：创建管道
		var err error
		stdinPipe, err = cmdObj.StdinPipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stdoutPipe, err = cmdObj.StdoutPipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stderrPipe, err = cmdObj.StderrPipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
	} else {
		// 处理自定义描述符（简化实现）
		var err error
		stdinPipe, err = cmdObj.StdinPipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stdoutPipe, err = cmdObj.StdoutPipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stderrPipe, err = cmdObj.StderrPipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
	}

	// 启动进程
	err := cmdObj.Start()
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 创建进程信息对象
	procInfo := NewProcessInfo(cmdObj, cmd)
	realPID := cmdObj.Process.Pid

	// 创建进程资源类，使用真实的系统进程ID作为资源ID
	// Resource 字段存储 ProcessInfo，包含 cmdObj 和命令信息
	resourceClass := core.NewResourceClass("process", procInfo, realPID)

	// 创建进程资源对象（使用 ResourceValue，嵌入 ClassValue）
	procResource := core.NewResourceValue(resourceClass, ctx)

	// 设置 pipes（存储管道信息）
	// 注意：这里简化实现，实际应该创建可读写的管道对象
	pipes.SetProperty("0", data.NewIntValue(realPID)) // stdin (使用 PID 作为标识)
	pipes.SetProperty("1", data.NewIntValue(realPID)) // stdout
	pipes.SetProperty("2", data.NewIntValue(realPID)) // stderr

	// 如果 pipes 参数是通过引用传递的，需要更新原对象
	if pipesValue != nil {
		if obj, ok := pipesValue.(*data.ObjectValue); ok {
			obj.SetProperty("0", data.NewIntValue(realPID))
			obj.SetProperty("1", data.NewIntValue(realPID))
			obj.SetProperty("2", data.NewIntValue(realPID))
		}
	}

	// 在后台等待进程结束并更新状态
	go func() {
		err := cmdObj.Wait()
		exitCode := -1
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}
		} else {
			exitCode = 0
		}
		procInfo.SetRunning(false)
		procInfo.SetExitCode(exitCode)
		// 关闭管道
		if stdinPipe != nil {
			stdinPipe.Close()
		}
		if stdoutPipe != nil {
			stdoutPipe.Close()
		}
		if stderrPipe != nil {
			stderrPipe.Close()
		}
	}()

	return procResource, nil
}

func (f *ProcOpenFunction) GetName() string {
	return "proc_open"
}

func (f *ProcOpenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "command", 0, nil, nil),
		node.NewParameter(nil, "descriptorspec", 1, node.NewNullLiteral(nil), nil),
		node.NewParameterReference(nil, "pipes", 2, data.Mixed{}),
		node.NewParameter(nil, "cwd", 3, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "env_vars", 4, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "options", 5, node.NewNullLiteral(nil), nil),
	}
}

func (f *ProcOpenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "command", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "descriptorspec", 1, data.NewBaseType("array")),
		node.NewVariable(nil, "pipes", 2, data.NewBaseType("array")),
		node.NewVariable(nil, "cwd", 3, data.NewBaseType("string")),
		node.NewVariable(nil, "env_vars", 4, data.NewBaseType("array")),
		node.NewVariable(nil, "options", 5, data.NewBaseType("array")),
	}
}
