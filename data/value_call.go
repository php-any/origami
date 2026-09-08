package data

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
	if c == nil || c.Value == nil {
		return ""
	}
	if gn, ok := c.Value.(GetName); ok {
		return gn.GetName()
	}
	return "Closure"
}

func (c *FuncValue) AsBool() (bool, error) {
	return true, nil
}

// GetMethod 使 FuncValue 支持 Closure 实例方法 (bindTo, call 等)
func (c *FuncValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "bindto", "bindTo":
		return &funcBindToMethod{closure: c}, true
	case "call":
		return &funcCallMethod{closure: c}, true
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

	boundThis := classValueFromBindObject(newThis)
	scopeClass := scopeClassFromBindArg(newScope)
	if scopeClass == "" && boundThis != nil {
		scopeClass = boundThis.Class.GetName()
	}

	if scopeClass != "" || boundThis != nil {
		return NewBoundFuncValue(m.closure.Value, scopeClass, boundThis), nil
	}
	return m.closure, nil
}

func classValueFromBindObject(v Value) *ClassValue {
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *ClassValue:
		return o
	case *ThisValue:
		return o.ClassValue
	}
	return nil
}

func scopeClassFromBindArg(v Value) string {
	if v == nil {
		return ""
	}
	switch o := v.(type) {
	case *StringValue:
		return o.Value
	case *ClassValue:
		return o.Class.GetName()
	case *ThisValue:
		if o.ClassValue != nil {
			return o.ClassValue.Class.GetName()
		}
	}
	return ""
}

// funcCallMethod 实现 Closure::call($newThis, ...$args)
type funcCallMethod struct {
	closure *FuncValue
}

func (m *funcCallMethod) GetName() string       { return "call" }
func (m *funcCallMethod) GetModifier() Modifier { return ModifierPublic }
func (m *funcCallMethod) GetIsStatic() bool     { return false }
func (m *funcCallMethod) GetReturnType() Types  { return nil }
func (m *funcCallMethod) GetParams() []GetValue {
	return []GetValue{
		NewParameter("newThis", 0),
		NewParameters("args", 1),
	}
}
func (m *funcCallMethod) GetVariables() []Variable {
	return []Variable{
		NewVariable("newThis", 0, nil),
		NewVariable("args", 1, nil),
	}
}
func (m *funcCallMethod) Call(ctx Context) (GetValue, Control) {
	newThis, _ := ctx.GetIndexValue(0)
	boundThis := classValueFromBindObject(newThis)
	scopeClass := ""
	if boundThis != nil {
		scopeClass = boundThis.Class.GetName()
	}

	args := closureCallArgs(ctx)
	vars := make([]Variable, len(args))
	for i := range args {
		vars[i] = NewVariable("", i, nil)
	}
	callCtx := ctx.CreateContext(vars)
	for i, arg := range args {
		callCtx.SetIndexZVal(i, NewZVal(arg))
	}

	return NewBoundFuncValue(m.closure.Value, scopeClass, boundThis).Call(callCtx)
}

func closureCallArgs(ctx Context) []Value {
	if flat := ctx.GetFlatCallArgs(); len(flat) > 1 {
		return flat[1:]
	}
	if argsVal, ok := ctx.GetIndexValue(1); ok && argsVal != nil {
		if arr, isArr := argsVal.(*ArrayValue); isArr {
			out := make([]Value, 0, len(arr.List))
			for _, zv := range arr.List {
				if zv != nil {
					out = append(out, zv.Value)
				}
			}
			return out
		}
		return []Value{argsVal}
	}
	var out []Value
	for i := 1; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok {
			break
		}
		out = append(out, v)
	}
	return out
}

// BoundFuncValue 表示通过 Closure::bind()/bindTo() 绑定了 $this 或作用域的闭包。
type BoundFuncValue struct {
	FuncValue
	ScopeClass  string
	BoundObject *ClassValue // bindTo($newThis, ...) 绑定的 $this 对象
}

func NewBoundFuncValue(v FuncStmt, scopeClass string, boundThis *ClassValue) *BoundFuncValue {
	return &BoundFuncValue{
		FuncValue:   FuncValue{Value: v},
		ScopeClass:  scopeClass,
		BoundObject: boundThis,
	}
}

func (b *BoundFuncValue) Call(ctx Context) (GetValue, Control) {
	return b.Value.Call(&BoundContext{
		Context:    ctx,
		ScopeClass: b.ScopeClass,
		BoundThis:  b.BoundObject,
	})
}

// BoundContext 携带 Closure::bind()/bindTo() 绑定信息。
type BoundContext struct {
	Context
	ScopeClass string
	BoundThis  *ClassValue
}

func (bc *BoundContext) CreateContext(vars []Variable) Context {
	return &BoundContext{
		Context:    bc.Context.CreateContext(vars),
		ScopeClass: bc.ScopeClass,
		BoundThis:  bc.BoundThis,
	}
}

func (bc *BoundContext) CreateBaseContext() Context {
	return &BoundContext{
		Context:    bc.Context.CreateBaseContext(),
		ScopeClass: bc.ScopeClass,
		BoundThis:  bc.BoundThis,
	}
}

// FindBoundContext 沿上下文链查找 Closure::bind/bindTo 注入的 BoundContext。
func FindBoundContext(ctx Context) *BoundContext {
	for c := ctx; c != nil; c = parentContext(c) {
		if bc, ok := c.(*BoundContext); ok {
			return bc
		}
	}
	return nil
}

func parentContext(ctx Context) Context {
	switch c := ctx.(type) {
	case *BoundContext:
		return c.Context
	case *ClassMethodContext:
		if c.ClassValue != nil {
			return c.ClassValue.Context
		}
	case *ClassValue:
		return c.Context
	}
	return nil
}
