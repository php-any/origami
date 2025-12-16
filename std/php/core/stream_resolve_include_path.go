package core

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StreamResolveIncludePathFunction 实现 stream_resolve_include_path 函数
// 该函数在 include_path 中查找文件，返回文件的绝对路径，如果没找到则返回 false
type StreamResolveIncludePathFunction struct{}

func NewStreamResolveIncludePathFunction() data.FuncStmt {
	return &StreamResolveIncludePathFunction{}
}

func (f *StreamResolveIncludePathFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取文件名参数
	filenameValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	var filename string
	switch p := filenameValue.(type) {
	case data.AsString:
		filename = p.AsString()
	default:
		filename = filenameValue.AsString()
	}

	// 如果文件名为空或者是绝对路径，直接检查是否存在
	if filename == "" {
		return data.NewBoolValue(false), nil
	}

	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			// 文件存在，返回绝对路径
			absPath, _ := filepath.Abs(filename)
			return data.NewStringValue(absPath), nil
		}
		return data.NewBoolValue(false), nil
	}

	// 尝试相对当前工作目录查找
	if _, err := os.Stat(filename); err == nil {
		absPath, _ := filepath.Abs(filename)
		return data.NewStringValue(absPath), nil
	}

	// 获取 include_path（简化实现，可以从环境变量或配置获取）
	// 这里使用当前目录作为默认的 include_path
	includePaths := []string{
		".",
	}

	// 尝试在 include_path 中查找文件
	for _, includePath := range includePaths {
		testPath := filepath.Join(includePath, filename)
		if _, err := os.Stat(testPath); err == nil {
			absPath, _ := filepath.Abs(testPath)
			return data.NewStringValue(absPath), nil
		}
	}

	// 文件未找到
	return data.NewBoolValue(false), nil
}

func (f *StreamResolveIncludePathFunction) GetName() string {
	return "stream_resolve_include_path"
}

func (f *StreamResolveIncludePathFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *StreamResolveIncludePathFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
