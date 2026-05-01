package data

import "fmt"

func NewFuncValue(v FuncStmt) *FuncValue {
	return &FuncValue{
		Value: v,
	}
}

type FuncValue struct {
	Value FuncStmt
}

func (c *FuncValue) GetValue(ctx Context) (GetValue, Control) {
	return c, nil
}

func (c *FuncValue) Call(ctx Context) (GetValue, Control) {
	return c.Value.Call(ctx)
}

func (c *FuncValue) AsString() string {
	return fmt.Sprintf("%v", c.Value)
}

func (c *FuncValue) AsBool() (bool, error) {
	return true, nil
}

// GetMethod 使 FuncValue 支持 Closure 实例方法 (bindTo, call 等)
func (c *FuncValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "bindto", "bindTo":
		return &funcBindToMethod{closure: c}, true
	}
	return nil, false
}

// funcBindToMethod 实现 Closure::bindTo($newThis, $newScope)
type funcBindToMethod struct {
	closure *FuncValue
}

func (m *funcBindToMethod) GetName() string       { return "bindTo" }
func (m *funcBindToMethod) GetModifier() Modifier { return ModifierPublic }
func (m *funcBindToMethod) GetIsStatic() bool     { return false }
func (m *funcBindToMethod) GetReturnType() Types  { return nil }
func (m *funcBindToMethod) GetParams() []GetValue {
	return []GetValue{
		NewParameter("newThis", 0),
		NewParameter("newScope", 1),
	}
}
func (m *funcBindToMethod) GetVariables() []Variable {
	return []Variable{
		NewVariable("newThis", 0, nil),
		NewVariable("newScope", 1, nil),
	}
}
func (m *funcBindToMethod) Call(ctx Context) (GetValue, Control) {
	newThis, _ := ctx.GetIndexValue(0)
	newScope, _ := ctx.GetIndexValue(1)

	// 提取 scope 类名
	var scopeClass string
	if sv, ok := newScope.(*StringValue); ok {
		scopeClass = sv.Value
	} else if cv, ok := newScope.(*ClassValue); ok {
		scopeClass = cv.Class.GetName()
	} else if _, ok := newScope.(*NullValue); ok {
		// scope = null → 无作用域绑定
	} else if newThis != nil {
		// 如果 newScope 是对象，用它的类名
		if cv, ok := newThis.(*ClassValue); ok {
			scopeClass = cv.Class.GetName()
		}
	}

	if scopeClass != "" {
		return NewBoundFuncValue(m.closure.Value, scopeClass), nil
	}
	return m.closure, nil
}

// BoundFuncValue 表示通过 Closure::bind() 绑定了作用域的闭包
// 在此闭包内可以访问绑定类的私有成员
type BoundFuncValue struct {
	FuncValue
	ScopeClass string // 绑定的类名，闭包可以访问该类的私有成员
}

func NewBoundFuncValue(v FuncStmt, scopeClass string) *BoundFuncValue {
	return &BoundFuncValue{
		FuncValue:  FuncValue{Value: v},
		ScopeClass: scopeClass,
	}
}

func (b *BoundFuncValue) Call(ctx Context) (GetValue, Control) {
	return b.Value.Call(&BoundContext{Context: ctx, ScopeClass: b.ScopeClass})
}

// BoundContext 携带 Closure::bind() 绑定作用域的上下文
type BoundContext struct {
	Context
	ScopeClass string
}

func (bc *BoundContext) CreateContext(vars []Variable) Context {
	return &BoundContext{
		Context:    bc.Context.CreateContext(vars),
		ScopeClass: bc.ScopeClass,
	}
}

func (bc *BoundContext) CreateBaseContext() Context {
	return &BoundContext{
		Context:    bc.Context.CreateBaseContext(),
		ScopeClass: bc.ScopeClass,
	}
}
