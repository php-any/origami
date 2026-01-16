package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsCallableFunction 实现 is_callable 函数
type IsCallableFunction struct{}

func NewIsCallableFunction() data.FuncStmt {
	return &IsCallableFunction{}
}

func (f *IsCallableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	// syntaxOnlyValue, _ := ctx.GetIndexValue(1)
	// callableNameValue, _ := ctx.GetIndexValue(2)

	// Check if value is callable
	// In Origami, we might check if it implements CallableValue or similar.
	// Or if it's a string representing a function, or array [obj, method], or closure.

	if value == nil {
		return data.NewBoolValue(false), nil
	}

	// 1. Closure / Function object
	if _, ok := value.(data.CallableValue); ok {
		return data.NewBoolValue(true), nil
	}

	if _, ok := value.(data.GetName); ok {
		return data.NewBoolValue(true), nil
	}

	// 2. String (function name)
	if str, ok := value.(data.AsString); ok {
		name := str.AsString()
		// Check if function exists
		_, exists := ctx.GetVM().GetFunc(name)
		if exists {
			return data.NewBoolValue(true), nil
		}
		// Also check static method string "Class::method"
		// Not implemented in simple check, but standard PHP supports it.
	}

	// 3. Array [obj|class, method]
	if arr, ok := value.(*data.ArrayValue); ok {
		if len(arr.Value) == 2 {
			objOrClass := arr.Value[0]
			methodNameVal := arr.Value[1]

			if methodName, ok := methodNameVal.(data.AsString); ok {
				method := methodName.AsString()

				// [object, method]
				if obj, ok := objOrClass.(*data.ClassValue); ok {
					// Check if object has method
					if _, found := obj.GetMethod(method); found {
						return data.NewBoolValue(true), nil
					}
				}

				// [class_string, method] (Static call)
				if classStr, ok := objOrClass.(data.AsString); ok {
					className := classStr.AsString()
					// Check if class exists and has method
					if class, _ := ctx.GetVM().GetClass(className); class != nil {
						if _, found := class.GetMethod(method); found {
							return data.NewBoolValue(true), nil
						}
					}
				}
			}
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *IsCallableFunction) GetName() string {
	return "is_callable"
}

func (f *IsCallableFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "syntax_only", 1, node.NewBooleanLiteral(nil, false), nil),
		node.NewParameter(nil, "callable_name", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *IsCallableFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "syntax_only", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "callable_name", 2, data.NewBaseType("string")),
	}
}
