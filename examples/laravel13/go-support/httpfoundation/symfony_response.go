package httpfoundation

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnSymfonyResponse = "Symfony\\Component\\HttpFoundation\\Response"

var responseStatusTexts = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Content Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Content",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

var responseHTTPConstants = map[string]int{
	"HTTP_CONTINUE":                             100,
	"HTTP_SWITCHING_PROTOCOLS":                  101,
	"HTTP_PROCESSING":                           102,
	"HTTP_EARLY_HINTS":                          103,
	"HTTP_OK":                                   200,
	"HTTP_CREATED":                              201,
	"HTTP_ACCEPTED":                             202,
	"HTTP_NON_AUTHORITATIVE_INFORMATION":        203,
	"HTTP_NO_CONTENT":                           204,
	"HTTP_RESET_CONTENT":                        205,
	"HTTP_PARTIAL_CONTENT":                      206,
	"HTTP_MULTI_STATUS":                         207,
	"HTTP_ALREADY_REPORTED":                     208,
	"HTTP_IM_USED":                              226,
	"HTTP_MULTIPLE_CHOICES":                     300,
	"HTTP_MOVED_PERMANENTLY":                    301,
	"HTTP_FOUND":                                302,
	"HTTP_SEE_OTHER":                            303,
	"HTTP_NOT_MODIFIED":                         304,
	"HTTP_USE_PROXY":                            305,
	"HTTP_RESERVED":                             306,
	"HTTP_TEMPORARY_REDIRECT":                   307,
	"HTTP_PERMANENTLY_REDIRECT":                 308,
	"HTTP_BAD_REQUEST":                          400,
	"HTTP_UNAUTHORIZED":                         401,
	"HTTP_PAYMENT_REQUIRED":                     402,
	"HTTP_FORBIDDEN":                            403,
	"HTTP_NOT_FOUND":                            404,
	"HTTP_METHOD_NOT_ALLOWED":                   405,
	"HTTP_NOT_ACCEPTABLE":                       406,
	"HTTP_PROXY_AUTHENTICATION_REQUIRED":        407,
	"HTTP_REQUEST_TIMEOUT":                      408,
	"HTTP_CONFLICT":                             409,
	"HTTP_GONE":                                 410,
	"HTTP_LENGTH_REQUIRED":                      411,
	"HTTP_PRECONDITION_FAILED":                  412,
	"HTTP_REQUEST_ENTITY_TOO_LARGE":             413,
	"HTTP_REQUEST_URI_TOO_LONG":                 414,
	"HTTP_UNSUPPORTED_MEDIA_TYPE":               415,
	"HTTP_REQUESTED_RANGE_NOT_SATISFIABLE":      416,
	"HTTP_EXPECTATION_FAILED":                   417,
	"HTTP_I_AM_A_TEAPOT":                        418,
	"HTTP_MISDIRECTED_REQUEST":                  421,
	"HTTP_UNPROCESSABLE_ENTITY":                 422,
	"HTTP_LOCKED":                               423,
	"HTTP_FAILED_DEPENDENCY":                    424,
	"HTTP_TOO_EARLY":                            425,
	"HTTP_UPGRADE_REQUIRED":                     426,
	"HTTP_PRECONDITION_REQUIRED":                428,
	"HTTP_TOO_MANY_REQUESTS":                    429,
	"HTTP_REQUEST_HEADER_FIELDS_TOO_LARGE":      431,
	"HTTP_UNAVAILABLE_FOR_LEGAL_REASONS":        451,
	"HTTP_INTERNAL_SERVER_ERROR":                500,
	"HTTP_NOT_IMPLEMENTED":                      501,
	"HTTP_BAD_GATEWAY":                          502,
	"HTTP_SERVICE_UNAVAILABLE":                  503,
	"HTTP_GATEWAY_TIMEOUT":                      504,
	"HTTP_VERSION_NOT_SUPPORTED":                505,
	"HTTP_VARIANT_ALSO_NEGOTIATES_EXPERIMENTAL": 506,
	"HTTP_INSUFFICIENT_STORAGE":                 507,
	"HTTP_LOOP_DETECTED":                        508,
	"HTTP_NOT_EXTENDED":                         510,
	"HTTP_NETWORK_AUTHENTICATION_REQUIRED":      511,
}

