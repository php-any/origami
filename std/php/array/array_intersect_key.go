package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIntersectKeyFunction implements array_intersect_key
type ArrayIntersectKeyFunction struct{}

func NewArrayIntersectKeyFunction() data.FuncStmt {
	return &ArrayIntersectKeyFunction{}
}

// extractKeys 从任意 PHP 数组值中提取所有 key 的集合
func extractKeys(v data.Value) map[string]bool {
	keys := make(map[string]bool)
	switch arr := v.(type) {
	case *data.ArrayValue:
		for idx, zv := range arr.List {
			if zv.Name != "" {
				keys[zv.Name] = true
			} else {
				keys[data.NewIntValue(idx).AsString()] = true
			}
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, _ data.Value) bool {
			keys[key] = true
			return true
		})
	}
	return keys
}

func (f *ArrayIntersectKeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	firstVal, _ := ctx.GetIndexValue(0)
	if firstVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 收集其余数组的 key 集合（index 1 是 variadic 打包的数组）
	var allKeys []map[string]bool
	if restVal, ok := ctx.GetIndexValue(1); ok && restVal != nil {
		if restArr, ok := restVal.(*data.ArrayValue); ok {
			for _, zv := range restArr.List {
				allKeys = append(allKeys, extractKeys(zv.Value))
			}
		}
	}

	// 辅助：检查 key 是否在所有集合中
	inAll := func(key string) bool {
		for _, ks := range allKeys {
			if !ks[key] {
				return false
			}
		}
		return true
	}

	// 根据 firstVal 类型构建结果
	switch first := firstVal.(type) {
	case *data.ArrayValue:
		result := data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
		for idx, zv := range first.List {
			key := zv.Name
			if key == "" {
				key = fmt.Sprintf("%d", idx)
			}
			if inAll(key) {
				result.List = append(result.List, &data.ZVal{Name: zv.Name, Value: zv.Value})
			}
		}
		return result, nil
	case *data.ObjectValue:
		result := data.NewObjectValue()
		first.RangeProperties(func(key string, val data.Value) bool {
			if inAll(key) {
				result.SetProperty(key, val)
			}
			return true
		})
		return result, nil
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}
}

func (f *ArrayIntersectKeyFunction) GetName() string { return "array_intersect_key" }

func (f *ArrayIntersectKeyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameters(nil, "arrays", 1, nil, nil),
	}
}

func (f *ArrayIntersectKeyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "arrays", 1, data.NewBaseType("array")),
	}
}
