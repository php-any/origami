package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// NewClassExistsFunction 创建 class_exists 函数
// PHP 语义：
// class_exists(string $class, bool $autoload = true): bool
// $autoload = false 时，仅检查内存中是否已定义，不尝试自动加载
func NewClassExistsFunction() data.FuncStmt {
	return &ClassExistsFunction{}
}

type ClassExistsFunction struct{}

func (f *ClassExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, _ := utils.ConvertFromIndex[string](ctx, 0)

	// 读取第二个参数 $autoload (默认 true)
	autoload, _ := utils.ConvertFromIndex[bool](ctx, 1)

	vm := ctx.GetVM()
	// 在 VM 中检查类是否存在
	_, exist := vm.GetClass(className)
	if exist {
		return data.NewBoolValue(true), nil
	}

	// 如果 $autoload 为 false，直接返回 false，不尝试加载
	if !autoload {
		return data.NewBoolValue(false), nil
	}

	// $autoload 为 true 时，尝试加载类
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return data.NewBoolValue(false), nil
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
		node.NewParameter(nil, "autoload", 1, data.NewBoolValue(true), data.Bool{}),
	}
}

func (f *ClassExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.String{}),
		node.NewVariable(nil, "autoload", 1, data.Bool{}),
	}
}
