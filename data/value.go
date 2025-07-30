package data

// Value 表示值接口
type Value interface {
	GetValue
	AsString() string
}

// CallableValue 表示可调用值接口
type CallableValue interface {
	Value
	// Call 调用函数
	Call(args ...Value) (Value, Control)
	// IsMethod 是否为方法
	IsMethod() bool
	// GetMethodName 获取方法名
	GetMethodName() string
}

type GetProperty interface {
	GetProperty(name string) (Value, bool)
}

type GetMethod interface {
	GetMethod(name string) (Method, bool)
}
