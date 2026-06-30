package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RequestParamClass #[RequestParam] 注解 — 声明参数来自 URL 查询参数
// 用于方法级: #[RequestParam('keyword', 'page', 'size')]
type RequestParamClass struct {
	node.Node
	source    *RequestParamSource
	construct data.Method
}

func (r *RequestParamClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := &RequestParamSource{}
	return data.NewClassValue(&RequestParamClass{
		source:    src,
		construct: &RequestParamConstructMethod{source: src},
	}, ctx.CreateBaseContext()), nil
}

func (r *RequestParamClass) GetName() string    { return "Net\\Annotation\\RequestParam" }
func (r *RequestParamClass) GetExtend() *string { return nil }
func (r *RequestParamClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetMethod}
}
func (r *RequestParamClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (r *RequestParamClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (r *RequestParamClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return r.construct, true
	}
	return nil, false
}
func (r *RequestParamClass) GetMethods() []data.Method { return []data.Method{r.construct} }
func (r *RequestParamClass) GetConstruct() data.Method { return r.construct }
func (r *RequestParamClass) ParamNames() []string {
	if r.source != nil {
		return r.source.names
	}
	return nil
}

type RequestParamSource struct{ names []string }

type RequestParamConstructMethod struct{ source *RequestParamSource }

func (m *RequestParamConstructMethod) GetName() string            { return "__construct" }
func (m *RequestParamConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RequestParamConstructMethod) GetIsStatic() bool          { return false }
func (m *RequestParamConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "names", 0, nil, nil),
	}
}
func (m *RequestParamConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "names", 0, nil)}
}
func (m *RequestParamConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (m *RequestParamConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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
