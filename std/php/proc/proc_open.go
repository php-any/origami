package proc

import (
	"io"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/stream"
)

// ProcOpenFunction 实现 proc_open 函数
// 执行命令并打开文件指针
type ProcOpenFunction struct{}

func NewProcOpenFunction() data.FuncStmt {
	return &ProcOpenFunction{}
}

func (f *ProcOpenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取命令参数：string 或 list（PHP proc_open(['bin', 'arg'], ...)）
	cmdValue, _ := ctx.GetIndexValue(0)
	if cmdValue == nil {
		return data.NewBoolValue(false), nil
	}

	var cmdObj *exec.Cmd
	switch c := cmdValue.(type) {
	case *data.ArrayValue:
		parts := c.ToValueList()
		if len(parts) == 0 {
			return data.NewBoolValue(false), nil
		}
		name := parts[0].AsString()
		args := make([]string, 0, len(parts)-1)
		for _, p := range parts[1:] {
			args = append(args, p.AsString())
		}
		if name == "" {
			return data.NewBoolValue(false), nil
		}
		cmdObj = exec.Command(name, args...)
	default:
		var cmd string
		if s, ok := cmdValue.(data.AsString); ok {
			cmd = s.AsString()
		} else {
			cmd = cmdValue.AsString()
		}
		if cmd == "" {
			return data.NewBoolValue(false), nil
		}
		cmdObj = shellCommand(cmd)
	}

	// 获取描述符数组（可选）
	// PHP 格式: [0 => ['pipe', 'r'], 1 => ['pipe', 'w'], 2 => ['pipe', 'w']]
	descriptorspecValue, _ := ctx.GetIndexValue(1)
	descriptorspec := make(map[int][]interface{})
	if descriptorspecValue != nil {
		if obj, ok := descriptorspecValue.(*data.ObjectValue); ok {
			obj.RangeProperties(func(key string, value data.Value) bool {
				fd, err := strconv.Atoi(key)
				if err != nil {
					return true
				}
				if arr, ok := value.(*data.ArrayValue); ok && len(arr.List) >= 2 {
					valueList := arr.ToValueList()
					descType := valueList[0].AsString()
					descMode := valueList[1].AsString()
					if descType == "pipe" {
						descriptorspec[fd] = []interface{}{descType, descMode}
					}
				}
				return true
			})
		} else if arr, ok := descriptorspecValue.(*data.ArrayValue); ok {
			// 保留 PHP 数组键（1/2），不能用 ToValueList 的 0..n-1 下标
			for i, zval := range arr.List {
				if zval == nil || zval.Value == nil {
					continue
				}
				fd := i
				if zval.Name != "" {
					if parsed, err := strconv.Atoi(zval.Name); err == nil {
						fd = parsed
					}
				}
				arrVal, ok := zval.Value.(*data.ArrayValue)
				if !ok || len(arrVal.List) < 2 {
					continue
				}
				arrValList := arrVal.ToValueList()
				descType := arrValList[0].AsString()
				descMode := arrValList[1].AsString()
				if descType == "pipe" {
					descriptorspec[fd] = []interface{}{descType, descMode}
				}
			}
		}
	}

	// 获取管道数组（可选，用于返回文件指针）
	// 对于引用参数，需要获取 ZVal 引用以便直接更新
	// 参数索引 2 对应 pipes 参数（引用参数）
	pipesZVal := ctx.GetIndexZVal(2)
	if pipesZVal == nil {
		return data.NewBoolValue(false), nil
	}
	// PHP：proc_open 会重建 $pipes；勿复用上次已 fclose 的流资源对象
	pipes := data.NewObjectValue()
	pipesZVal.Value = pipes

	// 处理描述符
	var stdoutPipe, stderrPipe io.ReadCloser
	var err error

	// 根据描述符配置创建管道
	if len(descriptorspec) > 0 {
		if desc, ok := descriptorspec[1]; ok && len(desc) >= 2 && desc[1] == "w" {
			// stdout (1) - 读取管道
			stdoutPipe, err = cmdObj.StdoutPipe()
			if err != nil {
				return data.NewBoolValue(false), nil
			}
		}
		if desc, ok := descriptorspec[2]; ok && len(desc) >= 2 && desc[1] == "w" {
			// stderr (2) - 读取管道
			stderrPipe, err = cmdObj.StderrPipe()
			if err != nil {
				return data.NewBoolValue(false), nil
			}
		}
	} else {
		// 如果没有指定描述符，默认创建所有管道
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
	err = cmdObj.Start()
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 创建进程信息对象
	procInfo := NewProcessInfo(cmdObj, cmdObj.Path)
	realPID := cmdObj.Process.Pid

	// 创建进程资源类，使用真实的系统进程ID作为资源ID
	// Resource 字段存储 ProcessInfo，包含 cmdObj 和命令信息
	resourceClass := core.NewResourceClass("process", procInfo, realPID)

	// 创建进程资源对象（使用 ResourceValue，嵌入 ClassValue）
	procResource := core.NewResourceValue(resourceClass, ctx)

	// 设置 pipes（存储管道信息）
	// 创建流资源对象并放入 pipes 对象
	// stdout (1) - 读取管道
	if stdoutPipe != nil {
		stdoutStreamInfo := stream.NewStreamInfoFromReader(stdoutPipe, "r")
		stdoutFd := realPID*10 + 1 // 生成一个唯一的文件描述符
		stdoutResourceClass := core.NewResourceClass("stream", stdoutStreamInfo, stdoutFd)
		stdoutResource := core.NewResourceValue(stdoutResourceClass, ctx)
		pipes.SetProperty("1", stdoutResource)
	}

	// stderr (2) - 读取管道
	if stderrPipe != nil {
		stderrStreamInfo := stream.NewStreamInfoFromReader(stderrPipe, "r")
		stderrFd := realPID*10 + 2 // 生成一个唯一的文件描述符
		stderrResourceClass := core.NewResourceClass("stream", stderrStreamInfo, stderrFd)
		stderrResource := core.NewResourceValue(stderrResourceClass, ctx)
		pipes.SetProperty("2", stderrResource)
	}

	// 更新引用参数的 ZVal.Value（显式重新赋值，确保引用参数被正确更新）
	pipesZVal.Value = pipes

	// 在后台等待进程结束并更新状态
	// 注意：不在这里关闭管道，因为 stream_get_contents 需要读取管道数据
	// 管道应该在 proc_close 或流关闭时关闭
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
		// 不在这里关闭管道，让 stream_get_contents 可以读取数据
		// 管道会在 proc_close 或流关闭时关闭
	}()

	return procResource, nil
}

func shellCommand(cmd string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", cmd)
	}
	return exec.Command("sh", "-c", cmd)
}

func (f *ProcOpenFunction) GetName() string {
	return "proc_open"
}

func (f *ProcOpenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "command", 0, nil, nil),
		node.NewParameter(nil, "descriptorspec", 1, node.NewNullLiteral(nil), nil),
		node.NewParameterReference(nil, "pipes", 2, nil, data.Mixed{}),
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
