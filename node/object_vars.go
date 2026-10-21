package node

import "github.com/php-any/origami/data"

// VisibleObjectVars reads initialized instance storage without invoking magic
// getters. Its visibility is the caller's lexical scope.
func VisibleObjectVars(ctx data.Context, object *data.ClassValue) *data.ArrayValue {
	out := data.NewArrayValueFromSlots(nil)
	object.RangeProperties(func(key string, value data.Value) bool {
		property, declared := object.GetPropertyStmt(key)
		if declared {
			if property.GetIsStatic() || !propertyAccessible(ctx, object, property, property.GetName()) {
				return true
			}
			key = property.GetName()
		}
		if value != nil {
			out.SetStringKey(key, value)
		}
		return true
	})
	return out
}
