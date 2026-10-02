package runtime_test

import (
	"strings"
	"testing"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
)

func TestAutoloadLeadingSeparatorInRequestVM(t *testing.T) {
	runtime.ClearAutoLoad()
	t.Cleanup(runtime.ClearAutoLoad)
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	// Repeat on the same AST cache with fresh request declarations.
	for i := 0; i < 2; i++ {
		request := runtime.NewRequestVM(base).(*runtime.RequestVM)
		var output strings.Builder
		request.SetOutputWriter(func(s string) { output.WriteString(s) })
		if _, ctl := request.LoadAndRun("../tests/php/autoload_leading_separator_test.php"); ctl != nil {
			t.Fatalf("request %d: %s", i, ctl.AsString())
		}
		if got := output.String(); got != "autoload leading separator OK\n" {
			t.Fatalf("request %d output = %q", i, got)
		}
		if _, found := base.GetClass("AutoloadLeadingSeparator\\Admin"); found {
			t.Fatal("autoloaded class escaped the request")
		}
		if _, found := base.GetInterface("AutoloadLeadingSeparator\\Contract"); found {
			t.Fatal("autoloaded interface escaped the request")
		}
		if callbacks := runtime.GetAutoLoad(); len(callbacks) != 0 {
			t.Fatal("autoload callback was not unregistered")
		}
	}
}
