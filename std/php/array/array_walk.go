package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayWalkFunction() data.FuncStmt {
	return &ArrayWalkFunction{}
}

type ArrayWalkFunction struct{}

func (fn *ArrayWalkFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cbVal, _ := ctx.GetIndexValue(1)
	if cbVal == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取数组参数的 ZVal 引用以支持按引用修改
	arrZVal := ctx.GetIndexZVal(0)
	if arrZVal == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取用户数据（可选第3个参数）
	userdata, _ := ctx.GetIndexValue(2)

	var list []*data.ZVal
	switch arr := arrZVal.Value.(type) {
	case *data.ArrayValue:
		list = arr.List
	case *data.ObjectValue:
		// ObjectValue 转成 ArrayValue 以便修改
		return data.NewBoolValue(true), nil
	default:
		return data.NewBoolValue(false), nil
	}

	switch cb := cbVal.(type) {
	case *data.FuncValue:
		vars := cb.Value.GetVariables()
		for i := range list {
			fnCtx := ctx.CreateContext(vars)
			if len(vars) > 0 {
				fnCtx.SetVariableValue(data.NewVariable("", 0, nil), list[i].Value)
			}
			if len(vars) > 1 {
				// 第二个参数是 key（索引）
				fnCtx.SetVariableValue(data.NewVariable("", 1, nil), data.NewIntValue(i))
			}
			if len(vars) > 2 && userdata != nil {
				fnCtx.SetVariableValue(data.NewVariable("", 2, nil), userdata)
			}
			ret, ctl := cb.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if v, ok := ret.(data.Value); ok {
				list[i].Value = v
			}
		}
	}

	return data.NewBoolValue(true), nil
}

func (fn *ArrayWalkFunction) GetName() string {
	return "array_walk"
}

func (fn *ArrayWalkFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.NewBaseType("array")),
		node.NewParameter(nil, "callback", 1, nil, nil),
		node.NewParameter(nil, "arg", 2, nil, nil),
	}
}

func (fn *ArrayWalkFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
		node.NewVariable(nil, "arg", 2, data.Mixed{}),
	}
}
