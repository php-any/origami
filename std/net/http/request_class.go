package http

import (
	"crypto/tls"
	"errors"
	"io"
	"mime/multipart"
	httpsrc "net/http"
	"net/url"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewRequestClass() data.ClassStmt {
	return &RequestClass{
		source:             nil,
		cookiesNamed:       &RequestCookiesNamedMethod{source: nil},
		protoAtLeast:       &RequestProtoAtLeastMethod{source: nil},
		withContext:        &RequestWithContextMethod{source: nil},
		writeProxy:         &RequestWriteProxyMethod{source: nil},
		context:            &RequestContextMethod{source: nil},
		cookie:             &RequestCookieMethod{source: nil},
		formValue:          &RequestFormValueMethod{source: nil},
		clone:              &RequestCloneMethod{source: nil},
		parseForm:          &RequestParseFormMethod{source: nil},
		referer:            &RequestRefererMethod{source: nil},
		setPathValue:       &RequestSetPathValueMethod{source: nil},
		pathValue:          &RequestPathValueMethod{source: nil},
		userAgent:          &RequestUserAgentMethod{source: nil},
		parseMultipartForm: &RequestParseMultipartFormMethod{source: nil},
		write:              &RequestWriteMethod{source: nil},
		setBasicAuth:       &RequestSetBasicAuthMethod{source: nil},
		formFile:           &RequestFormFileMethod{source: nil},
		multipartReader:    &RequestMultipartReaderMethod{source: nil},
		postFormValue:      &RequestPostFormValueMethod{source: nil},
		basicAuth:          &RequestBasicAuthMethod{source: nil},
		cookies:            &RequestCookiesMethod{source: nil},
		addCookie:          &RequestAddCookieMethod{source: nil},
	}
}

func NewRequestClassFrom(source *httpsrc.Request) data.ClassStmt {
	return &RequestClass{
		source:             source,
		multipartReader:    &RequestMultipartReaderMethod{source: source},
		pathValue:          &RequestPathValueMethod{source: source},
		setPathValue:       &RequestSetPathValueMethod{source: source},
		addCookie:          &RequestAddCookieMethod{source: source},
		basicAuth:          &RequestBasicAuthMethod{source: source},
		formFile:           &RequestFormFileMethod{source: source},
		formValue:          &RequestFormValueMethod{source: source},
		parseMultipartForm: &RequestParseMultipartFormMethod{source: source},
		write:              &RequestWriteMethod{source: source},
		protoAtLeast:       &RequestProtoAtLeastMethod{source: source},
		cookies:            &RequestCookiesMethod{source: source},
		parseForm:          &RequestParseFormMethod{source: source},
		withContext:        &RequestWithContextMethod{source: source},
		clone:              &RequestCloneMethod{source: source},
		postFormValue:      &RequestPostFormValueMethod{source: source},
		context:            &RequestContextMethod{source: source},
		referer:            &RequestRefererMethod{source: source},
		setBasicAuth:       &RequestSetBasicAuthMethod{source: source},
		cookiesNamed:       &RequestCookiesNamedMethod{source: source},
		cookie:             &RequestCookieMethod{source: source},
		userAgent:          &RequestUserAgentMethod{source: source},
		writeProxy:         &RequestWriteProxyMethod{source: source},
	}
}

type RequestClass struct {
	node.Node
	source             *httpsrc.Request
	multipartReader    data.Method
	pathValue          data.Method
	basicAuth          data.Method
	formValue          data.Method
	setBasicAuth       data.Method
	userAgent          data.Method
	write              data.Method
	referer            data.Method
	clone              data.Method
	parseMultipartForm data.Method
	postFormValue      data.Method
	protoAtLeast       data.Method
	parseForm          data.Method
	writeProxy         data.Method
	cookiesNamed       data.Method
	formFile           data.Method
	setPathValue       data.Method
	addCookie          data.Method
	context            data.Method
	withContext        data.Method
	cookie             data.Method
	cookies            data.Method
}

func (s *RequestClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewRequestClassFrom(&httpsrc.Request{}), ctx.CreateBaseContext()), nil
}

