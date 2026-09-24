package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayMapFunction 实现 PHP 内置函数 array_map
// array_map(callable $callback, array $array, array ...$arrays): array
//
// PHP 语义：仅传入一个数组时保留键；多个数组时结果为顺序整数键。
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
	var rawArrays []data.Value
	if paramsVal != nil {
		if paramsArr, ok := paramsVal.(*data.ArrayValue); ok {
			for _, z := range paramsArr.List {
				if z == nil || z.Value == nil {
					continue
				}
				rawArrays = append(rawArrays, z.Value)
			}
		}
	}

	if len(rawArrays) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// PHP：只传入一个数组时保留全部键（含字符串键与稀疏整数键）。
	// ComponentAttributeBag::merge 对 assoc 默认数组做 array_map 后 array_merge；
	// 若重编号成 0/1/2，HTML 会变成 class="" 0="fi-icon-btn …"。
	if len(rawArrays) == 1 {
		switch src := rawArrays[0].(type) {
		case *data.ObjectValue:
			return f.mapObjectPreserveKeys(ctx, cbVal, src)
		case *data.ArrayValue:
			return f.mapArrayPreserveKeys(ctx, cbVal, src)
		}
	}

	arrayLists := make([][]data.Value, 0, len(rawArrays))
	for _, raw := range rawArrays {
		vals := toValueList(raw)
		if vals != nil {
			arrayLists = append(arrayLists, vals)
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

		mapped, ctl := f.invokeCallback(ctx, cbVal, args)
		if ctl != nil {
			return nil, ctl
		}
		results = append(results, mapped)
	}

	return data.NewArrayValue(results), nil
}

func (f *ArrayMapFunction) mapArrayPreserveKeys(ctx data.Context, cbVal data.Value, av *data.ArrayValue) (data.GetValue, data.Control) {
	if av == nil {
		return data.NewArrayValue(nil), nil
	}
	out := make([]*data.ZVal, 0, len(av.List))
	for _, z := range av.List {
		if z == nil {
			out = append(out, nil)
			continue
		}
		arg := z.Value
		if arg == nil {
			arg = data.NewNullValue()
		}
		mapped, ctl := f.invokeCallback(ctx, cbVal, []data.Value{arg})
		if ctl != nil {
			return nil, ctl
		}
		out = append(out, data.CopyZValKeepName(z, mapped))
	}
	return &data.ArrayValue{List: out}, nil
}

func (f *ArrayMapFunction) mapObjectPreserveKeys(ctx data.Context, cbVal data.Value, ov *data.ObjectValue) (data.GetValue, data.Control) {
	out := data.NewObjectValue()
	var failed data.Control
	ov.RangeProperties(func(key string, value data.Value) bool {
		mapped, ctl := f.invokeCallback(ctx, cbVal, []data.Value{value})
		if ctl != nil {
			failed = ctl
			return false
		}
		out.SetProperty(key, mapped)
		return true
	})
	if failed != nil {
		return nil, failed
	}
	return out, nil
}

func (f *ArrayMapFunction) invokeCallback(ctx data.Context, cbVal data.Value, args []data.Value) (data.Value, data.Control) {
	switch cb := cbVal.(type) {
	case *data.BoundFuncValue:
		return f.callFuncStmt(ctx, cb.Value, args)
	case *data.FuncValue:
		return f.callFuncStmt(ctx, cb.Value, args)
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
						return v, nil
					}
					return data.NewNullValue(), nil
				}
			}
		}
		return data.NewNullValue(), nil
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
		return ret, nil
	default:
		funcName := cbVal.AsString()
		fnStmt, exists := ctx.GetVM().GetFunc(funcName)
		if !exists {
			return data.NewNullValue(), nil
		}
		return f.callFuncStmt(ctx, fnStmt, args)
	}
}

func (f *ArrayMapFunction) callFuncStmt(ctx data.Context, fn data.FuncStmt, args []data.Value) (data.Value, data.Control) {
	vars := fn.GetVariables()
	fnCtx := ctx.CreateContext(vars)
	for ai := 0; ai < len(vars) && ai < len(args); ai++ {
		fnCtx.SetVariableValue(data.NewVariable("", ai, nil), args[ai])
	}
	ret, ctl := fn.Call(fnCtx)
	if ctl != nil {
		return nil, ctl
	}
	if v, ok := ret.(data.Value); ok {
		return v, nil
	}
	return data.NewNullValue(), nil
}

func (f *ArrayMapFunction) GetName() string {
	return "array_map"
}

var arrayMapFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "callback", 0, nil, nil),
	node.NewParameters(nil, "arrays", 1, nil, nil),
}

func (f *ArrayMapFunction) GetParams() []data.GetValue {
	return arrayMapFunctionGetParams
}

var arrayMapFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "callback", 0, data.Mixed{}),
	node.NewVariable(nil, "arrays", 1, data.Mixed{}),
}

func (f *ArrayMapFunction) GetVariables() []data.Variable {
	return arrayMapFunctionGetVariables
}
