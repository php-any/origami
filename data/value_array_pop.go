package data

type ArrayValuePop struct {
	source *ArrayValue
}

// Call 实现数组的 pop 方法
// 移除并返回数组的最后一个元素，如果数组为空则返回 null
func (a *ArrayValuePop) Call(ctx Context) (GetValue, Control) {
	slot := a.source.PopSlot()
	if slot == nil {
		return NewNullValue(), nil
	}
	return slot.ReadValue(), nil
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

func (a *ArrayValuePop) GetReturnType() Types {
	return nil
}
