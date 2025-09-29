package data

type ArrayValueFlatMap struct {
	source []Value
}

// Call 实现数组的 flatMap 方法
// 首先使用映射函数映射每个元素，然后将结果扁平化一层，返回一个新数组
func (a *ArrayValueFlatMap) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewArrayValue(a.source), nil
	}

	// 检查回调函数是否可调用
	callable, ok := callback.(CallableValue)
	if !ok {
		return NewArrayValue(a.source), nil
	}

	// 创建结果数组
	var result []Value

	// 遍历数组元素并应用回调函数
	for i, element := range a.source {
		// 调用回调函数，传递元素、索引和数组
		mappedValue, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
		if ctl != nil {
			return nil, ctl
		}

		// 如果映射结果是数组，则展开一层
		if arrayResult, ok := mappedValue.(*ArrayValue); ok {
			result = append(result, arrayResult.Value...)
		} else {
			result = append(result, mappedValue.(Value))
		}
	}

	return NewArrayValue(result), nil
}

func (a *ArrayValueFlatMap) GetName() string {
	return "flatMap"
}

func (a *ArrayValueFlatMap) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueFlatMap) GetIsStatic() bool {
	return false
}

func (a *ArrayValueFlatMap) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueFlatMap) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}

func (a *ArrayValueFlatMap) GetReturnType() Types {
	return Arrays{}
}
