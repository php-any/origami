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
type SymbolTable struct {
	mu       sync.Mutex
	byName   sync.Map
	names    []Symbol
	snapshot atomic.Pointer[symbolSnapshot]
}

var Symbols = NewSymbolTable()

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{names: []Symbol{{}}}
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
	if id, ok := s.byName.Load(name); ok {
		return id.(SymbolID), true
	}
	id, ok := s.byName.Load(symbolKey(name, false))
	if !ok {
		return 0, false
	}
	return id.(SymbolID), true
}
func (s *SymbolTable) Intern(name string) SymbolID         { return s.intern(name, false) }
func (s *SymbolTable) InternProperty(name string) SymbolID { return s.intern(name, true) }
func (s *SymbolTable) intern(name string, sensitive bool) SymbolID {
	key := symbolKey(name, sensitive)
	if id, ok := s.byName.Load(key); ok {
		return id.(SymbolID)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.byName.Load(key); ok {
		return id.(SymbolID)
	}
	id := SymbolID(len(s.names))
	s.names = append(s.names, Symbol{display: strings.TrimPrefix(name, "\\"), key: unique.Make(key)})
	s.snapshot.Store(&symbolSnapshot{names: s.names})
	s.byName.Store(key, id)
	if !sensitive {
		s.byName.Store(strings.TrimPrefix(name, "\\"), id)
	}
	return id
}
func (s *SymbolTable) Name(id SymbolID) string {
	names := s.snapshot.Load().names
	if int(id) >= len(names) {
		return ""
	}
	return names[id].display
}
