package data

type ArrayValuePush struct {
	source *[]Value
}

// Call 实现数组的 push 方法
// 将一个或多个元素添加到数组的末尾，并返回新的数组长度
func (a *ArrayValuePush) Call(ctx Context) (GetValue, Control) {
	// 获取所有参数
	for _, argument := range a.GetParams() {
		var args []Value
		// argument data.Parameters
		argv, _ := argument.GetValue(ctx)
		if ar, ok := argv.(*ArrayValue); ok {
			for _, v := range ar.Value {
				args = append(args, v)
			}
		}
		// 将参数添加到数组末尾
		*a.source = append(*a.source, args...)
	}

	// 返回新的数组长度
	return NewIntValue(len(*a.source)), nil
}

func (a *ArrayValuePush) GetName() string {
	return "push"
}

func (a *ArrayValuePush) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValuePush) GetIsStatic() bool {
	return false
}

func (a *ArrayValuePush) GetParams() []GetValue {
	return []GetValue{
		NewParameters("items", 0),
	}
}

func (a *ArrayValuePush) GetVariables() []Variable {
	return []Variable{
		NewVariable("items", 0, nil),
	}
}

func (a *ArrayValuePush) GetReturnType() Types {
	return nil
}
