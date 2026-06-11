package container

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewComponentClass() data.ClassStmt {
	return &ComponentClass{construct: &ComponentConstructMethod{}}
}

type ComponentClass struct {
	node.Node
	construct data.Method
}

func (c *ComponentClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&ComponentClass{construct: &ComponentConstructMethod{}}, ctx.CreateBaseContext()), nil
}
func (c *ComponentClass) GetName() string    { return "Container\\Component" }
func (c *ComponentClass) GetExtend() *string { return nil }
func (c *ComponentClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (c *ComponentClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *ComponentClass) GetPropertyList() []data.Property           { return nil }
func (c *ComponentClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *ComponentClass) GetMethods() []data.Method { return []data.Method{c.construct} }
func (c *ComponentClass) GetConstruct() data.Method { return c.construct }

type ComponentConstructMethod struct{}

func (m *ComponentConstructMethod) GetName() string            { return "__construct" }
func (m *ComponentConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ComponentConstructMethod) GetIsStatic() bool          { return false }
func (m *ComponentConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}
func (m *ComponentConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}
func (m *ComponentConstructMethod) GetReturnType() data.Types { return nil }
func (m *ComponentConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, registerClassAnnotation(ctx, LifetimeTransient)
}

func NewSingletonAnnotationClass() data.ClassStmt {
	return &SingletonAnnotationClass{construct: &SingletonAnnotationConstructMethod{}}
}

type SingletonAnnotationClass struct {
	node.Node
	construct data.Method
}

func (c *SingletonAnnotationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&SingletonAnnotationClass{construct: &SingletonAnnotationConstructMethod{}}, ctx.CreateBaseContext()), nil
}
func (c *SingletonAnnotationClass) GetName() string    { return "Container\\Singleton" }
func (c *SingletonAnnotationClass) GetExtend() *string { return nil }
func (c *SingletonAnnotationClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (c *SingletonAnnotationClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *SingletonAnnotationClass) GetPropertyList() []data.Property           { return nil }
func (c *SingletonAnnotationClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *SingletonAnnotationClass) GetMethods() []data.Method { return []data.Method{c.construct} }
func (c *SingletonAnnotationClass) GetConstruct() data.Method { return c.construct }

type SingletonAnnotationConstructMethod struct{}

func (m *SingletonAnnotationConstructMethod) GetName() string            { return "__construct" }
func (m *SingletonAnnotationConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SingletonAnnotationConstructMethod) GetIsStatic() bool          { return false }
func (m *SingletonAnnotationConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}
func (m *SingletonAnnotationConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}
func (m *SingletonAnnotationConstructMethod) GetReturnType() data.Types { return nil }
func (m *SingletonAnnotationConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, registerClassAnnotation(ctx, LifetimeSingleton)
}

func NewScopedAnnotationClass() data.ClassStmt {
	return &ScopedAnnotationClass{construct: &ScopedAnnotationConstructMethod{}}
}

type ScopedAnnotationClass struct {
	node.Node
	construct data.Method
}

func (c *ScopedAnnotationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&ScopedAnnotationClass{construct: &ScopedAnnotationConstructMethod{}}, ctx.CreateBaseContext()), nil
}
func (c *ScopedAnnotationClass) GetName() string    { return "Container\\Scoped" }
func (c *ScopedAnnotationClass) GetExtend() *string { return nil }
func (c *ScopedAnnotationClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (c *ScopedAnnotationClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *ScopedAnnotationClass) GetPropertyList() []data.Property           { return nil }
func (c *ScopedAnnotationClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *ScopedAnnotationClass) GetMethods() []data.Method { return []data.Method{c.construct} }
func (c *ScopedAnnotationClass) GetConstruct() data.Method { return c.construct }

type ScopedAnnotationConstructMethod struct{}

func (m *ScopedAnnotationConstructMethod) GetName() string            { return "__construct" }
func (m *ScopedAnnotationConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ScopedAnnotationConstructMethod) GetIsStatic() bool          { return false }
func (m *ScopedAnnotationConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}
func (m *ScopedAnnotationConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}
func (m *ScopedAnnotationConstructMethod) GetReturnType() data.Types { return nil }
func (m *ScopedAnnotationConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, registerClassAnnotation(ctx, LifetimeScoped)
}
