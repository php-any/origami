package runtime

import (
	"github.com/php-any/origami/parser"
	"strings"
	"testing"
)

func TestClassDeclarationConstraints(t *testing.T) {
	cases := []struct{ source, error string }{
		{`final class Base {} class Child extends Base {}`, "cannot extend final"},
		{`class Base { final public function value() {} } class Child extends Base { public function value() {} }`, "Cannot override final method"},
		{`class Base { final public static function value() {} } class Child extends Base { public static function value() {} }`, "Cannot override final method"},
		{`readonly class Base {} class Child extends Base {}`, "Readonly and non-readonly"},
		{`class Base {} readonly class Child extends Base {}`, "Readonly and non-readonly"},
		{`class Invalid { public readonly $value; }`, "must have a type"},
		{`class Invalid { public readonly int $value = 1; }`, "cannot have default"},
		{`class Invalid { public static readonly int $value; }`, "cannot be readonly"},
	}
	for _, test := range cases {
		t.Run(test.error+test.source[:10], func(t *testing.T) {
			p := parser.NewParser()
			vm := NewVM(p)
			program, ctl := p.ParseString(test.source, "declaration.php")
			if ctl == nil {
				_, ctl = program.GetValue(vm.CreateContext(p.GetVariables()))
			}
			if ctl == nil || !strings.Contains(ctl.AsString(), test.error) {
				t.Fatalf("expected %q, got %v", test.error, ctl)
			}
		})
	}
}
