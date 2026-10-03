package runtime

import (
	"fmt"
	"github.com/php-any/origami/parser"
	"os"
	"sync"
	"testing"
)

func TestRequestEnvironmentIsolation(t *testing.T) {
	const key = "ORIGAMI_REQUEST_ENV_TEST"
	t.Setenv(key, "process")
	base := NewVM(parser.NewParser()).(*VM)
	startup := "startup"
	if !base.SetPHPEnvironment(key, &startup) {
		t.Fatal("startup set")
	}
	frozen := NewRequestVM(base).(*RequestVM)
	later := "later"
	base.SetPHPEnvironment(key, &later)
	if got, _ := frozen.LookupPHPEnvironment(key); got != "startup" {
		t.Fatal("startup snapshot mutated", got)
	}
	var wg sync.WaitGroup
	errors := make(chan string, 32)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			vm := NewRequestVM(base).(*RequestVM)
			value := fmt.Sprint(i)
			vm.SetPHPEnvironment(key, &value)
			value = "caller mutation"
			if got, _ := vm.LookupPHPEnvironment(key); got != fmt.Sprint(i) {
				errors <- "request value escaped"
			}
			all := vm.PHPEnvironment()
			all[key] = "snapshot mutation"
			if got, _ := vm.LookupPHPEnvironment(key); got != fmt.Sprint(i) {
				errors <- "environment snapshot escaped"
			}
			vm.SetPHPEnvironment(key, nil)
			if _, ok := vm.LookupPHPEnvironment(key); ok {
				errors <- "removal ignored"
			}
			if vm.SetPHPEnvironment("bad=name", &value) || vm.SetPHPEnvironment("", nil) {
				errors <- "invalid name accepted"
			}
		}(i)
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		t.Error(err)
	}
	if got, _ := base.LookupPHPEnvironment(key); got != "later" {
		t.Fatal("request leaked to base", got)
	}
	if os.Getenv(key) != "process" {
		t.Fatal("VM mutated process environment")
	}
}
