package kit

import (
	"fmt"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MacroStore 按类名存放宏（对齐 Laravel Macroable trait）。
type MacroStore struct {
	mu     sync.RWMutex
	macros map[string]map[string]data.Value // class -> lower(name) -> callable
}

var globalMacros = &MacroStore{macros: map[string]map[string]data.Value{}}

// RegisterMacroable 给 methods 表挂上 macro / hasMacro / flushMacros / __call / __callStatic。
func RegisterMacroable(methods map[string]data.Method, className string) {
	methods["macro"] = StaticMethod("macro", []string{"name", "macro"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroRegister(ctx, className)
	})
	methods["hasmacro"] = StaticMethod("hasMacro", []string{"name"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		name := ""
		if v := Arg(ctx, 0); v != nil {
			name = v.AsString()
		}
		return data.NewBoolValue(globalMacros.Has(className, name)), nil
	})
	methods["flushmacros"] = StaticMethod("flushMacros", nil, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		globalMacros.Flush(className)
		return data.NewNullValue(), nil
	})
	methods["mixin"] = StaticMethod("mixin", []string{"mixin", "replace"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroMixin(ctx, className)
	})
	methods["__callstatic"] = StaticMethod("__callStatic", []string{"method", "parameters"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroCall(ctx, className, true)
	})
	methods["__call"] = InstanceMethod("__call", []string{"method", "parameters"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroCall(ctx, className, false)
	})
}

func macroRegister(ctx data.Context, className string) (data.GetValue, data.Control) {
	nameV := Arg(ctx, 0)
	macroV := Arg(ctx, 1)
	if nameV == nil || macroV == nil {
		return data.NewNullValue(), nil
	}
	globalMacros.Set(className, nameV.AsString(), macroV)
	return data.NewNullValue(), nil
}

// macroMixin 对齐 Macroable::mixin()：把 $mixin 的 public/protected 方法逐个 invoke，
// 返回的闭包注册成本类的宏。返回 void（vendor 同样没有 return）。
func macroMixin(ctx data.Context, className string) (data.GetValue, data.Control) {
	replace := true
	if v := Arg(ctx, 1); v != nil && !IsNull(v) {
		replace = Truthy(v)
	}
	ctl, _ := InvokeMixinMethods(ctx, Arg(ctx, 0), replace,
		func(name string) bool { return globalMacros.Has(className, name) },
		func(name string, macro data.Value) { globalMacros.Set(className, name, macro) },
	)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

// InvokeMixinMethods 反射遍历 mixin 对象的 public/protected 方法并逐个 invoke，
// 把返回的闭包交给 set；replace 为 false 时 has(name) 命中的方法跳过。
//
// 对应 (new ReflectionClass($mixin))->getMethods(IS_PUBLIC | IS_PROTECTED) 的语义：
// 含继承链上的方法，不含 private；静态方法也在内（ReflectionMethod::invoke 对静态方法同样有效）。
// 非对象的 $mixin 返回错误，交由调用方决定抛什么异常。
func InvokeMixinMethods(
	ctx data.Context,
	mixin data.Value,
	replace bool,
	has func(string) bool,
	set func(string, data.Value),
) (data.Control, error) {
	cv, ok := Unwrap(mixin).(*data.ClassValue)
	if !ok || cv == nil || cv.Class == nil {
		return nil, fmt.Errorf("mixin is not an object")
	}
	vm := cv.GetVM()
	// seen 去重：子类覆盖父类同名方法时，getMethods() 只应产出子类那一份。
	seen := make(map[string]struct{})
	// 上溯深度做保护，避免继承链自引用时死循环。
	for stmt, depth := cv.Class, 0; stmt != nil && depth < 32; depth++ {
		for _, m := range stmt.GetMethods() {
			if m == nil {
				continue
			}
			name := m.GetName()
			key := data.MethodLookupKey(name)
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			if m.GetModifier() == data.ModifierPrivate {
				continue
			}
			if !replace && has(name) {
				continue
			}
			// 不能经 withArgs 造上下文：CreateContext 返回的 ClassMethodContext
			// 才带 $this，mixin 方法体里可能读 $this 或其它属性。
			nctx := cv.CreateContext(m.GetVariables())
			if ctl := data.BindDeclaredArgs(nctx, m, nil); ctl != nil {
				return ctl, nil
			}
			ret, ctl := m.Call(nctx)
			if ctl != nil {
				return ctl, nil
			}
			if v, ok := ret.(data.Value); ok && v != nil {
				set(name, v)
			}
		}
		extend := stmt.GetExtend()
		if extend == nil || vm == nil {
			break
		}
		parent, ctl := vm.GetOrLoadClass(*extend)
		if ctl != nil {
			return ctl, nil
		}
		stmt = parent
	}
	return nil, nil
}

func macroCall(ctx data.Context, className string, static bool) (data.GetValue, data.Control) {
	nameV := Arg(ctx, 0)
	paramsV := Arg(ctx, 1)
	if nameV == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Method %s:: does not exist.", className), "BadMethodCallException")
	}
	name := nameV.AsString()
	macro, ok := globalMacros.Get(className, name)
	if !ok {
		// 对齐 Macroable::__call/__callStatic：抛 BadMethodCallException，不是泛型 Exception。
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Method %s::%s does not exist.", className, name), "BadMethodCallException")
	}
	args := []data.Value{}
	if av, ok := paramsV.(*data.ArrayValue); ok && av != nil {
		args = av.ToValueList()
	}
	// 对齐 Macroable::__call / __callStatic：
	//
	//	$macro = $macro->bindTo($this, static::class);   // 静态时 bindTo(null, static::class)
	//	return $macro(...$parameters);
	//
	// 闭包宏的 $this 必须绑到接收者，不能当作第一个位置实参塞进去。
	switch m := macro.(type) {
	case *data.FuncValue:
		if static {
			return Call(ctx, data.NewBoundFuncValue(m.Value, className, nil), args...)
		}
		if recv := Receiver(ctx); recv != nil {
			return Call(ctx, data.NewBoundFuncValue(m.Value, className, recv), args...)
		}
		return Call(ctx, data.NewBoundFuncValue(m.Value, className, nil), args...)
	case *data.BoundFuncValue:
		if !static {
			if recv := Receiver(ctx); recv != nil {
				return Call(ctx, data.NewBoundFuncValue(m.Value, className, recv), args...)
			}
		}
		return Call(ctx, macro, args...)
	}
	return Call(ctx, macro, args...)
}

