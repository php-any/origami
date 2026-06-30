package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RequestBodyClass #[RequestBody] 注解 — 声明参数来自请求体
// 用于方法级: #[RequestBody('data')] 或 #[RequestBody('user')]
type RequestBodyClass struct {
	node.Node
	source    *RequestBodySource
	construct data.Method
}

func (r *RequestBodyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := &RequestBodySource{}
	return data.NewClassValue(&RequestBodyClass{
		source:    src,
		construct: &RequestBodyConstructMethod{source: src},
	}, ctx.CreateBaseContext()), nil
}

func (r *RequestBodyClass) GetName() string    { return "Net\\Annotation\\RequestBody" }
func (r *RequestBodyClass) GetExtend() *string { return nil }
func (r *RequestBodyClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetMethod}
}
func (r *RequestBodyClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (r *RequestBodyClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (r *RequestBodyClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return r.construct, true
	}
	return nil, false
}
func (r *RequestBodyClass) GetMethods() []data.Method { return []data.Method{r.construct} }
func (r *RequestBodyClass) GetConstruct() data.Method { return r.construct }
func (r *RequestBodyClass) ParamNames() []string {
	if r.source != nil {
		return r.source.names
	}
	return nil
}

type RequestBodySource struct{ names []string }

type RequestBodyConstructMethod struct{ source *RequestBodySource }

func (m *RequestBodyConstructMethod) GetName() string            { return "__construct" }
func (m *RequestBodyConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RequestBodyConstructMethod) GetIsStatic() bool          { return false }
func (m *RequestBodyConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "names", 0, nil, nil),
	}
}
func (m *RequestBodyConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "names", 0, nil)}
}
func (m *RequestBodyConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (m *RequestBodyConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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
