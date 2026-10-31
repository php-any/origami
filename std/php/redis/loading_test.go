package redis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/redis"
)

func hostVM() *runtime.VM {
	vm := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	return vm
}
func TestRedisExplicitRegistration(t *testing.T) {
	vm := hostVM()
	check := func(vm *runtime.VM, want bool) {
		t.Helper()
		for _, name := range []string{"Redis", "RedisException"} {
			if _, ok := vm.GetClass(name); ok != want {
				t.Fatalf("%s registered=%v, want %v", name, ok, want)
			}
		}
		if _, ok := vm.GetConstant("REDIS_VERSION"); ok != want {
			t.Fatal("version registration")
		}
		fn, _ := vm.GetFunc("extension_loaded")
		ctx := vm.CreateContext(fn.GetVariables())
		ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue("redis")))
		v, ctl := fn.Call(ctx)
		if ctl != nil || v.(*data.BoolValue).Value != want {
			t.Fatalf("extension_loaded: %v, %v", v, ctl)
		}
	}
	check(vm, false)
	redis.Load(vm)
	check(vm, true)
	check(hostVM(), false)
}

type userRedisClass struct{ data.StdClass }

func (*userRedisClass) GetName() string { return "Redis" }
func TestUserRedisClassDoesNotAdvertiseExtension(t *testing.T) {
	vm := hostVM()
	vm.AddClass(&userRedisClass{})
	fn, _ := vm.GetFunc("extension_loaded")
	ctx := vm.CreateContext(fn.GetVariables())
	ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue("redis")))
	v, ctl := fn.Call(ctx)
	if ctl != nil || v.(*data.BoolValue).Value {
		t.Fatal("a user Redis class must not advertise ext-redis")
	}
}
func runPHP(t *testing.T, name string) {
	t.Helper()
	vm := hostVM()
	redis.Load(vm)
	path, e := filepath.Abs("../../../tests/php/" + name)
	if e != nil {
		t.Fatal(e)
	}
	if _, ctl := vm.LoadAndRun(path); ctl != nil {
		t.Fatalf("%s: %v", name, ctl)
	}
}
func TestRedisPHPContract(t *testing.T) { runPHP(t, "redis_driver_contract_test.php") }
func TestRedisIntegration(t *testing.T) {
	if os.Getenv("ORIGAMI_REDIS_TEST_PORT") == "" {
		t.Skip("set ORIGAMI_REDIS_TEST_PORT for real Redis integration")
	}
	runPHP(t, "redis_driver_integration_test.php")
	runPHP(t, "redis_phpredis_features_test.php")
}
