package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewInArrayFunction() data.FuncStmt {
	return &InArrayFunction{
		params: []data.GetValue{
			node.NewParameter(nil, "needle", 0, nil, nil),
			node.NewParameter(nil, "haystack", 1, nil, nil),
			node.NewParameter(nil, "strict", 2, node.NewNullLiteral(nil), nil),
		},
		vars: []data.Variable{
			node.NewVariable(nil, "needle", 0, data.NewBaseType("mixed")),
			node.NewVariable(nil, "haystack", 1, data.NewBaseType("array")),
			node.NewVariable(nil, "strict", 2, data.NewBaseType("bool")),
		},
	}
}

type InArrayFunction struct {
	params []data.GetValue
	vars   []data.Variable
}

func (f *InArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	needleValue, _ := ctx.GetIndexValue(0)
	haystackValue, _ := ctx.GetIndexValue(1)
	strictValue, _ := ctx.GetIndexValue(2)

	if haystackValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 收集所有值（支持 ArrayValue / ObjectValue / ClassValue）
	var valueList []data.Value
	if arrayVal, ok := haystackValue.(*data.ArrayValue); ok {
		valueList = arrayVal.ToValueList()
	} else if classVal, ok := haystackValue.(*data.ClassValue); ok {
		classVal.RangeProperties(func(key string, v data.Value) bool {
			valueList = append(valueList, v)
			return true
		})
	} else {
		return data.NewBoolValue(false), nil
	}

	if needleValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 处理严格模式
	strict := false
	if strictValue != nil {
		if _, ok := strictValue.(*data.NullValue); !ok {
			if strictBool, ok := strictValue.(data.AsBool); ok {
				if s, err := strictBool.AsBool(); err == nil {
					strict = s
				}
			}
		}
	}

	// Enum cases are singleton objects. Their diagnostic AsString() only names
	// the enum type, so it cannot distinguish two cases of the same enum.
	needleEnum := inArrayEnumObject(needleValue)

	// 在数组中查找
	for _, val := range valueList {
		if strict {
			if valueEqualStrict(needleValue, val) {
				return data.NewBoolValue(true), nil
			}
		} else {
			if needleEnum != nil {
				if candidate := inArrayEnumObject(val); candidate != nil && needleEnum.ObjectValue == candidate.ObjectValue {
					return data.NewBoolValue(true), nil
				}
				continue
			}
			if inArrayEnumObject(val) != nil {
				continue
			}
			if needleValue.AsString() == val.AsString() {
				return data.NewBoolValue(true), nil
			}
		}
	}

	return data.NewBoolValue(false), nil
}

// EnumParser represents unit and backed enums as subclasses of BackedEnum.
func inArrayEnumObject(value data.Value) *data.ClassValue {
	var object *data.ClassValue
	switch v := value.(type) {
	case *data.ClassValue:
		object = v
	case *data.ThisValue:
		object = v.ClassValue
	default:
		return nil
	}
	if object == nil || object.Class == nil {
		return nil
	}
	if metadata, ok := object.Class.(interface{ DeclarationFlags() data.ClassFlags }); ok && metadata.DeclarationFlags()&data.ClassEnum != 0 {
		return object
	}
	parent := object.Class.GetExtend()
	if parent != nil && (*parent == "BackedEnum" || *parent == "\\BackedEnum" || *parent == "UnitEnum" || *parent == "\\UnitEnum") {
		return object
	}
	return nil
}

func (f *InArrayFunction) GetName() string {
	return "in_array"
}

func (f *InArrayFunction) GetParams() []data.GetValue {
	return f.params
}

func (f *InArrayFunction) GetVariables() []data.Variable {
	return f.vars
}

// valueEqualStrict 实现 PHP === 语义
func valueEqualStrict(a, b data.Value) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch va := a.(type) {
	case *data.StringValue:
		if vb, ok := b.(*data.StringValue); ok {
			return va.Value == vb.Value
		}
	case *data.IntValue:
		if vb, ok := b.(*data.IntValue); ok {
			return va.Value == vb.Value
		}
	case *data.FloatValue:
		if vb, ok := b.(*data.FloatValue); ok {
			return va.Value == vb.Value
		}
	case *data.BoolValue:
		if vb, ok := b.(*data.BoolValue); ok {
			return va.Value == vb.Value
		}
	case *data.NullValue:
		_, ok := b.(*data.NullValue)
		return ok
	default:
		return a == b
	}
	return false
}
