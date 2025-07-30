package data

type ArrayValueSplice struct {
	source *[]Value
}

// Call 实现数组的 splice 方法
// 通过删除现有元素和/或添加新元素来更改数组的内容，返回被删除的元素数组
func (a *ArrayValueSplice) Call(ctx Context) (GetValue, Control) {
	// 获取参数
	start := 0
	deleteCount := len(*a.source)

	// 获取 start 参数
	if startArg, ok := ctx.GetIndexValue(0); ok {
		if startInt, ok := startArg.(AsInt); ok {
			if s, err := startInt.AsInt(); err == nil {
				start = s
			}
		}
	}

	// 获取 deleteCount 参数
	if deleteCountArg, ok := ctx.GetIndexValue(1); ok {
		if deleteCountInt, ok := deleteCountArg.(AsInt); ok {
			if d, err := deleteCountInt.AsInt(); err == nil {
				deleteCount = d
			}
		}
	}

	// 处理负数索引
	if start < 0 {
		start = len(*a.source) + start
	}

	// 边界检查
	if start < 0 {
		start = 0
	}
	if start > len(*a.source) {
		start = len(*a.source)
	}
	if deleteCount < 0 {
		deleteCount = 0
	}
	if start+deleteCount > len(*a.source) {
		deleteCount = len(*a.source) - start
	}

	// 获取要删除的元素
	deletedElements := make([]Value, deleteCount)
	copy(deletedElements, (*a.source)[start:start+deleteCount])

	// 获取要插入的元素
	var insertElements []Value
	for i := 2; ; i++ {
		arg, ok := ctx.GetIndexValue(i)
		if !ok {
			break
		}
		insertElements = append(insertElements, arg)
	}

	// 执行 splice 操作
	newArray := make([]Value, 0, len(*a.source)-deleteCount+len(insertElements))
	newArray = append(newArray, (*a.source)[:start]...)
	newArray = append(newArray, insertElements...)
	newArray = append(newArray, (*a.source)[start+deleteCount:]...)
	*a.source = newArray

	// 返回被删除的元素
	return NewArrayValue(deletedElements), nil
}

func (a *ArrayValueSplice) GetName() string {
	return "splice"
}

func (a *ArrayValueSplice) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueSplice) GetIsStatic() bool {
	return false
}

func (a *ArrayValueSplice) GetParams() []GetValue {
	return []GetValue{
		NewParameter("start", 0),
		NewParameter("deleteCount", 1),
		NewParameters("items", 2),
	}
}

func (a *ArrayValueSplice) GetVariables() []Variable {
	return []Variable{
		NewVariable("start", 0, nil),
		NewVariable("deleteCount", 1, nil),
		NewVariable("items", 2, nil),
	}
}
