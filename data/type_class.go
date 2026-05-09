package data

import "strings"

type Class struct {
	Name string
}

func (i Class) Is(value Value) bool {
	switch c := value.(type) {
	case *ClassValue:
		return isClassValueInstanceOf(i.Name, c.Class, c.GetVM())
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
	case *ThrowValue:
		// ThrowValue 在 PHP 层面也是可抛出的，应能被 catch 捕获。
		// 有 Object 时直接复用 ClassValue 的类型检查逻辑；
		// 无 Object 时（Go 内部错误）也视为可抛出的异常。
		if c.Object != nil {
			return isClassValueInstanceOf(i.Name, c.Object.Class, c.Object.GetVM())
		}
		// Go 内部错误：允许被 catch (Throwable / Exception / Error) 捕获
		baseName := i.Name
		if idx := strings.LastIndex(i.Name, "\\"); idx >= 0 {
			baseName = i.Name[idx+1:]
		}
		if baseName == "Throwable" || baseName == "Exception" || baseName == "Error" {
			return true
		}
	}

	return false
}

// isClassValueInstanceOf 检查一个 ClassStmt 是否实现了目标类型（类名或接口名）
func isClassValueInstanceOf(target string, class ClassStmt, vm VM) bool {
	if target == class.GetName() {
		return true
	}
	for _, s := range class.GetImplements() {
		if target == s {
			return true
		} else if interfaceExtends(vm, s, target) {
			return true
		}
	}
	return extendISClass(target, class.GetExtend(), vm)
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
			vm.ThrowControl(acl)
		}
		iface, ok = vm.GetInterface(ifaceName)
		if !ok {
			return false
		}
	}

	if iface.GetName() == target {
		return true
	}

	visited := make(map[string]bool)
	var queue []string
	queue = append(queue, iface.GetExtends()...)

	for len(queue) > 0 {
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
