package data

type ArrayValuePop struct {
	source *[]Value
}

// Call 实现数组的 pop 方法
// 移除并返回数组的最后一个元素，如果数组为空则返回 null
func (a *ArrayValuePop) Call(ctx Context) (GetValue, Control) {
	if len(*a.source) == 0 {
		return NewNullValue(), nil
	}

	// 获取并移除最后一个元素
	lastElement := (*a.source)[len(*a.source)-1]
	*a.source = (*a.source)[:len(*a.source)-1]

	return lastElement, nil
}

func (a *ArrayValuePop) GetName() string {
	return "pop"
}

func (a *ArrayValuePop) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValuePop) GetIsStatic() bool {
	return false
}

func (a *ArrayValuePop) GetParams() []GetValue {
	return []GetValue{}
}

func (a *ArrayValuePop) GetVariables() []Variable {
	return []Variable{}
}
