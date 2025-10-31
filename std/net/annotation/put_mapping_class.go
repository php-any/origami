package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PutMappingClass PutMapping注解类
type PutMappingClass struct {
	node.Node
	source    *PutMapping
	construct data.Method
}

func (p *PutMappingClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newPutMapping()
	return data.NewClassValue(&PutMappingClass{
		source:    source,
		construct: &PutMappingConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (p *PutMappingClass) GetName() string { return "Net\\Annotation\\PutMapping" }

func (p *PutMappingClass) GetExtend() *string {
	return nil
}

func (p *PutMappingClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (p *PutMappingClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (p *PutMappingClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (p *PutMappingClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return p.construct, true
	}
	return nil, false
}

func (p *PutMappingClass) GetMethods() []data.Method { return []data.Method{p.construct} }

func (p *PutMappingClass) GetConstruct() data.Method {
	return p.construct
}

// Path 便捷访问当前实例的路径值
func (p *PutMappingClass) Path() string {
	if p.source != nil {
		return p.source.path
	}
	return ""
}

// PutMapping 映射实例
type PutMapping struct {
	path   string
	target interface{} // 被注解的反射实例
}

func newPutMapping() *PutMapping {
	return &PutMapping{}
}

// PutMappingConstructMethod 构造函数
type PutMappingConstructMethod struct {
	mapping *PutMapping
}

func (m *PutMappingConstructMethod) GetName() string {
	return "__construct"
}

func (m *PutMappingConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *PutMappingConstructMethod) GetIsStatic() bool {
	return false
}

func (m *PutMappingConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, data.NewStringValue("/"), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}

func (m *PutMappingConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *PutMappingConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *PutMappingConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	vv, _ := ctx.GetIndexValue(0)
	if v, ok := vv.(*data.StringValue); ok {
		m.mapping.path = v.AsString()
	}
	return data.NewStringValue("PutMapping annotation constructed"), nil
}

// 删除对脚本域可见的 path/process/mapping 方法，仅保留构造
