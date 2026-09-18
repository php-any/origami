package node

import "github.com/php-any/origami/data"

func storeClassStatic(class data.ClassStmt, name string, value data.Value) bool {
	if class == nil {
		return false
	}
	if data.StoreRequestStatic(class.GetName(), name, value) {
		return true
	}
	switch c := class.(type) {
	case *ClassStatement:
		c.StaticProperty.Store(name, value)
		return true
	case *AbstractClassStatement:
		c.StaticProperty.Store(name, value)
		return true
	case *ClassGeneric:
		c.StaticProperty.Store(name, value)
		return true
	}
	return false
}
