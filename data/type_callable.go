package data

type Callable struct {
}

func (c Callable) Is(value Value) bool {
	switch value.(type) {
	case *FuncValue:
		return true
	}
	return false
}

func (c Callable) String() string {
	return "callable"
}
