package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

func TryGetCallClassName(call data.GetValue) string {
	switch c := call.(type) {
	case *VariableExpression:
		return fmt.Sprintf("$%s", c.Name)
	case *data.ClassValue:
		return TryGetCallClassName(c.Class)
	case data.ClassStmt:
		return c.GetName()
	}

	return "TODO"
}
