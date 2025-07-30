package data

type ArrayValueIncludes struct {
	source []Value
}

// Call 实现数组的 includes 方法
// 判断数组是否包含指定的元素，如果包含则返回 true，否则返回 false
func (a *ArrayValueIncludes) Call(ctx Context) (GetValue, Control) {
	// 获取要查找的元素
	searchElement, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewBoolValue(false), nil
	}

	// 获取起始索引参数
	fromIndex := 0
	if fromIndexArg, ok := ctx.GetIndexValue(1); ok {
		if fromIndexInt, ok := fromIndexArg.(AsInt); ok {
			if f, err := fromIndexInt.AsInt(); err == nil {
				fromIndex = f
			}
		}
	}

	// 处理负数索引
	if fromIndex < 0 {
		fromIndex = len(a.source) + fromIndex
	}

	// 边界检查
	if fromIndex < 0 {
		fromIndex = 0
	}
	if fromIndex >= len(a.source) {
		return NewBoolValue(false), nil
	}

	// 查找元素
	for i := fromIndex; i < len(a.source); i++ {
		if a.source[i].AsString() == searchElement.AsString() {
			return NewBoolValue(true), nil
		}
	}

	return NewBoolValue(false), nil
}

func (a *ArrayValueIncludes) GetName() string {
	return "includes"
}

func (a *ArrayValueIncludes) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueIncludes) GetIsStatic() bool {
	return false
}

func (a *ArrayValueIncludes) GetParams() []GetValue {
	return []GetValue{
		NewParameter("searchElement", 0),
		NewParameter("fromIndex", 1),
	}
}

func (a *ArrayValueIncludes) GetVariables() []Variable {
	return []Variable{
		NewVariable("searchElement", 0, nil),
		NewVariable("fromIndex", 1, nil),
	}
}
