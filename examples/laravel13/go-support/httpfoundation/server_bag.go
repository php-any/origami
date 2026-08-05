package httpfoundation

import (
	"encoding/base64"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ServerBagClass 实现 Symfony\Component\HttpFoundation\ServerBag。
type ServerBagClass struct {
	node.Node
	source     *ParamBagData
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewServerBagClass() data.ClassStmt {
	return NewServerBagClassFrom(nil)
}

func NewServerBagClassFrom(source *ParamBagData) data.ClassStmt {
	c := &ServerBagClass{
		source:     source,
		properties: []data.Property{protectedArrayProp("parameters")},
	}
	c.methods, c.methodList = serverBagMethods()
	return c
}

func (c *ServerBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := c.source
	if src == nil {
		src = newParamBagData()
	} else {
		src = src.clone()
	}
	return data.NewProxyValue(NewServerBagClassFrom(src), ctx.CreateBaseContext()), nil
}

func (c *ServerBagClass) GetName() string {
	return fqnServerBag
}
func (c *ServerBagClass) GetExtend() *string {
	parent := fqnParameterBag
	return &parent
}
func (c *ServerBagClass) GetImplements() []string          { return nil }
func (c *ServerBagClass) GetSource() any                   { return c.source }
func (c *ServerBagClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *ServerBagClass) GetPropertyList() []data.Property { return c.properties }
func (c *ServerBagClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *ServerBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *ServerBagClass) GetMethods() []data.Method { return c.methodList }

func serverBagMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("parameters", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("parameters", 0, nil)},
			nil, parameterBagConstruct),
		pubMethod("getHeaders", nil, nil, data.NewBaseType("array"), serverBagGetHeaders),
	}
	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func serverBagGetHeaders(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	headers := map[string]data.Value{}
	if store == nil {
		return assocMapToArrayValue(headers), nil
	}

	for _, pair := range store.orderedPairs() {
		key, value := pair.key, pair.val
		if strings.HasPrefix(key, "HTTP_") {
			headers[key[5:]] = value
		} else if (key == "CONTENT_TYPE" || key == "CONTENT_LENGTH" || key == "CONTENT_MD5") && value.AsString() != "" {
			headers[key] = value
		}
	}

	if user, ok := store.get("PHP_AUTH_USER"); ok {
		headers["PHP_AUTH_USER"] = user
		if pw, ok := store.get("PHP_AUTH_PW"); ok {
			headers["PHP_AUTH_PW"] = pw
		} else {
			headers["PHP_AUTH_PW"] = data.NewStringValue("")
		}
	} else {
		var authorizationHeader string
		if v, ok := store.get("HTTP_AUTHORIZATION"); ok {
			authorizationHeader = v.AsString()
		} else if v, ok := store.get("REDIRECT_HTTP_AUTHORIZATION"); ok {
			authorizationHeader = v.AsString()
		}
		if authorizationHeader != "" {
			lower := strings.ToLower(authorizationHeader)
			if strings.HasPrefix(lower, "basic ") {
				decoded, err := base64.StdEncoding.DecodeString(authorizationHeader[6:])
				if err == nil {
					parts := strings.SplitN(string(decoded), ":", 2)
					if len(parts) == 2 {
						headers["PHP_AUTH_USER"] = data.NewStringValue(parts[0])
						headers["PHP_AUTH_PW"] = data.NewStringValue(parts[1])
					}
				}
			} else if _, hasDigest := store.get("PHP_AUTH_DIGEST"); !hasDigest && strings.HasPrefix(lower, "digest ") {
				headers["PHP_AUTH_DIGEST"] = data.NewStringValue(authorizationHeader)
				store.set("PHP_AUTH_DIGEST", data.NewStringValue(authorizationHeader))
			} else if strings.HasPrefix(lower, "bearer ") {
				headers["AUTHORIZATION"] = data.NewStringValue(authorizationHeader)
			}
		}
	}

	if _, ok := headers["AUTHORIZATION"]; ok {
		return assocMapToArrayValue(headers), nil
	}
	if user, ok := headers["PHP_AUTH_USER"]; ok {
		pw := ""
		if p, ok := headers["PHP_AUTH_PW"]; ok {
			pw = p.AsString()
		}
		headers["AUTHORIZATION"] = data.NewStringValue("Basic " + base64.StdEncoding.EncodeToString([]byte(user.AsString()+":"+pw)))
	} else if digest, ok := headers["PHP_AUTH_DIGEST"]; ok {
		headers["AUTHORIZATION"] = digest
	}
	return assocMapToArrayValue(headers), nil
}
