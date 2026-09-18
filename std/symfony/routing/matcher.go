package routing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

const (
	urlMatcherName               = "Symfony\\Component\\Routing\\Matcher\\UrlMatcher"
	compiledUrlMatcherName       = "Symfony\\Component\\Routing\\Matcher\\CompiledUrlMatcher"
	matcherDumperName            = "Symfony\\Component\\Routing\\Matcher\\Dumper\\MatcherDumper"
	compiledUrlMatcherDumperName = "Symfony\\Component\\Routing\\Matcher\\Dumper\\CompiledUrlMatcherDumper"
)

func newUrlMatcherClass() *rtClass {
	c := newRtClass(urlMatcherName, nil, []string{urlMatcherInterfaceName, requestMatcherIfaceName}, []data.Property{
		protProp("routes", data.NewNullValue()),
		protProp("context", data.NewNullValue()),
		protProp("allow", emptyArray()),
		protProp("allowSchemes", emptyArray()),
		protProp("request", data.NewNullValue()),
	})
	c.constInt("REQUIREMENT_MATCH", 0)
	c.constInt("REQUIREMENT_MISMATCH", 1)
	c.constInt("ROUTE_MATCH", 2)
	c.add(methDef("__construct",
		[]data.GetValue{
			param("routes", 0, nil, nil),
			param("context", 1, nil, nil),
		},
		[]data.Variable{
			variable("routes", 0, nil),
			variable("context", 1, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			setProp(cv, "routes", arg(ctx, 0))
			setProp(cv, "context", arg(ctx, 1))
			setProp(cv, "allow", phpList())
			setProp(cv, "allowSchemes", phpList())
			return data.NewNullValue(), nil
		}))
	c.add(meth("setContext", []string{"context"}, func(ctx data.Context) (data.GetValue, data.Control) {
		setProp(rtSelf(ctx), "context", arg(ctx, 0))
		return data.NewNullValue(), nil
	}))
	c.add(meth("getContext", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "context"), nil
	}))
	c.add(meth("match", []string{"pathinfo"}, urlMatcherMatch))
	c.add(meth("matchRequest", []string{"request"}, urlMatcherMatchRequest))
	return c
}

func urlMatcherMatch(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	setProp(cv, "allow", phpList())
	setProp(cv, "allowSchemes", phpList())
	pathinfo := phpRawURLDecode(argString(ctx, 0, "/"))
	if pathinfo == "" {
		pathinfo = "/"
	}
	ret, ctl := matchCollection(ctx, cv, pathinfo)
	if ctl != nil {
		return nil, ctl
	}
	if ret != nil && len(ret.List) > 0 {
		return ret, nil
	}
	allow := toStringSlice(prop(cv, "allow"))
	allowSchemes := toStringSlice(prop(cv, "allowSchemes"))
	if pathinfo == "/" && len(allow) == 0 && len(allowSchemes) == 0 {
		return nil, throwNamed(ctx, exNoConfiguration, "")
	}
	if len(allow) > 0 {
		uniq := uniqueStrings(allow)
		return nil, throwNamedArgs(ctx, exMethodNotAllowed, stringsToArray(uniq))
	}
	return nil, throwNamed(ctx, exResourceNotFound, fmt.Sprintf("No routes found for \"%s\".", pathinfo))
}

func urlMatcherMatchRequest(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	req := arg(ctx, 0)
	setProp(cv, "request", req)
	orig := prop(cv, "context")
	origCV := asClassValue(orig)
	if origCV != nil {
		cloned, ctl := cloneRequestContext(ctx, origCV)
		if ctl != nil {
			return nil, ctl
		}
		reqCV := asClassValue(req)
		if reqCV != nil {
			if _, ctl := callNamed(cloned, "fromRequest", reqCV); ctl != nil {
				return nil, ctl
			}
		}
		setProp(cv, "context", cloned)
		defer func() {
			setProp(cv, "context", orig)
			setProp(cv, "request", data.NewNullValue())
		}()
	}
	path := "/"
	if reqCV := asClassValue(req); reqCV != nil {
		if v, _ := callNamed(reqCV, "getPathInfo"); v != nil {
			path = v.AsString()
		}
	}
	return urlMatcherMatchArgs(ctx, path)
}

