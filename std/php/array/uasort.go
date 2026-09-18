package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UasortFunction 实现 uasort 函数
// uasort(array &$array, callable $callback): bool
// 使用用户自定义比较函数排序，并保持键名关联。
type UasortFunction struct{}

func NewUasortFunction() data.FuncStmt {
	return &UasortFunction{}
}

func (f *UasortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	var callbackVars []data.Variable
	switch cb := callbackValue.(type) {
	case *data.FuncValue:
		callbackVars = cb.Value.GetVariables()
	}

	sort.SliceStable(arrayRef.List, func(i, j int) bool {
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
						return compareLess(v)
					}
				}
				return false
			}
			if ret != nil {
				if v, ok := ret.(data.Value); ok {
					return compareLess(v)
				}
			}
		}
		return false
	})

	// uasort 保留键名，不重索引
	return data.NewBoolValue(true), nil
}

func compareLess(v data.Value) bool {
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value < 0
	}
	if fv, ok := v.(*data.FloatValue); ok {
		return fv.Value < 0
	}
	if as, ok := v.(data.AsInt); ok {
		if i, err := as.AsInt(); err == nil {
			return i < 0
		}
	}
	return false
}

func (f *UasortFunction) GetName() string {
	return "uasort"
}

func (f *UasortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "callback", 1, nil, data.Mixed{}),
	}
}

func (f *UasortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
	}
}
