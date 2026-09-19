package node

import (
	"fmt"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
	"github.com/php-any/origami/utils"
)

// ClassStatement 表示类定义语句
type ClassStatement struct {
	*Node           `pp:"-"`
	Name            string                   // 类名
	Extends         *string                  // 父类名
	Implements      []string                 // 实现的接口列表
	StaticProperty  sync.Map                 // 静态属性存储（运行时值）
	PropertiesIndex []string                 // 属性列表
	Properties      map[string]data.Property // 属性列表
	Methods         map[string]data.Method   // 方法列表
	StaticMethods   map[string]data.Method   // 静态方法列表
	methodsLower    map[string]data.Method   // 方法名小写索引（PHP 方法名不区分大小写）
	staticLower     map[string]data.Method
	indexedMethodN  int
	methodIdxMu     sync.RWMutex
	lookupCache     *data.MethodLookupCache
	Annotations     []*data.ClassValue // 类注解列表

	// 构造函数
	Construct data.Method

	// IsAbstract 为 true 表示 abstract class（允许未实现接口/抽象方法）
	IsAbstract bool

	// AnnotationsApplied 标记类注解是否已在执行期应用过，避免普通 require 重复执行时重复注册 @Route/@Command
	AnnotationsApplied bool

	// Traits 保存类直接使用的 trait 名（含命名空间全限定名），供 class_uses() 返回。
	Traits []string

	// DeferredTraits 保存解析期未能加载（依赖运行期 require/autoload 的 trait）的 trait 名。
	// 这类 trait 无法在解析期合并，需在运行期（ClassRegisterStmt）加载后合并进类。
	DeferredTraits       []string
	DeferredTraitAliases []data.TraitAlias
	// DeferredTraitsMerged 标记运行期是否已合并过延迟 trait，避免重复合并。
	DeferredTraitsMerged bool

	// StaticProperties 保存静态属性/常量声明（解析期收集，值延迟到首次访问时求值）。
	// 原生 PHP 的类常量是惰性求值的，允许前向引用（如 const A = [self::B]; const B = 1;），
	// 因此不能在解析期按声明顺序立即求值。
	// 注意：求值不加锁（sync.Mutex 不可重入，前向引用求值会递归调用 GetStaticProperty，
	// 加锁会死锁）。PHP 常量表达式无副作用，并发下重复求值结果一致，幂等无害。
	StaticProperties      map[string]data.Property
	StaticPropertiesIndex []string

	staticPropertyCtx data.Context // 惰性求值上下文（解析期缓存的 ClassValue）
}

// GetValue 获取类定义语句的值
func (c *ClassStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if !c.IsAbstract {
		if acl := ValidateConcreteClassAbstractMethods(ctx.GetVM(), c); acl != nil {
			return nil, acl
		}
	}
	object := data.NewClassValue(c, ctx)

	for _, property := range c.Properties {
		if property.GetIsStatic() {
			continue
		}
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
		ext, acl := vm.GetOrLoadClass(*c.Extends)
		if acl != nil {
			return nil, acl
		}
		// 初始化父类的属性（包括所有继承链上的父类）
		last := ext
		for {
			for _, property := range last.GetPropertyList() {
				// 如果子类已经设置了该属性，跳过
				if _, exists := c.Properties[property.GetName()]; exists {
					continue
				}
				// 初始化有默认值的属性
				def := property.GetDefaultValue()
				if def != nil {
					v, ctl := def.GetValue(object)
					if ctl != nil {
						return nil, ctl
					}
					object.SetProperty(property.GetName(), v.(data.Value))
				}
			}
			// 继续处理父类的父类
			if last.GetExtend() == nil {
				break
			}
			next, acl := vm.GetOrLoadClass(*last.GetExtend())
			if acl != nil {
				return nil, acl
			}
			last = next
		}
		_, acl = ext.GetValue(object)
		if acl != nil {
			return nil, acl
		}
	}

	return object, nil
}

