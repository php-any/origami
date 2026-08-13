package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// GetResourceTypeFunction 实现 get_resource_type 函数
// 对资源返回其资源类型字符串；非资源返回 "Unknown"
type GetResourceTypeFunction struct{}

func NewGetResourceTypeFunction() data.FuncStmt {
	return &GetResourceTypeFunction{}
}

func (f *GetResourceTypeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewStringValue("unknown"), nil
	}

	// 直接是 ResourceValue 的情况
	if res, ok := value.(*core.ResourceValue); ok {
		tp := res.GetResourceType()
		if tp == "" {
			return data.NewStringValue("unknown"), nil
		}
		return data.NewStringValue(tp), nil
	}

	// ClassValue 且实现 Resource 接口的情况
	if classValue, ok := value.(*data.ClassValue); ok {
		if classValue.Class != nil {
			if resourceClass, ok := classValue.Class.(*core.ResourceClass); ok {
				tp := resourceClass.GetResourceType()
				if tp == "" {
					return data.NewStringValue("unknown"), nil
				}
				return data.NewStringValue(tp), nil
			}
		}
	}

	return data.NewStringValue("unknown"), nil
}

func (f *GetResourceTypeFunction) GetName() string {
	return "get_resource_type"
}

func (f *GetResourceTypeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *GetResourceTypeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
