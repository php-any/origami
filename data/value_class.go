package data

import (
	"context"
	"errors"
	"fmt"
)

func NewClassValue(class ClassStmt, ctx Context) *ClassValue {
	var vm VM
	if ctx != nil {
		vm = ctx.GetVM()
	}
	return &ClassValue{
		ObjectValue: NewObjectValue(),
		Class:       class,
		Context:     ctx,
		vm:          vm,
	}
}

type AsClass interface {
	AsObject
}

type ClassValue struct {
	Context
	*ObjectValue
	Class ClassStmt
	// vm 钉在实例上：方法帧的 pooled Context 回收后 Context.GetVM() 会变 nil，
	// parent:: / 父类属性查找仍需 VM。SetVM(nil) 不得清掉该字段。
	vm VM
}

func (c *ClassValue) GetName() string {
	return c.Class.GetName()
}

func (c *ClassValue) GetValue(ctx Context) (GetValue, Control) {
	return c, nil
}

// AsBool：任意对象实例（含空 stdClass）在 PHP 中均为 true。
// 不可落入嵌入的 ObjectValue.AsBool（关联数组语义）。
func (c *ClassValue) AsBool() (bool, error) {
	return c != nil, nil
}

func (c *ClassValue) AsString() string {
	// PHP：对象转字符串应走 __toString；无则抛错。此处仅作调试/回退表示，
	// 禁止递归展开属性（Schema 等对象图有环，会导致栈溢出）。
	if c == nil || c.Class == nil {
		return ""
	}
	return fmt.Sprintf("Object(%s)", c.Class.GetName())
}

func (c *ClassValue) GetPropertyStmt(name string) (Property, bool) {
	if c == nil || c.Class == nil {
		return nil, false
	}
	if v, ok := c.Class.GetProperty(name); ok {
		return v, true
	}

	vm := c.GetVM()
	// 执行父级
	last := c.Class
	for vm != nil && last != nil && last.GetExtend() != nil {
		ext := last.GetExtend()
		if ext == nil || *ext == "" {
			break
		}
		next, acl := vm.GetOrLoadClass(*ext)
		if acl != nil || next == nil {
			return nil, false
		}

		property, ok := next.GetProperty(name)
		if ok {
			return property, true
		}
		last = next
	}

	return nil, false
}

func (c *ClassValue) GetProperty(name string) (Value, Control) {
	// 实例属性优先从 ObjectValue 读取（含声明属性被写入后的值）
	if c.ObjectValue != nil && c.ObjectValue.HasProperty(name) {
		return c.ObjectValue.GetProperty(name)
	}

	stmt, ok := c.GetPropertyStmt(name)
	if ok {
		// 必须用 ClassValue 自身作 Context，ClassProperty 通过 GetName 从 ObjectValue 取/初始化
		gv, acl := stmt.GetValue(c)
		if acl != nil {
			return nil, acl
		}
		if gv != nil {
			if val, ok := gv.(Value); ok {
				return val, nil
			}
		}
	}

	return NewNullValue(), nil
}

func (c *ClassValue) GetPropertyZVal(name string) (*ZVal, Control) {
	if c == nil || c.ObjectValue == nil {
		return nil, NewErrorThrow(nil, errors.New("Using $this when not in object context"))
	}
	v, ok := c.property.GetZVal(name)
	if ok && v != nil {
		return v, nil
	}
	// 声明属性带默认值时先物化（protected $componentStack = []），
	// 禁止先写成 null：array_pop($this->stack) 是引用读取，会把默认数组冲掉。
	if stmt, found := c.GetPropertyStmt(name); found {
		if _, acl := stmt.GetValue(c); acl != nil {
			return nil, acl
		}
		v, ok = c.property.GetZVal(name)
		if ok && v != nil {
			return v, nil
		}
	}
	c.SetProperty(name, NewNullValue())
	v, _ = c.property.GetZVal(name)
	return v, nil
}

func (c *ClassValue) GetMethod(name string) (Method, bool) {
	key := MethodLookupKey(name)
	if holder, ok := c.Class.(MethodLookupCacher); ok {
		if cache := holder.MethodLookupCache(); cache != nil {
			if m, found, hit, gen := cache.Lookup(key); hit {
				return m, found
			} else {
				m, found := c.lookupMethodUncached(name)
				cache.Store(key, m, found, gen)
				return m, found
			}
		}
	}
	return c.lookupMethodUncached(name)
}