func urlMatcherMatchArgs(ctx data.Context, pathinfo string) (data.GetValue, data.Control) {
	fnCtx := ctx
	if cm, ok := ctx.(*data.ClassMethodContext); ok {
		fnCtx = cm
	}
	_ = fnCtx
	cv := rtSelf(ctx)
	setProp(cv, "allow", phpList())
	setProp(cv, "allowSchemes", phpList())
	pathinfo = phpRawURLDecode(pathinfo)
	if pathinfo == "" {
		pathinfo = "/"
	}
	ret, ctl := matchCollection(ctx, cv, pathinfo)
	if ctl != nil {
		return nil, ctl
	}
	if ret != nil && len(ret.List) > 0 {
		return ret, nil
	}
	allow := toStringSlice(prop(cv, "allow"))
	if len(allow) > 0 {
		return nil, throwNamedArgs(ctx, exMethodNotAllowed, stringsToArray(uniqueStrings(allow)))
	}
	if pathinfo == "/" {
		return nil, throwNamed(ctx, exNoConfiguration, "")
	}
	return nil, throwNamed(ctx, exResourceNotFound, fmt.Sprintf("No routes found for \"%s\".", pathinfo))
}

func matchCollection(ctx data.Context, matcher *data.ClassValue, pathinfo string) (*data.ArrayValue, data.Control) {
	context := asClassValue(prop(matcher, "context"))
	method := "GET"
	if context != nil {
		method = strings.ToUpper(propString(context, "method", "GET"))
	}
	if method == "HEAD" {
		method = "GET"
	}
	trimmed := strings.TrimRight(pathinfo, "/")
	if trimmed == "" {
		trimmed = "/"
	}

	routesCV := asClassValue(prop(matcher, "routes"))
	if routesCV == nil {
		return phpList(), nil
	}
	all, ctl := callNamed(routesCV, "all")
	if ctl != nil {
		return nil, ctl
	}
	allArr := valueToArray(all)
	if allArr == nil || len(allArr.List) == 0 {
		return phpList(), nil
	}

	host := ""
	scheme := "http"
	if context != nil {
		host = propString(context, "host", "")
		scheme = propString(context, "scheme", "http")
	}

	for i, z := range allArr.List {
		if z == nil {
			continue
		}
		name := arrayKey(z, i)
		route := asClassValue(z.Value)
		if route == nil {
			continue
		}
		compiledVal, ctl := callNamed(route, "compile")
		if ctl != nil {
			return nil, ctl
		}
		compiled := asClassValue(compiledVal)
		if compiled == nil {
			continue
		}
		staticPrefix := strings.TrimRight(propString(compiled, "staticPrefix", ""), "/")
		if staticPrefix != "" && !strings.HasPrefix(trimmed, staticPrefix) {
			continue
		}
		regex := propString(compiled, "regex", "")
		if regex == "" {
			continue
		}
		named, _, ok := pregMatchGroups(regex, pathinfo)
		if !ok {
			continue
		}
		hostRegex := propString(compiled, "hostRegex", "")
		hostNamed := map[string]string{}
		if hostRegex != "" {
			hn, _, hostOK := pregMatchGroups(hostRegex, host)
			if !hostOK {
				continue
			}
			hostNamed = hn
		}
		attrs := phpList()
		defaults := propArray(route, "defaults")
		canonical := ""
		if v, ok := assocGet(defaults, "_canonical_route"); ok && !isNull(v) {
			canonical = v.AsString()
		}
		routeName := name
		if canonical != "" {
			routeName = canonical
		}
		for _, z := range defaults.List {
			if z == nil || z.Name == "" || z.Name == "_canonical_route" {
				continue
			}
			assocSet(attrs, z.Name, z.Value)
		}
		for k, v := range named {
			if k == "" {
				continue
			}
			assocSet(attrs, k, data.NewStringValue(v))
		}
		for k, v := range hostNamed {
			if k == "" {
				continue
			}
			assocSet(attrs, k, data.NewStringValue(v))
		}
		assocSet(attrs, "_route", data.NewStringValue(routeName))

		if cond := propString(route, "condition", ""); cond != "" {
			continue
		}
		schemes := toStringSlice(prop(route, "schemes"))
		if len(schemes) > 0 && !inStringSlice(schemes, scheme) {
			cur := propArray(matcher, "allowSchemes")
			for _, s := range schemes {
				cur.List = append(cur.List, data.NewZVal(data.NewStringValue(s)))
			}
			setProp(matcher, "allowSchemes", cur)
			continue
		}
		required := toStringSlice(prop(route, "methods"))
		if len(required) > 0 && !inStringSlice(required, method) {
			cur := propArray(matcher, "allow")
			for _, m := range required {
				cur.List = append(cur.List, data.NewZVal(data.NewStringValue(m)))
			}
			setProp(matcher, "allow", cur)
			continue
		}
		return attrs, nil
	}
	return phpList(), nil
}

