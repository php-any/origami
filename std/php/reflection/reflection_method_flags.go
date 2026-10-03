package reflection

import "github.com/php-any/origami/data"

type ReflectionMethodFlagMethod struct {
	name string
	flag data.MethodFlags
}

func (m *ReflectionMethodFlagMethod) GetName() string             { return m.name }
func (*ReflectionMethodFlagMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (*ReflectionMethodFlagMethod) GetIsStatic() bool             { return false }
func (*ReflectionMethodFlagMethod) GetParams() []data.GetValue    { return nil }
func (*ReflectionMethodFlagMethod) GetVariables() []data.Variable { return nil }
func (*ReflectionMethodFlagMethod) GetReturnType() data.Types     { return data.TypeBool }
func (m *ReflectionMethodFlagMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	return data.NewBoolValue(method != nil && data.MethodDeclarationFlags(method)&m.flag != 0), nil
}
