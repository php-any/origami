package data

// InterfaceStmt 表示接口定义
type InterfaceStmt interface {
	GetFrom() From
	GetValue(ctx Context) (GetValue, Control)
	GetName() string
	// GetExtend() *string // 父接口名
	GetExtends() []string                 // 父接口名
	GetMethod(name string) (Method, bool) // 方法列表
	GetMethods() []Method                 // 方法列表
}

// GetReturnType 表示可以获取返回类型的接口
type GetReturnType interface {
	GetReturnType() Types // 返回类型
}

type ClassGeneric interface {
	ClassStmt

	Clone(map[string]Types) ClassGeneric
	GenericList() []Types
}

// Iterator 迭代器值约束
type Iterator interface {
	// Current 返回当前元素
	Current(ctx Context) (Value, Control)
	// Key 返回当前元素的键
	Key(ctx Context) (Value, Control)
	// Next 将指针向前移动到下一个元素
	Next(ctx Context) Control
	// Rewind 将指针重置到第一个元素
	Rewind(ctx Context) (Value, Control)
	// Valid 检查当前位置是否有效（是否还有元素）
	Valid(ctx Context) (Value, Control)
}

// Generator 生成器约束
type Generator interface {
	Value
	Iterator
	// Send 向暂停中的生成器传递一个值
	Send(ctx Context, value Value) Control
	// Throw 向暂停中的生成器抛出一个异常
	Throw(ctx Context) Control
	// GetReturn 获取生成器函数的 return 值。
	GetReturn(ctx Context) (Value, Control)
}

type GetName interface {
	GetName() string
}

type GetZVal interface {
	GetZVal(ctx Context) (*ZVal, Control)
}
