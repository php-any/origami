package annotation

import (
	"github.com/php-any/origami/data"
	http2 "github.com/php-any/origami/std/net/http"
)

type RegisterRoute struct {
	vm data.VM
}

func (r *RegisterRoute) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return http2.DispatchHTTPRoutes(r.vm, ctx)
}
