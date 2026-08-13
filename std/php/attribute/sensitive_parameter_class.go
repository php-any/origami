package attribute

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SensitiveParameterClass PHP 8.2+ 原生 SensitiveParameter 类
// 用于标记敏感参数（如密码、API密钥等）
type SensitiveParameterClass struct {
	node.Node
	construct data.Method
}

func NewSensitiveParameterClass() *SensitiveParameterClass {
	return &SensitiveParameterClass{
		construct: &SensitiveParameterConstructMethod{},
	}
}

func (c *SensitiveParameterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *SensitiveParameterClass) GetName() string {
	return "SensitiveParameter"
}

func (c *SensitiveParameterClass) GetExtend() *string {
	return nil
}

func (c *SensitiveParameterClass) GetImplements() []string {
	return []string{}
}

func (c *SensitiveParameterClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *SensitiveParameterClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *SensitiveParameterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *SensitiveParameterClass) GetMethods() []data.Method {
	return []data.Method{
		c.construct,
	}
}

func (c *SensitiveParameterClass) GetConstruct() data.Method {
	return c.construct
}

// SensitiveParameterConstructMethod SensitiveParameter 构造函数
// SensitiveParameter::__construct()
type SensitiveParameterConstructMethod struct{}

func (m *SensitiveParameterConstructMethod) GetName() string {
	return "__construct"
}

func (m *SensitiveParameterConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *SensitiveParameterConstructMethod) GetIsStatic() bool {
	return false
}

func (m *SensitiveParameterConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *SensitiveParameterConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *SensitiveParameterConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *SensitiveParameterConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// SensitiveParameter 构造函数不需要特殊处理，只是标记参数为敏感
	return data.NewNullValue(), nil
}
