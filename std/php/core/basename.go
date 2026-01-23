package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// BasenameFunction 实现 basename 函数
// basename(string $path, string $suffix = ""): string
// 返回路径中的文件名部分，如果提供了 suffix 参数，会移除该后缀
type BasenameFunction struct{}

func NewBasenameFunction() data.FuncStmt {
	return &BasenameFunction{}
}

func (f *BasenameFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathValue, _ := ctx.GetIndexValue(0)
	suffixValue, _ := ctx.GetIndexValue(1) // 可选的 suffix 参数

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		if pathValue == nil {
			return data.NewStringValue(""), nil
		}
		path = pathValue.AsString()
	}

	// 如果路径为空，返回空字符串
	if path == "" {
		return data.NewStringValue(""), nil
	}

	// PHP 的 basename 在 Windows 路径下会同时识别 "/" 和 "\" 作为分隔符
	// 因此这里同时处理这两种分隔符
	// 1. 先移除末尾的路径分隔符（无论是 "/" 还是 "\"）
	normalizedPath := strings.TrimRight(path, "/\\")

	// 如果路径为空（只有斜杠），返回空字符串
	if normalizedPath == "" {
		return data.NewStringValue(""), nil
	}

	// 获取最后一个路径分隔符后的部分（支持 "/" 和 "\"）
	lastSlash := strings.LastIndexAny(normalizedPath, "/\\")
	var basename string
	if lastSlash == -1 {
		// 如果没有找到正斜杠，返回整个字符串
		basename = normalizedPath
	} else {
		basename = normalizedPath[lastSlash+1:]
	}

	// 如果提供了 suffix 参数，移除该后缀
	if suffixValue != nil {
		suffix := suffixValue.AsString()
		if suffix != "" && strings.HasSuffix(basename, suffix) {
			basename = strings.TrimSuffix(basename, suffix)
		}
	}

	return data.NewStringValue(basename), nil
}

func (f *BasenameFunction) GetName() string {
	return "basename"
}

func (f *BasenameFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, nil),
		node.NewParameter(nil, "suffix", 1, data.NewStringValue(""), nil),
	}
}

func (f *BasenameFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "suffix", 1, data.NewBaseType("string")),
	}
}
