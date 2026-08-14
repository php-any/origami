package runtime

import (
	"strings"

	"github.com/php-any/origami/data"
)

func (vm *VM) WriteOutput(s string) {
	// 无缓冲时（最常见的 echo 场景）走零锁快速路径。
	if !vm.hasOutputBuffer.Load() {
		data.WriteOutput(s)
		return
	}
	vm.mu.Lock()
	if n := len(vm.outputBuffers); n > 0 {
		vm.outputBuffers[n-1].WriteString(s)
		vm.mu.Unlock()
		return
	}
	vm.mu.Unlock()
	data.WriteOutput(s)
}

func (vm *VM) StartOutputBuffer() {
	vm.mu.Lock()
	vm.outputBuffers = append(vm.outputBuffers, &strings.Builder{})
	vm.hasOutputBuffer.Store(true)
	vm.mu.Unlock()
}

// syncBufferFlag 在修改缓冲栈后同步 hasOutputBuffer 原子标记。
func (vm *VM) syncBufferFlag() {
	vm.hasOutputBuffer.Store(len(vm.outputBuffers) > 0)
}

func (vm *VM) CleanOutputBuffer() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	content := vm.outputBuffers[n-1].String()
	vm.outputBuffers = vm.outputBuffers[:n-1]
	vm.syncBufferFlag()
	return content, true
}

// FlushOutputBuffer 弹出并返回栈顶缓冲内容，并把内容写出到上一层（或最终输出）。
func (vm *VM) FlushOutputBuffer() (string, bool) {
	vm.mu.Lock()
	n := len(vm.outputBuffers)
	if n == 0 {
		vm.mu.Unlock()
		return "", false
	}
	content := vm.outputBuffers[n-1].String()
	vm.outputBuffers = vm.outputBuffers[:n-1]
	vm.syncBufferFlag()
	vm.mu.Unlock()
	// 写出到上一层缓冲或最终输出。
	vm.WriteOutput(content)
	return content, true
}

// CleanCurrentBuffer 清空栈顶缓冲内容但不结束缓冲。
func (vm *VM) CleanCurrentBuffer() bool {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return false
	}
	vm.outputBuffers[n-1] = &strings.Builder{}
	return true
}

// OutputBufferLength 返回栈顶缓冲的字节长度（无缓冲返回 false）。
func (vm *VM) OutputBufferLength() (int, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return 0, false
	}
	return vm.outputBuffers[n-1].Len(), true
}

// OutputBufferStatus 返回缓冲层状态列表；full=true 返回全部层，否则仅最顶层。
func (vm *VM) OutputBufferStatus(full bool) []data.OutputBufferStatusInfo {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return nil
	}
	start := 0
	if !full {
		start = n - 1
	}
	result := make([]data.OutputBufferStatusInfo, 0, n-start)
	for i := start; i < n; i++ {
		result = append(result, data.OutputBufferStatusInfo{
			Level:      i + 1,
			Type:       1,
			Flags:      0,
			ChunkSize:  0,
			BufferSize: vm.outputBuffers[i].Len(),
			Name:       "default output handler",
		})
	}
	return result
}

// ListOutputHandlers 返回所有激活缓冲的处理器名。
func (vm *VM) ListOutputHandlers() []string {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return nil
	}
	result := make([]string, 0, n)
	for range vm.outputBuffers {
		result = append(result, "default output handler")
	}
	return result
}

// SetImplicitFlush 设置/清除隐式刷新标志。
func (vm *VM) SetImplicitFlush(on bool) {
	vm.mu.Lock()
	vm.implicitFlush = on
	vm.mu.Unlock()
}

// IsImplicitFlush 返回当前隐式刷新标志。
func (vm *VM) IsImplicitFlush() bool {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.implicitFlush
}

func (vm *VM) OutputBufferContents() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	return vm.outputBuffers[n-1].String(), true
}

func (vm *VM) OutputBufferLevel() int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return len(vm.outputBuffers)
}

