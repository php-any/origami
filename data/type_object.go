package data

type Object struct {
}

func (i Object) Is(value Value) bool {
	switch value.(type) {
	case *ObjectValue:
		return true
	case *ClassValue:
		return true
	}
	return false
}

func (i Object) String() string {
	return "object"
}
