package data

type ArrayValueEvery struct {
	source []Value
}

// Call 实现数组的 every 方法
// 检查数组中的所有元素是否都满足回调函数的条件，如果所有元素都满足则返回 true，否则返回 false
func (a *ArrayValueEvery) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewBoolValue(true), nil
	}

	// 检查回调函数是否可调用
	callable, ok := callback.(CallableValue)
	if !ok {
		return NewBoolValue(true), nil
	}

	// 遍历数组元素并检查是否都满足条件
	for i, element := range a.source {
		// 调用回调函数，传递元素、索引和数组
		testResult, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
		if ctl != nil {
			return nil, ctl
		}

		// 检查回调函数返回的结果是否为 true
		if boolResult, ok := testResult.(AsBool); ok {
			if isTrue, err := boolResult.AsBool(); err != nil || !isTrue {
				return NewBoolValue(false), nil
			}
		} else {
			// 如果返回值不是布尔类型，默认为 false
			return NewBoolValue(false), nil
		}
	}

	return NewBoolValue(true), nil
}

func (a *ArrayValueEvery) GetName() string {
	return "every"
}

func (a *ArrayValueEvery) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueEvery) GetIsStatic() bool {
	return false
}

func (a *ArrayValueEvery) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueEvery) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}

func (a *ArrayValueEvery) GetReturnType() Types {
	return Bool{}
}
