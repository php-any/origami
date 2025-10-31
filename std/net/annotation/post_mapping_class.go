package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PostMapping 注解（标记 POST 路由）
type PostMappingClass struct {
	node.Node
	source    *PostMapping
	construct data.Method
}

func (p *PostMappingClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newPostMapping()
	return data.NewClassValue(&PostMappingClass{
		source:    source,
		construct: &PostMappingConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (p *PostMappingClass) GetName() string                            { return "Net\\Annotation\\PostMapping" }
func (p *PostMappingClass) GetExtend() *string                         { return nil }
func (p *PostMappingClass) GetImplements() []string                    { return []string{node.TypeFeature} }
func (p *PostMappingClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (p *PostMappingClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (p *PostMappingClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return p.construct, true
	}
	return nil, false
}
func (p *PostMappingClass) GetMethods() []data.Method { return []data.Method{p.construct} }
func (p *PostMappingClass) GetConstruct() data.Method { return p.construct }

// Path 便捷访问当前实例的路径值
func (p *PostMappingClass) Path() string {
	if p.source != nil {
		return p.source.path
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

// 删除对脚本域可见的 path/process/mapping 方法，仅保留构造
