package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewRtrimFunction() data.FuncStmt {
	return &RtrimFunction{}
}

type RtrimFunction struct{}

func (f *RtrimFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	charlistValue, _ := ctx.GetIndexValue(1)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	// 处理字符列表参数
	if charlistValue != nil {
		// 检查是否为 NullValue
		if _, ok := charlistValue.(*data.NullValue); !ok {
			charlist := charlistValue.AsString()
			if charlist != "" {
				// 使用自定义字符列表，只去除右侧
				return data.NewStringValue(strings.TrimRight(str, charlist)), nil
			}
		}
	}

	// 默认去除右侧空白字符（包括空格、制表符、换行符等）
	return data.NewStringValue(strings.TrimRight(str, " \t\n\r\v\f")), nil
}

func (f *RtrimFunction) GetName() string {
	return "rtrim"
}

func (f *RtrimFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "charlist", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *RtrimFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "charlist", 1, data.NewBaseType("string")),
	}
}
