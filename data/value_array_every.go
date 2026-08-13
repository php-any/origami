package data

type ArrayValueEvery struct {
	source []*ZVal
}

// Call 实现数组的 every 方法
// 检查数组中的所有元素是否都满足回调函数的条件，如果所有元素都满足则返回 true，否则返回 false
func (a *ArrayValueEvery) Call(ctx Context) (GetValue, Control) {
	// 获取回调函数参数
	callback, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewBoolValue(true), nil
	}

	// 将 source 转换为 []Value 用于 NewArrayValue
	tempArray := &ArrayValue{List: a.source}
	sourceValues := tempArray.ToValueList()

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
			testResult, ctl := callable.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			tr := testResult.(Value)
			if boolResult, ok := tr.(AsBool); ok {
				if isTrue, err := boolResult.AsBool(); err != nil || !isTrue {
					return NewBoolValue(false), nil
				}
			} else {
				return NewBoolValue(false), nil
			}
		}
		return NewBoolValue(true), nil
	case CallableValue:
		// 遍历数组元素并检查是否都满足条件
		for i, zval := range a.source {
			element := zval.Value
			// 调用回调函数，传递元素、索引和数组
			testResult, ctl := callable.Call(element, NewIntValue(i), NewArrayValue(sourceValues))
			if ctl != nil {
				return nil, ctl
			}
			if boolResult, ok := testResult.(AsBool); ok {
				if isTrue, err := boolResult.AsBool(); err != nil || !isTrue {
					return NewBoolValue(false), nil
				}
			} else {
				return NewBoolValue(false), nil
			}
		}
		return NewBoolValue(true), nil
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
