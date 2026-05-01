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
