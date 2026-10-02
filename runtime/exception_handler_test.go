package runtime

import (
	"errors"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestExceptionHandlerRequestIsolation(t *testing.T) {
	for _, activeHTTP := range []bool{false, true} {
		t.Run(map[bool]string{false: "RequestVM", true: "HTTPBaseRegistration"}[activeHTTP], func(t *testing.T) {
			base := NewVM(parser.NewParser()).(*VM)
			startup := data.NewFuncValue(&shutdownTestCallback{fn: func(data.Context) {}})
			base.SetExceptionHandler(startup)
			var group sync.WaitGroup
			for i := 0; i < 16; i++ {
				group.Add(1)
				go func() {
					defer group.Done()
					if activeHTTP {
						defer BeginRequestOutput()()
					}
					request := NewRequestVM(base).(*RequestVM)
					calls := 0
					callback := data.NewFuncValue(&shutdownTestCallback{fn: func(ctx data.Context) {
						if ctx.GetVM() != request {
							t.Error("handler used worker VM")
						}
						calls++
					}})
					var old data.Value
					if activeHTTP {
						old = base.SetExceptionHandler(callback)
					} else {
						old = request.SetExceptionHandler(callback)
					}
					if old != startup || request.GetExceptionHandler() != callback {
						t.Error("registration escaped request")
					}
					thrown := data.NewErrorThrow(nil, errors.New("request failure"))
					if handled, ctl := request.HandleUnhandledException(thrown); !handled || ctl != nil || calls != 1 {
						t.Errorf("handled=%v control=%v calls=%d", handled, ctl, calls)
					}
					request.SetExceptionHandler(nil)
					request.RestoreExceptionHandler()
					if request.GetExceptionHandler() != callback {
						t.Error("did not restore request callback")
					}
					request.RestoreExceptionHandler()
					if request.GetExceptionHandler() != startup {
						t.Error("did not restore startup callback")
					}
					request.RestoreExceptionHandler()
					if request.GetExceptionHandler() != nil || !request.RestoreExceptionHandler() {
						t.Error("empty restore failed")
					}
				}()
			}
			group.Wait()
			if base.GetExceptionHandler() != startup {
				t.Fatal("worker handler changed")
			}
		})
	}
}

func TestExceptionHandlerBoundaryDoesNotReenter(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	thrown := data.NewErrorThrow(nil, errors.New("failure"))
	calls := 0
	base.SetExceptionHandler(data.NewFuncValue(&shutdownTestCallback{fn: func(data.Context) {
		calls++
		if handled, ctl := base.HandleUnhandledException(thrown); handled || ctl != thrown {
			t.Error("handler reentered itself")
		}
	}}))
	compileFatal := data.NewErrorThrow(nil, errors.New("compile fatal")).(*data.ThrowValue)
	compileFatal.PHPCompileFatal = true
	if handled, ctl := base.HandleUnhandledException(compileFatal); handled || ctl != compileFatal {
		t.Fatal("compile fatal reached handler")
	}
	for i := 0; i < 2; i++ {
		if handled, ctl := base.HandleUnhandledException(thrown); !handled || ctl != nil {
			t.Fatal("handler failed")
		}
	}
	if calls != 2 {
		t.Fatalf("calls=%d", calls)
	}
}