func uniqueStrings(ss []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func newCompiledUrlMatcherClass() *rtClass {
	c := newRtClass(compiledUrlMatcherName, strPtr(urlMatcherName), []string{urlMatcherInterfaceName, requestMatcherIfaceName}, []data.Property{
		privProp("matchHost", data.NewBoolValue(false)),
		privProp("staticRoutes", emptyArray()),
		privProp("regexpList", emptyArray()),
		privProp("dynamicRoutes", emptyArray()),
		privProp("checkCondition", data.NewNullValue()),
		protProp("context", data.NewNullValue()),
		protProp("request", data.NewNullValue()),
		protProp("allow", emptyArray()),
		protProp("allowSchemes", emptyArray()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{
			param("compiledRoutes", 0, nil, nil),
			param("context", 1, nil, nil),
		},
		[]data.Variable{
			variable("compiledRoutes", 0, nil),
			variable("context", 1, nil),
		},
		compiledMatcherConstruct))
	c.add(meth("match", []string{"pathinfo"}, compiledMatcherMatch))
	c.add(meth("setContext", []string{"context"}, func(ctx data.Context) (data.GetValue, data.Control) {
		setProp(rtSelf(ctx), "context", arg(ctx, 0))
		return data.NewNullValue(), nil
	}))
	c.add(meth("getContext", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "context"), nil
	}))
	c.add(meth("matchRequest", []string{"request"}, compiledMatcherMatchRequest))
	return c
}

func compiledMatcherConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	compiled := argArray(ctx, 0)
	setProp(cv, "context", arg(ctx, 1))
	getIdx := func(i int) data.Value {
		if compiled == nil || i >= len(compiled.List) || compiled.List[i] == nil {
			return data.NewNullValue()
		}
		return compiled.List[i].Value
	}
	matchHost := false
	if v := getIdx(0); v != nil {
		if bv, ok := v.(data.AsBool); ok {
			matchHost, _ = bv.AsBool()
		}
	}
	setProp(cv, "matchHost", data.NewBoolValue(matchHost))
	setProp(cv, "staticRoutes", valueToArray(getIdx(1)))
	setProp(cv, "regexpList", valueToArray(getIdx(2)))
	setProp(cv, "dynamicRoutes", valueToArray(getIdx(3)))
	setProp(cv, "checkCondition", getIdx(4))
	return data.NewNullValue(), nil
}

