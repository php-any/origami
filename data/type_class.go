package data

type Class struct {
	Name string
}

func (i Class) Is(value Value) bool {
	switch c := value.(type) {
	case *ClassValue:
		if i.Name == c.Class.GetName() {
			return true
		}
		for _, s := range c.Class.GetImplements() {
			// 直接实现的接口
			if i.Name == s {
				return true
			}
			// 接口继承：类实现的接口可能继承了目标接口
			if interfaceExtends(c.GetVM(), s, i.Name) {
				return true
			}
		}
		return extendISClass(i.Name, c.Class.GetExtend(), c.GetVM())
	case *ThisValue:
		if i.Name == c.Class.GetName() {
			return true
		}
		for _, s := range c.Class.GetImplements() {
			if i.Name == s {
				return true
			}
			if interfaceExtends(c.GetVM(), s, i.Name) {
				return true
			}
		}
		return extendISClass(i.Name, c.Class.GetExtend(), c.GetVM())
	case *ArrayValue:
		if i.Name == "iterable" {
			return true
		}
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
				if check == s || interfaceExtends(vm, s, check) {
					return true
				}
			}
			extend = c.GetExtend()
		}
	}
	return false
}

// interfaceExtends 检查接口 ifaceName 是否直接或通过继承链实现了目标接口 target。
// 用于类型系统中判断：类实现的接口本身也继承了其他接口的情况。
func interfaceExtends(vm VM, ifaceName, target string) bool {
	if ifaceName == target {
		return true
	}

	if vm == nil {
		return false
	}

	iface, ok := vm.GetInterface(ifaceName)
	if !ok {
		return false
	}

	// 直接同名
	if iface.GetName() == target {
		return true
	}

	// 沿着接口的 extends 链向上查找
	ext := iface.GetExtend()
	for ext != nil {
		if *ext == target {
			return true
		}
		parent, ok := vm.GetInterface(*ext)
		if !ok {
			break
		}
		if parent.GetName() == target {
			return true
		}
		ext = parent.GetExtend()
	}

	return false
}
