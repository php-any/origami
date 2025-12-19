package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RenameFunction 实现 rename 函数
// 对标 PHP rename: 重命名/移动文件，成功返回 true，失败返回 false（不抛异常）
type RenameFunction struct{}

func NewRenameFunction() data.FuncStmt {
	return &RenameFunction{}
}

func (f *RenameFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取源文件和目标文件
	oldValue, _ := ctx.GetIndexValue(0)
	newValue, _ := ctx.GetIndexValue(1)
	if oldValue == nil || newValue == nil {
		return data.NewBoolValue(false), nil
	}

	var oldName, newName string
	if s, ok := oldValue.(data.AsString); ok {
		oldName = s.AsString()
	} else {
		oldName = oldValue.AsString()
	}
	if s, ok := newValue.(data.AsString); ok {
		newName = s.AsString()
	} else {
		newName = newValue.AsString()
	}

	if oldName == "" || newName == "" {
		return data.NewBoolValue(false), nil
	}

	if err := os.Rename(oldName, newName); err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *RenameFunction) GetName() string {
	return "rename"
}

func (f *RenameFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "oldname", 0, nil, nil),
		node.NewParameter(nil, "newname", 1, nil, nil),
	}
}

func (f *RenameFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "oldname", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "newname", 1, data.NewBaseType("string")),
	}
}
