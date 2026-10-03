package reflection

import "github.com/php-any/origami/data"

type ReturnsReferenceMethod struct{}

func (*ReturnsReferenceMethod) GetName() string               { return "returnsReference" }
func (*ReturnsReferenceMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (*ReturnsReferenceMethod) GetIsStatic() bool             { return false }
func (*ReturnsReferenceMethod) GetParams() []data.GetValue    { return nil }
func (*ReturnsReferenceMethod) GetVariables() []data.Variable { return nil }
func (*ReturnsReferenceMethod) GetReturnType() data.Types     { return data.TypeBool }
func (*ReturnsReferenceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var declaration data.FuncStmt
	if object, ok := ctx.(*data.ClassMethodContext); ok && object.ObjectValue != nil {
		properties := object.ObjectValue
		if value, ok := reflectionProperty(properties, "_function").(*data.FuncValue); ok {
			declaration = value.Value
		}
		if name, ok := reflectionProperty(properties, "_funcName").(*data.StringValue); ok {
			declaration, _ = ctx.GetVM().GetFunc(name.Value)
		}
	}
	if declaration == nil {
		_, _, method := getReflectionMethodInfo(ctx)
		declaration = method
	}
	if ref, ok := declaration.(interface{ ReturnsByReference() bool }); ok {
		return data.NewBoolValue(ref.ReturnsByReference()), nil
	}
	return data.NewBoolValue(false), nil
}