// MergeDeferredTraits 在运行期合并解析期未能加载的 trait（依赖 require/autoload）。
// 幂等：已合并过则不重复。别名在全部 trait 合并后应用。
func (c *ClassStatement) MergeDeferredTraits(vm data.VM) data.Control {
	if len(c.DeferredTraits) == 0 || c.DeferredTraitsMerged {
		return nil
	}
	var stillDeferred []string
	for _, traitName := range c.DeferredTraits {
		trait, acl := vm.GetOrLoadClass(traitName)
		if acl != nil {
			stillDeferred = append(stillDeferred, traitName)
			continue
		}
		if trait == nil {
			return utils.NewThrowf("trait %s 不存在", traitName)
		}
		if acl := mergeTraitStmt(vm, c, trait); acl != nil {
			return acl
		}
	}
	if len(stillDeferred) > 0 {
		c.DeferredTraits = stillDeferred
		return nil
	}
	c.DeferredTraits = nil
	c.DeferredTraitsMerged = true
	// 应用别名：必须从 trait 取原方法，避免类覆盖同名方法后别名指错（如 __call as macroCall）
	for _, alias := range c.DeferredTraitAliases {
		if method, ok := findDeferredTraitMethod(vm, c.Traits, alias, false); ok {
			if _, exists := c.Methods[alias.Alias]; !exists {
				c.Methods[alias.Alias] = method
			}
		}
		if method, ok := findDeferredTraitMethod(vm, c.Traits, alias, true); ok {
			if _, exists := c.StaticMethods[alias.Alias]; !exists {
				c.StaticMethods[alias.Alias] = method
			}
		}
	}
	c.DeferredTraitAliases = nil
	c.invalidateMethodLookups()
	return nil
}

func findDeferredTraitMethod(vm data.VM, traitNames []string, alias data.TraitAlias, static bool) (data.Method, bool) {
	candidates := traitNames
	if alias.Trait != "" {
		candidates = []string{alias.Trait}
		for _, name := range traitNames {
			if name == alias.Trait || strings.HasSuffix(name, "\\"+alias.Trait) {
				candidates = []string{name}
				break
			}
		}
	}
	for _, traitName := range candidates {
		trait, acl := vm.GetOrLoadClass(traitName)
		if acl != nil || trait == nil {
			continue
		}
		if static {
			if cs, ok := trait.(*ClassStatement); ok {
				if method, ok := cs.StaticMethods[alias.Method]; ok {
					return method, true
				}
			}
			if gsm, ok := trait.(data.GetStaticMethod); ok {
				if method, ok := gsm.GetStaticMethod(alias.Method); ok {
					return method, true
				}
			}
			continue
		}
		for _, method := range trait.GetMethods() {
			if method.GetName() == alias.Method {
				return method, true
			}
		}
	}
	return nil, false
}

// CopyTraitInstanceMethods 把 trait 的实例方法并入 dst。
// 必须按 Methods 的 map 键复制：trait 内 `use T { foo as bar }` 会让同一方法占两个键，
// 若只按 method.GetName() 合并，别名键会丢（Filament Notification::getBaseIconColor）。
func CopyTraitInstanceMethods(dst map[string]data.Method, trait data.ClassStmt) {
	if dst == nil || trait == nil {
		return
	}
	if cs, ok := trait.(*ClassStatement); ok {
		for name, method := range cs.Methods {
			if _, exists := dst[name]; !exists {
				dst[name] = method
			}
		}
		return
	}
	for _, method := range trait.GetMethods() {
		if method == nil {
			continue
		}
		name := method.GetName()
		if _, exists := dst[name]; !exists {
			dst[name] = method
		}
	}
}

