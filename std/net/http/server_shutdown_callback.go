package http

import (
	"context"
	"time"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

type serverShutdownCallback struct {
	srv *httpsrc.Server
}

func newServerShutdownCallback(srv *httpsrc.Server) data.Value {
	return data.NewFuncValue(&serverShutdownCallback{srv: srv})
}

func (f *serverShutdownCallback) Call(ctx data.Context) (data.GetValue, data.Control) {
	shutdown, cancel := context.WithTimeout(ctx.GoContext(), 5*time.Second)
	defer cancel()
	_ = f.srv.Shutdown(shutdown)
	return nil, nil
}

func (f *serverShutdownCallback) GetName() string            { return "" }
func (f *serverShutdownCallback) GetParams() []data.GetValue { return nil }
func (f *serverShutdownCallback) GetVariables() []data.Variable {
	return nil
}
