package data

type Class struct {
	Name string
}

func (i Class) Is(value Value) bool {
	if c, ok := value.(*ClassValue); ok {
		if i.Name == c.Class.GetName() {
			return true
		}
	}
	return false
}

func (i Class) String() string {
	return i.Name
}
