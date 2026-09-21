package support

import (
	"html"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// registerHelpers 仅在 laravel13 / vendoraccel 路径注册，避免污染默认 zy 解释器。
func registerHelpers(vm data.VM) {
	for _, fn := range []data.FuncStmt{
		newHelperFunc("collect", []string{"value"}, helperCollect),
		newHelperFunc("data_get", []string{"target", "key", "default"}, helperDataGet),
		newHelperFunc("data_set", []string{"target", "key", "value", "overwrite"}, helperDataSet),
		newValueHelper(),
		newWithHelper(),
		newHelperFunc("filled", []string{"value"}, helperFilled),
		newHelperFunc("blank", []string{"value"}, helperBlank),
		newHelperFunc("class_basename", []string{"class"}, helperClassBasename),
		newEHelper(),
	} {
		vm.AddFunc(fn)
	}
}

func newHelperFunc(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) data.FuncStmt {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		if name == "data_set" && i == 0 {
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

// newWithHelper 对齐 Laravel with($value, $callback = null)。
func newWithHelper() data.FuncStmt {
	return &helperFunc{
		name: "with",
		params: []data.GetValue{
			node.NewParameter(nil, "value", 0, nil, nil),
			node.NewParameter(nil, "callback", 1, data.NewNullValue(), nil),
		},
		vars: []data.Variable{
			node.NewVariable(nil, "value", 0, nil),
			node.NewVariable(nil, "callback", 1, nil),
		},
		fn: helperWith,
	}
}

// newValueHelper 对齐 Laravel value($value, ...$args)。
func newValueHelper() data.FuncStmt {
	return &helperFunc{
		name: "value",
		params: []data.GetValue{
			node.NewParameter(nil, "value", 0, nil, nil),
			node.NewParameters(nil, "args", 1, nil, nil),
		},
		vars: []data.Variable{
			node.NewVariable(nil, "value", 0, nil),
			node.NewVariable(nil, "args", 1, nil),
		},
		fn: helperValue,
	}
}

// newEHelper 对齐 Laravel e()：Blade `echo e(...)` 不再解释 helpers.php。
func newEHelper() data.FuncStmt {
	return &helperFunc{
		name: "e",
		params: []data.GetValue{
			node.NewParameter(nil, "value", 0, nil, nil),
			node.NewParameter(nil, "doubleEncode", 1, data.NewBoolValue(true), nil),
		},
		vars: []data.Variable{
			node.NewVariable(nil, "value", 0, nil),
			node.NewVariable(nil, "doubleEncode", 1, nil),
		},
		fn: helperE,
	}
}

func helperE(ctx data.Context) (data.GetValue, data.Control) {
	v := unwrapValue(ctxIndexValue(ctx, 0))
	if v == nil {
		return data.NewStringValue(""), nil
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return data.NewStringValue(""), nil
	}

	if cv, ok := v.(*data.ClassValue); ok {
		if classIsA(cv, "Illuminate\\Contracts\\Support\\DeferringDisplayableValue") {
			ret, ctl := callClassMethod(cv, "resolveDisplayableValue")
			if ctl != nil {
				return nil, ctl
			}
			if rv, ok := ret.(data.Value); ok {
				v = unwrapValue(rv)
				cv, _ = v.(*data.ClassValue)
			} else {
				v = data.NewNullValue()
				cv = nil
			}
		}
		if cv != nil {
			if classIsA(cv, "Illuminate\\Contracts\\Support\\Htmlable") {
				ret, ctl := callClassMethod(cv, "toHtml")
				if ctl != nil {
					return nil, ctl
				}
				if ret == nil {
					return data.NewStringValue(""), nil
				}
				if rv, ok := ret.(data.Value); ok {
					if _, isNull := rv.(*data.NullValue); isNull {
						return data.NewStringValue(""), nil
					}
					return data.NewStringValue(rv.AsString()), nil
				}
				return data.NewStringValue(""), nil
			}
			if classIsA(cv, "BackedEnum") {
				pv, ctl := cv.GetProperty("value")
				if ctl != nil {
					return nil, ctl
				}
				v = unwrapValue(pv)
			}
		}
	}

	s := ""
	if val, ok := v.(data.Value); ok {
		if _, isNull := val.(*data.NullValue); isNull {
			return data.NewStringValue(""), nil
		}
		if sv, ok := val.(*data.StringValue); ok {
			s = sv.Value
		} else {
			s = val.AsString()
		}
	}
	if sv, ok := v.(*data.StringValue); ok && !stringsContainsHTMLSpecial(s) {
		return sv, nil
	}
	if !stringsContainsHTMLSpecial(s) {
		return data.NewStringValue(s), nil
	}
	return data.NewStringValue(html.EscapeString(s)), nil
}

func stringsContainsHTMLSpecial(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&', '<', '>', '"', '\'':
			return true
		}
	}
	return false
}

func callClassMethod(cv *data.ClassValue, name string) (data.GetValue, data.Control) {
	if cv == nil {
		return nil, nil
	}
	m, ok := cv.GetMethod(name)
	if !ok {
		return nil, nil
	}
	return m.Call(cv.CreateContext(m.GetVariables()))
}

func classIsA(cv *data.ClassValue, name string) bool {
	if cv == nil || cv.Class == nil {
		return false
	}
	if cv.Class.GetName() == name {
		return true
	}
	for _, iface := range cv.Class.GetImplements() {
		if iface == name {
			return true
		}
	}
	vm := cv.GetVM()
	if vm == nil {
		return false
	}
	last := cv.Class
	for last != nil && last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok || next == nil {
			break
		}
		if next.GetName() == name {
			return true
		}
		for _, iface := range next.GetImplements() {
			if iface == name {
				return true
			}
		}
		last = next
	}
	return false
}

