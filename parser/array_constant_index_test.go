package parser

import (
	"testing"

	"github.com/php-any/origami/node"
)

func TestConstantIndexBeforeAssignment(t *testing.T) {
	p := NewParser()
	program, ctl := p.ParseString("<?php $a[PHP_INT_MAX] = 'max';", "constant-index.php")
	if ctl != nil {
		t.Fatal(ctl)
	}
	assign, ok := program.Statements[0].(*node.BinaryAssign)
	if !ok {
		t.Fatalf("constant array-key assignment parsed as %T", program.Statements[0])
	}
	index, ok := assign.Left.(*node.IndexExpression)
	if !ok {
		t.Fatalf("assignment target parsed as %T", assign.Left)
	}
	constant, ok := index.Index.(*node.ConstantName)
	if !ok || constant.Name != "PHP_INT_MAX" {
		t.Fatalf("array key was consumed as a declaration: %T", index.Index)
	}
}
