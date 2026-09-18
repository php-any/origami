package routing

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/php-any/origami/data"
)

const urlGeneratorName = "Symfony\\Component\\Routing\\Generator\\UrlGenerator"

var decodedPathChars = map[string]string{
	"%2F":   "/",
	"%252F": "%2F",
	"%40":   "@",
	"%3A":   ":",
	"%3B":   ";",
	"%2C":   ",",
	"%3D":   "=",
	"%2B":   "+",
	"%21":   "!",
	"%2A":   "*",
	"%7C":   "|",
}

var decodedQueryChars = map[string]string{
	"%2F":   "/",
	"%252F": "%2F",
	"%3F":   "?",
	"%40":   "@",
	"%3A":   ":",
	"%21":   "!",
	"%3B":   ";",
	"%2C":   ",",
	"%2A":   "*",
}

func newUrlGeneratorClass() *rtClass {
	c := newRtClass(urlGeneratorName, nil, []string{urlGeneratorIfaceName, configurableReqName}, []data.Property{
		protProp("routes", data.NewNullValue()),
		protProp("context", data.NewNullValue()),
		protProp("strictRequirements", data.NewBoolValue(true)),
		privProp("defaultLocale", data.NewNullValue()),
	})
	c.constInt("ABSOLUTE_URL", refAbsoluteURL)
	c.constInt("ABSOLUTE_PATH", refAbsolutePath)
	c.constInt("RELATIVE_PATH", refRelativePath)
	c.constInt("NETWORK_PATH", refNetworkPath)
	c.add(methDef("__construct",
		[]data.GetValue{
			param("routes", 0, nil, nil),
			param("context", 1, nil, nil),
			param("logger", 2, data.NewNullValue(), nil),
			param("defaultLocale", 3, data.NewNullValue(), nil),
		},
		[]data.Variable{
			variable("routes", 0, nil),
			variable("context", 1, nil),
			variable("logger", 2, nil),
			variable("defaultLocale", 3, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			setProp(cv, "routes", arg(ctx, 0))
			setProp(cv, "context", arg(ctx, 1))
			setProp(cv, "strictRequirements", data.NewBoolValue(true))
			loc := arg(ctx, 3)
			if loc == nil {
				loc = data.NewNullValue()
			}
			setProp(cv, "defaultLocale", loc)
			return data.NewNullValue(), nil
		}))
	c.add(meth("setContext", []string{"context"}, func(ctx data.Context) (data.GetValue, data.Control) {
		setProp(rtSelf(ctx), "context", arg(ctx, 0))
		return data.NewNullValue(), nil
	}))
	c.add(meth("getContext", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "context"), nil
	}))
	c.add(meth("setStrictRequirements", []string{"enabled"}, func(ctx data.Context) (data.GetValue, data.Control) {
		v := arg(ctx, 0)
		if v == nil || isNull(v) {
			setProp(rtSelf(ctx), "strictRequirements", data.NewNullValue())
		} else {
			setProp(rtSelf(ctx), "strictRequirements", v)
		}
		return data.NewNullValue(), nil
	}))
	c.add(meth("isStrictRequirements", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "strictRequirements"), nil
	}))
	c.add(meth("generate", []string{"name", "parameters", "referenceType"}, generatorGenerate))
	c.add(staticMeth("getRelativePath", []string{"basePath", "targetPath"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(getRelativePath(argString(ctx, 0, ""), argString(ctx, 1, ""))), nil
	}))
	return c
}

