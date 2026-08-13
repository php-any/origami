package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IniSetFunction 实现 PHP 内置函数 ini_set
//
//	ini_set(string $option, string|int|float|bool|null $value): string|false
//
// 设置指定配置选项的值，返回该选项的旧值；若配置项不存在则返回 false。
type IniSetFunction struct{}

func NewIniSetFunction() data.FuncStmt {
	return &IniSetFunction{}
}

func (f *IniSetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	optionVal, _ := ctx.GetIndexValue(0)
	valueVal, _ := ctx.GetIndexValue(1)

	if optionVal == nil {
		return data.NewBoolValue(false), nil
	}

	option := optionVal.AsString()

	var newValue string
	if valueVal != nil {
		newValue = valueVal.AsString()
	}

	old, hadOld := IniSet(option, newValue)
	if !hadOld {
		// 配置项之前不存在，返回空字符串（PHP 行为：ini_set 首次设置也返回 ""）
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(old), nil
}

func (f *IniSetFunction) GetName() string {
	return "ini_set"
}

func (f *IniSetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "option", 0, nil, data.String{}),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}

func (f *IniSetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "option", 0, data.String{}),
		node.NewVariable(nil, "value", 1, data.NewBaseType("mixed")),
	}
}