// SymfonyResponseClass 实现 Symfony\Component\HttpFoundation\Response。
type SymfonyResponseClass struct {
	node.Node
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

// NewSymfonyResponseClass 导出 Symfony Response 构造器。
func NewSymfonyResponseClass() data.ClassStmt {
	c := &SymfonyResponseClass{
		properties: []data.Property{
			publicProp("headers", nil),
			protectedProp("content", data.NewStringValue("")),
			protectedProp("version", data.NewStringValue("1.0")),
			protectedProp("statusCode", data.NewIntValue(200)),
			publicProp("statusText", data.NewStringValue("OK")),
			protectedProp("charset", data.NewNullValue()),
		},
	}
	c.methods, c.methodList = symfonyResponseMethods()
	return c
}

func (c *SymfonyResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *SymfonyResponseClass) GetName() string         { return fqnSymfonyResponse }
func (c *SymfonyResponseClass) GetExtend() *string      { return nil }
func (c *SymfonyResponseClass) GetImplements() []string { return nil }
func (c *SymfonyResponseClass) GetConstruct() data.Method {
	return c.methods["__construct"]
}
func (c *SymfonyResponseClass) GetPropertyList() []data.Property { return c.properties }
func (c *SymfonyResponseClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *SymfonyResponseClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *SymfonyResponseClass) GetMethods() []data.Method { return c.methodList }

// GetStaticProperty 暴露 HTTP_* 状态码常量与 $statusTexts。
func (c *SymfonyResponseClass) GetStaticProperty(name string) (data.Value, bool) {
	if code, ok := responseHTTPConstants[name]; ok {
		return data.NewIntValue(code), true
	}
	if name == "statusTexts" {
		list := make([]*data.ZVal, 0, len(responseStatusTexts))
		for code, text := range responseStatusTexts {
			list = append(list, &data.ZVal{Name: data.IntArrayKeyName(code), Value: data.NewStringValue(text)})
		}
		return &data.ArrayValue{List: list}, true
	}
	return nil, false
}

func symfonyResponseMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("content", 0, data.NewStringValue(""), nil),
				param("status", 1, data.NewIntValue(200), nil),
				param("headers", 2, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				variable("content", 0, nil),
				variable("status", 1, nil),
				variable("headers", 2, nil),
			},
			nil, symfonyResponseConstruct),
		pubMethod("__toString", nil, nil, data.NewBaseType("string"), symfonyResponseToString),
		pubMethod("prepare",
			[]data.GetValue{param("request", 0, nil, nil)},
			[]data.Variable{variable("request", 0, nil)},
			nil, symfonyResponsePrepare),
		pubMethod("sendHeaders",
			[]data.GetValue{param("statusCode", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("statusCode", 0, nil)},
			nil, symfonyResponseSendHeaders),
		pubMethod("sendContent", nil, nil, nil, symfonyResponseSendContent),
		pubMethod("send",
			[]data.GetValue{param("flush", 0, data.NewBoolValue(true), nil)},
			[]data.Variable{variable("flush", 0, nil)},
			nil, symfonyResponseSend),
		pubMethod("setContent",
			[]data.GetValue{param("content", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("content", 0, nil)},
			nil, symfonyResponseSetContent),
		pubMethod("getContent", nil, nil, nil, symfonyResponseGetContent),
		pubMethod("setProtocolVersion",
			[]data.GetValue{param("version", 0, nil, nil)},
			[]data.Variable{variable("version", 0, nil)},
			nil, symfonyResponseSetProtocolVersion),
		pubMethod("getProtocolVersion", nil, nil, data.NewBaseType("string"), symfonyResponseGetProtocolVersion),
		pubMethod("setStatusCode",
			[]data.GetValue{
				param("code", 0, nil, nil),
				param("text", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("code", 0, nil),
				variable("text", 1, nil),
			},
			nil, symfonyResponseSetStatusCode),
		pubMethod("getStatusCode", nil, nil, data.NewBaseType("int"), symfonyResponseGetStatusCode),
		pubMethod("setCharset",
			[]data.GetValue{param("charset", 0, nil, nil)},
			[]data.Variable{variable("charset", 0, nil)},
			nil, symfonyResponseSetCharset),
		pubMethod("getCharset", nil, nil, nil, symfonyResponseGetCharset),
		pubMethod("isCacheable", nil, nil, data.NewBaseType("bool"), symfonyResponseIsCacheable),
		pubMethod("isFresh", nil, nil, data.NewBaseType("bool"), symfonyResponseIsFresh),
		pubMethod("isValidateable", nil, nil, data.NewBaseType("bool"), symfonyResponseIsValidateable),
		pubMethod("setPrivate", nil, nil, nil, symfonyResponseSetPrivate),
		pubMethod("setPublic", nil, nil, nil, symfonyResponseSetPublic),
		pubMethod("setImmutable",
			[]data.GetValue{param("immutable", 0, data.NewBoolValue(true), nil)},
			[]data.Variable{variable("immutable", 0, nil)},
			nil, symfonyResponseSetImmutable),
		pubMethod("isImmutable", nil, nil, data.NewBaseType("bool"), symfonyResponseIsImmutable),
		pubMethod("mustRevalidate", nil, nil, data.NewBaseType("bool"), symfonyResponseMustRevalidate),
		pubMethod("getDate", nil, nil, nil, symfonyResponseGetDate),
		pubMethod("setDate",
			[]data.GetValue{param("date", 0, nil, nil)},
			[]data.Variable{variable("date", 0, nil)},
			nil, symfonyResponseSetDate),
		pubMethod("getAge", nil, nil, data.NewBaseType("int"), symfonyResponseGetAge),
		pubMethod("expire", nil, nil, nil, symfonyResponseExpire),
		pubMethod("getExpires", nil, nil, nil, symfonyResponseGetExpires),
		pubMethod("setExpires",
			[]data.GetValue{param("date", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("date", 0, nil)},
			nil, symfonyResponseSetExpires),
		pubMethod("getMaxAge", nil, nil, nil, symfonyResponseGetMaxAge),
		pubMethod("setMaxAge",
			[]data.GetValue{param("value", 0, nil, nil)},
			[]data.Variable{variable("value", 0, nil)},
			nil, symfonyResponseSetMaxAge),
		pubMethod("setStaleIfError",
			[]data.GetValue{param("value", 0, nil, nil)},
			[]data.Variable{variable("value", 0, nil)},
			nil, symfonyResponseSetStaleIfError),
		pubMethod("setStaleWhileRevalidate",
			[]data.GetValue{param("value", 0, nil, nil)},
			[]data.Variable{variable("value", 0, nil)},
			nil, symfonyResponseSetStaleWhileRevalidate),
		pubMethod("setSharedMaxAge",
			[]data.GetValue{param("value", 0, nil, nil)},
			[]data.Variable{variable("value", 0, nil)},
			nil, symfonyResponseSetSharedMaxAge),
		pubMethod("getTtl", nil, nil, nil, symfonyResponseGetTtl),
		pubMethod("setTtl",
			[]data.GetValue{param("seconds", 0, nil, nil)},
			[]data.Variable{variable("seconds", 0, nil)},
			nil, symfonyResponseSetTtl),
		pubMethod("setClientTtl",
			[]data.GetValue{param("seconds", 0, nil, nil)},
			[]data.Variable{variable("seconds", 0, nil)},
			nil, symfonyResponseSetClientTtl),
		pubMethod("getLastModified", nil, nil, nil, symfonyResponseGetLastModified),
		pubMethod("setLastModified",
			[]data.GetValue{param("date", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("date", 0, nil)},
			nil, symfonyResponseSetLastModified),
		pubMethod("getEtag", nil, nil, nil, symfonyResponseGetEtag),
		pubMethod("setEtag",
			[]data.GetValue{
				param("etag", 0, data.NewNullValue(), nil),
				param("weak", 1, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("etag", 0, nil),
				variable("weak", 1, nil),
			},
			nil, symfonyResponseSetEtag),
		pubMethod("setCache",
			[]data.GetValue{param("options", 0, nil, nil)},
			[]data.Variable{variable("options", 0, nil)},
			nil, symfonyResponseSetCache),
		pubMethod("setNotModified", nil, nil, nil, symfonyResponseSetNotModified),
		pubMethod("hasVary", nil, nil, data.NewBaseType("bool"), symfonyResponseHasVary),
		pubMethod("getVary", nil, nil, data.NewBaseType("array"), symfonyResponseGetVary),
		pubMethod("setVary",
			[]data.GetValue{
				param("headers", 0, nil, nil),
				param("replace", 1, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				variable("headers", 0, nil),
				variable("replace", 1, nil),
			},
			nil, symfonyResponseSetVary),
		pubMethod("isNotModified",
			[]data.GetValue{param("request", 0, nil, nil)},
			[]data.Variable{variable("request", 0, nil)},
			data.NewBaseType("bool"), symfonyResponseIsNotModified),
		pubMethod("isInvalid", nil, nil, data.NewBaseType("bool"), symfonyResponseIsInvalid),
		pubMethod("isInformational", nil, nil, data.NewBaseType("bool"), symfonyResponseIsInformational),
		pubMethod("isSuccessful", nil, nil, data.NewBaseType("bool"), symfonyResponseIsSuccessful),
		pubMethod("isRedirection", nil, nil, data.NewBaseType("bool"), symfonyResponseIsRedirection),
		pubMethod("isClientError", nil, nil, data.NewBaseType("bool"), symfonyResponseIsClientError),
		pubMethod("isServerError", nil, nil, data.NewBaseType("bool"), symfonyResponseIsServerError),
		pubMethod("isOk", nil, nil, data.NewBaseType("bool"), symfonyResponseIsOk),
		pubMethod("isForbidden", nil, nil, data.NewBaseType("bool"), symfonyResponseIsForbidden),
		pubMethod("isNotFound", nil, nil, data.NewBaseType("bool"), symfonyResponseIsNotFound),
		pubMethod("isRedirect",
			[]data.GetValue{param("location", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("location", 0, nil)},
			data.NewBaseType("bool"), symfonyResponseIsRedirect),
		pubMethod("isEmpty", nil, nil, data.NewBaseType("bool"), symfonyResponseIsEmpty),
		pubMethod("closeOutputBuffers",
			[]data.GetValue{
				param("targetLevel", 0, nil, nil),
				param("flush", 1, nil, nil),
			},
			[]data.Variable{
				variable("targetLevel", 0, nil),
				variable("flush", 1, nil),
			},
			nil, symfonyResponseCloseOutputBuffers),
		pubMethod("setContentSafe",
			[]data.GetValue{param("safe", 0, data.NewBoolValue(true), nil)},
			[]data.Variable{variable("safe", 0, nil)},
			nil, symfonyResponseSetContentSafe),
	}
	// mark closeOutputBuffers static
	for i, m := range list {
		if m.GetName() == "closeOutputBuffers" {
			bm := m.(*bagMethod)
			bm.static = true
			list[i] = bm
		}
	}
	out := make(map[string]data.Method, len(list))
	for _, m := range list {
		out[m.GetName()] = m
	}
	return out, list
}

func symfonyResponseConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Response::__construct: invalid context"))
	}
	content, _ := ctx.GetIndexValue(0)
	status := intParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)

	headers := createResponseHeaders(ctx, headersVal)
	responseSetProp(cv, "headers", headers)

	if content == nil || func() bool { _, ok := content.(*data.NullValue); return ok }() {
		responseSetProp(cv, "content", data.NewStringValue(""))
	} else if rendered, ok := tryRender(ctx, content); ok {
		// Laravel 的 View 等 Renderable 对象必须先执行 render()，
		// 不能使用 ClassValue.AsString() 的对象调试文本。
		responseSetProp(cv, "content", data.NewStringValue(rendered))
	} else {
		responseSetProp(cv, "content", data.NewStringValue(content.AsString()))
	}

	if ctl := applyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	responseSetProp(cv, "version", data.NewStringValue("1.0"))
	return nil, nil
}

func applyStatusCode(cv *data.ClassValue, code int, text data.Value) data.Control {
	responseSetProp(cv, "statusCode", data.NewIntValue(code))
	if code < 100 || code >= 600 {
		return throwNamed("InvalidArgumentException", `The HTTP status code "%d" is not valid.`, code)
	}
	if text != nil {
		if _, ok := text.(*data.NullValue); !ok {
			responseSetProp(cv, "statusText", data.NewStringValue(text.AsString()))
			return nil
		}
	}
	responseSetProp(cv, "statusText", data.NewStringValue(statusTextFor(code)))
	return nil
}

func symfonyResponseToString(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headers := responseHeaders(cv)
	headerStr := ""
	if headers != nil {
		headerStr = headers.AsString()
	}
	line := fmt.Sprintf("HTTP/%s %d %s\r\n%s\r\n%s",
		responseVersion(cv), responseStatusCode(cv), responseStatusText(cv), headerStr, responseContent(cv))
	return data.NewStringValue(line), nil
}

func symfonyResponsePrepare(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headers := responseHeaders(cv)
	reqVal, _ := ctx.GetIndexValue(0)
	req, _ := reqVal.(*data.ClassValue)

	code := responseStatusCode(cv)
	informational := code >= 100 && code < 200
	empty := code == 204 || code == 304

	if informational || empty {
		responseSetProp(cv, "content", data.NewStringValue(""))
		respHeaderRemove(headers, "Content-Type")
		respHeaderRemove(headers, "Content-Length")
	} else {
		if !respHeaderHas(headers, "Content-Type") && req != nil {
			format, _ := callObjMethod(req, "getRequestFormat", data.NewNullValue())
			if format != nil {
				if _, isNull := format.(*data.NullValue); !isNull {
					mime, _ := callObjMethod(req, "getMimeType", data.NewStringValue(format.(data.Value).AsString()))
					if mime != nil {
						if _, isNull := mime.(*data.NullValue); !isNull && mime.(data.Value).AsString() != "" {
							respHeaderSet(headers, "Content-Type", []string{mime.(data.Value).AsString()}, true)
						}
					}
				}
			}
		}
		charset, ok := responseCharset(cv)
		if !ok || charset == "" {
			charset = "utf-8"
		}
		if !respHeaderHas(headers, "Content-Type") {
			respHeaderSet(headers, "Content-Type", []string{"text/html; charset=" + charset}, true)
		} else {
			ct := respHeaderGet(headers, "Content-Type", "")
			if strings.HasPrefix(strings.ToLower(ct), "text/") && !strings.Contains(strings.ToLower(ct), "charset") {
				respHeaderSet(headers, "Content-Type", []string{ct + "; charset=" + charset}, true)
			}
		}
		if respHeaderHas(headers, "Transfer-Encoding") {
			respHeaderRemove(headers, "Content-Length")
		}
		if req != nil {
			isHead, _ := callObjMethod(req, "isMethod", data.NewStringValue("HEAD"))
			if valueIsTruthy(isHead) {
				length := respHeaderGet(headers, "Content-Length", "")
				responseSetProp(cv, "content", data.NewStringValue(""))
				if length != "" {
					respHeaderSet(headers, "Content-Length", []string{length}, true)
				}
			}
		}
	}

	if req != nil {
		serverProp, _ := req.GetProperty("server")
		proto := "HTTP/1.1"
		if serverCV, ok := serverProp.(*data.ClassValue); ok {
			got, _ := callObjMethod(serverCV, "get", data.NewStringValue("SERVER_PROTOCOL"), data.NewStringValue("HTTP/1.0"))
			if got != nil {
				proto = got.(data.Value).AsString()
			}
		}
		if proto != "HTTP/1.0" {
			responseSetProp(cv, "version", data.NewStringValue("1.1"))
		}
	}

	if responseVersion(cv) == "1.0" && strings.Contains(respHeaderGet(headers, "Cache-Control", ""), "no-cache") {
		respHeaderSet(headers, "Pragma", []string{"no-cache"}, true)
		respHeaderSet(headers, "Expires", []string{"-1"}, true)
	}

	return responseSelf(ctx), nil
}

func symfonyResponseSendHeaders(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headers := responseHeaders(cv)
	statusCode := responseStatusCode(cv)
	if n, ok := optionalIntParam(ctx, 0); ok {
		statusCode = n
	}

	informational := statusCode >= 100 && statusCode < 200
	if informational {
		return responseSelf(ctx), nil
	}

	// Prefer preserve-case listing when available.
	all := GetHeaderBagAll(headers)
	if rh := ResponseHeaderBagFrom(headers); rh != nil {
		rh.mu.RLock()
		for lk, vs := range rh.headers {
			name := rh.headerNames[lk]
			if name == "" {
				name = lk
			}
			if strings.EqualFold(name, "set-cookie") {
				continue
			}
			replace := strings.EqualFold(name, "Content-Type")
			for _, v := range vs {
				val := ""
				if v != nil {
					val = *v
				}
				SendHeaderContext(ctx, name+": "+val, replace, statusCode)
				replace = false
			}
		}
		// cookies
		for _, byPath := range rh.cookies {
			for _, byName := range byPath {
				for _, cookie := range byName {
					if cookie != nil {
						SendHeaderContext(ctx, "Set-Cookie: "+cookie.String(), false, statusCode)
					}
				}
			}
		}
		rh.mu.RUnlock()
	} else {
		for name, values := range all {
			if strings.EqualFold(name, "set-cookie") {
				continue
			}
			replace := strings.EqualFold(name, "Content-Type")
			for _, value := range values {
				SendHeaderContext(ctx, name+": "+value, replace, statusCode)
				replace = false
			}
		}
	}

	SendHeaderContext(ctx, fmt.Sprintf("HTTP/%s %d %s", responseVersion(cv), statusCode, responseStatusText(cv)), true, statusCode)
	return responseSelf(ctx), nil
}

func symfonyResponseSendContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	data.EmitOutput(ctx, responseContent(cv))
	return responseSelf(ctx), nil
}

func symfonyResponseSend(ctx data.Context) (data.GetValue, data.Control) {
	if _, ctl := symfonyResponseSendHeaders(ctx); ctl != nil {
		return nil, ctl
	}
	if _, ctl := symfonyResponseSendContent(ctx); ctl != nil {
		return nil, ctl
	}
	return responseSelf(ctx), nil
}

func symfonyResponseSetContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	content, _ := ctx.GetIndexValue(0)
	if content == nil {
		responseSetProp(cv, "content", data.NewStringValue(""))
	} else if _, ok := content.(*data.NullValue); ok {
		responseSetProp(cv, "content", data.NewStringValue(""))
	} else if rendered, ok := tryRender(ctx, content); ok {
		responseSetProp(cv, "content", data.NewStringValue(rendered))
	} else {
		responseSetProp(cv, "content", data.NewStringValue(content.AsString()))
	}
	return responseSelf(ctx), nil
}

func symfonyResponseGetContent(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(responseContent(responseClassValue(ctx))), nil
}

func symfonyResponseSetProtocolVersion(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := ctx.GetIndexValue(0)
	if v != nil {
		responseSetProp(cv, "version", data.NewStringValue(v.AsString()))
	}
	return responseSelf(ctx), nil
}

func symfonyResponseGetProtocolVersion(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(responseVersion(responseClassValue(ctx))), nil
}

func symfonyResponseSetStatusCode(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	code := intParam(ctx, 0, 200)
	text, _ := ctx.GetIndexValue(1)
	if ctl := applyStatusCode(cv, code, text); ctl != nil {
		return nil, ctl
	}
	return responseSelf(ctx), nil
}

func symfonyResponseGetStatusCode(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(responseStatusCode(responseClassValue(ctx))), nil
}

func symfonyResponseSetCharset(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := ctx.GetIndexValue(0)
	if v != nil {
		responseSetProp(cv, "charset", data.NewStringValue(v.AsString()))
	}
	return responseSelf(ctx), nil
}

func symfonyResponseGetCharset(ctx data.Context) (data.GetValue, data.Control) {
	s, ok := responseCharset(responseClassValue(ctx))
	if !ok {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(s), nil
}

func symfonyResponseIsCacheable(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	code := responseStatusCode(cv)
	switch code {
	case 200, 203, 300, 301, 302, 404, 410:
	default:
		return data.NewBoolValue(false), nil
	}
	headers := responseHeaders(cv)
	if respHasCC(headers, "no-store") {
		return data.NewBoolValue(false), nil
	}
	if _, ok := respGetCC(headers, "private"); ok {
		return data.NewBoolValue(false), nil
	}
	validateable := respHeaderHas(headers, "Last-Modified") || respHeaderHas(headers, "ETag")
	fresh := false
	if ttl, ctl := computeTtl(cv); ctl == nil && ttl != nil && *ttl > 0 {
		fresh = true
	}
	return data.NewBoolValue(validateable || fresh), nil
}

func symfonyResponseIsFresh(ctx data.Context) (data.GetValue, data.Control) {
	ttl, ctl := computeTtl(responseClassValue(ctx))
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(ttl != nil && *ttl > 0), nil
}

func symfonyResponseIsValidateable(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	return data.NewBoolValue(respHeaderHas(headers, "Last-Modified") || respHeaderHas(headers, "ETag")), nil
}

func symfonyResponseSetPrivate(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	respRemoveCC(headers, "public")
	respAddCC(headers, "private", true)
	return responseSelf(ctx), nil
}

func symfonyResponseSetPublic(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	respAddCC(headers, "public", true)
	respRemoveCC(headers, "private")
	return responseSelf(ctx), nil
}

func symfonyResponseSetImmutable(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	if boolParam(ctx, 0, true) {
		respAddCC(headers, "immutable", true)
	} else {
		respRemoveCC(headers, "immutable")
	}
	return responseSelf(ctx), nil
}

func symfonyResponseIsImmutable(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(respHasCC(responseHeaders(responseClassValue(ctx)), "immutable")), nil
}

func symfonyResponseMustRevalidate(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	return data.NewBoolValue(respHasCC(headers, "must-revalidate") || respHasCC(headers, "proxy-revalidate")), nil
}

func dateStringValue(t time.Time) data.Value {
	return data.NewStringValue(respFormatHTTPDate(t))
}

func headerDateOrNull(headers *data.ClassValue, key string) (data.GetValue, data.Control) {
	if !respHeaderHas(headers, key) {
		return data.NewNullValue(), nil
	}
	s := respHeaderGet(headers, key, "")
	if t, ok := respParseHTTPDate(s); ok {
		return dateStringValue(t), nil
	}
	return data.NewNullValue(), nil
}

func symfonyResponseGetDate(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	if !respHeaderHas(headers, "Date") {
		now := time.Now().UTC()
		respHeaderSet(headers, "Date", []string{respFormatHTTPDate(now)}, true)
		return dateStringValue(now), nil
	}
	return headerDateOrNull(headers, "Date")
}

func parseDateArg(v data.Value) (time.Time, bool) {
	if v == nil {
		return time.Time{}, false
	}
	if _, ok := v.(*data.NullValue); ok {
		return time.Time{}, false
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return time.Unix(int64(n), 0).UTC(), true
		}
	}
	s := v.AsString()
	if t, ok := respParseHTTPDate(s); ok {
		return t, true
	}
	if n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
		return time.Unix(n, 0).UTC(), true
	}
	return time.Time{}, false
}

