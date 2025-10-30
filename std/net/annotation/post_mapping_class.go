package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PostMapping 注解（标记 POST 路由）
type PostMappingClass struct {
	node.Node
	process   data.Method
	mapping   data.Method
	construct data.Method
	pathMeth  data.Method
}

func (p *PostMappingClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newPostMapping()
	return data.NewClassValue(&PostMappingClass{
		process:   &PostMappingProcessMethod{source},
		mapping:   &PostMappingMappingMethod{source},
		construct: &PostMappingConstructMethod{source},
		pathMeth:  &PostMappingPathMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (p *PostMappingClass) GetName() string                            { return "Net\\Annotation\\PostMapping" }
func (p *PostMappingClass) GetExtend() *string                         { return nil }
func (p *PostMappingClass) GetImplements() []string                    { return []string{node.TypeFeature} }
func (p *PostMappingClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (p *PostMappingClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (p *PostMappingClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return p.process, true
	case "mapping":
		return p.mapping, true
	case "__construct":
		return p.construct, true
	case "path":
		return p.pathMeth, true
	}
	return nil, false
}
func (p *PostMappingClass) GetMethods() []data.Method {
	return []data.Method{p.process, p.mapping, p.construct, p.pathMeth}
}
func (p *PostMappingClass) GetConstruct() data.Method { return p.construct }

// Path 便捷访问当前实例的路径值
func (p *PostMappingClass) Path() string {
	if pm, ok := p.process.(*PostMappingProcessMethod); ok && pm.mapping != nil {
		return pm.mapping.path
	}
	return ""
}

type PostMapping struct {
	path   string
	target interface{}
}

func newPostMapping() *PostMapping { return &PostMapping{} }

type PostMappingConstructMethod struct{ mapping *PostMapping }

func (m *PostMappingConstructMethod) GetName() string            { return "__construct" }
func (m *PostMappingConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *PostMappingConstructMethod) GetIsStatic() bool          { return false }
func (m *PostMappingConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, data.NewStringValue("/"), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}
func (m *PostMappingConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}
func (m *PostMappingConstructMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *PostMappingConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	vv, _ := ctx.GetIndexValue(0)
	if v, ok := vv.(*data.StringValue); ok {
		m.mapping.path = v.AsString()
	}
	return data.NewStringValue("PostMapping annotation constructed"), nil
}

// PostMappingPathMethod 暴露 path 给控制器注解读取
type PostMappingPathMethod struct{ mapping *PostMapping }

func (m *PostMappingPathMethod) GetName() string            { return "path" }
func (m *PostMappingPathMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *PostMappingPathMethod) GetIsStatic() bool          { return false }
func (m *PostMappingPathMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *PostMappingPathMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *PostMappingPathMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *PostMappingPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.mapping.path), nil
}

type PostMappingProcessMethod struct{ mapping *PostMapping }

func (m *PostMappingProcessMethod) GetName() string               { return "process" }
func (m *PostMappingProcessMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PostMappingProcessMethod) GetIsStatic() bool             { return false }
func (m *PostMappingProcessMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *PostMappingProcessMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *PostMappingProcessMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *PostMappingProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue("PostMapping processed"), nil
}

type PostMappingMappingMethod struct{ mapping *PostMapping }

func (m *PostMappingMappingMethod) GetName() string               { return "mapping" }
func (m *PostMappingMappingMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *PostMappingMappingMethod) GetIsStatic() bool             { return false }
func (m *PostMappingMappingMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *PostMappingMappingMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *PostMappingMappingMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *PostMappingMappingMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue("POST mapping registered"), nil
}
