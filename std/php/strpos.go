package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrposFunction() data.FuncStmt {
	return &StrposFunction{}
}

type StrposFunction struct{}

func (f *StrposFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystackValue, _ := ctx.GetIndexValue(0)
	needleValue, _ := ctx.GetIndexValue(1)
	offsetValue, _ := ctx.GetIndexValue(2)

	if haystackValue == nil || needleValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 转换为字符串
	haystack := haystackValue.AsString()
	needle := needleValue.AsString()

	// 如果 needle 为空字符串，返回 false
	if needle == "" {
		return data.NewBoolValue(false), nil
	}

	// 处理偏移量
	offset := 0
	if offsetValue != nil {
		if offsetInt, ok := offsetValue.(data.AsInt); ok {
			if o, err := offsetInt.AsInt(); err == nil {
				offset = o
			}
		}
	}

	// 如果偏移量为负数，从末尾开始计算
	if offset < 0 {
		offset = len(haystack) + offset
		if offset < 0 {
			offset = 0
		}
	}

	// 检查偏移量是否超出范围
	if offset >= len(haystack) {
		return data.NewBoolValue(false), nil
	}

	// 查找子字符串位置
	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return data.NewBoolValue(false), nil
	}

	return data.NewIntValue(pos + offset), nil
}

func (f *StrposFunction) GetName() string {
	return "strpos"
}

func (f *StrposFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *StrposFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
	}
}
