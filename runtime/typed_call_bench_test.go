package runtime

import (
	"testing"
	"unsafe"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

var typedCallBenchSink data.GetValue

func BenchmarkTypedCall(b *testing.B) {
	cases := []struct{ name, source string }{
		{"untyped", "function value($number) { return $number; } value(42);"},
		{"exact_int", "function value(int $number): int { return $number; } value(42);"},
		{"weak_string", "function value(string $number): string { return $number; } value(42);"},
		{"method_exact", "class Sample { function value(int $number): int { return $number; } } $object = new Sample(); $object->value(42);"},
	}
	for _, test := range cases {
		b.Run(test.name, func(b *testing.B) {
			p := parser.NewParser()
			vm := NewVM(p)
			program, ctl := p.ParseString(test.source, "typed_bench.php")
			if ctl != nil {
				b.Fatal(ctl.AsString())
			}
			ctx := vm.CreateContext(p.GetVariables())
			if _, ctl := program.GetValue(ctx); ctl != nil {
				b.Fatal(ctl.AsString())
			}
			call := program.Statements[len(program.Statements)-1]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				value, ctl := call.GetValue(ctx)
				if ctl != nil {
					b.Fatal(ctl.AsString())
				}
				typedCallBenchSink = value
			}
		})
	}
}

func TestTypeModeStorageSize(t *testing.T) {
	if unsafe.Sizeof(uintptr(0)) != 8 {
		t.Skip("amd64 storage comparison")
	}
	t.Logf("Context=%d TokenFrom=%d", unsafe.Sizeof(Context{}), unsafe.Sizeof(node.TokenFrom{}))
	if unsafe.Sizeof(Context{}) > 152 || unsafe.Sizeof(node.TokenFrom{}) > 56 {
		t.Fatal("type mode expanded every runtime frame or AST source location")
	}
}
