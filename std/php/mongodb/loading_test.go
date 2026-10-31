package mongodb_test

import (
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/mongodb"
)

func newHostVM() *runtime.VM {
	vm := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	return vm
}

func TestMongoDBExplicitRegistration(t *testing.T) {
	vm := newHostVM()
	check := func(vm *runtime.VM, want bool) {
		t.Helper()
		for _, name := range []string{"MongoDB\\Driver\\Manager", "MongoDB\\BSON\\ObjectId"} {
			if _, ok := vm.GetClass(name); ok != want {
				t.Fatalf("class %s registered=%v, want %v", name, ok, want)
			}
		}
		if _, ok := vm.GetFunc("MongoDB\\BSON\\fromPHP"); ok != want {
			t.Fatalf("BSON functions registered=%v, want %v", ok, want)
		}
		if _, ok := vm.GetInterface("MongoDB\\BSON\\Type"); ok != want {
			t.Fatalf("BSON interface registered=%v, want %v", ok, want)
		}
		if _, ok := vm.GetConstant("MONGODB_VERSION"); ok != want {
			t.Fatalf("version constant registered=%v, want %v", ok, want)
		}
		fn, _ := vm.GetFunc("extension_loaded")
		ctx := vm.CreateContext(fn.GetVariables())
		ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue("mongodb")))
		v, ctl := fn.Call(ctx)
		if ctl != nil || v.(*data.BoolValue).Value != want {
			t.Fatalf("extension_loaded: %v, %v; want %v", v, ctl, want)
		}
	}
	check(vm, false)
	mongodb.Load(vm)
	check(vm, true)
	check(newHostVM(), false)
}

func TestMongoDBPHPRegression(t *testing.T) {
	vm := newHostVM()
	mongodb.Load(vm)
	path, err := filepath.Abs("../../../tests/php/mongodb_bson_driver_test.php")
	if err != nil {
		t.Fatal(err)
	}
	if _, ctl := vm.LoadAndRun(path); ctl != nil {
		t.Fatalf("PHP regression: %v", ctl)
	}
}
