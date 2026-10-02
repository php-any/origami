package runtime_test

import (
	"strings"
	"testing"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
)

func TestStaticEnumDefaultUsesCurrentRequest(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	for i := 0; i < 2; i++ {
		request := runtime.NewRequestVM(base).(*runtime.RequestVM)
		var output strings.Builder
		request.SetOutputWriter(func(value string) { output.WriteString(value) })
		if _, ctl := request.LoadAndRun("../tests/php/static_enum_property_default_test.php"); ctl != nil {
			t.Fatalf("request %d: %s", i, ctl.AsString())
		}
		if got := output.String(); got != "static enum property default: OK\n" {
			t.Fatalf("request %d output = %q", i, got)
		}
		if _, found := base.GetClass("StaticEnumPropertyDefault\\LateAlignment"); found {
			t.Fatal("enum declaration escaped the request")
		}
	}
}
