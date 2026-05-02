package array

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayDiffUkeyFunction 实现 array_diff_ukey 函数
type ArrayDiffUkeyFunction struct{}

func NewArrayDiffUkeyFunction() data.FuncStmt {
	return &ArrayDiffUkeyFunction{}
}

func (f *ArrayDiffUkeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrValue, _ := ctx.GetIndexValue(0)
	if arrValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 收集第一个数组的键
	var baseKeys map[string]bool
	switch v := arrValue.(type) {
	case *data.ObjectValue:
		baseKeys = make(map[string]bool)
		v.RangeProperties(func(key string, _ data.Value) bool {
			baseKeys[key] = true
			return true
		})
	case *data.ArrayValue:
		baseKeys = make(map[string]bool)
		for i, z := range v.List {
			_ = i
			baseKeys[z.Value.AsString()] = true
		}
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取回调函数
	callbackVal, _ := ctx.GetIndexValue(1)
	if callbackVal == nil {
		return arrValue, nil
	}

	// 收集后续数组的键（用于排除）
	// 取所有后续数组的键的并集
	allOtherKeys := make(map[string]bool)
	for j := 2; ; j++ {
		otherVal, ok := ctx.GetIndexValue(j)
		if !ok || otherVal == nil {
			break
		}
		switch v := otherVal.(type) {
		case *data.ObjectValue:
			v.RangeProperties(func(key string, _ data.Value) bool {
				allOtherKeys[key] = true
				return true
			})
		case *data.ArrayValue:
			for _, z := range v.List {
				allOtherKeys[z.Value.AsString()] = true
			}
		}
	}

	// 使用回调函数比较键
	result := data.NewObjectValue()
	switch v := arrValue.(type) {
	case *data.ObjectValue:
		v.RangeProperties(func(key string, val data.Value) bool {
			found := false
			for otherKey := range allOtherKeys {
				// 调用回调: callback(key, otherKey)
				if callCmp(callbackVal, ctx, key, otherKey) == 0 {
					found = true
					break
				}
			}
			if !found {
				result.SetProperty(key, val)
			}
			return true
		})
	case *data.ArrayValue:
		// For indexed arrays, convert to ObjectValue result
		for _, z := range v.List {
			key := z.Value.AsString()
			found := false
			for otherKey := range allOtherKeys {
				if callCmp(callbackVal, ctx, key, otherKey) == 0 {
					found = true
					break
				}
			}
			if !found {
				result.SetProperty(key, z.Value)
			}
		}
	}

	return result, nil
}

func callCmp(callbackVal data.Value, ctx data.Context, a, b string) int {
	if fv, ok := callbackVal.(*data.FuncValue); ok {
		// 创建调用上下文
		fnCtx := ctx.CreateContext(fv.Value.GetVariables())
		params := fv.Value.GetParams()
		if len(params) >= 2 {
			if v, ok := params[0].(data.Variable); ok {
				v.SetValue(fnCtx, data.NewStringValue(strings.ToLower(a)))
			}
			if v, ok := params[1].(data.Variable); ok {
				v.SetValue(fnCtx, data.NewStringValue(strings.ToLower(b)))
			}
		}
		ret, ctl := fv.Value.Call(fnCtx)
		if ctl != nil {
			return 1
		}
		if iv, ok := ret.(data.AsInt); ok {
			i, _ := iv.AsInt()
			return i
		}
	}
	return strings.Compare(a, b)
}

func (f *ArrayDiffUkeyFunction) GetName() string {
	return "array_diff_ukey"
}

func (f *ArrayDiffUkeyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, data.Mixed{}),
	}
}

func (f *ArrayDiffUkeyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.Mixed{}),
	}
}
