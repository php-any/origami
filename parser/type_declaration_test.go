package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"testing"
)

func TestCompactTypeDeclarations(t *testing.T) {
	program, ctl := NewParser().ParseString("function declaration(int|string $value): mixed { return $value; }", "types.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	function := program.Statements[0].(*node.FunctionStatement)
	if function.Ret != data.TypeMixed {
		t.Fatal("explicit mixed confused with no declaration")
	}
	parameter := function.Params[0].(data.Parameter)
	ref, ok := parameter.GetType().(data.TypeRef)
	if !ok || ref.Kind() != data.TypeKindUnion || ref.Members().Len() != 2 {
		t.Fatalf("parameter has legacy type %T", parameter.GetType())
	}
}
func TestReturnDeclarationCompileFailures(t *testing.T) {
	for _, source := range []string{
		"function f(): void { return null; }",
		"function f(): void { if (false) { return 1; } }",
		"function f(): never { return; }",
		"function f(): mixed { return; }",
		"function f(): ?string { return; }",
		"$f = fn():void => 1;",
		"function f(): int, string { return [1, 'value']; }",
		"function f(array<string> $value) {}",
		"function f(): array<string> { return []; }",
	} {
		t.Run(source, func(t *testing.T) {
			_, ctl := NewParser().ParseString(source, "invalid.php")
			failure, ok := ctl.(*data.ThrowValue)
			if !ok || !failure.PHPCompileFatal {
				t.Fatalf("expected compile fatal: %v", ctl)
			}
		})
	}
	for _, source := range []string{
		"function f() { return; }",
		"function f(): void { return; }",
		"function f(): void { $inner = function(): int { return 1; }; }",
		"function f(): never { throw null; }",
		"function f(): true { return true; }",
	} {
		if _, ctl := NewParser().ParseString(source, "valid.php"); ctl != nil {
			t.Fatal(ctl.AsString())
		}
	}
}

// Resolved call targets are runtime links, not nested declarations. Traversing
// them validates another function's returns against the caller and can cycle.
func TestReturnValidationIgnoresResolvedCallTargets(t *testing.T) {
	target := node.NewFunctionStatement(nil, "target", nil, []data.GetValue{&node.ReturnStatement{Value: data.NewIntValue(1)}}, nil, data.TypeInt, false)
	call := &node.CallExpression{Fun: target}
	target.Body = append(target.Body, call)
	if ctl := validateReturnBody([]data.GetValue{call}, data.TypeVoid); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if ctl := validateReturnDeclarations(&node.Program{Statements: []data.GetValue{target}}); ctl != nil {
		t.Fatal(ctl.AsString())
	}
}
