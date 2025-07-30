package data

type String struct {
}

func (i String) Is(value Value) bool {
	switch value.(type) {
	case *StringValue:
		return true
	}
	return false
}

func (i String) String() string {
	return "string"
}
