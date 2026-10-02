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
)

// Execute generated constructors so loss of declaration metadata changes an
// observable PHP result instead of merely checking emitted strings.
func TestGeneratedStrictDeclarations(t *testing.T) {
	cases := []struct{ name, source, expected string }{
		{"function", "declare(strict_types=1); function value(): string { return 42; } return value();", "TypeError"},
		{"method", "declare(strict_types=1); class Sample { function value(): string { return 42; } } return (new Sample())->value();", "TypeError"},
		{"closure", "declare(strict_types=1); $f = static function &(): string { return 42; }; return $f();", "TypeError"},
		{"weak_closure", "$f = static function (): string { return 42; }; return $f();", "*data.StringValue:42"},
	}
	var generated strings.Builder
	var calls strings.Builder
	for _, test := range cases {
		p := parser.NewParser()
		runtime.NewVM(p)
		program, ctl := p.ParseString(test.source, test.name+".php")
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
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
)
func run(build func() (data.GetValue, []data.Variable), expected string) {
 vm := runtime.NewVM(parser.NewParser())
 program, vars := build()
 value, ctl := program.GetValue(vm.CreateContext(vars))
 actual := ""
 if failure, ok := ctl.(*data.ThrowValue); ok { actual = failure.GetName() } else if ctl != nil { panic(ctl.AsString()) } else { actual = fmt.Sprintf("%T:%s", value, value.(data.Value).AsString()) }
 if actual != expected { panic(fmt.Sprintf("got %s, want %s", actual, expected)) }
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