func (c *ClassValue) lookupMethodUncached(name string) (Method, bool) {
	if c == nil || c.Class == nil {
		return nil, false
	}
	// 记录从类自身命中的抽象方法作为兜底：若抽象方法已被父类具体实现覆盖，
	// 则应解析到父类的具体实现（与 PHP 一致，例如 trait 声明 abstract 方法
	// 而父类提供了实现的情形）。
	var fallback Method
	if fn, ok := c.Class.GetMethod(name); ok && fn != nil {
		if _, isAbstract := fn.(AbstractMethodMarker); isAbstract {
			fallback = fn
		} else {
			return fn, true
		}
	}

	vm := c.GetVM()
	// 执行父级
	last := c.Class
	for vm != nil && last != nil && last.GetExtend() != nil {
		ext := last.GetExtend()
		next, acl := vm.GetOrLoadClass(*ext)
		if acl != nil || next == nil {
			break
		}

		fn, ok := next.GetMethod(name)
		if ok && fn != nil {
			if _, isAbstract := fn.(AbstractMethodMarker); isAbstract {
				if fallback == nil {
					fallback = fn
				}
				last = next
				continue
			}
			return fn, true
		}
		last = next
	}
	// 若父类链上仅存在抽象方法（未被具体实现覆盖），返回该抽象方法，
	// 由调用方在真正调用时抛出“抽象方法不能调用”的错误。
	if fallback != nil {
		return fallback, true
	}

	// PHP 允许在实例上调用静态方法：$obj->staticMethod()
	if gsm, ok := c.Class.(GetStaticMethod); ok {
		if m, ok := gsm.GetStaticMethod(name); ok {
			return m, true
		}
	}
	// 也在父类中查找静态方法
	last2 := c.Class
	for vm != nil && last2 != nil && last2.GetExtend() != nil {
		ext := last2.GetExtend()
		next, acl := vm.GetOrLoadClass(*ext)
		if acl != nil || next == nil {
			break
		}
		if gsm, ok := next.(GetStaticMethod); ok {
			if m, ok := gsm.GetStaticMethod(name); ok {
				return m, true
			}
		}
		last2 = next
	}

	return nil, false
}

func (c *ClassValue) GetProperties() map[string]Value {
	result := make(map[string]Value)

	// 首先获取实例属性（从 ObjectValue 继承）
	if c.ObjectValue != nil {
		instanceProps := c.ObjectValue.GetProperties()
		for name, value := range instanceProps {
			result[name] = value
		}
	}

	// 然后获取类定义的属性
	classProps := c.Class.GetPropertyList()
	for _, prop := range classProps {
		// 如果实例中没有这个属性，则使用类定义的默认值
		if _, exists := result[prop.GetName()]; !exists {
			defaultValue := prop.GetDefaultValue()
			if defaultValue != nil {
				value, _ := defaultValue.GetValue(c.Context)
				if value != nil {
					if val, ok := value.(Value); ok {
						result[prop.GetName()] = val
					}
				}
			} else {
				// 如果没有默认值，使用 null
				result[prop.GetName()] = NewNullValue()
			}
		}
	}

	// 处理继承的属性
	vm := c.GetVM()
	last := c.Class
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok {
			break
		}

		parentProps := next.GetPropertyList()
		for _, prop := range parentProps {
			// 只添加非私有属性，且实例中没有的属性
			if prop.GetModifier() != ModifierPrivate {
				if _, exists := result[prop.GetName()]; !exists {
					defaultValue := prop.GetDefaultValue()
					if defaultValue != nil {
						value, _ := defaultValue.GetValue(c.Context)
						if value != nil {
							if val, ok := value.(Value); ok {
								result[prop.GetName()] = val
							}
						}
					} else {
						result[prop.GetName()] = NewNullValue()
						c.SetProperty(prop.GetName(), result[prop.GetName()]) // 需要引用起来
					}
				}
			}
		}
		last = next
	}

	return result
}

// RangeProperties 按插入顺序遍历所有属性
// 使用此方法可保证遍历顺序与插入顺序一致，避免 Go map 遍历顺序随机的问题
func (c *ClassValue) RangeProperties(fn func(key string, value Value) bool) {
	if c == nil || c.ObjectValue == nil {
		return
	}
	c.ObjectValue.RangeProperties(fn)
}

func (c *ClassValue) CreateContext(vars []Variable) Context {
	// 符号表从对象已有的执行上下文长出来（剥掉 BoundContext），不绕回 VM.CreateContext。
	inner := unwrapBoundContext(c.Context).CreateContext(vars)
	inner.SetVM(c.vm)
	return &ClassMethodContext{
		ClassValue:  c.CloneWithContext(inner),
		StaticClass: nil,
	}
}

func (c *ClassValue) CloneWithContext(ctx Context) *ClassValue {
	return &ClassValue{
		ObjectValue: c.ObjectValue,
		Class:       c.Class,
		Context:     ctx,
		vm:          c.vm,
	}
}

