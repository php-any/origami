package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestLimitsAndContextCancellation(t *testing.T) {
	for _, scenario := range []string{"size", "idle", "context", "write"} {
		t.Run(scenario, func(t *testing.T) {
			result := make(chan string, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, e := (&ws.Upgrader{}).Upgrade(w, r, nil)
				if e != nil {
					result <- e.Error()
					return
				}
				defer conn.Close()
				base := runtime.NewVM(parser.NewParser())
				vm := runtime.NewRequestVM(base).(*runtime.RequestVM)
				restore := vm.Activate()
				defer restore()
				defer vm.RunShutdownCallbacks()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				vm.BindContext(ctx)
				class := NewConnClassFrom(conn)
				m, _ := class.GetMethod("configure")
				scope := vm.CreateContext(m.GetVariables())
				scope.SetIndexZVal(0, data.NewZVal(data.NewArrayValueFromSlots([]*data.ZVal{data.NewNamedZVal("maxMessageBytes", data.NewIntValue(16)), data.NewNamedZVal("readTimeoutMs", data.NewIntValue(50))})))
				if _, ctl := m.Call(scope); ctl != nil {
					result <- ctl.AsString()
					return
				}
				if scenario == "context" {
					cancel()
					conn.SetReadDeadline(time.Now().Add(time.Second))
				}
				if scenario == "write" {
					m, _ = class.GetMethod("writeText")
					scope = vm.CreateContext(m.GetVariables())
					scope.SetIndexZVal(0, data.NewZVal(data.NewStringValue(strings.Repeat("x", 17))))
				} else {
					m, _ = class.GetMethod("readText")
					scope = vm.CreateContext(m.GetVariables())
				}
				_, ctl := m.Call(scope)
				if ctl == nil {
					result <- "operation unexpectedly succeeded"
				} else {
					result <- "OK"
				}
			}))
			defer server.Close()
			client, _, e := ws.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
			if e != nil {
				t.Fatal(e)
			}
			defer client.Close()
			if scenario == "size" {
				if e := client.WriteMessage(ws.TextMessage, []byte(strings.Repeat("x", 17))); e != nil {
					t.Fatal(e)
				}
			}
			select {
			case got := <-result:
				if got != "OK" {
					t.Fatal(got)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("connection did not terminate")
			}
		})
	}
}
