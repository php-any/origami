package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	netdata "github.com/php-any/origami/std/net/data"
	validationannotation "github.com/php-any/origami/std/validation/annotation"
)

func TestHomeControllerHelloQueryReturn(t *testing.T) {
	root := filepath.Join("..", "..", "..", "examples", "spring")
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	netannotation.Load(vm)
	validationannotation.Load(vm)
	Load(vm)
	ctx := vm.CreateContext(nil)

	prog, acl := p.ParseFile(filepath.Join(root, "src/Controller/HomeController.php"))
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("load: %v", acl)
	}

	cls, acl := vm.GetOrLoadClass("Spring\\Controller\\HomeController")
	if acl != nil {
		t.Fatalf("class: %v", acl)
	}
	method, ok := cls.GetMethod("hello")
	if !ok {
		t.Fatal("hello method missing")
	}
	for i, p := range method.GetParams() {
		param, ok := p.(*node.Parameter)
		if !ok {
			t.Fatalf("param %d type %T", i, p)
		}
		t.Logf("param %s annotations=%d", param.Name, len(param.Annotations))
	}

	spec, acl := netannotation.AnalyzeHandlerParams(method, vm, "/hello")
	if acl != nil {
		t.Fatalf("analyze: %v", acl)
	}
	if len(spec.Params) != 2 {
		t.Fatalf("expected 2 bindings, got %+v", spec.Params)
	}

	inst, acl := node.InstantiateController(cls, ctx)
	if acl != nil {
		t.Fatalf("instantiate: %v", acl)
	}

	req := httptest.NewRequest(http.MethodGet, "/hello?name=ab&age=25", nil)
	_, requestClass := beginRequest(req)
	reqProxy := data.NewProxyValue(requestClass, ctx)
	rec := httptest.NewRecorder()
	_, responseClass := beginResponse(rec, req)
	resProxy := data.NewProxyValue(responseClass, ctx)

	rt := netdata.Route{
		Method:      "GET",
		Path:        "/hello",
		Target:      method,
		Receiver:    inst,
		HandlerSpec: spec,
	}

	ret, acl := executeControllerMethod(vm, ctx, rt, reqProxy, resProxy)
	if acl != nil {
		t.Fatalf("execute: %v", acl)
	}
	t.Logf("ret=%#v", ret)
	if acl := writeHandlerReturnValue(resProxy, ret); acl != nil {
		t.Fatalf("write return: %v", acl)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() == 0 {
		t.Fatal("empty response body")
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v raw=%s", err, rec.Body.String())
	}
	if body["message"] != "success" {
		t.Fatalf("body=%#v", body)
	}
}
