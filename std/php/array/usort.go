package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UsortFunction 实现 usort 函数
// usort(array &$array, callable $callback): bool
// 使用用户自定义的比较函数对数组进行排序
type UsortFunction struct{}

func NewUsortFunction() data.FuncStmt {
	return &UsortFunction{}
}

func (f *UsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	callbackValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil || callbackValue == nil {
		return data.NewBoolValue(false), nil
	}

	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if len(arrayRef.List) <= 1 {
		return data.NewBoolValue(true), nil
	}

	// 获取回调函数变量信息
	var callbackVars []data.Variable
	switch cb := callbackValue.(type) {
	case *data.FuncValue:
		callbackVars = cb.Value.GetVariables()
	}

	sort.Slice(arrayRef.List, func(i, j int) bool {
		// 调用回调函数
		fnCtx := ctx.CreateContext(callbackVars)
		if len(callbackVars) > 0 {
			fnCtx.SetIndexZVal(0, data.NewZVal(arrayRef.List[i].Value))
		}
		if len(callbackVars) > 1 {
			fnCtx.SetIndexZVal(1, data.NewZVal(arrayRef.List[j].Value))
		}

		switch cb := callbackValue.(type) {
		case *data.FuncValue:
			ret, ctl := cb.Call(fnCtx)
			if ctl != nil {
				if rv, ok := ctl.(data.ReturnControl); ok {
					if v, ok := rv.ReturnValue().(data.Value); ok {
						if iv, ok := v.(*data.IntValue); ok {
							return iv.Value < 0
						}
						if fv, ok := v.(*data.FloatValue); ok {
							return fv.Value < 0
						}
					}
				}
				return false
			}
			if ret != nil {
				if v, ok := ret.(data.Value); ok {
					if iv, ok := v.(*data.IntValue); ok {
						return iv.Value < 0
					}
					if fv, ok := v.(*data.FloatValue); ok {
						return fv.Value < 0
					}
				}
			}
		}
		return false
	})

	// 重新索引（清空字符串键名，整数键从0开始）
	for i, zval := range arrayRef.List {
		zval.Name = ""
		_ = i
	}

	return data.NewBoolValue(true), nil
}

func (f *UsortFunction) GetName() string {
	return "usort"
}

func (f *UsortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
		node.NewParameter(nil, "callback", 1, nil, data.Mixed{}),
	}
}

func (f *UsortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
	}
}
