package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewInjectClass() data.ClassStmt {
	return &InjectClass{construct: &InjectConstructMethod{}}
}

type InjectClass struct {
	node.Node
	construct data.Method
}

func (c *InjectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&InjectClass{construct: &InjectConstructMethod{}}, ctx.CreateBaseContext()), nil
}
func (c *InjectClass) GetName() string    { return "Container\\Inject" }
func (c *InjectClass) GetExtend() *string { return nil }
func (c *InjectClass) GetImplements() []string {
	return []string{node.TypeMacro, node.TypeTargetParameter}
}
func (c *InjectClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *InjectClass) GetPropertyList() []data.Property           { return nil }
func (c *InjectClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *InjectClass) GetMethods() []data.Method { return []data.Method{c.construct} }
func (c *InjectClass) GetConstruct() data.Method { return c.construct }

type InjectConstructMethod struct{}

func (m *InjectConstructMethod) GetName() string            { return "__construct" }
func (m *InjectConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *InjectConstructMethod) GetIsStatic() bool          { return false }
func (m *InjectConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "service", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}
func (m *InjectConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "service", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}
func (m *InjectConstructMethod) GetReturnType() data.Types { return nil }
func (m *InjectConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	service, _ := annotationStringArg(ctx, 0)
	param, className, acl := annotationTargetParameter(ctx)
	if acl != nil {
		return nil, acl
	}
	if className == "" {
		return nil, utils.NewThrow(errors.New("Container\\Inject 缺少所属类信息"))
	}
	metadataMarkConstructorInject(className, param.Index, param.Name, service, true)
	return nil, nil
}
