package http

import (
	nethttp "net/http"

	"github.com/php-any/origami/data"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

func requestClassValue(ctx data.Context) *data.ClassValue {
	return httpfoundation.RequestClassValue(ctx)
}

func valueToAssocMap(v data.Value) (map[string]data.Value, error) {
	return httpfoundation.ValueToAssocMap(v)
}

func isNull(v data.GetValue) bool {
	return httpfoundation.IsNull(v)
}

func serverFromRequest(r *nethttp.Request) map[string]data.Value {
	return httpfoundation.ServerFromRequest(r)
}

func assocMapToArrayValue(m map[string]data.Value) *data.ArrayValue {
	return httpfoundation.AssocMapToArrayValue(m)
}

func requestHeader(ctx data.Context, key string) string {
	return httpfoundation.RequestHeader(ctx, key)
}

func requestPath(ctx data.Context) string { return httpfoundation.RequestPath(ctx) }

func requestBaseURL(ctx data.Context) string { return httpfoundation.RequestBaseURL(ctx) }

func requestHost(ctx data.Context) string { return httpfoundation.RequestHost(ctx) }

func requestHTTPHost(ctx data.Context) string { return httpfoundation.RequestHTTPHost(ctx) }

func requestSecure(ctx data.Context) bool { return httpfoundation.RequestSecure(ctx) }

func GetParamBagMap(cv *data.ClassValue) map[string]data.Value {
	return httpfoundation.GetParamBagMap(cv)
}

func SetParamBagMap(cv *data.ClassValue, params map[string]data.Value) {
	httpfoundation.SetParamBagMap(cv, params)
}

func GetHeaderBagAll(cv *data.ClassValue) map[string][]string {
	return httpfoundation.GetHeaderBagAll(cv)
}

func valueStrings(value data.Value) []string {
	return httpfoundation.ValueStrings(value)
}

func bodyFromRequest(r *nethttp.Request) string {
	return httpfoundation.BodyFromRequest(r)
}
