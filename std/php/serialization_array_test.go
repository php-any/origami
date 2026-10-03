package php

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestUnserializeStringByteLengths(t *testing.T) {
	function := NewUnserializeFunction()
	vm := runtime.NewVM(parser.NewParser())
	for _, test := range []struct{ raw, expected string }{
		{`s:0:"";`, ""},
		{`s:3:"a"b";`, "a\"b"},
		{"s:3:\"a\x00b\";", "a\x00b"},
		{`s:6:"中文";`, "中文"},
		{`s:14:"__origami_o:{}";`, "__origami_o:{}"},
	} {
		ctx := vm.CreateContext(function.GetVariables())
		ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue(test.raw)))
		value, ctl := function.Call(ctx)
		text, ok := value.(*data.StringValue)
		if ctl != nil || !ok || text.Value != test.expected {
			t.Fatalf("%q = %v, %v; want %q", test.raw, value, ctl, test.expected)
		}
	}
	for _, raw := range []string{`s:9:"x";`, `s:-1:"x";`, `s:1:"xy";`, `s:99999999999999999999999:"x";`} {
		ctx := vm.CreateContext(function.GetVariables())
		ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue(raw)))
		value, ctl := function.Call(ctx)
		failed, ok := value.(*data.BoolValue)
		if ctl != nil || !ok || failed.Value {
			t.Fatalf("accepted invalid byte length: %q", raw)
		}
	}
}
