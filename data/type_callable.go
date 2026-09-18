package data

type Callable struct {
}

func (c Callable) Is(value Value) bool {
	switch v := value.(type) {
	case *FuncValue, *ArrayValue:
		return true
	case *StringValue:
		return true
	case *ClassValue:
		// PHP：实现 __invoke 的对象是 callable
		if _, ok := v.GetMethod("__invoke"); ok {
			return true
		}
		if v.Class != nil {
			if _, ok := v.Class.GetMethod("__invoke"); ok {
				return true
			}
		}
	}
	return false
}

func (c Callable) String() string {
	return "callable"
}
