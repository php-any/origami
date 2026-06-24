package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewBindClass() data.ClassStmt {
	return &BindClass{construct: &BindConstructMethod{}}
}

type BindClass struct {
	node.Node
	construct data.Method
}

func (c *BindClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&BindClass{construct: &BindConstructMethod{}}, ctx.CreateBaseContext()), nil
}
func (c *BindClass) GetName() string    { return "Container\\Bind" }
func (c *BindClass) GetExtend() *string { return nil }
func (c *BindClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (c *BindClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *BindClass) GetPropertyList() []data.Property           { return nil }
func (c *BindClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *BindClass) GetMethods() []data.Method { return []data.Method{c.construct} }
func (c *BindClass) GetConstruct() data.Method { return c.construct }

type BindConstructMethod struct{}

func (m *BindConstructMethod) GetName() string            { return "__construct" }
func (m *BindConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *BindConstructMethod) GetIsStatic() bool          { return false }
func (m *BindConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "abstract", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewAnnotationTargetParameter(nil, 1),
	}
}
func (m *BindConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "abstract", 0, nil),
		node.NewAnnotationTargetVariable(nil, 1),
	}
}
func (m *BindConstructMethod) GetReturnType() data.Types { return nil }
func (m *BindConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	abstract, acl := annotationStringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	if abstract == "" {
		return nil, utils.NewThrow(errors.New("Bind 缺少 abstract 参数"))
	}
	cls, acl := annotationTargetClass(ctx)
	if acl != nil {
		return nil, acl
	}
	if e := activeEngine(ctx); e != nil {
		e.Bind(abstract, cls.Name)
	}
	return nil, nil
}
