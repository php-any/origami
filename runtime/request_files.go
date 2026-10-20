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
	snapshot := &requestFileSnapshot{loaded: make(map[string]struct{}, len(vm.phpFileCache)), results: make(map[string]data.GetValue, len(vm.includeOnceResults))}
	for file := range vm.phpFileCache {
		snapshot.loaded[file] = struct{}{}
	}
	for file, result := range vm.includeOnceResults {
		snapshot.results[file] = result
	}
	vm.requestFiles.Store(snapshot)
	return snapshot
}
