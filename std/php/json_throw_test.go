package php

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func BenchmarkJSONThrowBoundary(b *testing.B) {
	vm := runtime.NewVM(parser.NewParser())
	Load(vm)
	ctx := vm.CreateContext(NewJsonDecodeFunction().GetVariables())
	ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue(`{"name":"Ada","times":2}`)))
	ctx.SetIndexZVal(1, data.NewZVal(data.NewBoolValue(true)))
	ctx.SetIndexZVal(2, data.NewZVal(data.NewIntValue(512)))
	ctx.SetIndexZVal(3, data.NewZVal(data.NewIntValue(0)))
	fn := NewJsonDecodeFunction()
	if _, control := fn.Call(ctx); control != nil {
		b.Fatal(control.AsString())
	}
	value := data.NewArrayValue([]data.Value{data.NewStringValue("Ada"), data.NewIntValue(2)})
	b.Run("decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = fn.Call(ctx)
		}
	})
	b.Run("encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = JsonEncode(ctx, value)
		}
	})
}
