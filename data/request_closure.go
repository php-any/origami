package data

// RequestClosureBinder is a setup-time boundary. Executing a closure has no
// scope lookup: worker policies rebind its captured objects before handling.
type RequestClosureBinder interface {
	RequestScopeObjects() []*ClassValue
	BindRequestScope(Context, map[*ObjectValue]*ClassValue) FuncStmt
}

// RequestCaptureScope is shared by all callbacks copied into one request.
// It preserves aliases and cycles, including references shared by callbacks.
// It runs only during setup; normal closure calls perform no scope lookup.
type RequestCaptureScope struct {
	Context     Context
	Objects     map[*ObjectValue]*ClassValue
	ScopeObject func(*ClassValue) *ClassValue
	ScopeArray  func(*ArrayValue) *ArrayValue
	ScopeSlot   func(*ZVal) *ZVal
	ScopeValue  func(Value) Value
	values      map[Value]Value
	closures    map[FuncStmt]FuncStmt
	slots       map[*ZVal]*ZVal
	arrayScope  *ArrayOverlayScope
}

type requestCaptureBinder interface {
	BindRequestCapture(*RequestCaptureScope) FuncStmt
}

func NewRequestCaptureScope(ctx Context, objects map[*ObjectValue]*ClassValue) *RequestCaptureScope {
	s := &RequestCaptureScope{Context: ctx, Objects: objects, values: map[Value]Value{}, closures: map[FuncStmt]FuncStmt{}, slots: map[*ZVal]*ZVal{}, arrayScope: NewArrayOverlayScope()}
	s.arrayScope.refs, s.arrayScope.bindValue = s.slots, s.Bind
	return s
}
func (s *RequestCaptureScope) Closure(source FuncStmt) FuncStmt        { return s.closures[source] }
func (s *RequestCaptureScope) RememberClosure(source, target FuncStmt) { s.closures[source] = target }
func (s *RequestCaptureScope) BindSlot(slot *ZVal) *ZVal {
	if s.ScopeSlot != nil {
		return s.ScopeSlot(slot)
	}
	return s.arrayScope.copySlot(slot, nil)
}
func (s *RequestCaptureScope) bindClosure(source FuncStmt) FuncStmt {
	if clone := s.closures[source]; clone != nil {
		return clone
	}
	if binder, ok := source.(requestCaptureBinder); ok {
		return binder.BindRequestCapture(s)
	}
	if binder, ok := source.(RequestClosureBinder); ok {
		if s.ScopeObject != nil {
			for _, object := range binder.RequestScopeObjects() {
				s.ScopeObject(object)
			}
		}
		target := binder.BindRequestScope(s.Context, s.Objects)
		s.RememberClosure(source, target)
		return target
	}
	return source
}
func (s *RequestCaptureScope) Bind(value Value) Value {
	if value == nil {
		return nil
	}
	if bound, ok := s.values[value]; ok {
		return bound
	}
	s.values[value] = value
	var bound Value = value
	switch v := value.(type) {
	case *ClassValue:
		if v.ObjectValue != nil {
			if target := s.Objects[v.ObjectValue]; target != nil {
				bound = target
			} else if s.ScopeObject != nil {
				bound = s.ScopeObject(v)
			}
		}
	case *ThisValue:
		if v.ObjectValue != nil {
			if target := s.Objects[v.ObjectValue]; target != nil {
				bound = NewThisValue(target)
			} else if s.ScopeObject != nil {
				bound = NewThisValue(s.ScopeObject(v.ClassValue))
			}
		}
	case *ArrayValue:
		if s.ScopeArray != nil {
			bound = s.ScopeArray(v)
			break
		}
		clone := &ArrayValue{
			flatArrayStore:        FlatArrayStore{entries: make([]*ZVal, v.Len()), appendKeyKnown: v.appendKeyKnown, intKeySeen: v.intKeySeen, nextIntKey: v.nextIntKey, iterator: v.iterator},
			IndirectOverloadClass: v.IndirectOverloadClass,
		}
		s.values[value] = clone
		for i, slot := range v.Range() {
			clone.ReplaceSlot(i, s.BindSlot(slot))
		}
		bound = clone
	case *FuncValue:
		clone := NewFuncValue(v.Value)
		s.values[value] = clone
		clone.Value = s.bindClosure(v.Value)
		bound = clone
	case *BoundFuncValue:
		owner := v.BoundObject
		if owner != nil {
			if target := s.Objects[owner.ObjectValue]; target != nil {
				owner = target
			} else if s.ScopeObject != nil {
				owner = s.ScopeObject(owner)
			}
		}
		clone := NewBoundFuncValue(v.Value, v.ScopeClass, owner)
		s.values[value] = clone
		clone.Value = s.bindClosure(v.Value)
		bound = clone
	case *ArraySlotRef:
		clone := &ArraySlotRef{}
		s.values[value] = clone
		clone.Slot = s.BindSlot(v.Slot)
		bound = clone
	default:
		if s.ScopeValue != nil {
			bound = s.ScopeValue(value)
		}
	}
	s.values[value] = bound
	return bound
}
