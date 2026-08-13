package annotation_test

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	netdata "github.com/php-any/origami/std/net/data"
)

func TestAnalyzeHandlerParams(t *testing.T) {
	method := node.NewMethod(nil, "login", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "request", 0, nil, data.NewBaseType("Spring\\DTO\\Request\\LoginRequest")),
			node.NewParameter(nil, "response", 1, nil, data.NewBaseType("Net\\Http\\Response")),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "request", 0, nil),
			node.NewVariable(nil, "response", 1, nil),
		},
		data.NewBaseType("void"),
	)

	spec, acl := netannotation.AnalyzeHandlerParams(method, nil, "/api/auth/login")
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if len(spec.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(spec.Params))
	}
	if spec.Params[0].Source != netdata.SourceDTO {
		t.Fatalf("param0 source want DTO, got %v", spec.Params[0].Source)
	}
	if spec.Params[1].Source != netdata.SourceResponse {
		t.Fatalf("param1 source want Response, got %v", spec.Params[1].Source)
	}
}

func TestAnalyzeHandlerParamsLogout(t *testing.T) {
	method := node.NewMethod(nil, "logout", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "response", 0, nil, data.NewBaseType("Net\\Http\\Response")),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "response", 0, nil),
		},
		data.NewBaseType("void"),
	)

	spec, acl := netannotation.AnalyzeHandlerParams(method, nil, "/api/auth/logout")
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if len(spec.Params) != 1 || spec.Params[0].Source != netdata.SourceResponse {
		t.Fatalf("logout binding failed: %+v", spec.Params)
	}
}

func TestAnalyzeHandlerParamsUntypedRejected(t *testing.T) {
	method := node.NewMethod(nil, "profile", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "r", 0, nil, nil),
			node.NewParameter(nil, "w", 1, nil, nil),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "r", 0, nil),
			node.NewVariable(nil, "w", 1, nil),
		},
		nil,
	)

	_, acl := netannotation.AnalyzeHandlerParams(method, nil, "/api/profile")
	if acl == nil {
		t.Fatal("expected error for untyped params")
	}
}

func TestAnalyzeHandlerParamsPathVariable(t *testing.T) {
	method := node.NewMethod(nil, "getProduct", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "id", 0, nil, data.NewBaseType("int")),
			node.NewParameter(nil, "response", 1, nil, data.NewBaseType("Net\\Http\\Response")),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "id", 0, nil),
			node.NewVariable(nil, "response", 1, nil),
		},
		data.NewBaseType("void"),
	)

	spec, acl := netannotation.AnalyzeHandlerParams(method, nil, "/api/product/{id}")
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if spec.Params[0].Source != netdata.SourcePath || spec.Params[0].PathKey != "id" {
		t.Fatalf("path binding failed: %+v", spec.Params[0])
	}
	if spec.Params[0].TypeFQN != "int" {
		t.Fatalf("path type = %q", spec.Params[0].TypeFQN)
	}
}

func TestAnalyzeHandlerParamsScalarQueryBinding(t *testing.T) {
	method := node.NewMethod(nil, "hello", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "name", 0, nil, data.NewBaseType("string")),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "name", 0, nil),
		},
		nil,
	)

	spec, acl := netannotation.AnalyzeHandlerParams(method, nil, "/hello")
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if spec.Params[0].Source != netdata.SourceQuery || spec.Params[0].QueryKey != "name" {
		t.Fatalf("query binding failed: %+v", spec.Params[0])
	}
}

func TestAnalyzeHandlerParamsQueryValidationAnnotations(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
namespace Test\Controller;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;
class HelloController {
    public function hello(
        #[NotBlank(message: "name required")]
        #[Size(min: 1, max: 64)]
        string $name
    ): void {}
}
`
	prog, acl := p.ParseString(src, "ctrl.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("run: %v", acl)
	}

	cls, acl := vm.GetOrLoadClass("Test\\Controller\\HelloController")
	if acl != nil {
		t.Fatalf("class: %v", acl)
	}
	method, ok := cls.GetMethod("hello")
	if !ok {
		t.Fatal("method hello not found")
	}

	spec, acl := netannotation.AnalyzeHandlerParams(method, vm, "/hello")
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if !spec.Params[0].Validate {
		t.Fatal("expected Validate=true for query param with constraints")
	}
	if len(spec.Params[0].Constraints) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(spec.Params[0].Constraints))
	}
}

func TestValidationConstraintsDetection(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
namespace Test\DTO;
use Validation\Annotation\Size;
class LoginRequest {
    #[Size(min: 2, max: 100)]
    public string $username;
}
`
	prog, acl := p.ParseString(src, "test.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("run: %v", acl)
	}

	method := node.NewMethod(nil, "login", "public", false,
		[]data.GetValue{
			node.NewParameter(nil, "request", 0, nil, data.NewBaseType("Test\\DTO\\LoginRequest")),
			node.NewParameter(nil, "response", 1, nil, data.NewBaseType("Net\\Http\\Response")),
		},
		nil,
		[]data.Variable{
			node.NewVariable(nil, "request", 0, nil),
			node.NewVariable(nil, "response", 1, nil),
		},
		data.NewBaseType("void"),
	)

	spec, acl := netannotation.AnalyzeHandlerParams(method, vm, "/login")
	if acl != nil {
		t.Fatalf("unexpected acl: %v", acl)
	}
	if !spec.Params[0].Validate {
		t.Fatal("expected Validate=true for DTO with Size constraint")
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
		t.Fatalf("unexpected acl: %v", acl)
	}
	if spec.Params[0].Source != netdata.SourceFormFile || spec.Params[0].FileKey != "avatar" {
		t.Fatalf("uploaded file binding failed: %+v", spec.Params[0])
	}
}
