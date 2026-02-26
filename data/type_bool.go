package data

type Bool struct {
}

func (i Bool) Is(value Value) bool {
	// 原生 PHP 中，bool 类型声明会对标量做弱类型转换：
	// int/string 等在传给 bool 参数或作为 bool 返回值时都是允许的。
	// 这里遵循同样的语义：只要实现了 AsBool 接口，就认为可以作为 bool 使用。
	if _, ok := value.(*BoolValue); ok {
		return true
	}
	if _, ok := value.(AsBool); ok {
		return true
	}
	return false
}

func (i Bool) String() string {
	return "bool"
}
