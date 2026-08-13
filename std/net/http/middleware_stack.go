package http

import (
	"sort"
	"time"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type middlewareEntry struct {
	priority int
	fn       MiddlewareFunc
}

func applyMiddlewares(final httpsrc.Handler, entries []middlewareEntry) httpsrc.Handler {
	if len(entries) == 0 {
		return final
	}
	sorted := append([]middlewareEntry{}, entries...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].priority < sorted[j].priority
	})
	h := final
	for i := len(sorted) - 1; i >= 0; i-- {
		h = sorted[i].fn(h)
	}
	return h
}

func (s *ServerClass) finalizeHandler(h httpsrc.Handler) httpsrc.Handler {
	h = applyMiddlewares(h, s.middlewares)
	h = withResponseFormatter(s, h)
	return withErrorHandler(s, h)
}

type errorHandlerSlot struct {
	fn  data.FuncStmt
	ctx data.Context
}

func withResponseFormatter(server *ServerClass, next httpsrc.Handler) httpsrc.Handler {
	slot := server.formatHandler
	if slot == nil {
		return next
	}
	return httpsrc.HandlerFunc(func(w httpsrc.ResponseWriter, r *httpsrc.Request) {
		attachRequestFormatter(r, slot)
		next.ServeHTTP(w, r)
	})
}

func withErrorHandler(server *ServerClass, next httpsrc.Handler) httpsrc.Handler {
	if server.errorHandler == nil {
		return next
	}
	return httpsrc.HandlerFunc(func(w httpsrc.ResponseWriter, r *httpsrc.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				invokeErrorHandler(server, w, r, rec)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func invokeErrorHandler(server *ServerClass, w httpsrc.ResponseWriter, r *httpsrc.Request, recovered any) {
	slot := server.errorHandler
	if slot == nil || slot.fn == nil {
		panic(recovered)
	}

	rw, response := beginResponse(w, r)
	defer rw.commitPending()
	r, request := beginRequest(r)

	vars := slot.fn.GetVariables()
	if len(vars) < 3 {
		panic(recovered)
	}

	mctx := slot.ctx.CreateContext(vars)
	mctx.SetVariableValue(vars[0], data.NewProxyValue(request, mctx))
	mctx.SetVariableValue(vars[1], data.NewProxyValue(response, mctx))
	mctx.SetVariableValue(vars[2], errorValueFromRecovered(recovered))

	if _, acl := slot.fn.Call(mctx); acl != nil {
		panic(acl)
	}
}

func errorValueFromRecovered(recovered any) data.Value {
	switch v := recovered.(type) {
	case *data.ThrowValue:
		if err := v.GetError(); err != nil {
			return data.NewStringValue(err.Error())
		}
		return data.NewStringValue("unknown error")
	case data.Control:
		return data.NewStringValue("request error")
	default:
		return data.NewAnyValue(recovered)
	}
}

func cookieFromPHP(ctx data.Context) (*httpsrc.Cookie, error) {
	if cookie, err := utils.ConvertFromIndex[*httpsrc.Cookie](ctx, 0); err == nil {
		return cookie, nil
	}

	name, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, err
	}
	value, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, err
	}

	c := &httpsrc.Cookie{Name: name, Value: value}
	if opt, ok := ctx.GetIndexValue(2); ok {
		applyCookieOptions(c, opt)
	}
	return c, nil
}

func applyCookieOptions(c *httpsrc.Cookie, opt data.GetValue) {
	arr, ok := opt.(*data.ArrayValue)
	if !ok {
		return
	}
	for _, z := range arr.List {
		if z == nil {
			continue
		}
		key := z.Name
		if key == "" {
			continue
		}
		val := z.Value
		switch key {
		case "path", "Path":
			c.Path = val.AsString()
		case "domain", "Domain":
			c.Domain = val.AsString()
		case "expires", "Expires":
			if t, err := utils.Convert[time.Time](val); err == nil {
				c.Expires = t
			}
		case "maxAge", "MaxAge":
			if n, err := utils.Convert[int](val); err == nil {
				c.MaxAge = n
			}
		case "secure", "Secure":
			if b, err := utils.Convert[bool](val); err == nil {
				c.Secure = b
			}
		case "httpOnly", "HttpOnly":
			if b, err := utils.Convert[bool](val); err == nil {
				c.HttpOnly = b
			}
		case "sameSite", "SameSite":
			c.SameSite = parseSameSite(val.AsString())
		}
	}
}

func parseSameSite(v string) httpsrc.SameSite {
	switch v {
	case "Strict", "strict":
		return httpsrc.SameSiteStrictMode
	case "None", "none":
		return httpsrc.SameSiteNoneMode
	default:
		return httpsrc.SameSiteLaxMode
	}
}
