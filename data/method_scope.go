package data

// MethodDeclaringClass preserves the lexical scope of inherited methods.
// It is shared by explicit calls and implicit language conversions.
func MethodDeclaringClass(vm VM, class ClassStmt, methodName string) ClassStmt {
	if class == nil {
		return nil
	}
	if methodName == "__construct" && vm != nil {
		for class.GetExtend() != nil {
			parent, ctl := vm.GetOrLoadClass(*class.GetExtend())
			if ctl != nil || parent == nil {
				break
			}
			method, _ := class.GetMethod(methodName)
			inherited, _ := parent.GetMethod(methodName)
			if method == nil || method != inherited {
				break
			}
			class = parent
		}
		return class
	}
	if declaresMethod(class, methodName) || vm == nil {
		return class
	}
	last := class
	for last.GetExtend() != nil {
		parentName := last.GetExtend()
		if *parentName == "" {
			break
		}
		parent, ctl := vm.GetOrLoadClass(*parentName)
		if ctl != nil || parent == nil {
			break
		}
		if declaresMethod(parent, methodName) {
			return parent
		}
		last = parent
	}
	return class
}

func declaresMethod(class ClassStmt, name string) bool {
	if method, found := class.GetMethod(name); found && method != nil {
		return true
	}
	if static, ok := class.(GetStaticMethod); ok {
		if method, found := static.GetStaticMethod(name); found && method != nil {
			return true
		}
	}
	return false
}
