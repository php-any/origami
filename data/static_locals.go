package data

import "sync"

// StaticLocals 函数/方法内 static 局部变量存储（跨调用持久）
type StaticLocals struct {
	mu   sync.Mutex
	Vals map[int]Value
}

func NewStaticLocals() *StaticLocals {
	return &StaticLocals{Vals: make(map[int]Value)}
}

func (s *StaticLocals) Get(index int) (Value, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.Vals[index]
	return v, ok
}

func (s *StaticLocals) Init(index int, val Value) Value {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.Vals[index]; ok {
		return v
	}
	s.Vals[index] = val
	return val
}
