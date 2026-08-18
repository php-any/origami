package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassGetStartLineMethod 实现 ReflectionClass::getStartLine(): int|false
type ReflectionClassGetStartLineMethod struct{}

func (m *ReflectionClassGetStartLineMethod) GetName() string               { return "getStartLine" }
func (m *ReflectionClassGetStartLineMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ReflectionClassGetStartLineMethod) GetIsStatic() bool             { return false }
func (m *ReflectionClassGetStartLineMethod) GetParams() []data.GetValue    { return nil }
func (m *ReflectionClassGetStartLineMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassGetStartLineMethod) GetReturnType() data.Types     { return data.Mixed{} }

func (m *ReflectionClassGetStartLineMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}
	return sourceStartLineFrom(classStmt), nil
}
