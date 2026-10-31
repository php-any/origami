package container_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
)

func TestConstructorAnnotations(t *testing.T) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	defer vm.RunShutdownCallbacks()
	_, control := vm.LoadAndRun(filepath.Join("..", "..", "tests", "php", "container_constructor_annotations_test.php"))
	if control != nil {
		t.Fatalf("constructor annotation regression: %v", control)
	}
}

func TestInjectionRejectsNonConstructorParameters(t *testing.T) {
	for name, source := range map[string]string{
		"method":   `<?php class InvalidInjectMethod { public function handle(#[\Container\Annotation\Inject('probe')] object $value): void {} }`,
		"function": `<?php function invalidInjectFunction(#[\Container\Annotation\Named('probe')] object $value): void {}`,
		"closure":  `<?php class InvalidInjectClosure { public function __construct() { $fn = function (#[\Container\Annotation\Inject('probe')] object $value): void {}; } }`,
	} {
		t.Run(name, func(t *testing.T) {
			p := parser.NewParser()
			vm := runtime.NewVM(p).(*runtime.VM)
			std.Load(vm)
			php.Load(vm)
			defer vm.RunShutdownCallbacks()
			_, control := p.ParseString(source, name+".php")
			thrown, ok := control.(data.ThrowControl)
			if !ok || !strings.Contains(thrown.GetError().GetMessage(), "require a constructor parameter") {
				t.Fatalf("expected constructor target error, got %v", control)
			}
		})
	}
}
