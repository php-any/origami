package data

type Float struct {
}

func (i Float) Is(value Value) bool {
	if _, ok := value.(AsFloat); ok {
		return true
	}
	return false
}

func (i Float) String() string {
	return "float"
}
