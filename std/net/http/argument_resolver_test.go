package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netdata "github.com/php-any/origami/std/net/data"
	"github.com/php-any/origami/std/validation"
)

func TestResolveHandlerArgsNoParams(t *testing.T) {
	args, acl := resolveHandlerArgs(nil, nil, netdata.HandlerSpec{}, nil, nil)
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if args == nil {
		t.Fatal("expected non-nil empty slice for zero-param handler")
	}
	if len(args) != 0 {
		t.Fatalf("want empty args, got len=%d", len(args))
	}
}

func TestResolvePathParamInt(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	ctx := vm.CreateContext(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/product/42", nil)
	req.SetPathValue("id", "42")
	_, requestClass := beginRequest(req)
	reqProxy := data.NewProxyValue(requestClass, ctx)

	spec := netdata.HandlerSpec{
		Params: []netdata.ParamBinding{
			{Name: "id", TypeFQN: "int", Source: netdata.SourcePath, Index: 0, PathKey: "id"},
		},
	}

	args, acl := resolveHandlerArgs(vm, ctx, spec, reqProxy, nil)
	if acl != nil {
		t.Fatalf("resolve failed: %v", acl)
	}
	iv, ok := args[0].(*data.IntValue)
	if !ok || iv.Value != 42 {
		t.Fatalf("want int 42, got %#v", args[0])
	}
}

func TestResolveQueryParamString(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	ctx := vm.CreateContext(nil)

	req := httptest.NewRequest(http.MethodGet, "/hello?name=Origami", nil)
	_, requestClass := beginRequest(req)
	reqProxy := data.NewProxyValue(requestClass, ctx)

	spec := netdata.HandlerSpec{
		Params: []netdata.ParamBinding{
			{Name: "name", TypeFQN: "string", Source: netdata.SourceQuery, Index: 0, QueryKey: "name"},
		},
	}

	args, acl := resolveHandlerArgs(vm, ctx, spec, reqProxy, nil)
	if acl != nil {
		t.Fatalf("resolve failed: %v", acl)
	}
	sv, ok := args[0].(*data.StringValue)
	if !ok || sv.AsString() != "Origami" {
		t.Fatalf("want name Origami, got %#v", args[0])
	}
}

func TestResolveQueryParamValidationFailure(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	ctx := vm.CreateContext(nil)

	src := `<?php
class C {
    public function hello(
        #[Validation\Annotation\NotBlank(message: "name required")]
        string $name
    ): void {}
}
`
	prog, acl := p.ParseString(src, "c.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("run: %v", acl)
	}
	cls, acl := vm.GetOrLoadClass("C")
	if acl != nil {
		t.Fatalf("class: %v", acl)
	}
	method, _ := cls.GetMethod("hello")
	param := method.GetParams()[0].(*node.Parameter)
	constraints := validation.ConstraintAnnotations(param.Annotations)
	if len(constraints) == 0 {
		t.Fatal("expected parameter constraint annotations")
	}

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	_, requestClass := beginRequest(req)
	reqProxy := data.NewProxyValue(requestClass, ctx)
	rec := httptest.NewRecorder()
	_, responseClass := beginResponse(rec, req)
	resProxy := data.NewProxyValue(responseClass, ctx)

	spec := netdata.HandlerSpec{
		Params: []netdata.ParamBinding{
			{
				Name: "name", TypeFQN: "string", Source: netdata.SourceQuery,
				Index: 0, QueryKey: "name", Validate: true, Constraints: constraints,
			},
		},
	}

	args, acl := resolveHandlerArgs(vm, ctx, spec, reqProxy, resProxy)
	if args != nil {
		t.Fatalf("expected nil args on validation failure, got %#v", args)
	}
	if acl != nil {
		t.Fatalf("unexpected acl after validation response: %v", acl)
	}
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestValidateScalarBindingNotBlank(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	ctx := vm.CreateContext(nil)

	src := `<?php
class C {
    public function hello(
        #[Validation\Annotation\NotBlank(message: "name required")]
        string $name
    ): void {}
}
`
	prog, acl := p.ParseString(src, "c.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("run: %v", acl)
	}
	cls, _ := vm.GetOrLoadClass("C")
	method, _ := cls.GetMethod("hello")
	param := method.GetParams()[0].(*node.Parameter)
	constraints := validation.ConstraintAnnotations(param.Annotations)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, responseClass := beginResponse(rec, req)
	resProxy := data.NewProxyValue(responseClass, ctx)

	binding := netdata.ParamBinding{
		Name: "name", Validate: true, Constraints: constraints,
	}
	done, acl := validateScalarBinding(resProxy, binding, data.NewStringValue(""))
	if !done {
		t.Fatal("expected validation to stop handler")
	}
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestResolveQueryParamIntNotBlankValidationFailure(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	ctx := vm.CreateContext(nil)

	src := `<?php
class C {
    public function hello(
        #[Validation\Annotation\NotBlank(message: "age 不能为空")]
        int $age
    ): void {}
}
`
	prog, acl := p.ParseString(src, "c.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("run: %v", acl)
	}
	cls, acl := vm.GetOrLoadClass("C")
	if acl != nil {
		t.Fatalf("class: %v", acl)
	}
	method, _ := cls.GetMethod("hello")
	param := method.GetParams()[0].(*node.Parameter)
	constraints := validation.ConstraintAnnotations(param.Annotations)

	req := httptest.NewRequest(http.MethodGet, "/hello?name=test", nil)
	_, requestClass := beginRequest(req)
	reqProxy := data.NewProxyValue(requestClass, ctx)
	rec := httptest.NewRecorder()
	_, responseClass := beginResponse(rec, req)
	resProxy := data.NewProxyValue(responseClass, ctx)

	spec := netdata.HandlerSpec{
		Params: []netdata.ParamBinding{
			{
				Name: "age", TypeFQN: "int", Source: netdata.SourceQuery,
				Index: 0, QueryKey: "age", Validate: true, Constraints: constraints,
			},
		},
	}

	args, acl := resolveHandlerArgs(vm, ctx, spec, reqProxy, resProxy)
	if args != nil {
		t.Fatalf("expected nil args on validation failure, got %#v", args)
	}
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}
