package data

type ArrayValueFindIndex struct {
	source []Value
}

// Call 实现数组的 findIndex 方法
// 返回数组中第一个满足回调函数条件的元素的索引，如果没有找到则返回 -1
func (a *ArrayValueFindIndex) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewIntValue(-1), nil
	}

	// 检查回调函数是否可调用
	callable, ok := callback.(CallableValue)
	if !ok {
		return NewIntValue(-1), nil
	}

	// 遍历数组元素并查找第一个满足条件的元素的索引
	for i, element := range a.source {
		// 调用回调函数，传递元素、索引和数组
		testResult, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
		if ctl != nil {
			return nil, ctl
		}

		// 检查回调函数返回的结果是否为 true
		if boolResult, ok := testResult.(AsBool); ok {
			if isTrue, err := boolResult.AsBool(); err == nil && isTrue {
				return NewIntValue(i), nil
			}
		}
	}

	return NewIntValue(-1), nil
}

func (a *ArrayValueFindIndex) GetName() string {
	return "findIndex"
}

func (a *ArrayValueFindIndex) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueFindIndex) GetIsStatic() bool {
	return false
}

func (a *ArrayValueFindIndex) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueFindIndex) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}

func (a *ArrayValueFindIndex) GetReturnType() Types {
	return Int{}
}
