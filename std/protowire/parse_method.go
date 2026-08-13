package protowire

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ParseMethod implements Protowire::parse()
// Parse binary protobuf data. If a class name is given, uses @Field
// annotations on the class properties to deserialize into an object.
// Otherwise returns raw field array.
type ParseMethod struct{}

func NewParseMethod() data.Method {
	return &ParseMethod{}
}

func (m *ParseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	dataVal, ok := ctx.GetIndexValue(0)
	if !ok || dataVal == nil {
		return data.NewNullValue(), nil
	}
	inputStr, ok := dataVal.(data.AsString)
	if !ok {
		return data.NewNullValue(), nil
	}

	arg1Val, _ := ctx.GetIndexValue(1)

	// Determine: is arg1 a class name (string) or options (array)?
	className := ""
	opts := &ParseOptions{}

	if arg1Val != nil {
		if _, isNull := arg1Val.(*data.NullValue); !isNull {
			if _, isStr := arg1Val.(*data.StringValue); isStr {
				// arg1 is a string → className
				className = arg1Val.AsString()
				if optsVal, ok2 := ctx.GetIndexValue(2); ok2 && optsVal != nil {
					if err := parseOptionsFromPHPValue(optsVal, opts); err != nil {
						return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::parse: %w", err))
					}
				}
			} else {
				// arg1 is NOT a string → treat as options
				if err := parseOptionsFromPHPValue(arg1Val, opts); err != nil {
					return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::parse: %w", err))
				}
			}
		}
	}

	if className == "" {
		// No className → raw parse
		return rawParseFields([]byte(inputStr.AsString()), opts)
	}

	// Annotation-based deserialization
	vm := ctx.GetVM()
	classStmt, _ := vm.GetClass(className)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::parse: class %q not found", className))
	}

	fieldMap := readFieldAnnotations(classStmt)
	rawFields, err := ParseRawFields([]byte(inputStr.AsString()), opts)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::parse: %w", err))
	}

	instance := data.NewClassValue(classStmt, nil)
	for _, f := range rawFields {
		if propName, found := fieldMap[f.Number]; found {
			instance.SetProperty(propName, goValueToPHPValue(f.Value))
		}
	}
	return instance, nil
}

func (m *ParseMethod) GetName() string {
	return "parse"
}

func (m *ParseMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ParseMethod) GetIsStatic() bool {
	return true
}

func (m *ParseMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, data.String{}),
		node.NewParameter(nil, "classNameOrOptions", 1, data.NewNullValue(), nil),
		node.NewParameter(nil, "options", 2, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *ParseMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.String{}),
		node.NewVariable(nil, "classNameOrOptions", 1, nil),
		node.NewVariable(nil, "options", 2, data.NewBaseType("array")),
	}
}

func (m *ParseMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// rawParseFields parses binary data and returns the raw field array.
func rawParseFields(input []byte, opts *ParseOptions) (data.GetValue, data.Control) {
	fields, err := ParseRawFields(input, opts)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::parse: %w", err))
	}
	return fieldsToPHPArray(fields), nil
}

// readFieldAnnotations scans a ClassStmt's properties for @Field annotations
// and returns a map of fieldNumber → propertyName.
func readFieldAnnotations(classStmt data.ClassStmt) map[int32]string {
	fieldMap := make(map[int32]string)
	for _, prop := range classStmt.GetPropertyList() {
		if cp, ok := prop.(*node.ClassProperty); ok {
			for _, ann := range cp.Annotations {
				if ann == nil || ann.Class == nil {
					continue
				}
				if ann.Class.GetName() == "Protowire\\Annotation\\Field" {
					props := ann.GetProperties()
					numVal, hasNum := props["number"]
					if !hasNum || numVal == nil {
						continue
					}
					ai, ok := numVal.(data.AsInt)
					if !ok {
						continue
					}
					n, err := ai.AsInt()
					if err != nil {
						continue
					}
					fieldMap[int32(n)] = cp.GetName()
				}
			}
		}
	}
	return fieldMap
}
