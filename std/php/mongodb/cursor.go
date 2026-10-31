package mongodb

import (
	"context"
	"runtime"
	"time"

	"github.com/php-any/origami/data"
	oruntime "github.com/php-any/origami/runtime"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cursorState struct {
	cursor                  *mongo.Cursor
	owner                   *managerState // Retain the client until this cursor has been collected.
	mapping                 typeMap
	current                 data.Value
	started, done, iterator bool
	index                   int
}

func (s *cursorState) close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.cursor.Close(ctx)
	runtime.KeepAlive(s.owner)
}
func newCursor(ctx data.Context, c *mongo.Cursor, owner *managerState) *data.ClassValue {
	s := &cursorState{cursor: c, owner: owner, index: -1}
	runtime.SetFinalizer(s, (*cursorState).close)
	return object(ctx, driverNS+"Cursor", s)
}
func (s *cursorState) advance(ctx data.Context) (data.GetValue, data.Control) {
	s.started = true
	s.index++
	s.current = nil
	if s.done {
		return data.NewNullValue(), nil
	}
	if !s.cursor.Next(oruntime.RequestContext()) {
		s.done = true
		err := s.cursor.Err()
		s.close()
		if err != nil {
			return driverError(ctx, err)
		}
		return data.NewNullValue(), nil
	}
	var d bson.D
	if err := s.cursor.Decode(&d); err != nil {
		s.done = true
		s.close()
		return driverError(ctx, err)
	}
	v, err := decodeValue(ctx, d, s.mapping, true)
	if err != nil {
		s.done = true
		s.close()
		return invalid(err)
	}
	s.current = v
	return data.NewNullValue(), nil
}
func cursorClasses() []*nativeClass {
	c := &nativeClass{name: driverNS + "Cursor", interfaces: []string{"IteratorAggregate"}}
	c.methods = []*nativeMethod{
		method("__construct", nil, func(ctx data.Context) (data.GetValue, data.Control) {
			return logic("Cursor cannot be constructed directly")
		}),
		method("setTypeMap", []string{"typemap"}, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*cursorState](ctx)
			if !ok {
				return logic("uninitialized Cursor")
			}
			t, err := parseTypeMap(arg(ctx, 0))
			if err != nil {
				return invalid(err)
			}
			s.mapping = t
			if s.current != nil {
				var d bson.D
				if err = s.cursor.Decode(&d); err != nil {
					return driverError(ctx, err)
				}
				s.current, err = decodeValue(ctx, d, t, true)
				if err != nil {
					return invalid(err)
				}
			}
			return data.NewNullValue(), nil
		}),
		method("toArray", nil, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*cursorState](ctx)
			if !ok {
				return logic("uninitialized Cursor")
			}
			if s.index > 0 {
				return logic("Cursors cannot rewind after starting iteration")
			}
			if !s.started {
				if _, ctl := s.advance(ctx); ctl != nil {
					return nil, ctl
				}
			}
			a := data.NewArrayValue(nil).(*data.ArrayValue)
			for s.current != nil {
				a.AppendValue(s.current)
				if _, ctl := s.advance(ctx); ctl != nil {
					return nil, ctl
				}
			}
			return a, nil
		}),
		method("getIterator", nil, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*cursorState](ctx)
			if !ok {
				return logic("uninitialized Cursor")
			}
			if s.iterator {
				return logic("Cursors cannot yield multiple iterators")
			}
			s.iterator = true
			return object(ctx, "Origami\\MongoDB\\CursorIterator", s), nil
		}),
		method("isDead", nil, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*cursorState](ctx)
			if !ok {
				return logic("uninitialized Cursor")
			}
			return data.NewBoolValue(s.cursor.ID() == 0), nil
		}),
	}
	it := &nativeClass{name: "Origami\\MongoDB\\CursorIterator", interfaces: []string{"Iterator"}}
	for _, name := range []string{"rewind", "next", "current", "key", "valid"} {
		it.methods = append(it.methods, method(name, nil, func(ctx data.Context) (data.GetValue, data.Control) {
			s, ok := state[*cursorState](ctx)
			if !ok {
				return logic("uninitialized cursor iterator")
			}
			switch name {
			case "rewind":
				if s.index > 0 {
					return logic("Cursors cannot rewind after starting iteration")
				}
				if !s.started {
					return s.advance(ctx)
				}
			case "next":
				return s.advance(ctx)
			case "valid":
				return data.NewBoolValue(s.current != nil), nil
			case "current":
				if s.current != nil {
					return s.current, nil
				}
			case "key":
				if s.current != nil {
					return data.NewIntValue(s.index), nil
				}
			}
			return data.NewNullValue(), nil
		}))
	}
	return []*nativeClass{c, it}
}
