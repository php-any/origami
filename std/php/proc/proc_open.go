package proc

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/stream"
	"github.com/php-any/origami/utils"
)

// ProcOpenFunction 实现 proc_open 函数
// 执行命令并打开文件指针
type ProcOpenFunction struct{}

func NewProcOpenFunction() data.FuncStmt {
	return &ProcOpenFunction{}
}

func (f *ProcOpenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	request := ctx.GoContext()
	if request.Err() != nil {
		panic(data.ErrRequestCanceled)
	}
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
		cmdObj = exec.CommandContext(request, name, args...)
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
		cmdObj = shellCommand(request, cmd)
	}
	cmdObj.WaitDelay = 2 * time.Second

	// 获取描述符数组（可选）
	// PHP 格式: [0 => ['pipe', 'r'], 1 => ['pipe', 'w'], 2 => ['pipe', 'w']]
	descriptorspecValue, _ := ctx.GetIndexValue(1)
	descriptorspec := make(map[int][]interface{})
	if descriptorspecValue != nil {
		if arr, ok := descriptorspecValue.(*data.ArrayValue); ok {
			for arraySlots129, // 保留 PHP 数组键（1/2），不能用 ToValueList 的 0..n-1 下标
				i := arr.View(), 0; i < arraySlots129.Len(); i++ {
				zval := arraySlots129.At(i)
				if zval == nil || zval.ReadValue() == nil {
					continue
				}
				fd := i
				if zval.Name != "" {
					if parsed, err := strconv.Atoi(zval.Name); err == nil {
						fd = parsed
					}
				}
				arrVal, ok := zval.ReadValue().(*data.ArrayValue)
				if !ok || arrVal.Len() < 2 {
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
	pipes := data.NewArrayValue(nil).(*data.ArrayValue)
	pipesZVal.StoreRaw(pipes)

	// 处理描述符
	var stdoutPipe, stderrPipe io.ReadCloser
	var stdoutWriter, stderrWriter *os.File
	var err error
	started := false
	defer func() {
		if stdoutWriter != nil {
			_ = stdoutWriter.Close()
		}
		if stderrWriter != nil {
			_ = stderrWriter.Close()
		}
		if !started {
			if stdoutPipe != nil {
				_ = stdoutPipe.Close()
			}
			if stderrPipe != nil {
				_ = stderrPipe.Close()
			}
		}
	}()

	// 根据描述符配置创建管道
	if len(descriptorspec) > 0 {
		if desc, ok := descriptorspec[1]; ok && len(desc) >= 2 && desc[1] == "w" {
			// stdout (1) - 读取管道
			var reader *os.File
			reader, stdoutWriter, err = os.Pipe()
			if err != nil {
				return data.NewBoolValue(false), nil
			}
			stdoutPipe = reader
			cmdObj.Stdout = stdoutWriter
		}
		if desc, ok := descriptorspec[2]; ok && len(desc) >= 2 && desc[1] == "w" {
			// stderr (2) - 读取管道
			var reader *os.File
			reader, stderrWriter, err = os.Pipe()
			if err != nil {
				return data.NewBoolValue(false), nil
			}
			stderrPipe = reader
			cmdObj.Stderr = stderrWriter
		}
	} else {
		// 如果没有指定描述符，默认创建所有管道
		var stdoutReader *os.File
		stdoutReader, stdoutWriter, err = os.Pipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stdoutPipe = stdoutReader
		cmdObj.Stdout = stdoutWriter
		var stderrReader *os.File
		stderrReader, stderrWriter, err = os.Pipe()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stderrPipe = stderrReader
		cmdObj.Stderr = stderrWriter
	}

	// 处理 cwd 参数（索引 3）：PHP 允许指定子进程工作目录
	if cwdValue, _ := ctx.GetIndexValue(3); cwdValue != nil {
		if _, isNull := cwdValue.(*data.NullValue); !isNull {
			cmdObj.Dir = cwdValue.AsString()
		}
	}

	cmdObj.Env = node.EnvironmentEntries(ctx)
	// An explicit env_vars array replaces the inherited environment, even empty.
	if envValue, _ := ctx.GetIndexValue(4); envValue != nil {
		if _, isNull := envValue.(*data.NullValue); !isNull {
			if envArr, ok := envValue.(*data.ArrayValue); ok {
				env := make([]string, 0, envArr.Len())
				for arraySlots130, arrayPosition130 := envArr.View(), 0; arrayPosition130 < arraySlots130.Len(); arrayPosition130++ {
					item := arraySlots130.At(arrayPosition130)
					if item.ReadValue() == nil {
						continue
					}
					env = append(env, item.Name+"="+item.ReadValue().AsString())
				}
				cmdObj.Env = env
			}
		}
	}

	// 启动进程
	finishTree, err := utils.StartCommandTree(cmdObj)
	if err != nil {
		if request.Err() != nil {
			panic(data.ErrRequestCanceled)
		}
		return data.NewBoolValue(false), nil
	}
	started = true
	if stdoutWriter != nil {
		_ = stdoutWriter.Close()
	}
	if stderrWriter != nil {
		_ = stderrWriter.Close()
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
		procInfo.AddPipe(stdoutStreamInfo)
		stdoutFd := realPID*10 + 1 // 生成一个唯一的文件描述符
		stdoutResourceClass := core.NewResourceClass("stream", stdoutStreamInfo, stdoutFd)
		stdoutResource := core.NewResourceValue(stdoutResourceClass, ctx)
		pipes.SetIntKey(1, stdoutResource)
	}

	// stderr (2) - 读取管道
	if stderrPipe != nil {
		stderrStreamInfo := stream.NewStreamInfoFromReader(stderrPipe, "r")
		procInfo.AddPipe(stderrStreamInfo)
		stderrFd := realPID*10 + 2 // 生成一个唯一的文件描述符
		stderrResourceClass := core.NewResourceClass("stream", stderrStreamInfo, stderrFd)
		stderrResource := core.NewResourceValue(stderrResourceClass, ctx)
		pipes.SetIntKey(2, stderrResource)
	}

	// 更新引用参数的 ZVal.Value（显式重新赋值，确保引用参数被正确更新）
	pipesZVal.StoreRaw(pipes)

	// 在后台等待进程结束并更新状态
	// 注意：不在这里关闭管道，因为 stream_get_contents 需要读取管道数据
	// 管道应该在 proc_close 或流关闭时关闭
	go func() {
		defer finishTree()
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
		procInfo.markDone()
	}()

	return procResource, nil
}

func shellCommand(ctx context.Context, cmd string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", cmd)
	}
	return exec.CommandContext(ctx, "sh", "-c", cmd)
}

func (f *ProcOpenFunction) GetName() string {
	return "proc_open"
}

var procOpenFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "command", 0, nil, nil),
	node.NewParameter(nil, "descriptorspec", 1, node.NewNullLiteral(nil), nil),
	node.NewParameterReference(nil, "pipes", 2, nil, data.Mixed{}),
	node.NewParameter(nil, "cwd", 3, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "env_vars", 4, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "options", 5, node.NewNullLiteral(nil), nil),
}

func (f *ProcOpenFunction) GetParams() []data.GetValue {
	return procOpenFunctionGetParams
}

var procOpenFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "command", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "descriptorspec", 1, data.NewBaseType("array")),
	node.NewVariable(nil, "pipes", 2, data.NewBaseType("array")),
	node.NewVariable(nil, "cwd", 3, data.NewBaseType("string")),
	node.NewVariable(nil, "env_vars", 4, data.NewBaseType("array")),
	node.NewVariable(nil, "options", 5, data.NewBaseType("array")),
}

func (f *ProcOpenFunction) GetVariables() []data.Variable {
	return procOpenFunctionGetVariables
}
