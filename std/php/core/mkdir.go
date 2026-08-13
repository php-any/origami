package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MkdirFunction 实现 PHP 内置函数 mkdir
//
//	mkdir(string $directory, int $permissions = 0777, bool $recursive = false): bool
//
// 创建目录。$recursive 为 true 时递归创建多级目录。成功返回 true，失败返回 false。
type MkdirFunction struct{}

func NewMkdirFunction() data.FuncStmt {
	return &MkdirFunction{}
}

func (f *MkdirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dirValue, _ := ctx.GetIndexValue(0)
	if dirValue == nil {
		return data.NewBoolValue(false), nil
	}

	directory := dirValue.AsString()
	if directory == "" {
		return data.NewBoolValue(false), nil
	}

	// 权限参数，默认 0777
	permissions := os.FileMode(0777)
	if permValue, _ := ctx.GetIndexValue(1); permValue != nil {
		if asInt, ok := permValue.(data.AsInt); ok {
			if perm, err := asInt.AsInt(); err == nil {
				permissions = os.FileMode(perm)
			}
		}
	}

	// 递归参数，默认 false
	recursive := false
	if recValue, _ := ctx.GetIndexValue(2); recValue != nil {
		if asBool, ok := recValue.(data.AsBool); ok {
			if r, err := asBool.AsBool(); err == nil {
				recursive = r
			}
		}
	}

	var err error
	if recursive {
		err = os.MkdirAll(directory, permissions)
	} else {
		err = os.Mkdir(directory, permissions)
	}

	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *MkdirFunction) GetName() string {
	return "mkdir"
}

func (f *MkdirFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, nil),
		node.NewParameter(nil, "permissions", 1, node.NewIntLiteral(nil, "0777"), nil),
		node.NewParameter(nil, "recursive", 2, nil, nil),
	}
}

func (f *MkdirFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "permissions", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "recursive", 2, data.NewBaseType("bool")),
	}
}
