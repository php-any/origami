package data

type ArrayValueFlatMap struct {
	source []*ZVal
}

// Call 实现数组的 flatMap 方法
// 首先使用映射函数映射每个元素，然后将结果扁平化一层，返回一个新数组
func (a *ArrayValueFlatMap) Call(ctx Context) (GetValue, Control) {
	// 将 source 转换为 []Value 用于 NewArrayValue
	tempArray := &ArrayValue{List: a.source}
	sourceValues := tempArray.ToValueList()

	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewArrayValue(sourceValues), nil
	}

	// 创建结果数组
	var result []Value

	switch callable := callback.(type) {
	case *FuncValue:
		vars := callable.Value.GetVariables()
		fnCtx := ctx.CreateContext(vars)
		for i, zval := range a.source {
			element := zval.Value
			args := []Value{element, NewIntValue(i), NewArrayValue(sourceValues)}
			for ai := 0; ai < len(vars) && ai < len(args); ai++ {
				fnCtx.SetVariableValue(NewVariable("", ai, nil), args[ai])
			}
			mappedValue, ctl := callable.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			mv := mappedValue.(Value)
			if arrayResult, ok := mv.(*ArrayValue); ok {
				// 将 List 转换为 []Value
				result = append(result, arrayResult.ToValueList()...)
			} else {
				result = append(result, mv)
			}
		}
		return NewArrayValue(result), nil
	case CallableValue:
		// 遍历数组元素并应用回调函数
		for i, zval := range a.source {
			element := zval.Value
			// 调用回调函数，传递元素、索引和数组
			mappedValue, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(sourceValues))
			if ctl != nil {
				return nil, ctl
			}
			// 如果映射结果是数组，则展开一层
			if arrayResult, ok := mappedValue.(*ArrayValue); ok {
				// 将 List 转换为 []Value
				result = append(result, arrayResult.ToValueList()...)
			} else {
				result = append(result, mappedValue)
			}
		}
		return NewArrayValue(result), nil
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
