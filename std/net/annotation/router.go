package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
	http2 "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/utils"
	"net/http"
)

type RegisterRoute struct {
	vm *runtime.TempVM
}

func (r *RegisterRoute) GetValue(ctx data.Context) (v data.GetValue, acl data.Control) {
	routes := r.vm.Cache

	mux := http.NewServeMux()
	for _, rt := range routes {
		mux.HandleFunc(rt.Method+" "+rt.Path, func(w http.ResponseWriter, r *http.Request) {
			request := http2.NewRequestClassFrom(r)
			response := http2.NewResponseWriterClassFrom(w)

			ctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, ctx))
			ctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, ctx))

			v, acl = rt.Target.Call(ctx)
		})
	}

	request, err := utils.ConvertFromIndex[*http.Request](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	response, err := utils.ConvertFromIndex[http.ResponseWriter](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	// 手动调用 mux.ServeHTTP，触发路由和处理
	mux.ServeHTTP(response, request)
	return nil, nil
}
