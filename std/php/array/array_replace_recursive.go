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

	// 第一个数组是基础数组
	base := paramsList[0]
	replacements := paramsList[1:]

	result := deepCopy(base)
	for _, replacement := range replacements {
		result = recursiveReplace(result, replacement)
	}

	return result, nil
}

func recursiveReplace(base, replacement data.Value) data.Value {
	baseObj, baseIsObj := base.(*data.ObjectValue)
	replObj, replIsObj := replacement.(*data.ObjectValue)

	if baseIsObj && replIsObj {
		result := data.NewObjectValue()
		for key, val := range baseObj.GetProperties() {
			result.SetProperty(key, val)
		}
		for key, val := range replObj.GetProperties() {
			baseVal, _ := baseObj.GetProperty(key)
			_, isNull := baseVal.(*data.NullValue)
			if !isNull {
				if _, isArr := baseVal.(*data.ArrayValue); isArr {
					if _, isArr2 := val.(*data.ArrayValue); isArr2 {
						result.SetProperty(key, recursiveReplace(baseVal, val))
						continue
					}
				}
				if _, isObj := baseVal.(*data.ObjectValue); isObj {
					if _, isObj2 := val.(*data.ObjectValue); isObj2 {
						result.SetProperty(key, recursiveReplace(baseVal, val))
						continue
					}
				}
			}
			result.SetProperty(key, val)
		}
		return result
	}

	baseArr, baseIsArr := base.(*data.ArrayValue)
	replArr, replIsArr := replacement.(*data.ArrayValue)
	if baseIsArr && replIsArr {
		baseVals := baseArr.ToValueList()
		replVals := replArr.ToValueList()
		result := make([]data.Value, len(baseVals))
		copy(result, baseVals)
		for i, val := range replVals {
			if i < len(result) {
				result[i] = recursiveReplace(result[i], val)
			} else {
				result = append(result, val)
			}
		}
		return data.NewArrayValue(result)
	}

	return replacement
}

func deepCopy(v data.Value) data.Value {
	switch val := v.(type) {
	case *data.ObjectValue:
		result := data.NewObjectValue()
		for key, prop := range val.GetProperties() {
			result.SetProperty(key, deepCopy(prop))
		}
		return result
	case *data.ArrayValue:
		vals := val.ToValueList()
		result := make([]data.Value, len(vals))
		for i, item := range vals {
			result[i] = deepCopy(item)
		}
		return data.NewArrayValue(result)
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