func compiledMatcherMatchRequest(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	req := arg(ctx, 0)
	setProp(cv, "request", req)
	orig := prop(cv, "context")
	if origCV := asClassValue(orig); origCV != nil {
		cloned, ctl := cloneRequestContext(ctx, origCV)
		if ctl != nil {
			return nil, ctl
		}
		if reqCV := asClassValue(req); reqCV != nil {
			if _, ctl := callNamed(cloned, "fromRequest", reqCV); ctl != nil {
				return nil, ctl
			}
		}
		setProp(cv, "context", cloned)
		defer func() {
			setProp(cv, "context", orig)
			setProp(cv, "request", data.NewNullValue())
		}()
	}
	path := "/"
	if reqCV := asClassValue(req); reqCV != nil {
		if v, _ := callNamed(reqCV, "getPathInfo"); v != nil {
			path = v.AsString()
		}
	}
	path = phpRawURLDecode(path)
	if path == "" {
		path = "/"
	}
	ret, ctl := compiledDoMatch(ctx, cv, path)
	if ctl != nil {
		return nil, ctl
	}
	if ret != nil && len(ret.List) > 0 {
		return ret, nil
	}
	allow := toStringSlice(prop(cv, "allow"))
	if len(allow) > 0 {
		return nil, throwNamedArgs(ctx, exMethodNotAllowed, stringsToArray(uniqueStrings(allow)))
	}
	return nil, throwNamed(ctx, exResourceNotFound, fmt.Sprintf("No routes found for \"%s\".", path))
}

func compiledMatcherMatch(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	pathinfo := phpRawURLDecode(argString(ctx, 0, "/"))
	if pathinfo == "" {
		pathinfo = "/"
	}
	ret, ctl := compiledDoMatch(ctx, cv, pathinfo)
	if ctl != nil {
		return nil, ctl
	}
	if ret != nil && len(ret.List) > 0 {
		return ret, nil
	}
	allow := toStringSlice(prop(cv, "allow"))
	if len(allow) > 0 {
		return nil, throwNamedArgs(ctx, exMethodNotAllowed, stringsToArray(uniqueStrings(allow)))
	}
	if pathinfo == "/" {
		return nil, throwNamed(ctx, exNoConfiguration, "")
	}
	return nil, throwNamed(ctx, exResourceNotFound, fmt.Sprintf("No routes found for \"%s\".", pathinfo))
}

