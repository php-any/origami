package data

type Class struct {
	Name string
}

func (i Class) Is(value Value) bool {
	if c, ok := value.(*ClassValue); ok {
		if i.Name == c.Class.GetName() {
			return true
		}
		for _, s := range c.Class.GetImplements() {
			if i.Name == s {
				return true
			}
		}
		return extendISClass(i.Name, c.Class.GetExtend(), c.GetVM())
	} else if c, ok := value.(*ThisValue); ok {
		if i.Name == c.Class.GetName() {
			return true
		}
		for _, s := range c.Class.GetImplements() {
			if i.Name == s {
				return true
			}
		}
		return extendISClass(i.Name, c.Class.GetExtend(), c.GetVM())
	}
	return false
}

func (i Class) String() string {
	return i.Name
}

func extendISClass(check string, extend *string, vm VM) bool {
	for extend != nil {
		c, ok := vm.GetClass(*extend)
		extend = nil
		if ok {
			if check == c.GetName() {
				return true
			}
			for _, s := range c.GetImplements() {
				if check == s {
					return true
				}
			}
			extend = c.GetExtend()
		}
	}
	return false
}
