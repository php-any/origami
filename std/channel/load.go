package channel

import (
	"github.com/php-any/origami/data"
)

// Load 加载 channel 模块
func Load(vm data.VM) data.Control {
	// 注册 Channel 类
	channelClass := NewChannelClass()
	return vm.AddClass(channelClass)
}
