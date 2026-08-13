package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ConstraintParam 约束注解构造参数定义。
type ConstraintParam struct {
	Name       string
	Type       data.Types
	DefaultVal data.GetValue
}

// ConstraintSpec 校验约束注解元数据。
type ConstraintSpec struct {
	FullName   string
	Repeatable bool
	Params     []ConstraintParam
}

// ConstraintClass 通用校验约束注解（Name、Min、Max、Email 等）。
type ConstraintClass struct {
	node.Node
	spec      ConstraintSpec
	state     map[string]data.GetValue
	construct data.Method
}

func newConstraintClass(spec ConstraintSpec) *ConstraintClass {
	c := &ConstraintClass{
		spec:  spec,
		state: make(map[string]data.GetValue),
	}
	c.construct = &constraintConstructMethod{class: c}
	return c
}

func (c *ConstraintClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	inst := newConstraintClass(c.spec)
	return data.NewClassValue(inst, ctx.CreateBaseContext()), nil
}

func (c *ConstraintClass) GetName() string { return c.spec.FullName }

func (c *ConstraintClass) GetExtend() *string { return nil }

func (c *ConstraintClass) GetImplements() []string {
	impls := []string{
		node.TypeFeature,
		node.TypeTargetProperty,
		node.TypeTargetParameter,
	}
	if c.spec.Repeatable {
		impls = append(impls, node.TypeRepeatable)
	}
	return impls
}

func (c *ConstraintClass) GetProperty(_ string) (data.Property, bool) { return nil, false }

func (c *ConstraintClass) GetPropertyList() []data.Property { return nil }

func (c *ConstraintClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}

func (c *ConstraintClass) GetMethods() []data.Method { return []data.Method{c.construct} }

func (c *ConstraintClass) GetConstruct() data.Method { return c.construct }

// Spec 返回约束定义。
func (c *ConstraintClass) Spec() ConstraintSpec { return c.spec }

// State 返回构造时写入的参数值（供编译模式序列化）。
func (c *ConstraintClass) State() map[string]data.GetValue {
	out := make(map[string]data.GetValue, len(c.state))
	for k, v := range c.state {
		out[k] = v
	}
	return out
}

type constraintConstructMethod struct {
	class *ConstraintClass
}

func (m *constraintConstructMethod) GetName() string            { return "__construct" }
func (m *constraintConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *constraintConstructMethod) GetIsStatic() bool          { return false }

func (m *constraintConstructMethod) GetParams() []data.GetValue {
	params := make([]data.GetValue, 0, len(m.class.spec.Params)+1)
	for i, p := range m.class.spec.Params {
		params = append(params, node.NewParameter(nil, p.Name, i, p.DefaultVal, p.Type))
	}
	params = append(params, node.NewAnnotationTargetParameter(nil, len(params)))
	return params
}

func (m *constraintConstructMethod) GetVariables() []data.Variable {
	vars := make([]data.Variable, 0, len(m.class.spec.Params)+1)
	for i, p := range m.class.spec.Params {
		vars = append(vars, node.NewVariable(nil, p.Name, i, nil))
	}
	vars = append(vars, node.NewAnnotationTargetVariable(nil, len(vars)))
	return vars
}

func (m *constraintConstructMethod) GetReturnType() data.Types { return data.NewBaseType("string") }

func (m *constraintConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	for i, p := range m.class.spec.Params {
		if v, ok := ctx.GetIndexValue(i); ok && v != nil {
			m.class.state[p.Name] = v
		}
	}
	return data.NewStringValue(m.class.spec.FullName + " annotation constructed"), nil
}
