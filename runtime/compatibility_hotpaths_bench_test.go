package runtime

import (
	"github.com/php-any/origami/parser"
	"testing"
)

func BenchmarkCompatibilityHotPaths(b *testing.B) {
	cases := []struct{ name, source string }{
		{"property_read", "class Sample { public int $number = 42; } $object = new Sample(); $object->number;"},
		{"property_write", "class Sample { public int $number = 42; } $object = new Sample(); $object->number = 43;"},
		{"property_reference", "class Sample { public int $number = 42; } $object = new Sample(); $ref =& $object->number;"},
		{"property_coalesce", "class Sample { public int $number = 42; } $object = new Sample(); $object->number ?? 0;"},
		{"property_isset", "class Sample { public int $number = 42; } $object = new Sample(); isset($object->number);"},
		{"foreach_value", "$items = [1,2,3,4,5,6,7,8]; foreach ($items as $value) { $copy = $value; }"},
		{"counter_loop_8", "$copy = 0; for ($i = 0; $i < 8; $i++) { $copy = $i; }"},
	}
	for _, test := range cases {
		b.Run(test.name, func(b *testing.B) {
			p := parser.NewParser()
			vm := NewVM(p)
			program, ctl := p.ParseString(test.source, "compat_bench.php")
			if ctl != nil {
				b.Fatal(ctl.AsString())
			}
			ctx := vm.CreateContext(p.GetVariables())
			if _, ctl := program.GetValue(ctx); ctl != nil {
				b.Fatal(ctl.AsString())
			}
			operation := program.Statements[len(program.Statements)-1]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				value, ctl := operation.GetValue(ctx)
				if ctl != nil {
					b.Fatal(ctl.AsString())
				}
				typedCallBenchSink = value
			}
		})
	}
}
