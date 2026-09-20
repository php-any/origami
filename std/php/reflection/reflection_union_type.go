package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionUnionTypeClass 对齐 PHP ReflectionUnionType（PHP 8.0+）。
type ReflectionUnionTypeClass struct {
	node.Node
}

func (c *ReflectionUnionTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ReflectionUnionTypeClass) GetName() string { return "ReflectionUnionType" }

func (c *ReflectionUnionTypeClass) GetExtend() *string {
	parent := "ReflectionType"
	return &parent
}

func (c *ReflectionUnionTypeClass) GetImplements() []string { return nil }

func (c *ReflectionUnionTypeClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ReflectionUnionTypeClass) GetPropertyList() []data.Property { return nil }

func (c *ReflectionUnionTypeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &reflectionCompoundConstructMethod{}, true
	case "getTypes":
		return &ReflectionUnionTypeGetTypesMethod{}, true
	case "__toString":
		return &ReflectionTypeToStringMethod{}, true
	case "allowsNull":
		return &ReflectionTypeAllowsNullMethod{}, true
	}
	return nil, false
}

func (c *ReflectionUnionTypeClass) GetMethods() []data.Method {
	return []data.Method{
		&reflectionCompoundConstructMethod{},
		&ReflectionUnionTypeGetTypesMethod{},
		&ReflectionTypeToStringMethod{},
		&ReflectionTypeAllowsNullMethod{},
	}
}

func (c *ReflectionUnionTypeClass) GetConstruct() data.Method {
	return &reflectionCompoundConstructMethod{}
}

func newReflectionUnionType(ctx data.Context, u data.UnionType) *data.ClassValue {
	return newCompoundReflectionType(ctx, &ReflectionUnionTypeClass{}, u.Types, u.String())
}

type ReflectionUnionTypeGetTypesMethod struct{}

func (m *ReflectionUnionTypeGetTypesMethod) GetName() string            { return "getTypes" }
func (m *ReflectionUnionTypeGetTypesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionUnionTypeGetTypesMethod) GetIsStatic() bool          { return false }
func (m *ReflectionUnionTypeGetTypesMethod) GetReturnType() data.Types  { return data.Arrays{} }
func (m *ReflectionUnionTypeGetTypesMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionUnionTypeGetTypesMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionUnionTypeGetTypesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return getCompoundMemberTypes(ctx), nil
}

type reflectionCompoundConstructMethod struct{}

func (m *reflectionCompoundConstructMethod) GetName() string            { return "__construct" }
func (m *reflectionCompoundConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *reflectionCompoundConstructMethod) GetIsStatic() bool          { return false }
func (m *reflectionCompoundConstructMethod) GetReturnType() data.Types  { return nil }
func (m *reflectionCompoundConstructMethod) GetParams() []data.GetValue { return nil }
func (m *reflectionCompoundConstructMethod) GetVariables() []data.Variable {
	return nil
}
func (m *reflectionCompoundConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

func newCompoundReflectionType(ctx data.Context, class data.ClassStmt, members []data.Types, typeName string) *data.ClassValue {
	typeValue := data.NewClassValue(class, ctx.CreateBaseContext())
	named := make([]data.Value, 0, len(members))
	for _, t := range members {
		if t == nil {
			continue
		}
		named = append(named, newReflectionNamedType(ctx, t))
	}
	typeValue.ObjectValue.SetProperty("_typeName", data.NewStringValue(typeName))
	typeValue.ObjectValue.SetProperty("_allowsNull", data.NewBoolValue(compoundAllowsNull(members)))
	typeValue.ObjectValue.SetProperty("_memberTypes", data.NewArrayValue(named))
	return typeValue
}

func getCompoundMemberTypes(ctx data.Context) data.Value {
	empty := data.NewArrayValue(nil)
	objCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || objCtx.ObjectValue == nil {
		return empty
	}
	props := objCtx.ObjectValue.GetProperties()
	v, has := props["_memberTypes"]
	if !has || v == nil {
		return empty
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return av
	}
	return empty
}
