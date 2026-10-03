package data

import (
	"fmt"
	"testing"
)

var arrayCreationSink Value
var thisCreationSink *ThisValue

func BenchmarkSmallArrayCreation(b *testing.B) {
	for _, size := range []int{0, 1, 3, 8, 32} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			values := make([]Value, size)
			for i := range values {
				values[i] = NewIntValue(i)
			}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				arrayCreationSink = NewArrayValue(values)
			}
		})
	}
}

func BenchmarkThisIdentity(b *testing.B) {
	object := &ClassValue{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		thisCreationSink = NewThisValue(object)
	}
}
