package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func classDeclarationFlags(class data.ClassStmt) data.ClassFlags {
	if metadata, ok := class.(interface{ DeclarationFlags() data.ClassFlags }); ok {
		return metadata.DeclarationFlags()
	}
	return 0
}

// Declaration checks run once on registration, outside calls/property reads.
func ValidateClassHierarchy(vm data.VM, class *ClassStatement) data.Control {
	if class.Extends == nil {
		return nil
	}
	parent, ctl := vm.GetOrLoadClass(*class.Extends)
	if ctl != nil {
		return ctl
	}
	if parent == nil {
		return nil
	}
	if classDeclarationFlags(parent)&data.ClassFinal != 0 {
		return data.NewCompileFatal(class.GetFrom(), fmt.Sprintf("Class %s cannot extend final class %s", class.GetName(), parent.GetName()))
	}
	if (class.Flags&data.ClassReadonly != 0) != (classDeclarationFlags(parent)&data.ClassReadonly != 0) {
		return data.NewCompileFatal(class.GetFrom(), "Readonly and non-readonly classes cannot inherit each other")
	}
	var failure data.Control
	forEachClassMethod(class, func(method data.Method) bool {
		for owner := parent; owner != nil; {
			inherited, found := owner.GetMethod(method.GetName())
			if !found {
				if static, ok := owner.(data.GetStaticMethod); ok {
					inherited, found = static.GetStaticMethod(method.GetName())
				}
			}
			if found && data.MethodDeclarationFlags(inherited)&data.MethodFinal != 0 && method != inherited && (inherited.GetModifier() != data.ModifierPrivate || method.GetName() == "__construct") {
				failure = data.NewCompileFatal(class.GetFrom(), fmt.Sprintf("Cannot override final method %s::%s()", owner.GetName(), method.GetName()))
				return false
			}
			if owner.GetExtend() == nil {
				break
			}
			owner, _ = vm.GetClass(*owner.GetExtend())
		}
		return true
	})
	return failure
}
