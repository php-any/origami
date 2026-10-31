package script

import (
	"context"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

type callback struct{ t *testing.T }

func (f *callback) GetName() string               { return "isolated" }
func (f *callback) GetParams() []data.GetValue    { return nil }
func (f *callback) GetVariables() []data.Variable { return nil }
func (f *callback) Call(c data.Context) (data.GetValue, data.Control) {
	vm := c.GetVM().(*runtime.RequestVM)
	if vm.EnsureGlobalZVal("probe").ReadValue().AsString() != "" {
		f.t.Error("state leaked")
	}
	vm.EnsureGlobalZVal("probe").StoreRaw(data.NewStringValue("private"))
	if c.GoContext().Err() != context.Canceled {
		f.t.Error("cancellation was not bound")
	}
	return nil, nil
}
func TestInvokeIsolationAndCancellation(t *testing.T) {
	base := runtime.NewVM(parser.NewParser())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for i := 0; i < 2; i++ {
		if e := Invoke(base, &callback{t: t}, ctx, func(data.Context) []data.Value { return nil }); e != nil {
			t.Fatal(e)
		}
	}
}

// Bound closures must keep their explicit visibility scope across a request VM.
type boundProbe struct{ callback }

func (f *boundProbe) Call(ctx data.Context) (data.GetValue, data.Control) {
	bound, ok := ctx.(*data.BoundContext)
	if !ok || !bound.ExplicitBind || bound.ScopeClass != "Probe" {
		f.t.Errorf("bound closure context lost: %T", ctx)
	}
	return nil, nil
}
func TestBoundCallback(t *testing.T) {
	base := runtime.NewVM(parser.NewParser())
	ctx := base.CreateContext([]data.Variable{data.NewVariable("callback", 0, nil)})
	value := data.NewBoundFuncValue(&boundProbe{callback{t: t}}, "Probe", nil)
	ctx.SetIndexZVal(0, data.NewZVal(value))
	fn, err := Callback(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := Invoke(base, fn, context.Background(), func(data.Context) []data.Value { return nil }); err != nil {
		t.Fatal(err)
	}
}

type noParametersProbe struct{ callback }

func (f *noParametersProbe) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("local", 0, nil)}
}
func (f *noParametersProbe) Call(ctx data.Context) (data.GetValue, data.Control) {
	local, _ := ctx.GetIndexValue(0)
	if _, ok := local.(*data.NullValue); !ok {
		f.t.Errorf("extra argument overwrote a local: %T", local)
	}
	if len(ctx.GetFlatCallArgs()) != 1 {
		f.t.Error("extra argument was not retained for func_get_args")
	}
	return nil, nil
}
func TestExtraCallbackArgumentsDoNotOverwriteLocals(t *testing.T) {
	base := runtime.NewVM(parser.NewParser())
	err := Invoke(base, &noParametersProbe{callback{t: t}}, context.Background(), func(data.Context) []data.Value { return []data.Value{data.NewStringValue("connection")} })
	if err != nil {
		t.Fatal(err)
	}
}
