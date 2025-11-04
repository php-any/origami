package data

type ArrayValueFilter struct {
	source []Value
}

// Call 实现数组的 filter 方法
// 遍历数组中的每个元素，调用回调函数，返回所有使回调函数返回 true 的元素组成的新数组
func (a *ArrayValueFilter) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewArrayValue(a.source), nil
	}

	// 创建结果数组
	var result []Value

	switch callable := callback.(type) {
	case *FuncValue:
		vars := callable.Value.GetVariables()
		fnCtx := ctx.CreateContext(vars)
		for i, element := range a.source {
			args := []Value{element, NewIntValue(i), NewArrayValue(a.source)}
			for ai := 0; ai < len(vars) && ai < len(args); ai++ {
				fnCtx.SetVariableValue(NewVariable("", ai, nil), args[ai])
			}
			filterResult, ctl := callable.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			fv := filterResult.(Value)
			if boolResult, ok := fv.(AsBool); ok {
				if isTrue, err := boolResult.AsBool(); err == nil && isTrue {
					result = append(result, element)
				}
			}
		}
		return NewArrayValue(result), nil
	case CallableValue:
		// 遍历数组元素并应用回调函数
		for i, element := range a.source {
			// 调用回调函数，传递元素、索引和数组
			filterResult, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(a.source))
			if ctl != nil {
				return nil, ctl
			}

			// 检查回调函数返回的结果是否为 true
			if boolResult, ok := filterResult.(AsBool); ok {
				if isTrue, err := boolResult.AsBool(); err == nil && isTrue {
					result = append(result, element)
				}
			}
		}
		return NewArrayValue(result), nil
	}

	return NewArrayValue(result), nil
}

func (a *ArrayValueFilter) GetName() string {
	return "filter"
}

func (a *ArrayValueFilter) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueFilter) GetIsStatic() bool {
	return false
}

func (a *ArrayValueFilter) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
	}
}

func (a *ArrayValueFilter) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
	}
}

func (a *ArrayValueFilter) GetReturnType() Types {
	return Arrays{}
}
