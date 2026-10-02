package main

import (
	"testing"

	"github.com/php-any/origami/data"
)

func TestBladeIconAttributes(t *testing.T) {
	vm, _ := buildVM()
	value, ctl := vm.LoadAndRun("tests/origami/blade_icon_attributes.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if result, ok := value.(*data.BoolValue); !ok || !result.Value {
		t.Fatalf("icon regression returned %v", value)
	}
}