func symfonyResponseSetDate(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := ctx.GetIndexValue(0)
	t, ok := parseDateArg(v)
	if !ok {
		t = time.Now().UTC()
	}
	respHeaderSet(responseHeaders(cv), "Date", []string{respFormatHTTPDate(t)}, true)
	return responseSelf(ctx), nil
}

func computeAge(cv *data.ClassValue) int {
	headers := responseHeaders(cv)
	if age := respHeaderGet(headers, "Age", ""); age != "" {
		if n, err := strconv.Atoi(age); err == nil {
			return n
		}
	}
	dateStr := respHeaderGet(headers, "Date", "")
	if t, ok := respParseHTTPDate(dateStr); ok {
		a := int(time.Now().UTC().Sub(t).Seconds())
		if a < 0 {
			return 0
		}
		return a
	}
	return 0
}

func symfonyResponseGetAge(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(computeAge(responseClassValue(ctx))), nil
}

func symfonyResponseExpire(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	ttl, _ := computeTtl(cv)
	if ttl != nil && *ttl > 0 {
		maxAge, _ := computeMaxAge(cv)
		if maxAge != nil {
			respHeaderSet(responseHeaders(cv), "Age", []string{strconv.Itoa(*maxAge)}, true)
			respHeaderRemove(responseHeaders(cv), "Expires")
		}
	}
	return responseSelf(ctx), nil
}

