package kit

import (
	"fmt"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MacroStore 按类名存放宏（对齐 Laravel Macroable trait）。
type MacroStore struct {
	mu    sync.RWMutex
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

func macroCall(ctx data.Context, className string, static bool) (data.GetValue, data.Control) {
	nameV := Arg(ctx, 0)
	paramsV := Arg(ctx, 1)
	if nameV == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s:: does not exist.", className))
	}
	name := nameV.AsString()
	macro, ok := globalMacros.Get(className, name)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s::%s does not exist.", className, name))
	}
	args := []data.Value{}
	if av, ok := paramsV.(*data.ArrayValue); ok && av != nil {
		args = av.ToValueList()
	}
	if !static {
		if recv := Receiver(ctx); recv != nil {
			args = append([]data.Value{recv}, args...)
		}
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
	s.macros[class][strings.ToLower(name)] = v
}

func (s *MacroStore) Get(class, name string) (data.Value, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.macros[class]
	if !ok {
		return nil, false
	}
	v, ok := m[strings.ToLower(name)]
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
