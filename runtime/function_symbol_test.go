package runtime

import (
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

func TestFunctionSymbolsUseExecutingRequest(t *testing.T) {
	p := parser.NewParser()
	base := NewVM(p)
	program, ctl := p.ParseString("return vAlUe();", "shared_function_call.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func(value int) {
			defer group.Done()
			request := NewRequestVM(base)
			function := node.NewFunctionStatement(nil, "Value", nil, []data.GetValue{node.NewReturnStatement(nil, data.NewIntValue(value))}, nil, nil, false)
			if ctl := request.AddFunc(function); ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			result, ctl := program.GetValue(request.CreateContext(nil))
			if ctl != nil || result.(*data.IntValue).Value != value {
				t.Errorf("shared call returned %v / %v; want %d", result, ctl, value)
			}
		}(i)
	}
	group.Wait()
	if _, found := base.GetFunc("value"); found {
		t.Fatal("request function was published into the base VM")
	}
}
