package data

type String struct {
}

func (i String) Is(value Value) bool {
	switch value.(type) {
	case *StringValue:
		return true
	case *IntValue, *FloatValue, *BoolValue:
		// PHP 弱类型模式下，int/float/bool 可隐式转换为 string 参数/返回值
		return true
	}
	return false
}

func (i String) String() string {
	return "string"
}
