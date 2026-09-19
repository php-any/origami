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
		newHelperFunc("value", []string{"value", "args"}, helperValue),
		newHelperFunc("with", []string{"value", "callback"}, helperWith),
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
	arg, hasArg := ctx.GetIndexValue(1)
	switch cb := v.(type) {
	case *data.FuncValue:
		if cb == nil || cb.Value == nil {
			return v, nil
		}
		if hasArg && arg != nil && !isNull(arg) {
			if av, ok := arg.(data.Value); ok {
				return callValue(ctx, cb, av)
			}
		}
		return cb.Value.Call(ctx.CreateContext(cb.Value.GetVariables()))
	case *data.BoundFuncValue:
		if cb == nil {
			return v, nil
		}
		if hasArg && arg != nil && !isNull(arg) {
			if av, ok := arg.(data.Value); ok {
				return callValue(ctx, cb, av)
			}
		}
		return cb.Call(ctx)
	}
	return v, nil
}

func helperWith(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	if cb != nil && !isNull(cb) {
		_, ctl := callValue(ctx, cb, v)
		if ctl != nil {
			return nil, ctl
		}
	}
	return v, nil
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
