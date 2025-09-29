package data

type ArrayValueSlice struct {
	source []Value
}

// Call 实现数组的 slice 方法
// 返回数组的一个浅拷贝，从 start 到 end（不包括 end）的元素组成的新数组
func (a *ArrayValueSlice) Call(ctx Context) (GetValue, Control) {
	// 获取参数
	start := 0
	end := len(a.source)

	// 获取 start 参数
	if startArg, ok := ctx.GetIndexValue(0); ok {
		if startInt, ok := startArg.(AsInt); ok {
			if s, err := startInt.AsInt(); err == nil {
				start = s
			}
		}
	}

	// 获取 end 参数
	if endArg, ok := ctx.GetIndexValue(1); ok {
		if endInt, ok := endArg.(AsInt); ok {
			if e, err := endInt.AsInt(); err == nil {
				end = e
			}
		}
	}

	// 处理负数索引
	if start < 0 {
		start = len(a.source) + start
	}
	if end < 0 {
		end = len(a.source) + end
	}

	// 边界检查
	if start < 0 {
		start = 0
	}
	if end > len(a.source) {
		end = len(a.source)
	}
	if start > end {
		start = end
	}

	// 返回切片
	return NewArrayValue(a.source[start:end]), nil
}

func (a *ArrayValueSlice) GetName() string {
	return "slice"
}

func (a *ArrayValueSlice) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueSlice) GetIsStatic() bool {
	return false
}

func (a *ArrayValueSlice) GetParams() []GetValue {
	return []GetValue{
		NewParameter("start", 0),
		NewParameter("end", 1),
	}
}

func (a *ArrayValueSlice) GetVariables() []Variable {
	return []Variable{
		NewVariable("start", 0, nil),
		NewVariable("end", 1, nil),
	}
}

func (a *ArrayValueSlice) GetReturnType() Types {
	return Arrays{}
}
