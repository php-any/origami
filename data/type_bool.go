package data

type Bool struct {
}

func (i Bool) Is(value Value) bool {
	if _, ok := value.(AsBool); ok {
		return true
	}
	return false
}

func (i Bool) String() string {
	return "bool"
}