func symfonyResponseGetExpires(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	if !respHeaderHas(headers, "Expires") {
		return data.NewNullValue(), nil
	}
	s := respHeaderGet(headers, "Expires", "")
	if t, ok := respParseHTTPDate(s); ok {
		return dateStringValue(t), nil
	}
	// invalid date => past
	return dateStringValue(time.Now().UTC().Add(-172800 * time.Second)), nil
}

func symfonyResponseSetExpires(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := ctx.GetIndexValue(0)
	if v == nil || func() bool { _, ok := v.(*data.NullValue); return ok }() {
		respHeaderRemove(responseHeaders(cv), "Expires")
		return responseSelf(ctx), nil
	}
	t, ok := parseDateArg(v)
	if !ok {
		t = time.Now().UTC()
	}
	respHeaderSet(responseHeaders(cv), "Expires", []string{respFormatHTTPDate(t)}, true)
	return responseSelf(ctx), nil
}

func computeMaxAge(cv *data.ClassValue) (*int, data.Control) {
	headers := responseHeaders(cv)
	if v, ok := respGetCC(headers, "s-maxage"); ok {
		switch t := v.(type) {
		case int:
			return &t, nil
		case string:
			if n, err := strconv.Atoi(t); err == nil {
				return &n, nil
			}
		}
	}
	if v, ok := respGetCC(headers, "max-age"); ok {
		switch t := v.(type) {
		case int:
			return &t, nil
		case string:
			if n, err := strconv.Atoi(t); err == nil {
				return &n, nil
			}
		}
	}
	if respHeaderHas(headers, "Expires") {
		expStr := respHeaderGet(headers, "Expires", "")
		dateStr := respHeaderGet(headers, "Date", "")
		exp, ok1 := respParseHTTPDate(expStr)
		date, ok2 := respParseHTTPDate(dateStr)
		if !ok2 {
			date = time.Now().UTC()
		}
		if ok1 {
			n := int(exp.Sub(date).Seconds())
			if n < 0 {
				n = 0
			}
			return &n, nil
		}
	}
	return nil, nil
}

