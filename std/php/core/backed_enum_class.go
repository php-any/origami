package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// BackedEnumClass 提供一个最小实现的 \BackedEnum 基类
//
// 设计目标：
//   - 让 `instanceof \BackedEnum` 能够工作
//   - 提供 `public string $value` 属性
//   - 提供 `__construct($value)` 构造函数，将入参写入 `$this->value`
//
// 具体枚举类型（如 Status）由 enum 语法解析为继承此类的普通类来实现。
type BackedEnumClass struct {
	node.Node

	valueProp data.Property
	ctor      data.Method
}

func (c *BackedEnumClass) ensureInit() {
	if c.valueProp == nil {
		c.valueProp = node.NewProperty(
			nil,
			"value",
			"public",
			false,
			nil,
			data.NewBaseType("string"),
		)
	}
	if c.ctor == nil {
		c.ctor = &BackedEnumConstructMethod{}
	}
}

func (c *BackedEnumClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	c.ensureInit()
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *BackedEnumClass) GetFrom() data.From { return c.Node.GetFrom() }

func (c *BackedEnumClass) GetName() string { return "BackedEnum" }

func (c *BackedEnumClass) GetExtend() *string        { return nil }
func (c *BackedEnumClass) GetImplements() []string   { return nil }
func (c *BackedEnumClass) GetConstruct() data.Method { c.ensureInit(); return c.ctor }
func (c *BackedEnumClass) GetMethods() []data.Method { c.ensureInit(); return []data.Method{c.ctor} }
func (c *BackedEnumClass) GetPropertyList() []data.Property {
	c.ensureInit()
	return []data.Property{c.valueProp}
}

func (c *BackedEnumClass) GetProperty(name string) (data.Property, bool) {
	c.ensureInit()
	if name == "value" {
		return c.valueProp, true
	}
	return nil, false
}

func (c *BackedEnumClass) GetMethod(name string) (data.Method, bool) {
	c.ensureInit()
	if name == "__construct" {
		return c.ctor, true
	}
	return nil, false
}

// BackedEnumConstructMethod 实现 __construct，用于设置 $value
type BackedEnumConstructMethod struct{}

func (m *BackedEnumConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return nil, nil
	}

	// 将参数写入当前对象的 value 属性
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		if v, ok2 := val.(data.Value); ok2 {
			if acl := objCtx.SetProperty("value", v); acl != nil {
				return nil, acl
			}
		}
	}

	return nil, nil
}

func (m *BackedEnumConstructMethod) GetName() string            { return "__construct" }
func (m *BackedEnumConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *BackedEnumConstructMethod) GetIsStatic() bool          { return false }

func (m *BackedEnumConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		// 允许任意类型的枚举底层值
		node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
	}
}

func (m *BackedEnumConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
	}
}

func (m *BackedEnumConstructMethod) GetReturnType() data.Types { return nil }
