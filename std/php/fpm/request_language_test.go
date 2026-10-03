package fpm_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/fpm"
)

func TestSharedProgramHasRequestOwnedDeclarationsAndGlobals(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	file := filepath.Join(t.TempDir(), "entry.php")
	source := `<?php
class RequestOnlyClass {}
function requestOnlyFunction() { global $hits; return ++$hits; }
define('REQUEST_ONLY_CONSTANT', 42);
$hits = 0;
return requestOnlyFunction();`
	if err := os.WriteFile(file, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	check := func() {
		request := fpm.New(base, nil)
		if _, found := request.GetClass("RequestOnlyClass"); found {
			t.Error("class available before request file ran")
		}
		value, ctl := request.LoadAndRun(file)
		if ctl != nil {
			t.Error(ctl.AsString())
			return
		}
		if got := value.(data.Value).AsString(); got != "1" {
			t.Errorf("hits=%s", got)
		}
		if got := request.EnsureGlobalZVal("hits").ReadValue().AsString(); got != "1" {
			t.Errorf("global hits=%s", got)
		}
		if _, found := request.GetClass("RequestOnlyClass"); !found {
			t.Error("request class missing")
		}
		if _, found := request.GetFunc("requestOnlyFunction"); !found {
			t.Error("request function missing")
		}
		if value, found := request.GetConstant("REQUEST_ONLY_CONSTANT"); !found || value.AsString() != "42" {
			t.Error("request constant missing")
		}
	}
	check()
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() { defer group.Done(); check() }()
	}
	group.Wait()
	if _, found := base.GetClass("RequestOnlyClass"); found {
		t.Fatal("request class published into VM")
	}
	if _, found := base.GetFunc("requestOnlyFunction"); found {
		t.Fatal("request function published into VM")
	}
	if _, found := base.GetConstant("REQUEST_ONLY_CONSTANT"); found {
		t.Fatal("request constant published into VM")
	}
}

func TestRequestFileAndCompiledExecutionStateIsolated(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	first := runtime.NewRequestVM(base).(*runtime.RequestVM)
	second := fpm.New(base, nil)
	file := filepath.Join(t.TempDir(), "compiled.php")
	base.RegisterCompiledFile(file, func() (data.GetValue, []data.Variable) { return data.NewIntValue(7), nil })
	for _, request := range []*runtime.RequestVM{first, second} {
		value, ctl := request.RunCompiledFile(file)
		if ctl != nil || value.(data.Value).AsString() != "7" {
			t.Fatalf("compiled result=%v control=%v", value, ctl)
		}
		if value, ctl := request.RunCompiledFile(file); value != nil || ctl != nil {
			t.Fatal("compiled file ran twice")
		}
	}
	if base.GetPhpFileCache(file) {
		t.Fatal("request include state escaped to VM")
	}
	first.SetIncludeOnceResult("one.php", data.NewIntValue(99))
	if _, found := second.GetIncludeOnceResult("one.php"); found {
		t.Fatal("include result escaped request")
	}
	if _, found := base.GetIncludeOnceResult("one.php"); found {
		t.Fatal("include result escaped to VM")
	}
	base.SetConstant("PHP_SAPI", data.NewStringValue("cli"))
	if ctl := first.SetConstant("PHP_SAPI", data.NewStringValue("cli-server")); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if value, _ := second.GetConstant("PHP_SAPI"); value.AsString() != "cli" {
		t.Fatal("SAPI constant escaped request")
	}
}

func TestRequestTemplateAndCompileLoadStayInRequest(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	file := filepath.Join(t.TempDir(), "template.php")
	if err := os.WriteFile(file, []byte("<?php class TemplateOnlyClass {} define('TEMPLATE_ONLY', $label); return $label;"), 0600); err != nil {
		t.Fatal(err)
	}
	first, second := fpm.New(base, nil), fpm.New(base, nil)
	props := data.NewArrayValue(nil).(*data.ArrayValue)
	props.SetStringKey("label", data.NewStringValue("template"))
	value, ctl := first.ParseFile(file, props)
	if ctl != nil || value.AsString() != "template" {
		t.Fatalf("template=%v control=%v", value, ctl)
	}
	if _, ok := first.GetClass("TemplateOnlyClass"); !ok {
		t.Fatal("template declaration missing")
	}
	if _, ok := base.GetClass("TemplateOnlyClass"); ok {
		t.Fatal("template declaration escaped request")
	}
	if _, ok := base.GetConstant("TEMPLATE_ONLY"); ok {
		t.Fatal("template constant escaped request")
	}
	if _, ok := second.GetConstant("TEMPLATE_ONLY"); ok {
		t.Fatal("template constant escaped to second request")
	}
	if ctl := second.CompileLoad(file); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if _, ok := second.GetClass("TemplateOnlyClass"); !ok {
		t.Fatal("compile-only request declaration missing")
	}
	if !second.GetPhpFileCache(file) || base.GetPhpFileCache(file) {
		t.Fatal("compile-load file state escaped request")
	}
	if _, ok := second.GetConstant("TEMPLATE_ONLY"); ok {
		t.Fatal("compile-load executed top-level PHP")
	}
}
