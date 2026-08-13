package node

import "github.com/php-any/origami/data"

// LookupStaticProperty 在类、父类及 implements 的接口继承链上查找类常量/静态属性。
func LookupStaticProperty(vm data.VM, class data.ClassStmt, name string) (data.Value, bool) {
	if class == nil {
		return nil, false
	}
	if v, ok := staticPropertyOn(class, name); ok {
		return v, true
	}
	for _, ifaceName := range class.GetImplements() {
		if v, ok := lookupStaticPropertyOnInterface(vm, ifaceName, name); ok {
			return v, true
		}
	}
	last := class
	for last.GetExtend() != nil {
		parentName := last.GetExtend()
		if parentName == nil || *parentName == "" {
			break
		}
		parent, acl := vm.GetOrLoadClass(*parentName)
		if acl != nil || parent == nil {
			break
		}
		if v, ok := staticPropertyOn(parent, name); ok {
			return v, true
		}
		for _, ifaceName := range parent.GetImplements() {
			if v, ok := lookupStaticPropertyOnInterface(vm, ifaceName, name); ok {
				return v, true
			}
		}
		last = parent
	}
	return nil, false
}

func staticPropertyOn(class data.ClassStmt, name string) (data.Value, bool) {
	if gsp, ok := class.(data.GetStaticProperty); ok {
		return gsp.GetStaticProperty(name)
	}
	return nil, false
}

func lookupStaticPropertyOnInterface(vm data.VM, ifaceName, name string) (data.Value, bool) {
	iface, acl := vm.LoadPkg(ifaceName)
	if acl != nil || iface == nil {
		return nil, false
	}
	if gsp, ok := iface.(data.GetStaticProperty); ok {
		if v, ok := gsp.GetStaticProperty(name); ok {
			return v, true
		}
	}
	if is, ok := iface.(data.InterfaceStmt); ok {
		for _, parent := range is.GetExtends() {
			if v, ok := lookupStaticPropertyOnInterface(vm, parent, name); ok {
				return v, true
			}
		}
	}
	return nil, false
}
