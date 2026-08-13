package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
)

func TestRequestBindQueryIntDTO(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
class UserListQuery {
    public int $min_age = 0;
    public int $limit = 10;
}
`
	prog, acl := p.ParseString(src, "dto.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("load class: %v", acl)
	}

	req := httptest.NewRequest(http.MethodGet, "/users?min_age=25&limit=5", nil)
	_, requestClass := beginRequest(req)
	reqCV := data.NewProxyValue(requestClass, ctx)

	bindMethod, _ := reqCV.GetMethod("bind")
	vars := bindMethod.GetVariables()
	bindCtx := reqCV.CreateContext(vars)
	bindCtx.SetVariableValue(vars[0], data.NewStringValue("UserListQuery"))

	dto, acl := bindMethod.Call(bindCtx)
	if acl != nil {
		t.Fatalf("bind failed: %v", acl)
	}
	cv, ok := dto.(*data.ClassValue)
	if !ok {
		t.Fatalf("expected ClassValue, got %T", dto)
	}

	minAgeProp, ok := cv.GetPropertyStmt("min_age")
	if !ok {
		t.Fatal("min_age property not found")
	}
	minAgeGV, acl := minAgeProp.GetValue(cv)
	if acl != nil {
		t.Fatalf("read min_age: %v", acl)
	}
	minAgeInt, ok := minAgeGV.(*data.IntValue)
	if !ok || minAgeInt.Value != 25 {
		t.Fatalf("min_age want 25, got %#v", minAgeGV)
	}

	limitProp, ok := cv.GetPropertyStmt("limit")
	if !ok {
		t.Fatal("limit property not found")
	}
	limitGV, acl := limitProp.GetValue(cv)
	if acl != nil {
		t.Fatalf("read limit: %v", acl)
	}
	limitInt, ok := limitGV.(*data.IntValue)
	if !ok || limitInt.Value != 5 {
		t.Fatalf("limit want 5, got %#v", limitGV)
	}
}

func TestRequestBindQueryPreservesDefaults(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
class UserListQuery {
    public int $min_age = 0;
    public int $limit = 10;
}
`
	prog, acl := p.ParseString(src, "dto.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("load class: %v", acl)
	}

	req := httptest.NewRequest(http.MethodGet, "/users?min_age=25", nil)
	_, requestClass := beginRequest(req)
	reqCV := data.NewProxyValue(requestClass, ctx)

	bindMethod, _ := reqCV.GetMethod("bind")
	vars := bindMethod.GetVariables()
	bindCtx := reqCV.CreateContext(vars)
	bindCtx.SetVariableValue(vars[0], data.NewStringValue("UserListQuery"))

	dto, acl := bindMethod.Call(bindCtx)
	if acl != nil {
		t.Fatalf("bind failed: %v", acl)
	}
	cv := dto.(*data.ClassValue)
	limitProp, ok := cv.GetPropertyStmt("limit")
	if !ok {
		t.Fatal("limit property not found")
	}
	limitGV, acl := limitProp.GetValue(cv)
	if acl != nil {
		t.Fatalf("read limit: %v", acl)
	}
	limitInt := limitGV.(*data.IntValue)
	if limitInt.Value != 10 {
		t.Fatalf("limit default want 10, got %d", limitInt.Value)
	}
}