// mergeTraitStmt 将一个已加载的 trait 合并进类（实例/静态方法、属性、静态属性）。
func mergeTraitStmt(vm data.VM, class *ClassStatement, trait data.ClassStmt) data.Control {
	CopyTraitInstanceMethods(class.Methods, trait)
	if cs, ok := trait.(*ClassStatement); ok {
		for methodName, method := range cs.StaticMethods {
			if _, exists := class.StaticMethods[methodName]; !exists {
				class.StaticMethods[methodName] = method
			}
		}
	}
	traitProperties := trait.GetPropertyList()
	for _, property := range traitProperties {
		propertyName := property.GetName()
		if _, exists := class.Properties[propertyName]; !exists {
			if property.GetIsStatic() {
				defaultValue := property.GetDefaultValue()
				if defaultValue != nil {
					classVal := data.NewClassValue(class, vm.CreateContext([]data.Variable{}))
					v, acl := defaultValue.GetValue(classVal)
					if acl != nil {
						return acl
					}
					class.StaticProperty.Store(propertyName, v)
				} else {
					class.StaticProperty.Store(propertyName, data.NewNullValue())
				}
			} else {
				class.Properties[propertyName] = property
				class.PropertiesIndex = append(class.PropertiesIndex, propertyName)
			}
		}
	}
	if cs, ok := trait.(*ClassStatement); ok {
		// 已求值槽位（trait 静态属性曾通过惰性求值初始化过）
		cs.StaticProperty.Range(func(key, value any) bool {
			if _, exists := class.StaticProperty.Load(key); !exists {
				class.StaticProperty.Store(key, value.(data.Value))
			}
			return true
		})
		// 惰性声明：trait 静态属性惰性化后可能尚未求值（StaticProperty 为空），
		// 需把声明复制到类，使 static::$x 访问能按惰性求值在类上下文初始化。
		for name, prop := range cs.StaticProperties {
			if _, exists := class.StaticProperties[name]; exists {
				continue
			}
			if _, has := class.StaticProperty.Load(name); has {
				continue
			}
			class.StaticProperties[name] = prop
			class.StaticPropertiesIndex = append(class.StaticPropertiesIndex, name)
		}
	}
	class.invalidateMethodLookups()
	return nil
}

func (c *ClassStatement) GetConstruct() data.Method {
	return c.Construct
}

