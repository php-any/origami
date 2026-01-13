package core

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RealpathFunction 实现 realpath 函数
// realpath(string $path): string|false
// 返回规范化的绝对路径名，如果路径不存在则返回 false
type RealpathFunction struct{}

func NewRealpathFunction() data.FuncStmt {
	return &RealpathFunction{}
}

func (f *RealpathFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathValue, _ := ctx.GetIndexValue(0)

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		if pathValue == nil {
			return data.NewBoolValue(false), nil
		}
		path = pathValue.AsString()
	}

	// 如果路径为空，返回 false
	if path == "" {
		return data.NewBoolValue(false), nil
	}

	// 检查路径是否存在
	if _, err := os.Stat(path); err != nil {
		// 路径不存在，返回 false
		return data.NewBoolValue(false), nil
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 解析符号链接（如果存在）
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// 如果解析符号链接失败，使用绝对路径
		realPath = absPath
	}

	// 规范化路径（移除多余的斜杠和 `.`、`..` 组件）
	realPath = filepath.Clean(realPath)

	return data.NewStringValue(realPath), nil
}

func (f *RealpathFunction) GetName() string {
	return "realpath"
}

func (f *RealpathFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, nil),
	}
}

func (f *RealpathFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
	}
}
