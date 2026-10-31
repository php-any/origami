package websocket

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/net/internal/script"
	"github.com/php-any/origami/utils"
)

func NewConnClass() data.ClassStmt {
	return &ConnClass{source: nil}
}

func NewConnClassFrom(source *websocket.Conn) data.ClassStmt {
	return &ConnClass{source: source}
}

type ConnClass struct {
	stop func() bool
	node.Node
	source *websocket.Conn
	limits map[string]int
}

func (s *ConnClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewConnClassFrom(nil), ctx.CreateBaseContext()), nil
}

func (s *ConnClass) GetName() string         { return "Net\\Websocket\\Conn" }
func (s *ConnClass) GetExtend() *string      { return nil }
func (s *ConnClass) GetImplements() []string { return nil }
func (s *ConnClass) AsString() string        { return "Conn{}" }
func (s *ConnClass) GetSource() any          { return s.source }

func (s *ConnClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "configure":
		return s.configureMethod(), true
	case "readText":
		return &ConnReadTextMethod{source: s.source, limits: s.limits}, true
	case "writeText":
		return &ConnWriteTextMethod{source: s.source, limits: s.limits}, true
	case "close":
		return &ConnCloseMethod{source: s.source, stop: s.stop}, true
	}
	return nil, false
}

func (s *ConnClass) GetMethods() []data.Method {
	return []data.Method{
		s.configureMethod(),
		&ConnReadTextMethod{source: s.source, limits: s.limits},
		&ConnWriteTextMethod{source: s.source, limits: s.limits},
		&ConnCloseMethod{source: s.source, stop: s.stop},
	}
}

func (s *ConnClass) GetConstruct() data.Method { return nil }

func (s *ConnClass) GetProperty(name string) (data.Property, bool) { return nil, false }

func (s *ConnClass) GetPropertyList() []data.Property { return []data.Property{} }

type ConnReadTextMethod struct {
	limits map[string]int
	source *websocket.Conn
}

func (h *ConnReadTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return nil, utils.NewThrow(errors.New("websocket 连接未初始化"))
	}
	if h.limits != nil {
		if e := h.source.SetReadDeadline(time.Now().Add(time.Duration(h.limits["readTimeoutMs"]) * time.Millisecond)); e != nil {
			return nil, utils.NewThrow(e)
		}
	}
	_, payload, err := h.source.ReadMessage()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewStringValue(string(payload)), nil
}

func (h *ConnReadTextMethod) GetName() string            { return "readText" }
func (h *ConnReadTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnReadTextMethod) GetIsStatic() bool          { return false }

var connReadTextMethodGetParams = []data.GetValue{}

func (h *ConnReadTextMethod) GetParams() []data.GetValue { return connReadTextMethodGetParams }

var connReadTextMethodGetVariables = []data.Variable{}

func (h *ConnReadTextMethod) GetVariables() []data.Variable {
	return connReadTextMethodGetVariables
}
func (h *ConnReadTextMethod) GetReturnType() data.Types { return data.NewBaseType("string") }

type ConnWriteTextMethod struct {
	limits map[string]int
	source *websocket.Conn
}

func (h *ConnWriteTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return nil, utils.NewThrow(errors.New("websocket 连接未初始化"))
	}
	msg, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if h.limits != nil {
		if len(msg) > h.limits["maxMessageBytes"] {
			return nil, utils.NewThrowf("message exceeds maxMessageBytes")
		}
		if e := h.source.SetWriteDeadline(time.Now().Add(time.Duration(h.limits["writeTimeoutMs"]) * time.Millisecond)); e != nil {
			return nil, utils.NewThrow(e)
		}
	}
	if err = h.source.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ConnWriteTextMethod) GetName() string            { return "writeText" }
func (h *ConnWriteTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnWriteTextMethod) GetIsStatic() bool          { return false }

var connWriteTextMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "message", 0, nil, data.String{}),
}

func (h *ConnWriteTextMethod) GetParams() []data.GetValue {
	return connWriteTextMethodGetParams
}

var connWriteTextMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "message", 0, nil),
}

func (h *ConnWriteTextMethod) GetVariables() []data.Variable {
	return connWriteTextMethodGetVariables
}
func (h *ConnWriteTextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

type ConnCloseMethod struct {
	stop   func() bool
	source *websocket.Conn
}

func (h *ConnCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return nil, utils.NewThrow(errors.New("websocket 连接未初始化"))
	}
	if h.stop != nil {
		h.stop()
	}
	if err := h.source.Close(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ConnCloseMethod) GetName() string            { return "close" }
func (h *ConnCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnCloseMethod) GetIsStatic() bool          { return false }

var connCloseMethodGetParams = []data.GetValue{}

func (h *ConnCloseMethod) GetParams() []data.GetValue { return connCloseMethodGetParams }

var connCloseMethodGetVariables = []data.Variable{}

func (h *ConnCloseMethod) GetVariables() []data.Variable {
	return connCloseMethodGetVariables
}
func (h *ConnCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

func (s *ConnClass) configureMethod() data.Method {
	return &script.Method{Name: "configure", Args: []string{"options"}, Run: func(c data.Context) (data.GetValue, error) {
		if s.source == nil {
			return nil, fmt.Errorf("uninitialized websocket")
		}
		v, e := script.Options(c, 0, map[string]int{"maxMessageBytes": 1048576, "readTimeoutMs": 10000, "writeTimeoutMs": 10000})
		if e == nil {
			s.limits = v
			s.source.SetReadLimit(int64(v["maxMessageBytes"]))
			if s.stop != nil {
				s.stop()
			}
			s.stop = context.AfterFunc(c.GoContext(), func() { s.source.Close() })
		}
		return nil, e
	}}
}
