package php

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewVarDumpFunction 创建 var_dump 函数
func NewVarDumpFunction() data.FuncStmt {
	return &VarDumpFunction{}
}

// VarDumpFunction 实现 PHP 风格的 var_dump
// 为了简单，当前实现主要是把参数按顺序打印出来，
// 字符串按内容打印，其它类型用 fmt.Println 输出其值。
type VarDumpFunction struct{}

func (f *VarDumpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 参数定义里只有一个 node.Parameters，取到实际传入的所有参数
	for _, argument := range f.GetParams() {
		argv, _ := argument.GetValue(ctx)

		switch temp := argv.(type) {
		case data.Variable:
			// 如果是变量，取出真实值再打印
			v, acl := ctx.GetVariableValue(temp)
			if acl != nil {
				return nil, acl
			}
			switch arg := v.(type) {
			case data.AsString:
				fmt.Println(arg.AsString())
			default:
				fmt.Println(arg)
			}
		default:
			// 直接是值，按类型打印
			switch arg := temp.(type) {
			case data.AsString:
				fmt.Println(arg.AsString())
			default:
				fmt.Println(arg)
			}
		}
	}

	return nil, nil
}

func (f *VarDumpFunction) GetName() string {
	return "var_dump"
}

func (f *VarDumpFunction) GetParams() []data.GetValue {
	// 使用 node.Parameters 支持任意数量的参数
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

func (f *VarDumpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}
