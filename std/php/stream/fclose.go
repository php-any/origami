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

	res, ok := streamValue.(*core.ResourceValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	resource := res.GetResource()
	if resource == nil {
		return data.NewBoolValue(false), nil
	}

	var err error
	switch info := resource.(type) {
	case *StreamInfo:
		err = info.Close()
	case *StreamInfoFromReader:
		// proc_open 管道
		err = info.Close()
	default:
		if closer, ok := resource.(interface{ Close() error }); ok {
			err = closer.Close()
		} else {
			return data.NewBoolValue(false), nil
		}
	}
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
