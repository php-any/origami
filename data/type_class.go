package data

type Class struct{ Name string }

func (i Class) Is(value Value) bool {
	return NominalValueMatches(value, i.Name, nil)
}

// NominalValueMatches is the runtime predicate shared by type hints, catch,
// instanceof and builtin reflection. It never autoloads.
func NominalValueMatches(value Value, name string, ctx Context) bool {
	if TypeNameEqual(name, "iterable") {
		if _, ok := value.(*ArrayValue); ok {
			return true
		}
		name = "Traversable"
	}
	var vm VM
	if ctx != nil {
		vm = ctx.GetVM()
	}
	switch c := value.(type) {
	case *ClassValue:
		if vm == nil {
			vm = c.GetVM()
		}
		return NominalIsA(c.Class, name, vm)
	case *ThisValue:
		if vm == nil {
			vm = c.GetVM()
		}
		return NominalIsA(c.Class, name, vm)
	case *ThrowValue:
		if c.Object != nil {
			if vm == nil {
				vm = c.Object.GetVM()
			}
			return NominalIsA(c.Object.Class, name, vm)
		}
		actual := c.Name
		if actual == "" {
			actual = "Exception"
		}
		if vm != nil {
			if class, found := vm.GetClass(actual); found {
				return NominalIsA(class, name, vm)
			}
		}
		return InternalTypeIsA(actual, name)
	case *FuncValue, *BoundFuncValue:
		return TypeNameEqual(name, "Closure")
	case Generator:
		return InternalTypeIsA("Generator", name)
	}
	return false
}

func (i Class) String() string { return i.Name }
