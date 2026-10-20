package runtime

import (
	"context"
	"os"
	"sync/atomic"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/perfmon"
	"github.com/php-any/origami/utils"
)

type parsedPHPFile struct {
	program        data.GetValue
	vars           []data.Variable
	classes        []data.ClassStmt
	interfaces     []data.InterfaceStmt
	sourceSize     int64
	sourceModified int64
}

type parsedFileFlight struct {
	done        chan struct{}
	owner       uint64
	entry       *parsedPHPFile
	control     data.Control
	completed   bool
	invalidated atomic.Bool
}

// 已发布的一代不被改写。命中只读 entries；第一次解析用比较交换登记 flight 和结果。
type parseGeneration struct {
	entries map[string]*parsedPHPFile
	flights map[string]*parsedFileFlight
}

// Swapping the root separates hot-reload generations. A parse already in
// progress may finish for its existing callers but cannot refill the new cache.
type parsedFileCache struct {
	current atomic.Pointer[parseGeneration]
}

func (cache *parsedFileCache) generation() *parseGeneration {
	if gen := cache.current.Load(); gen != nil {
		return gen
	}
	empty := &parseGeneration{entries: map[string]*parsedPHPFile{}, flights: map[string]*parsedFileFlight{}}
	if cache.current.CompareAndSwap(nil, empty) {
		return empty
	}
	return cache.current.Load()
}

func copyParsedEntries(source map[string]*parsedPHPFile) map[string]*parsedPHPFile {
	copied := make(map[string]*parsedPHPFile, len(source)+1)
	for file, entry := range source {
		copied[file] = entry
	}
	return copied
}

func copyParsedFlights(source map[string]*parsedFileFlight) map[string]*parsedFileFlight {
	copied := make(map[string]*parsedFileFlight, len(source)+1)
	for file, flight := range source {
		copied[file] = flight
	}
	return copied
}

func waitForRequestLoad(done <-chan struct{}) {
	request := RequestContext()
	data.CheckRequest(request)
	select {
	case <-done:
		data.CheckRequest(request)
	case <-request.Done():
		panic(data.ErrRequestCanceled)
	}
}

func (cache *parsedFileCache) load(request context.Context, file string, parse func() (*parsedPHPFile, data.Control)) (*parsedPHPFile, data.Control) {
	for {
		data.CheckRequest(request)
		gen := cache.generation()
		if cached, ok := gen.entries[file]; ok {
			return cached, nil
		}
		if flight, ok := gen.flights[file]; ok {
			if flight.owner == goid() {
				return nil, utils.NewThrowf("Recursive parse of %s", file)
			}
			select {
			case <-flight.done:
				data.CheckRequest(request)
			case <-request.Done():
				panic(data.ErrRequestCanceled)
			}
			if flight.completed && !flight.invalidated.Load() {
				return flight.entry, flight.control
			}
			continue
		}
		flight := &parsedFileFlight{done: make(chan struct{}), owner: goid()}
		flights := copyParsedFlights(gen.flights)
		flights[file] = flight
		if !cache.current.CompareAndSwap(gen, &parseGeneration{entries: gen.entries, flights: flights}) {
			continue
		}
		return cache.finishParse(request, file, flight, parse)
	}
}

func (cache *parsedFileCache) finishParse(request context.Context, file string, flight *parsedFileFlight, parse func() (*parsedPHPFile, data.Control)) (*parsedPHPFile, data.Control) {
	defer func() {
		cache.removeFlight(file, flight)
		close(flight.done)
	}()
	if cached, ok := cache.generation().entries[file]; ok {
		flight.entry, flight.completed = cached, true
		return cached, nil
	}
	entry, control := parse()
	data.CheckRequest(request)
	if control == nil && !flight.invalidated.Load() {
		cache.publish(file, entry)
		if flight.invalidated.Load() {
			cache.deleteEntry(file, entry)
		}
	}
	flight.entry, flight.control, flight.completed = entry, control, true
	return entry, control
}

func (cache *parsedFileCache) publish(file string, entry *parsedPHPFile) {
	for {
		gen := cache.generation()
		if current, ok := gen.entries[file]; ok && current == entry {
			return
		}
		entries := copyParsedEntries(gen.entries)
		entries[file] = entry
		if cache.current.CompareAndSwap(gen, &parseGeneration{entries: entries, flights: gen.flights}) {
			return
		}
	}
}

