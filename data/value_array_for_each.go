package data

type ArrayValueForEach struct {
	source []Value
}

// Call 实现数组的 forEach 方法
// 对数组中的每个元素执行一次提供的回调函数
func (a *ArrayValueForEach) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewNullValue(), nil
	}

	// 检查回调函数是否可调用
	callable, ok := callback.(CallableValue)
	if !ok {
		return NewNullValue(), nil
	}

	// 遍历数组元素
	for i, element := range a.source {
		// 调用回调函数，传递元素、索引和数组
		_, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
		if ctl != nil {
			return nil, ctl
		}
	}

	return NewNullValue(), nil
}

func (a *ArrayValueForEach) GetName() string {
	return "forEach"
}

func (a *ArrayValueForEach) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueForEach) GetIsStatic() bool {
	return false
}

func (a *ArrayValueForEach) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueForEach) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}
