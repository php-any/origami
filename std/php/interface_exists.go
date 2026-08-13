package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewInterfaceExistsFunction 创建 interface_exists 函数
// PHP 语义：
// interface_exists(string $interface, bool $autoload = true): bool
// 这里与 class_exists 的实现方式保持一致：
// - 先检查 VM 中是否已有该接口
// - 若不存在，则通过 GetOrLoadInterface 触发自动加载
// - 忽略 $autoload 参数（总是按 autoload=true 处理）
func NewInterfaceExistsFunction() data.FuncStmt {
	return &InterfaceExistsFunction{}
}

type InterfaceExistsFunction struct{}

func (f *InterfaceExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 读取第一个参数 $interface
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "interface", 0, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}
	interfaceName := v.(data.AsString).AsString()

	vm := ctx.GetVM()

	// 1. 先直接在 VM 中检查接口是否已注册
	if _, exist := vm.GetInterface(interfaceName); exist {
		return data.NewBoolValue(true), nil
	}

	// 2. 若未找到，则尝试通过自动加载机制加载接口定义
	iface, acl := vm.GetOrLoadInterface(interfaceName)
	if acl != nil {
		// 若自动加载阶段抛出“找不到接口”之类的错误，按 interface_exists 语义视为不存在
		// 避免将内部加载错误暴露给用户代码
		return data.NewBoolValue(false), nil
	}
	if iface == nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *InterfaceExistsFunction) GetName() string {
	return "interface_exists"
}

func (f *InterfaceExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "interface", 0, nil, data.String{}),
	}
}

func (f *InterfaceExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "interface", 0, data.String{}),
	}
}