func (cache *parsedFileCache) deleteEntry(file string, entry *parsedPHPFile) {
	for {
		gen := cache.generation()
		current, ok := gen.entries[file]
		if !ok || (entry != nil && current != entry) {
			return
		}
		entries := copyParsedEntries(gen.entries)
		delete(entries, file)
		if cache.current.CompareAndSwap(gen, &parseGeneration{entries: entries, flights: gen.flights}) {
			return
		}
	}
}

func (cache *parsedFileCache) removeFlight(file string, flight *parsedFileFlight) {
	for {
		gen := cache.generation()
		current, ok := gen.flights[file]
		if !ok || current != flight {
			return
		}
		flights := copyParsedFlights(gen.flights)
		delete(flights, file)
		if cache.current.CompareAndSwap(gen, &parseGeneration{entries: gen.entries, flights: flights}) {
			return
		}
	}
}

func (cache *parsedFileCache) invalidate(file string) {
	if gen := cache.current.Load(); gen != nil {
		if flight := gen.flights[file]; flight != nil {
			flight.invalidated.Store(true)
		}
	}
	cache.deleteEntry(file, nil)
}

// InvalidateParsedFile runs on source writes, leaving cached reads unchanged.
func (vm *VM) InvalidateParsedFile(file string) {
	file = normalizePhpFilePath(file)
	vm.parsedFiles.Load().invalidate(file)
}
func (vm *RequestVM) InvalidateParsedFile(file string) { vm.Base.InvalidateParsedFile(file) }

// Includes already stat their source; reuse that result to detect writes made
// by external tools without adding a filesystem operation to cache hits.
func (vm *VM) ValidateParsedFileVersion(file string, info os.FileInfo) {
	file = normalizePhpFilePath(file)
	cache := vm.parsedFiles.Load()
	if entry, ok := cache.generation().entries[file]; ok &&
		(entry.sourceSize != info.Size() || entry.sourceModified != info.ModTime().UnixNano()) {
		vm.InvalidateParsedFile(file)
	}
}
func (vm *RequestVM) ValidateParsedFileVersion(file string, info os.FileInfo) {
	vm.Base.ValidateParsedFileVersion(file, info)
}

// ParseFileCached 解析 PHP 文件并缓存 AST（进程级，按规范化路径去重）。
// 多次执行同一入口/被 include 的文件时不应重复读盘与词法/语法分析。
func (vm *VM) ParseFileCached(file string) (data.GetValue, []data.Variable, data.Control) {
	return vm.parseFileCachedFor(file, nil)
}

func (vm *VM) parseFileCachedFor(file string, request *RequestVM) (data.GetValue, []data.Variable, data.Control) {
	file = normalizePhpFilePath(file)
	if file == "" {
		return nil, nil, nil
	}
	cache := vm.parsedFiles.Load()
	if cached, ok := cache.generation().entries[file]; ok {
		perfmon.NoteParse(file, true, 0)
		if request != nil {
			request.registerParsedDeclarations(cached)
		}
		return cached.program, cached.vars, nil
	}

	cached, acl := cache.load(RequestContext(), file, func() (*parsedPHPFile, data.Control) {
		t0 := perfmon.Now()
		p := vm.parser.Clone()
		// Parsing may register declarations for forward references. Keep those
		// registrations out of the VM catalog when compiling a request file.
		var declarations *RequestVM
		if request != nil {
			declarations = NewRequestVM(vm).(*RequestVM)
			p.SetVM(declarations)
			declarations.parser = p
		}
		info, _ := os.Stat(file)
		program, acl := p.ParseFile(file)
		perfmon.NoteParse(file, false, perfmon.Since(t0))
		if acl != nil {
			return nil, acl
		}
		entry := &parsedPHPFile{program: program, vars: p.GetVariables()}
		if info != nil {
			entry.sourceSize, entry.sourceModified = info.Size(), info.ModTime().UnixNano()
		}
		if declarations != nil {
			entry.classes = declarations.AddedClasses()
			for _, declaration := range declarations.addedInterfaces {
				entry.interfaces = append(entry.interfaces, declaration)
			}
		}
		p.ReleaseTokenBuffer()
		return entry, nil
	})
	if acl != nil {
		return nil, nil, acl
	}
	if request != nil {
		request.registerParsedDeclarations(cached)
	}
	return cached.program, cached.vars, nil
}

// ClearParsedFileCache 清空已缓存的 AST，供热重载使用。
func (vm *VM) ClearParsedFileCache() {
	vm.parsedFiles.Store(&parsedFileCache{})
}
