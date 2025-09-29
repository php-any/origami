package data

type ArrayValueReverse struct {
	source []Value
}

// Call 实现数组的 reverse 方法
// 反转数组中元素的顺序，并返回反转后的数组
func (a *ArrayValueReverse) Call(ctx Context) (GetValue, Control) {
	// 反转数组
	for i, j := 0, len(a.source)-1; i < j; i, j = i+1, j-1 {
		(a.source)[i], (a.source)[j] = (a.source)[j], (a.source)[i]
	}

	// 返回反转后的数组
	return NewArrayValue(a.source), nil
}

func (a *ArrayValueReverse) GetName() string {
	return "reverse"
}

func (a *ArrayValueReverse) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueReverse) GetIsStatic() bool {
	return false
}

func (a *ArrayValueReverse) GetParams() []GetValue {
	return []GetValue{}
}

func (a *ArrayValueReverse) GetVariables() []Variable {
	return []Variable{}
}

func (a *ArrayValueReverse) GetReturnType() Types {
	return Arrays{}
}
