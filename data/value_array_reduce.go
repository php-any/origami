package data

type ArrayValueReduce struct {
	source []*ZVal
}

// Call 实现数组的 reduce 方法
// 将数组中的所有元素通过回调函数累积为单个值
func (a *ArrayValueReduce) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewNullValue(), nil
	}

	// 将 source 转换为 []Value 用于 NewArrayValue
	tempArray := &ArrayValue{List: a.source}
	sourceValues := tempArray.ToValueList()

	// 同时支持 *FuncValue 与 CallableValue
	switch callable := callback.(type) {
	case *FuncValue:
		vars := callable.Value.GetVariables()
		fnCtx := ctx.CreateContext(vars)

		// 获取初始值参数
		var accumulator Value
		if initialValue, ok := ctx.GetIndexValue(1); ok {
			accumulator = initialValue
			for i, zval := range a.source {
				element := zval.Value
				args := []Value{accumulator, element, NewIntValue(i), NewArrayValue(sourceValues)}
				for ai := 0; ai < len(vars) && ai < len(args); ai++ {
					fnCtx.SetVariableValue(NewVariable("", ai, nil), args[ai])
				}
				reduceResult, ctl := callable.Value.Call(fnCtx)
				if ctl != nil {
					return nil, ctl
				}
				accumulator = reduceResult.(Value)
			}
			return accumulator, nil
		} else {
			if len(a.source) == 0 {
				return NewNullValue(), nil
			}
			accumulator = a.source[0].Value
			for i := 1; i < len(a.source); i++ {
				element := a.source[i].Value
				args := []Value{accumulator, element, NewIntValue(i), NewArrayValue(sourceValues)}
				for ai := 0; ai < len(vars) && ai < len(args); ai++ {
					fnCtx.SetVariableValue(NewVariable("", ai, nil), args[ai])
				}
				reduceResult, ctl := callable.Value.Call(fnCtx)
				if ctl != nil {
					return nil, ctl
				}
				accumulator = reduceResult.(Value)
			}
			return accumulator, nil
		}

	case CallableValue:
		// 获取初始值参数
		var accumulator Value
		if initialValue, ok := ctx.GetIndexValue(1); ok {
			accumulator = initialValue
		} else {
			// 如果没有提供初始值，使用第一个元素作为初始值
			if len(a.source) == 0 {
				return NewNullValue(), nil
			}
			accumulator = a.source[0].Value
			// 从第二个元素开始遍历
			for i := 1; i < len(a.source); i++ {
				element := a.source[i].Value
				// 调用回调函数，传递累积值、当前元素、索引和数组
				reduceResult, ctl := callable.Call(accumulator, element, NewIntValue(i), NewArrayValue(sourceValues))
				if ctl != nil {
					return nil, ctl
				}
				accumulator = reduceResult
			}
			return accumulator, nil
		}

		// 从第一个元素开始遍历
		for i, zval := range a.source {
			element := zval.Value
			// 调用回调函数，传递累积值、当前元素、索引和数组
			reduceResult, ctl := callable.Call(accumulator, element, NewIntValue(i), NewArrayValue(sourceValues))
			if ctl != nil {
				return nil, ctl
			}
			accumulator = reduceResult
		}

		return accumulator, nil
	}

	return NewNullValue(), nil
}

func (a *ArrayValueReduce) GetName() string {
	return "reduce"
}

func (a *ArrayValueReduce) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueReduce) GetIsStatic() bool {
	return false
}

func (a *ArrayValueReduce) GetParams() []GetValue {
	return []GetValue{
		NewParameter("callback", 0),
		NewParameter("initialValue", 1),
	}
}

func (a *ArrayValueReduce) GetVariables() []Variable {
	return []Variable{
		NewVariable("callback", 0, nil),
		NewVariable("initialValue", 1, nil),
	}
}

func (a *ArrayValueReduce) GetReturnType() Types {
	return NewNullableType(NewBaseType(""))
}
