package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// GetDebugTypeFunction 实现 get_debug_type 函数
// get_debug_type(mixed $value): string
// 返回变量的调试类型信息（PHP 8.0+）
// 比 gettype() 更详细，能够区分不同的类型
type GetDebugTypeFunction struct{}

func NewGetDebugTypeFunction() data.FuncStmt {
	return &GetDebugTypeFunction{}
}

func (f *GetDebugTypeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "value", 0, data.Mixed{}))
	if ctl != nil {
		return nil, ctl
	}

	internalV, intervalCtl := v.GetValue(ctx)
	if intervalCtl != nil {
		return nil, intervalCtl
	}

	tp := "unknown"
	switch val := internalV.(type) {
	case *data.ArrayValue:
		tp = "array"
	case *data.BoolValue:
		tp = "bool"
	case *data.ClassValue:
		// 对于对象，返回具体的类名
		tp = val.Class.GetName()
	case *data.FloatValue:
		tp = "float"
	case *data.IntValue:
		tp = "int" // get_debug_type 返回 "int" 而不是 "integer"
	case *data.ObjectValue:
		tp = "object"
	case *data.StringValue:
		tp = "string"
	case *data.NullValue:
		tp = "null"
	case *data.AnyValue:
		tp = "any"
	case *core.ResourceValue:
		// 对于资源，返回资源类型和状态
		resourceType := val.GetResourceType()
		resource := val.GetResource()

		// 检查资源是否已关闭
		isClosed := false
		if streamInfo, ok := resource.(interface{ IsClosed() bool }); ok {
			isClosed = streamInfo.IsClosed()
		}

		if resourceType != "" {
			if isClosed {
				tp = "resource (closed)"
			} else {
				tp = "resource (" + resourceType + ")"
			}
		} else {
			if isClosed {
				tp = "resource (closed)"
			} else {
				tp = "resource"
			}
		}
	default:
		// 未知类型，返回 "unknown"
		if val == nil {
			tp = "null"
		} else {
			tp = "unknown"
		}
	}

	return data.NewStringValue(tp), intervalCtl
}

func (f *GetDebugTypeFunction) GetName() string {
	return "get_debug_type"
}

func (f *GetDebugTypeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
	}
}

func (f *GetDebugTypeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
	}
}
