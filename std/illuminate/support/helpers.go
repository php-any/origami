package support

import (
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
	v, _ := ctx.GetIndexValue(0)
	if v == nil || isNull(v) {
		return data.NewBoolValue(true), nil
	}
	if s, ok := v.(*data.StringValue); ok {
		return data.NewBoolValue(stringsTrim(s.Value) == ""), nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return data.NewBoolValue(len(toEntries(av)) == 0), nil
	}
	return data.NewBoolValue(!truthy(v)), nil
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
