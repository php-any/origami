package data

type ArrayValueReverse struct {
	source *ArrayValue
}

// Call 实现数组的 reverse 方法
// 反转数组中元素的顺序，并返回反转后的数组
func (a *ArrayValueReverse) Call(ctx Context) (GetValue, Control) {
	// 反转数组
	slots := a.source.slots()
	for i, j := 0, len(slots)-1; i < j; i, j = i+1, j-1 {
		slots[i], slots[j] = slots[j], slots[i]
	}

	a.source.ReplaceAll(slots)

	// 返回反转后的数组
	values := make([]Value, len(slots))
	for i, zval := range slots {
		values[i] = zval.Value
	}
	return NewArrayValue(values), nil
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
