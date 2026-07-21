package container

import (
	"sync"

	"github.com/php-any/origami/data"
)

var (
	applicationContainerMu sync.Mutex
	applicationContainer   data.GetValue
)

func setApplicationContainer(c data.GetValue) func() {
	applicationContainerMu.Lock()
	prev := applicationContainer
	applicationContainer = c
	applicationContainerMu.Unlock()
	return func() {
		applicationContainerMu.Lock()
		applicationContainer = prev
		applicationContainerMu.Unlock()
	}
}

// ApplicationContainer 返回 #[Application] 扫描期绑定的 Container 实例；未在扫描中时为 nil。
func ApplicationContainer() data.GetValue {
	applicationContainerMu.Lock()
	c := applicationContainer
	applicationContainerMu.Unlock()
	return c
}
