package runtime

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

func TestObjectHandleOutlivesPooledMethodContext(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	class := node.NewClassStatement(nil, "WorkerObject", "", nil, nil, nil)
	object := data.NewClassValue(class, vm.CreateContext(nil))
	frame := object.CreateContext(nil).(*data.ClassMethodContext)
	borrowed := frame.ClassValue
	frame.ReleasePooled()
	// Simulate a native fluent method returning its frame's ClassValue. Creating
	// the next method must use the object's stable context, even if the pool hands
	// back the exact same frame which the previous invocation used.
	for i := 0; i < 100; i++ {
		next := borrowed.CreateContext(nil).(*data.ClassMethodContext)
		if next.GetVM() != vm || next.InstanceIdentity() != object {
			t.Fatal("method handle retained a recycled context")
		}
		next.ReleasePooled()
	}
	request := object.CloneRequestScoped(vm.CreateContext(nil))
	if request.InstanceIdentity() != request || request.ObjectValue == object.ObjectValue || request.Class != object.Class {
		t.Fatal("request proxy did not preserve metadata and distinct instance identity")
	}
	if request.CreateContext(nil).(*data.ClassMethodContext).InstanceIdentity() != request {
		t.Fatal("request method escaped to the original instance")
	}
}

func TestRequestProxyPreservesNominalTypes(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	root := node.NewInterfaceStatement(nil, "RootContract", nil, nil)
	leaf := node.NewInterfaceStatement(nil, "LeafContract", []string{"RootContract"}, nil)
	if ctl := base.AddInterface(root); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if ctl := base.AddInterface(leaf); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	class := node.NewClassStatement(nil, "WorkerObject", "", []string{"LeafContract"}, nil, nil)
	object := data.NewClassValue(class, base.CreateContext(nil))
	first := object.CloneRequestScoped(NewRequestVM(base).CreateContext(nil))
	second := object.CloneRequestScoped(NewRequestVM(base).CreateContext(nil))
	for _, value := range []*data.ClassValue{object, first, second} {
		for _, name := range []string{"workerobject", "ROOTCONTRACT", "\\leafcontract"} {
			if !(data.Class{Name: name}).Is(value) {
				t.Fatalf("proxy lost nominal type %s", name)
			}
		}
	}
	if first.InstanceIdentity() == second.InstanceIdentity() || first.Class != second.Class {
		t.Fatal("proxies must share type metadata but have separate identities")
	}
	request := NewRequestVM(base)
	if ctl := request.AddInterface(node.NewInterfaceStatement(nil, "RequestContract", []string{"RootContract"}, nil)); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if _, ok := request.GetInterface("\\REQUESTCONTRACT"); !ok {
		t.Fatal("request interface lookup is case sensitive")
	}
	if _, ok := base.GetInterface("RequestContract"); ok {
		t.Fatal("request declaration leaked into base")
	}
}

func BenchmarkObjectMethodContext(b *testing.B) {
	vm := NewVM(parser.NewParser()).(*VM)
	class := node.NewClassStatement(nil, "WorkerObject", "", nil, nil, nil)
	object := data.NewClassValue(class, vm.CreateContext(nil))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		object.CreateContext(nil).(*data.ClassMethodContext).ReleasePooled()
	}
}

func BenchmarkRequestObjectScope(b *testing.B) {
	vm := NewVM(parser.NewParser()).(*VM)
	ctx := vm.CreateContext(nil)
	class := node.NewClassStatement(nil, "WorkerObject", "", nil, nil, nil)
	object := data.NewClassValue(class, ctx)
	for i := 0; i < 40; i++ {
		_ = object.ObjectValue.SetProperty(string(rune('A'+i)), data.NewArrayValue([]data.Value{data.NewIntValue(i)}))
	}
	b.Run("overlay", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			object.CloneRequestScoped(ctx)
		}
	})
	b.Run("deep-clone", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			object.CloneSandbox(ctx)
		}
	})
}

func TestMethodFramePreservesRequestVM(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	request := NewRequestVM(base)
	class := node.NewClassStatement(nil, "WorkerObject", "", nil, nil, nil)
	object := data.NewClassValue(class, base.CreateContext(nil))
	frame := data.WrapMethodFrame(request.CreateContext(nil), object, class, class)
	if frame.GetVM() != request || frame.Context.GetVM() != request {
		t.Fatal("method frame escaped to worker VM")
	}
	if frame.InstanceIdentity() != object {
		t.Fatal("caller frame changed PHP object identity")
	}
}
