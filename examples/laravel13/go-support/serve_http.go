package gosupport

import (
	"fmt"
	"os"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

// RunHTTPServe 跳过 artisan / Console Kernel（不扫描全部 Artisan 命令），
// 只走 Composer autoload + bootstrap/app.php + HTTP Kernel。
// 对应 `go run . serve`，list/migrate 等仍走官方 artisan。
func RunHTTPServe(vm *runtime.VM) error {
	BootLog("HTTP serve path (skip artisan command discovery)")
	t := time.Now()
	if _, ctl := vm.LoadAndRun("vendor/autoload.php"); ctl != nil {
		return fmt.Errorf("autoload.php: %s", ctl.AsString())
	}
	BootLog("vendor/autoload.php " + time.Since(t).Truncate(time.Millisecond).String())

	t = time.Now()
	appVal, ctl := vm.LoadAndRun("bootstrap/app.php")
	if ctl != nil {
		return fmt.Errorf("bootstrap/app.php: %s", ctl.AsString())
	}
	app, ok := appVal.(*data.ClassValue)
	if !ok {
		return fmt.Errorf("bootstrap/app.php 未返回 Application（%T）", appVal)
	}
	BootLog("bootstrap/app.php " + time.Since(t).Truncate(time.Millisecond).String())

	host, port, err := serveAddress(os.Args)
	if err != nil {
		return err
	}
	return runLaravelHTTPServer(host, port, vm, app)
}