func (s *RequestClass) GetName() string         { return "Net\\Http\\Request" }
func (s *RequestClass) GetExtend() *string      { return nil }
func (s *RequestClass) GetImplements() []string { return nil }
func (s *RequestClass) AsString() string        { return "Request{}" }
func (s *RequestClass) GetSource() any          { return s.source }
func (s *RequestClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "context":
		return s.context, true
	case "parseForm":
		return s.parseForm, true
	case "cookie":
		return s.cookie, true
	case "formValue":
		return s.formValue, true
	case "writeProxy":
		return s.writeProxy, true
	case "addCookie":
		return s.addCookie, true
	case "cookies":
		return s.cookies, true
	case "userAgent":
		return s.userAgent, true
	case "postFormValue":
		return s.postFormValue, true
	case "multipartReader":
		return s.multipartReader, true
	case "setPathValue":
		return s.setPathValue, true
	case "pathValue":
		return s.pathValue, true
	case "referer":
		return s.referer, true
	case "withContext":
		return s.withContext, true
	case "protoAtLeast":
		return s.protoAtLeast, true
	case "parseMultipartForm":
		return s.parseMultipartForm, true
	case "setBasicAuth":
		return s.setBasicAuth, true
	case "write":
		return s.write, true
	case "cookiesNamed":
		return s.cookiesNamed, true
	case "formFile":
		return s.formFile, true
	case "clone":
		return s.clone, true
	case "basicAuth":
		return s.basicAuth, true
	}
	return nil, false
}

func (s *RequestClass) GetMethods() []data.Method {
	return []data.Method{
		s.userAgent,
		s.cookiesNamed,
		s.pathValue,
		s.setPathValue,
		s.context,
		s.postFormValue,
		s.protoAtLeast,
		s.cookies,
		s.parseMultipartForm,
		s.addCookie,
		s.clone,
		s.formFile,
		s.basicAuth,
		s.cookie,
		s.referer,
		s.writeProxy,
		s.parseForm,
		s.withContext,
		s.write,
		s.formValue,
		s.multipartReader,
		s.setBasicAuth,
	}
}

func (s *RequestClass) GetConstruct() data.Method { return nil }

func (s *RequestClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "Method":
		return node.NewProperty(nil, "Method", "public", true, data.NewAnyValue(s.source.Method)), true
	case "URL":
		return node.NewProperty(nil, "URL", "public", true, data.NewAnyValue(s.source.URL)), true
	case "Proto":
		return node.NewProperty(nil, "Proto", "public", true, data.NewAnyValue(s.source.Proto)), true
	case "ProtoMajor":
		return node.NewProperty(nil, "ProtoMajor", "public", true, data.NewAnyValue(s.source.ProtoMajor)), true
	case "ProtoMinor":
		return node.NewProperty(nil, "ProtoMinor", "public", true, data.NewAnyValue(s.source.ProtoMinor)), true
	case "Header":
		return node.NewProperty(nil, "Header", "public", true, data.NewAnyValue(s.source.Header)), true
	case "Body":
		return node.NewProperty(nil, "Body", "public", true, data.NewAnyValue(s.source.Body)), true
	case "GetBody":
		return node.NewProperty(nil, "GetBody", "public", true, data.NewAnyValue(s.source.GetBody)), true
	case "ContentLength":
		return node.NewProperty(nil, "ContentLength", "public", true, data.NewAnyValue(s.source.ContentLength)), true
	case "TransferEncoding":
		return node.NewProperty(nil, "TransferEncoding", "public", true, data.NewAnyValue(s.source.TransferEncoding)), true
	case "Close":
		return node.NewProperty(nil, "Close", "public", true, data.NewAnyValue(s.source.Close)), true
	case "Host":
		return node.NewProperty(nil, "Host", "public", true, data.NewAnyValue(s.source.Host)), true
	case "Form":
		return node.NewProperty(nil, "Form", "public", true, data.NewAnyValue(s.source.Form)), true
	case "PostForm":
		return node.NewProperty(nil, "PostForm", "public", true, data.NewAnyValue(s.source.PostForm)), true
	case "MultipartForm":
		return node.NewProperty(nil, "MultipartForm", "public", true, data.NewAnyValue(s.source.MultipartForm)), true
	case "Trailer":
		return node.NewProperty(nil, "Trailer", "public", true, data.NewAnyValue(s.source.Trailer)), true
	case "RemoteAddr":
		return node.NewProperty(nil, "RemoteAddr", "public", true, data.NewAnyValue(s.source.RemoteAddr)), true
	case "RequestURI":
		return node.NewProperty(nil, "RequestURI", "public", true, data.NewAnyValue(s.source.RequestURI)), true
	case "TLS":
		return node.NewProperty(nil, "TLS", "public", true, data.NewAnyValue(s.source.TLS)), true
	case "Cancel":
		return node.NewProperty(nil, "Cancel", "public", true, data.NewAnyValue(s.source.Cancel)), true
	case "Response":
		return node.NewProperty(nil, "Response", "public", true, data.NewAnyValue(s.source.Response)), true
	case "Pattern":
		return node.NewProperty(nil, "Pattern", "public", true, data.NewAnyValue(s.source.Pattern)), true
	}
	return nil, false
}

