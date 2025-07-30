package data

// Generic 泛型
type Generic struct {
	Name  string
	Types []Types
}

func (i Generic) Is(value Value) bool {
	switch value.(type) {
	case *ArrayValue:
		return true
	}
	return false
}

func (i Generic) String() string {
	return "array"
}
