package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"log"
	"net/http"
)

func newServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

type Server struct {
	port int
	addr string
	mux  *http.ServeMux
}

func (server *Server) Construct(ctx data.Context, port int, addr string) {
	server.port = port
	server.addr = addr
	server.mux = http.NewServeMux()
}

func (server *Server) Get(ctx data.Context, path data.StringValue, handler data.FuncValue) {
	server.mux.HandleFunc(path.AsString(), func(w http.ResponseWriter, r *http.Request) {
		request := data.NewClassValue(NewRequestClass(w, r), ctx)
		response := data.NewClassValue(NewResponseClass(w, r), ctx)

		fnCtx := ctx.CreateContext(handler.Value.GetVariables())

		fnCtx.SetVariableValue(data.NewVariable("request", 0, nil), request)
		fnCtx.SetVariableValue(data.NewVariable("response", 1, nil), response)

		_, acl := handler.Call(fnCtx)
		if acl != nil {
			ctx.GetVM().ThrowControl(acl)
		}
	})
}

func (server *Server) Run(ctx data.Context) {
	log.Println(fmt.Sprintf("listen at %s:%d\n ", server.addr, server.port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", server.addr, server.port), server.mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
