package runtime

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

var arrayReferenceBenchSink data.GetValue

func BenchmarkArrayIndexReference(b *testing.B) {
	vm := NewVM(parser.NewParser()).(*VM)
	variable := node.NewVariable(nil, "items", 0, nil)
	ctx := vm.CreateContext([]data.Variable{variable})
	array := data.NewArrayValue(nil).(*data.ArrayValue)
	for i := 0; i < 128; i++ {
		array.SetStringKey("key"+data.IntArrayKeyName(i), data.NewIntValue(i))
	}
	ctx.SetIndexZVal(0, data.NewZVal(array))
	expr := node.NewValueReference(nil, node.NewIndexExpression(nil, variable, data.NewStringValue("key127")))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, ctl := expr.GetValue(ctx)
		if ctl != nil {
			b.Fatal(ctl)
		}
		arrayReferenceBenchSink = value
	}
}

func BenchmarkArrayPackedKeys(b *testing.B) {
	value := data.NewIntValue(1)
	values := make([]data.Value, 128)
	for i := range values {
		values[i] = value
	}
	array := data.NewArrayValue(values).(*data.ArrayValue)
	_, _ = array.FindSlotByIntKey(127)
	b.Run("lookup", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			slot, _ := array.FindSlotByIntKey(127)
			arrayReferenceBenchSink = slot.ReadValue()
		}
	})
	b.Run("replace", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			array.SetIntKey(127, value)
		}
	})
	b.Run("append128", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			a := data.NewArrayValue(nil).(*data.ArrayValue)
			for n := 0; n < 128; n++ {
				a.AppendValue(value)
			}
			arrayReferenceBenchSink = a
		}
	})
}
