package fpm_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/fpm"
)

func TestExceptionHandlerCallableBoundary(t *testing.T) {
	cases := []struct{ name, callback, expected string }{
		{"function", "'namedHandler'", "function:failure"},
		{"instance", "[$object, 'instanceHandler']", "instance:failure"},
		{"static_array", "[1 => 'staticHandler', 0 => 'HandlerFixture']", "static:failure"},
		{"static_string", "'HandlerFixture::staticHandler'", "static:failure"},
		{"invokable", "$object", "invoke:failure"},
		{"bound_closure", "(function (Throwable $e) { echo $this->name . ':' . $e->getMessage(); })->bindTo($object, 'HandlerFixture')", "bound:failure"},
		{"no_parameters", "function () { echo 'no-parameters'; }", "no-parameters"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
			php.Load(base)
			file := filepath.Join(t.TempDir(), "handler.php")
			source := `<?php
function namedHandler(Throwable $e) { echo 'function:' . $e->getMessage(); }
class HandlerFixture {
 public $name = 'bound';
 public function instanceHandler(Throwable $e) { echo 'instance:' . $e->getMessage(); }
 public static function staticHandler(Throwable ...$errors) { echo 'static:' . $errors[0]->getMessage(); }
 public function __invoke(Throwable $e) { echo 'invoke:' . $e->getMessage(); }
}
$object = new HandlerFixture();
set_exception_handler(` + test.callback + `);
try { throw new Exception('caught'); } catch (Exception $e) {}
throw new Exception('failure');`
			if err := os.WriteFile(file, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			var out strings.Builder
			request := fpm.New(base, func(s string) { out.WriteString(s) })
			_, thrown := request.LoadAndRun(file)
			if thrown == nil || out.Len() != 0 {
				t.Fatal("caught/uncaught exception dispatched inside script")
			}
			if handled, ctl := request.HandleUnhandledException(thrown); !handled || ctl != nil {
				t.Fatalf("handled=%v control=%v", handled, ctl)
			}
			if out.String() != test.expected {
				t.Fatalf("got %q, want %q", out.String(), test.expected)
			}
			if base.GetExceptionHandler() != nil {
				t.Fatal("request registration changed base VM")
			}
		})
	}
}

func TestFPMExceptionHandlerRegistrationIsolated(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	first := data.NewStringValue("startup-1")
	second := data.NewStringValue("startup-2")
	base.SetExceptionHandler(first)
	base.SetExceptionHandler(second)
	var group sync.WaitGroup
	for i := 0; i < 16; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			request := fpm.New(base, nil)
			local := data.NewStringValue("request")
			if old := request.SetExceptionHandler(local); old != second {
				t.Error("startup handler not inherited")
			}
			request.RestoreExceptionHandler()
			if request.GetExceptionHandler() != second {
				t.Error("local stack not restored")
			}
			request.RestoreExceptionHandler()
			if request.GetExceptionHandler() != first {
				t.Error("startup stack not restored")
			}
		}()
	}
	group.Wait()
	if base.GetExceptionHandler() != second {
		t.Fatal("worker stack changed")
	}
}

func TestExceptionHandlerFailurePropagates(t *testing.T) {
	for _, test := range []struct{ name, callback, control string }{
		{"throws", "function (Throwable $e) { echo 'once'; throw new Exception('handler failed'); }", "Exception"},
		{"type_error", "function (int $e) { echo 'wrong'; }", "TypeError"},
		{"missing_parameter", "function (Throwable $e, $required) { echo 'wrong'; }", "ArgumentCountError"},
		{"exit", "function (Throwable $e) { exit(7); }", "exit"},
	} {
		t.Run(test.name, func(t *testing.T) {
			base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
			php.Load(base)
			file := filepath.Join(t.TempDir(), "handler-failure.php")
			if err := os.WriteFile(file, []byte("<?php set_exception_handler("+test.callback+"); throw new Exception('original');"), 0600); err != nil {
				t.Fatal(err)
			}
			var out strings.Builder
			request := fpm.New(base, func(s string) { out.WriteString(s) })
			_, original := request.LoadAndRun(file)
			if original == nil {
				t.Fatal("missing original exception")
			}
			handled, ctl := request.HandleUnhandledException(original)
			if !handled || ctl == nil || ctl == original {
				t.Fatal("handler failure not propagated")
			}
			if test.control == "exit" {
				exit, ok := ctl.(interface {
					IsExit() bool
					GetCode() int
				})
				if !ok || !exit.IsExit() || exit.GetCode() != 7 {
					t.Fatalf("wrong exit control: %v", ctl)
				}
			} else if thrown, ok := ctl.(*data.ThrowValue); !ok || thrown.GetName() != test.control {
				t.Fatalf("wrong exception: %v", ctl)
			}
			want := ""
			if test.name == "throws" {
				want = "once"
			}
			if out.String() != want {
				t.Fatalf("handler output %q", out.String())
			}
		})
	}
}
