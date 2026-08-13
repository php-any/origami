package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IniGetFunction 实现 PHP 内置函数 ini_get
//
//	ini_get(string $option): string|false
//
// 获取指定配置选项的当前值；若配置项不存在则返回 false。
type IniGetFunction struct{}

func NewIniGetFunction() data.FuncStmt {
	return &IniGetFunction{}
}

func (f *IniGetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	optionVal, _ := ctx.GetIndexValue(0)
	if optionVal == nil {
		return data.NewBoolValue(false), nil
	}

	option := optionVal.AsString()

	value, ok := IniGet(option)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(value), nil
}

func (f *IniGetFunction) GetName() string {
	return "ini_get"
}

func (f *IniGetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "option", 0, nil, data.String{}),
	}
}

func (f *IniGetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "option", 0, data.String{}),
	}
}
