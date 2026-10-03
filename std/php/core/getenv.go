package core

import (
	"strings"

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
		return allEnvironmentVariables(ctx), nil
	}
	if _, isNull := nameValue.(*data.NullValue); isNull {
		return allEnvironmentVariables(ctx), nil
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

	// PHP 环境独立于可写的 $_ENV / $_SERVER 数组。
	value, exists := node.LookupEnvVar(ctx, name)
	if !exists {
		// 环境变量不存在，返回 false
		return data.NewBoolValue(false), nil
	}

	// 环境变量存在，返回其值（即使是空字符串也返回）
	return data.NewStringValue(value), nil
}

func allEnvironmentVariables(ctx data.Context) *data.ArrayValue {
	// PHP 7.1+：无参数 getenv() 返回全部环境变量。
	list := make([]*data.ZVal, 0)
	for _, entry := range node.EnvironmentEntries(ctx) {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		zv := data.NewZVal(data.NewStringValue(parts[1]))
		zv.Name = parts[0]
		list = append(list, zv)
	}
	return data.NewArrayValueFromSlots(list)
}

func (f *GetenvFunction) GetName() string {
	return "getenv"
}

var getenvFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewNullableType(data.String{})),
	node.NewParameter(nil, "local_only", 1, data.NewBoolValue(false), data.Bool{}),
}

func (f *GetenvFunction) GetParams() []data.GetValue {
	return getenvFunctionGetParams
}

var getenvFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "name", 0, data.String{}),
	node.NewVariable(nil, "local_only", 1, data.Bool{}),
}

func (f *GetenvFunction) GetVariables() []data.Variable {
	return getenvFunctionGetVariables
}
