package httpfoundation

import (
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
)

func TestJsonResponseEncodesNestedJsonSerializable(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	vm.AddClass(NewJsonResponseClass())

	src := `<?php
class JsonRespNested_Box implements JsonSerializable {
    public function jsonSerialize(): mixed {
        return [["id" => 1]];
    }
}
`
	prog, acl := p.ParseString(src, "json_resp_nested.php")
	if acl != nil {
		t.Fatalf("parse: %v", acl)
	}
	ctx := vm.CreateContext(nil)
	if _, acl = prog.GetValue(ctx); acl != nil {
		t.Fatalf("load class: %v", acl)
	}
	stmt, ok := vm.GetClass("JsonRespNested_Box")
	if !ok {
		t.Fatal("JsonRespNested_Box not registered")
	}
	box := data.NewClassValue(stmt, ctx)
	payload := &data.ArrayValue{List: []*data.ZVal{
		data.NewNamedZVal("entries", box),
		data.NewNamedZVal("status", data.NewStringValue("enabled")),
	}}
	resp := data.NewClassValue(NewJsonResponseClass(), ctx)
	if _, ctl := jsonResponseSetDataWith(resp, payload); ctl != nil {
		t.Fatalf("setData: %v", ctl)
	}
	content, _ := resp.GetProperty("content")
	if content == nil {
		t.Fatal("missing content")
	}
	s := content.AsString()
	if strings.Contains(s, "Object(") {
		t.Fatalf("encoded Object dump: %s", s)
	}
	if !strings.Contains(s, `"entries"`) || !strings.Contains(s, `"id"`) || !strings.Contains(s, `"status"`) {
		t.Fatalf("unexpected json: %s", s)
	}
}
