package data

import "sync"

// StaticLocals 函数/方法内 static 局部变量存储（跨调用、递归共享同一 ZVal）。
type StaticLocals struct {
	mu    sync.Mutex
	Slots map[int]*ZVal
}

func NewStaticLocals() *StaticLocals {
	return &StaticLocals{Slots: make(map[int]*ZVal)}
}

func (s *StaticLocals) Get(index int) (Value, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	z, ok := s.Slots[index]
	if !ok || z == nil {
		return nil, false
	}
	return z.Value, true
}

func (s *StaticLocals) Slot(index int) *ZVal {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Slots[index]
}

func (s *StaticLocals) Init(index int, val Value) Value {
	return s.EnsureSlot(index, val).Value
}

// EnsureSlot 返回（必要时创建）指定下标的共享 ZVal。递归调用必须共用此指针。
func (s *StaticLocals) EnsureSlot(index int, val Value) *ZVal {
	if s == nil {
		return NewZVal(val)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if z, ok := s.Slots[index]; ok && z != nil {
		return z
	}
	z := &ZVal{Value: val, Defined: true}
	s.Slots[index] = z
	return z
}

// Update 在 static 变量已注册后同步最新值（用于 ++/-- 等修改）。
func (s *StaticLocals) Update(index int, val Value) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if z, ok := s.Slots[index]; ok && z != nil {
		z.Value = val
		z.Defined = true
	}
}
