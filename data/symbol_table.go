package data

import (
	"strings"
	"sync"
	"sync/atomic"
	"unique"
)

type Symbol struct {
	display string
	key     unique.Handle[string]
}
type symbolSnapshot struct{ names []Symbol }

// Symbols are append-only. A published prefix is never changed.
//
// 名字表用「不可变 map + 原子指针发布」而不是 sync.Map：查表是热路径
// （instanceof / GetFunc 每次请求数百次），而 sync.Map 会把 string 键装箱成 any，
// 再走 runtime.nilinterhash 计算接口哈希——实测占全部 CPU 的 4%。
// 普通 map 的 mapaccess2_faststr 直接按字符串哈希，且不产生接口值。
// 写入只在出现新符号时发生（启动期注册为主），代价是复制一张几千项的 map。
type SymbolTable struct {
	mu       sync.Mutex
	table    atomic.Pointer[map[string]SymbolID]
	names    []Symbol
	snapshot atomic.Pointer[symbolSnapshot]
}

var Symbols = NewSymbolTable()

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{names: []Symbol{{}}}
	empty := map[string]SymbolID{}
	s.table.Store(&empty)
	s.snapshot.Store(&symbolSnapshot{names: s.names})
	return s
}
func symbolKey(name string, sensitive bool) string {
	if sensitive {
		return "property:" + name
	}
	return nominalKey(strings.TrimPrefix(name, "\\"))
}
func (s *SymbolTable) Lookup(name string) (SymbolID, bool) {
	m := s.table.Load()
	if id, ok := (*m)[name]; ok {
		return id, true
	}
	// 去掉前导命名空间分隔符后的原名也在表里，这一步不需要规范化、不产生分配。
	trimmed := strings.TrimPrefix(name, "\\")
	if trimmed != name {
		if id, ok := (*m)[trimmed]; ok {
			return id, true
		}
	}
	// 只剩大小写差异时才需要构造规范化键（会分配）。
	key := nominalKey(trimmed)
	if key == trimmed {
		return 0, false
	}
	if id, ok := (*m)[key]; ok {
		return id, true
	}
	return 0, false
}
func (s *SymbolTable) Intern(name string) SymbolID         { return s.intern(name, false) }
func (s *SymbolTable) InternProperty(name string) SymbolID { return s.intern(name, true) }
func (s *SymbolTable) intern(name string, sensitive bool) SymbolID {
	key := symbolKey(name, sensitive)
	if id, ok := (*s.table.Load())[key]; ok {
		return id
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current := *s.table.Load()
	if id, ok := current[key]; ok {
		return id
	}
	id := SymbolID(len(s.names))
	s.names = append(s.names, Symbol{display: strings.TrimPrefix(name, "\\"), key: unique.Make(key)})
	s.snapshot.Store(&symbolSnapshot{names: s.names})
	next := make(map[string]SymbolID, len(current)+2)
	for k, v := range current {
		next[k] = v
	}
	next[key] = id
	if !sensitive {
		next[strings.TrimPrefix(name, "\\")] = id
	}
	s.table.Store(&next)
	return id
}
func (s *SymbolTable) Name(id SymbolID) string {
	names := s.snapshot.Load().names
	if int(id) >= len(names) {
		return ""
	}
	return names[id].display
}
