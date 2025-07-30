package data

// GetValue 表示可以获取值的节点
type GetValue interface {
	// GetValue 获取节点的值
	GetValue(ctx Context) (GetValue, Control)
}
