package data

type Int struct {
}

func (i Int) Is(value Value) bool {
	switch value.(type) {
	case *IntValue:
		return true
	}
	return false
}

func (i Int) String() string {
	return "int"
}
