package data

type ArrayValueSome struct {
	source []Value
}

// Call 实现数组的 some 方法
// 检查数组中是否至少有一个元素满足回调函数的条件，如果有则返回 true，否则返回 false
func (a *ArrayValueSome) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewBoolValue(false), nil
	}

	// 检查回调函数是否可调用
	callable, ok := callback.(CallableValue)
	if !ok {
		return NewBoolValue(false), nil
	}

	// 遍历数组元素并检查是否有元素满足条件
	for i, element := range a.source {
		// 调用回调函数，传递元素、索引和数组
		testResult, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
		if ctl != nil {
			return nil, ctl
		}

		// 检查回调函数返回的结果是否为 true
		if boolResult, ok := testResult.(AsBool); ok {
			if isTrue, err := boolResult.AsBool(); err == nil && isTrue {
				return NewBoolValue(true), nil
			}
		}
	}

	return NewBoolValue(false), nil
}

func (a *ArrayValueSome) GetName() string {
	return "some"
}

func (a *ArrayValueSome) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueSome) GetIsStatic() bool {
	return false
}

func (a *ArrayValueSome) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueSome) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}

func (a *ArrayValueSome) GetReturnType() Types {
	return Bool{}
}
