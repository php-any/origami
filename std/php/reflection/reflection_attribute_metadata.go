package reflection

import "github.com/php-any/origami/data"

type ReflectionAttributeMetadataMethod struct {
	ReflectionAttributeGetNameMethod
	name string
}

func (m *ReflectionAttributeMetadataMethod) GetName() string { return m.name }
func (m *ReflectionAttributeMetadataMethod) GetReturnType() data.Types {
	if m.name == "getTarget" {
		return data.Int{}
	}
	return data.Bool{}
}
func (m *ReflectionAttributeMetadataMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	object := ctx.(*data.ClassMethodContext)
	class := object.Class.(*ReflectionAttributeClass)
	if m.name == "getTarget" {
		return data.NewIntValue(class.target), nil
	}
	return data.NewBoolValue(class.repeated), nil
}
