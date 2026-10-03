package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"reflect"
)

// This compile-time walk handles nested branches but stops at nested function
// boundaries. Return declarations are validated even for unreachable statements.
func validateReturnBody(body []data.GetValue, ret data.Types) data.Control {
	if ret == nil || ret == data.TypeInvalid {
		return nil
	}
	generator := false
	_ = walkSyntax(reflect.ValueOf(body), func(value any) (bool, data.Control) {
		switch value.(type) {
		case *node.FunctionStatement, *node.LambdaExpression, *node.ClassMethod, *node.ClassRegisterStmt:
			return false, nil
		case *node.YieldStatement, *node.YieldFromStatement:
			generator = true
		}
		return true, nil
	})
	if generator {
		// A generator's return declaration describes the Generator object,
		// not values passed to getReturn(). Bare return is legal here.
		ref := data.DeclaredTypeRef(ret)
		if ref.Kind() == data.TypeKindBuiltin && ref != data.TypeMixed && ref != data.TypeObject && ref != data.TypeIterable {
			return data.NewCompileFatal(nil, "Generator return type must be a supertype of Generator")
		}
		return nil
	}
	return walkSyntax(reflect.ValueOf(body), func(value any) (bool, data.Control) {
		switch statement := value.(type) {
		case *node.FunctionStatement, *node.LambdaExpression, *node.ClassMethod, *node.ClassRegisterStmt:
			return false, nil
		case *node.ReturnStatement:
			if ret == data.TypeNever {
				return false, data.NewCompileFatal(statement.GetFrom(), "A never-returning function must not return")
			}
			if ret == data.TypeVoid && statement.Value != nil {
				return false, data.NewCompileFatal(statement.GetFrom(), "A void function must not return a value")
			}
			if ret != data.TypeVoid && statement.Value == nil {
				return false, data.NewCompileFatal(statement.GetFrom(), "A function with a return type must return a value")
			}
		}
		return true, nil
	})
}

func validateReturnDeclarations(program *node.Program) data.Control {
	return walkSyntax(reflect.ValueOf(program), func(value any) (bool, data.Control) {
		switch function := value.(type) {
		case *node.LambdaExpression:
			return true, validateReturnBody(function.Body, function.Ret)
		case *node.FunctionStatement:
			return true, validateReturnBody(function.Body, function.Ret)
		case *node.ClassMethod:
			return true, validateReturnBody(function.Body, function.Ret)
		}
		return true, nil
	})
}
func walkSyntax(value reflect.Value, visit func(any) (bool, data.Control)) data.Control {
	return walkSyntaxOnce(value, visit, make(map[any]struct{}))
}

func walkSyntaxOnce(value reflect.Value, visit func(any) (bool, data.Control), seen map[any]struct{}) data.Control {
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		return walkSyntaxOnce(value.Elem(), visit, seen)
	}
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}
	if value.CanInterface() {
		if _, ok := value.Interface().(data.Value); ok {
			return nil
		}
		if value.Kind() == reflect.Ptr {
			pointer := value.Interface()
			if _, visited := seen[pointer]; visited {
				return nil
			}
			seen[pointer] = struct{}{}
		}
		descend, ctl := visit(value.Interface())
		if ctl != nil || !descend {
			return ctl
		}
	}
	switch value.Kind() {
	case reflect.Ptr:
		if !value.IsNil() {
			return walkSyntaxOnce(value.Elem(), visit, seen)
		}
	case reflect.Struct:
		// Only syntax records belong to this traversal. Runtime values, type
		// metadata and synchronization state may hold unrelated cyclic graphs.
		if value.Type().PkgPath() != "github.com/php-any/origami/node" {
			return nil
		}
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			if !field.IsExported() || field.Name == "Node" || field.Name == "FuncStmt" || field.Name == "Fun" {
				continue
			}
			if ctl := walkSyntaxOnce(value.Field(i), visit, seen); ctl != nil {
				return ctl
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if ctl := walkSyntaxOnce(value.Index(i), visit, seen); ctl != nil {
				return ctl
			}
		}
	case reflect.Map:
		iterator := value.MapRange()
		for iterator.Next() {
			if ctl := walkSyntaxOnce(iterator.Value(), visit, seen); ctl != nil {
				return ctl
			}
		}
	}
	return nil
}
