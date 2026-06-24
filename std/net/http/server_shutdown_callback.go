package http

import (
	"context"

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
	_ = f.srv.Shutdown(context.Background())
	return nil, nil
}

func (f *serverShutdownCallback) GetName() string            { return "" }
func (f *serverShutdownCallback) GetParams() []data.GetValue { return nil }
func (f *serverShutdownCallback) GetVariables() []data.Variable {
	return nil
}
