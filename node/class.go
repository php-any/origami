package node

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
	"sync"
)

// ClassStatement 表示类定义语句
type ClassStatement struct {
	*Node          `pp:"-"`
	Name           string                   // 类名
	Extends        *string                  // 父类名
	Implements     []string                 // 实现的接口列表
	StaticProperty sync.Map                 // 静态属性列表
	Properties     map[string]data.Property // 属性列表
	Methods        map[string]data.Method   // 方法列表
	StaticMethods  map[string]data.Method   // 静态方法列表
	Annotations    []*data.ClassValue       // 类注解列表

	// 构造函数
	Construct data.Method
}

// GetValue 获取类定义语句的值
func (c *ClassStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	object := data.NewClassValue(c, ctx)

	for _, property := range c.Properties {
		def := property.GetDefaultValue()
		if def == nil {
			continue
		}
		v, ctl := def.GetValue(object)
		if ctl != nil {
			return nil, ctl
		}
		object.SetProperty(property.GetName(), v.(data.Value))
	}
	if c.Extends != nil {
		vm := object.GetVM()
		ext, ok := vm.GetClass(*c.Extends)
		if ok {
			_, ctl := ext.GetValue(object)
			if ctl != nil {
				return nil, ctl
			}
		}
	}

	return object, nil
}

func (c *ClassStatement) GetConstruct() data.Method {
	return c.Construct
}

// NewClassStatement 创建一个新的类定义语句
func NewClassStatement(from data.From, name string, extends string, implements []string, properties map[string]data.Property, methods map[string]data.Method) *ClassStatement {
	class := &ClassStatement{
		Node:       NewNode(from),
		Name:       name,
		Extends:    &extends,
		Implements: implements,
		Properties: properties,
		Methods:    methods,
	}
	if extends == "" {
		class.Extends = nil
	}
	if construct, ok := class.GetMethod(token.ConstructName); ok {
		class.Construct = construct
	}
	return class
}

// GetName 返回类名
func (c *ClassStatement) GetName() string {
	return c.Name
}

func (c *ClassStatement) GetExtend() *string {
	return c.Extends
}

// GetImplements 返回实现的接口列表
func (c *ClassStatement) GetImplements() []string {
	return c.Implements
}

func (c *ClassStatement) AddAnnotations(a *data.ClassValue) {
	if c.Annotations == nil {
		c.Annotations = []*data.ClassValue{}
	}
	c.Annotations = append(c.Annotations, a)
}

func (c *ClassStatement) GetProperties() map[string]data.Property {
	return c.Properties
}

func (c *ClassStatement) GetProperty(name string) (data.Property, bool) {
	if f, ok := c.Properties[name]; ok {
		return f, true
	}
	return nil, false
}

func (c *ClassStatement) GetMethod(name string) (data.Method, bool) {
	if f, ok := c.Methods[name]; ok {
		return f, true
	}
	return nil, false
}

func (c *ClassStatement) GetMethods() []data.Method {
	var methods []data.Method
	for _, f := range c.Methods {
		methods = append(methods, f)
	}
	return methods
}

func (c *ClassStatement) GetStaticProperty(name string) (data.Value, bool) {
	if f, ok := c.StaticProperty.Load(name); ok {
		return f.(data.Value), true
	}
	return nil, false
}

func (c *ClassStatement) GetStaticMethod(name string) (data.Method, bool) {
	if f, ok := c.StaticMethods[name]; ok {
		return f, true
	}
	return nil, false
}

type ClassProperty struct {
	*Node        `pp:"-"`
	Name         string             // 属性名
	Modifier     data.Modifier      // 访问修饰符
	IsStatic     bool               // 是否是静态属性
	DefaultValue data.GetValue      // 默认值
	Annotations  []*data.ClassValue // 属性注解列表
}

func (p *ClassProperty) GetIndex() int {
	panic("属性使用哈希实现")
}

func (p *ClassProperty) GetType() data.Types {
	//TODO implement me
	panic("implement me")
}

func (p *ClassProperty) SetValue(ctx data.Context, value data.Value) data.Control {
	p.DefaultValue = value
	return nil
}

// NewProperty 创建一个新的属性
func NewProperty(from data.From, name string, modifier string, isStatic bool, defaultValue data.GetValue) *ClassProperty {
	if name[0:1] == "$" {
		name = name[1:]
	}
	return &ClassProperty{
		Node:         NewNode(from),
		Name:         name,
		Modifier:     data.NewModifier(modifier),
		IsStatic:     isStatic,
		DefaultValue: defaultValue,
	}
}

func (p *ClassProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	v, acl := ctx.GetVariableValue(p)
	if v != nil {
		return v, acl
	} else {
		v = data.NewNullValue()
	}
	if p.DefaultValue != nil {
		v, acl := p.DefaultValue.GetValue(ctx)
		if v != nil {
			ctx.SetVariableValue(p, v.(data.Value))
			return v, acl
		}
		return v, acl
	}
	return v, acl
}

