package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UnlinkFunction 实现 unlink 函数
// 对标 PHP unlink: 删除文件，成功返回 true，失败返回 false（不抛异常）
type UnlinkFunction struct{}

func NewUnlinkFunction() data.FuncStmt {
	return &UnlinkFunction{}
}

func (f *UnlinkFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 尝试删除文件
	if err := os.Remove(path); err != nil {
		// 与 is_dir/is_file 风格保持一致：失败返回 false，不抛异常
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *UnlinkFunction) GetName() string {
	return "unlink"
}

func (f *UnlinkFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *UnlinkFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
