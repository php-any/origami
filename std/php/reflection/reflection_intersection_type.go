package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionIntersectionTypeClass 对齐 PHP ReflectionIntersectionType（PHP 8.1+）。
type ReflectionIntersectionTypeClass struct {
	node.Node
}

func (c *ReflectionIntersectionTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ReflectionIntersectionTypeClass) GetName() string { return "ReflectionIntersectionType" }

func (c *ReflectionIntersectionTypeClass) GetExtend() *string {
	parent := "ReflectionType"
	return &parent
}

func (c *ReflectionIntersectionTypeClass) GetImplements() []string { return nil }

func (c *ReflectionIntersectionTypeClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ReflectionIntersectionTypeClass) GetPropertyList() []data.Property { return nil }

func (c *ReflectionIntersectionTypeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &reflectionCompoundConstructMethod{}, true
	case "getTypes":
		return &ReflectionIntersectionTypeGetTypesMethod{}, true
	case "__toString":
		return &ReflectionTypeToStringMethod{}, true
	case "allowsNull":
		return &ReflectionTypeAllowsNullMethod{}, true
	}
	return nil, false
}

func (c *ReflectionIntersectionTypeClass) GetMethods() []data.Method {
	return []data.Method{
		&reflectionCompoundConstructMethod{},
		&ReflectionIntersectionTypeGetTypesMethod{},
		&ReflectionTypeToStringMethod{},
		&ReflectionTypeAllowsNullMethod{},
	}
}

func (c *ReflectionIntersectionTypeClass) GetConstruct() data.Method {
	return &reflectionCompoundConstructMethod{}
}

func newReflectionIntersectionType(ctx data.Context, i data.IntersectionType) *data.ClassValue {
	return newCompoundReflectionType(ctx, &ReflectionIntersectionTypeClass{}, i.Types, i.String())
}

type ReflectionIntersectionTypeGetTypesMethod struct{}

func (m *ReflectionIntersectionTypeGetTypesMethod) GetName() string { return "getTypes" }
func (m *ReflectionIntersectionTypeGetTypesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionIntersectionTypeGetTypesMethod) GetIsStatic() bool          { return false }
func (m *ReflectionIntersectionTypeGetTypesMethod) GetReturnType() data.Types  { return data.Arrays{} }
func (m *ReflectionIntersectionTypeGetTypesMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionIntersectionTypeGetTypesMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionIntersectionTypeGetTypesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return getCompoundMemberTypes(ctx), nil
}
