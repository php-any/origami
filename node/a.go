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
	case *data.ThisValue:
		return TryGetCallClassName(c.ClassValue)
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
	case *ReturnStatement:
		return fmt.Sprintf("return %s", TryGetCallClassName(c.Value))
	case *TernaryExpression:
		return fmt.Sprintf("%s ? %s : %s", TryGetCallClassName(c.Condition), TryGetCallClassName(c.TrueValue), TryGetCallClassName(c.FalseValue))
	case *CallObjectMethod:
		return fmt.Sprintf("%s->%s", TryGetCallClassName(c.Object), c.Method)
	case *BinaryAssignVariable:
		return fmt.Sprintf("%s = %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *CallExpression:
		return fmt.Sprintf("%s(%s)", "call", c.FunName)
	}

	return "TODO"
}
