package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetMappingClass GetMapping注解类
type GetMappingClass struct {
	node.Node
	source    *GetMapping
	construct data.Method
}

func (g *GetMappingClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newGetMapping()
	return data.NewClassValue(&GetMappingClass{
		source:    source,
		construct: &GetMappingConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (g *GetMappingClass) GetName() string { return "Net\\Annotation\\GetMapping" }

func (g *GetMappingClass) GetExtend() *string {
	return nil
}

func (g *GetMappingClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (g *GetMappingClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (g *GetMappingClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (g *GetMappingClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return g.construct, true
	}
	return nil, false
}

func (g *GetMappingClass) GetMethods() []data.Method { return []data.Method{g.construct} }

func (g *GetMappingClass) GetConstruct() data.Method {
	return g.construct
}

// Path 便捷访问当前实例的路径值
func (g *GetMappingClass) Path() string {
	if g.source != nil {
		return g.source.path
	}
	return ""
}

// GetMapping 映射实例
type GetMapping struct {
	path   string
	target interface{} // 被注解的反射实例
}

func newGetMapping() *GetMapping {
	return &GetMapping{}
}

// GetMappingConstructMethod 构造函数
type GetMappingConstructMethod struct {
	mapping *GetMapping
}

func (m *GetMappingConstructMethod) GetName() string {
	return "__construct"
}

func (m *GetMappingConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GetMappingConstructMethod) GetIsStatic() bool {
	return false
}

func (m *GetMappingConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, data.NewStringValue("/"), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}

func (m *GetMappingConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *GetMappingConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *GetMappingConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	vv, _ := ctx.GetIndexValue(0)
	if v, ok := vv.(*data.StringValue); ok {
		m.mapping.path = v.AsString()
	}
	return data.NewStringValue("GetMapping annotation constructed"), nil
}

// GetMappingPathMethod 暴露 path 给控制器注解读取
// 删除对脚本域可见的 path/process/mapping 方法，仅保留构造

// GetMappingProcessMethod 处理注解的方法
// 保留最小接口，仅通过构造接收参数

//（已移除）

// GetMappingMappingMethod 注册映射的方法
//（已移除）
