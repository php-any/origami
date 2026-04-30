package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayReplaceFunction() data.FuncStmt {
	return &ArrayReplaceFunction{}
}

type ArrayReplaceFunction struct{}

func (f *ArrayReplaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// array_replace: 用后面的数组覆盖前面的数组中对应的键
	// 与 array_merge 不同，int 键也会被覆盖而不是重新编号

	// 判断第一个参数类型
	first := paramsList[0]

	switch base := first.(type) {
	case *data.ObjectValue:
		// 关联数组：创建结果副本
		result := data.NewObjectValue()
		base.RangeProperties(func(key string, val data.Value) bool {
			result.SetProperty(key, val)
			return true
		})

		// 用后续数组覆盖
		for i := 1; i < len(paramsList); i++ {
			switch v := paramsList[i].(type) {
			case *data.ObjectValue:
				v.RangeProperties(func(key string, val data.Value) bool {
					result.SetProperty(key, val)
					return true
				})
			case *data.ArrayValue:
				values := v.ToValueList()
				for j, val := range values {
					result.SetProperty(fmt.Sprintf("%d", j), val)
				}
			}
		}
		return result, nil

	case *data.ArrayValue:
		// 纯列表数组：转为 ObjectValue 以支持键覆盖
		result := data.NewObjectValue()
		values := base.ToValueList()
		for i, val := range values {
			result.SetProperty(fmt.Sprintf("%d", i), val)
		}

		// 用后续数组覆盖
		for i := 1; i < len(paramsList); i++ {
			switch v := paramsList[i].(type) {
			case *data.ObjectValue:
				v.RangeProperties(func(key string, val data.Value) bool {
					result.SetProperty(key, val)
					return true
				})
			case *data.ArrayValue:
				vals := v.ToValueList()
				for j, val := range vals {
					result.SetProperty(fmt.Sprintf("%d", j), val)
				}
			}
		}
		return result, nil

	default:
		return data.NewArrayValue([]data.Value{}), nil
	}
}

func (f *ArrayReplaceFunction) GetName() string {
	return "array_replace"
}

func (f *ArrayReplaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayReplaceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
