package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PathVariableClass #[PathVariable] 注解 — 声明参数来自路由路径
// 可用于方法级: #[PathVariable('id', 'postId')]
// 也可用于参数级: #[PathVariable]（需语言层支持存储参数注解，暂由方法级声明替代）
type PathVariableClass struct {
	node.Node
	source    *PathVariableSource
	construct data.Method
}

func (p *PathVariableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := &PathVariableSource{}
	return data.NewClassValue(&PathVariableClass{
		source:    src,
		construct: &PathVariableConstructMethod{source: src},
	}, ctx.CreateBaseContext()), nil
}

func (p *PathVariableClass) GetName() string    { return "Net\\Annotation\\PathVariable" }
func (p *PathVariableClass) GetExtend() *string { return nil }
func (p *PathVariableClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetMethod}
}
func (p *PathVariableClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (p *PathVariableClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (p *PathVariableClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return p.construct, true
	}
	return nil, false
}
func (p *PathVariableClass) GetMethods() []data.Method { return []data.Method{p.construct} }
func (p *PathVariableClass) GetConstruct() data.Method { return p.construct }
func (p *PathVariableClass) ParamNames() []string {
	if p.source != nil {
		return p.source.names
	}
	return nil
}

type PathVariableSource struct{ names []string }

type PathVariableConstructMethod struct{ source *PathVariableSource }

func (m *PathVariableConstructMethod) GetName() string            { return "__construct" }
func (m *PathVariableConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *PathVariableConstructMethod) GetIsStatic() bool          { return false }
func (m *PathVariableConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "names", 0, nil, nil), // 可变参数
	}
}
func (m *PathVariableConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "names", 0, nil)}
}
func (m *PathVariableConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (m *PathVariableConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 收集所有传入的参数名
	for i := 0; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok {
			break
		}
		if s, ok2 := v.(data.AsString); ok2 {
			m.source.names = append(m.source.names, s.AsString())
		}
	}
	return data.NewNullValue(), nil
}