func helperCollect(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return newCollectionInstance(ctx, v)
}

func helperDataGet(ctx data.Context) (data.GetValue, data.Control) {
	return arrGet(ctx)
}

func helperDataSet(ctx data.Context) (data.GetValue, data.Control) {
	return arrSet(ctx)
}

func helperValue(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return laravelValue(ctx, v, variadicValues(ctx, 1)...)
}

// helperWith 对齐 Laravel with($value, $callback = null)：
// 无回调返回 $value；有回调必须返回 $callback($value)，不能丢掉回调结果。
// Handler::shouldntReport 依赖 with(Limit::none(), fn => false) 得到 false。
func helperWith(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		v = data.NewNullValue()
	}
	cb, _ := ctx.GetIndexValue(1)
	if cb == nil || isNull(cb) {
		return v, nil
	}
	return callValue(ctx, cb, v)
}

// laravelValue 对齐 Laravel value($value, ...$args)：仅 Closure 会被调用。
func laravelValue(ctx data.Context, v data.Value, args ...data.Value) (data.GetValue, data.Control) {
	if v == nil {
		return data.NewNullValue(), nil
	}
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return callValue(ctx, v, args...)
	default:
		return v, nil
	}
}

func variadicValues(ctx data.Context, index int) []data.Value {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil || isNull(v) {
		return nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return av.ToValueList()
	}
	return []data.Value{v}
}

func helperFilled(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	b, ctl := helperBlank(ctx)
	if ctl != nil {
		return nil, ctl
	}
	bv, _ := b.(data.AsBool).AsBool()
	_ = v
	return data.NewBoolValue(!bv), nil
}

func helperBlank(ctx data.Context) (data.GetValue, data.Control) {
	v := unwrapValue(ctxIndexValue(ctx, 0))
	if v == nil || isNull(v) {
		return data.NewBoolValue(true), nil
	}
	// Laravel blank()：is_numeric / is_bool 一律非 blank，filled(0) 为 true。
	// 不能走 PHP empty/truthy，否则仪表盘 Stat 的 0 会被当成占位符藏掉。
	switch v.(type) {
	case *data.IntValue, *data.FloatValue, *data.BoolValue:
		return data.NewBoolValue(false), nil
	}
	if s, ok := v.(*data.StringValue); ok {
		return data.NewBoolValue(stringsTrim(s.Value) == ""), nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return data.NewBoolValue(len(toEntries(av)) == 0), nil
	}
	return data.NewBoolValue(!truthy(v)), nil
}

func ctxIndexValue(ctx data.Context, i int) data.Value {
	if ctx == nil {
		return nil
	}
	v, _ := ctx.GetIndexValue(i)
	return v
}

func helperClassBasename(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	s := classBasenameSource(v)
	if i := lastSlash(s); i >= 0 {
		s = s[i+1:]
	}
	return data.NewStringValue(s), nil
}

// classBasenameSource 对齐 Laravel helpers.php 的 class_basename：
// is_object($class) ? get_class($class) : $class，再 basename。
// $this 是 *ThisValue，不能走 AsString()（会得到 Object(Foo\Bar)，basename 变成 "Bar)"）。
func classBasenameSource(v data.Value) string {
	if v == nil {
		return ""
	}
	if name := objectClassName(v); name != "" {
		return name
	}
	return v.AsString()
}

func objectClassName(v data.Value) string {
	switch t := v.(type) {
	case *data.ThisValue:
		if t != nil && t.Class != nil {
			return t.Class.GetName()
		}
	case *data.ClassValue:
		if t != nil && t.Class != nil {
			return t.Class.GetName()
		}
	}
	return ""
}

func stringsTrim(s string) string {
	return trimSpace(s)
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end {
		c := s[start]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		start++
	}
	for end > start {
		c := s[end-1]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		end--
	}
	return s[start:end]
}

func lastSlash(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\\' || s[i] == '/' {
			return i
		}
	}
	return -1
}
