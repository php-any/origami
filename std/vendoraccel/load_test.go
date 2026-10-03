package vendoraccel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/laravel/serve"
	"testing"
)

func TestRegistrationGateRejectsUncertifiedLoaderAtomically(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	g := &registrationGate{VM: vm}
	serve.Load(g)
	if ctl := g.AddClass(&data.StdClass{}); ctl == nil {
		t.Fatal("uncertified class admitted")
	}
	if ctl := g.commit(); ctl == nil {
		t.Fatal("ignored registration failure committed")
	}
	for _, name := range []string{"stdClass", "Illuminate\\Foundation\\Console\\ServeCommand"} {
		if _, found := vm.GetClass(name); found {
			t.Fatalf("partial registration of %s", name)
		}
	}
}

func TestRegistrationGateRejectsOtherDeclarationKinds(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	for _, register := range []func(*registrationGate) data.Control{
		func(g *registrationGate) data.Control {
			return g.AddInterface(node.NewInterfaceStatement(nil, "Uncertified", nil, nil))
		},
		func(g *registrationGate) data.Control {
			class := node.NewClassStatement(nil, "Illuminate\\Foundation\\Console\\ServeCommand", "", nil, nil, nil)
			return g.AddClass(class)
		},
		func(g *registrationGate) data.Control {
			return g.AddFunc(node.NewFunctionStatement(nil, "uncertified", nil, nil, nil, nil, false))
		},
		func(g *registrationGate) data.Control { return g.RegisterFunction("uncertified", func() {}) },
		func(g *registrationGate) data.Control { return g.RegisterReflectClass("uncertified", struct{}{}) },
	} {
		g := &registrationGate{VM: vm}
		if register(g) == nil || g.commit() == nil {
			t.Fatal("uncertified registration committed")
		}
	}
	Load(vm)
	if _, found := vm.GetClass("Illuminate\\Foundation\\Console\\ServeCommand"); !found {
		t.Fatal("host adapter missing")
	}
}
