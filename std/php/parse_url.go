package php

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewParseUrlFunction() data.FuncStmt {
	return &ParseUrlFunction{}
}

type ParseUrlFunction struct{}

func (fn *ParseUrlFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	component, _ := ctx.GetIndexValue(1)
	raw, ctl := node.ValueToDisplayString(ctx, value)
	if ctl != nil {
		return nil, ctl
	}
	// PHP replaces ASCII control bytes in URL components with underscores.
	raw = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return '_'
		}
		return r
	}, raw)
	parsed, err := url.Parse(raw)
	if err != nil || parsed == nil {
		return data.NewBoolValue(false), nil
	}
	result := data.NewArrayValueFromSlots(nil)
	set := func(key, value string) { result.SetStringKey(key, data.NewStringValue(value)) }
	if parsed.Scheme != "" {
		set("scheme", parsed.Scheme)
	}
	if parsed.Host != "" {
		host := parsed.Hostname()
		if strings.HasPrefix(parsed.Host, "[") {
			host = "[" + host + "]"
		}
		set("host", host)
	}
	if port := parsed.Port(); port != "" {
		number, err := strconv.Atoi(port)
		if err != nil || number > 65535 {
			return data.NewBoolValue(false), nil
		}
		result.SetStringKey("port", data.NewIntValue(number))
	}
	if parsed.User != nil {
		set("user", parsed.User.Username())
		if password, exists := parsed.User.Password(); exists {
			set("pass", password)
		}
	}
	path := parsed.EscapedPath()
	if parsed.Opaque != "" {
		path = parsed.Opaque
	}
	if path != "" || raw == "" {
		set("path", path)
	}
	if parsed.RawQuery != "" || parsed.ForceQuery {
		set("query", parsed.RawQuery)
	}
	if strings.Contains(raw, "#") {
		set("fragment", parsed.Fragment)
	}
	index := -1
	if integer, ok := component.(data.AsInt); ok {
		index, _ = integer.AsInt()
	}
	if index < 0 {
		return result, nil
	}
	names := [...]string{"scheme", "host", "port", "user", "pass", "path", "query", "fragment"}
	if index >= len(names) {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("parse_url(): Argument #2 ($component) must be a valid component"), "ValueError")
	}
	if slot, _ := result.LookupZValByStringKey(names[index]); slot != nil {
		return slot.ReadValue(), nil
	}
	return data.NewNullValue(), nil
}

func (fn *ParseUrlFunction) GetName() string {
	return "parse_url"
}

var parseUrlFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "url", 0, nil, nil),
	node.NewParameter(nil, "component", 1, data.NewIntValue(-1), nil),
}

func (fn *ParseUrlFunction) GetParams() []data.GetValue {
	return parseUrlFunctionGetParams
}

var parseUrlFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "url", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "component", 1, data.NewBaseType("int")),
}

func (fn *ParseUrlFunction) GetVariables() []data.Variable {
	return parseUrlFunctionGetVariables
}
