package reflection

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestUserAnnotationArgumentsMixedPositionalNamed(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	ctx := vm.CreateContext(nil)
	decl := node.NewAnnotation(nil, "Mark", []data.GetValue{
		node.NewNamedArgument(nil, "count", node.NewIntLiteral(nil, "7")),
	})
	decl.Arguments = append([]data.GetValue{node.NewStringLiteral(nil, "mix")}, decl.Arguments...)

	arr, ctl := userAnnotationArguments(ctx, decl)
	if ctl != nil {
		t.Fatalf("userAnnotationArguments: %v", ctl)
	}
	av := arr.(*data.ArrayValue)
	if slot, ok := av.LookupZValByStringKey("count"); !ok || slot.ReadValue().AsString() != "7" {
		t.Fatalf("count=%#v", slot)
	}
	if av.Len() < 1 || av.At(0).ReadValue().AsString() != "mix" {
		t.Fatalf("positional=%v", av.ToValueList())
	}
}

func TestClassStatementStoresAnnotationDeclaration(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	decl := node.NewAnnotation(nil, "Table", []data.GetValue{
		node.NewNamedArgument(nil, "name", node.NewStringLiteral(nil, "users")),
	})
	inst := data.NewClassValue(&node.ClassStatement{Name: "Table"}, vm.CreateContext(nil))
	cls := &node.ClassStatement{Name: "App\\Entity"}
	cls.AddAnnotations(inst, decl)
	if len(cls.AnnotationDecls) != 1 || cls.AnnotationDecls[0] != decl {
		t.Fatalf("AnnotationDecls not stored: %#v", cls.AnnotationDecls)
	}
	got := node.AnnotationDeclAt(cls.AnnotationDecls, 0)
	if got != decl {
		t.Fatalf("AnnotationDeclAt mismatch")
	}
}
