package data

import "strings"

// TypeNameEqual compares PHP class names: ASCII case insensitive, with an
// optional leading namespace separator. Non-ASCII bytes remain significant.
func TypeNameEqual(a, b string) bool {
	if a == b {
		return true
	}
	if len(a) != 0 && a[0] == '\\' {
		a = a[1:]
	}
	if len(b) != 0 && b[0] == '\\' {
		b = b[1:]
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
	}
	return true
}

// OriginalClass preserves the identity of aliases without changing lookup names.
type OriginalClass interface{ OriginalClass() ClassStmt }

func originalClass(c ClassStmt) ClassStmt {
	for c != nil {
		alias, ok := c.(OriginalClass)
		if !ok {
			break
		}
		c = alias.OriginalClass()
	}
	return c
}

func canonicalTarget(target string, vm VM) string {
	target = strings.TrimPrefix(target, "\\")
	if vm != nil {
		if c, ok := vm.GetClass(target); ok && c != nil {
			return originalClass(c).GetName()
		}
	}
	return target
}

// NominalIsA examines already linked metadata. Type comparisons never autoload.
func NominalIsA(class ClassStmt, target string, vm VM) bool {
	if class == nil {
		return false
	}
	if class.GetName() == target {
		return true
	}
	// The usual parent check needs no graph traversal or registry lookup.
	if parent := class.GetExtend(); parent != nil && *parent == target {
		return true
	}
	if nominalIsA(class, target, vm) {
		return true
	}
	canonical := canonicalTarget(target, vm)
	return !TypeNameEqual(canonical, target) && nominalIsA(class, canonical, vm)
}

// SameNominalClass excludes aliases of the source from strict subclass checks.
func SameNominalClass(class ClassStmt, target string, vm VM) bool {
	if class == nil {
		return false
	}
	if TypeNameEqual(class.GetName(), target) {
		return true
	}
	return TypeNameEqual(originalClass(class).GetName(), canonicalTarget(target, vm))
}

func nominalIsA(class ClassStmt, target string, vm VM) bool {
	if class == nil {
		return false
	}
	class = originalClass(class)
	var local [16]string
	visited := local[:0]
	for class != nil {
		name := class.GetName()
		if name == target || TypeNameEqual(name, target) {
			return true
		}
		for _, seen := range visited {
			if TypeNameEqual(seen, name) {
				return false
			}
		}
		visited = append(visited, name)
		for _, impl := range class.GetImplements() {
			if InterfaceIsA(impl, target, vm) {
				return true
			}
		}
		parent := class.GetExtend()
		if parent == nil {
			return builtinTypeIsA(name, target)
		}
		if *parent == target || TypeNameEqual(*parent, target) {
			return true
		}
		if vm == nil {
			return builtinTypeIsA(*parent, target)
		}
		next, ok := vm.GetClass(*parent)
		if !ok {
			return builtinTypeIsA(*parent, target)
		}
		class = originalClass(next)
	}
	return false
}

// InterfaceIsA supports deep and diamond inheritance and terminates on cycles.
func InterfaceIsA(source, target string, vm VM) bool {
	var local [16]string
	return interfaceIsA(source, target, vm, local[:0])
}

func interfaceIsA(source, target string, vm VM, path []string) bool {
	if source == target || TypeNameEqual(source, target) {
		return true
	}
	for _, seen := range path {
		if TypeNameEqual(seen, source) {
			return false
		}
	}
	if vm == nil {
		return builtinTypeIsA(source, target)
	}
	iface, ok := vm.GetInterface(source)
	if !ok {
		return builtinTypeIsA(source, target)
	}
	path = append(path, source)
	if len(iface.GetExtends()) == 0 {
		return builtinTypeIsA(iface.GetName(), target)
	}
	for _, parent := range iface.GetExtends() {
		if interfaceIsA(parent, target, vm, path) {
			return true
		}
	}
	return false
}

var internalTypeParents = map[string]string{
	"Error": "Throwable", "ValueError": "Error", "TypeError": "Error",
	"ArgumentCountError": "TypeError", "ArithmeticError": "Error",
	"DivisionByZeroError": "ArithmeticError", "UnhandledMatchError": "Error",
	"CompileError": "Error", "ParseError": "CompileError", "AssertionError": "Error", "FiberError": "Error",
	"Exception": "Throwable", "ErrorException": "Exception",
	"RuntimeException": "Exception", "LogicException": "Exception",
	"InvalidArgumentException": "LogicException", "UnexpectedValueException": "RuntimeException",
	"BadFunctionCallException": "LogicException", "BadMethodCallException": "BadFunctionCallException",
	"DomainException": "LogicException", "OverflowException": "RuntimeException",
	"RangeException": "RuntimeException", "UnderflowException": "RuntimeException",
	"LengthException": "LogicException", "OutOfBoundsException": "RuntimeException",
	"OutOfRangeException": "LogicException", "PDOException": "RuntimeException",
	"ReflectionException": "Exception", "JsonException": "Exception",
	"Generator": "Iterator", "Iterator": "Traversable", "IteratorAggregate": "Traversable",
}

func builtinTypeIsA(source, target string) bool {
	for parent := internalTypeParents[source]; parent != ""; parent = internalTypeParents[parent] {
		if TypeNameEqual(parent, target) {
			return true
		}
	}
	return false
}

// InternalTypeIsA applies the same ancestry rules to internal throws and native
// values which do not yet carry a ClassStmt. Namespace suffixes are not aliases.
func InternalTypeIsA(source, target string) bool {
	if TypeNameEqual(source, target) {
		return true
	}
	source = strings.TrimPrefix(source, "\\")
	if _, ok := internalTypeParents[source]; !ok {
		for name := range internalTypeParents {
			if TypeNameEqual(name, source) {
				source = name
				break
			}
		}
	}
	return builtinTypeIsA(source, target)
}
