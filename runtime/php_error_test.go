package runtime

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"sync"
	"testing"
)

func TestPHPErrorStateRequestIsolation(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	base.PHPErrorState().SetReporting(123)
	base.PHPErrorState().SetLast(&data.PHPErrorInfo{Message: "startup"})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			request := NewRequestVM(base).(*RequestVM)
			state := request.PHPErrorState()
			if state.Reporting() != 123 || state.Last() != nil {
				t.Error("startup diagnostic leaked or reporting lost")
			}
			state.SetReporting(i)
			state.SetLast(&data.PHPErrorInfo{Type: i, Message: "request"})
			if state.Reporting() != i || state.Last().Type != i {
				t.Error("request diagnostic changed")
			}
		}(i)
	}
	wg.Wait()
	if base.PHPErrorState().Reporting() != 123 || base.PHPErrorState().Last().Message != "startup" {
		t.Fatal("request mutated worker")
	}
}