// NewClassStatement 创建一个新的类定义语句
func NewClassStatement(from data.From, name string, extends string, implements []string, properties []data.Property, methods map[string]data.Method) *ClassStatement {
	propertiesIndex := make([]string, len(properties))
	propertiesMap := make(map[string]data.Property, len(properties))
	for i, property := range properties {
		propertiesIndex[i] = property.GetName()
		propertiesMap[property.GetName()] = property
	}

	class := &ClassStatement{
		Node:            NewNode(from),
		Name:            name,
		Extends:         &extends,
		Implements:      implements,
		PropertiesIndex: propertiesIndex,
		Properties:      propertiesMap,
		Methods:         methods,
		lookupCache:     data.NewMethodLookupCache(),
	}
	if extends == "" {
		class.Extends = nil
	}
	if construct, ok := methods[token.ConstructName]; ok && construct != nil {
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

func (c *ClassStatement) GetPropertyList() []data.Property {
	properties := make([]data.Property, len(c.Properties))
	for i, prop := range c.PropertiesIndex {
		properties[i] = c.Properties[prop]
	}
	return properties
}

func (c *ClassStatement) GetProperty(name string) (data.Property, bool) {
	if f, ok := c.Properties[name]; ok {
		return f, true
	}
	return nil, false
}

func (c *ClassStatement) MethodLookupCache() *data.MethodLookupCache {
	if c.lookupCache == nil {
		c.lookupCache = data.NewMethodLookupCache()
	}
	return c.lookupCache
}

func (c *ClassStatement) invalidateMethodLookups() {
	c.methodIdxMu.Lock()
	c.methodsLower = nil
	c.staticLower = nil
	c.indexedMethodN = -1
	c.methodIdxMu.Unlock()
	if c.lookupCache != nil {
		c.lookupCache.Invalidate()
	}
}

func (c *ClassStatement) rebuildMethodIndexLocked() {
	n := len(c.Methods) + len(c.StaticMethods)
	if c.methodsLower != nil && c.indexedMethodN == n {
		return
	}
	c.methodsLower = make(map[string]data.Method, len(c.Methods)+1)
	for key, f := range c.Methods {
		if f != nil {
			c.methodsLower[strings.ToLower(key)] = f
		}
	}
	if c.Construct != nil {
		lk := strings.ToLower(token.ConstructName)
		if _, ok := c.methodsLower[lk]; !ok {
			c.methodsLower[lk] = c.Construct
		}
	}
	c.staticLower = make(map[string]data.Method, len(c.StaticMethods))
	for key, f := range c.StaticMethods {
		if f != nil {
			c.staticLower[strings.ToLower(key)] = f
		}
	}
	c.indexedMethodN = n
}

func (c *ClassStatement) lookupMethodLower(name string) (data.Method, bool) {
	n := len(c.Methods) + len(c.StaticMethods)
	c.methodIdxMu.RLock()
	if c.methodsLower != nil && c.indexedMethodN == n {
		f, ok := c.methodsLower[strings.ToLower(name)]
		c.methodIdxMu.RUnlock()
		return f, ok && f != nil
	}
	c.methodIdxMu.RUnlock()

	c.methodIdxMu.Lock()
	c.rebuildMethodIndexLocked()
	f, ok := c.methodsLower[strings.ToLower(name)]
	c.methodIdxMu.Unlock()
	return f, ok && f != nil
}

func (c *ClassStatement) lookupStaticLower(name string) (data.Method, bool) {
	n := len(c.Methods) + len(c.StaticMethods)
	c.methodIdxMu.RLock()
	if c.staticLower != nil && c.indexedMethodN == n {
		f, ok := c.staticLower[strings.ToLower(name)]
		c.methodIdxMu.RUnlock()
		return f, ok && f != nil
	}
	c.methodIdxMu.RUnlock()

	c.methodIdxMu.Lock()
	c.rebuildMethodIndexLocked()
	f, ok := c.staticLower[strings.ToLower(name)]
	c.methodIdxMu.Unlock()
	return f, ok && f != nil
}

func (c *ClassStatement) GetMethod(name string) (data.Method, bool) {
	if f, ok := c.Methods[name]; ok && f != nil {
		return f, true
	}
	if name != "" {
		if f, ok := c.lookupMethodLower(name); ok {
			return f, true
		}
	}
	if name == token.ConstructName && c.Construct != nil {
		return c.Construct, true
	}
	return nil, false
}

func (c *ClassStatement) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(c.Methods)+len(c.StaticMethods)+1)
	for _, f := range c.Methods {
		methods = append(methods, f)
	}
	for _, f := range c.StaticMethods {
		methods = append(methods, f)
	}
	// 如果构造函数存在且不在方法列表中，添加它
	if c.Construct != nil {
		hasConstruct := false
		for _, method := range methods {
			if method.GetName() == token.ConstructName {
				hasConstruct = true
				break
			}
		}
		if !hasConstruct {
			methods = append(methods, c.Construct)
		}
	}
	return methods
}

func (c *ClassStatement) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := data.LoadRequestStatic(c.GetName(), name); ok {
		return v, true
	}
	if f, ok := c.StaticProperty.Load(name); ok {
		return data.CowRequestStatic(c.GetName(), name, f.(data.Value)), true
	}
	// 惰性初始化：声明列表中存在但尚未求值的静态属性/常量。
	// 前向引用（const A = [self::B]; const B = 1;）在求值 A 时会递归调用
	// GetStaticProperty("B")，因此这里不能加锁（sync.Mutex 不可重入会死锁）。
	// PHP 常量表达式无副作用，重复求值结果一致。
	if prop, ok := c.StaticProperties[name]; ok {
		v, acl := c.initStaticProperty(prop)
		if acl == nil && v != nil {
			return data.CowRequestStatic(c.GetName(), name, v), true
		}
	}
	return nil, false
}

