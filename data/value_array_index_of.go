package data

type ArrayValueIndexOf struct {
	source []Value
}

// Call 实现数组的 indexOf 方法
// 返回数组中第一个与指定元素相等的元素的索引，如果没找到则返回 -1
func (a *ArrayValueIndexOf) Call(ctx Context) (GetValue, Control) {
	// 获取要查找的元素
	searchElement, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewIntValue(-1), nil
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
		return NewIntValue(-1), nil
	}

	// 查找元素
	for i := fromIndex; i < len(a.source); i++ {
		if a.source[i].AsString() == searchElement.AsString() {
			return NewIntValue(i), nil
		}
	}

	return NewIntValue(-1), nil
}

func (a *ArrayValueIndexOf) GetName() string {
	return "indexOf"
}

func (a *ArrayValueIndexOf) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueIndexOf) GetIsStatic() bool {
	return false
}

func (a *ArrayValueIndexOf) GetParams() []GetValue {
	return []GetValue{
		NewParameter("searchElement", 0),
		NewParameter("fromIndex", 1),
	}
}

func (a *ArrayValueIndexOf) GetVariables() []Variable {
	return []Variable{
		NewVariable("searchElement", 0, nil),
		NewVariable("fromIndex", 1, nil),
	}
}
