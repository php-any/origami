package data

import "errors"

type ArrayValueMap struct {
	source []Value
}

// Call 实现数组的 map 方法
// 创建一个新数组，其结果是该数组中的每个元素调用一次提供的回调函数后的返回值
func (a *ArrayValueMap) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewArrayValue(a.source), nil
	}

	// 创建结果数组
	result := make([]Value, len(a.source))

	switch callable := callback.(type) {
	case *FuncValue:
		// 使用函数定义的变量创建调用上下文，并按顺序写入参数：element, index, array
		vars := callable.Value.GetVariables()
		fnCtx := ctx.CreateContext(vars)
		for i, element := range a.source {
			args := []Value{element, NewIntValue(i), NewArrayValue(a.source)}
			for ai := 0; ai < len(vars) && ai < len(args); ai++ {
				fnCtx.SetVariableValue(NewVariable("", ai, nil), args[ai])
			}
			ret, ctl := callable.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			result[i] = ret.(Value)
		}

		return NewArrayValue(result), nil

	case CallableValue:
		// 遍历数组元素并应用回调函数
		for i, element := range a.source {
			// 调用回调函数，传递元素、索引和数组
			mappedValue, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
			if ctl != nil {
				return nil, ctl
			}
			result[i] = mappedValue
		}

		return NewArrayValue(result), nil
	}

	return nil, NewErrorThrow(nil, errors.New("call func failed"))
}

func (a *ArrayValueMap) GetName() string {
	return "map"
}

func (a *ArrayValueMap) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueMap) GetIsStatic() bool {
	return false
}

func (a *ArrayValueMap) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueMap) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}

func (a *ArrayValueMap) GetReturnType() Types {
	return Arrays{}
}