func computeTtl(cv *data.ClassValue) (*int, data.Control) {
	maxAge, ctl := computeMaxAge(cv)
	if ctl != nil {
		return nil, ctl
	}
	if maxAge == nil {
		return nil, nil
	}
	n := *maxAge - computeAge(cv)
	if n < 0 {
		n = 0
	}
	return &n, nil
}

func symfonyResponseGetMaxAge(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := computeMaxAge(responseClassValue(ctx))
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(*v), nil
}

func symfonyResponseSetMaxAge(ctx data.Context) (data.GetValue, data.Control) {
	respAddCC(responseHeaders(responseClassValue(ctx)), "max-age", intParam(ctx, 0, 0))
	return responseSelf(ctx), nil
}

func symfonyResponseSetStaleIfError(ctx data.Context) (data.GetValue, data.Control) {
	respAddCC(responseHeaders(responseClassValue(ctx)), "stale-if-error", intParam(ctx, 0, 0))
	return responseSelf(ctx), nil
}

func symfonyResponseSetStaleWhileRevalidate(ctx data.Context) (data.GetValue, data.Control) {
	respAddCC(responseHeaders(responseClassValue(ctx)), "stale-while-revalidate", intParam(ctx, 0, 0))
	return responseSelf(ctx), nil
}