// GetName 返回属性名
func (p *ClassProperty) GetName() string {
	return p.Name
}

// GetModifier 返回访问修饰符
func (p *ClassProperty) GetModifier() data.Modifier {
	return p.Modifier
}

// GetDefaultValue 返回默认值
func (p *ClassProperty) GetDefaultValue() data.GetValue {
	return p.DefaultValue
}

func (p *ClassProperty) GetIsStatic() bool {
	return p.IsStatic
}

func (p *ClassProperty) AddAnnotations(a *data.ClassValue) {
	if p.Annotations == nil {
		p.Annotations = []*data.ClassValue{}
	}
	p.Annotations = append(p.Annotations, a)
}

type ClassMethod struct {
	*Node       `pp:"-"`
	Name        string          // 方法名
	Modifier    data.Modifier   // 访问修饰符
	IsStatic    bool            // 是否是静态方法
	Params      []data.GetValue // 参数列表
	Body        []Statement     // 方法体
	vars        []data.Variable
	Annotations []*data.ClassValue // 方法注解列表
	Ret         data.Types         // 返回类型
}

func (m *ClassMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	//TODO implement me
	panic("implement me")
}

// NewMethod 创建一个新的方法
func NewMethod(from data.From, name string, modifier string, isStatic bool, params []data.GetValue, body []Statement, vars []data.Variable, ret data.Types) data.Method {
	return &ClassMethod{
		Node:     NewNode(from),
		Name:     name,
		Modifier: data.NewModifier(modifier),
		IsStatic: isStatic,
		Params:   params,
		Body:     body,
		vars:     vars,
		Ret:      ret,
	}
}

func (m *ClassMethod) AddAnnotations(a *data.ClassValue) {
	if m.Annotations == nil {
		m.Annotations = []*data.ClassValue{}
	}
	m.Annotations = append(m.Annotations, a)
}

func (m *ClassMethod) GetIsStatic() bool {
	return false
}

// GetName 返回方法名
func (m *ClassMethod) GetName() string {
	return m.Name
}

// GetModifier 返回访问修饰符
func (m *ClassMethod) GetModifier() data.Modifier {
	return m.Modifier
}

// GetParams 返回参数列表
func (m *ClassMethod) GetParams() []data.GetValue {
	return m.Params
}

func (m *ClassMethod) GetVariables() []data.Variable {
	return m.vars
}

// GetReturnType 返回方法返回类型
func (m *ClassMethod) GetReturnType() data.Types {
	return m.Ret
}

func (m *ClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var ctl data.Control
	for _, statement := range m.Body {
		v, ctl = statement.GetValue(ctx)
		if ctl != nil {
			switch rv := ctl.(type) {
			case data.ReturnControl:
				return rv.ReturnValue(), nil
			}
			return nil, ctl
		}
	}

	return v, nil
}

// 检查 source 是否实现了(继承了) target 类或接口
func checkClassIs(ctx data.Context, source data.ClassStmt, target string) bool {
	if source.GetName() == target {
		return true
	} else {
		if source.GetImplements() != nil {
			for _, impl := range source.GetImplements() {
				if impl == target {
					return true
				}
				// 检查接口继承
				if vm := ctx.GetVM(); vm != nil {
					if interfaceStmt, ok := vm.GetInterface(impl); ok {
						if checkInterfaceIs(ctx, interfaceStmt, target) {
							return true
						}
					}
				}
			}
		}

		if source.GetExtend() != nil {
			vm := ctx.GetVM()
			// 执行父级
			last := source
			for last.GetExtend() != nil || last.GetImplements() != nil {
				if last.GetImplements() != nil {
					for _, impl := range last.GetImplements() {
						if impl == target {
							return true
						}
						// 检查接口继承
						if interfaceStmt, ok := vm.GetInterface(impl); ok {
							if checkInterfaceIs(ctx, interfaceStmt, target) {
								return true
							}
						}
					}
				}
				if last.GetExtend() != nil {
					if *last.GetExtend() == target {
						return true
					}
					next, ok := vm.GetClass(*(last.GetExtend()))
					if ok && checkClassIs(ctx, next, target) {
						return true
					} else if ok {
						last = next
					}
				}
				return false
			}
		}
	}

	return false
}

// 检查接口是否继承了目标接口
func checkInterfaceIs(ctx data.Context, source data.InterfaceStmt, target string) bool {
	if source.GetName() == target {
		return true
	}

	if source.GetExtend() != nil {
		vm := ctx.GetVM()
		if interfaceStmt, ok := vm.GetInterface(*source.GetExtend()); ok {
			return checkInterfaceIs(ctx, interfaceStmt, target)
		}
	}

	return false
}

type AddAnnotations interface {
	AddAnnotations(a *data.ClassValue)
}
