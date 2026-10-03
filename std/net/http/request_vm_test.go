package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

type requestVMCallback struct {
	variables []data.Variable
	call      func(data.Context)
}

func (f *requestVMCallback) GetName() string               { return "request_vm_callback" }
func (f *requestVMCallback) GetParams() []data.GetValue    { return nil }
func (f *requestVMCallback) GetVariables() []data.Variable { return f.variables }
func (f *requestVMCallback) Call(ctx data.Context) (data.GetValue, data.Control) {
	f.call(ctx)
	return data.NewNullValue(), nil
}

func TestNormalAndHotHandlerUseOneRequestVM(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	root := base.CreateContext(nil)
	variables := []data.Variable{data.NewVariable("request", 0, nil), data.NewVariable("response", 1, nil)}
	var middlewareVM *runtime.RequestVM
	shutdown := 0
	middleware := &requestVMCallback{variables: append(append([]data.Variable(nil), variables...), data.NewVariable("next", 2, nil)), call: func(ctx data.Context) {
		middlewareVM = ctx.GetVM().(*runtime.RequestVM)
		middlewareVM.EnsureGlobalZVal("token").StoreRaw(data.NewStringValue("request"))
		next, _ := ctx.GetIndexValue(2)
		call := ctx.CreateContext(next.(*data.FuncValue).Value.GetVariables())
		request, _ := ctx.GetIndexValue(0)
		response, _ := ctx.GetIndexValue(1)
		call.SetIndexZVal(0, data.NewZVal(request))
		call.SetIndexZVal(1, data.NewZVal(response))
		if _, ctl := next.(*data.FuncValue).Call(call); ctl != nil {
			t.Error(ctl.AsString())
		}
	}}
	callback := &requestVMCallback{variables: variables, call: func(ctx data.Context) {
		vm := ctx.GetVM().(*runtime.RequestVM)
		if vm != middlewareVM {
			t.Error("middleware and handler used different request VMs")
		}
		if vm.EnsureGlobalZVal("token").ReadValue().AsString() != "request" {
			t.Error("middleware state missing")
		}
		vm.AddShutdownCallback(data.NewFuncValue(&requestVMCallback{call: func(data.Context) { shutdown++ }}))
	}}
	wrap, err := newMiddleware(middleware, root)
	if err != nil {
		t.Fatal(err)
	}
	for _, handler := range []http.Handler{Handler{Value: callback, Ctx: root}, HotHandler{Value: callback, Ctx: root}} {
		wrap(handler).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	}
	if shutdown != 2 {
		t.Fatalf("shutdown=%d", shutdown)
	}
	if base.EnsureGlobalZVal("token").ReadValue().AsString() != "" {
		t.Fatal("HTTP state escaped to VM")
	}
}
