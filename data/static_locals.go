package data

import "sync"

// StaticLocals 函数/方法内 static 局部变量存储（跨调用持久）。
// 内部使用 *ZVal 而非 Value，使上下文与 store 共享同一 ZVal 指针，
// 递归调用时 inner 修改静态变量的值能立即被 outer 看到。
type StaticLocals struct {
	mu   sync.Mutex
	Vals map[int]*ZVal
}

func NewStaticLocals() *StaticLocals {
	return &StaticLocals{Vals: make(map[int]*ZVal)}
}

// GetZVal 返回指定索引的 ZVal 指针（nil 表示不存在）。
func (s *StaticLocals) GetZVal(index int) *ZVal {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Vals[index]
}

func (s *StaticLocals) Get(index int) (Value, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if zv := s.Vals[index]; zv != nil {
		return zv.Value, true
	}
	return nil, false
}

func (s *StaticLocals) Init(index int, val Value) Value {
	s.mu.Lock()
	defer s.mu.Unlock()
	if zv := s.Vals[index]; zv != nil {
		return zv.Value
	}
	zv := &ZVal{Value: val}
	s.Vals[index] = zv
	return val
}

// InitZVal 初始化 ZVal 并返回其指针，供上下文直接共享。
func (s *StaticLocals) InitZVal(index int, val Value) *ZVal {
	s.mu.Lock()
	defer s.mu.Unlock()
	if zv := s.Vals[index]; zv != nil {
		return zv
	}
	zv := &ZVal{Value: val}
	s.Vals[index] = zv
	return zv
}

func (s *StaticLocals) Set(index int, val Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if zv := s.Vals[index]; zv != nil {
		zv.Value = val
		return
	}
	s.Vals[index] = &ZVal{Value: val}
}

func (s *StaticLocals) SetNoLock(index int, val Value) {
	if zv := s.Vals[index]; zv != nil {
		zv.Value = val
		return
	}
	s.Vals[index] = &ZVal{Value: val}
}

func (s *StaticLocals) Range(fn func(index int, val Value) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, zv := range s.Vals {
		if zv != nil {
			if !fn(k, zv.Value) {
				break
			}
		}
	}
}

func (s *StaticLocals) CollectIndices() []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	indices := make([]int, 0, len(s.Vals))
	for k := range s.Vals {
		indices = append(indices, k)
	}
	return indices
}
