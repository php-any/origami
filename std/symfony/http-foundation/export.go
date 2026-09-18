package httpfoundation

import (
	"net/http"

	"github.com/php-any/origami/data"
)

// 导出给 std/illuminate/http 等跨包子模块使用的薄包装（本包内部仍用小写实现）。

func PubMethod(name string, params []data.GetValue, vars []data.Variable, ret data.Types, fn func(ctx data.Context) (data.GetValue, data.Control)) data.Method {
	return pubMethod(name, params, vars, ret, fn)
}

func Param(name string, index int, def data.GetValue, ty data.Types) data.GetValue {
	return param(name, index, def, ty)
}

func Variable(name string, index int, ty data.Types) data.Variable {
	return variable(name, index, ty)
}

func PublicProp(name string, def data.GetValue) data.Property {
	return publicProp(name, def)
}

func ResponseClassValue(ctx data.Context) *data.ClassValue {
	return responseClassValue(ctx)
}

func ResponseSelf(ctx data.Context) data.GetValue {
	return responseSelf(ctx)
}

func ResponseContent(cv *data.ClassValue) string {
	return responseContent(cv)
}

func ResponseStatusCode(cv *data.ClassValue) int {
	return responseStatusCode(cv)
}

func ResponseStatusText(cv *data.ClassValue) string {
	return responseStatusText(cv)
}

func ResponseHeaders(cv *data.ClassValue) *data.ClassValue {
	return responseHeaders(cv)
}

func ResponseSetProp(cv *data.ClassValue, name string, value data.Value) {
	responseSetProp(cv, name, value)
}

func HeadersMapFromValue(v data.Value) map[string][]string {
	return headersMapFromValue(v)
}

func HeadersFromRequest(request *http.Request) map[string][]string {
	return headersFromRequest(request)
}

func CreateResponseHeaders(ctx data.Context, v data.Value) *data.ClassValue {
	return createResponseHeaders(ctx, v)
}

func ApplyStatusCode(cv *data.ClassValue, code int, text data.Value) data.Control {
	return applyStatusCode(cv, code, text)
}

func RespHeaderSet(headers *data.ClassValue, key string, values []string, replace bool) {
	respHeaderSet(headers, key, values, replace)
}

func RespHeaderRemove(headers *data.ClassValue, key string) {
	respHeaderRemove(headers, key)
}

func CallObjMethod(obj *data.ClassValue, name string, args ...data.Value) (data.GetValue, data.Control) {
	return callObjMethod(obj, name, args...)
}

func ThrowNamed(name string, format string, args ...any) data.Control {
	return throwNamed(name, format, args...)
}

func IntParam(ctx data.Context, index int, def int) int {
	return intParam(ctx, index, def)
}

func BoolParam(ctx data.Context, index int, def bool) bool {
	return boolParam(ctx, index, def)
}

func OptionalStringParam(ctx data.Context, index int) (string, bool, bool) {
	return optionalStringParam(ctx, index)
}

func ValueToAssocMap(v data.Value) (map[string]data.Value, error) {
	return valueToAssocMap(v)
}

func SetBagCookie(src *ResponseHeaderBagData, cookie *BagCookie) {
	setBagCookie(src, cookie)
}

// SetMethodModifier 将 PubMethod 创建的方法改为 protected（Illuminate ResponseTrait）。
func SetMethodModifier(m data.Method, mod data.Modifier) bool {
	bm, ok := m.(*bagMethod)
	if !ok {
		return false
	}
	bm.modifier = mod
	return true
}

func RequestClassValue(ctx data.Context) *data.ClassValue {
	return requestClassValue(ctx)
}

func ServerFromRequest(request *http.Request) map[string]data.Value {
	return serverFromRequest(request)
}

func IsNull(v data.GetValue) bool {
	return isNull(v)
}

func AssocMapToArrayValue(m map[string]data.Value) *data.ArrayValue {
	return assocMapToArrayValue(m)
}

func RequestHeader(ctx data.Context, key string) string { return requestHeader(ctx, key) }
func RequestPath(ctx data.Context) string               { return requestPath(ctx) }
func RequestBaseURL(ctx data.Context) string            { return requestBaseURL(ctx) }
func RequestHost(ctx data.Context) string               { return requestHost(ctx) }
func RequestHTTPHost(ctx data.Context) string           { return requestHTTPHost(ctx) }
func RequestSecure(ctx data.Context) bool               { return requestSecure(ctx) }

func BodyFromRequest(request *http.Request) string { return bodyFromRequest(request) }

func ValueStrings(value data.Value) []string { return valueStrings(value) }

func RequestCachedContent(cv *data.ClassValue) string {
	if cv == nil {
		return ""
	}
	if class, ok := cv.Class.(*SymfonyRequestClass); ok && class.source != nil {
		return class.source.content
	}
	return ""
}

// FqnSymfonyResponse 供 Illuminate\Http\Response 声明父类。
const FqnSymfonyResponse = fqnSymfonyResponse

const (
	FqnFile               = fqnFile
	FqnUploadedFile       = fqnUploadedFile
	FqnRequestStack       = fqnRequestStack
	FqnStreamedResponse   = fqnStreamedResponse
	FqnBinaryFileResponse = fqnBinaryFileResponse
)
