package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	illuminatehttp "github.com/php-any/origami/std/laravel/framework/illuminate/http"
)

func TestOfficialRequestCapture(t *testing.T) {
	base, _ := buildVM()
	if _, ctl := base.LoadAndRun("vendor/autoload.php"); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	vm := runtime.NewRequestVM(base).(*runtime.RequestVM)
	httpRequest := httptest.NewRequest("GET", "http://localhost/__runtime/order?item=one", nil)
	vm.BindHTTP(httpRequest, httptest.NewRecorder())
	ctx := vm.CreateContext(nil)
	request, ctl := illuminatehttp.NewIlluminateRequestValue(ctx, httpRequest)
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if _, ok := request.Class.(*node.ClassStatement); !ok {
		t.Fatalf("request replaced by %T", request.Class)
	}
	for method, want := range map[string]string{"getRequestUri": "/__runtime/order?item=one", "getPathInfo": "/__runtime/order", "path": "__runtime/order", "method": "GET", "host": "localhost"} {
		call := node.NewObjectMethod(nil, request, method, nil)
		got, ctl := call.GetValue(ctx)
		if ctl != nil {
			t.Errorf("%s: %s", method, ctl.AsString())
			continue
		}
		if got.(data.Value).AsString() != want {
			t.Errorf("%s=%q want %q", method, got.(data.Value).AsString(), want)
		}
	}
}

func TestOfficialRequestBodyAndHeaders(t *testing.T) {
	base, _ := buildVM()
	if _, ctl := base.LoadAndRun("vendor/autoload.php"); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	for _, fixture := range []struct {
		contentType, body, method string
		want                      data.Value
	}{
		{"application/json", `{"item":"json","nested":{"value":7}}`, "POST", data.NewStringValue("json")},
		{"application/x-www-form-urlencoded", "item=first&item=last&nested[value]=7", "POST", data.NewStringValue("last")},
	} {
		vm := runtime.NewRequestVM(base).(*runtime.RequestVM)
		raw := httptest.NewRequest(fixture.method, "https://example.test:8443/input?query=one", strings.NewReader(fixture.body))
		raw.Header.Set("Content-Type", fixture.contentType)
		raw.RemoteAddr = "127.0.0.1:4567"
		vm.BindHTTP(raw, httptest.NewRecorder())
		ctx := vm.CreateContext(nil)
		request, ctl := illuminatehttp.NewIlluminateRequestValue(ctx, raw)
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		for _, method := range []string{"getContent", "getContent"} {
			value, ctl := node.NewObjectMethod(nil, request, method, nil).GetValue(ctx)
			if ctl != nil || value.(data.Value).AsString() != fixture.body {
				t.Fatalf("%s body=%v control=%v", fixture.contentType, value, ctl)
			}
		}
		value, ctl := node.NewObjectMethod(nil, request, "input", []data.GetValue{data.NewStringValue("item")}).GetValue(ctx)
		if ctl != nil || value.(data.Value).AsString() != fixture.want.AsString() {
			t.Fatalf("%s input=%v control=%v", fixture.contentType, value, ctl)
		}
		server, ctl := node.NewServerVariable(nil).GetValue(ctx)
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		for name, want := range map[string]string{"SERVER_NAME": "example.test", "SERVER_PORT": "8443", "REMOTE_ADDR": "127.0.0.1", "REMOTE_PORT": "4567", "CONTENT_TYPE": fixture.contentType, "HTTPS": "on"} {
			slot, exists := server.(*data.ArrayValue).LookupZValByStringKey(name)
			if !exists || slot.ReadValue().AsString() != want {
				t.Errorf("%s metadata %s=%v want %s", fixture.contentType, name, slot, want)
			}
		}
	}
}