// CallMacro 公开宏调用（供 Stringable 等先代理再回退宏）。
func CallMacro(ctx data.Context, className string, static bool) (data.GetValue, data.Control) {
	return macroCall(ctx, className, static)
}

func (s *MacroStore) Set(class, name string, v data.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.macros[class] == nil {
		s.macros[class] = map[string]data.Value{}
	}
	s.macros[class][data.MethodLookupKey(name)] = v
}

func (s *MacroStore) Get(class, name string) (data.Value, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.macros[class]
	if !ok {
		return nil, false
	}
	v, ok := m[data.MethodLookupKey(name)]
	return v, ok
}

func (s *MacroStore) Has(class, name string) bool {
	_, ok := s.Get(class, name)
	return ok
}

func (s *MacroStore) Flush(class string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.macros, class)
}

// HelperFunc 注册全局函数。
func HelperFunc(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) data.FuncStmt {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		ps[i] = node.NewParameter(nil, p, i, nil, nil)
		vs[i] = node.NewVariable(nil, p, i, nil)
	}
	return &helperFunc{name: name, params: ps, vars: vs, fn: fn}
}

// HelperFuncRef 第一个参数按引用（data_set）。
func HelperFuncRef(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) data.FuncStmt {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		if i == 0 {
			ps[i] = node.NewParameterReference(nil, p, i, nil, nil)
		} else {
			ps[i] = node.NewParameter(nil, p, i, nil, nil)
		}
		vs[i] = node.NewVariable(nil, p, i, nil)
	}
	return &helperFunc{name: name, params: ps, vars: vs, fn: fn}
}

type helperFunc struct {
	name   string
	params []data.GetValue
	vars   []data.Variable
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (f *helperFunc) Call(ctx data.Context) (data.GetValue, data.Control) { return f.fn(ctx) }
func (f *helperFunc) GetName() string                                     { return f.name }
func (f *helperFunc) GetParams() []data.GetValue                          { return f.params }
func (f *helperFunc) GetVariables() []data.Variable                       { return f.vars }
