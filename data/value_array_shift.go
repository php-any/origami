package data

type ArrayValueShift struct {
	source *ArrayValue
}

// Call 实现数组的 shift 方法
// 移除并返回数组的第一个元素，如果数组为空则返回 null
func (a *ArrayValueShift) Call(ctx Context) (GetValue, Control) {
	slot := a.source.ShiftSlot()
	if slot == nil {
		return NewNullValue(), nil
	}
	return slot.Value, nil
}

func (a *ArrayValueShift) GetName() string {
	return "shift"
}

func (a *ArrayValueShift) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueShift) GetIsStatic() bool {
	return false
}

func (a *ArrayValueShift) GetParams() []GetValue {
	return []GetValue{}
}

func (a *ArrayValueShift) GetVariables() []Variable {
	return []Variable{}
}

func (a *ArrayValueShift) GetReturnType() Types {
	return nil
}
