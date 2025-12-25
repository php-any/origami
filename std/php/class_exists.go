package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewClassExistsFunction 创建 class_exists 函数
// PHP 语义：
// class_exists(string $class, bool $autoload = true): bool
// 这里先忽略 $autoload 参数，直接检查 VM 中是否已注册该类
func NewClassExistsFunction() data.FuncStmt {
	return &ClassExistsFunction{}
}

type ClassExistsFunction struct{}

func (f *ClassExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 读取第一个参数 $class
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "class", 0, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}
	className := v.(data.AsString).AsString()

	vm := ctx.GetVM()
	// 在 VM 中检查类是否存在
	_, exist := vm.GetClass(className)
	if exist {
		return data.NewBoolValue(true), nil
	}
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *ClassExistsFunction) GetName() string {
	return "class_exists"
}

func (f *ClassExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, data.String{}),
	}
}

func (f *ClassExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.String{}),
	}
}
