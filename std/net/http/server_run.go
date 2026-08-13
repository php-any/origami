package http

import (
	"errors"
	"fmt"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
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
	ctx.GetVM().AddShutdownCallback(newServerShutdownCallback(srv))

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, httpsrc.ErrServerClosed) {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ServerRunMethod) GetName() string            { return "run" }
func (h *ServerRunMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerRunMethod) GetIsStatic() bool          { return false }
func (h *ServerRunMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (h *ServerRunMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ServerRunMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
