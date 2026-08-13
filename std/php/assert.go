package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// AssertFunction 实现 assert()（供 vendor 库中的运行时类型检查使用）。
type AssertFunction struct{}

func NewAssertFunction() data.FuncStmt {
	return &AssertFunction{}
}

func (f *AssertFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cond, _ := ctx.GetIndexValue(0)
	if cond == nil {
		return data.NewBoolValue(false), nil
	}
	ok := false
	if b, err := cond.(data.AsBool).AsBool(); err == nil {
		ok = b
	}
	if !ok {
		msg := "Assertion failed"
		if desc, _ := ctx.GetIndexValue(1); desc != nil {
			if s := desc.AsString(); s != "" {
				msg = s
			}
		}
		return nil, data.NewErrorThrow(nil, data.NewError(nil, msg, nil))
	}
	return data.NewBoolValue(true), nil
}

func (f *AssertFunction) GetName() string { return "assert" }

func (f *AssertFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "assertion", 0, nil, nil),
		node.NewParameter(nil, "description", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *AssertFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "assertion", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "description", 1, data.NewBaseType("mixed")),
	}
}