func symfonyResponseSetSharedMaxAge(ctx data.Context) (data.GetValue, data.Control) {
	if _, ctl := symfonyResponseSetPublic(ctx); ctl != nil {
		return nil, ctl
	}
	respAddCC(responseHeaders(responseClassValue(ctx)), "s-maxage", intParam(ctx, 0, 0))
	return responseSelf(ctx), nil
}

func symfonyResponseGetTtl(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := computeTtl(responseClassValue(ctx))
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(*v), nil
}

func symfonyResponseSetTtl(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	seconds := intParam(ctx, 0, 0)
	// setSharedMaxAge(age + seconds)
	age := computeAge(cv)
	ctx2 := cv.CreateContext([]data.Variable{variable("value", 0, nil)})
	ctx2.SetVariableValue(variable("value", 0, nil), data.NewIntValue(age+seconds))
	return symfonyResponseSetSharedMaxAge(ctx2)
}

func symfonyResponseSetClientTtl(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	seconds := intParam(ctx, 0, 0)
	age := computeAge(cv)
	respAddCC(responseHeaders(cv), "max-age", age+seconds)
	return responseSelf(ctx), nil
}

func symfonyResponseGetLastModified(ctx data.Context) (data.GetValue, data.Control) {
	return headerDateOrNull(responseHeaders(responseClassValue(ctx)), "Last-Modified")
}

