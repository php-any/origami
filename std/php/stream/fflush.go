package stream

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// FflushFunction 实现 fflush 函数
type FflushFunction struct{}

func NewFflushFunction() data.FuncStmt {
	return &FflushFunction{}
}

func (f *FflushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	streamValue, _ := ctx.GetIndexValue(0)
	if streamValue == nil {
		return data.NewBoolValue(false), nil
	}

	var streamInfo *StreamInfo
	if res, ok := streamValue.(*core.ResourceValue); ok {
		resource := res.GetResource()
		if resource == nil {
			return data.NewBoolValue(false), nil
		}
		if info, ok := resource.(*StreamInfo); ok {
			streamInfo = info
		} else {
			return data.NewBoolValue(false), nil
		}
	} else {
		return data.NewBoolValue(false), nil
	}

	if streamInfo.IsClosed() {
		return data.NewBoolValue(false), nil
	}

	if err := streamInfo.Flush(); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *FflushFunction) GetName() string {
	return "fflush"
}

func (f *FflushFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
	}
}

func (f *FflushFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
	}
}
