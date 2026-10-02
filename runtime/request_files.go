package runtime

import "github.com/php-any/origami/data"

// A request inherits the completed bootstrap's include set. The snapshot is
// immutable; execution only adds to RequestVM's own maps. It is shared by all
// requests of the same bootstrap generation rather than copied per request.
type requestFileSnapshot struct {
	loaded  map[string]struct{}
	results map[string]data.GetValue
}

func (vm *VM) snapshotRequestFiles() *requestFileSnapshot {
	if snapshot := vm.requestFiles.Load(); snapshot != nil {
		return snapshot
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if snapshot := vm.requestFiles.Load(); snapshot != nil {
		return snapshot
	}
	snapshot := &requestFileSnapshot{loaded: make(map[string]struct{}), results: make(map[string]data.GetValue)}
	vm.phpFileCache.Range(func(key, _ any) bool { snapshot.loaded[key.(string)] = struct{}{}; return true })
	vm.includeOnceResults.Range(func(key, value any) bool { snapshot.results[key.(string)] = value.(data.GetValue); return true })
	vm.requestFiles.Store(snapshot)
	return snapshot
}
