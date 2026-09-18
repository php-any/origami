package data

import "strings"

type Class struct {
	Name string
}

func (i Class) Is(value Value) bool {
	switch c := value.(type) {
	case *ClassValue:
		if i.Name == "iterable" {
			return isIterableClassValue(c)
		}
		return isClassValueInstanceOf(i.Name, c.Class, c.GetVM())
	case *ThisValue:
		if i.Name == "iterable" {
			return isIterableClassStmt(c.Class, c.GetVM())
		}
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
		// 无 Object 时（Go 内部 NewErrorThrowByName）按 Name 对齐 Error/Exception 继承链。
		if c.Object != nil {
			return isClassValueInstanceOf(i.Name, c.Object.Class, c.Object.GetVM())
		}
		catchBase := i.Name
		if idx := strings.LastIndex(i.Name, "\\"); idx >= 0 {
			catchBase = i.Name[idx+1:]
		}
		if catchBase == "Throwable" || catchBase == "Exception" || catchBase == "Error" {
			return true
		}
		throwBase := c.Name
		if throwBase == "" {
			throwBase = "Exception"
		}
		if idx := strings.LastIndex(throwBase, "\\"); idx >= 0 {
			throwBase = throwBase[idx+1:]
		}
		return throwBase == catchBase || phpInternalThrowExtends(throwBase, catchBase)
	}

	return false
}

// isIterableClassValue 对齐 PHP iterable：array | Traversable（含 Generator）
func isIterableClassValue(c *ClassValue) bool {
	if c == nil || c.Class == nil {
		return false
	}
	return isIterableClassStmt(c.Class, c.GetVM())
}

func isIterableClassStmt(class ClassStmt, vm VM) bool {
	if class == nil {
		return false
	}
	name := class.GetName()
	if name == "Generator" || strings.HasSuffix(name, "\\Generator") {
		return true
	}
	for _, impl := range class.GetImplements() {
		base := impl
		if idx := strings.LastIndex(impl, "\\"); idx >= 0 {
			base = impl[idx+1:]
		}
		switch base {
		case "Iterator", "IteratorAggregate", "Traversable", "Generator":
			return true
		}
		if interfaceExtends(vm, impl, "Traversable") ||
			interfaceExtends(vm, impl, "Iterator") ||
			interfaceExtends(vm, impl, "IteratorAggregate") {
			return true
		}
	}
	return false
}

func normalizeTypeClassName(name string) string {
	return strings.TrimPrefix(name, "\\")
}

// isClassValueInstanceOf 检查一个 ClassStmt 是否实现了目标类型（类名或接口名）
func isClassValueInstanceOf(target string, class ClassStmt, vm VM) bool {
	target = normalizeTypeClassName(target)
	if class == nil {
		return false
	}
	if target == normalizeTypeClassName(class.GetName()) {
		return true
	}
	for _, s := range class.GetImplements() {
		if target == normalizeTypeClassName(s) {
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
	check = normalizeTypeClassName(check)
	for extend != nil {
		parent := normalizeTypeClassName(*extend)
		if check == parent {
			return true
		}
		if vm == nil {
			return builtinExceptionExtends(parent, check)
		}
		c, ok := vm.GetClass(parent)
		if !ok {
			c, ok = vm.GetClass(*extend)
		}
		extend = nil
		if ok {
			if check == normalizeTypeClassName(c.GetName()) {
				return true
			}
			for _, s := range c.GetImplements() {
				if check == normalizeTypeClassName(s) || interfaceExtends(vm, s, check) {
					return true
				}
			}
			extend = c.GetExtend()
			if extend == nil && builtinExceptionExtends(parent, check) {
				return true
			}
		} else if builtinExceptionExtends(parent, check) {
			return true
		}
	}
	return false
}

func builtinExceptionExtends(parent, target string) bool {
	parent = normalizeTypeClassName(parent)
	target = normalizeTypeClassName(target)
	switch parent {
	case "PDOException", "RuntimeException":
		return target == "RuntimeException" || target == "Exception" || target == "Throwable"
	case "Exception":
		return target == "Throwable"
	case "Error", "ValueError", "TypeError":
		return target == "Error" || target == "Throwable"
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

// phpInternalThrowExtends 判断 Go 侧按名字抛出的异常是否属于 catch 类型（无 Object 时）。
func phpInternalThrowExtends(throwName, catchName string) bool {
	parent := map[string]string{
		"Error":                 "Throwable",
		"ValueError":            "Error",
		"TypeError":             "Error",
		"ArgumentCountError":    "TypeError",
		"ArithmeticError":       "Error",
		"DivisionByZeroError":   "ArithmeticError",
		"UnhandledMatchError":   "Error",
		"ParseError":            "Error",
		"CompileError":          "Error",
		"Exception":             "Throwable",
		"ErrorException":        "Exception",
		"RuntimeException":      "Exception",
		"LogicException":        "Exception",
		"InvalidArgumentException": "LogicException",
		"UnexpectedValueException": "RuntimeException",
		"BadMethodCallException": "LogicException",
		"DomainException":       "LogicException",
		"OverflowException":     "RuntimeException",
		"RangeException":        "RuntimeException",
		"UnderflowException":    "RuntimeException",
		"LengthException":       "LogicException",
		"OutOfBoundsException":  "RuntimeException",
		"OutOfRangeException":   "LogicException",
		"PDOException":          "RuntimeException",
		"ReflectionException":   "Exception",
		"JsonException":         "Exception",
	}
	cur := throwName
	for i := 0; i < 8 && cur != ""; i++ {
		if cur == catchName {
			return true
		}
		next, ok := parent[cur]
		if !ok {
			return false
		}
		cur = next
	}
	return false
}
