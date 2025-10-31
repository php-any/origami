package http

import (
	"errors"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// 辅助函数：创建属性值，减少重复代码
func createPropertyValue(source *httpsrc.Request, name string, value interface{}) data.Property {
	return node.NewProperty(nil, name, "public", true, data.NewAnyValue(value))
}

// 辅助函数：统一的错误处理
func handlePropertyError(err error, propertyName string) data.Control {
	return utils.NewThrowf("设置属性 %s 失败: %v", propertyName, err)
}

// 辅助函数：安全的类型转换
func safeConvert[T any](value data.Value, target *T) error {
	converted, err := utils.Convert[T](value)
	if err != nil {
		return err
	}
	*target = converted
	return nil
}

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
	// 原有的方法
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
	case "method":
		return &RequestMethodMethod{source: s.source}, true
	case "url":
		return &RequestUrlMethod{source: s.source}, true
	case "fullUrl":
		return &RequestFullUrlMethod{source: s.source}, true
	case "path":
		return &RequestPathMethod{source: s.source}, true
	case "query":
		return &RequestQueryMethod{source: s.source}, true
	case "header":
		return &RequestHeaderMethod{source: s.source}, true
	case "ip":
		return &RequestIpMethod{source: s.source}, true
	case "has":
		return &RequestHasMethod{source: s.source}, true
	case "input":
		return &RequestInputMethod{source: s.source}, true
	case "only":
		return &RequestOnlyMethod{source: s.source}, true
	case "except":
		return &RequestExceptMethod{source: s.source}, true
	case "all":
		return &RequestAllMethod{source: s.source}, true
	case "file":
		return &RequestFileMethod{source: s.source}, true
	case "isMethod":
		return &RequestIsMethodMethod{source: s.source}, true
	case "isSecure":
		return &RequestIsSecureMethod{source: s.source}, true
	case "bind":
		return &RequestBindMethod{source: s.source}, true
	case "body":
		return &RequestBodyMethod{source: s.source}, true
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
	// 所有数据访问都通过方法进行
	return nil, false
}

func (s *RequestClass) GetPropertyList() []data.Property {
	// 所有数据访问都通过方法进行
	return []data.Property{}
}

func (s *RequestClass) SetProperty(name string, value data.Value) data.Control {
	return data.NewErrorThrow(nil, errors.New("request 对象是只读的，不允许设置属性"))
}
