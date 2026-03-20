package data

// ZVal 模仿 PHP 的 zval 结构
type ZVal struct {
	Name  string // 变量名称，用于 extract 等按名称操作
	Value Value
}

// NewZVal 创建一个新的 ZVal
func NewZVal(v Value) *ZVal {
	return &ZVal{
		Value: v,
	}
}

// NewNamedZVal 创建一个带名称的 ZVal
func NewNamedZVal(name string, v Value) *ZVal {
	return &ZVal{
		Name:  name,
		Value: v,
	}
}

type ZValGetter interface {
	GetZVal(v Variable) *ZVal
}
