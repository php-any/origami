package data

// InterfaceStmt 表示接口定义
type InterfaceStmt interface {
	GetFrom() From
	GetValue(ctx Context) (GetValue, Control)
	GetName() string
	GetExtend() *string                   // 父接口名
	GetMethod(name string) (Method, bool) // 方法列表
	GetMethods() []Method                 // 方法列表
}

// GetReturnType 表示可以获取返回类型的接口
type GetReturnType interface {
	GetReturnType() Types // 返回类型
}
