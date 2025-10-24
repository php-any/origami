package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewIdClass 创建 Id 注解类
func NewIdClass() *IdClass {
	return &IdClass{}
}

// IdClass Id注解类 - 特性注解
type IdClass struct {
	node.Node
	construct data.Method
}

func (i *IdClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建新的注解实例
	instance := &IdClass{}

	// 创建构造函数，传入注解实例引用
	instance.construct = &IdConstructMethod{idClass: instance}

	return data.NewClassValue(instance, ctx.CreateBaseContext()), nil
}

func (i *IdClass) GetName() string {
	return "Database\\Annotation\\Id"
}

func (i *IdClass) GetExtend() *string {
	return nil
}

func (i *IdClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (i *IdClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (i *IdClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (i *IdClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return i.construct, true
	}
	return nil, false
}

func (i *IdClass) GetMethods() []data.Method {
	return []data.Method{
		i.construct,
	}
}

func (i *IdClass) GetConstruct() data.Method {
	return i.construct
}

// IdConstructMethod 构造函数 - 特性注解只接收注解参数
type IdConstructMethod struct {
	idClass *IdClass // 引用注解实例
}

func (m *IdConstructMethod) GetName() string {
	return "__construct"
}

func (m *IdConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *IdConstructMethod) GetIsStatic() bool {
	return false
}

func (m *IdConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *IdConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *IdConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *IdConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 主键注解构造函数：不需要参数，只负责初始化
	// 构造函数只负责初始化，返回 nil 表示成功
	return nil, nil
}
