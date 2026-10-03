package runtime_test

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestDeclarationDNFAndCompileFailures(t *testing.T) {
	valid := []string{
		"interface A {} interface B {} interface C {} function f((A&B)|C|null $x): (A&B)|C|null { return $x; }",
		"interface A {} interface B {} class C { public A&B $x; public function f((A&B)|null &$x): A&B { return $x; } }",
		"interface A {} interface B {} interface C { public function f(A&B $x): (A&B)|null; }",
		"trait A { public function f(self $x): parent { return $x; } }",
		"class A { public function f(self $x): self { return $x; } } function f() { return; }",
		"class A {} class B extends A { public parent $x; public function f(): parent { return $this; } }",
	}
	invalid := []string{
		"class A { public const int VALUE = '1'; }", "class A { public const string VALUE = 1; }",
		"function f(int|int $x) {}", "function f(bool|false $x) {}", "function f(true|false $x) {}",
		"function f(mixed|int $x) {}", "function f(?mixed $x) {}", "function f(?null $x) {}",
		"function f(void $x) {}", "function f(never $x) {}", "function f(int&A $x) {}",
		"function f(self&A $x) {}", "function f(A&B|C $x) {}", "function f(A|B&C $x) {}",
		"function f((A&B) $x) {}", "function f((A|B)&C $x) {}", "function f((A&B)|A $x) {}",
		"function f((A&B)|(B&A) $x) {}", "function f((A&B&C)|(A&B) $x) {}",
		"function f(object|A $x) {}", "function f(iterable|array $x) {}", "function f(iterable|Traversable $x) {}",
		"function f(self $x) {}", "function f(): static { return null; }", "class A { public callable $x; }",
		"class A { public function f(static $x) {} }", "class A { public function f(): parent { return $this; } }",
		"class A {} function f(): self { return null; }",
		"class A { public function __construct(public callable $x) {} }",
	}
	for _, source := range valid {
		t.Run("valid "+source, func(t *testing.T) {
			p := parser.NewParser()
			vm := runtime.NewVM(p)
			p.SetVM(vm)
			if _, ctl := p.ParseString(source, "declaration.php"); ctl != nil {
				t.Fatal(ctl.AsString())
			}
		})
	}
	for _, source := range invalid {
		t.Run("invalid "+source, func(t *testing.T) {
			p := parser.NewParser()
			vm := runtime.NewVM(p)
			p.SetVM(vm)
			_, ctl := p.ParseString(source, "declaration.php")
			if failure, ok := ctl.(*data.ThrowValue); !ok || !failure.PHPCompileFatal {
				t.Fatalf("expected compile fatal, got %v", ctl)
			}
		})
	}
}