// SetStaticPropertyContext 设置惰性求值上下文（解析期缓存 ClassValue，使 self::/parent::/static:: 可用）。
func (c *ClassStatement) SetStaticPropertyContext(ctx data.Context) {
	c.staticPropertyCtx = ctx
}

// initStaticProperty 求值静态属性/常量的默认值并缓存到 StaticProperty。
func (c *ClassStatement) initStaticProperty(prop data.Property) (data.Value, data.Control) {
	def := prop.GetDefaultValue()
	if def == nil {
		c.StaticProperty.Store(prop.GetName(), data.NewNullValue())
		return data.NewNullValue(), nil
	}
	ctx := c.staticPropertyCtx
	if ctx == nil {
		return nil, nil
	}
	v, acl := def.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	val, ok := v.(data.Value)
	if !ok || val == nil {
		val = data.NewNullValue()
	}
	c.StaticProperty.Store(prop.GetName(), val)
	return val, nil
}

func (c *ClassStatement) GetStaticMethod(name string) (data.Method, bool) {
	if f, ok := c.StaticMethods[name]; ok && f != nil {
		return f, true
	}
	if name == "" {
		return nil, false
	}
	return c.lookupStaticLower(name)
}

type ClassProperty struct {
	*Node        `pp:"-"`
	Name         string             // 属性名
	Modifier     data.Modifier      // 访问修饰符
	IsStatic     bool               // 是否是静态属性
	IsReadonly   bool               // 是否是只读属性
	IsPromoted   bool               // 是否是构造函数参数属性提升
	DefaultValue data.GetValue      // 默认值
	Annotations  []*data.ClassValue // 属性注解列表
	Type         data.Types         // 属性类型
}

func (p *ClassProperty) GetIndex() int {
	panic("属性使用哈希实现")
}

func (p *ClassProperty) GetZVal(object data.GetPropertyZVal) (*data.ZVal, data.Control) {
	return object.GetPropertyZVal(p.Name)
}

func (p *ClassProperty) GetType() data.Types {
	return p.Type
}

func (p *ClassProperty) SetType(t data.Types) {
	p.Type = t
}

func (p *ClassProperty) SetValue(ctx data.Context, value data.Value) data.Control {
	p.DefaultValue = value
	return nil
}

// NewProperty 创建一个新的属性
func NewProperty(from data.From, name string, modifier string, isStatic bool, defaultValue data.GetValue, tys ...data.Types) *ClassProperty {
	return NewPropertyWithReadonly(from, name, modifier, isStatic, false, defaultValue, tys...)
}

// NewPropertyWithReadonly 创建一个新的属性（支持 readonly）
func NewPropertyWithReadonly(from data.From, name string, modifier string, isStatic bool, isReadonly bool, defaultValue data.GetValue, tys ...data.Types) *ClassProperty {
	return NewPropertyWithPromoted(from, name, modifier, isStatic, isReadonly, false, defaultValue, tys...)
}