func generatorGenerate(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	name := argString(ctx, 0, "")
	parameters := cloneArray(argArray(ctx, 1))
	refType := argInt(ctx, 2, refAbsolutePath)
	routes := asClassValue(prop(cv, "routes"))
	if routes == nil {
		return nil, throwNamed(ctx, exRouteNotFound, fmt.Sprintf("Unable to generate a URL for the named route \"%s\" as such route does not exist.", name))
	}
	locale := ""
	if v, ok := assocGet(parameters, "_locale"); ok && !isNull(v) {
		locale = v.AsString()
	} else if context := asClassValue(prop(cv, "context")); context != nil {
		if pv, ok := assocGet(propArray(context, "parameters"), "_locale"); ok && !isNull(pv) {
			locale = pv.AsString()
		}
	}
	if locale == "" {
		if v := prop(cv, "defaultLocale"); v != nil && !isNull(v) {
			locale = v.AsString()
		}
	}
	var route *data.ClassValue
	if locale != "" {
		loc := locale
		for loc != "" {
			got, ctl := callNamed(routes, "get", data.NewStringValue(name+"."+loc))
			if ctl != nil {
				return nil, ctl
			}
			if r := asClassValue(got); r != nil {
				canon := routeGetDefault(r, "_canonical_route")
				if !isNull(canon) && canon.AsString() == name {
					route = r
					break
				}
			}
			if i := strings.IndexByte(loc, '_'); i >= 0 {
				loc = loc[:i]
			} else {
				break
			}
		}
	}
	if route == nil {
		got, ctl := callNamed(routes, "get", data.NewStringValue(name))
		if ctl != nil {
			return nil, ctl
		}
		route = asClassValue(got)
	}
	if route == nil {
		return nil, throwNamed(ctx, exRouteNotFound, fmt.Sprintf("Unable to generate a URL for the named route \"%s\" as such route does not exist.", name))
	}
	compiledVal, ctl := callNamed(route, "compile")
	if ctl != nil {
		return nil, ctl
	}
	compiled := asClassValue(compiledVal)
	if compiled == nil {
		return nil, throwNamed(ctx, exRouteNotFound, fmt.Sprintf("Unable to generate a URL for the named route \"%s\" as such route does not exist.", name))
	}
	defaults := propArray(route, "defaults")
	variables := toStringSlice(prop(compiled, "variables"))
	if assocHas(defaults, "_canonical_route") && assocHas(defaults, "_locale") {
		if !inStringSlice(variables, "_locale") {
			assocUnset(parameters, "_locale")
		} else if !assocHas(parameters, "_locale") {
			if v, ok := assocGet(defaults, "_locale"); ok {
				assocSet(parameters, "_locale", v)
			}
		}
	}
	url, ctl := doGenerate(ctx, cv, route, compiled, parameters, name, refType)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(url), nil
}

