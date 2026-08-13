package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// IsResourceFunction 实现 is_resource 函数
type IsResourceFunction struct{}

func NewIsResourceFunction() data.FuncStmt {
	return &IsResourceFunction{}
}

func (f *IsResourceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查是否为 ResourceValue（它嵌入了 ClassValue）
	if _, ok := value.(*core.ResourceValue); ok {
		return data.NewBoolValue(true), nil
	}

	// 或者检查是否为 ClassValue，并且实现了 Resource 接口
	if classValue, ok := value.(*data.ClassValue); ok {
		if classValue.Class != nil {
			implements := classValue.Class.GetImplements()
			for _, impl := range implements {
				if impl == "Resource" {
					return data.NewBoolValue(true), nil
				}
			}
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *IsResourceFunction) GetName() string {
	return "is_resource"
}

func (f *IsResourceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsResourceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