func compiledDoMatch(ctx data.Context, matcher *data.ClassValue, pathinfo string) (*data.ArrayValue, data.Control) {
	setProp(matcher, "allow", phpList())
	setProp(matcher, "allowSchemes", phpList())
	trimmed := strings.TrimRight(pathinfo, "/")
	if trimmed == "" {
		trimmed = "/"
	}
	context := asClassValue(prop(matcher, "context"))
	method := "GET"
	scheme := "http"
	host := ""
	if context != nil {
		method = strings.ToUpper(propString(context, "method", "GET"))
		scheme = propString(context, "scheme", "http")
		host = strings.ToLower(propString(context, "host", ""))
	}
	canonical := method
	if method == "HEAD" {
		canonical = "GET"
	}

	staticRoutes := valueToArray(prop(matcher, "staticRoutes"))
	if staticRoutes == nil {
		staticRoutes = phpList()
	}
	if bucket, ok := assocGet(staticRoutes, trimmed); ok {
		list := valueToArray(bucket)
		for _, z := range list.List {
			if z == nil {
				continue
			}
			row := valueToArray(z.Value)
			if row == nil || len(row.List) < 5 {
				continue
			}
			ret := valueToArray(row.List[0].Value)
			requiredHost := ""
			if row.List[1] != nil && row.List[1].Value != nil && !isNull(row.List[1].Value) {
				requiredHost = row.List[1].Value.AsString()
			}
			if requiredHost != "" && requiredHost != host {
				if !strings.HasPrefix(requiredHost, "{") || !pregMatch(requiredHost, host) {
					continue
				}
			}
			requiredMethods := valueToArray(row.List[2].Value)
			requiredSchemes := valueToArray(row.List[3].Value)
			if requiredSchemes != nil && len(requiredSchemes.List) > 0 && !assocHas(requiredSchemes, scheme) {
				continue
			}
			if requiredMethods != nil && len(requiredMethods.List) > 0 &&
				!assocHas(requiredMethods, canonical) && !assocHas(requiredMethods, method) {
				for _, mz := range requiredMethods.List {
					if mz != nil && mz.Name != "" {
						cur := propArray(matcher, "allow")
						cur.List = append(cur.List, data.NewZVal(data.NewStringValue(mz.Name)))
						setProp(matcher, "allow", cur)
					}
				}
				continue
			}
			return cloneArray(ret), nil
		}
	}

	regexpList := valueToArray(prop(matcher, "regexpList"))
	dynamicRoutes := valueToArray(prop(matcher, "dynamicRoutes"))
	if regexpList == nil || len(regexpList.List) == 0 {
		return phpList(), nil
	}
	matchedPath := pathinfo
	if bv, ok := prop(matcher, "matchHost").(data.AsBool); ok {
		if b, _ := bv.AsBool(); b {
			matchedPath = host + "." + pathinfo
		}
	}
	for i, z := range regexpList.List {
		if z == nil || z.Value == nil {
			continue
		}
		regex := z.Value.AsString()
		named, indexed, ok := pregMatchGroups(regex, matchedPath)
		if !ok {
			continue
		}
		mark := arrayKey(z, i)
		if m, ok := named["MARK"]; ok && m != "" {
			mark = m
		}
		bucket, ok := assocGet(dynamicRoutes, mark)
		if !ok {
			if n, err := strconv.Atoi(mark); err == nil {
				if dynamicRoutes != nil && n < len(dynamicRoutes.List) && dynamicRoutes.List[n] != nil {
					bucket = dynamicRoutes.List[n].Value
					ok = true
				}
			}
		}
		if !ok {
			continue
		}
		list := valueToArray(bucket)
		if list == nil {
			continue
		}
		for _, rz := range list.List {
			if rz == nil {
				continue
			}
			row := valueToArray(rz.Value)
			if row == nil || len(row.List) < 2 {
				continue
			}
			ret := valueToArray(row.List[0].Value)
			vars := valueToArray(row.List[1].Value)
			out := cloneArray(ret)
			if vars != nil {
				for vi, vz := range vars.List {
					if vz == nil || vz.Value == nil {
						continue
					}
					name := vz.Value.AsString()
					if vi+1 < len(indexed) {
						assocSet(out, name, data.NewStringValue(indexed[vi+1]))
					}
				}
			}
			requiredMethods := phpList()
			requiredSchemes := phpList()
			if len(row.List) > 2 && row.List[2] != nil {
				requiredMethods = valueToArray(row.List[2].Value)
			}
			if len(row.List) > 3 && row.List[3] != nil {
				requiredSchemes = valueToArray(row.List[3].Value)
			}
			if requiredSchemes != nil && len(requiredSchemes.List) > 0 && !assocHas(requiredSchemes, scheme) {
				continue
			}
			if requiredMethods != nil && len(requiredMethods.List) > 0 &&
				!assocHas(requiredMethods, canonical) && !assocHas(requiredMethods, method) {
				continue
			}
			return out, nil
		}
	}
	return phpList(), nil
}

func newMatcherDumperClass() *rtClass {
	c := newRtClass(matcherDumperName, nil, []string{matcherDumperIfaceName}, []data.Property{
		privProp("routes", data.NewNullValue()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{param("routes", 0, nil, nil)},
		[]data.Variable{variable("routes", 0, nil)},
		func(ctx data.Context) (data.GetValue, data.Control) {
			setProp(rtSelf(ctx), "routes", arg(ctx, 0))
			return data.NewNullValue(), nil
		}))
	c.add(meth("getRoutes", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "routes"), nil
	}))
	c.add(meth("dump", []string{"options"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue("<?php\nreturn [];\n"), nil
	}))
	return c
}

