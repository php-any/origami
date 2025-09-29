package data

type ArrayValueFlat struct {
	source []Value
}

// Call 实现数组的 flat 方法
// 将嵌套数组扁平化，返回一个新数组，其中所有子数组元素都被递归地连接到指定深度
func (a *ArrayValueFlat) Call(ctx Context) (GetValue, Control) {
	// 获取深度参数，默认为 1
	depth := 1
	if depthArg, ok := ctx.GetIndexValue(0); ok {
		if depthInt, ok := depthArg.(AsInt); ok {
			if d, err := depthInt.AsInt(); err == nil {
				depth = d
			}
		}
	}

	// 递归扁平化数组
	result := a.flattenArray(a.source, depth)
	return NewArrayValue(result), nil
}

// flattenArray 递归扁平化数组的辅助函数
func (a *ArrayValueFlat) flattenArray(arr []Value, depth int) []Value {
	if depth <= 0 {
		return arr
	}

	var result []Value
	for _, element := range arr {
		// 如果元素是数组且深度大于0，则递归扁平化
		if arrayElement, ok := element.(*ArrayValue); ok && depth > 0 {
			flattened := a.flattenArray(arrayElement.Value, depth-1)
			result = append(result, flattened...)
		} else {
			result = append(result, element)
		}
	}

	return result
}

func (a *ArrayValueFlat) GetName() string {
	return "flat"
}

func (a *ArrayValueFlat) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueFlat) GetIsStatic() bool {
	return false
}

func (a *ArrayValueFlat) GetParams() []GetValue {
	return []GetValue{
		NewParameter("depth", 0),
	}
}

func (a *ArrayValueFlat) GetVariables() []Variable {
	return []Variable{
		NewVariable("depth", 0, nil),
	}
}

func (a *ArrayValueFlat) GetReturnType() Types {
	return Arrays{}
}
