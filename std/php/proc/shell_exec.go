package proc

import (
	"bytes"
	"os/exec"

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

	cmdObj := exec.Command("sh", "-c", cmd)
	var stdout bytes.Buffer
	cmdObj.Stdout = &stdout

	err := cmdObj.Run()
	if err != nil {
		// shell_exec 失败时返回 null
		return data.NewNullValue(), nil
	}

	return data.NewStringValue(stdout.String()), nil
}

func (f *ShellExecFunction) GetName() string {
	return "shell_exec"
}

func (f *ShellExecFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "command", 0, nil, nil),
	}
}

func (f *ShellExecFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "command", 0, data.NewBaseType("string")),
	}
}
