package data

type ArrayValueConcat struct {
	source []*ZVal
}

// Call 实现数组的 concat 方法
// 合并两个或多个数组，返回一个新数组，包含所有数组的元素
func (a *ArrayValueConcat) Call(ctx Context) (GetValue, Control) {
	// 创建结果数组，先复制原数组
	result := make([]*ZVal, len(a.source))
	copy(result, a.source)

	// 获取所有参数并添加到结果数组
	for _, argument := range a.GetParams() {
		argv, _ := argument.GetValue(ctx)
		if ar, ok := argv.(*ArrayValue); ok {
			result = append(result, ar.List...)
		} else {
			result = append(result, NewZVal(argv.(Value)))
		}
	}

	// 转换为 []Value 用于 NewArrayValue
	values := make([]Value, len(result))
	for i, zval := range result {
		values[i] = zval.Value
	}
	return NewArrayValue(values), nil
}

func (a *ArrayValueConcat) GetName() string {
	return "concat"
}

func (a *ArrayValueConcat) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueConcat) GetIsStatic() bool {
	return false
}

func (a *ArrayValueConcat) GetParams() []GetValue {
	return []GetValue{
		NewParameters("items", 0),
	}
}

func (a *ArrayValueConcat) GetVariables() []Variable {
	return []Variable{
		NewVariable("items", 0, nil),
	}
}

func (a *ArrayValueConcat) GetReturnType() Types {
	return Arrays{}
}