func doGenerate(ctx data.Context, gen, route, compiled *data.ClassValue, parameters *data.ArrayValue, name string, refType int) (string, data.Control) {
	queryParameters := phpList()
	if q, ok := assocGet(parameters, "_query"); ok {
		if av, ok := q.(*data.ArrayValue); ok {
			queryParameters = av
			assocUnset(parameters, "_query")
		} else {
			return "", throwNamed(ctx, exInvalidParameter, "Parameter \"_query\" must be an array of query parameters.")
		}
	}
	defaults := propArray(route, "defaults")
	context := asClassValue(prop(gen, "context"))
	merged := cloneArray(defaults)
	if context != nil {
		for _, z := range propArray(context, "parameters").List {
			if z != nil && z.Name != "" {
				assocSet(merged, z.Name, z.Value)
			}
		}
	}
	for i, z := range parameters.List {
		if z == nil {
			continue
		}
		assocSet(merged, arrayKey(z, i), z.Value)
	}
	variables := toStringSlice(prop(compiled, "variables"))
	var missing []string
	for _, v := range variables {
		if !assocHas(merged, v) {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		vals := make([]data.Value, len(missing))
		for i, m := range missing {
			vals[i] = data.NewStringValue(m)
		}
		return "", throwNamedArgs(ctx, exMissingMandatory, data.NewStringValue(name), phpList(vals...))
	}

	tokens := propArray(compiled, "tokens")
	url := ""
	optional := true
	strict := true
	if sv := prop(gen, "strictRequirements"); sv != nil && !isNull(sv) {
		if bv, ok := sv.(data.AsBool); ok {
			strict, _ = bv.AsBool()
		}
	} else if isNull(prop(gen, "strictRequirements")) {
		strict = false
	}
	for _, z := range tokens.List {
		if z == nil {
			continue
		}
		tok, _ := z.Value.(*data.ArrayValue)
		if tok == nil || len(tok.List) < 2 {
			continue
		}
		kind := tok.List[0].Value.AsString()
		if kind == "variable" {
			varName := ""
			if len(tok.List) > 3 && tok.List[3] != nil {
				varName = tok.List[3].Value.AsString()
			}
			important := false
			if len(tok.List) > 5 && tok.List[5] != nil {
				if bv, ok := tok.List[5].Value.(data.AsBool); ok {
					important, _ = bv.AsBool()
				}
			}
			val, _ := assocGet(merged, varName)
			def, hasDef := assocGet(defaults, varName)
			sameDefault := hasDef && val != nil && def != nil && val.AsString() == def.AsString()
			if !optional || important || !hasDef || (val != nil && !isNull(val) && !sameDefault) {
				req := ""
				if len(tok.List) > 2 && tok.List[2] != nil {
					req = tok.List[2].Value.AsString()
				}
				given := ""
				if val != nil {
					given = val.AsString()
				}
				if req != "" && !isNull(prop(gen, "strictRequirements")) {
					pat := "#^(?:" + req + ")$#i"
					if len(tok.List) > 4 && tok.List[4] != nil {
						if bv, ok := tok.List[4].Value.(data.AsBool); ok {
							if u, _ := bv.AsBool(); u {
								pat += "u"
							}
						}
					}
					if !pregMatch(pat, given) {
						if strict {
							return "", throwNamed(ctx, exInvalidParameter, fmt.Sprintf("Parameter \"%s\" for route \"%s\" must match \"%s\" (\"%s\" given) to generate a corresponding URL.", varName, name, req, given))
						}
						return "", nil
					}
				}
				sep := tok.List[1].Value.AsString()
				url = sep + given + url
				optional = false
			}
		} else {
			url = tok.List[1].Value.AsString() + url
			optional = false
		}
	}
	if url == "" {
		url = "/"
	}
	url = applyDecoded(phpRawURLEncode(url), decodedPathChars)
	if strings.Contains(url, "/.") {
		parts := strings.Split(url, "/")
		for i, p := range parts {
			if p == "." {
				parts[i] = "%2E"
			} else if p == ".." {
				parts[i] = "%2E%2E"
			}
		}
		url = strings.Join(parts, "/")
	}

	schemeAuthority := ""
	host := ""
	scheme := "http"
	if context != nil {
		host = propString(context, "host", "")
		scheme = propString(context, "scheme", "http")
	}
	requiredSchemes := toStringSlice(prop(route, "schemes"))
	if len(requiredSchemes) > 0 && !inStringSlice(requiredSchemes, scheme) {
		refType = refAbsoluteURL
		scheme = requiredSchemes[0]
	}
	hostTokens := propArray(compiled, "hostTokens")
	if len(hostTokens.List) > 0 {
		routeHost := ""
		for _, z := range hostTokens.List {
			tok, _ := z.Value.(*data.ArrayValue)
			if tok == nil || len(tok.List) < 2 {
				continue
			}
			if tok.List[0].Value.AsString() == "variable" {
				varName := tok.List[3].Value.AsString()
				val, _ := assocGet(merged, varName)
				s := ""
				if val != nil {
					s = val.AsString()
				}
				routeHost = tok.List[1].Value.AsString() + s + routeHost
			} else {
				routeHost = tok.List[1].Value.AsString() + routeHost
			}
		}
		if routeHost != host {
			host = routeHost
			if refType != refAbsoluteURL {
				refType = refNetworkPath
			}
		}
	}
	if refType == refAbsoluteURL || refType == refNetworkPath {
		if host != "" || (scheme != "" && scheme != "http" && scheme != "https") {
			port := ""
			if context != nil {
				if scheme == "http" && propInt(context, "httpPort", 80) != 80 {
					port = fmt.Sprintf(":%d", propInt(context, "httpPort", 80))
				} else if scheme == "https" && propInt(context, "httpsPort", 443) != 443 {
					port = fmt.Sprintf(":%d", propInt(context, "httpsPort", 443))
				}
			}
			if refType == refNetworkPath || scheme == "" {
				schemeAuthority = "//"
			} else {
				schemeAuthority = scheme + "://"
			}
			schemeAuthority += host + port
		}
	}
	if refType == refRelativePath {
		base := "/"
		if context != nil {
			base = propString(context, "pathInfo", "/")
		}
		url = getRelativePath(base, url)
	} else {
		baseUrl := ""
		if context != nil {
			baseUrl = propString(context, "baseUrl", "")
		}
		url = schemeAuthority + baseUrl + url
	}

	extra := phpList()
	varSet := map[string]bool{}
	for _, v := range variables {
		varSet[v] = true
	}
	for _, z := range parameters.List {
		if z == nil || z.Name == "" || varSet[z.Name] {
			continue
		}
		if dv, ok := assocGet(defaults, z.Name); ok && dv != nil && z.Value != nil && dv.AsString() == z.Value.AsString() {
			continue
		}
		assocSet(extra, z.Name, z.Value)
	}
	for _, z := range queryParameters.List {
		if z != nil && z.Name != "" {
			assocSet(extra, z.Name, z.Value)
		}
	}
	fragment := ""
	if v, ok := assocGet(defaults, "_fragment"); ok && !isNull(v) {
		fragment = v.AsString()
	}
	if v, ok := assocGet(extra, "_fragment"); ok {
		if !isNull(v) {
			fragment = v.AsString()
		}
		assocUnset(extra, "_fragment")
	}
	if len(extra.List) > 0 {
		if q := httpBuildQueryRFC3986(extra); q != "" {
			url += "?" + applyDecoded(q, decodedQueryChars)
		}
	}
	if fragment != "" {
		url += "#" + applyDecoded(phpRawURLEncode(fragment), decodedQueryChars)
	}
	return url, nil
}

func applyDecoded(s string, table map[string]string) string {
	for k, v := range table {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func httpBuildQueryRFC3986(arr *data.ArrayValue) string {
	if arr == nil || len(arr.List) == 0 {
		return ""
	}
	var parts []string
	var walk func(prefix string, v data.Value)
	walk = func(prefix string, v data.Value) {
		if av, ok := v.(*data.ArrayValue); ok {
			for i, z := range av.List {
				if z == nil {
					continue
				}
				key := arrayKey(z, i)
				next := prefix + "[" + url.QueryEscape(key) + "]"
				if prefix == "" {
					next = url.QueryEscape(key)
				}
				walk(next, z.Value)
			}
			return
		}
		val := ""
		if v != nil && !isNull(v) {
			val = v.AsString()
		}
		parts = append(parts, prefix+"="+phpRawURLEncode(val))
	}
	for i, z := range arr.List {
		if z == nil {
			continue
		}
		walk(url.QueryEscape(arrayKey(z, i)), z.Value)
	}
	return strings.Join(parts, "&")
}

func getRelativePath(basePath, targetPath string) string {
	if basePath == targetPath {
		return ""
	}
	trimSlash := func(s string) string {
		if s != "" && s[0] == '/' {
			return s[1:]
		}
		return s
	}
	sourceDirs := strings.Split(trimSlash(basePath), "/")
	targetDirs := strings.Split(trimSlash(targetPath), "/")
	if len(sourceDirs) > 0 {
		sourceDirs = sourceDirs[:len(sourceDirs)-1]
	}
	targetFile := ""
	if len(targetDirs) > 0 {
		targetFile = targetDirs[len(targetDirs)-1]
		targetDirs = targetDirs[:len(targetDirs)-1]
	}
	i := 0
	for i < len(sourceDirs) && i < len(targetDirs) && sourceDirs[i] == targetDirs[i] {
		i++
	}
	sourceDirs = sourceDirs[i:]
	targetDirs = targetDirs[i:]
	targetDirs = append(targetDirs, targetFile)
	path := strings.Repeat("../", len(sourceDirs)) + strings.Join(targetDirs, "/")
	if path == "" || path[0] == '/' {
		return "./" + path
	}
	colon := strings.IndexByte(path, ':')
	slash := strings.IndexByte(path, '/')
	if colon >= 0 && (slash < 0 || colon < slash) {
		return "./" + path
	}
	return path
}