func symfonyResponseSetLastModified(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := ctx.GetIndexValue(0)
	if v == nil || func() bool { _, ok := v.(*data.NullValue); return ok }() {
		respHeaderRemove(responseHeaders(cv), "Last-Modified")
		return responseSelf(ctx), nil
	}
	t, ok := parseDateArg(v)
	if !ok {
		t = time.Now().UTC()
	}
	respHeaderSet(responseHeaders(cv), "Last-Modified", []string{respFormatHTTPDate(t)}, true)
	return responseSelf(ctx), nil
}

func symfonyResponseGetEtag(ctx data.Context) (data.GetValue, data.Control) {
	headers := responseHeaders(responseClassValue(ctx))
	if !respHeaderHas(headers, "ETag") {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(respHeaderGet(headers, "ETag", "")), nil
}

func symfonyResponseSetEtag(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headers := responseHeaders(cv)
	etag, _ := ctx.GetIndexValue(0)
	weak := boolParam(ctx, 1, false)
	if etag == nil || func() bool { _, ok := etag.(*data.NullValue); return ok }() {
		respHeaderRemove(headers, "ETag")
		return responseSelf(ctx), nil
	}
	s := etag.AsString()
	if !strings.HasPrefix(s, "\"") {
		s = "\"" + s + "\""
	}
	if weak {
		s = "W/" + s
	}
	respHeaderSet(headers, "ETag", []string{s}, true)
	return responseSelf(ctx), nil
}

var cacheControlOptionKeys = map[string]bool{
	"must_revalidate": true, "no_cache": true, "no_store": true, "no_transform": true,
	"public": true, "private": true, "proxy_revalidate": true, "max_age": true, "s_maxage": true,
	"stale_if_error": true, "stale_while_revalidate": true, "immutable": true,
	"last_modified": true, "etag": true,
}

func symfonyResponseSetCache(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	optsVal, _ := ctx.GetIndexValue(0)
	opts, err := valueToAssocMap(optsVal)
	if err != nil {
		return nil, throwNamed("InvalidArgumentException", "%s", err.Error())
	}
	for k := range opts {
		if !cacheControlOptionKeys[k] {
			return nil, throwNamed("InvalidArgumentException", `Response does not support the following options: "%s".`, k)
		}
	}
	if v, ok := opts["etag"]; ok {
		respHeaderSet(responseHeaders(cv), "ETag", []string{v.AsString()}, true)
	}
	if v, ok := opts["last_modified"]; ok {
		if t, ok := parseDateArg(v); ok {
			respHeaderSet(responseHeaders(cv), "Last-Modified", []string{respFormatHTTPDate(t)}, true)
		}
	}
	if v, ok := opts["max_age"]; ok {
		if n, err := strconv.Atoi(v.AsString()); err == nil {
			respAddCC(responseHeaders(cv), "max-age", n)
		}
	}
	if v, ok := opts["s_maxage"]; ok {
		if n, err := strconv.Atoi(v.AsString()); err == nil {
			_, _ = symfonyResponseSetPublic(ctx)
			respAddCC(responseHeaders(cv), "s-maxage", n)
		}
	}
	if v, ok := opts["stale_while_revalidate"]; ok {
		if n, err := strconv.Atoi(v.AsString()); err == nil {
			respAddCC(responseHeaders(cv), "stale-while-revalidate", n)
		}
	}
	if v, ok := opts["stale_if_error"]; ok {
		if n, err := strconv.Atoi(v.AsString()); err == nil {
			respAddCC(responseHeaders(cv), "stale-if-error", n)
		}
	}
	boolDirs := []string{"must_revalidate", "no_cache", "no_store", "no_transform", "proxy_revalidate", "immutable"}
	for _, d := range boolDirs {
		if v, ok := opts[d]; ok {
			dir := strings.ReplaceAll(d, "_", "-")
			truthy := valueIsTruthy(v)
			if truthy {
				respAddCC(responseHeaders(cv), dir, true)
			} else {
				respRemoveCC(responseHeaders(cv), dir)
			}
		}
	}
	if v, ok := opts["public"]; ok {
		if valueIsTruthy(v) {
			_, _ = symfonyResponseSetPublic(ctx)
		} else {
			_, _ = symfonyResponseSetPrivate(ctx)
		}
	}
	if v, ok := opts["private"]; ok {
		if valueIsTruthy(v) {
			_, _ = symfonyResponseSetPrivate(ctx)
		} else {
			_, _ = symfonyResponseSetPublic(ctx)
		}
	}
	return responseSelf(ctx), nil
}

func symfonyResponseSetNotModified(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	if ctl := applyStatusCode(cv, 304, nil); ctl != nil {
		return nil, ctl
	}
	responseSetProp(cv, "content", data.NewStringValue(""))
	headers := responseHeaders(cv)
	for _, h := range []string{"Allow", "Content-Encoding", "Content-Language", "Content-Length", "Content-MD5", "Content-Type", "Last-Modified"} {
		respHeaderRemove(headers, h)
	}
	return responseSelf(ctx), nil
}

func symfonyResponseHasVary(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(respHeaderHas(responseHeaders(responseClassValue(ctx)), "Vary")), nil
}

func symfonyResponseGetVary(ctx data.Context) (data.GetValue, data.Control) {
	items := respHeaderGetAll(responseHeaders(responseClassValue(ctx)), "Vary")
	out := []data.Value{}
	for _, item := range items {
		for _, part := range strings.FieldsFunc(item, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' }) {
			if part != "" {
				out = append(out, data.NewStringValue(part))
			}
		}
	}
	return data.NewArrayValue(out), nil
}

func symfonyResponseSetVary(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headersVal, _ := ctx.GetIndexValue(0)
	replace := boolParam(ctx, 1, true)
	var values []string
	if arr, ok := headersVal.(*data.ArrayValue); ok {
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				values = append(values, z.Value.AsString())
			}
		}
	} else if headersVal != nil {
		values = []string{headersVal.AsString()}
	}
	respHeaderSet(responseHeaders(cv), "Vary", values, replace)
	return responseSelf(ctx), nil
}

