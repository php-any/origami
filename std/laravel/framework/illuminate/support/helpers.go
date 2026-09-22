package support

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

func registerHelpers(vm data.VM) {
	vm.AddFunc(newValueHelper())
	vm.AddFunc(newWithHelper())
	vm.AddFunc(kit.HelperFunc("filled", []string{"value"}, helperFilled))
	vm.AddFunc(kit.HelperFunc("blank", []string{"value"}, helperBlank))
	vm.AddFunc(kit.HelperFunc("class_basename", []string{"class"}, helperClassBasename))
}

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

func helperValue(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return laravelValue(ctx, v, variadicValues(ctx, 1)...)
}

func helperWith(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		v = data.NewNullValue()
	}
	cb, _ := ctx.GetIndexValue(1)
	if cb == nil || kit.IsNull(cb) {
		return v, nil
	}
	return kit.Call(ctx, cb, v)
}

func laravelValue(ctx data.Context, v data.Value, args ...data.Value) (data.GetValue, data.Control) {
	if v == nil {
		return data.NewNullValue(), nil
	}
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return kit.Call(ctx, v, args...)
	default:
		return v, nil
	}
}

func variadicValues(ctx data.Context, index int) []data.Value {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil || kit.IsNull(v) {
		return nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return av.ToValueList()
	}
	return []data.Value{v}
}

func helperFilled(ctx data.Context) (data.GetValue, data.Control) {
	b, ctl := helperBlank(ctx)
	if ctl != nil {
		return nil, ctl
	}
	bv, _ := b.(data.AsBool).AsBool()
	return data.NewBoolValue(!bv), nil
}

func helperBlank(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Unwrap(kit.Arg(ctx, 0))
	if v == nil || kit.IsNull(v) {
		return data.NewBoolValue(true), nil
	}
	switch v.(type) {
	case *data.IntValue, *data.FloatValue, *data.BoolValue:
		return data.NewBoolValue(false), nil
	}
	if s, ok := v.(*data.StringValue); ok {
		return data.NewBoolValue(trimSpace(s.Value) == ""), nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return data.NewBoolValue(len(kit.Entries(av)) == 0), nil
	}
	return data.NewBoolValue(!kit.Truthy(v)), nil
}

func helperClassBasename(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	s := classBasenameSource(v)
	if i := lastSlash(s); i >= 0 {
		s = s[i+1:]
	}
	return data.NewStringValue(s), nil
}

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
