package routing

import (
	"github.com/php-any/origami/data"
)

const compiledRouteName = "Symfony\\Component\\Routing\\CompiledRoute"

func newCompiledRouteClass() *rtClass {
	c := newRtClass(compiledRouteName, nil, nil, []data.Property{
		privProp("staticPrefix", data.NewStringValue("")),
		privProp("regex", data.NewStringValue("")),
		privProp("tokens", emptyArray()),
		privProp("pathVariables", emptyArray()),
		privProp("hostRegex", data.NewNullValue()),
		privProp("hostTokens", emptyArray()),
		privProp("hostVariables", emptyArray()),
		privProp("variables", emptyArray()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{
			param("staticPrefix", 0, nil, nil),
			param("regex", 1, nil, nil),
			param("tokens", 2, emptyArray(), nil),
			param("pathVariables", 3, emptyArray(), nil),
			param("hostRegex", 4, data.NewNullValue(), nil),
			param("hostTokens", 5, emptyArray(), nil),
			param("hostVariables", 6, emptyArray(), nil),
			param("variables", 7, emptyArray(), nil),
		},
		[]data.Variable{
			variable("staticPrefix", 0, nil),
			variable("regex", 1, nil),
			variable("tokens", 2, nil),
			variable("pathVariables", 3, nil),
			variable("hostRegex", 4, nil),
			variable("hostTokens", 5, nil),
			variable("hostVariables", 6, nil),
			variable("variables", 7, nil),
		},
		compiledRouteConstruct))
	c.add(meth("getStaticPrefix", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "staticPrefix", "")), nil
	}))
	c.add(meth("getRegex", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "regex", "")), nil
	}))
	c.add(meth("getHostRegex", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		v := prop(rtSelf(ctx), "hostRegex")
		if isNull(v) || v.AsString() == "" {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(v.AsString()), nil
	}))
	c.add(meth("getTokens", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "tokens"), nil
	}))
	c.add(meth("getHostTokens", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "hostTokens"), nil
	}))
	c.add(meth("getVariables", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "variables"), nil
	}))
	c.add(meth("getPathVariables", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "pathVariables"), nil
	}))
	c.add(meth("getHostVariables", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "hostVariables"), nil
	}))
	c.add(meth("__serialize", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		return phpAssoc(
			"vars", prop(cv, "variables"),
			"path_prefix", prop(cv, "staticPrefix"),
			"path_regex", prop(cv, "regex"),
			"path_tokens", prop(cv, "tokens"),
			"path_vars", prop(cv, "pathVariables"),
			"host_regex", prop(cv, "hostRegex"),
			"host_tokens", prop(cv, "hostTokens"),
			"host_vars", prop(cv, "hostVariables"),
		), nil
	}))
	c.add(meth("unserialize", []string{"data"}, compiledRouteUnserialize))
	c.add(meth("__unserialize", []string{"data"}, compiledRouteUnserialize))
	return c
}

func compiledRouteConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	setProp(cv, "staticPrefix", data.NewStringValue(argString(ctx, 0, "")))
	setProp(cv, "regex", data.NewStringValue(argString(ctx, 1, "")))
	setProp(cv, "tokens", argArray(ctx, 2))
	setProp(cv, "pathVariables", argArray(ctx, 3))
	hostRegex := arg(ctx, 4)
	if hostRegex == nil || isNull(hostRegex) || hostRegex.AsString() == "" {
		setProp(cv, "hostRegex", data.NewNullValue())
	} else {
		setProp(cv, "hostRegex", data.NewStringValue(hostRegex.AsString()))
	}
	setProp(cv, "hostTokens", argArray(ctx, 5))
	setProp(cv, "hostVariables", argArray(ctx, 6))
	setProp(cv, "variables", argArray(ctx, 7))
	return data.NewNullValue(), nil
}

func compiledRouteUnserialize(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	dataArr := argArray(ctx, 0)
	if v, ok := assocGet(dataArr, "vars"); ok {
		setProp(cv, "variables", valueToArray(v))
	}
	if v, ok := assocGet(dataArr, "path_prefix"); ok {
		setProp(cv, "staticPrefix", v)
	}
	if v, ok := assocGet(dataArr, "path_regex"); ok {
		setProp(cv, "regex", v)
	}
	if v, ok := assocGet(dataArr, "path_tokens"); ok {
		setProp(cv, "tokens", valueToArray(v))
	}
	if v, ok := assocGet(dataArr, "path_vars"); ok {
		setProp(cv, "pathVariables", valueToArray(v))
	}
	if v, ok := assocGet(dataArr, "host_regex"); ok {
		setProp(cv, "hostRegex", v)
	}
	if v, ok := assocGet(dataArr, "host_tokens"); ok {
		setProp(cv, "hostTokens", valueToArray(v))
	}
	if v, ok := assocGet(dataArr, "host_vars"); ok {
		setProp(cv, "hostVariables", valueToArray(v))
	}
	return data.NewNullValue(), nil
}

func newCompiledRoute(ctx data.Context, staticPrefix, regex string, tokens, pathVars data.Value, hostRegex string, hostTokens, hostVars, variables data.Value) (*data.ClassValue, data.Control) {
	var host data.Value = data.NewNullValue()
	if hostRegex != "" {
		host = data.NewStringValue(hostRegex)
	}
	if tokens == nil {
		tokens = emptyArray()
	}
	if pathVars == nil {
		pathVars = emptyArray()
	}
	if hostTokens == nil {
		hostTokens = emptyArray()
	}
	if hostVars == nil {
		hostVars = emptyArray()
	}
	if variables == nil {
		variables = emptyArray()
	}
	return instantiate(ctx, compiledRouteName,
		data.NewStringValue(staticPrefix),
		data.NewStringValue(regex),
		tokens,
		pathVars,
		host,
		hostTokens,
		hostVars,
		variables,
	)
}
