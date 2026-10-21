package data

// MethodDelegate supports native object handlers such as SPL outer iterators.
// Delegation does not add declarations to class reflection.
type MethodDelegate interface {
	MethodDelegate(object *ClassValue) *ClassValue
}

func DelegatedMethodReceiver(object *ClassValue) *ClassValue {
	for class := object.Class; class != nil; {
		if handler, ok := class.(MethodDelegate); ok {
			return handler.MethodDelegate(object)
		}
		parent := class.GetExtend()
		if parent == nil {
			break
		}
		class, _ = object.GetVM().GetClass(*parent)
	}
	return nil
}

func ObjectMethodExists(object *ClassValue, name string) bool {
	seen := map[*ClassValue]bool{}
	for object != nil && !seen[object] {
		seen[object] = true
		if _, ok := object.GetMethod(name); ok {
			return true
		}
		object = DelegatedMethodReceiver(object)
	}
	return false
}
