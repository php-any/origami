package data

// ZVal 模仿 PHP 的 zval 结构
type ZVal struct {
	Value Value
}

// NewZVal 创建一个新的 ZVal
func NewZVal(v Value) *ZVal {
	return &ZVal{
		Value: v,
	}
}

type ZValGetter interface {
	GetZVal(v Variable) *ZVal
}
