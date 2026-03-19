package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetClassFunction 实现 get_class 函数
// get_class(?object $object = null): string|false
// 返回对象的类名。如果 object 不是对象则返回 false。
// 如果省略了参数，函数返回当前类名。
type GetClassFunction struct{}

func NewGetClassFunction() data.FuncStmt {
	return &GetClassFunction{}
}

func (f *GetClassFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数（可选）
	objectValue, hasObject := ctx.GetIndexValue(0)

	// 如果没有参数，尝试从上下文中获取 $this
	if !hasObject || objectValue == nil {
		// 尝试获取 $this
		thisVar := node.NewVariable(nil, "this", 0, nil)
		thisValue, ctl := thisVar.GetValue(ctx)
		if ctl == nil && thisValue != nil {
			if classValue, ok := thisValue.(*data.ClassValue); ok {
				// 返回当前类的类名
				return data.NewStringValue(classValue.Class.GetName()), nil
			}
		}
		// 如果没有 $this 或不是对象，返回 false
		return data.NewBoolValue(false), nil
	}

	// 检查是否为 null
	if _, isNull := objectValue.(*data.NullValue); isNull {
		// 如果参数是 null，尝试从上下文中获取 $this
		thisVar := node.NewVariable(nil, "this", 0, nil)
		thisValue, ctl := thisVar.GetValue(ctx)
		if ctl == nil && thisValue != nil {
			if classValue, ok := thisValue.(*data.ClassValue); ok {
				// 返回当前类的类名
				return data.NewStringValue(classValue.Class.GetName()), nil
			}
		}
		// 如果没有 $this 或不是对象，返回 false
		return data.NewBoolValue(false), nil
	}

	// 检查是否为 ClassValue（对象）
	if classValue, ok := objectValue.(data.GetName); ok {
		// 返回对象的类名
		return data.NewStringValue(classValue.GetName()), nil
	}

	// 不是对象，返回 false
	return data.NewBoolValue(false), nil
}

func (f *GetClassFunction) GetName() string {
	return "get_class"
}

func (f *GetClassFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, node.NewNullLiteral(nil), nil),
	}
}

func (f *GetClassFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.NewBaseType("object|null")),
	}
}
