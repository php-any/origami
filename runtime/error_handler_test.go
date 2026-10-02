package runtime

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"sync"
	"testing"
)

func TestErrorHandlerRequestStacksIsolated(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	startup := data.NewStringValue("startup")
	base.SetErrorHandler(startup)
	for _, http := range []bool{false, true} {
		var group sync.WaitGroup
		for i := 0; i < 16; i++ {
			group.Add(1)
			go func() {
				defer group.Done()
				if http {
					defer BeginRequestOutput()()
				}
				request := NewRequestVM(base).(*RequestVM)
				local := data.NewStringValue("request")
				var old data.Value
				if http {
					old = base.SetErrorHandler(local)
				} else {
					old = request.SetErrorHandler(local)
				}
				if old != startup || request.GetErrorHandler() != local {
					t.Error("error handler escaped request")
				}
				request.RestoreErrorHandler()
				if request.GetErrorHandler() != startup {
					t.Error("startup error handler not restored")
				}
				request.RestoreErrorHandler()
				if request.GetErrorHandler() != nil || !request.RestoreErrorHandler() {
					t.Error("empty error handler restore failed")
				}
			}()
		}
		group.Wait()
	}
	if base.GetErrorHandler() != startup {
		t.Fatal("VM handler changed")
	}
}
