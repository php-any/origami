package runtime

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"testing"
)

func TestTickStateRequestIsolation(t *testing.T) {
	base := NewVM(parser.NewParser())
	first := NewRequestVM(base).(data.TickProvider)
	second := NewRequestVM(base).(data.TickProvider)
	first.PHPTickState().Count = 7
	if second.PHPTickState().Count != 0 || base.(data.TickProvider).PHPTickState().Count != 0 {
		t.Fatal("tick counters leaked between requests")
	}
}
