package data

import "sort"

type ArrayValueSort struct {
	source *[]Value
}

// Call 实现数组的 sort 方法
// 对数组元素进行排序，默认按字符串比较排序，并返回排序后的数组
func (a *ArrayValueSort) Call(ctx Context) (GetValue, Control) {
	// 创建数组的副本进行排序
	sortedArray := make([]Value, len(*a.source))
	copy(sortedArray, *a.source)

	// 使用字符串比较进行排序
	sort.Slice(sortedArray, func(i, j int) bool {
		return sortedArray[i].AsString() < sortedArray[j].AsString()
	})

	// 更新原数组
	*a.source = sortedArray

	// 返回排序后的数组
	return NewArrayValue(sortedArray), nil
}

func (a *ArrayValueSort) GetName() string {
	return "sort"
}

func (a *ArrayValueSort) GetModifier() Modifier {
	return ModifierPublic
}

func (a *ArrayValueSort) GetIsStatic() bool {
	return false
}

func (a *ArrayValueSort) GetParams() []GetValue {
	return []GetValue{}
}

func (a *ArrayValueSort) GetVariables() []Variable {
	return []Variable{}
}

func (a *ArrayValueSort) GetReturnType() Types {
	return Arrays{}
}
