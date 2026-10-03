package reflection

import "github.com/php-any/origami/data"

// ReflectionClassIsEnumMethod 实现 PHP 8.1 ReflectionClass::isEnum。
type ReflectionClassIsEnumMethod struct{}

func (m *ReflectionClassIsEnumMethod) GetName() string            { return "isEnum" }
func (m *ReflectionClassIsEnumMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionClassIsEnumMethod) GetIsStatic() bool          { return false }
func (m *ReflectionClassIsEnumMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionClassIsEnumMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionClassIsEnumMethod) GetReturnType() data.Types { return data.Bool{} }

func (m *ReflectionClassIsEnumMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(reflectionClassFlags(ctx)&data.ClassEnum != 0), nil
}
