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
	_, class := getReflectionClassInfo(ctx)
	for class != nil {
		for _, implemented := range class.GetImplements() {
			if implemented == "UnitEnum" || implemented == "BackedEnum" {
				return data.NewBoolValue(true), nil
			}
		}
		parent := class.GetExtend()
		if parent == nil {
			break
		}
		if *parent == "UnitEnum" || *parent == "BackedEnum" {
			return data.NewBoolValue(true), nil
		}
		next, control := ctx.GetVM().GetOrLoadClass(*parent)
		if control != nil {
			return data.NewBoolValue(false), nil
		}
		class = next
	}
	return data.NewBoolValue(false), nil
}
