package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DeleteMappingClass DeleteMapping注解类
type DeleteMappingClass struct {
	node.Node
	source    *DeleteMapping
	construct data.Method
}

func (d *DeleteMappingClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newDeleteMapping()
	return data.NewClassValue(&DeleteMappingClass{
		source:    source,
		construct: &DeleteMappingConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (d *DeleteMappingClass) GetName() string { return "Net\\Annotation\\DeleteMapping" }

func (d *DeleteMappingClass) GetExtend() *string {
	return nil
}

func (d *DeleteMappingClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (d *DeleteMappingClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (d *DeleteMappingClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (d *DeleteMappingClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return d.construct, true
	}
	return nil, false
}

func (d *DeleteMappingClass) GetMethods() []data.Method { return []data.Method{d.construct} }

func (d *DeleteMappingClass) GetConstruct() data.Method {
	return d.construct
}

// Path 便捷访问当前实例的路径值
func (d *DeleteMappingClass) Path() string {
	if d.source != nil {
		return d.source.path
	}
	return ""
}

// DeleteMapping 映射实例
type DeleteMapping struct {
	path   string
	target interface{} // 被注解的反射实例
}

func newDeleteMapping() *DeleteMapping {
	return &DeleteMapping{}
}

// DeleteMappingConstructMethod 构造函数
type DeleteMappingConstructMethod struct {
	mapping *DeleteMapping
}

func (m *DeleteMappingConstructMethod) GetName() string {
	return "__construct"
}

func (m *DeleteMappingConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *DeleteMappingConstructMethod) GetIsStatic() bool {
	return false
}

func (m *DeleteMappingConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, data.NewStringValue("/"), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}

func (m *DeleteMappingConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *DeleteMappingConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *DeleteMappingConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	vv, _ := ctx.GetIndexValue(0)
	if v, ok := vv.(*data.StringValue); ok {
		m.mapping.path = v.AsString()
	}
	return data.NewStringValue("DeleteMapping annotation constructed"), nil
}

// 删除对脚本域可见的 path/process/mapping 方法，仅保留构造
