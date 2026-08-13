package origami_runtime

import (
	"os"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
)

// runPhp 用完整 std 环境运行一段 PHP 代码，返回 Control。
// 若 PHP 抛异常/错误，LoadAndRun 会返回非 nil 的 Control。
func runPhp(code string) data.Control {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	httplib.Load(vm)

	path := "/tmp/origami_runtime_test.php"
	_ = os.WriteFile(path, []byte("<?php\n"+code+"\n"), 0o644)
	_, ctl := vm.LoadAndRun(path)
	return ctl
}

// assertRunOK 断言 PHP 代码正常执行完成（无异常抛出）。
func assertRunOK(t *testing.T, code string) {
	t.Helper()
	if ctl := runPhp(code); ctl != nil {
		t.Fatalf("PHP 执行失败: %v", ctl)
	}
}

func TestDateIntervalISOParsing(t *testing.T) {
	code := `
$iv = new DateInterval('P1D');
if ($iv->d !== 1) { throw new Exception('P1D days != 1, got: '.$iv->d); }
$iv2 = new DateInterval('P1Y2M3DT4H5M6S');
if ($iv2->y !== 1 || $iv2->m !== 2 || $iv2->d !== 3 || $iv2->h !== 4 || $iv2->i !== 5 || $iv2->s !== 6) {
    throw new Exception('full spec parse failed');
}
$iv3 = new DateInterval('P2W');
if ($iv3->d !== 14) { throw new Exception('P2W days != 14, got: '.$iv3->d); }
echo "OK";
`
	assertRunOK(t, code)
}

func TestDateIntervalDaysDefaultFalse(t *testing.T) {
	code := `
$iv = new DateInterval('P1D');
if ($iv->days !== false) { throw new Exception('days should default to false, got: '.var_export($iv->days, true)); }
echo "OK";
`
	assertRunOK(t, code)
}

func TestDateTimeAddSub(t *testing.T) {
	code := `
$dt = new DateTime('2020-11-29');
$dt->add(new DateInterval('P1D'));
if ($dt->format('Y-m-d') !== '2020-11-30') { throw new Exception('add 1 day failed: '.$dt->format('Y-m-d')); }
$dt->sub(new DateInterval('P1D'));
if ($dt->format('Y-m-d') !== '2020-11-29') { throw new Exception('sub 1 day failed'); }
$dt2 = new DateTime('2020-01-31');
$dt2->add(new DateInterval('P1M'));
echo $dt2->format('Y-m-d');
echo "OK";
`
	assertRunOK(t, code)
}

func TestDateTimeComparison(t *testing.T) {
	code := `
$a = new DateTime('2020-11-29');
$b = new DateTime('2020-12-24');
if (!($a < $b)) { throw new Exception('a < b 应为 true'); }
if ($a == $b) { throw new Exception('a == b 应为 false'); }
if ($a > $b) { throw new Exception('a > b 应为 false'); }
if ($a <=> $b !== -1) { throw new Exception('spaceship 应为 -1'); }
echo "OK";
`
	assertRunOK(t, code)
}

func TestClosureFromCallable(t *testing.T) {
	code := `
$fn = Closure::fromCallable('strtoupper');
if ($fn('hi') !== 'HI') { throw new Exception('fromCallable string failed'); }
class Foo { public function greet($name) { return "hi ".$name; } }
$obj = new Foo();
$cfn = Closure::fromCallable([$obj, 'greet']);
if ($cfn('world') !== 'hi world') { throw new Exception('fromCallable array failed'); }
echo "OK";
`
	assertRunOK(t, code)
}
