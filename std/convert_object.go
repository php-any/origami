package std

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ObjectFunction 实现 PHP 的 (object) 类型转换。
// PHP：(object)$array 得到 stdClass 实例（属性为数组键）；Blade $loop = (object)$last 依赖此语义。
type ObjectFunction struct{}

func NewObjectFunction() data.FuncStmt { return &ObjectFunction{} }

func (f *ObjectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return newStdClassInstance(ctx), nil
	}

	switch val := v.(type) {
	case *data.ClassValue:
		// 已是对象实例，原样返回（与 PHP 一致）
		return val, nil
	case *data.ThisValue:
		return val.ClassValue, nil
	case *data.ObjectValue:
		// Origami 内部关联数组容器：转为真正的 stdClass
		return objectValueToStdClass(ctx, val), nil
	case *data.ArrayValue:
		return arrayToStdClass(ctx, val), nil
	default:
		// 标量：包装为带 scalar 属性的 stdClass（对齐 PHP）
		obj := newStdClassInstance(ctx)
		if cv, ok := obj.(*data.ClassValue); ok {
			cv.SetProperty("scalar", v)
		}
		return obj, nil
	}
}

func newStdClassInstance(ctx data.Context) data.GetValue {
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewObjectValue()
	}
	cls, ok := vm.GetClass("stdClass")
	if !ok {
		return data.NewObjectValue()
	}
	return data.NewClassValue(cls, ctx.CreateBaseContext())
}

func arrayToStdClass(ctx data.Context, arr *data.ArrayValue) data.GetValue {
	obj := newStdClassInstance(ctx)
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		// 回退：无 stdClass 时用 ObjectValue
		ov := data.NewObjectValue()
		for i, z := range arr.List {
			key := fmt.Sprintf("%d", i)
			if z != nil && z.Name != "" {
				key = z.Name
			}
			if z == nil || z.Value == nil {
				ov.SetProperty(key, data.NewNullValue())
			} else {
				ov.SetProperty(key, z.Value)
			}
		}
		return ov
	}
	for i, z := range arr.List {
		key := fmt.Sprintf("%d", i)
		if z != nil && z.Name != "" {
			key = z.Name
		}
		if z == nil || z.Value == nil {
			cv.SetProperty(key, data.NewNullValue())
		} else {
			cv.SetProperty(key, z.Value)
		}
	}
	return cv
}

func objectValueToStdClass(ctx data.Context, ov *data.ObjectValue) data.GetValue {
	obj := newStdClassInstance(ctx)
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return ov
	}
	ov.RangeProperties(func(key string, val data.Value) bool {
		if val == nil {
			cv.SetProperty(key, data.NewNullValue())
		} else {
			cv.SetProperty(key, val)
		}
		return true
	})
	return cv
}

func (f *ObjectFunction) GetName() string { return "object" }

var objectFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *ObjectFunction) GetParams() []data.GetValue {
	return objectFunctionGetParams
}

var objectFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
}

func (f *ObjectFunction) GetVariables() []data.Variable {
	return objectFunctionGetVariables
}
