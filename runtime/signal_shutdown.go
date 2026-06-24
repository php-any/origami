package runtime

import (
	"os"
	"sync"
	"syscall"

	"github.com/php-any/origami/data"
	signalsrc "os/signal"
)

var shutdownSignalOnce sync.Once

// InstallShutdownSignalHandler 监听 SIGINT/SIGTERM，触发时执行 vm.RunShutdownCallbacks()。
// 使用独立 channel，与用户侧 Signal\notify 互不干扰。
func InstallShutdownSignalHandler(vm data.VM) {
	shutdownSignalOnce.Do(func() {
		go func() {
			ch := make(chan os.Signal, 1)
			signalsrc.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			vm.RunShutdownCallbacks()
		}()
	})
}
