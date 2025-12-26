package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewSubstrCountFunction() data.FuncStmt {
	return &SubstrCountFunction{}
}

type SubstrCountFunction struct{}

func (f *SubstrCountFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystackValue, _ := ctx.GetIndexValue(0)
	needleValue, _ := ctx.GetIndexValue(1)
	offsetValue, _ := ctx.GetIndexValue(2)
	lengthValue, _ := ctx.GetIndexValue(3)

	if haystackValue == nil || needleValue == nil {
		return data.NewIntValue(0), nil
	}

	// 转换为字符串
	haystack := haystackValue.AsString()
	needle := needleValue.AsString()

	// 如果 needle 为空字符串，返回 0
	if needle == "" {
		return data.NewIntValue(0), nil
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
	haystackLen := len(haystack)
	if offset < 0 {
		offset = haystackLen + offset
		if offset < 0 {
			offset = 0
		}
	}

	// 检查偏移量是否超出范围
	if offset >= haystackLen {
		return data.NewIntValue(0), nil
	}

	// 处理长度参数
	searchLen := haystackLen - offset
	if lengthValue != nil {
		// 检查是否是 null
		if _, isNull := lengthValue.(*data.NullValue); !isNull {
			if lengthInt, ok := lengthValue.(data.AsInt); ok {
				if l, err := lengthInt.AsInt(); err == nil {
					searchLen = l
					// 处理负数长度
					if searchLen < 0 {
						searchLen = haystackLen + searchLen - offset
						if searchLen < 0 {
							searchLen = 0
						}
					}
					// 确保不超过剩余长度
					if searchLen > haystackLen-offset {
						searchLen = haystackLen - offset
					}
				}
			}
		}
	}

	// 如果搜索长度为 0 或负数，返回 0
	if searchLen <= 0 {
		return data.NewIntValue(0), nil
	}

	// 获取要搜索的子字符串
	searchStr := haystack[offset : offset+searchLen]

	// 计算子字符串出现的次数
	count := strings.Count(searchStr, needle)

	return data.NewIntValue(count), nil
}

func (f *SubstrCountFunction) GetName() string {
	return "substr_count"
}

func (f *SubstrCountFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "length", 3, node.NewNullLiteral(nil), nil),
	}
}

func (f *SubstrCountFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 3, data.NewBaseType("int|null")),
	}
}
