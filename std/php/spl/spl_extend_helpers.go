package spl

import "github.com/php-any/origami/data"

// splExtendGetMethod 在子类本地方法未命中时，�?GetExtend 链查找父类方法（内置 SPL 子类用）�?
func splExtendGetMethod(c data.ClassStmt, name string, local func(string) (data.Method, bool)) (data.Method, bool) {
	if local != nil {
		if m, ok := local(name); ok && m != nil {
			return m, true
		}
	}
	parent := splExtendResolveParent(c)
	if parent == nil {
		return nil, false
	}
	return parent.GetMethod(name)
}

// splExtendGetMethods 合并子类方法与父类方法（同名时子类优先）�?
func splExtendGetMethods(c data.ClassStmt, local []data.Method) []data.Method {
	parent := splExtendResolveParent(c)
	if parent == nil {
		return local
	}
	parentMethods := parent.GetMethods()
	if len(local) == 0 {
		return parentMethods
	}
	seen := make(map[string]struct{}, len(local))
	for _, m := range local {
		seen[m.GetName()] = struct{}{}
	}
	merged := make([]data.Method, 0, len(local)+len(parentMethods))
	merged = append(merged, local...)
	for _, m := range parentMethods {
		if _, ok := seen[m.GetName()]; !ok {
			merged = append(merged, m)
		}
	}
	return merged
}

func splExtendResolveParent(c data.ClassStmt) data.ClassStmt {
	if c == nil || c.GetExtend() == nil {
		return nil
	}
	switch *c.GetExtend() {
	case "SplDoublyLinkedList":
		return NewSplDoublyLinkedListClass()
	case "SplHeap":
		return NewSplHeapClass()
	case "ArrayIterator":
		return NewArrayIteratorClass()
	default:
		return nil
	}
}
