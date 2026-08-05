package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodHasReturnTypeMethod 实现 ReflectionMethod::hasReturnType
type ReflectionMethodHasReturnTypeMethod struct{}

func (m *ReflectionMethodHasReturnTypeMethod) GetName() string { return "hasReturnType" }

func (m *ReflectionMethodHasReturnTypeMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ReflectionMethodHasReturnTypeMethod) GetIsStatic() bool { return false }

func (m *ReflectionMethodHasReturnTypeMethod) GetParams() []data.GetValue { return []data.GetValue{} }

func (m *ReflectionMethodHasReturnTypeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ReflectionMethodHasReturnTypeMethod) GetReturnType() data.Types { return data.Bool{} }

func (m *ReflectionMethodHasReturnTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(method.GetReturnType() != nil), nil
}
