package runtime

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

type executionContextKey struct{}
type executionContextProbe struct{}

func (*executionContextProbe) GetName() string               { return "execution_context_probe" }
func (*executionContextProbe) GetParams() []data.GetValue    { return nil }
func (*executionContextProbe) GetVariables() []data.Variable { return nil }
func (*executionContextProbe) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GoContext().Value(executionContextKey{}).(string)
	return data.NewStringValue(value), nil
}

func TestExecutionContextInheritance(t *testing.T) {
	vm := NewVM(parser.NewParser())
	vm.AddFunc(&executionContextProbe{})
	parent := vm.CreateContext(nil).(*Context)
	goContext, cancel := context.WithCancel(context.WithValue(context.Background(), executionContextKey{}, "bound"))
	defer cancel()
	parent.SetGoContext(goContext)
	for _, child := range []data.Context{parent.CreateContext(nil), parent.CreateBaseContext()} {
		if child.GoContext().Value(executionContextKey{}) != "bound" {
			t.Fatal("child did not inherit execution context")
		}
		data.ReleaseContext(child)
	}
	file := filepath.Join(t.TempDir(), "included.php")
	if err := os.WriteFile(file, []byte("<?php return execution_context_probe();"), 0600); err != nil {
		t.Fatal(err)
	}
	value, ctl := vm.LoadInCallerContext(parent, file)
	if returned, ok := ctl.(data.ReturnControl); ok {
		value, ctl = returned.ReturnValue(), nil
	}
	if ctl != nil || value.(data.AsString).AsString() != "bound" {
		t.Fatalf("include lost context: %v %v", value, ctl)
	}
	cancel()
	if parent.CreateContext(nil).GoContext().Err() != context.Canceled {
		t.Fatal("child lost cancellation")
	}
	parent.ReleasePooled()
	if vm.CreateContext(nil).GoContext().Value(executionContextKey{}) != nil {
		t.Fatal("pooled frame retained execution values")
	}
}
