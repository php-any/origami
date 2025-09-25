package http

import (
	"errors"
	"net/http"

	"github.com/php-any/origami/data"
)

func newHandler(v data.FuncStmt, ctx data.Context) (Handler, error) {
	if len(v.GetVariables()) < 2 {
		return Handler{}, errors.New("invalid variable definition")
	}

	return Handler{Value: v, Ctx: ctx.CreateContext(v.GetVariables())}, nil
}

type Handler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := f.Ctx.CreateContext(f.Value.GetVariables())

	request := NewRequestClassFrom(r)
	response := NewResponseWriterClassFrom(w)

	ctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, ctx))
	ctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, ctx))

	f.Value.Call(ctx)
}
