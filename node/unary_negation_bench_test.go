package node

import (
	"github.com/php-any/origami/data"
	"testing"
)

var unaryBenchSink data.GetValue

func BenchmarkUnaryIntegerNegation(b *testing.B) {
	expression := NewUnaryExpression(nil, "-", data.NewIntValue(42))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		value, ctl := expression.GetValue(nil)
		if ctl != nil {
			b.Fatal(ctl.AsString())
		}
		unaryBenchSink = value
	}
}
