package http

import (
	"errors"
	"fmt"
	"net"
	httpsrc "net/http"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/net/internal/script"
	"github.com/php-any/origami/utils"
)

type ServerRunMethod struct {
	server *ServerClass
}

func (h *ServerRunMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	srv := &httpsrc.Server{
		Addr:    h.server.Host + ":" + fmt.Sprintf("%d", h.server.Port),
		Handler: h.server.source,
	}
	if l := h.server.limits; l != nil {
		srv.Handler = limitedHandler(srv.Handler, l)
		srv.MaxHeaderBytes = l["maxHeaderBytes"]
		srv.ReadHeaderTimeout = time.Duration(l["readHeaderTimeoutMs"]) * time.Millisecond
		srv.ReadTimeout = time.Duration(l["readTimeoutMs"]) * time.Millisecond
		srv.WriteTimeout = time.Duration(l["writeTimeoutMs"]) * time.Millisecond
		srv.IdleTimeout = time.Duration(l["idleTimeoutMs"]) * time.Millisecond
	}
	ctx.GetVM().AddShutdownCallback(newServerShutdownCallback(srv))

	ln, err := net.Listen("tcp", net.JoinHostPort(h.server.Host, fmt.Sprint(h.server.Port)))
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if h.server.limits != nil {
		ln = script.LimitListener(ln, h.server.limits["maxConnections"])
	}
	defer ln.Close()
	err = srv.Serve(ln)
	if err != nil && !errors.Is(err, httpsrc.ErrServerClosed) {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ServerRunMethod) GetName() string            { return "run" }
func (h *ServerRunMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerRunMethod) GetIsStatic() bool          { return false }

var serverRunMethodGetParams = []data.GetValue{}

func (h *ServerRunMethod) GetParams() []data.GetValue {
	return serverRunMethodGetParams
}

var serverRunMethodGetVariables = []data.Variable{}

func (h *ServerRunMethod) GetVariables() []data.Variable {
	return serverRunMethodGetVariables
}
func (h *ServerRunMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
