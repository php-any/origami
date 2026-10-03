package main

import (
	"context"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
)

func TestInferenceLeavesExecutableDeclarationsUnchanged(t *testing.T) {
	variable := node.NewVariable(nil, "item", 0, nil)
	first := node.NewBinaryAssign(nil, variable, &node.NewExpression{ClassName: "First"})
	second := node.NewBinaryAssign(nil, variable, &node.NewExpression{ClassName: "Second"})
	document := &DocumentInfo{}
	left := NewLspContext(context.Background(), nil)
	right := NewLspContext(context.Background(), nil)
	visit := func(*LspContext, data.GetValue, data.GetValue) bool { return true }
	document.foreachNode(left, first, nil, visit)
	document.foreachNode(left, second, nil, visit)
	document.foreachNode(right, second, nil, visit)
	if variable.Type != data.TypeInvalid || variable.GetType() != nil {
		t.Fatal("analysis changed the executable variable declaration")
	}
	inferred, ok := left.GetVariableType("$item").(*data.LspTypes)
	if !ok || len(inferred.Types) != 2 {
		t.Fatalf("left inference = %#v", left.GetVariableType("item"))
	}
	if got := getClassNameFromType(right.GetVariableType("item")); got != "Second" {
		t.Fatalf("document inference leaked: %q", got)
	}
}

func TestCompletionReadsCompactDeclaration(t *testing.T) {
	variable := node.NewVariable(nil, "item", 0, data.NewDeclaredType("?Example\\Service"))
	typeOf := getTypeFromLeftNode(variable, nil, nil, "", defines.Position{})
	if got := getClassNameFromType(typeOf); got != "Example\\Service" {
		t.Fatalf("compact nullable type = %T %v", typeOf, typeOf)
	}
}

func TestLspScopeKeepsUnspecifiedRuntimeType(t *testing.T) {
	scope := NewLspScope(nil, "global", "global", "types.php")
	if typ := scope.AddVariable("item", nil, nil).GetType(); typ != nil {
		t.Fatalf("scope injected analysis as a declaration: %v", typ)
	}
}
