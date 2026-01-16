package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetenvFunction 实现 getenv 函数
// 获取环境变量的值
type GetenvFunction struct{}

func NewGetenvFunction() data.FuncStmt {
	return &GetenvFunction{}
}

func (f *GetenvFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数：环境变量名
	nameValue, _ := ctx.GetIndexValue(0)
	if nameValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 将参数转换为字符串
	var name string
	if str, ok := nameValue.(data.AsString); ok {
		name = str.AsString()
	} else {
		name = nameValue.AsString()
	}

	// 检查环境变量名是否为空
	if name == "" {
		return data.NewBoolValue(false), nil
	}

	// 获取环境变量的值，使用 LookupEnv 来区分"不存在"和"值为空字符串"
	value, exists := os.LookupEnv(name)
	if !exists {
		// 环境变量不存在，返回 false
		return data.NewBoolValue(false), nil
	}

	// 环境变量存在，返回其值（即使是空字符串也返回）
	return data.NewStringValue(value), nil
}

func (f *GetenvFunction) GetName() string {
	return "getenv"
}

func (f *GetenvFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

func (f *GetenvFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}
