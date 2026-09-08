package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayReplaceRecursiveFunction() data.FuncStmt {
	return &ArrayReplaceRecursiveFunction{}
}

type ArrayReplaceRecursiveFunction struct{}

func (f *ArrayReplaceRecursiveFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsValue, _ := ctx.GetIndexValue(0)
	if paramsValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsArray, ok := paramsValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsList := paramsArray.ToValueList()
	if len(paramsList) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	result := deepCopyPreserveKeys(paramsList[0])
	for _, replacement := range paramsList[1:] {
		result = recursiveReplaceByKeys(result, replacement)
	}

	return result, nil
}

func recursiveReplaceByKeys(base, replacement data.Value) data.Value {
	baseObj, baseIsObj := base.(*data.ObjectValue)
	replObj, replIsObj := replacement.(*data.ObjectValue)

	if baseIsObj && replIsObj {
		result := data.NewObjectValue()
		baseObj.RangeProperties(func(key string, val data.Value) bool {
			result.SetProperty(key, val)
			return true
		})
		replObj.RangeProperties(func(key string, val data.Value) bool {
			baseVal, _ := baseObj.GetProperty(key)
			if baseVal != nil {
				if _, isNull := baseVal.(*data.NullValue); !isNull {
					if bothArraysOrObjects(baseVal, val) {
						result.SetProperty(key, recursiveReplaceByKeys(baseVal, val))
						return true
					}
				}
			}
			result.SetProperty(key, val)
			return true
		})
		return result
	}

	baseArr, baseIsArr := base.(*data.ArrayValue)
	replArr, replIsArr := replacement.(*data.ArrayValue)
	if baseIsArr && replIsArr {
		out := data.CloneArrayValue(baseArr)
		for i, z := range replArr.List {
			if z == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = data.IntArrayKeyName(i)
			}
			if existing, ok := out.LookupZValByStringKey(key); ok && existing != nil && bothArraysOrObjects(existing.Value, z.Value) {
				existing.Value = recursiveReplaceByKeys(existing.Value, z.Value)
				continue
			}
			setArrayNamedValue(out, key, z.Value)
		}
		return out
	}

	return replacement
}

func bothArraysOrObjects(a, b data.Value) bool {
	_, aArr := a.(*data.ArrayValue)
	_, bArr := b.(*data.ArrayValue)
	if aArr && bArr {
		return true
	}
	_, aObj := a.(*data.ObjectValue)
	_, bObj := b.(*data.ObjectValue)
	return aObj && bObj
}

func deepCopyPreserveKeys(v data.Value) data.Value {
	switch val := v.(type) {
	case *data.ObjectValue:
		result := data.NewObjectValue()
		val.RangeProperties(func(key string, prop data.Value) bool {
			result.SetProperty(key, deepCopyPreserveKeys(prop))
			return true
		})
		return result
	case *data.ArrayValue:
		cloned := data.CloneArrayValue(val)
		for _, z := range cloned.List {
			if z != nil {
				z.Value = deepCopyPreserveKeys(z.Value)
			}
		}
		return cloned
	default:
		return v
	}
}

func (f *ArrayReplaceRecursiveFunction) GetName() string {
	return "array_replace_recursive"
}

func (f *ArrayReplaceRecursiveFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayReplaceRecursiveFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
