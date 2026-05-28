package websocket

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewConnClass() data.ClassStmt {
	return &ConnClass{source: nil}
}

func NewConnClassFrom(source *websocket.Conn) data.ClassStmt {
	return &ConnClass{source: source}
}

type ConnClass struct {
	node.Node
	source *websocket.Conn
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
	case "readText":
		return &ConnReadTextMethod{source: s.source}, true
	case "writeText":
		return &ConnWriteTextMethod{source: s.source}, true
	case "close":
		return &ConnCloseMethod{source: s.source}, true
	}
	return nil, false
}

func (s *ConnClass) GetMethods() []data.Method {
	return []data.Method{
		&ConnReadTextMethod{source: s.source},
		&ConnWriteTextMethod{source: s.source},
		&ConnCloseMethod{source: s.source},
	}
}

func (s *ConnClass) GetConstruct() data.Method { return nil }

func (s *ConnClass) GetProperty(name string) (data.Property, bool) { return nil, false }

func (s *ConnClass) GetPropertyList() []data.Property { return []data.Property{} }

type ConnReadTextMethod struct {
	source *websocket.Conn
}

func (h *ConnReadTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return nil, utils.NewThrow(errors.New("websocket 连接未初始化"))
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
func (h *ConnReadTextMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (h *ConnReadTextMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ConnReadTextMethod) GetReturnType() data.Types { return data.NewBaseType("string") }

type ConnWriteTextMethod struct {
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
	if err = h.source.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ConnWriteTextMethod) GetName() string            { return "writeText" }
func (h *ConnWriteTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnWriteTextMethod) GetIsStatic() bool          { return false }
func (h *ConnWriteTextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "message", 0, nil, data.String{}),
	}
}
func (h *ConnWriteTextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "message", 0, nil),
	}
}
func (h *ConnWriteTextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

type ConnCloseMethod struct {
	source *websocket.Conn
}

func (h *ConnCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return nil, utils.NewThrow(errors.New("websocket 连接未初始化"))
	}
	if err := h.source.Close(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ConnCloseMethod) GetName() string            { return "close" }
func (h *ConnCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnCloseMethod) GetIsStatic() bool          { return false }
func (h *ConnCloseMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (h *ConnCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ConnCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
