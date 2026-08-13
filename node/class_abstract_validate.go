package node

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

// ValidateConcreteClassAbstractMethods 检查非抽象类是否仍含未实现的抽象/接口方法（与 PHP zend 一致）
func ValidateConcreteClassAbstractMethods(vm data.VM, class data.ClassStmt) data.Control {
	selfAbstract := abstractMethodsDeclaredOnClass(class)
	if len(selfAbstract) > 0 {
		msg := formatDeclaresAbstractMethodFatal(class.GetName(), selfAbstract)
		return data.NewCompileFatal(class.GetFrom(), msg)
	}

	missing, acl := collectUnimplementedAbstractMethods(vm, class)
	if acl != nil {
		return acl
	}
	if len(missing) == 0 {
		return nil
	}
	msg := formatAbstractMethodsFatal(class.GetName(), missing)
	return data.NewCompileFatal(class.GetFrom(), msg)
}

func abstractMethodsDeclaredOnClass(class data.ClassStmt) []string {
	var names []string
	for _, method := range class.GetMethods() {
		if _, ok := method.(*AbstractMethod); ok {
			names = append(names, method.GetName())
		}
	}
	if cs, ok := class.(*ClassStatement); ok {
		for _, method := range cs.StaticMethods {
			if _, ok := method.(*AbstractMethod); ok {
				names = append(names, method.GetName())
			}
		}
	}
	if cg, ok := class.(*ClassGeneric); ok {
		for _, method := range cg.StaticMethods {
			if _, ok := method.(*AbstractMethod); ok {
				names = append(names, method.GetName())
			}
		}
	}
	return names
}

func formatDeclaresAbstractMethodFatal(className string, methods []string) string {
	if len(methods) == 1 {
		return fmt.Sprintf(
			"Class %s declares abstract method %s() and must therefore be declared abstract",
			className, methods[0],
		)
	}
	names := make([]string, len(methods))
	for i, m := range methods {
		names[i] = m + "()"
	}
	return fmt.Sprintf(
		"Class %s declares abstract methods %s and must therefore be declared abstract",
		className, strings.Join(names, ", "),
	)
}

func formatAbstractMethodsFatal(className string, missing []string) string {
	n := len(missing)
	methodWord := "method"
	remainingWord := "method"
	if n != 1 {
		methodWord = "methods"
		remainingWord = "methods"
	}
	return fmt.Sprintf(
		"Class %s contains %d abstract %s and must therefore be declared abstract or implement the remaining %s (%s)",
		className, n, methodWord, remainingWord, strings.Join(missing, ", "),
	)
}

func collectUnimplementedAbstractMethods(vm data.VM, class data.ClassStmt) ([]string, data.Control) {
	seen := make(map[string]struct{})
	var missing []string
	add := func(entry string) {
		if _, ok := seen[entry]; ok {
			return
		}
		seen[entry] = struct{}{}
		missing = append(missing, entry)
	}

	for _, ifaceName := range class.GetImplements() {
		entries, acl := unimplementedInterfaceMethods(vm, class, ifaceName)
		if acl != nil {
			return nil, acl
		}
		for _, e := range entries {
			add(e)
		}
	}

	if class.GetExtend() != nil {
		parent, acl := vm.GetOrLoadClass(*class.GetExtend())
		if acl != nil {
			return nil, acl
		}
		entries, acl := unimplementedFromParentClass(vm, class, parent)
		if acl != nil {
			return nil, acl
		}
		for _, e := range entries {
			add(e)
		}
	}

	return missing, nil
}

func unimplementedFromParentClass(vm data.VM, class, parent data.ClassStmt) ([]string, data.Control) {
	var missing []string
	for _, method := range parent.GetMethods() {
		if _, ok := method.(*AbstractMethod); ok {
			if !classImplementsConcreteMethod(vm, class, method.GetName()) {
				missing = append(missing, parent.GetName()+"::"+method.GetName())
			}
		}
	}
	for _, ifaceName := range parent.GetImplements() {
		entries, acl := unimplementedInterfaceMethods(vm, class, ifaceName)
		if acl != nil {
			return nil, acl
		}
		missing = append(missing, entries...)
	}
	if parent.GetExtend() != nil {
		grand, acl := vm.GetOrLoadClass(*parent.GetExtend())
		if acl != nil {
			return nil, acl
		}
		entries, acl := unimplementedFromParentClass(vm, class, grand)
		if acl != nil {
			return nil, acl
		}
		missing = append(missing, entries...)
	}
	return missing, nil
}

func unimplementedInterfaceMethods(vm data.VM, class data.ClassStmt, ifaceName string) ([]string, data.Control) {
	iface, ok := vm.GetInterface(ifaceName)
	if !ok {
		var acl data.Control
		iface, acl = vm.GetOrLoadInterface(ifaceName)
		if acl != nil {
			return nil, acl
		}
	}
	var missing []string
	for _, method := range iface.GetMethods() {
		if !classImplementsConcreteMethod(vm, class, method.GetName()) {
			missing = append(missing, ifaceName+"::"+method.GetName())
		}
	}
	for _, parentIface := range iface.GetExtends() {
		entries, acl := unimplementedInterfaceMethods(vm, class, parentIface)
		if acl != nil {
			return nil, acl
		}
		missing = append(missing, entries...)
	}
	return missing, nil
}

func classImplementsConcreteMethod(vm data.VM, class data.ClassStmt, methodName string) bool {
	for c := class; c != nil; {
		if classDeclaresConcreteMethod(c, methodName) {
			return true
		}
		if c.GetExtend() == nil {
			break
		}
		parent, acl := vm.GetOrLoadClass(*c.GetExtend())
		if acl != nil || parent == nil {
			break
		}
		c = parent
	}
	return false
}

func classDeclaresConcreteMethod(c data.ClassStmt, methodName string) bool {
	if m, ok := c.GetMethod(methodName); ok && m != nil {
		if _, isAbstract := m.(*AbstractMethod); !isAbstract {
			return true
		}
	}
	if sg, ok := c.(interface {
		GetStaticMethod(string) (data.Method, bool)
	}); ok {
		if m, ok := sg.GetStaticMethod(methodName); ok && m != nil {
			if _, isAbstract := m.(*AbstractMethod); !isAbstract {
				return true
			}
		}
	}
	return false
}
