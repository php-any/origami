package stream

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/utils"
)

// FwriteFunction 实现 fwrite 函数
type FwriteFunction struct{}

func NewFwriteFunction() data.FuncStmt {
	return &FwriteFunction{}
}

func (f *FwriteFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取流资源
	streamValue, _ := ctx.GetIndexValue(0)
	if streamValue == nil {
		return data.NewIntValue(0), nil
	}

	// 从资源对象中获取 StreamInfo
	var streamInfo *StreamInfo
	if res, ok := streamValue.(*core.ResourceValue); ok {
		// 使用 ResourceValue 的 GetResource 方法
		resource := res.GetResource()
		if resource == nil {
			return data.NewIntValue(0), nil
		}
		if info, ok := resource.(*StreamInfo); ok {
			streamInfo = info
		} else {
			// 类型断言失败，resource 不是 *StreamInfo 类型
			return data.NewIntValue(0), nil
		}
	} else {
		// 不是 ResourceValue 类型
		return data.NewIntValue(0), nil
	}

	// 检查流是否已关闭
	if streamInfo.IsClosed() {
		return data.NewIntValue(0), nil
	}

	// 获取要写入的数据
	dataValue, _ := ctx.GetIndexValue(1)
	if _, ok := dataValue.(*data.BoolValue); ok {
		return nil, utils.NewThrow(errors.New("debug"))
	}
	if dataValue == nil {
		return data.NewIntValue(0), nil
	}

	var content string
	if s, ok := dataValue.(data.AsString); ok {
		content = s.AsString()
	} else {
		content = dataValue.AsString()
	}

	if content == "" {
		return data.NewIntValue(0), nil
	}

	// 获取可选的 length 参数
	var maxLength int = len(content)
	lengthValue, ok := ctx.GetIndexValue(2)
	if ok {
		if _, ok := lengthValue.(*data.NullValue); !ok {
			if intVal, ok := lengthValue.(data.AsInt); ok {
				if length, err := intVal.AsInt(); err == nil && length >= 0 {
					if length < maxLength {
						maxLength = length
					}
				}
			}
		}
	}

	// 写入数据
	bytes := []byte(content[:maxLength])
	if streamInfo == nil {
		return data.NewIntValue(0), nil
	}
	if streamInfo.File == nil {
		return data.NewIntValue(0), nil
	}

	n, err := streamInfo.Write(bytes)
	if err != nil {
		return data.NewIntValue(0), nil
	}

	return data.NewIntValue(n), nil
}

func (f *FwriteFunction) GetName() string {
	return "fwrite"
}

func (f *FwriteFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, nil),
		node.NewParameter(nil, "length", 2, nil, nil),
	}
}

func (f *FwriteFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
		node.NewVariable(nil, "data", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "length", 2, data.NewBaseType("int")),
	}
}
