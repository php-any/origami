package runtime

import (
	"context"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/perfmon"
	"github.com/php-any/origami/utils"
)

type parsedPHPFile struct {
	program    data.GetValue
	vars       []data.Variable
	classes    []data.ClassStmt
	interfaces []data.InterfaceStmt
}

type parsedFileFlight struct {
	done      chan struct{}
	owner     uint64
	entry     *parsedPHPFile
	control   data.Control
	completed bool
}

// Swapping the root separates hot-reload generations. A parse already in
// progress may finish for its existing callers but cannot refill the new cache.
type parsedFileCache struct {
	entries sync.Map
	flights sync.Map
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
		if cached, ok := syncMapLoad[*parsedPHPFile](&cache.entries, file); ok {
			return cached, nil
		}
		candidate := &parsedFileFlight{done: make(chan struct{}), owner: goid()}
		actual, loaded := cache.flights.LoadOrStore(file, candidate)
		flight := actual.(*parsedFileFlight)
		if loaded {
			if flight.owner == candidate.owner {
				return nil, utils.NewThrowf("Recursive parse of %s", file)
			}
			select {
			case <-flight.done:
				data.CheckRequest(request)
			case <-request.Done():
				panic(data.ErrRequestCanceled)
			}
			if flight.completed {
				return flight.entry, flight.control
			}
			// A canceled/panicking leader owns no result; another request retries.
			continue
		}
		return func() (*parsedPHPFile, data.Control) {
			defer func() { cache.flights.Delete(file); close(flight.done) }()
			// Another loader can publish between the first lookup and election.
			if cached, ok := syncMapLoad[*parsedPHPFile](&cache.entries, file); ok {
				flight.entry, flight.completed = cached, true
				return cached, nil
			}
			entry, control := parse()
			data.CheckRequest(request)
			if control == nil {
				cache.entries.Store(file, entry)
			}
			flight.entry, flight.control, flight.completed = entry, control, true
			return entry, control
		}()
	}
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
	if cached, ok := syncMapLoad[*parsedPHPFile](&cache.entries, file); ok {
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
		program, acl := p.ParseFile(file)
		perfmon.NoteParse(file, false, perfmon.Since(t0))
		if acl != nil {
			return nil, acl
		}
		entry := &parsedPHPFile{program: program, vars: p.GetVariables()}
		if declarations != nil {
			entry.classes = declarations.AddedClasses()
			for _, declaration := range declarations.addedInterfaces {
				entry.interfaces = append(entry.interfaces, declaration)
			}
		}
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