func symfonyResponseIsNotModified(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	reqVal, _ := ctx.GetIndexValue(0)
	req, _ := reqVal.(*data.ClassValue)
	if req == nil {
		return data.NewBoolValue(false), nil
	}
	cacheable, _ := callObjMethod(req, "isMethodCacheable")
	if !valueIsTruthy(cacheable) {
		return data.NewBoolValue(false), nil
	}

	notModified := false
	etag := respHeaderGet(responseHeaders(cv), "ETag", "")
	etagsVal, _ := callObjMethod(req, "getETags")
	if etagsVal != nil && etag != "" {
		matchEtag := etag
		if strings.HasPrefix(matchEtag, "W/") {
			matchEtag = matchEtag[2:]
		}
		if arr, ok := etagsVal.(*data.ArrayValue); ok {
			for _, z := range arr.List {
				if z == nil || z.Value == nil {
					continue
				}
				ifNone := z.Value.AsString()
				if strings.HasPrefix(ifNone, "W/") {
					ifNone = ifNone[2:]
				}
				if ifNone == matchEtag || ifNone == "*" {
					notModified = true
					break
				}
			}
		}
	} else {
		reqHeaders, _ := req.GetProperty("headers")
		modifiedSince := ""
		if rh, ok := reqHeaders.(*data.ClassValue); ok {
			got, _ := callObjMethod(rh, "get", data.NewStringValue("If-Modified-Since"))
			if got != nil {
				if _, isNull := got.(*data.NullValue); !isNull {
					modifiedSince = got.(data.Value).AsString()
				}
			}
		}
		lastModified := respHeaderGet(responseHeaders(cv), "Last-Modified", "")
		if modifiedSince != "" && lastModified != "" {
			t1, ok1 := respParseHTTPDate(modifiedSince)
			t2, ok2 := respParseHTTPDate(lastModified)
			if ok1 && ok2 && !t1.Before(t2) {
				notModified = true
			}
		}
	}

	if notModified {
		_, _ = symfonyResponseSetNotModified(ctx)
	}
	return data.NewBoolValue(notModified), nil
}

func symfonyResponseIsInvalid(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code < 100 || code >= 600), nil
}
func symfonyResponseIsInformational(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code >= 100 && code < 200), nil
}
func symfonyResponseIsSuccessful(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code >= 200 && code < 300), nil
}
func symfonyResponseIsRedirection(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code >= 300 && code < 400), nil
}
func symfonyResponseIsClientError(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code >= 400 && code < 500), nil
}
func symfonyResponseIsServerError(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code >= 500 && code < 600), nil
}
func symfonyResponseIsOk(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(responseStatusCode(responseClassValue(ctx)) == 200), nil
}
func symfonyResponseIsForbidden(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(responseStatusCode(responseClassValue(ctx)) == 403), nil
}
func symfonyResponseIsNotFound(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(responseStatusCode(responseClassValue(ctx)) == 404), nil
}

func symfonyResponseIsRedirect(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	code := responseStatusCode(cv)
	ok := false
	switch code {
	case 201, 301, 302, 303, 307, 308:
		ok = true
	}
	if !ok {
		return data.NewBoolValue(false), nil
	}
	loc, present, ok := optionalStringParam(ctx, 0)
	if !present || !ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(loc == respHeaderGet(responseHeaders(cv), "Location", "")), nil
}

func symfonyResponseIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	code := responseStatusCode(responseClassValue(ctx))
	return data.NewBoolValue(code == 204 || code == 304), nil
}

func symfonyResponseCloseOutputBuffers(ctx data.Context) (data.GetValue, data.Control) {
	// Origami 无完整 PHP output buffering 栈，保留空实现以兼容调用。
	return data.NewNullValue(), nil
}

func symfonyResponseSetContentSafe(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headers := responseHeaders(cv)
	safe := boolParam(ctx, 0, true)
	if safe {
		respHeaderSet(headers, "Preference-Applied", []string{"safe"}, true)
	} else if respHeaderGet(headers, "Preference-Applied", "") == "safe" {
		respHeaderRemove(headers, "Preference-Applied")
	}
	respHeaderSet(headers, "Vary", []string{"Prefer"}, false)
	return data.NewNullValue(), nil
}
