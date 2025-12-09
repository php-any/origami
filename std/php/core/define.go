package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Load 注册 PHP 核心相关函数
func Load(vm data.VM) {
	vm.AddFunc(NewDefineFunction())
}

func NewDefineFunction() data.FuncStmt {
	return &DefineFunction{}
}

// DefineFunction 实现 define(name, value [, case_insensitive]) 的基础支持
// 当前实现只保证调用成功，返回 true，用于兼容常见启动常量定义。
type DefineFunction struct {
	data.Function
}

func (f *DefineFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 尝试读取参数，保持与 PHP 一致的签名
	_, acl := ctx.GetVariableValue(node.NewVariable(nil, "name", 0, data.String{}))
	if acl != nil {
		return data.NewBoolValue(false), acl
	}
	_, acl = ctx.GetVariableValue(node.NewVariable(nil, "value", 1, nil))
	if acl != nil {
		return data.NewBoolValue(false), acl
	}
	// 第三个参数 case_insensitive 忽略
	return data.NewBoolValue(true), nil
}

func (f *DefineFunction) GetName() string {
	return "define"
}

func (f *DefineFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}

func (f *DefineFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
		node.NewVariable(nil, "value", 1, nil),
	}
}
