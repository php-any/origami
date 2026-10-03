package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/laravel/httpkernel"
	"github.com/php-any/origami/std/laravel/serve"
)

func TestLoginRequestIsolation(t *testing.T) {
	vm, _ := buildVM()
	value, ctl := vm.LoadAndRun("tests/origami/runtime_bootstrap.php")
	if ctl != nil {
		if thrown, ok := ctl.(*data.ThrowValue); ok {
			t.Fatalf("%s\n%+v", ctl.AsString(), thrown.StackFrames)
		}
		t.Fatal(ctl.AsString())
	}
	app := value.(*data.ClassValue)
	handler, err := serve.NewHTTPHandler(vm, app)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	client := &http.Client{Timeout: 35 * time.Second}
	livewireScript := regexp.MustCompile(`<script\s+src="[^"]*/livewire-[0-9a-f]{8}/livewire(?:\.min)?\.js(?:\?[^"]*)?"`)
	check := func() {
		t.Helper()
		response, err := client.Get(server.URL + "/admin/login")
		if err != nil {
			t.Error(err)
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil || response.StatusCode != 200 {
			t.Errorf("login status=%d body=%s error=%v", response.StatusCode, body[:min(len(body), 600)], err)
			return
		}
		if !bytes.Contains(body, []byte("wire:snapshot=")) || !livewireScript.Match(body) {
			t.Errorf("login assets incomplete: snapshot=%v script=%v bytes=%d", bytes.Contains(body, []byte("wire:snapshot=")), livewireScript.Match(body), len(body))
		}
	}
	check()
	func() {
		restore := runtime.BeginRequestOutput()
		defer restore()
		kernel, ctl := httpkernel.Resolve(app)
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		scoped := httpkernel.Sandbox(runtime.NewRequestVM(vm).CreateContext(nil), kernel)
		value, ctl := scoped.GetProperty("app")
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		request, ok := value.(*data.ClassValue)
		if !ok {
			t.Fatalf("sandbox app %T", value)
		}
		if request.ObjectValue == app.ObjectValue {
			t.Fatal("sandbox reused worker app")
		}
		raw, _ := request.GetProperty("bindings")
		bindings := raw.(*data.ArrayValue)
		if slot, ok := bindings.LookupZValByStringKey("Illuminate\\Foundation\\Mix"); ok {
			entry := slot.ReadValue().(*data.ArrayValue)
			concrete, _ := entry.LookupZValByStringKey("concrete")
			closure := concrete.ReadValue().(*data.FuncValue).Value.(data.RequestClosureBinder)
			for _, owner := range closure.RequestScopeObjects() {
				if owner.GetName() == request.GetName() && owner.ObjectValue != request.ObjectValue {
					t.Fatal("factory retained the worker application")
				}
			}
		}
	}()
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() { defer group.Done(); check() }()
	}
	group.Wait()
}
