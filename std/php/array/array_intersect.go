package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIntersectFunction 实现 PHP 内置函数 array_intersect
// array_intersect(array $array, array ...$arrays): array
// 返回第一个数组中，出现在所有其余数组中的值。
// 当前实现：
// - 对索引数组：返回 ArrayValue（不保留原整数键，仅保证值集合正确）
// - 对关联数组：返回 ObjectValue，并保留第一个数组的字符串键
type ArrayIntersectFunction struct{}

func NewArrayIntersectFunction() data.FuncStmt {
	return &ArrayIntersectFunction{}
}

func (f *ArrayIntersectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取所有参数（打包在 ParametersTODO 中）
	paramsVal, _ := ctx.GetIndexValue(0)
	if paramsVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsArr, ok := paramsVal.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsList := paramsArr.ToValueList()
	if len(paramsList) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 第一个数组作为基准
	baseVal := paramsList[0]

	// 将后续所有数组参数的值集合为 map[string]struct{}，使用字符串值作为比较键
	var otherSets []map[string]struct{}

	for idx := 1; idx < len(paramsList); idx++ {
		v := paramsList[idx]
		if v == nil {
			continue
		}

		set := make(map[string]struct{})

		switch av := v.(type) {
		case *data.ArrayValue:
			for _, z := range av.List {
				set[z.Value.AsString()] = struct{}{}
			}
		case *data.ObjectValue:
			props := av.GetProperties()
			for _, val := range props {
				set[val.AsString()] = struct{}{}
			}
		default:
			// 非数组参数忽略（PHP 会发 warning，这里简化）
		}

		otherSets = append(otherSets, set)
	}

	if len(otherSets) == 0 {
		// 没有其它数组，PHP 会返回第一个数组本身，这里简化为空数组
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 计算交集：第一个数组中的值必须出现在所有其它集合中
	switch v := baseVal.(type) {
	case *data.ArrayValue:
		values := v.ToValueList()
		result := make([]data.Value, 0, len(values))
		for _, val := range values {
			valStr := val.AsString()
			inAll := true
			for _, set := range otherSets {
				if _, ok := set[valStr]; !ok {
					inAll = false
					break
				}
			}
			if inAll {
				result = append(result, val)
			}
		}
		return data.NewArrayValue(result), nil

	case *data.ObjectValue:
		resultObj := data.NewObjectValue()
		// 使用 RangeProperties 保证遍历顺序与插入顺序一致
		v.RangeProperties(func(k string, val data.Value) bool {
			valStr := val.AsString()
			inAll := true
			for _, set := range otherSets {
				if _, ok := set[valStr]; !ok {
					inAll = false
					break
				}
			}
			if inAll {
				resultObj.SetProperty(k, val)
			}
			return true
		})
		return resultObj, nil

	default:
		return data.NewArrayValue([]data.Value{}), nil
	}
}

func (f *ArrayIntersectFunction) GetName() string {
	return "array_intersect"
}

func (f *ArrayIntersectFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		// 使用 ParametersTODO 收集所有数组参数
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayIntersectFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
