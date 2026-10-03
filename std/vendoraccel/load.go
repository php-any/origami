package vendoraccel

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/serve"
)

// Load registers only the process/HTTP host adapter. Uncertified vendor
// replacements are never executed or registered: Composer's PHP declarations
// remain authoritative, including files autoload and Reflection metadata.
// A replacement may be enabled only after its versioned differential contract
// proves signatures, visibility, Reflection and serialization compatibility.
func Load(vm data.VM) {
	gate := &registrationGate{VM: vm}
	serve.Load(gate)
	if ctl := gate.commit(); ctl != nil {
		vm.ThrowControl(ctl)
	}
}

// Registration is staged: even a loader that ignores AddClass's control cannot
// partially install an uncertified replacement. Only the concrete SAPI adapter
// is admitted. Adding a name to a loader is insufficient to enable a vendor class.
type registrationGate struct {
	data.VM
	host    data.ClassStmt
	failure data.Control
}

func (g *registrationGate) reject(name string) data.Control {
	ctl := data.NewErrorThrow(nil, fmt.Errorf("uncertified native vendor registration: %s", name))
	if g.failure == nil {
		g.failure = ctl
	}
	return ctl
}
func (g *registrationGate) AddClass(class data.ClassStmt) data.Control {
	if _, ok := class.(*serve.ServeCommandClass); !ok || class.GetName() != "Illuminate\\Foundation\\Console\\ServeCommand" || g.host != nil {
		return g.reject(class.GetName())
	}
	g.host = class
	return nil
}
func (g *registrationGate) AddInterface(class data.InterfaceStmt) data.Control {
	return g.reject(class.GetName())
}
func (g *registrationGate) AddFunc(fn data.FuncStmt) data.Control { return g.reject(fn.GetName()) }
func (g *registrationGate) RegisterFunction(name string, fn interface{}) data.Control {
	return g.reject(name)
}
func (g *registrationGate) RegisterReflectClass(name string, instance interface{}) data.Control {
	return g.reject(name)
}
func (g *registrationGate) commit() data.Control {
	if g.failure != nil {
		return g.failure
	}
	if g.host != nil {
		return g.VM.AddClass(g.host)
	}
	return nil
}
