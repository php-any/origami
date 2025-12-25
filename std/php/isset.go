package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIssetFunction() data.FuncStmt {
	return &IssetFunction{}
}

type IssetFunction struct{}

func (f *IssetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	params := f.GetParams()
	if len(params) == 0 {
		return data.NewBoolValue(false), nil
	}

	// 获取第一个参数的值
	varParam := params[0]
	varValue, ctl := varParam.GetValue(ctx)
	if ctl != nil {
		return data.NewBoolValue(false), nil
	}

	switch varValue.(type) {
	case *data.NullValue:
		return data.NewBoolValue(false), nil
	}

	// 值已设置且不为 null
	return data.NewBoolValue(true), nil
}

func (f *IssetFunction) GetName() string {
	return "isset"
}

func (f *IssetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "var", 0, nil, data.Mixed{}),
	}
}

func (f *IssetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "var", 0, data.Mixed{}),
	}
}
