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
			for arraySlots95, arrayPosition95 := paramsArr.View(), 0; arrayPosition95 < arraySlots95.Len(); arrayPosition95++ {
				z := arraySlots95.At(arrayPosition95)
				if z == nil || z.ReadValue() == nil {
					continue
				}
				rawArrays = append(rawArrays, z.ReadValue())
			}
		}
	}

	if len(rawArrays) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}
	for _, value := range rawArrays {
		if _, ok := value.(*data.ArrayValue); !ok {
			return nil, throwMustBeArray("array_map", value)
		}
	}

	// PHP：只传入一个数组时保留全部键（含字符串键与稀疏整数键）。
	// ComponentAttributeBag::merge 对 assoc 默认数组做 array_map 后 array_merge；
	// 若重编号成 0/1/2，HTML 会变成 class="" 0="fi-icon-btn …"。
	if len(rawArrays) == 1 {
		switch src := rawArrays[0].(type) {

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

	// PHP pads shorter arrays with null up to the longest input.
	minLen := len(arrayLists[0])
	for _, a := range arrayLists[1:] {
		if len(a) > minLen {
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
			} else {
				args = append(args, data.NewNullValue())
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
	out := make([]*data.ZVal, 0, av.Len())
	for arraySlots96, arrayPosition96 := av.View(), 0; arrayPosition96 < arraySlots96.Len(); arrayPosition96++ {
		z := arraySlots96.At(arrayPosition96)
		if z == nil {
			out = append(out, nil)
			continue
		}
		arg := z.ReadValue()
		if arg == nil {
			arg = data.NewNullValue()
		}
		mapped, ctl := f.invokeCallback(ctx, cbVal, []data.Value{arg})
		if ctl != nil {
			return nil, ctl
		}
		out = append(out, data.CopyZValKeepName(z, mapped))
	}
	return data.NewArrayValueFromSlots(out), nil
}

func (f *ArrayMapFunction) invokeCallback(ctx data.Context, cbVal data.Value, args []data.Value) (data.Value, data.Control) {
	if _, isNull := cbVal.(*data.NullValue); isNull {
		if len(args) == 1 {
			return args[0], nil
		}
		return data.NewArrayValue(args), nil
	}
	switch cb := cbVal.(type) {
	case *data.BoundFuncValue:
		return f.callFuncStmt(ctx, cb.Value, args, cb)
	case *data.FuncValue:
		return f.callFuncStmt(ctx, cb.Value, args, nil)
	case *data.ArrayValue:
		// PHP 数组可调用: [$obj, 'method']
		if cb.Len() == 2 {
			objVal := cb.At(0).ReadValue()
			methodVal := cb.At(1).ReadValue()
			if obj, ok := objVal.(data.GetMethod); ok {
				methodName := methodVal.AsString()
				if method, has := obj.GetMethod(methodName); has {
					varies := method.GetVariables()
					fnCtx := ctx.CreateContext(varies)
					if instance, ok := objVal.(*data.ClassValue); ok {
						fnCtx = data.WrapMethodFrame(fnCtx, instance, instance.Class, instance.Class)
					}
					defer data.ReleaseContext(fnCtx)
					fnCtx.SetStrictTypes(false)
					if ctl := data.BindDeclaredArgs(fnCtx, method, args); ctl != nil {
						return nil, ctl
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
		return f.callFuncStmt(ctx, fnStmt, args, nil)
	}
}

func (f *ArrayMapFunction) callFuncStmt(ctx data.Context, fn data.FuncStmt, args []data.Value, bound *data.BoundFuncValue) (data.Value, data.Control) {
	vars := fn.GetVariables()
	fnCtx := ctx.CreateContext(vars)
	defer data.ReleaseContext(fnCtx)
	fnCtx.SetStrictTypes(false)
	if ctl := data.BindDeclaredArgs(fnCtx, fn, args); ctl != nil {
		return nil, ctl
	}
	var ret data.GetValue
	var ctl data.Control
	if bound != nil {
		ret, ctl = bound.Call(fnCtx)
	} else {
		ret, ctl = fn.Call(fnCtx)
	}
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