// NewPropertyWithPromoted 创建一个新的属性（支持 readonly 和 promoted）
func NewPropertyWithPromoted(from data.From, name string, modifier string, isStatic bool, isReadonly bool, isPromoted bool, defaultValue data.GetValue, tys ...data.Types) *ClassProperty {
	if name[0:1] == "$" {
		name = name[1:]
	}
	var ty data.Types
	if len(tys) > 0 {
		ty = tys[0]
	}

	return &ClassProperty{
		Node:         NewNode(from),
		Name:         name,
		Modifier:     data.NewModifier(modifier),
		IsStatic:     isStatic,
		IsReadonly:   isReadonly,
		IsPromoted:   isPromoted,
		DefaultValue: defaultValue,
		Type:         ty,
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
			if c, ok := ctx.(data.SetProperty); ok {
				c.SetProperty(p.Name, v.(data.Value))
			} else {
				ctx.SetVariableValue(p, v.(data.Value))
			}
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
	Body        []data.GetValue // 方法体
	vars        []data.Variable
	Annotations []*data.ClassValue // 方法注解列表
	Ret         data.Types         // 返回类型
	IsGenerator bool               // 是否是生成器方法（含 yield）
}

func (m *ClassMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	//TODO implement me
	panic("implement me")
}

// NewMethod 创建一个新的方法
func NewMethod(from data.From, name string, modifier string, isStatic bool, params []data.GetValue, body []data.GetValue, vars []data.Variable, ret data.Types) data.Method {
	return &ClassMethod{
		Node:        NewNode(from),
		Name:        name,
		Modifier:    data.NewModifier(modifier),
		IsStatic:    isStatic,
		Params:      params,
		Body:        body,
		vars:        vars,
		Ret:         ret,
		IsGenerator: containsYield(body),
	}
}

func (m *ClassMethod) AddAnnotations(a *data.ClassValue) {
	if m.Annotations == nil {
		m.Annotations = []*data.ClassValue{}
	}
	m.Annotations = append(m.Annotations, a)
}

func (m *ClassMethod) GetIsStatic() bool {
	return m.IsStatic
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
	// 不再无条件 BindStaticLocals。static 局部变量由 StaticVarStatement 惰性绑定。

	// PHP 语义：如果方法是 generator（含 yield），调用时立即返回 Generator 对象，不执行方法体
	if m.IsGenerator {
		markContextEscaped(ctx)
		generator := NewFuncYieldStackState(ctx, m, m.Body, 0, nil, nil)
		generatorClass := NewGeneratorClass(generator)
		return generatorClass.GetValue(ctx)
	}

	// 调用深度限制，防止无限递归导致栈溢出。
	// Filament/Livewire + debug DOMDocument 校验的合法嵌套可超过 500。
	frame := data.CallFrame{Function: m.Name}
	if m.IsStatic {
		frame.Type = "::"
	} else {
		frame.Type = "->"
	}
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		frame.Class = cmc.Class.GetName()
	}
	if from := m.GetFrom(); from != nil {
		frame.File = from.GetSource()
		line, _ := from.GetStartPosition()
		frame.Line = line + 1
	}
	depth := phpCallEnter(ctx, frame)
	if depth > 2500 {
		phpCallLeave(ctx)
		className := ""
		if cmc, ok := ctx.(*data.ClassMethodContext); ok {
			className = cmc.Class.GetName() + "::"
		}
		return nil, data.NewErrorThrow(m.GetFrom(), fmt.Errorf("方法 %s%s 调用深度超过限制(%d)", className, m.Name, depth))
	}
	defer phpCallLeave(ctx)

	var ctl data.Control
	for bodyIndex := 0; bodyIndex < len(m.Body); bodyIndex++ {
		statement := m.Body[bodyIndex]
		_, ctl = statement.GetValue(ctx)
		if ctl != nil {
			switch rv := ctl.(type) {
			case data.ExitControl:
				return nil, ctl
			case data.ReturnControl:
				ret := rv.ReturnValue()
				if m.Ret == nil {
					return ret, nil // 不判断类型
				}
				if m.Ret.Is(ret) {
					return ret, nil
				}
				// 允许 null 返回（PHP 兼容：方法可能隐式返回 null）
				if ret == nil {
					return data.NewNullValue(), nil
				}
				if _, isNull := ret.(*data.NullValue); isNull {
					return data.NewNullValue(), nil
				}
				// 声明返回 string 时，允许返回带 __toString 的对象并自动转为字符串（与 PHP 一致）
				if m.Ret != nil && m.Ret.String() == "string" {
					if obj, ok := ret.(*data.ClassValue); ok {
						if toStr, ok := obj.GetMethod("__toString"); ok && toStr != nil {
							fnCtx := obj.CreateContext(toStr.GetVariables())
							fnCtx.SetCallArgs([]data.GetValue{})
							val, ctl := toStr.Call(fnCtx)
							if ctl == nil && val != nil {
								return val, nil
							}
						}
					}
				}
				return nil, data.NewErrorThrow(m.GetFrom(), fmt.Errorf("方法(%s)返回值类型错误; 期望 %s, 实际 %T", m.Name, m.Ret.String(), ret))
			case data.GotoControl:
				offset, acl := resolveGotoBodyIndex(m.from, m.Body, rv)
				if acl != nil {
					return nil, acl
				}
				bodyIndex = offset - 1
				continue
			case LabelControl:
				continue
			case data.YieldControl:
				// Generator 方法：将执行状态保存为生成器
				generator := rv.CreateStackState(ctx, m, m.Body, bodyIndex)
				generatorClass := NewGeneratorClass(generator)
				return generatorClass.GetValue(ctx)
			case data.YieldValueControl:
				// Generator 方法：将执行状态保存为生成器
				generator := NewFuncYieldStackState(ctx, m, m.Body, bodyIndex+1, rv.GetYieldKey(), rv.GetYieldValue())
				generatorClass := NewGeneratorClass(generator)
				return generatorClass.GetValue(ctx)
			case data.AddStack:
				if tv, ok := rv.(*data.ThrowValue); ok && tv.PHPUncaughtError {
					return nil, ctl
				}
				if c, ok := statement.(GetFrom); ok {
					rv.AddStackWithInfo(c.GetFrom(), "body", TryGetCallClassName(statement))
				}
				if c, ok := ctx.(data.GetName); ok {
					rv.AddStackWithInfo(m.from, c.GetName(), m.GetName())
				} else {
					rv.AddStackWithInfo(m.from, "", m.GetName())
				}
			}
			return nil, ctl
		}
	}

	persistStaticLocals(ctx, m.vars)
	// PHP：方法没有 return 时返回 null（Livewire ViewContext::extractFromEnvironment 依赖此语义）。
	return data.NewNullValue(), nil
}

