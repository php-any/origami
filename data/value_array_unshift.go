package data

type ArrayValueUnshift struct {
	source *[]Value
}

// Call 实现数组的 unshift 方法
// 将一个或多个元素添加到数组的开头，并返回新的数组长度
func (a *ArrayValueUnshift) Call(ctx Context) (GetValue, Control) {
	// 获取所有参数
	var args []Value
	for _, argument := range a.GetParams() {
		// argument data.Parameters
		argv, _ := argument.GetValue(ctx)
		if ar, ok := argv.(*ArrayValue); ok {
			for _, v := range ar.Value {
				args = append(args, v)
			}
		}
	}

	// 将参数添加到数组开头
	*a.source = append(args, *a.source...)

	// 返回新的数组长度
	return NewIntValue(len(*a.source)), nil
}

func (a *ArrayValueUnshift) GetName() string {
	return "unshift"
}

func (a *ArrayValueUnshift) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueUnshift) GetIsStatic() bool {
	return false
}

func (a *ArrayValueUnshift) GetParams() []GetValue {
	return []GetValue{
		NewParameters("items", 0),
	}
}

func (a *ArrayValueUnshift) GetVariables() []Variable {
	return []Variable{
		NewVariable("items", 0, nil),
	}
}

func (a *ArrayValueUnshift) GetReturnType() Types {
	return nil
}
