package runtime_test

import (
	"github.com/php-any/origami/data"
	"strings"
	"sync"
	"testing"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
)

func TestAutoloadLeadingSeparatorInRequestVM(t *testing.T) {
	runtime.ClearAutoLoad()
	t.Cleanup(runtime.ClearAutoLoad)
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	// Repeat on the same AST cache with fresh request declarations.
	for i := 0; i < 2; i++ {
		request := runtime.NewRequestVM(base).(*runtime.RequestVM)
		var output strings.Builder
		request.SetOutputWriter(func(s string) { output.WriteString(s) })
		if _, ctl := request.LoadAndRun("../tests/php/autoload_leading_separator_test.php"); ctl != nil {
			t.Fatalf("request %d: %s", i, ctl.AsString())
		}
		if got := output.String(); got != "autoload leading separator OK\n" {
			t.Fatalf("request %d output = %q", i, got)
		}
		if _, found := base.GetClass("AutoloadLeadingSeparator\\Admin"); found {
			t.Fatal("autoloaded class escaped the request")
		}
		if _, found := base.GetInterface("AutoloadLeadingSeparator\\Contract"); found {
			t.Fatal("autoloaded interface escaped the request")
		}
		if callbacks := runtime.GetAutoLoad(); len(callbacks) != 0 {
			t.Fatal("autoload callback was not unregistered")
		}
	}
}

func TestDefaultAutoloadAndIniRequestIsolation(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			defer runtime.BeginRequestOutput()()
			request := runtime.NewRequestVM(base).(*runtime.RequestVM)
			request.SetOutputWriter(func(string) {})
			if _, ctl := request.LoadAndRun("../tests/php/default_autoload_test.php"); ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			request.SetAutoloadExtensions(".request")
			request.StorePHPIni("precision", "3")
			if value, found := request.LookupPHPIni("precision"); !found || value != "3" {
				t.Error("request ini write was lost")
			}
		}()
	}
	group.Wait()
	if _, found := base.GetClass("DefaultAutoload\\CaseDefaultLoader"); found {
		t.Fatal("default autoload class leaked to worker")
	}
	if base.AutoloadExtensions() != ".inc,.php" {
		t.Fatal("autoload extensions leaked to worker")
	}
	if _, found := base.LookupPHPIni("include_path"); found {
		t.Fatal("include_path leaked to worker")
	}
	if _, found := base.LookupPHPIni("precision"); found {
		t.Fatal("precision leaked to worker")
	}
}

func TestCachedDeclarationTraitsAndConstructorsAreRequestOwned(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func(marker int) {
			defer group.Done()
			defer runtime.BeginRequestOutput()()
			request := runtime.NewRequestVM(base).(*runtime.RequestVM)
			ctx := request.CreateContext(nil)
			ctx.SetVariableByName("request_marker", data.NewIntValue(marker))
			value, ctl := request.LoadInCallerContext(ctx, "../tests/php/request_cached_declaration_test.php")
			if ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			object, ok := value.(*data.ClassValue)
			if !ok {
				t.Errorf("cached declaration returned %T", value)
				return
			}
			declaration, _ := request.GetClass("CachedScopeChild")
			if object.Class != declaration {
				t.Error("constructor retained a different request declaration")
			}
			value, ctl = object.GetProperty("marker")
			if ctl != nil || value.(*data.IntValue).Value != marker {
				t.Error("inherited constructor used another request")
			}
		}(i + 1)
	}
	group.Wait()
	for _, name := range []string{"CachedScopeParent", "CachedScopeChild", "CachedScopeTrait"} {
		if _, found := base.GetClass(name); found {
			t.Fatalf("%s leaked into worker", name)
		}
	}
}

func TestRetainedClosureScopesUnregisteredReceiver(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	value, ctl := base.LoadAndRun("../tests/php/request_closure_receiver_test.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	callback := value.(*data.FuncValue)
	owner := callback.Value.(data.RequestClosureBinder).RequestScopeObjects()[0]
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			defer runtime.BeginRequestOutput()()
			request := runtime.NewRequestVM(base).(*runtime.RequestVM)
			bound := request.RequestObjectScope().Bind(callback).(*data.FuncValue)
			result, ctl := bound.Value.Call(request.CreateContext(bound.Value.GetVariables()))
			if ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			if result.(*data.IntValue).Value != 1 {
				t.Error("closure receiver escaped request scope")
			}
		}()
	}
	group.Wait()
	count, ctl := owner.GetProperty("count")
	if ctl != nil || count.(*data.IntValue).Value != 0 {
		t.Fatal("worker receiver was modified")
	}
}

func TestAutoloadRegistryRequestIdentityAndIsolation(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	value, ctl := base.LoadAndRun("../tests/php/autoload_registry_test.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	parent := value.(*data.ClassValue)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			defer runtime.BeginRequestOutput()()
			request := runtime.NewRequestVM(base).(*runtime.RequestVM)
			ctx := request.CreateContext(nil)
			originals := request.AutoloadOriginals()
			if len(originals) != 1 {
				t.Errorf("callbacks %d", len(originals))
				return
			}
			if loaded, ctl := runtime.CallAutoLoad("RequestRegistryMissing", ctx); loaded || ctl != nil {
				t.Error("missing autoload result", ctl)
				return
			}
			array := originals[0].(*data.ArrayValue)
			owner := array.At(0).ReadValue().(*data.ClassValue)
			if owner.ObjectValue == parent.ObjectValue {
				t.Error("autoload receiver retained parent")
			}
			names, _ := owner.GetProperty("names")
			if names.(*data.ArrayValue).Len() != 2 {
				t.Error("autoload callback writes did not reach original callable receiver")
			}
			if !request.UnregisterAutoload(originals[0]) || len(request.AutoloadCallbacks()) != 0 {
				t.Error("request unregister failed")
			}
		}()
	}
	group.Wait()
	names, _ := parent.GetProperty("names")
	if names.(*data.ArrayValue).Len() != 1 || len(base.AutoloadCallbacks()) != 1 {
		t.Fatal("request autoload state escaped")
	}
}
