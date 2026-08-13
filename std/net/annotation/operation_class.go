package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	netdata "github.com/php-any/origami/std/net/data"
)

// OperationClass @Operation 注解类（参考 Spring / OpenAPI @Operation）
//
// 用于描述接口操作的元数据，可与 @GetMapping 等映射注解一起标注在控制器方法上。
//
// 用法:
//
//	#[Operation(summary: "获取商品列表", description: "返回所有商品", tags: ["products"])]
//	#[GetMapping(path: "/products")]
//	public function listProducts(): void { }
type OperationClass struct {
	node.Node
	source    *Operation
	construct data.Method
}

func (o *OperationClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newOperation()
	return data.NewClassValue(&OperationClass{
		source:    source,
		construct: &OperationConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (o *OperationClass) GetName() string { return "Net\\Annotation\\Operation" }

func (o *OperationClass) GetExtend() *string {
	return nil
}

func (o *OperationClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetMethod}
}

func (o *OperationClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (o *OperationClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (o *OperationClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return o.construct, true
	}
	return nil, false
}

func (o *OperationClass) GetMethods() []data.Method { return []data.Method{o.construct} }

func (o *OperationClass) GetConstruct() data.Method {
	return o.construct
}

func (o *OperationClass) Summary() string     { return o.source.summary }
func (o *OperationClass) Description() string { return o.source.description }
func (o *OperationClass) Tags() []string      { return append([]string(nil), o.source.tags...) }
func (o *OperationClass) OperationID() string { return o.source.operationId }
func (o *OperationClass) Deprecated() bool    { return o.source.deprecated }
func (o *OperationClass) Hidden() bool        { return o.source.hidden }

// Info 转为路由注册使用的元数据结构。
func (o *OperationClass) Info() netdata.OperationInfo {
	if o.source == nil {
		return netdata.OperationInfo{}
	}
	return netdata.OperationInfo{
		Summary:     o.source.summary,
		Description: o.source.description,
		Tags:        append([]string(nil), o.source.tags...),
		OperationID: o.source.operationId,
		Deprecated:  o.source.deprecated,
		Hidden:      o.source.hidden,
	}
}

// FindOperationAnnotation 从方法注解列表中查找 @Operation。
func FindOperationAnnotation(annotations []*data.ClassValue) *OperationClass {
	for _, ann := range annotations {
		if oc, ok := ann.Class.(*OperationClass); ok {
			return oc
		}
	}
	return nil
}

// Operation 映射实例
type Operation struct {
	summary     string
	description string
	tags        []string
	operationId string
	deprecated  bool
	hidden      bool
	target      interface{}
}

func newOperation() *Operation {
	return &Operation{}
}

// OperationConstructMethod 构造函数
type OperationConstructMethod struct {
	operation *Operation
}

func (m *OperationConstructMethod) GetName() string {
	return "__construct"
}

func (m *OperationConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *OperationConstructMethod) GetIsStatic() bool {
	return false
}

func (m *OperationConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "summary", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "description", 1, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "tags", 2, data.NewNullValue(), nil),
		node.NewParameter(nil, "operationId", 3, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "deprecated", 4, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "hidden", 5, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewAnnotationTargetParameter(nil, 6),
	}
}

func (m *OperationConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "summary", 0, nil),
		node.NewVariable(nil, "description", 1, nil),
		node.NewVariable(nil, "tags", 2, nil),
		node.NewVariable(nil, "operationId", 3, nil),
		node.NewVariable(nil, "deprecated", 4, nil),
		node.NewVariable(nil, "hidden", 5, nil),
		node.NewAnnotationTargetVariable(nil, 6),
	}
}

func (m *OperationConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *OperationConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if v, ok := ctx.GetIndexValue(0); ok {
		m.operation.summary = v.AsString()
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		m.operation.description = v.AsString()
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		m.operation.tags = parseStringSlice(v)
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		m.operation.operationId = v.AsString()
	}
	if v, ok := ctx.GetIndexValue(4); ok {
		m.operation.deprecated = parseBoolValue(v, false)
	}
	if v, ok := ctx.GetIndexValue(5); ok {
		m.operation.hidden = parseBoolValue(v, false)
	}
	return data.NewStringValue("Operation annotation constructed"), nil
}

func parseStringSlice(v data.GetValue) []string {
	if v == nil {
		return nil
	}
	if sv, ok := v.(*data.StringValue); ok {
		if s := sv.AsString(); s != "" {
			return []string{s}
		}
		return nil
	}
	arr, ok := v.(*data.ArrayValue)
	if !ok {
		return nil
	}
	tags := make([]string, 0, len(arr.List))
	for _, item := range arr.ToValueList() {
		if s := item.AsString(); s != "" {
			tags = append(tags, s)
		}
	}
	return tags
}

func parseBoolValue(v data.GetValue, defaultVal bool) bool {
	if v == nil {
		return defaultVal
	}
	if bv, ok := v.(*data.BoolValue); ok {
		b, err := bv.AsBool()
		if err == nil {
			return b
		}
	}
	return defaultVal
}
