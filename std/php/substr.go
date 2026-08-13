package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewSubstrFunction() data.FuncStmt {
	return &SubstrFunction{}
}

type SubstrFunction struct{}

func (f *SubstrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	startValue, _ := ctx.GetIndexValue(1)
	lengthValue, _ := ctx.GetIndexValue(2)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	// 获取起始位置
	start := 0
	if startValue != nil {
		if startInt, ok := startValue.(data.AsInt); ok {
			if s, err := startInt.AsInt(); err == nil {
				start = s
			}
		}
	}

	// 处理负数起始位置
	strLen := len(str)
	if start < 0 {
		start = strLen + start
		if start < 0 {
			start = 0
		}
	}

	// 如果起始位置超出范围
	if start >= strLen {
		return data.NewStringValue(""), nil
	}

	// 处理长度参数
	if lengthValue == nil {
		// 没有长度参数，返回从起始位置到末尾
		return data.NewStringValue(str[start:]), nil
	}

	// 检查是否为 NullValue
	if _, ok := lengthValue.(*data.NullValue); ok {
		// 没有长度参数，返回从起始位置到末尾
		return data.NewStringValue(str[start:]), nil
	}

	length := 0
	if lengthInt, ok := lengthValue.(data.AsInt); ok {
		if l, err := lengthInt.AsInt(); err == nil {
			length = l
		}
	}

	// 处理负数长度
	if length < 0 {
		end := strLen + length
		if end < start {
			return data.NewStringValue(""), nil
		}
		return data.NewStringValue(str[start:end]), nil
	}

	// 计算结束位置
	end := start + length
	if end > strLen {
		end = strLen
	}

	return data.NewStringValue(str[start:end]), nil
}

func (f *SubstrFunction) GetName() string {
	return "substr"
}

func (f *SubstrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "start", 1, nil, nil),
		node.NewParameter(nil, "length", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *SubstrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "start", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 2, data.NewBaseType("int")),
	}
}
