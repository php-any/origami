package compile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
)

// Execute generated constructors so loss of declaration metadata changes an
// observable PHP result instead of merely checking emitted strings.
func TestGeneratedStrictDeclarations(t *testing.T) {
	cases := []struct{ name, source, expected string }{
		{"function", "declare(strict_types=1); function value(): string { return 42; } return value();", "TypeError"},
		{"method", "declare(strict_types=1); class Sample { function value(): string { return 42; } } return (new Sample())->value();", "TypeError"},
		{"closure", "declare(strict_types=1); $f = static function &(): string { return 42; }; return $f();", "TypeError"},
		{"weak_closure", "$f = static function (): string { return 42; }; return $f();", "*data.StringValue:42"},
		{"compound_closure", "$f = static function (): int|float { return '2.5'; }; return $f();", "*data.FloatValue:2.5"},
		{"decoded_string", `$f = static function (): string { return "'quoted'\nline"; }; return $f();`, "*data.StringValue:'quoted'\nline"},
		{"typed_property", "class Sample { public string $value; } $object = new Sample(); $object->value = 42; return $object->value;", "*data.StringValue:42"},
		{"final_class", "final class Sample {} return (new ReflectionClass('Sample'))->isFinal();", "*data.BoolValue:true"},
		{"readonly_class", "readonly class Sample { public int $value; } return (new ReflectionClass('Sample'))->isReadOnly();", "*data.BoolValue:true"},
		{"final_method", "class Sample { final public function value(): int { return 1; } } return (new ReflectionMethod('Sample', 'value'))->isFinal();", "*data.BoolValue:true"},
		{"reference_method", "class Sample { public static string $value = 'before'; public function &value(): string { return self::$value; } } $object = new Sample(); $value =& $object->value(); $value = 'after'; return Sample::$value;", "*data.StringValue:after"},
		{"static_interface", "interface Contract { public static function &value(): string; } $method = new ReflectionMethod('Contract', 'value'); return $method->isStatic() && $method->returnsReference();", "*data.BoolValue:true"},
		{"abstract_method", "abstract class Sample { abstract public function value(): int; } return (new ReflectionMethod('Sample', 'value'))->isAbstract();", "*data.BoolValue:true"},
		{"class_constant", "class Sample { const LABEL = 'class'; } return Sample::LABEL;", "*data.StringValue:class"},
		{"interface_constant", "interface Contract { const LABEL = 'interface'; } return Contract::LABEL;", "*data.StringValue:interface"},
		{"enum", "enum Sample: int { case One = 1; case Two = 2; } return Sample::from(2) === Sample::cases()[1];", "*data.BoolValue:true"},
		{"function_attribute", "#[Attribute] class Marker { public function __construct(public string $label) {} } #[Marker(label: 'ok')] function value() {} return (new ReflectionFunction('value'))->getAttributes()[0]->newInstance()->label;", "*data.StringValue:ok"},
	}
	var generated strings.Builder
	var calls strings.Builder
	for _, test := range cases {
		p := parser.NewParser()
		base := runtime.NewVM(p).(*runtime.VM)
		php.Load(base)
		program, ctl := p.ParseString(test.source, test.name+".php")
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		augmentProgramASTFromBase(program, base, test.name+".php")
		g := NewGenerator()
		code, err := g.Generate(ParsedFile{Path: test.name + ".php", Program: program, Variables: p.GetVariables()})
		if err != nil {
			t.Fatal(err)
		}
		generated.WriteString(code)
		fmt.Fprintf(&calls, "run(%s, %q)\n", g.funcNameForPath(test.name+".php"), test.expected)
	}
	source := `package main
import (
 "fmt"
 "github.com/php-any/origami/data"
 "github.com/php-any/origami/node"
 "github.com/php-any/origami/parser"
 "github.com/php-any/origami/runtime"
 "github.com/php-any/origami/std/php"
 "github.com/php-any/origami/std/php/core"
 "github.com/php-any/origami/std/php/attribute"
)
var _ = core.BackedEnumCasesMethod{}
var _ = attribute.NewAttributeClass
func run(build func() (data.GetValue, []data.Variable), expected string) {
 vm := runtime.NewVM(parser.NewParser())
 php.Load(vm)
 program, vars := build()
 value, ctl := program.GetValue(vm.CreateContext(vars))
 actual := ""
 if failure, ok := ctl.(*data.ThrowValue); ok { actual = failure.GetName() } else if ctl != nil { panic(ctl.AsString()) } else { actual = fmt.Sprintf("%T:%s", value, value.(data.Value).AsString()) }
 if actual != expected { panic(fmt.Sprintf("got %s, want %s; control %v", actual, expected, ctl)) }
}
` + generated.String() + "func main() {\n" + calls.String() + "}\n"
	path := filepath.Join(t.TempDir(), "generated.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "run", path)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated execution: %v\n%s\n%s", err, output, source)
	}
}
