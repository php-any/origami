package support

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func testVM(t *testing.T) data.VM {
	t.Helper()
	p := parser.NewParser()
	vm := runtime.NewVM(p)
	Load(vm)
	return vm
}

func TestStrKebabNative(t *testing.T) {
	vm := testVM(t)
	cls, ok := vm.GetClass("Illuminate\\Support\\Str")
	if !ok {
		t.Fatal("Str overlay missing")
	}
	gsm, ok := cls.(data.GetStaticMethod)
	if !ok {
		t.Fatal("Str is not static")
	}
	m, ok := gsm.GetStaticMethod("kebab")
	if !ok {
		t.Fatal("kebab missing")
	}
	ctx := vm.CreateContext(m.GetVariables())
	if ctl := m.GetVariables()[0].SetValue(ctx, data.NewStringValue("fullWidth")); ctl != nil {
		t.Fatal(ctl)
	}
	got, ctl := m.Call(ctx)
	if ctl != nil {
		t.Fatal(ctl)
	}
	if got.(data.Value).AsString() != "full-width" {
		t.Fatalf("got %v", got)
	}
}

func TestStrOrderedUuidNative(t *testing.T) {
	vm := testVM(t)
	cls, ok := vm.GetClass("Illuminate\\Support\\Str")
	if !ok {
		t.Fatal("Str missing")
	}
	m, ok := cls.(data.GetStaticMethod).GetStaticMethod("orderedUuid")
	if !ok {
		t.Fatal("orderedUuid missing")
	}
	got, ctl := m.Call(vm.CreateContext(m.GetVariables()))
	if ctl != nil {
		t.Fatal(ctl)
	}
	cv, ok := got.(*data.ClassValue)
	if !ok {
		t.Fatalf("want uuid object, got %T", got)
	}
	toStr, ok := cv.GetMethod("__toString")
	if !ok {
		t.Fatal("__toString missing")
	}
	s, ctl := toStr.Call(cv.CreateContext(toStr.GetVariables()))
	if ctl != nil {
		t.Fatal(ctl)
	}
	if !uuidRFC4122Re.MatchString(s.(data.Value).AsString()) {
		t.Fatalf("not rfc4122: %s", s.(data.Value).AsString())
	}
}

func TestStrParseCallbackNative(t *testing.T) {
	vm := testVM(t)
	cls, _ := vm.GetClass("Illuminate\\Support\\Str")
	m, ok := cls.(data.GetStaticMethod).GetStaticMethod("parseCallback")
	if !ok {
		t.Fatal("parseCallback missing")
	}
	ctx := vm.CreateContext(m.GetVariables())
	if ctl := m.GetVariables()[0].SetValue(ctx, data.NewStringValue("Foo@bar")); ctl != nil {
		t.Fatal(ctl)
	}
	got, ctl := m.Call(ctx)
	if ctl != nil {
		t.Fatal(ctl)
	}
	av := got.(*data.ArrayValue)
	vals := av.ToValueList()
	if len(vals) != 2 || vals[0].AsString() != "Foo" || vals[1].AsString() != "bar" {
		t.Fatalf("got %#v", vals)
	}
}

func TestStrPluralNative(t *testing.T) {
	if englishPlural("frame") != "frames" {
		t.Fatalf("frame -> %s", englishPlural("frame"))
	}
	vm := testVM(t)
	cls, _ := vm.GetClass("Illuminate\\Support\\Str")
	m, ok := cls.(data.GetStaticMethod).GetStaticMethod("plural")
	if !ok {
		t.Fatal("plural missing")
	}
	ctx := vm.CreateContext(m.GetVariables())
	_ = m.GetVariables()[0].SetValue(ctx, data.NewStringValue("frame"))
	got, ctl := m.Call(ctx)
	if ctl != nil {
		t.Fatal(ctl)
	}
	if got.(data.Value).AsString() != "frames" {
		t.Fatalf("got %s", got.(data.Value).AsString())
	}
}

func TestBagExtractPropNamesNative(t *testing.T) {
	vm := testVM(t)
	cls, ok := vm.GetClass("Illuminate\\View\\ComponentAttributeBag")
	if !ok {
		t.Fatal("bag overlay missing")
	}
	gsm := cls.(data.GetStaticMethod)
	m, ok := gsm.GetStaticMethod("extractPropNames")
	if !ok {
		t.Fatal("extractPropNames missing")
	}
	ctx := vm.CreateContext(m.GetVariables())
	keys := data.NewArrayValue(nil).(*data.ArrayValue)
	keys.SetStringKey("fullWidth", data.NewBoolValue(false))
	if ctl := m.GetVariables()[0].SetValue(ctx, keys); ctl != nil {
		t.Fatal(ctl)
	}
	got, ctl := m.Call(ctx)
	if ctl != nil {
		t.Fatal(ctl)
	}
	av := got.(*data.ArrayValue)
	vals := av.ToValueList()
	if len(vals) < 2 {
		t.Fatalf("len=%d", len(vals))
	}
	found := map[string]bool{}
	for _, v := range vals {
		found[v.AsString()] = true
	}
	if !found["fullWidth"] || !found["full-width"] {
		t.Fatalf("props=%v", found)
	}
}