func newCompiledUrlMatcherDumperClass() *rtClass {
	c := newRtClass(compiledUrlMatcherDumperName, strPtr(matcherDumperName), []string{matcherDumperIfaceName}, []data.Property{
		privProp("routes", data.NewNullValue()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{param("routes", 0, nil, nil)},
		[]data.Variable{variable("routes", 0, nil)},
		func(ctx data.Context) (data.GetValue, data.Control) {
			setProp(rtSelf(ctx), "routes", arg(ctx, 0))
			return data.NewNullValue(), nil
		}))
	c.add(meth("getRoutes", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "routes"), nil
	}))
	c.add(meth("getCompiledRoutes", []string{"forDump"}, dumperGetCompiledRoutes))
	c.add(meth("dump", []string{"options"}, func(ctx data.Context) (data.GetValue, data.Control) {
		if _, ctl := dumperGetCompiledRoutes(ctx); ctl != nil {
			return nil, ctl
		}
		return data.NewStringValue("<?php\nreturn [];\n"), nil
	}))
	return c
}

func dumperGetCompiledRoutes(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	col := asClassValue(prop(cv, "routes"))
	if col == nil {
		return phpList(data.NewBoolValue(false), phpList(), phpList(), phpList(), data.NewNullValue()), nil
	}
	all, ctl := callNamed(col, "all")
	if ctl != nil {
		return nil, ctl
	}
	allArr := valueToArray(all)
	staticRoutes := &data.ArrayValue{}
	regexpList := &data.ArrayValue{}
	dynamicRoutes := &data.ArrayValue{}
	matchHost := false
	mark := 0
	if allArr != nil {
		for i, z := range allArr.List {
			if z == nil {
				continue
			}
			name := arrayKey(z, i)
			route := asClassValue(z.Value)
			if route == nil {
				continue
			}
			if propString(route, "host", "") != "" {
				matchHost = true
			}
			compiledVal, ctl := callNamed(route, "compile")
			if ctl != nil {
				return nil, ctl
			}
			compiled := asClassValue(compiledVal)
			if compiled == nil {
				continue
			}
			pathVars := propArray(compiled, "pathVariables")
			path := propString(route, "path", "/")
			hasTrailingSlash := path != "/" && strings.HasSuffix(path, "/")
			url := path
			if hasTrailingSlash {
				url = strings.TrimRight(path, "/")
				if url == "" {
					url = "/"
				}
			}
			defaults := cloneArray(propArray(route, "defaults"))
			assocUnset(defaults, "_canonical_route")
			ret := phpAssoc("_route", data.NewStringValue(name))
			for _, dz := range defaults.List {
				if dz != nil && dz.Name != "" {
					assocSet(ret, dz.Name, dz.Value)
				}
			}
			methods := flipStrings(toStringSlice(prop(route, "methods")))
			schemes := flipStrings(toStringSlice(prop(route, "schemes")))
			var methodsV, schemesV data.Value = data.NewNullValue(), data.NewNullValue()
			if methods != nil {
				methodsV = methods
			}
			if schemes != nil {
				schemesV = schemes
			}
			host := propString(route, "host", "")
			var hostV data.Value = data.NewNullValue()
			if host != "" {
				if len(propArray(compiled, "hostVariables").List) > 0 {
					hostV = data.NewStringValue(propString(compiled, "hostRegex", ""))
				} else {
					hostV = data.NewStringValue(strings.ToLower(host))
				}
			}
			if len(pathVars.List) == 0 {
				row := phpList(ret, hostV, methodsV, schemesV, data.NewBoolValue(hasTrailingSlash), data.NewBoolValue(false), data.NewNullValue())
				bucket, _ := assocGet(staticRoutes, url)
				bArr, _ := bucket.(*data.ArrayValue)
				if bArr == nil {
					bArr = phpList()
				}
				bArr.List = append(bArr.List, data.NewZVal(row))
				assocSet(staticRoutes, url, bArr)
				continue
			}
			regex := propString(compiled, "regex", "")
			assocSet(regexpList, strconv.Itoa(mark), data.NewStringValue(regex))
			row := phpList(ret, pathVars, methodsV, schemesV, data.NewBoolValue(hasTrailingSlash), data.NewBoolValue(false), data.NewNullValue())
			assocSet(dynamicRoutes, strconv.Itoa(mark), phpList(row))
			mark++
		}
	}
	return phpList(
		data.NewBoolValue(matchHost),
		staticRoutes,
		regexpList,
		dynamicRoutes,
		data.NewNullValue(),
	), nil
}
