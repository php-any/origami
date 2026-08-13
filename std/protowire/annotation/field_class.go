package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// NewFieldClass 创建 Field 注解类
func NewFieldClass() *FieldClass {
	return &FieldClass{}
}

// FieldClass 是 @Field(number: int, type: int, encoding?: string) 注解，
// 用于声明属性对应的 protobuf 字段编号、wire type 和可选的编码方式。
//
// encoding 可选值:
//   - "" (默认): 整数编码
//   - "float": float32 (配合 PROTOWIRE_FIXED32)
//   - "double": float64 (配合 PROTOWIRE_FIXED64)
//   - "zigzag": sint32/sint64 (配合 PROTOWIRE_VARINT)
//   - "packed": packed repeated (配合 PROTOWIRE_LENGTH_DELIMITED)
//   - "message": 嵌套消息序列化 (配合 PROTOWIRE_LENGTH_DELIMITED)
type FieldClass struct {
	node.Node
	construct data.Method
	number    int
	wireType  int
	encoding  string
}

func (c *FieldClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	instance := &FieldClass{}
	instance.construct = &FieldConstructMethod{fieldClass: instance}
	return data.NewClassValue(instance, ctx.CreateBaseContext()), nil
}

func (c *FieldClass) GetName() string {
	return "Protowire\\Annotation\\Field"
}

func (c *FieldClass) GetExtend() *string {
	return nil
}

func (c *FieldClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetProperty}
}

func (c *FieldClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "number":
		return node.NewProperty(nil, "number", "public", false, data.NewIntValue(c.number)), true
	case "type":
		return node.NewProperty(nil, "type", "public", false, data.NewIntValue(c.wireType)), true
	case "encoding":
		return node.NewProperty(nil, "encoding", "public", false, data.NewStringValue(c.encoding)), true
	}
	return nil, false
}

func (c *FieldClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "number", "public", false, data.NewIntValue(c.number)),
		node.NewProperty(nil, "type", "public", false, data.NewIntValue(c.wireType)),
		node.NewProperty(nil, "encoding", "public", false, data.NewStringValue(c.encoding)),
	}
}

func (c *FieldClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *FieldClass) GetMethods() []data.Method {
	return []data.Method{c.construct}
}

func (c *FieldClass) GetConstruct() data.Method {
	return c.construct
}

type FieldConstructMethod struct {
	fieldClass *FieldClass
}

func (m *FieldConstructMethod) GetName() string            { return "__construct" }
func (m *FieldConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FieldConstructMethod) GetIsStatic() bool          { return false }
func (m *FieldConstructMethod) GetReturnType() data.Types  { return nil }
func (m *FieldConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "number", 0, data.NewNullValue(), data.NewBaseType("int")),
		node.NewParameter(nil, "type", 1, data.NewNullValue(), data.NewBaseType("int")),
		node.NewParameter(nil, "encoding", 2, data.NewStringValue(""), data.NewBaseType("string")),
	}
}
func (m *FieldConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "number", 0, nil),
		node.NewVariable(nil, "type", 1, nil),
		node.NewVariable(nil, "encoding", 2, nil),
	}
}

func (m *FieldConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("@Field 缺少 number 参数"))
	}
	number := 0
	if ai, ok := a0.(data.AsInt); ok {
		n, err := ai.AsInt()
		if err != nil {
			return nil, utils.NewThrow(errors.New("@Field number 参数必须为整数"))
		}
		number = n
	}

	wireType := 0
	if a1, ok := ctx.GetIndexValue(1); ok {
		if ai, ok := a1.(data.AsInt); ok {
			n, err := ai.AsInt()
			if err != nil {
				return nil, utils.NewThrow(errors.New("@Field type 参数必须为整数"))
			}
			wireType = n
		}
	}

	encoding := ""
	if a2, ok := ctx.GetIndexValue(2); ok {
		if sv, ok := a2.(*data.StringValue); ok {
			encoding = sv.AsString()
		}
	}

	if m.fieldClass != nil {
		m.fieldClass.number = number
		m.fieldClass.wireType = wireType
		m.fieldClass.encoding = encoding
	}
	return nil, nil
}
