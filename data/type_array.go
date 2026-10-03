package data

type Arrays struct {
}

func (i Arrays) Is(value Value) bool {
	switch value.(type) {
	case *ArrayValue:
		return true
	}
	return false
}

func (i Arrays) String() string {
	return "array"
}
