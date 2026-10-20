package data

import (
	"strings"
	"sync/atomic"
	"unique"
)

type Symbol struct {
	display string
	key     unique.Handle[string]
}

// 已发布的一代只增不改。查找读这一代的 map；新名字用比较交换挂上新的一代。
type symbolGeneration struct {
	table map[string]SymbolID
	names []Symbol
}

// Symbols are append-only. A published prefix is never changed.
//
// 名字表用「不可变 map + 原子指针发布」而不是 sync.Map：查表是热路径
// （instanceof / GetFunc 每次请求数百次），而 sync.Map 会把 string 键装箱成 any，
// 再走 runtime.nilinterhash 计算接口哈希——实测占全部 CPU 的 4%。
// 普通 map 的 mapaccess2_faststr 直接按字符串哈希，且不产生接口值。
// 写入只在出现新符号时发生（启动期注册为主），代价是复制一张几千项的 map。
type SymbolTable struct {
	current atomic.Pointer[symbolGeneration]
}

var Symbols = NewSymbolTable()

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{}
	s.current.Store(&symbolGeneration{table: map[string]SymbolID{}, names: []Symbol{{}}})
	return s
}
func symbolKey(name string, sensitive bool) string {
	if sensitive {
		return "property:" + name
	}
	return nominalKey(strings.TrimPrefix(name, "\\"))
}
func (s *SymbolTable) Lookup(name string) (SymbolID, bool) {
	m := s.current.Load().table
	if id, ok := m[name]; ok {
		return id, true
	}
	// 去掉前导命名空间分隔符后的原名也在表里，这一步不需要规范化、不产生分配。
	trimmed := strings.TrimPrefix(name, "\\")
	if trimmed != name {
		if id, ok := m[trimmed]; ok {
			return id, true
		}
	}
	// 只剩大小写差异时才需要构造规范化键（会分配）。
	key := nominalKey(trimmed)
	if key == trimmed {
		return 0, false
	}
	if id, ok := m[key]; ok {
		return id, true
	}
	return 0, false
}
func (s *SymbolTable) Intern(name string) SymbolID         { return s.intern(name, false) }
func (s *SymbolTable) InternProperty(name string) SymbolID { return s.intern(name, true) }
func (s *SymbolTable) intern(name string, sensitive bool) SymbolID {
	key := symbolKey(name, sensitive)
	for {
		gen := s.current.Load()
		if id, ok := gen.table[key]; ok {
			return id
		}
		id := SymbolID(len(gen.names))
		names := make([]Symbol, len(gen.names)+1)
		copy(names, gen.names)
		names[id] = Symbol{display: strings.TrimPrefix(name, "\\"), key: unique.Make(key)}
		table := make(map[string]SymbolID, len(gen.table)+2)
		for existing, value := range gen.table {
			table[existing] = value
		}
		table[key] = id
		if !sensitive {
			table[strings.TrimPrefix(name, "\\")] = id
		}
		if s.current.CompareAndSwap(gen, &symbolGeneration{table: table, names: names}) {
			return id
		}
	}
}
func (s *SymbolTable) Name(id SymbolID) string {
	names := s.current.Load().names
	if int(id) >= len(names) {
		return ""
	}
	return names[id].display
}
