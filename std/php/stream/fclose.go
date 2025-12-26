package stream

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// FcloseFunction 实现 fclose 函数
type FcloseFunction struct{}

func NewFcloseFunction() data.FuncStmt {
	return &FcloseFunction{}
}

func (f *FcloseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取流资源
	streamValue, _ := ctx.GetIndexValue(0)
	if streamValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 从资源对象中获取 StreamInfo
	var streamInfo *StreamInfo
	if res, ok := streamValue.(*core.ResourceValue); ok {
		resource := res.GetResource()
		if info, ok := resource.(*StreamInfo); ok {
			streamInfo = info
		} else {
			return data.NewBoolValue(false), nil
		}
	} else {
		return data.NewBoolValue(false), nil
	}

	// 关闭流
	err := streamInfo.Close()
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(true), nil
}

func (f *FcloseFunction) GetName() string {
	return "fclose"
}

func (f *FcloseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
	}
}

func (f *FcloseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
	}
}
