package validation_test

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/validation"
)

func TestValidateObjectLoginRequest(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
namespace Test\DTO;
use Validation\Annotation\Size;
use Validation\Annotation\Pattern;
class LoginRequest {
    #[Size(min: 2, max: 100, message: "username size invalid")]
    #[Pattern(regexp: "/^[a-zA-Z0-9_-]+$/", message: "username pattern invalid")]
    public string $username;
    #[Size(min: 6, max: 100)]
    public string $password;
}
`
	prog, acl := p.ParseString(src, "dto.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("load: %v", acl)
	}

	cls, acl := vm.GetOrLoadClass("Test\\DTO\\LoginRequest")
	if acl != nil {
		t.Fatalf("class: %v", acl)
	}
	inst, acl := cls.GetValue(ctx)
	if acl != nil {
		t.Fatalf("instance: %v", acl)
	}
	cv := inst.(*data.ClassValue)

	cv.SetProperty("username", data.NewStringValue("ab"))
	cv.SetProperty("password", data.NewStringValue("secret"))
	if v := validation.ValidateObject(cv); len(v) != 0 {
		t.Fatalf("expected pass, got %#v", v)
	}

	cv.SetProperty("username", data.NewStringValue("a"))
	if len(validation.ValidateObject(cv)) == 0 {
		t.Fatal("expected size violation")
	}

	cv.SetProperty("username", data.NewStringValue("bad!"))
	if len(validation.ValidateObject(cv)) == 0 {
		t.Fatal("expected pattern violation")
	}
}

func TestValidateConstraintsScalar(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
class C {
    public function f(
        #[Validation\Annotation\NotBlank(message: "name required")]
        #[Validation\Annotation\Size(min: 2, max: 10)]
        string $name
    ): void {}
}
`
	prog, acl := p.ParseString(src, "c.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("run: %v", acl)
	}
	cls, acl := vm.GetOrLoadClass("C")
	if acl != nil {
		t.Fatalf("class: %v", acl)
	}
	method, _ := cls.GetMethod("f")
	param := method.GetParams()[0].(*node.Parameter)
	constraints := validation.ConstraintAnnotations(param.Annotations)

	if v := validation.ValidateConstraints("name", "name", constraints, data.NewStringValue("")); len(v) == 0 {
		t.Fatal("expected NotBlank violation")
	}
	if v := validation.ValidateConstraints("name", "name", constraints, data.NewStringValue("a")); len(v) == 0 {
		t.Fatal("expected Size violation")
	}
	if v := validation.ValidateConstraints("name", "name", constraints, data.NewStringValue("ok")); len(v) != 0 {
		t.Fatalf("expected pass, got %#v", v)
	}
}

func TestValidatorClassValidate(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)

	src := `<?php
class SimpleDto {
    #[Validation\Annotation\NotBlank(message: "name required")]
    public string $name;
}
`
	prog, acl := p.ParseString(src, "dto.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("load: %v", acl)
	}

	cls, _ := vm.GetOrLoadClass("SimpleDto")
	inst, _ := cls.GetValue(ctx)
	cv := inst.(*data.ClassValue)
	cv.SetProperty("name", data.NewStringValue(""))

	validatorCls, _ := vm.GetOrLoadClass("Validation\\Validator")
	validatorInst, _ := validatorCls.GetValue(ctx)
	validatorCV := validatorInst.(*data.ClassValue)
	method, _ := validatorCV.GetMethod("validate")
	fnCtx := validatorCV.CreateContext(method.GetVariables())
	fnCtx.SetVariableValue(method.GetVariables()[0], cv)
	result, acl := method.Call(fnCtx)
	if acl != nil {
		t.Fatalf("validate call: %v", acl)
	}
	arr, ok := result.(*data.ArrayValue)
	if !ok || len(arr.List) == 0 {
		t.Fatalf("expected violations array, got %#v", result)
	}
}
