package php

import (
	"net/url"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewParseUrlFunction() data.FuncStmt {
	return &ParseUrlFunction{}
}

type ParseUrlFunction struct{}

func (fn *ParseUrlFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	urlStr, _ := ctx.GetIndexValue(0)
	component, _ := ctx.GetIndexValue(1)

	u := ""
	if s, ok := urlStr.(data.AsString); ok {
		u = s.AsString()
	}
	if u == "" {
		return data.NewBoolValue(false), nil
	}
	// Handle null values
	if _, isNull := urlStr.(*data.NullValue); isNull {
		return data.NewBoolValue(false), nil
	}

	parsed, err := url.Parse(u)
	if err != nil || parsed == nil {
		return data.NewBoolValue(false), nil
	}

	comp := -1 // PHP_URL_ALL
	if c, ok := component.(data.AsInt); ok {
		if n, err := c.AsInt(); err == nil {
			comp = n
		}
	}

	if comp >= 0 {
		// Return specific component as string or null
		switch comp {
		case 0: // PHP_URL_SCHEME
			if parsed.Scheme != "" {
				return data.NewStringValue(parsed.Scheme), nil
			}
		case 1: // PHP_URL_HOST
			if parsed.Host != "" {
				return data.NewStringValue(parsed.Host), nil
			}
		case 2: // PHP_URL_PORT
			if parsed.Port() != "" {
				return data.NewStringValue(parsed.Port()), nil
			}
		case 3: // PHP_URL_USER
			if parsed.User != nil {
				return data.NewStringValue(parsed.User.Username()), nil
			}
		case 4: // PHP_URL_PASS
			if parsed.User != nil {
				if pwd, ok2 := parsed.User.Password(); ok2 {
					return data.NewStringValue(pwd), nil
				}
			}
		case 5: // PHP_URL_PATH
			if parsed.Path != "" {
				return data.NewStringValue(parsed.Path), nil
			}
		case 6: // PHP_URL_QUERY
			if parsed.RawQuery != "" {
				return data.NewStringValue(parsed.RawQuery), nil
			}
		case 7: // PHP_URL_FRAGMENT
			if parsed.Fragment != "" {
				return data.NewStringValue(parsed.Fragment), nil
			}
		}
		return data.NewNullValue(), nil
	}

	// Return full array
	result := data.NewObjectValue()
	if parsed.Scheme != "" {
		result.SetProperty("scheme", data.NewStringValue(parsed.Scheme))
	}
	if parsed.Host != "" {
		result.SetProperty("host", data.NewStringValue(parsed.Host))
	}
	if port := parsed.Port(); port != "" {
		result.SetProperty("port", data.NewStringValue(port))
	}
	if parsed.User != nil {
		result.SetProperty("user", data.NewStringValue(parsed.User.Username()))
		if pwd, ok := parsed.User.Password(); ok {
			result.SetProperty("pass", data.NewStringValue(pwd))
		}
	}
	if parsed.Path != "" {
		result.SetProperty("path", data.NewStringValue(parsed.Path))
	}
	if parsed.RawQuery != "" {
		result.SetProperty("query", data.NewStringValue(parsed.RawQuery))
	}
	if parsed.Fragment != "" {
		result.SetProperty("fragment", data.NewStringValue(parsed.Fragment))
	}
	return result, nil
}

func (fn *ParseUrlFunction) GetName() string {
	return "parse_url"
}

func (fn *ParseUrlFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "url", 0, nil, nil),
		node.NewParameter(nil, "component", 1, data.NewIntValue(-1), nil),
	}
}

func (fn *ParseUrlFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "url", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "component", 1, data.NewBaseType("int")),
	}
}
