package node

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

func CheckDynamicPropertyCreation(object *data.ClassValue, name string, from data.From) data.Control {
	if classDeclarationFlags(object.Class)&(data.ClassReadonly|data.ClassEnum) != 0 {
		return data.NewErrorThrowByName(from, fmt.Errorf("Cannot create dynamic property %s::$%s", object.Class.GetName(), name), "Error")
	}
	return nil
}

// ReportDynamicPropertyCreation is called only after installing a new,
// undeclared PHP property. Native instance initialization uses raw storage.
func ReportDynamicPropertyCreation(ctx data.Context, object *data.ClassValue, name string, from data.From) data.Control {
	if ctx == nil || object.Class.GetName() == "__PHP_Incomplete_Class" || data.NominalIsA(object.Class, "stdClass", ctx.GetVM()) {
		return nil
	}
	for class := object.Class; class != nil; {
		if declaration := classStmtFromAny(class); declaration != nil {
			for _, annotation := range declaration.Annotations {
				if strings.EqualFold(annotation.Class.GetName(), "AllowDynamicProperties") {
					return nil
				}
			}
		}
		if class.GetExtend() == nil {
			break
		}
		class, _ = ctx.GetVM().GetClass(*class.GetExtend())
	}
	return data.EmitPHPError(ctx, 8192, fmt.Sprintf("Creation of dynamic property %s::$%s is deprecated", object.Class.GetName(), name), from)
}
