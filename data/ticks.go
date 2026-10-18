package data

// TickState belongs to one execution/request, never to the compiled AST.
type TickState struct {
	Callbacks []TickCallback
	Count     int
	Running   bool
}

type TickCallback struct {
	Callback   Value
	Invocation FuncStmt
}

type TickProvider interface{ PHPTickState() *TickState }

func (s *TickState) Advance(ctx Context, interval int) Control {
	if s.Running {
		return nil
	}
	s.Count++
	if s.Count < interval {
		return nil
	}
	s.Count = 0
	s.Running = true
	defer func() { s.Running = false }()
	// Callbacks may change the registration list while it is being dispatched.
	callbacks := append([]TickCallback(nil), s.Callbacks...)
	for _, callback := range callbacks {
		if _, ctl := callback.Invocation.Call(ctx); ctl != nil {
			return ctl
		}
	}
	return nil
}
