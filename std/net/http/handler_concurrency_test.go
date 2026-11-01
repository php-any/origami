package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

// Every request writes its PHP globals before any request can read them back.
// Run with -race to catch writes to obsolete process-global fallback caches too.
func TestConcurrentLanguageHandlersKeepRequestState(t *testing.T) {
	directory := t.TempDir()
	setup := filepath.Join(directory, "setup.php")
	read := filepath.Join(directory, "read.php")
	for path, source := range map[string]string{
		setup: `<?php $GLOBALS['token']=$_GET['token']; $_SESSION['token']=$_GET['token'];`,
		read:  `<?php return [$_GET['token'],$_POST['form'],$_REQUEST['form'],$_COOKIE['token'],$GLOBALS['token'],$_SESSION['token']];`,
	} {
		if err := os.WriteFile(path, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, hot := range []bool{false, true} {
		name := "normal"
		if hot {
			name = "hot"
		}
		t.Run(name, func(t *testing.T) {
			const count = 32
			base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
			root := base.CreateContext(nil)
			variables := []data.Variable{data.NewVariable("request", 0, nil), data.NewVariable("response", 1, nil)}
			var seen sync.Map
			var shutdowns atomic.Int32
			arrived := make(chan struct{}, count)
			release := make(chan struct{})
			middleware := &requestVMCallback{variables: append(append([]data.Variable(nil), variables...), data.NewVariable("next", 2, nil)), call: func(ctx data.Context) {
				vm := ctx.GetVM().(*runtime.RequestVM)
				token := vm.HTTPRequest().Header.Get("X-Request-ID")
				seen.Store(token, vm)
				vm.EnsureGlobalZVal("middleware_token").StoreRaw(data.NewStringValue(token))
				vm.AddShutdownCallback(data.NewFuncValue(&requestVMCallback{call: func(data.Context) { shutdowns.Add(1) }}))
				next, _ := ctx.GetIndexValue(2)
				call := vm.CreateContext(next.(*data.FuncValue).Value.GetVariables())
				for index := 0; index < 2; index++ {
					value, _ := ctx.GetIndexValue(index)
					call.SetIndexZVal(index, data.NewZVal(value))
				}
				if _, control := next.(*data.FuncValue).Call(call); control != nil {
					t.Error(control.AsString())
				}
			}}
			callback := &requestVMCallback{variables: variables, call: func(ctx data.Context) {
				vm := ctx.GetVM().(*runtime.RequestVM)
				token := vm.HTTPRequest().Header.Get("X-Request-ID")
				previous, found := seen.Load(token)
				if !found || previous != vm {
					t.Error("middleware and handler must share one RequestVM")
				}
				if vm.EnsureGlobalZVal("middleware_token").ReadValue().AsString() != token {
					t.Error("middleware state was lost")
				}
				if _, control := vm.LoadAndRun(setup); control != nil {
					t.Error(control.AsString())
				}
				arrived <- struct{}{}
				<-release
				result, control := vm.LoadAndRun(read)
				if control != nil {
					t.Error(control.AsString())
					return
				}
				values, ok := result.(*data.ArrayValue)
				if !ok || values.Len() != 6 {
					t.Errorf("unexpected PHP result: %T", result)
					return
				}
				for _, value := range values.Range() {
					if actual := value.ReadValue().AsString(); actual != token {
						t.Errorf("PHP state leaked: got %q, want %q", actual, token)
					}
				}
			}}
			wrap, err := newMiddleware(middleware, root)
			if err != nil {
				t.Fatal(err)
			}
			var handler http.Handler = Handler{Value: callback, Ctx: root}
			if hot {
				handler = HotHandler{Value: callback, Ctx: root}
			}
			handler = wrap(handler)
			var group sync.WaitGroup
			for i := 0; i < count; i++ {
				group.Add(1)
				go func(i int) {
					defer group.Done()
					token := fmt.Sprint(i)
					request := httptest.NewRequest("POST", "/?token="+token, strings.NewReader("form="+token))
					request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					request.Header.Set("Cookie", "token="+token)
					request.Header.Set("X-Request-ID", token)
					handler.ServeHTTP(httptest.NewRecorder(), request)
				}(i)
			}
			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()
			for i := 0; i < count; i++ {
				select {
				case <-arrived:
				case <-timer.C:
					close(release)
					group.Wait()
					t.Fatal("concurrent requests did not reach handler")
				}
			}
			close(release)
			group.Wait()
			if got := shutdowns.Load(); got != count {
				t.Fatalf("shutdown callbacks=%d, want %d", got, count)
			}
			if base.EnsureGlobalZVal("token").ReadValue().AsString() != "" || base.EnsureGlobalZVal("middleware_token").ReadValue().AsString() != "" {
				t.Fatal("request state escaped to the startup VM")
			}
		})
	}
}

func BenchmarkLanguageHTTPHandler(b *testing.B) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	variables := []data.Variable{data.NewVariable("request", 0, nil), data.NewVariable("response", 1, nil)}
	callback := &requestVMCallback{variables: variables, call: func(data.Context) {}}
	for _, hot := range []bool{false, true} {
		name := "normal"
		var handler http.Handler = Handler{Value: callback, Ctx: base.CreateContext(nil)}
		if hot {
			name = "hot"
			handler = HotHandler{Value: callback, Ctx: base.CreateContext(nil)}
		}
		b.Run(name, func(b *testing.B) {
			request := httptest.NewRequest("GET", "/", nil)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				handler.ServeHTTP(httptest.NewRecorder(), request)
			}
		})
	}
}
