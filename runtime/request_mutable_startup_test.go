package runtime_test

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"sync"
	"testing"
)

func TestStartupMutableConstantsAndIncludeResults(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	nested := data.NewArrayValue([]data.Value{data.NewIntValue(1)}).(*data.ArrayValue)
	parent := data.NewArrayValue([]data.Value{nested}).(*data.ArrayValue)
	base.SetConstant("STARTUP_ARRAY", parent)
	base.SetIncludeOnceResult("startup.php", parent)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			request := runtime.NewRequestVM(base).(*runtime.RequestVM)
			value, _ := request.GetConstant("STARTUP_ARRAY")
			result, _ := request.GetIncludeOnceResult("startup.php")
			if value != result {
				t.Error("startup alias lost")
			}
			child := value.(*data.ArrayValue).At(0).ReadValue().(*data.ArrayValue)
			child.SetIntKey(0, data.NewIntValue(9))
			if result.(*data.ArrayValue).At(0).ReadValue().(*data.ArrayValue).At(0).ReadValue().AsString() != "9" {
				t.Error("startup alias mutation lost")
			}
		}()
	}
	group.Wait()
	if nested.At(0).ReadValue().AsString() != "1" {
		t.Fatal("startup value mutated")
	}
}
