package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// TryGetCallClassName 尽可能尝试获取名称
func TryGetCallClassName(call data.GetValue) string {
	switch c := call.(type) {
	case *VariableExpression:
		return fmt.Sprintf("$%s", c.Name)
	case *data.ClassValue:
		return TryGetCallClassName(c.Class)
	case data.ClassStmt:
		return c.GetName()
	case *data.NullValue:
		return "null"
	case *CallStaticMethod:
		return fmt.Sprintf("%s::%s", TryGetCallClassName(c.stmt), c.Method)
	case *This:
		return "this"
	case *IndexExpression:
		return fmt.Sprintf("%s[%s]", TryGetCallClassName(c.Array), TryGetCallClassName(c.Index))
	case *CallObjectProperty:
		return fmt.Sprintf("%s->%s", TryGetCallClassName(c.Object), c.Property)
	}

	return "TODO"
}
