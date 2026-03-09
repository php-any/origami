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
		// 先直接比较 implements 列表中的字符串（不依赖接口是否已加载）
		for _, s := range c.Class.GetImplements() {
			if i.Name == s {
				return true
			} else if interfaceExtends(c.GetVM(), s, i.Name) {
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
		}
		for _, s := range c.Class.GetImplements() {
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
// 这里只基于 VM 中已注册的接口做判断，不负责触发自动加载；接口加载应在类/接口文件加载阶段完成。
func interfaceExtends(vm VM, ifaceName, target string) bool {
	if ifaceName == target {
		return true
	}

	if vm == nil {
		return false
	}

	iface, ok := vm.GetInterface(ifaceName)
	if !ok {
		_, acl := vm.LoadPkg(ifaceName)
		if acl != nil {
			vm.ThrowControl(acl) // TODO
		}
		iface, ok = vm.GetInterface(ifaceName)
		if !ok {
			return false
		}
	}

	// 直接同名
	if iface.GetName() == target {
		return true
	}

	// 沿着接口的 extends 链向上查找（支持多个父接口）
	visited := make(map[string]bool)
	var queue []string

	queue = append(queue, iface.GetExtends()...)

	for len(queue) > 0 {
		// 取出队首元素
		name := queue[0]
		queue = queue[1:]

		if name == target {
			return true
		}
		if visited[name] {
			continue
		}
		visited[name] = true

		parent, ok := vm.GetInterface(name)
		if !ok {
			continue
		}
		if parent.GetName() == target {
			return true
		}

		queue = append(queue, parent.GetExtends()...)
	}

	return false
}