// CloneSandbox 按 PHP clone 语义复制实例：数组属性按值拷贝，对象属性仍共享。
// 供 HTTP 每请求隔离 Application/Router 的 instances / currentRequest，不改共享单例服务。
func (c *ClassValue) CloneSandbox(ctx Context) *ClassValue {
	if c == nil {
		return nil
	}
	vm := c.vm
	if ctx != nil {
		if v := ctx.GetVM(); v != nil {
			vm = v
		}
	}
	obj := c.ObjectValue
	if obj != nil {
		obj = DeepCloneObjectValue(obj)
		obj.Context = ctx
	} else {
		obj = NewObjectValue()
	}
	return &ClassValue{
		ObjectValue: obj,
		Class:       c.Class,
		Context:     ctx,
		vm:          vm,
	}
}

// CloneSandboxKeys 浅拷贝对象属性表，仅对 deepKeys 中的数组/关联数组做深拷贝。
// Application 启动后 bindings/aliases 等只读；每请求只需隔离 instances 等可变槽，
// 全量 DeepClone 会在暖路径上白白拷贝数千绑定，把 /hello 压到几十 QPS。
func (c *ClassValue) CloneSandboxKeys(ctx Context, deepKeys []string) *ClassValue {
	if c == nil {
		return nil
	}
	vm := c.vm
	if ctx != nil {
		if v := ctx.GetVM(); v != nil {
			vm = v
		}
	}
	obj := c.ObjectValue
	if obj == nil {
		return &ClassValue{ObjectValue: NewObjectValue(), Class: c.Class, Context: ctx, vm: vm}
	}
	deep := make(map[string]struct{}, len(deepKeys))
	for _, k := range deepKeys {
		deep[k] = struct{}{}
	}
	clone := &ObjectValue{
		Value:                 obj.Value,
		Context:               ctx,
		property:              NewOrderedMap(),
		IndirectOverloadClass: obj.IndirectOverloadClass,
	}
	obj.property.Range(func(key string, value Value) bool {
		if _, ok := deep[key]; ok {
			clone.property.Set(key, deepCloneValue(value, 0))
		} else {
			clone.property.Set(key, value)
		}
		return true
	})
	return &ClassValue{
		ObjectValue: clone,
		Class:       c.Class,
		Context:     ctx,
		vm:          vm,
	}
}

func (c *ClassValue) withVM(vm VM) *ClassValue {
	if c == nil {
		return nil
	}
	if vm != nil {
		c.vm = vm
	}
	return c
}

func unwrapBoundContext(ctx Context) Context {
	for ctx != nil {
		if bc, ok := ctx.(*BoundContext); ok {
			ctx = bc.Context
			continue
		}
		return ctx
	}
	return ctx
}

func (c *ClassValue) SetVariableValue(variable Variable, value Value) Control {
	return c.SetProperty(variable.GetName(), value)
}

func (c *ClassValue) SetProperty(name string, value Value) Control {
	if c == nil || c.ObjectValue == nil {
		return NewErrorThrow(nil, errors.New("Using $this when not in object context"))
	}
	if set, ok := c.Class.(SetProperty); ok {
		return set.SetProperty(name, value)
	}
	if zv, ok := c.property.GetZVal(name); ok && zv != nil {
		CowAssign(zv, value)
		return nil
	}
	c.property.Set(name, CowAddRef(value))
	return nil
}

func (c *ClassValue) GetVariableValue(variable Variable) (Value, Control) {
	if c == nil || c.ObjectValue == nil {
		return nil, NewErrorThrow(nil, errors.New("Using $this when not in object context"))
	}
	return c.ObjectValue.GetVariableValue(variable)
}

func (c *ClassValue) ReturnSlot(v Value) ReturnControl {
	if c != nil && c.Context != nil {
		return c.Context.ReturnSlot(v)
	}
	return NewReturnControl(v)
}

func (c *ClassValue) GoContext() context.Context {
	return context.Background()
}

func (c *ClassValue) GetVM() VM {
	return c.vm
}

func (c *ClassValue) SetVM(vm VM) {
	if vm == nil {
		return
	}
	c.vm = vm
	c.Context.SetVM(vm)
}

// WriteOutput 把 echo 转到内层 Context 上的请求缓冲，避免落到共享 VM 再解析 goid。
func (c *ClassValue) WriteOutput(s string) {
	if c != nil && c.Context != nil {
		if sink, ok := c.Context.(OutputSink); ok {
			sink.WriteOutput(s)
			return
		}
		if vm := c.GetVM(); vm != nil {
			if sink, ok := vm.(OutputSink); ok {
				sink.WriteOutput(s)
				return
			}
		}
	}
	WriteOutput(s)
}

type ClassMethodContext struct {
	*ClassValue
	StaticClass ClassStmt // 运行时（后期）类结构，用于 static:: 后期静态绑定
	SelfClass   ClassStmt // 代码定义所在的类，用于 self:: 和 parent:: 解析（处理 trait 合并场景）
}

