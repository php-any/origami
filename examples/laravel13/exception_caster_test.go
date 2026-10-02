package main

import (
	"testing"

	"github.com/php-any/origami/data"
)

func TestSymfonyExceptionCaster(t *testing.T) {
	vm, _ := buildVM()
	value, ctl := vm.LoadAndRun("tests/origami/exception_caster.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if result, ok := value.(*data.BoolValue); !ok || !result.Value {
		t.Fatalf("caster regression returned %v", value)
	}
}
