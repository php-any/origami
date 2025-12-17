package data

type Bool struct {
}

func (i Bool) Is(value Value) bool {
	if _, ok := value.(*BoolValue); ok {
		return true
	}
	return false
}

func (i Bool) String() string {
	return "bool"
}