// WrapMethodFrame 在已有符号表外包一层方法身份，不新建符号表。
func WrapMethodFrame(inner Context, identity *ClassValue, self, static ClassStmt) *ClassMethodContext {
	if self == nil {
		self = identity.Class
	}
	if static == nil {
		static = identity.Class
	}
	inner.SetVM(identity.vm)
	return &ClassMethodContext{
		ClassValue:  identity.CloneWithContext(inner),
		SelfClass:   self,
		StaticClass: static,
	}
}

// NewStaticMethodContext 为静态方法包装已有帧。调用方传入的 inner 已是符号表。
func NewStaticMethodContext(inner Context, self ClassStmt, static ClassStmt) *ClassMethodContext {
	if static == nil {
		static = self
	}
	vm := inner.GetVM()
	inner.SetVM(vm)
	return &ClassMethodContext{
		ClassValue: &ClassValue{
			Class:   self,
			Context: inner,
			vm:      vm,
		},
		SelfClass:   self,
		StaticClass: static,
	}
}

func (c *ClassMethodContext) ReturnSlot(v Value) ReturnControl {
	return c.Context.ReturnSlot(v)
}

func (c *ClassMethodContext) CreateContext(vars []Variable) Context {
	nc := c.Context.CreateContext(vars)
	nc.SetVM(c.vm)
	return nc
}

func (c *ClassMethodContext) CreateBaseContext() Context {
	nc := c.Context.CreateBaseContext()
	nc.SetVM(c.vm)
	return nc
}

func (c *ClassMethodContext) SetVariableValue(variable Variable, value Value) Control {
	if pl, ok := variable.(PropertyLvalue); ok {
		return pl.SetValue(c, value)
	}
	return c.Context.SetVariableValue(variable, value)
}

func (c *ClassMethodContext) GetVariableValue(variable Variable) (Value, Control) {
	if _, ok := variable.(Property); ok {
		if c.ObjectValue == nil {
			return nil, NewErrorThrow(nil, errors.New("Using $this when not in object context"))
		}
		return c.ObjectValue.GetVariableValue(variable)
	}
	if pl, ok := variable.(PropertyLvalue); ok {
		gv, ctl := pl.GetValue(c)
		if ctl != nil {
			return nil, ctl
		}
		if val, ok := gv.(Value); ok {
			return val, nil
		}
		return NewNullValue(), nil
	}
	return c.Context.GetVariableValue(variable)
}

func (c *ClassMethodContext) GetIndexValue(index int) (Value, bool) {
	return c.Context.GetIndexValue(index)
}

func (c *ClassMethodContext) SetIndexZVal(index int, v *ZVal) {
	c.Context.SetIndexZVal(index, v)
}

func (c *ClassMethodContext) GetIndexZVal(index int) *ZVal {
	return c.Context.GetIndexZVal(index)
}

func (c *ClassMethodContext) BindStaticLocals(store *StaticLocals) {
	if b, ok := c.Context.(StaticLocalsBinder); ok {
		b.BindStaticLocals(store)
	}
}

func (c *ClassMethodContext) StaticLocalsStore() *StaticLocals {
	if b, ok := c.Context.(StaticLocalsBinder); ok {
		return b.StaticLocalsStore()
	}
	return nil
}

func (c *ClassMethodContext) GoContext() context.Context {
	return context.Background()
}

func (c *ClassMethodContext) ReleasePooled() {
	if r, ok := c.Context.(interface{ ReleasePooled() }); ok {
		r.ReleasePooled()
	}
}

func (c *ClassMethodContext) MarkEscaped() {
	if e, ok := c.Context.(ContextEscaper); ok {
		e.MarkEscaped()
	}
}

func (c *ClassMethodContext) IsEscaped() bool {
	if e, ok := c.Context.(ContextEscaper); ok {
		return e.IsEscaped()
	}
	return false
}

func (c *ClassMethodContext) EnterCall() int {
	return c.Context.(CallRecorder).EnterCall()
}

func (c *ClassMethodContext) LeaveCall() {
	c.Context.(CallRecorder).LeaveCall()
}

func (c *ClassMethodContext) PushCallFrame(frame CallFrame) {
	c.Context.(CallRecorder).PushCallFrame(frame)
}

func (c *ClassMethodContext) PopCallFrame() {
	c.Context.(CallRecorder).PopCallFrame()
}

func (c *ClassMethodContext) SnapshotCallStack() []CallFrame {
	return c.Context.(CallRecorder).SnapshotCallStack()
}

func (c *ClassValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalClass(c)
}

func (c *ClassValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalClass(data, c)
}

func (c *ClassValue) ToGoValue(serializer Serializer) (any, error) {
	return serializer.MarshalClass(c)
}
