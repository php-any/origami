package php

import (
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HashEqualsFunction 实现 hash_equals（时序安全比较）。
type HashEqualsFunction struct{}

func NewHashEqualsFunction() data.FuncStmt { return &HashEqualsFunction{} }

func (f *HashEqualsFunction) GetName() string { return "hash_equals" }

func (f *HashEqualsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "known_string", 0, nil, data.String{}),
		node.NewParameter(nil, "user_string", 1, nil, data.String{}),
	}
}

func (f *HashEqualsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "known_string", 0, data.String{}),
		node.NewVariable(nil, "user_string", 1, data.String{}),
	}
}

func (f *HashEqualsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	a, _ := ctx.GetIndexValue(0)
	b, _ := ctx.GetIndexValue(1)
	as, bs := "", ""
	if a != nil {
		as = a.AsString()
	}
	if b != nil {
		bs = b.AsString()
	}
	if len(as) != len(bs) {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(subtle.ConstantTimeCompare([]byte(as), []byte(bs)) == 1), nil
}

// SetCookieFunction 实现 setcookie（支持 expires 整数与 options 数组两种签名）。
type SetCookieFunction struct{}

func NewSetCookieFunction() data.FuncStmt { return &SetCookieFunction{} }

func (f *SetCookieFunction) GetName() string { return "setcookie" }

func (f *SetCookieFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
		node.NewParameter(nil, "value", 1, node.NewStringLiteral(nil, ""), data.String{}),
		node.NewParameter(nil, "expires_or_options", 2, node.NewIntLiteral(nil, "0"), data.Mixed{}),
		node.NewParameter(nil, "path", 3, node.NewStringLiteral(nil, ""), data.String{}),
		node.NewParameter(nil, "domain", 4, node.NewStringLiteral(nil, ""), data.String{}),
		node.NewParameter(nil, "secure", 5, node.NewBooleanLiteral(nil, false), data.Bool{}),
		node.NewParameter(nil, "httponly", 6, node.NewBooleanLiteral(nil, false), data.Bool{}),
	}
}

func (f *SetCookieFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
		node.NewVariable(nil, "value", 1, data.String{}),
		node.NewVariable(nil, "expires_or_options", 2, data.Mixed{}),
		node.NewVariable(nil, "path", 3, data.String{}),
		node.NewVariable(nil, "domain", 4, data.String{}),
		node.NewVariable(nil, "secure", 5, data.Bool{}),
		node.NewVariable(nil, "httponly", 6, data.Bool{}),
	}
}

func (f *SetCookieFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	nameVal, _ := ctx.GetIndexValue(0)
	valueVal, _ := ctx.GetIndexValue(1)
	name, value := "", ""
	if nameVal != nil {
		name = nameVal.AsString()
	}
	if valueVal != nil {
		value = valueVal.AsString()
	}
	if name == "" {
		return data.NewBoolValue(false), nil
	}

	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}

	optVal, _ := ctx.GetIndexValue(2)
	if optVal != nil {
		if applyCookieOptions(cookie, optVal) {
			// options 数组模式
		} else if asInt, ok := optVal.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil && n > 0 {
				cookie.Expires = time.Unix(int64(n), 0)
				cookie.MaxAge = int(time.Until(cookie.Expires).Seconds())
			} else if err == nil && n < 0 {
				cookie.MaxAge = -1
			}
			if pathVal, ok := ctx.GetIndexValue(3); ok && pathVal != nil {
				if p := pathVal.AsString(); p != "" {
					cookie.Path = p
				}
			}
			if domainVal, ok := ctx.GetIndexValue(4); ok && domainVal != nil {
				cookie.Domain = domainVal.AsString()
			}
			if secureVal, ok := ctx.GetIndexValue(5); ok && secureVal != nil {
				if b, ok := secureVal.(data.AsBool); ok {
					cookie.Secure, _ = b.AsBool()
				}
			}
			if httpOnlyVal, ok := ctx.GetIndexValue(6); ok && httpOnlyVal != nil {
				if b, ok := httpOnlyVal.(data.AsBool); ok {
					cookie.HttpOnly, _ = b.AsBool()
				}
			}
		}
	}

	writer := responseWriterFromContext(ctx)
	if writer == nil {
		return data.NewBoolValue(false), nil
	}
	http.SetCookie(writer, cookie)
	return data.NewBoolValue(true), nil
}

func applyCookieOptions(cookie *http.Cookie, optVal data.Value) bool {
	get := func(key string) (data.Value, bool) {
		switch o := optVal.(type) {
		case *data.ArrayValue:
			if z, ok := o.LookupZValByStringKey(key); ok && z != nil {
				return z.Value, true
			}
		case *data.ObjectValue:
			if v, ctl := o.GetProperty(key); ctl == nil && v != nil {
				if _, isNull := v.(*data.NullValue); !isNull {
					return v, true
				}
			}
		default:
			return nil, false
		}
		return nil, false
	}

	switch optVal.(type) {
	case *data.ArrayValue, *data.ObjectValue:
	default:
		return false
	}

	if v, ok := get("expires"); ok {
		if asInt, ok := v.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				if n > 0 {
					cookie.Expires = time.Unix(int64(n), 0)
					cookie.MaxAge = int(time.Until(cookie.Expires).Seconds())
				} else if n < 0 {
					cookie.MaxAge = -1
				}
			}
		}
	}
	if v, ok := get("path"); ok {
		if p := v.AsString(); p != "" {
			cookie.Path = p
		}
	}
	if v, ok := get("domain"); ok {
		cookie.Domain = v.AsString()
	}
	if v, ok := get("secure"); ok {
		if b, ok := v.(data.AsBool); ok {
			cookie.Secure, _ = b.AsBool()
		} else if asInt, ok := v.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				cookie.Secure = n != 0
			}
		}
	}
	if v, ok := get("httponly"); ok {
		if b, ok := v.(data.AsBool); ok {
			cookie.HttpOnly, _ = b.AsBool()
		} else if asInt, ok := v.(data.AsInt); ok {
			if n, err := asInt.AsInt(); err == nil {
				cookie.HttpOnly = n != 0
			}
		}
	}
	if v, ok := get("samesite"); ok {
		switch v.AsString() {
		case "None", "none":
			cookie.SameSite = http.SameSiteNoneMode
		case "Strict", "strict":
			cookie.SameSite = http.SameSiteStrictMode
		default:
			cookie.SameSite = http.SameSiteLaxMode
		}
	}
	return true
}
