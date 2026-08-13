package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	netdata "github.com/php-any/origami/std/net/data"
)

func TestResolveFormFileParam(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "demo")
	part, err := writer.CreateFormFile("avatar", "photo.png")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, "png-bytes")
	_ = writer.Close()

	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	Load(vm)
	ctx := vm.CreateContext(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/upload/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, requestClass := beginRequest(req)
	reqProxy := data.NewProxyValue(requestClass, ctx)
	rec := httptest.NewRecorder()
	_, responseClass := beginResponse(rec, req)
	resProxy := data.NewProxyValue(responseClass, ctx)

	spec := netdata.HandlerSpec{
		Params: []netdata.ParamBinding{
			{Name: "avatar", TypeFQN: uploadedFileTypeFQN, Source: netdata.SourceFormFile, Index: 0, FileKey: "avatar"},
		},
	}

	args, acl := resolveHandlerArgs(vm, ctx, spec, reqProxy, resProxy)
	if acl != nil {
		t.Fatalf("resolve failed: %v", acl)
	}
	fileProxy, ok := args[0].(*data.ProxyValue)
	if !ok {
		t.Fatalf("want UploadedFile proxy, got %#v", args[0])
	}
	nameMethod, _ := fileProxy.GetMethod("originalName")
	vars := nameMethod.GetVariables()
	callCtx := fileProxy.CreateContext(vars)
	nameVal, acl := nameMethod.Call(callCtx)
	if acl != nil {
		t.Fatalf("originalName: %v", acl)
	}
	if nameVal.(data.AsString).AsString() != "photo.png" {
		t.Fatalf("want photo.png, got %q", nameVal.(data.AsString).AsString())
	}
}

func TestBindMultipartDTO(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "readme")
	part, err := writer.CreateFormFile("file", "readme.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, "hello upload")
	_ = writer.Close()

	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	Load(vm)

	src := `<?php
use Net\Http\UploadedFile;
class UploadFileRequest {
    public string $title;
    public UploadedFile $file;
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

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, requestClass := beginRequest(req)
	reqCV := data.NewProxyValue(requestClass, ctx)

	bindMethod, _ := reqCV.GetMethod("bind")
	vars := bindMethod.GetVariables()
	bindCtx := reqCV.CreateContext(vars)
	bindCtx.SetVariableValue(vars[0], data.NewStringValue("UploadFileRequest"))

	dto, acl := bindMethod.Call(bindCtx)
	if acl != nil {
		t.Fatalf("bind failed: %v", acl)
	}
	cv, ok := dto.(*data.ClassValue)
	if !ok {
		t.Fatalf("want ClassValue, got %T", dto)
	}

	titleProp, ok := cv.GetPropertyStmt("title")
	if !ok {
		t.Fatal("title property not found")
	}
	titleGV, acl := titleProp.GetValue(cv)
	if acl != nil {
		t.Fatalf("read title: %v", acl)
	}
	if titleGV.(data.AsString).AsString() != "readme" {
		t.Fatalf("title want readme, got %q", titleGV.(data.AsString).AsString())
	}
	fileProp, ok := cv.GetPropertyStmt("file")
	if !ok {
		t.Fatal("file property not found")
	}
	fileGV, acl := fileProp.GetValue(cv)
	if acl != nil {
		t.Fatalf("read file: %v", acl)
	}
	fileProxy, ok := fileGV.(*data.ProxyValue)
	if !ok {
		t.Fatalf("file want proxy, got %T", fileGV)
	}
	nameMethod, _ := fileProxy.GetMethod("originalName")
	nameCtx := fileProxy.CreateContext(nameMethod.GetVariables())
	name, acl := nameMethod.Call(nameCtx)
	if acl != nil {
		t.Fatalf("originalName: %v", acl)
	}
	if name.(data.AsString).AsString() != "readme.txt" {
		t.Fatalf("file name want readme.txt, got %q", name.(data.AsString).AsString())
	}
}

func TestAnalyzeHandlerParamsUploadedFile(t *testing.T) {
	method := node.NewMethod(nil, "uploadAvatar", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "avatar", 0, nil, data.NewBaseType("Net\\Http\\UploadedFile")),
			node.NewParameter(nil, "response", 1, nil, data.NewBaseType("Net\\Http\\Response")),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "avatar", 0, nil),
			node.NewVariable(nil, "response", 1, nil),
		},
		data.NewBaseType("void"),
	)
	spec, acl := netannotation.AnalyzeHandlerParams(method, nil, "/api/upload/avatar")
	if acl != nil {
		t.Fatalf("analyze: %v", acl)
	}
	if len(spec.Params) != 2 {
		t.Fatalf("want 2 params, got %+v", spec.Params)
	}
	if spec.Params[0].Source != netdata.SourceFormFile || spec.Params[0].FileKey != "avatar" {
		t.Fatalf("avatar binding want SourceFormFile, got %+v", spec.Params[0])
	}
}