func (s *RequestClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "Method", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "URL", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Proto", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "ProtoMajor", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "ProtoMinor", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Header", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Body", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "GetBody", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "ContentLength", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "TransferEncoding", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Close", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Host", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Form", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "PostForm", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "MultipartForm", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Trailer", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "RemoteAddr", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "RequestURI", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "TLS", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Cancel", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Response", "public", true, data.NewAnyValue(nil)),
		node.NewProperty(nil, "Pattern", "public", true, data.NewAnyValue(nil)),
	}
}

func (s *RequestClass) SetProperty(name string, value data.Value) data.Control {
	if s.source == nil {
		return data.NewErrorThrow(nil, errors.New("无法设置属性，source 为 nil"))
	}

	switch name {
	case "Method":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Method = val
		return nil
	case "URL":
		val, err := utils.Convert[url.URL](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		converted := new(url.URL)
		*converted = val
		s.source.URL = converted
		return nil
	case "Proto":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Proto = val
		return nil
	case "ProtoMajor":
		val, err := utils.Convert[int](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.ProtoMajor = val
		return nil
	case "ProtoMinor":
		val, err := utils.Convert[int](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.ProtoMinor = val
		return nil
	case "Header":
		val, err := utils.Convert[httpsrc.Header](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Header = val
		return nil
	case "Body":
		val, err := utils.Convert[io.ReadCloser](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Body = val
		return nil
	case "GetBody":
		val, err := utils.Convert[func() (io.ReadCloser, error)](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.GetBody = val
		return nil
	case "ContentLength":
		val, err := utils.Convert[int64](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.ContentLength = val
		return nil
	case "TransferEncoding":
		val, err := utils.Convert[[]string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.TransferEncoding = val
		return nil
	case "Close":
		val, err := utils.Convert[bool](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Close = val
		return nil
	case "Host":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Host = val
		return nil
	case "Form":
		val, err := utils.Convert[url.Values](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Form = val
		return nil
	case "PostForm":
		val, err := utils.Convert[url.Values](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.PostForm = val
		return nil
	case "MultipartForm":
		val, err := utils.Convert[multipart.Form](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		converted := new(multipart.Form)
		*converted = val
		s.source.MultipartForm = converted
		return nil
	case "Trailer":
		val, err := utils.Convert[httpsrc.Header](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Trailer = val
		return nil
	case "RemoteAddr":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.RemoteAddr = val
		return nil
	case "RequestURI":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.RequestURI = val
		return nil
	case "TLS":
		val, err := utils.Convert[tls.ConnectionState](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		converted := new(tls.ConnectionState)
		*converted = val
		s.source.TLS = converted
		return nil
	case "Cancel":
		val, err := utils.Convert[chan struct{}](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Cancel = val
		return nil
	case "Response":
		val, err := utils.Convert[httpsrc.Response](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		converted := new(httpsrc.Response)
		*converted = val
		s.source.Response = converted
		return nil
	case "Pattern":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Pattern = val
		return nil
	default:
		return data.NewErrorThrow(nil, errors.New("属性不存在: "+name))
	}
}
