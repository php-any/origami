package routing

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/php-any/origami/data"
)

const routeName = "Symfony\\Component\\Routing\\Route"

func newRouteClass() *rtClass {
	c := newRtClass(routeName, nil, nil, []data.Property{
		privProp("path", data.NewStringValue("/")),
		privProp("host", data.NewStringValue("")),
		privProp("schemes", emptyArray()),
		privProp("methods", emptyArray()),
		privProp("defaults", emptyArray()),
		privProp("requirements", emptyArray()),
		privProp("options", emptyArray()),
		privProp("condition", data.NewStringValue("")),
		privProp("compiled", data.NewNullValue()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{
			param("path", 0, nil, nil),
			param("defaults", 1, emptyArray(), nil),
			param("requirements", 2, emptyArray(), nil),
			param("options", 3, emptyArray(), nil),
			param("host", 4, data.NewStringValue(""), nil),
			param("schemes", 5, emptyArray(), nil),
			param("methods", 6, emptyArray(), nil),
			param("condition", 7, data.NewStringValue(""), nil),
		},
		[]data.Variable{
			variable("path", 0, nil),
			variable("defaults", 1, nil),
			variable("requirements", 2, nil),
			variable("options", 3, nil),
			variable("host", 4, nil),
			variable("schemes", 5, nil),
			variable("methods", 6, nil),
			variable("condition", 7, nil),
		},
		routeConstruct))
	c.add(meth("getPath", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "path", "/")), nil
	}))
	c.add(meth("setPath", []string{"pattern"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeSetPath(cv, argString(ctx, 0, "/"))
		return cv, nil
	}))
	c.add(meth("getHost", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "host", "")), nil
	}))
	c.add(meth("setHost", []string{"pattern"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		pattern := argString(ctx, 0, "")
		routeSetHost(cv, pattern)
		return cv, nil
	}))
	c.add(meth("getSchemes", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "schemes"), nil
	}))
	c.add(meth("setSchemes", []string{"schemes"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeSetSchemes(cv, arg(ctx, 0))
		return cv, nil
	}))
	c.add(meth("hasScheme", []string{"scheme"}, func(ctx data.Context) (data.GetValue, data.Control) {
		ss := toStringSlice(prop(rtSelf(ctx), "schemes"))
		return data.NewBoolValue(inStringSlice(ss, strings.ToLower(argString(ctx, 0, "")))), nil
	}))
	c.add(meth("getMethods", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "methods"), nil
	}))
	c.add(meth("setMethods", []string{"methods"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeSetMethods(cv, arg(ctx, 0))
		return cv, nil
	}))
	c.add(meth("getOptions", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "options"), nil
	}))
	c.add(meth("setOptions", []string{"options"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeSetOptions(cv, argArray(ctx, 0))
		return cv, nil
	}))
	c.add(meth("addOptions", []string{"options"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeAddOptions(cv, argArray(ctx, 0))
		return cv, nil
	}))
	c.add(meth("setOption", []string{"name", "value"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		opts := cloneArray(propArray(cv, "options"))
		assocSet(opts, argString(ctx, 0, ""), arg(ctx, 1))
		setProp(cv, "options", opts)
		setProp(cv, "compiled", data.NewNullValue())
		return cv, nil
	}))
	c.add(meth("getOption", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		v, ok := assocGet(propArray(rtSelf(ctx), "options"), argString(ctx, 0, ""))
		if !ok {
			return data.NewNullValue(), nil
		}
		return v, nil
	}))
	c.add(meth("hasOption", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewBoolValue(assocHas(propArray(rtSelf(ctx), "options"), argString(ctx, 0, ""))), nil
	}))
	c.add(meth("getDefaults", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "defaults"), nil
	}))
	c.add(meth("setDefaults", []string{"defaults"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "defaults", phpList())
		routeAddDefaults(cv, argArray(ctx, 0))
		return cv, nil
	}))
	c.add(meth("addDefaults", []string{"defaults"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeAddDefaults(cv, argArray(ctx, 0))
		return cv, nil
	}))
	c.add(meth("getDefault", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return routeGetDefault(rtSelf(ctx), argString(ctx, 0, "")), nil
	}))
	c.add(meth("hasDefault", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewBoolValue(routeHasDefault(rtSelf(ctx), argString(ctx, 0, ""))), nil
	}))
	c.add(meth("setDefault", []string{"name", "default"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		routeSetDefault(cv, argString(ctx, 0, ""), arg(ctx, 1))
		return cv, nil
	}))
	c.add(meth("getRequirements", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "requirements"), nil
	}))
	c.add(meth("setRequirements", []string{"requirements"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "requirements", phpList())
		if ctl := routeAddRequirements(cv, argArray(ctx, 0)); ctl != nil {
			return nil, ctl
		}
		return cv, nil
	}))
	c.add(meth("addRequirements", []string{"requirements"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		if ctl := routeAddRequirements(cv, argArray(ctx, 0)); ctl != nil {
			return nil, ctl
		}
		return cv, nil
	}))
	c.add(meth("getRequirement", []string{"key"}, func(ctx data.Context) (data.GetValue, data.Control) {
		s := routeGetRequirement(rtSelf(ctx), argString(ctx, 0, ""))
		if s == "" && !assocHas(propArray(rtSelf(ctx), "requirements"), argString(ctx, 0, "")) {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(s), nil
	}))
	c.add(meth("hasRequirement", []string{"key"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewBoolValue(assocHas(propArray(rtSelf(ctx), "requirements"), argString(ctx, 0, ""))), nil
	}))
	c.add(meth("setRequirement", []string{"key", "regex"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		if ctl := routeSetRequirement(cv, argString(ctx, 0, ""), argString(ctx, 1, "")); ctl != nil {
			return nil, ctl
		}
		return cv, nil
	}))
	c.add(meth("getCondition", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "condition", "")), nil
	}))
	c.add(meth("setCondition", []string{"condition"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "condition", data.NewStringValue(argString(ctx, 0, "")))
		setProp(cv, "compiled", data.NewNullValue())
		return cv, nil
	}))
	c.add(meth("compile", nil, routeCompile))
	c.add(meth("__serialize", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		return phpAssoc(
			"path", prop(cv, "path"),
			"host", prop(cv, "host"),
			"defaults", prop(cv, "defaults"),
			"requirements", prop(cv, "requirements"),
			"options", prop(cv, "options"),
			"schemes", prop(cv, "schemes"),
			"methods", prop(cv, "methods"),
			"condition", prop(cv, "condition"),
			"compiled", prop(cv, "compiled"),
		), nil
	}))
	c.add(meth("__unserialize", []string{"data"}, routeUnserialize))
	return c
}

func routeConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	routeSetOptions(cv, argArray(ctx, 3))
	routeSetPath(cv, argString(ctx, 0, "/"))
	routeAddDefaults(cv, argArray(ctx, 1))
	if ctl := routeAddRequirements(cv, argArray(ctx, 2)); ctl != nil {
		return nil, ctl
	}
	host := arg(ctx, 4)
	if host == nil || isNull(host) {
		routeSetHost(cv, "")
	} else {
		routeSetHost(cv, host.AsString())
	}
	routeSetSchemes(cv, arg(ctx, 5))
	routeSetMethods(cv, arg(ctx, 6))
	cond := arg(ctx, 7)
	if cond == nil || isNull(cond) {
		setProp(cv, "condition", data.NewStringValue(""))
	} else {
		setProp(cv, "condition", data.NewStringValue(cond.AsString()))
	}
	return data.NewNullValue(), nil
}

func routeSetPath(cv *data.ClassValue, pattern string) {
	pattern = extractInlineDefaultsAndRequirements(cv, pattern)
	pattern = "/" + strings.TrimLeft(strings.TrimSpace(pattern), "/")
	setProp(cv, "path", data.NewStringValue(pattern))
	setProp(cv, "compiled", data.NewNullValue())
}

func routeSetHost(cv *data.ClassValue, pattern string) {
	pattern = extractInlineDefaultsAndRequirements(cv, pattern)
	setProp(cv, "host", data.NewStringValue(pattern))
	setProp(cv, "compiled", data.NewNullValue())
}

func routeSetSchemes(cv *data.ClassValue, v data.Value) {
	ss := toStringSlice(v)
	for i := range ss {
		ss[i] = strings.ToLower(ss[i])
	}
	setProp(cv, "schemes", stringsToArray(ss))
	setProp(cv, "compiled", data.NewNullValue())
}

func routeSetMethods(cv *data.ClassValue, v data.Value) {
	ss := toStringSlice(v)
	for i := range ss {
		ss[i] = strings.ToUpper(ss[i])
	}
	setProp(cv, "methods", stringsToArray(ss))
	setProp(cv, "compiled", data.NewNullValue())
}

func routeSetOptions(cv *data.ClassValue, extra *data.ArrayValue) {
	opts := phpAssoc("compiler_class", data.NewStringValue(routeCompilerName))
	if extra != nil {
		for i, z := range extra.List {
			if z == nil || z.Value == nil {
				continue
			}
			assocSet(opts, arrayKey(z, i), z.Value)
		}
	}
	setProp(cv, "options", opts)
	setProp(cv, "compiled", data.NewNullValue())
}

func routeAddOptions(cv *data.ClassValue, extra *data.ArrayValue) {
	opts := cloneArray(propArray(cv, "options"))
	if extra != nil {
		for i, z := range extra.List {
			if z == nil || z.Value == nil {
				continue
			}
			assocSet(opts, arrayKey(z, i), z.Value)
		}
	}
	setProp(cv, "options", opts)
	setProp(cv, "compiled", data.NewNullValue())
}

func routeIsLocalized(cv *data.ClassValue) bool {
	locale := routeGetDefault(cv, "_locale")
	canonical := routeGetDefault(cv, "_canonical_route")
	if isNull(locale) || isNull(canonical) {
		return false
	}
	req := routeGetRequirement(cv, "_locale")
	return req == regexp.QuoteMeta(locale.AsString())
}

