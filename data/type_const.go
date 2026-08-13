package data

type Const struct {
	// 保留常量类型
	MyType Types
}

// Is 永远不允许赋值
func (i Const) Is(value Value) bool {
	return false
}

func (i Const) String() string {
	return "const"
}
