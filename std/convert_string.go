package std

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStringFunction() data.FuncStmt { return &StringFunction{} }

type StringFunction struct{}

func (f *StringFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}

	// PHP (string) / 字符串强转：对象优先调用 __toString
	switch val := v.(type) {
	case *data.NullValue:
		return data.NewStringValue(""), nil
	case *data.BoolValue:
		if val.Value {
			return data.NewStringValue("1"), nil
		}
		return data.NewStringValue(""), nil
	case *data.StringValue:
		return val, nil
	case *data.ArrayValue:
		return data.NewStringValue("Array"), nil
	case *data.ThisValue:
		if val.ClassValue != nil {
			return stringFromClassValue(val.ClassValue)
		}
		return data.NewStringValue(""), nil
	case *data.ClassValue:
		return stringFromClassValue(val)
	case *data.ObjectValue:
		return data.NewStringValue("Object"), nil
	}

	if s, ok := v.(data.AsString); ok {
		return data.NewStringValue(s.AsString()), nil
	}
	return data.NewStringValue(v.AsString()), nil
}

func stringFromClassValue(val *data.ClassValue) (data.GetValue, data.Control) {
	if method, ok := val.GetMethod("__toString"); ok && method != nil {
		fnCtx := val.CreateContext(method.GetVariables())
		fnCtx.SetCallArgs([]data.GetValue{})
		ret, ctl := method.Call(fnCtx)
		if ctl != nil {
			return nil, ctl
		}
		if sv, ok := ret.(data.Value); ok {
			// 避免 __toString 误返回 $this 时再走 AsString 得到 Object(...)
			if _, isObj := sv.(*data.ClassValue); isObj {
				return data.NewStringValue(""), nil
			}
			if _, isThis := sv.(*data.ThisValue); isThis {
				return data.NewStringValue(""), nil
			}
			return data.NewStringValue(sv.AsString()), nil
		}
		if gv, ok := ret.(data.GetValue); ok {
			v2, ctl2 := gv.GetValue(fnCtx)
			if ctl2 != nil {
				return nil, ctl2
			}
			if sv, ok := v2.(data.Value); ok {
				return data.NewStringValue(sv.AsString()), nil
			}
		}
	}
	return data.NewStringValue("Object"), nil
}

func (f *StringFunction) GetName() string { return "string" }

var stringFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *StringFunction) GetParams() []data.GetValue {
	return stringFunctionGetParams
}

var stringFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
}

func (f *StringFunction) GetVariables() []data.Variable {
	return stringFunctionGetVariables
}