func routeAddDefaults(cv *data.ClassValue, extra *data.ArrayValue) {
	if extra == nil {
		return
	}
	defs := cloneArray(propArray(cv, "defaults"))
	localized := routeIsLocalized(cv)
	for i, z := range extra.List {
		if z == nil {
			continue
		}
		key := arrayKey(z, i)
		if key == "_locale" && localized {
			continue
		}
		val := z.Value
		if val == nil {
			val = data.NewNullValue()
		}
		assocSet(defs, key, val)
	}
	setProp(cv, "defaults", defs)
	setProp(cv, "compiled", data.NewNullValue())
}

func routeGetDefault(cv *data.ClassValue, name string) data.Value {
	v, ok := assocGet(propArray(cv, "defaults"), name)
	if !ok {
		return data.NewNullValue()
	}
	return v
}

func routeHasDefault(cv *data.ClassValue, name string) bool {
	return assocHas(propArray(cv, "defaults"), name)
}

func routeSetDefault(cv *data.ClassValue, name string, val data.Value) {
	if name == "_locale" && routeIsLocalized(cv) {
		return
	}
	defs := cloneArray(propArray(cv, "defaults"))
	assocSet(defs, name, val)
	setProp(cv, "defaults", defs)
	setProp(cv, "compiled", data.NewNullValue())
}

func routeAddRequirements(cv *data.ClassValue, extra *data.ArrayValue) data.Control {
	if extra == nil {
		return nil
	}
	reqs := cloneArray(propArray(cv, "requirements"))
	localized := routeIsLocalized(cv)
	for i, z := range extra.List {
		if z == nil || z.Value == nil {
			continue
		}
		key := arrayKey(z, i)
		if key == "_locale" && localized {
			continue
		}
		sanitized, err := sanitizeRequirement(key, z.Value.AsString())
		if err != nil {
			return data.NewErrorThrow(nil, err)
		}
		assocSet(reqs, key, data.NewStringValue(sanitized))
	}
	setProp(cv, "requirements", reqs)
	setProp(cv, "compiled", data.NewNullValue())
	return nil
}

func routeGetRequirement(cv *data.ClassValue, key string) string {
	v, ok := assocGet(propArray(cv, "requirements"), key)
	if !ok || v == nil || isNull(v) {
		return ""
	}
	return v.AsString()
}

func routeSetRequirement(cv *data.ClassValue, key, regex string) data.Control {
	if key == "_locale" && routeIsLocalized(cv) {
		return nil
	}
	sanitized, err := sanitizeRequirement(key, regex)
	if err != nil {
		return data.NewErrorThrow(nil, err)
	}
	reqs := cloneArray(propArray(cv, "requirements"))
	assocSet(reqs, key, data.NewStringValue(sanitized))
	setProp(cv, "requirements", reqs)
	setProp(cv, "compiled", data.NewNullValue())
	return nil
}

func routeCompile(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	compiled := prop(cv, "compiled")
	if compiled != nil && !isNull(compiled) {
		if asClassValue(compiled) != nil {
			return compiled, nil
		}
	}
	className := routeCompilerName
	if opt, ok := assocGet(propArray(cv, "options"), "compiler_class"); ok && opt != nil && !isNull(opt) {
		className = opt.AsString()
	}
	var out data.GetValue
	var ctl data.Control
	if className == routeCompilerName || className == "" {
		out, ctl = compileRouteValue(ctx, cv)
	} else {
		stmt, ok := ctx.GetVM().GetClass(className)
		if !ok {
			loaded, acl := ctx.GetVM().GetOrLoadClass(className)
			if acl != nil {
				return nil, acl
			}
			stmt = loaded
		}
		if stmt == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("compiler class %s not found", className))
		}
		var method data.Method
		if gsm, ok := stmt.(data.GetStaticMethod); ok {
			method, _ = gsm.GetStaticMethod("compile")
		}
		if method == nil {
			method, _ = stmt.GetMethod("compile")
		}
		if method == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("compiler class %s has no compile()", className))
		}
		fnCtx := ctx.CreateBaseContext().CreateContext(method.GetVariables())
		if len(method.GetVariables()) > 0 {
			_ = fnCtx.SetVariableValue(method.GetVariables()[0], cv)
		}
		out, ctl = method.Call(fnCtx)
	}
	if ctl != nil {
		return nil, ctl
	}
	if v, ok := out.(data.Value); ok {
		setProp(cv, "compiled", v)
	}
	return out, nil
}

func routeUnserialize(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	arr := argArray(ctx, 0)
	copyKeys := []string{"path", "host", "defaults", "requirements", "options", "schemes", "methods", "condition", "compiled"}
	for _, k := range copyKeys {
		if v, ok := assocGet(arr, k); ok {
			setProp(cv, k, v)
		}
	}
	return data.NewNullValue(), nil
}
