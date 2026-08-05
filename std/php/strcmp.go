package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrcmpFunction 实现 strcmp 全局函数
type StrcmpFunction struct{}

// NewStrcmpFunction 创建一个新的 strcmp 函数实例
func NewStrcmpFunction() data.FuncStmt {
	return &StrcmpFunction{}
}

// GetName 返回函数名
func (f *StrcmpFunction) GetName() string {
	return "strcmp"
}

// GetParams 返回参数列表
func (f *StrcmpFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str1", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "str2", 1, nil, data.Mixed{}),
	}
}

// GetVariables 返回变量列表
func (f *StrcmpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str1", 0, data.Mixed{}),
		node.NewVariable(nil, "str2", 1, data.Mixed{}),
	}
}

// Call 执行 strcmp 函数
// strcmp 用于二进制安全字符串比较，区分大小写
// 返回值：<0 如果 str1 < str2，>0 如果 str1 > str2，0 如果 str1 == str2
func (f *StrcmpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 str1 参数
	str1Value, _ := ctx.GetIndexValue(0)
	if str1Value == nil {
		return data.NewIntValue(0), nil
	}

	str1 := ""
	if s, ok := str1Value.(data.AsString); ok {
		str1 = s.AsString()
	} else {
		str1 = str1Value.AsString()
	}

	// 获取 str2 参数
	str2Value, _ := ctx.GetIndexValue(1)
	if str2Value == nil {
		return data.NewIntValue(0), nil
	}

	str2 := ""
	if s, ok := str2Value.(data.AsString); ok {
		str2 = s.AsString()
	} else {
		str2 = str2Value.AsString()
	}

	// 使用 strings.Compare 进行比较（Go 1.6+）
	// strings.Compare 返回：-1 (str1 < str2), 0 (str1 == str2), 1 (str1 > str2)
	result := strings.Compare(str1, str2)

	return data.NewIntValue(result), nil
}

// StrcasecmpFunction 实现 strcasecmp 全局函数（不区分大小写）
type StrcasecmpFunction struct{}

// NewStrcasecmpFunction 创建一个新的 strcasecmp 函数实例
func NewStrcasecmpFunction() data.FuncStmt {
	return &StrcasecmpFunction{}
}

// GetName 返回函数名
func (f *StrcasecmpFunction) GetName() string {
	return "strcasecmp"
}

// GetParams 返回参数列表
func (f *StrcasecmpFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str1", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "str2", 1, nil, data.Mixed{}),
	}
}

// GetVariables 返回变量列表
func (f *StrcasecmpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str1", 0, data.Mixed{}),
		node.NewVariable(nil, "str2", 1, data.Mixed{}),
	}
}

// Call 执行 strcasecmp 函数
// strcasecmp 用于二进制安全字符串比较，不区分大小写
func (f *StrcasecmpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 str1 参数
	str1Value, _ := ctx.GetIndexValue(0)
	if str1Value == nil {
		return data.NewIntValue(0), nil
	}

	str1 := ""
	if s, ok := str1Value.(data.AsString); ok {
		str1 = s.AsString()
	} else {
		str1 = str1Value.AsString()
	}

	// 获取 str2 参数
	str2Value, _ := ctx.GetIndexValue(1)
	if str2Value == nil {
		return data.NewIntValue(0), nil
	}

	str2 := ""
	if s, ok := str2Value.(data.AsString); ok {
		str2 = s.AsString()
	} else {
		str2 = str2Value.AsString()
	}

	// 转换为小写后比较
	str1Lower := strings.ToLower(str1)
	str2Lower := strings.ToLower(str2)
	result := strings.Compare(str1Lower, str2Lower)

	return data.NewIntValue(result), nil
}

// StrncasecmpFunction 实现 strncasecmp 全局函数（不区分大小写，按长度比较）
type StrncasecmpFunction struct{}

func NewStrncasecmpFunction() data.FuncStmt {
	return &StrncasecmpFunction{}
}

func (f *StrncasecmpFunction) GetName() string {
	return "strncasecmp"
}

func (f *StrncasecmpFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str1", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "str2", 1, nil, data.Mixed{}),
		node.NewParameter(nil, "length", 2, nil, data.Int{}),
	}
}

func (f *StrncasecmpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str1", 0, data.Mixed{}),
		node.NewVariable(nil, "str2", 1, data.Mixed{}),
		node.NewVariable(nil, "length", 2, data.Int{}),
	}
}

func (f *StrncasecmpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	str1Value, _ := ctx.GetIndexValue(0)
	str2Value, _ := ctx.GetIndexValue(1)
	lenValue, _ := ctx.GetIndexValue(2)

	str1 := ""
	if str1Value != nil {
		str1 = str1Value.AsString()
	}
	str2 := ""
	if str2Value != nil {
		str2 = str2Value.AsString()
	}
	n := 0
	if asInt, ok := lenValue.(data.AsInt); ok {
		if v, err := asInt.AsInt(); err == nil {
			n = v
		}
	}
	if n < 0 {
		n = 0
	}
	if n < len(str1) {
		str1 = str1[:n]
	}
	if n < len(str2) {
		str2 = str2[:n]
	}
	return data.NewIntValue(strings.Compare(strings.ToLower(str1), strings.ToLower(str2))), nil
}
