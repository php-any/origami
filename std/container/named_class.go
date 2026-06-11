package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewNamedClass() data.ClassStmt {
	return &NamedClass{construct: &NamedConstructMethod{}}
}

type NamedClass struct {
	node.Node
	construct data.Method
}

func (c *NamedClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&NamedClass{construct: &NamedConstructMethod{}}, ctx.CreateBaseContext()), nil
}
func (c *NamedClass) GetName() string    { return "Container\\Named" }
func (c *NamedClass) GetExtend() *string { return nil }
func (c *NamedClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetParameter}
}
func (c *NamedClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *NamedClass) GetPropertyList() []data.Property           { return nil }
func (c *NamedClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *NamedClass) GetMethods() []data.Method { return []data.Method{c.construct} }
func (c *NamedClass) GetConstruct() data.Method { return c.construct }

type NamedConstructMethod struct{}

func (m *NamedConstructMethod) GetName() string            { return "__construct" }
func (m *NamedConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *NamedConstructMethod) GetIsStatic() bool          { return false }
func (m *NamedConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}
func (m *NamedConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}
func (m *NamedConstructMethod) GetReturnType() data.Types { return nil }
func (m *NamedConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, acl := annotationStringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	if name == "" {
		return nil, utils.NewThrow(errors.New("Container\\Named 缺少 name 参数"))
	}
	param, className, acl := annotationTargetParameter(ctx)
	if acl != nil {
		return nil, acl
	}
	if className == "" {
		return nil, utils.NewThrow(errors.New("Container\\Named 缺少所属类信息"))
	}
	metadataMarkConstructorInject(className, param.Index, param.Name, name, false)
	return nil, nil
}
