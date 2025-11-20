package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewExplodeFunction() data.FuncStmt {
	return &ExplodeFunction{}
}

type ExplodeFunction struct{}

func (f *ExplodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	separatorValue, _ := ctx.GetIndexValue(0)
	stringValue, _ := ctx.GetIndexValue(1)
	limitValue, _ := ctx.GetIndexValue(2)

	if separatorValue == nil || stringValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	separator := separatorValue.AsString()
	str := stringValue.AsString()

	// 处理 limit 参数
	limit := -1
	if limitValue != nil {
		if _, ok := limitValue.(*data.NullValue); !ok {
			if limitInt, ok := limitValue.(data.AsInt); ok {
				if l, err := limitInt.AsInt(); err == nil {
					limit = l
				}
			}
		}
	}

	// 分割字符串
	var parts []string
	if limit < 0 {
		// 没有限制，分割所有
		parts = strings.Split(str, separator)
	} else if limit == 0 {
		// limit 为 0，返回包含原字符串的数组
		parts = []string{str}
	} else {
		// 限制分割次数
		parts = strings.SplitN(str, separator, limit)
	}

	// 转换为 Value 数组
	var values []data.Value
	for _, part := range parts {
		values = append(values, data.NewStringValue(part))
	}

	return data.NewArrayValue(values), nil
}

func (f *ExplodeFunction) GetName() string {
	return "explode"
}

func (f *ExplodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "separator", 0, nil, nil),
		node.NewParameter(nil, "string", 1, nil, nil),
		node.NewParameter(nil, "limit", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *ExplodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "separator", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "string", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "limit", 2, data.NewBaseType("int")),
	}
}
