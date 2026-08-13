package data

type Arrays struct {
}

func (i Arrays) Is(value Value) bool {
	switch value.(type) {
	case *ArrayValue:
		return true
	case *ObjectValue:
		// PHP 关联数组在 Origami 中可能用 ObjectValue 表示，仍视为 array 类型
		return true
	}
	return false
}

func (i Arrays) String() string {
	return "array"
}
