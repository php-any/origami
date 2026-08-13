package data

type Mixed struct{}

func (m Mixed) Is(value Value) bool {
	return true
}

func (m Mixed) String() string {
	return "mixed"
}
