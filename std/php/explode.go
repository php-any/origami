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

	// 处理 limit 参数（遵循 PHP explode 语义）
	hasLimit := false
	limit := 0
	if limitValue != nil {
		if _, ok := limitValue.(*data.NullValue); !ok {
			if limitInt, ok := limitValue.(data.AsInt); ok {
				if l, err := limitInt.AsInt(); err == nil {
					limit = l
					hasLimit = true
				}
			}
		}
	}

	// 分割字符串
	var parts []string
	if !hasLimit {
		// 未提供 limit 或为 null：分割所有
		parts = strings.Split(str, separator)
	} else if limit > 0 {
		// 正数 limit：与 PHP 一致，限制返回元素数量
		parts = strings.SplitN(str, separator, limit)
	} else if limit == 0 {
		// limit = 0 时视为 1，与 PHP 行为一致
		parts = strings.Split(str, separator)
	} else { // limit < 0
		// 负数 limit：返回除最后 -limit 个元素之外的所有元素
		all := strings.Split(str, separator)
		if len(all) == 1 {
			// PHP 特例：未找到分隔符且 limit<0 时返回空数组
			parts = []string{}
		} else {
			n := len(all) + limit // limit 为负数
			if n < 0 {
				n = 0
			}
			parts = all[:n]
		}
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
