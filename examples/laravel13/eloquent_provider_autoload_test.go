package main

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

func TestEloquentProviderRequestAutoload(t *testing.T) {
	vm, _ := buildVM()
	if _, ctl := vm.LoadAndRun("tests/origami/runtime_bootstrap.php"); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if _, found := vm.GetClass("App\\Models\\Admin"); found {
		t.Fatal("Admin must be unloaded to exercise request autoload")
	}
	request := runtime.NewRequestVM(vm)
	value, ctl := request.LoadAndRun("tests/origami/eloquent_provider_autoload.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	model, ok := value.(*data.ClassValue)
	if !ok || model.GetName() != "App\\Models\\Admin" {
		t.Fatalf("model = %v", value)
	}
}
