package proc

import (
	"bytes"
	"runtime"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ShellExecFunction 实现 shell_exec 函数
// 通过 shell 执行命令并返回完整输出字符串
type ShellExecFunction struct{}

func NewShellExecFunction() data.FuncStmt {
	return &ShellExecFunction{}
}

func (f *ShellExecFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cmdValue, _ := ctx.GetIndexValue(0)
	if cmdValue == nil {
		return data.NewNullValue(), nil
	}

	var cmd string
	if s, ok := cmdValue.(data.AsString); ok {
		cmd = s.AsString()
	} else {
		cmd = cmdValue.AsString()
	}

	if cmd == "" {
		return data.NewNullValue(), nil
	}

	cmdObj := shellCommand(ctx.GoContext(), cmd)
	cmdObj.WaitDelay = 2 * time.Second
	var stdout bytes.Buffer
	cmdObj.Stdout = &stdout

	_ = cmdObj.Run()
	if ctx.GoContext().Err() != nil {
		panic(data.ErrRequestCanceled)
	}
	output := stdout.String()
	if runtime.GOOS == "windows" {
		// PHP opens shell_exec's Windows pipe in CRT text mode: CRLF is
		// translated to LF and Ctrl-Z marks EOF. proc_open remains binary.
		if end := strings.IndexByte(output, 0x1a); end >= 0 {
			output = output[:end]
		}
		output = strings.ReplaceAll(output, "\r\n", "\n")
	}
	if len(output) == 0 {
		// PHP returns null when the command produces no output.
		return data.NewNullValue(), nil
	}

	return data.NewStringValue(output), nil
}

func (f *ShellExecFunction) GetName() string {
	return "shell_exec"
}

var shellExecFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "command", 0, nil, nil),
}

func (f *ShellExecFunction) GetParams() []data.GetValue {
	return shellExecFunctionGetParams
}

var shellExecFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "command", 0, data.NewBaseType("string")),
}

func (f *ShellExecFunction) GetVariables() []data.Variable {
	return shellExecFunctionGetVariables
}