// 检查 source 是否实现了(继承了) target 类或接口
func checkClassIs(ctx data.Context, source data.ClassStmt, target string) (bool, data.Control) {
	if source.GetName() == target {
		return true, nil
	} else {
		if source.GetImplements() != nil {
			for _, impl := range source.GetImplements() {
				if impl == target {
					return true, nil
				}
				// 检查接口继承
				if vm := ctx.GetVM(); vm != nil {
					if interfaceStmt, ok := vm.GetInterface(impl); ok {
						if checkInterfaceIs(ctx, interfaceStmt, target) {
							return true, nil
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
							return true, nil
						}
						// 检查接口继承
						if interfaceStmt, ok := vm.GetInterface(impl); ok {
							if checkInterfaceIs(ctx, interfaceStmt, target) {
								return true, nil
							}
						}
					}
				}
				if last.GetExtend() != nil {
					if *last.GetExtend() == target {
						return true, nil
					}
					next, acl := vm.GetOrLoadClass(*(last.GetExtend()))
					if acl != nil {
						return false, acl
					}
					ok, acl := checkClassIs(ctx, next, target)
					if acl != nil {
						return false, acl
					}
					if ok {
						return true, nil
					} else {
						last = next
					}
				}
				return false, nil
			}
		}
	}

	return false, nil
}

// 检查接口是否继承了目标接口
func checkInterfaceIs(ctx data.Context, source data.InterfaceStmt, target string) bool {
	if source.GetName() == target {
		return true
	}

	vm := ctx.GetVM()
	for _, parentName := range source.GetExtends() {
		if interfaceStmt, ok := vm.GetInterface(parentName); ok {
			if checkInterfaceIs(ctx, interfaceStmt, target) {
				return true
			}
		}
	}

	return false
}

type AddAnnotations interface {
	AddAnnotations(a *data.ClassValue)
}
