package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"

	httpsrc "net/http"
	"time"
)

func NewCookieClass() data.ClassStmt {
	return &CookieClass{
		source: nil,
		valid:  &CookieValidMethod{source: nil},
		string: &CookieStringMethod{source: nil},
	}
}

func NewCookieClassFrom(source *httpsrc.Cookie) data.ClassStmt {
	return &CookieClass{
		source: source,
		string: &CookieStringMethod{source: source},
		valid:  &CookieValidMethod{source: source},
	}
}

type CookieClass struct {
	node.Node
	source *httpsrc.Cookie
	string data.Method
	valid  data.Method
}

func (s *CookieClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewCookieClassFrom(&httpsrc.Cookie{}), ctx.CreateBaseContext()), nil
}

func (s *CookieClass) GetName() string         { return "http\\Cookie" }
func (s *CookieClass) GetExtend() *string      { return nil }
func (s *CookieClass) GetImplements() []string { return nil }
func (s *CookieClass) AsString() string        { return "Cookie{}" }
func (s *CookieClass) GetSource() any          { return s.source }
func (s *CookieClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "string":
		return s.string, true
	case "valid":
		return s.valid, true
	}
	return nil, false
}

func (s *CookieClass) GetMethods() []data.Method {
	return []data.Method{
		s.valid,
		s.string,
	}
}

func (s *CookieClass) GetConstruct() data.Method { return nil }

func (s *CookieClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "Name":
		return node.NewProperty(nil, "Name", "public", true, data.NewAnyValue(s.source.Name)), true
	case "Value":
		return node.NewProperty(nil, "Value", "public", true, data.NewAnyValue(s.source.Value)), true
	case "Quoted":
		return node.NewProperty(nil, "Quoted", "public", true, data.NewAnyValue(s.source.Quoted)), true
	case "Path":
		return node.NewProperty(nil, "Path", "public", true, data.NewAnyValue(s.source.Path)), true
	case "Domain":
		return node.NewProperty(nil, "Domain", "public", true, data.NewAnyValue(s.source.Domain)), true
	case "Expires":
		return node.NewProperty(nil, "Expires", "public", true, data.NewAnyValue(s.source.Expires)), true
	case "RawExpires":
		return node.NewProperty(nil, "RawExpires", "public", true, data.NewAnyValue(s.source.RawExpires)), true
	case "MaxAge":
		return node.NewProperty(nil, "MaxAge", "public", true, data.NewAnyValue(s.source.MaxAge)), true
	case "Secure":
		return node.NewProperty(nil, "Secure", "public", true, data.NewAnyValue(s.source.Secure)), true
	case "HttpOnly":
		return node.NewProperty(nil, "HttpOnly", "public", true, data.NewAnyValue(s.source.HttpOnly)), true
	case "SameSite":
		return node.NewProperty(nil, "SameSite", "public", true, data.NewAnyValue(s.source.SameSite)), true
	case "Partitioned":
		return node.NewProperty(nil, "Partitioned", "public", true, data.NewAnyValue(s.source.Partitioned)), true
	case "Raw":
		return node.NewProperty(nil, "Raw", "public", true, data.NewAnyValue(s.source.Raw)), true
	case "Unparsed":
		return node.NewProperty(nil, "Unparsed", "public", true, data.NewAnyValue(s.source.Unparsed)), true
	}
	return nil, false
}

func (s *CookieClass) GetProperties() map[string]data.Property {
	return map[string]data.Property{
		"Name":        node.NewProperty(nil, "Name", "public", true, data.NewAnyValue(nil)),
		"Value":       node.NewProperty(nil, "Value", "public", true, data.NewAnyValue(nil)),
		"Quoted":      node.NewProperty(nil, "Quoted", "public", true, data.NewAnyValue(nil)),
		"Path":        node.NewProperty(nil, "Path", "public", true, data.NewAnyValue(nil)),
		"Domain":      node.NewProperty(nil, "Domain", "public", true, data.NewAnyValue(nil)),
		"Expires":     node.NewProperty(nil, "Expires", "public", true, data.NewAnyValue(nil)),
		"RawExpires":  node.NewProperty(nil, "RawExpires", "public", true, data.NewAnyValue(nil)),
		"MaxAge":      node.NewProperty(nil, "MaxAge", "public", true, data.NewAnyValue(nil)),
		"Secure":      node.NewProperty(nil, "Secure", "public", true, data.NewAnyValue(nil)),
		"HttpOnly":    node.NewProperty(nil, "HttpOnly", "public", true, data.NewAnyValue(nil)),
		"SameSite":    node.NewProperty(nil, "SameSite", "public", true, data.NewAnyValue(nil)),
		"Partitioned": node.NewProperty(nil, "Partitioned", "public", true, data.NewAnyValue(nil)),
		"Raw":         node.NewProperty(nil, "Raw", "public", true, data.NewAnyValue(nil)),
		"Unparsed":    node.NewProperty(nil, "Unparsed", "public", true, data.NewAnyValue(nil)),
	}
}

func (s *CookieClass) SetProperty(name string, value data.Value) data.Control {
	if s.source == nil {
		return data.NewErrorThrow(nil, errors.New("无法设置属性，source 为 nil"))
	}

	switch name {
	case "Name":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Name = val
		return nil
	case "Value":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Value = val
		return nil
	case "Quoted":
		val, err := utils.Convert[bool](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Quoted = val
		return nil
	case "Path":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Path = val
		return nil
	case "Domain":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Domain = val
		return nil
	case "Expires":
		val, err := utils.Convert[time.Time](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Expires = val
		return nil
	case "RawExpires":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.RawExpires = val
		return nil
	case "MaxAge":
		val, err := utils.Convert[int](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.MaxAge = val
		return nil
	case "Secure":
		val, err := utils.Convert[bool](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Secure = val
		return nil
	case "HttpOnly":
		val, err := utils.Convert[bool](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.HttpOnly = val
		return nil
	case "SameSite":
		val, err := utils.Convert[httpsrc.SameSite](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.SameSite = val
		return nil
	case "Partitioned":
		val, err := utils.Convert[bool](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Partitioned = val
		return nil
	case "Raw":
		val, err := utils.Convert[string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Raw = val
		return nil
	case "Unparsed":
		val, err := utils.Convert[[]string](value)
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		s.source.Unparsed = val
		return nil
	default:
		return data.NewErrorThrow(nil, errors.New("属性不存在: "+name))
	}
}
