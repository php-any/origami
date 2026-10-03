package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ArrayFunction struct{}

func NewArrayFunction() data.FuncStmt { return &ArrayFunction{} }

func (f *ArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	switch v := a1.(type) {
	case *data.ArrayValue:
		return v, nil

	case *data.ClassValue:
		return castClassToArray(v)
	case *data.ThisValue:
		return castClassToArray(v.ClassValue)
	case *data.NullValue:
		// PHP: (array) null => []（空数组，不是 [null]）
		return data.NewArrayValue([]data.Value{}), nil
	case *data.StringValue, *data.IntValue, *data.FloatValue, *data.BoolValue:
		// PHP: (array) 标量 => array(0 => 标量)
		return data.NewArrayValue([]data.Value{a1}), nil
	default:
		return data.NewArrayValue([]data.Value{a1}), nil
	}
}

// Object casts include declared instance properties, with PHP's visibility
// prefixes. Do not change RangeProperties: callers such as foreach need their
// own visibility rules, and ordinary property access uses the unmangled names.
func castClassToArray(object *data.ClassValue) (data.GetValue, data.Control) {
	result := data.NewArrayValueFromSlots(nil)
	if object == nil {
		return result, nil
	}
	type declaredProperty struct {
		property data.Property
		key      string
	}
	indexes := make(map[string]int)
	var declarations []declaredProperty
	var classes []data.ClassStmt
	for class := object.Class; class != nil; {
		classes = append(classes, class)
		extend := class.GetExtend()
		if extend == nil || object.GetVM() == nil {
			break
		}
		parent, ctl := object.GetVM().GetOrLoadClass(*extend)
		if ctl != nil {
			return nil, ctl
		}
		class = parent
	}
	for i := len(classes) - 1; i >= 0; i-- {
		for _, property := range classes[i].GetPropertyList() {
			if property == nil || property.GetIsStatic() {
				continue
			}
			name := property.GetName()
			key := name
			switch property.GetModifier() {
			case data.ModifierProtected:
				key = "\x00*\x00" + name
			case data.ModifierPrivate:
				key = "\x00" + classes[i].GetName() + "\x00" + name
			}
			declaration := declaredProperty{property, key}
			if index, exists := indexes[name]; exists {
				declarations[index] = declaration
			} else {
				indexes[name] = len(declarations)
				declarations = append(declarations, declaration)
			}
		}
	}
	// Declared defaults can be lazy for native classes. Untyped properties with
	// no default are null; typed properties with no value remain absent.
	for _, declaration := range declarations {
		property := declaration.property
		name := property.GetName()
		value, ctl := object.ObjectValue.GetProperty(name)
		if !object.ObjectValue.HasProperty(name) {
			if def := property.GetDefaultValue(); def != nil {
				var raw data.GetValue
				raw, ctl = def.GetValue(object)
				value, _ = raw.(data.Value)
			} else if property.GetType() != nil {
				continue
			}
		}
		if ctl != nil {
			return nil, ctl
		}
		if value == nil {
			value = data.NewNullValue()
		}
		result.SetStringKey(declaration.key, value)
	}
	object.RangeProperties(func(key string, value data.Value) bool {
		if _, declared := indexes[key]; !declared {
			result.SetStringKey(key, value)
		}
		return true
	})
	return result, nil
}

func (f *ArrayFunction) GetName() string { return "array" }

var arrayFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "data", 0, nil, nil),
}

func (f *ArrayFunction) GetParams() []data.GetValue {
	return arrayFunctionGetParams
}

var arrayFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "data", 0, data.Mixed{}),
}

func (f *ArrayFunction) GetVariables() []data.Variable {
	return arrayFunctionGetVariables
}
