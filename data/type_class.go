package data

type Class struct{ Name string }

func (i Class) Is(value Value) bool {
	switch c := value.(type) {
	case *ClassValue:
		if i.Name == c.Class.GetName() {
			return true
		}
		if i.Name == "iterable" {
			return NominalIsA(c.Class, "Traversable", c.GetVM())
		}
		return NominalIsA(c.Class, i.Name, c.GetVM())
	case *ThisValue:
		if i.Name == c.Class.GetName() {
			return true
		}
		if i.Name == "iterable" {
			return NominalIsA(c.Class, "Traversable", c.GetVM())
		}
		return NominalIsA(c.Class, i.Name, c.GetVM())
	case *ArrayValue:
		return TypeNameEqual(i.Name, "iterable")
	case *ThrowValue:
		if c.Object != nil {
			return NominalIsA(c.Object.Class, i.Name, c.Object.GetVM())
		}
		name := c.Name
		if name == "" {
			name = "Exception"
		}
		return InternalTypeIsA(name, i.Name)
	case *FuncValue, *BoundFuncValue:
		return TypeNameEqual(i.Name, "Closure")
	case Generator:
		return TypeNameEqual(i.Name, "iterable") || InternalTypeIsA("Generator", i.Name)
	}
	return false
}

func (i Class) String() string { return i.Name }
