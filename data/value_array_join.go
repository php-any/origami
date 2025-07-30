package data

type ArrayValueJoin struct {
	source []Value
}

// Call 实现数组的 join 方法
// 将数组的所有元素转换为字符串并用指定的分隔符连接
func (a *ArrayValueJoin) Call(ctx Context) (GetValue, Control) {
	// 获取分隔符参数，默认为逗号
	separator := ","
	if sepArg, ok := ctx.GetIndexValue(0); ok {
		if sepStr, ok := sepArg.(AsString); ok {
			separator = sepStr.AsString()
		}
	}

	// 将数组元素转换为字符串并用分隔符连接
	var result string
	for i, value := range a.source {
		if i > 0 {
			result += separator
		}
		result += value.AsString()
	}

	return NewStringValue(result), nil
}

func (a *ArrayValueJoin) GetName() string {
	return "join"
}

func (a *ArrayValueJoin) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueJoin) GetIsStatic() bool {
	return false
}

func (a *ArrayValueJoin) GetParams() []GetValue {
	return []GetValue{
		NewParameter("separator", 0),
	}
}

func (a *ArrayValueJoin) GetVariables() []Variable {
	return []Variable{
		NewVariable("separator", 0, nil),
	}
}
