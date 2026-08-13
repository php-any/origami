package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RmdirFunction 实现 rmdir 函数
// 语义上与 unlink 类似：尝试删除给定路径（文件或空目录），成功返回 true，失败返回 false
// 注意：底层使用 os.Remove，非空目录会删除失败并返回 false
type RmdirFunction struct{}

func NewRmdirFunction() data.FuncStmt {
	return &RmdirFunction{}
}

func (f *RmdirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取路径参数
	pathValue, _ := ctx.GetIndexValue(0)
	if pathValue == nil {
		return data.NewBoolValue(false), nil
	}

	var path string
	if s, ok := pathValue.(data.AsString); ok {
		path = s.AsString()
	} else {
		path = pathValue.AsString()
	}

	// 空路径直接返回 false
	if path == "" {
		return data.NewBoolValue(false), nil
	}

	// 尝试删除目录（或文件）
	if err := os.Remove(path); err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *RmdirFunction) GetName() string {
	return "rmdir"
}

func (f *RmdirFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, nil),
	}
}

func (f *RmdirFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.NewBaseType("string")),
	}
}