func (vm *TempVM) WriteOutput(s string) {
	// 无缓冲时（最常见的 echo 场景）走零锁快速路径。
	if !vm.hasOutputBuffer.Load() {
		vm.Base.WriteOutput(s)
		return
	}
	vm.mu.Lock()
	if n := len(vm.outputBuffers); n > 0 {
		vm.outputBuffers[n-1].WriteString(s)
		vm.mu.Unlock()
		return
	}
	vm.mu.Unlock()
	vm.Base.WriteOutput(s)
}

func (vm *TempVM) StartOutputBuffer() {
	vm.mu.Lock()
	vm.outputBuffers = append(vm.outputBuffers, &strings.Builder{})
	vm.hasOutputBuffer.Store(true)
	vm.mu.Unlock()
}

// syncBufferFlag 在修改缓冲栈后同步 hasOutputBuffer 原子标记。
func (vm *TempVM) syncBufferFlag() {
	vm.hasOutputBuffer.Store(len(vm.outputBuffers) > 0)
}

func (vm *TempVM) CleanOutputBuffer() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	content := vm.outputBuffers[n-1].String()
	vm.outputBuffers = vm.outputBuffers[:n-1]
	vm.syncBufferFlag()
	return content, true
}

// FlushOutputBuffer 弹出并返回栈顶缓冲内容，并把内容写出到上一层（或最终输出）。
func (vm *TempVM) FlushOutputBuffer() (string, bool) {
	vm.mu.Lock()
	n := len(vm.outputBuffers)
	if n == 0 {
		vm.mu.Unlock()
		return "", false
	}
	content := vm.outputBuffers[n-1].String()
	vm.outputBuffers = vm.outputBuffers[:n-1]
	vm.syncBufferFlag()
	vm.mu.Unlock()
	// 写出到上一层缓冲或最终输出。
	vm.WriteOutput(content)
	return content, true
}

// CleanCurrentBuffer 清空栈顶缓冲内容但不结束缓冲。
func (vm *TempVM) CleanCurrentBuffer() bool {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return false
	}
	vm.outputBuffers[n-1] = &strings.Builder{}
	return true
}

// OutputBufferLength 返回栈顶缓冲的字节长度（无缓冲返回 false）。
func (vm *TempVM) OutputBufferLength() (int, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return 0, false
	}
	return vm.outputBuffers[n-1].Len(), true
}

// OutputBufferStatus 返回缓冲层状态列表；full=true 返回全部层，否则仅最顶层。
func (vm *TempVM) OutputBufferStatus(full bool) []data.OutputBufferStatusInfo {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return nil
	}
	start := 0
	if !full {
		start = n - 1
	}
	result := make([]data.OutputBufferStatusInfo, 0, n-start)
	for i := start; i < n; i++ {
		result = append(result, data.OutputBufferStatusInfo{
			Level:      i + 1,
			Type:       1,
			Flags:      0,
			ChunkSize:  0,
			BufferSize: vm.outputBuffers[i].Len(),
			Name:       "default output handler",
		})
	}
	return result
}

// ListOutputHandlers 返回所有激活缓冲的处理器名。
func (vm *TempVM) ListOutputHandlers() []string {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return nil
	}
	result := make([]string, 0, n)
	for range vm.outputBuffers {
		result = append(result, "default output handler")
	}
	return result
}

// SetImplicitFlush 设置/清除隐式刷新标志。
func (vm *TempVM) SetImplicitFlush(on bool) {
	vm.mu.Lock()
	vm.implicitFlush = on
	vm.mu.Unlock()
}

// IsImplicitFlush 返回当前隐式刷新标志。
func (vm *TempVM) IsImplicitFlush() bool {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.implicitFlush
}

func (vm *TempVM) OutputBufferContents() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	return vm.outputBuffers[n-1].String(), true
}

func (vm *TempVM) OutputBufferLevel() int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return len(vm.outputBuffers)
}

var (
	_ data.OutputSink       = (*VM)(nil)
	_ data.OutputBufferHost = (*VM)(nil)
	_ data.OutputSink       = (*TempVM)(nil)
	_ data.OutputBufferHost = (*TempVM)(nil)
)
