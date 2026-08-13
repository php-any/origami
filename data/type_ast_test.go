package data

import "testing"

func TestASTType(t *testing.T) {
	ty := AST{}
	if !ty.Is(NewNullValue()) {
		t.Fatal("AST should accept null")
	}
	if !ty.Is(NewAnyValue(nil)) {
		t.Fatal("AST should accept AnyValue")
	}
	if !ty.Is(NewASTValue(nil, nil)) {
		t.Fatal("AST should accept ASTValue")
	}
	if ty.Is(NewStringValue("x")) {
		t.Fatal("AST should reject string")
	}
	if ty.String() != "AstNode" {
		t.Fatalf("AST.String() = %q, want AstNode", ty.String())
	}
}
