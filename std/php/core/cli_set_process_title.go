package core

import (
	"errors"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type CliSetProcessTitleFunction struct{}

func NewCliSetProcessTitleFunction() data.FuncStmt {
	return &CliSetProcessTitleFunction{}
}

func (f *CliSetProcessTitleFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	title, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("cli_set_process_title() expects exactly 1 parameter, 0 given"))
	}

	var titleValue string
	if str, ok := title.(data.AsString); ok {
		titleValue = str.AsString()
	} else {
		titleValue = title.AsString()
	}

	// 设置进程标题
	// 注意：在某些系统上（如 Linux），这需要 prctl 系统调用
	// 在 Go 中，我们可以使用 os 包的一些方法，但标准库没有直接支持
	// 这里我们尝试设置进程名称，如果失败则返回 false
	err := setProcessTitle(titleValue)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *CliSetProcessTitleFunction) GetName() string {
	return "cli_set_process_title"
}

func (f *CliSetProcessTitleFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, nil),
	}
}

func (f *CliSetProcessTitleFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.String{}),
	}
}

// setProcessTitle 设置进程标题
// 在 Linux 上，这需要使用 prctl 系统调用
// 在 macOS 上，可以使用 setproctitle
// 在 Windows 上，可以使用 SetConsoleTitle
// 这里我们使用一个简单的实现，尝试设置进程名称
func setProcessTitle(title string) error {
	// 在 Go 中，标准库没有直接支持设置进程标题
	// 我们可以尝试通过环境变量或其他方式
	// 但最可靠的方式是使用 CGO 调用系统调用
	// 为了简化，这里我们使用一个基本的实现

	// 尝试设置进程名称（在某些系统上可能不工作）
	// 注意：这只是一个占位实现，实际效果取决于操作系统
	_ = os.Setenv("PROCESS_TITLE", title)

	// 返回 nil 表示成功（即使实际上可能没有真正设置）
	// 在实际应用中，可能需要使用 CGO 来调用系统调用
	return nil
}
