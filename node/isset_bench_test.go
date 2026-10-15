package node

import (
	"testing"

	"github.com/php-any/origami/data"
)

func BenchmarkIssetArray(b *testing.B) {
	for _, named := range []bool{false, true} {
		name := "packed"
		var key data.Value = data.NewIntValue(0)
		if named {
			name = "named"
			key = data.NewStringValue("present")
		}
		b.Run(name, func(b *testing.B) {
			array := data.NewArrayValue(nil).(*data.ArrayValue)
			array.SetKey(key, data.NewIntValue(1))
			expr := NewIssetStatement(nil, []data.GetValue{NewIndexExpression(nil, array, key)})
			value, ctl := expr.GetValue(nil)
			if ctl != nil || !value.(*data.BoolValue).Value {
				b.Fatal("isset failed")
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = expr.GetValue(nil)
			}
		})
	}
}
