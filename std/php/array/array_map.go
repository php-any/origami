package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayMapFunction 实现 PHP 内置函数 array_map
// array_map(callable $callback, array $array, array ...$arrays): array
type ArrayMapFunction struct{}

func NewArrayMapFunction() data.FuncStmt {
	return &ArrayMapFunction{}
}

// toValueList 将数组值（可能是 ArrayValue 或 ObjectValue）统一转换为 []data.Value 切片
// ObjectValue 使用 RangeProperties 保证顺序
func toValueList(v data.Value) []data.Value {
	switch arr := v.(type) {
	case *data.ArrayValue:
		return arr.ToValueList()
	case *data.ObjectValue:
		result := make([]data.Value, 0)
		arr.RangeProperties(func(key string, value data.Value) bool {
			result = append(result, value)
			return true
		})
		return result
	default:
		return nil
	}
}

func (f *ArrayMapFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cbVal, has := ctx.GetIndexValue(0)
	if !has || cbVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 收集所有数组参数（支持 ArrayValue 和 ObjectValue）
	paramsVal, _ := ctx.GetIndexValue(1)
	var arrayLists [][]data.Value
	if paramsVal != nil {
		if paramsArr, ok := paramsVal.(*data.ArrayValue); ok {
			for _, z := range paramsArr.List {
				if z.Value == nil {
					continue
				}
				vals := toValueList(z.Value)
				if vals != nil {
					arrayLists = append(arrayLists, vals)
				}
			}
		}
	}

	if len(arrayLists) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 计算最短数组长度
	minLen := len(arrayLists[0])
	for _, a := range arrayLists[1:] {
		if len(a) < minLen {
			minLen = len(a)
		}
	}
	if minLen == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	results := make([]data.Value, 0, minLen)

	for i := 0; i < minLen; i++ {
		args := make([]data.Value, 0, len(arrayLists))
		for _, arr := range arrayLists {
			if i < len(arr) {
				args = append(args, arr[i])
			}
		}

		switch cb := cbVal.(type) {
		case *data.FuncValue:
			vars := cb.Value.GetVariables()
			fnCtx := ctx.CreateContext(vars)
			for ai := 0; ai < len(vars) && ai < len(args); ai++ {
				fnCtx.SetVariableValue(data.NewVariable("", ai, nil), args[ai])
			}
			ret, ctl := cb.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if v, ok := ret.(data.Value); ok {
				results = append(results, v)
			} else {
				results = append(results, data.NewNullValue())
			}
		case *data.ArrayValue:
			// PHP 数组可调用: [$obj, 'method']
			if len(cb.List) == 2 {
				objVal := cb.List[0].Value
				methodVal := cb.List[1].Value
				if obj, ok := objVal.(data.GetMethod); ok {
					methodName := methodVal.AsString()
					if method, has := obj.GetMethod(methodName); has {
						varies := method.GetVariables()
						fnCtx := ctx.CreateContext(varies)
						for ai := 0; ai < len(varies) && ai < len(args); ai++ {
							fnCtx.SetVariableValue(varies[ai], args[ai])
						}
						ret, ctl := method.Call(fnCtx)
						if ctl != nil {
							return nil, ctl
						}
						if v, ok := ret.(data.Value); ok {
							results = append(results, v)
						} else {
							results = append(results, data.NewNullValue())
						}
						continue
					}
				}
			}
			results = append(results, data.NewNullValue())
		case data.CallableValue:
			var arg0, arg1, arg2 data.Value = data.NewNullValue(), data.NewNullValue(), data.NewNullValue()
			if len(args) > 0 {
				arg0 = args[0]
			}
			if len(args) > 1 {
				arg1 = args[1]
			}
			if len(args) > 2 {
				arg2 = args[2]
			}
			ret, ctl := cb.Call(arg0, arg1, arg2)
			if ctl != nil {
				return nil, ctl
			}
			results = append(results, ret)
		default:
			results = append(results, data.NewNullValue())
		}
	}

	return data.NewArrayValue(results), nil
}

func (f *ArrayMapFunction) GetName() string {
	return "array_map"
}

func (f *ArrayMapFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
		node.NewParameters(nil, "arrays", 1, nil, nil),
	}
}

func (f *ArrayMapFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
		node.NewVariable(nil, "arrays", 1, data.Mixed{}),
	}
}
