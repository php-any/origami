package routing

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/php-any/origami/data"
)

const (
	routeCompilerName       = "Symfony\\Component\\Routing\\RouteCompiler"
	routeCompilerSeparators = "/,;.:-_~+*=@|"
	variableMaxLength       = 32
)

var (
	compileVarRe   = regexp.MustCompile(`\{(!)?([A-Za-z0-9_\x80-\xFF]+)\}`)
	inlineVarRe    = regexp.MustCompile(`\{(!?)([A-Za-z0-9_\x80-\xFF]+)(?::([A-Za-z0-9_\x80-\xFF]+)(\.[A-Za-z0-9_\x80-\xFF]+)?)?(<[^>]*>)?(\?[^}]*)?\}`)
	digitStartRe   = regexp.MustCompile(`^\d`)
	nextVarStartRe = regexp.MustCompile(`^\{[A-Za-z0-9_\x80-\xFF]+\}`)
	highByteRe     = regexp.MustCompile(`[\x80-\xFF]`)
	utf8ReqRe      = regexp.MustCompile(`[\x80-\xFF]|\\[pPX]|\\x\{`)
)

func newRouteCompilerClass() *rtClass {
	c := newRtClass(routeCompilerName, nil, []string{routeCompilerIfaceName}, nil)
	c.constStr("SEPARATORS", routeCompilerSeparators)
	c.constInt("VARIABLE_MAXIMUM_LENGTH", variableMaxLength)
	c.add(staticMeth("compile", []string{"route"}, func(ctx data.Context) (data.GetValue, data.Control) {
		route := asClassValue(arg(ctx, 0))
		if route == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("RouteCompiler::compile expects a Route"))
		}
		return compileRouteValue(ctx, route)
	}))
	return c
}

func compileRouteValue(ctx data.Context, route *data.ClassValue) (data.GetValue, data.Control) {
	hostVariables := phpList()
	variables := phpList()
	var hostRegex string
	hostTokens := phpList()

	host := propString(route, "host", "")
	if host != "" {
		res, ctl := compilePattern(ctx, route, host, true)
		if ctl != nil {
			return nil, ctl
		}
		hostVariables = res.variables
		variables = cloneArray(hostVariables)
		hostTokens = res.tokens
		hostRegex = res.regex
	}

	locale := routeGetDefault(route, "_locale")
	canonical := routeGetDefault(route, "_canonical_route")
	if !isNull(locale) && !isNull(canonical) {
		req, _ := assocGet(propArray(route, "requirements"), "_locale")
		if req != nil && req.AsString() == regexp.QuoteMeta(locale.AsString()) {
			reqs := cloneArray(propArray(route, "requirements"))
			assocUnset(reqs, "_locale")
			setProp(route, "requirements", reqs)
			path := strings.ReplaceAll(propString(route, "path", "/"), "{_locale}", locale.AsString())
			setProp(route, "path", data.NewStringValue(path))
			setProp(route, "compiled", data.NewNullValue())
		}
	}

	path := propString(route, "path", "/")
	res, ctl := compilePattern(ctx, route, path, false)
	if ctl != nil {
		return nil, ctl
	}
	for _, z := range res.variables.List {
		if z != nil && z.Value != nil && z.Value.AsString() == "_fragment" {
			return nil, throwNamed(ctx, "InvalidArgumentException",
				fmt.Sprintf("Route pattern \"%s\" cannot contain \"_fragment\" as a path parameter.", path))
		}
	}

	allVars := cloneArray(variables)
	seen := map[string]bool{}
	for _, z := range allVars.List {
		if z != nil && z.Value != nil {
			seen[z.Value.AsString()] = true
		}
	}
	for _, z := range res.variables.List {
		if z == nil || z.Value == nil {
			continue
		}
		name := z.Value.AsString()
		if seen[name] {
			continue
		}
		seen[name] = true
		allVars.List = append(allVars.List, data.NewZVal(data.NewStringValue(name)))
	}

	return newCompiledRoute(ctx, res.staticPrefix, res.regex, res.tokens, res.variables, hostRegex, hostTokens, hostVariables, allVars)
}

type compiledPattern struct {
	staticPrefix string
	regex        string
	tokens       *data.ArrayValue
	variables    *data.ArrayValue
}

type tokenBuf struct {
	kind      string
	text      string
	regexp    string
	varName   string
	utf8      bool
	important bool
}

func (t tokenBuf) toValue() data.Value {
	if t.kind == "text" {
		return phpList(data.NewStringValue("text"), data.NewStringValue(t.text))
	}
	items := []data.Value{
		data.NewStringValue("variable"),
		data.NewStringValue(t.text),
		data.NewStringValue(t.regexp),
		data.NewStringValue(t.varName),
	}
	if t.utf8 || t.important {
		items = append(items, data.NewBoolValue(t.utf8))
	}
	if t.important {
		if len(items) == 4 {
			items = append(items, data.NewBoolValue(false))
		}
		items = append(items, data.NewBoolValue(true))
	}
	return phpList(items...)
}

func compilePattern(ctx data.Context, route *data.ClassValue, pattern string, isHost bool) (*compiledPattern, data.Control) {
	tokens := make([]tokenBuf, 0)
	variables := make([]string, 0)
	pos := 0
	defaultSep := "/"
	if isHost {
		defaultSep = "."
	}
	useUtf8 := utf8.ValidString(pattern)
	needsUtf8 := false
	if opt, ok := assocGet(propArray(route, "options"), "utf8"); ok && !isNull(opt) {
		needsUtf8 = valueIsTruthy(opt)
	}

	if !needsUtf8 && useUtf8 && highByteRe.MatchString(pattern) {
		return nil, throwNamed(ctx, "LogicException",
			fmt.Sprintf("Cannot use UTF-8 route patterns without setting the \"utf8\" option for route \"%s\".", propString(route, "path", pattern)))
	}
	if !useUtf8 && needsUtf8 {
		return nil, throwNamed(ctx, "LogicException",
			fmt.Sprintf("Cannot mix UTF-8 requirements with non-UTF-8 pattern \"%s\".", pattern))
	}

	matches := compileVarRe.FindAllStringSubmatchIndex(pattern, -1)
	for _, m := range matches {
		important := m[2] >= 0
		varName := pattern[m[4]:m[5]]
		precedingText := pattern[pos:m[0]]
		pos = m[1]

		precedingChar := ""
		if len(precedingText) > 0 {
			if useUtf8 {
				r, size := utf8.DecodeLastRuneInString(precedingText)
				if r != utf8.RuneError {
					precedingChar = string(r)
				} else {
					precedingChar = precedingText[len(precedingText)-size:]
				}
			} else {
				precedingChar = precedingText[len(precedingText)-1:]
			}
		}
		isSeparator := precedingChar != "" && strings.Contains(routeCompilerSeparators, precedingChar)

		if digitStartRe.MatchString(varName) {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Variable name \"%s\" cannot start with a digit in route pattern \"%s\". Please use a different name.", varName, pattern), "DomainException")
		}
		for _, v := range variables {
			if v == varName {
				return nil, throwNamed(ctx, "LogicException",
					fmt.Sprintf("Route pattern \"%s\" cannot reference variable name \"%s\" more than once.", pattern, varName))
			}
		}
		if len(varName) > variableMaxLength {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Variable name \"%s\" cannot be longer than %d characters in route pattern \"%s\". Please use a shorter name.", varName, variableMaxLength, pattern), "DomainException")
		}

		if isSeparator && precedingText != precedingChar {
			tokens = append(tokens, tokenBuf{kind: "text", text: precedingText[:len(precedingText)-len(precedingChar)]})
		} else if !isSeparator && precedingText != "" {
			tokens = append(tokens, tokenBuf{kind: "text", text: precedingText})
		}

		regexp := routeGetRequirement(route, varName)
		if regexp == "" {
			following := ""
			if pos < len(pattern) {
				following = pattern[pos:]
			}
			nextSep := findNextSeparator(following, useUtf8)
			regexp = "[^" + phpPregQuote(defaultSep)
			if defaultSep != nextSep && nextSep != "" {
				regexp += phpPregQuote(nextSep)
			}
			regexp += "]+"
			if (nextSep != "" && !nextVarStartRe.MatchString(following)) || following == "" {
				regexp += "+"
			}
		} else {
			if !utf8.ValidString(regexp) {
				useUtf8 = false
			} else if !needsUtf8 && utf8ReqRe.MatchString(regexp) {
				return nil, throwNamed(ctx, "LogicException",
					fmt.Sprintf("Cannot use UTF-8 route requirements without setting the \"utf8\" option for variable \"%s\" in pattern \"%s\".", varName, pattern))
			}
			if !useUtf8 && needsUtf8 {
				return nil, throwNamed(ctx, "LogicException",
					fmt.Sprintf("Cannot mix UTF-8 requirement with non-UTF-8 charset for variable \"%s\" in pattern \"%s\".", varName, pattern))
			}
			regexp = transformCapturingGroupsToNonCapturings(regexp)
		}

		sep := ""
		if isSeparator {
			sep = precedingChar
		}
		tokens = append(tokens, tokenBuf{
			kind:      "variable",
			text:      sep,
			regexp:    regexp,
			varName:   varName,
			important: important,
		})
		variables = append(variables, varName)
	}

	if pos < len(pattern) {
		tokens = append(tokens, tokenBuf{kind: "text", text: pattern[pos:]})
	}

	firstOptional := int(^uint(0) >> 1)
	if !isHost {
		for i := len(tokens) - 1; i >= 0; i-- {
			t := tokens[i]
			if t.kind == "variable" && !t.important && routeHasDefault(route, t.varName) {
				firstOptional = i
			} else {
				break
			}
		}
	}

	var rx strings.Builder
	for i := range tokens {
		rx.WriteString(computeRegexp(tokens, i, firstOptional))
	}
	regex := "{^" + rx.String() + "$}sD"
	if isHost {
		regex += "i"
	}
	if needsUtf8 {
		regex += "u"
		for i := range tokens {
			if tokens[i].kind == "variable" {
				tokens[i].utf8 = true
			}
		}
	}

	// reverse tokens
	rev := make([]data.Value, len(tokens))
	for i := range tokens {
		rev[len(tokens)-1-i] = tokens[i].toValue()
	}

	vars := make([]data.Value, len(variables))
	for i, v := range variables {
		vars[i] = data.NewStringValue(v)
	}

	return &compiledPattern{
		staticPrefix: determineStaticPrefix(route, tokens),
		regex:        regex,
		tokens:       phpList(rev...),
		variables:    phpList(vars...),
	}, nil
}

func determineStaticPrefix(route *data.ClassValue, tokens []tokenBuf) string {
	if len(tokens) == 0 {
		return ""
	}
	if tokens[0].kind != "text" {
		if routeHasDefault(route, tokens[0].varName) || tokens[0].text == "/" {
			return ""
		}
		return tokens[0].text
	}
	prefix := tokens[0].text
	if len(tokens) > 1 && tokens[1].text != "/" && !routeHasDefault(route, tokens[1].varName) {
		prefix += tokens[1].text
	}
	return prefix
}

func findNextSeparator(pattern string, useUtf8 bool) string {
	if pattern == "" {
		return ""
	}
	stripped := compileVarRe.ReplaceAllString(pattern, "")
	if stripped == "" {
		return ""
	}
	ch := stripped[:1]
	if useUtf8 {
		r, size := utf8.DecodeRuneInString(stripped)
		if r != utf8.RuneError {
			ch = stripped[:size]
		}
	}
	if strings.Contains(routeCompilerSeparators, ch) {
		return ch
	}
	return ""
}

func computeRegexp(tokens []tokenBuf, index, firstOptional int) string {
	token := tokens[index]
	if token.kind == "text" {
		return phpPregQuote(token.text)
	}
	if index == 0 && firstOptional == 0 {
		return phpPregQuote(token.text) + "(?P<" + token.varName + ">" + token.regexp + ")?"
	}
	regexp := phpPregQuote(token.text) + "(?P<" + token.varName + ">" + token.regexp + ")"
	if index >= firstOptional {
		regexp = "(?:" + regexp
		if len(tokens)-1 == index {
			n := len(tokens) - firstOptional
			if firstOptional == 0 {
				n--
			}
			if n < 0 {
				n = 0
			}
			regexp += strings.Repeat(")?", n)
		}
	}
	return regexp
}

func transformCapturingGroupsToNonCapturings(regexp string) string {
	b := []byte(regexp)
	for i := 0; i < len(b); i++ {
		if b[i] == '\\' {
			i++
			continue
		}
		if b[i] != '(' || i+2 >= len(b) {
			continue
		}
		i++
		if b[i] == '*' || b[i] == '?' {
			i++
			continue
		}
		b = append(b[:i], append([]byte("?:"), b[i:]...)...)
		i++
	}
	return string(b)
}

func extractInlineDefaultsAndRequirements(route *data.ClassValue, pattern string) string {
	if !strings.ContainsAny(pattern, "?<:") {
		return pattern
	}
	mapping, _ := assocGet(propArray(route, "defaults"), "_route_mapping")
	var mappingArr *data.ArrayValue
	if av, ok := mapping.(*data.ArrayValue); ok {
		mappingArr = cloneArray(av)
	} else {
		mappingArr = phpList()
	}
	changed := false
	out := inlineVarRe.ReplaceAllStringFunc(pattern, func(m string) string {
		sub := inlineVarRe.FindStringSubmatch(m)
		if sub == nil {
			return m
		}
		bang := sub[1]
		name := sub[2]
		mapped := sub[3]
		propPath := sub[4]
		req := sub[5]
		def := sub[6]
		if def != "" {
			if def == "?" {
				routeSetDefault(route, name, data.NewNullValue())
			} else {
				routeSetDefault(route, name, data.NewStringValue(def[1:]))
			}
		}
		if req != "" {
			routeSetRequirement(route, name, req[1:len(req)-1])
		}
		if mapped != "" {
			changed = true
			if propPath != "" {
				mappingArr = assocSet(mappingArr, name, phpList(data.NewStringValue(mapped), data.NewStringValue(propPath[1:])))
			} else {
				mappingArr = assocSet(mappingArr, name, data.NewStringValue(mapped))
			}
		}
		return "{" + bang + name + "}"
	})
	if changed && len(mappingArr.List) > 0 {
		routeSetDefault(route, "_route_mapping", mappingArr)
	}
	return out
}

func sanitizeRequirement(key, regex string) (string, error) {
	if regex != "" {
		if regex[0] == '^' {
			regex = regex[1:]
		} else if strings.HasPrefix(regex, `\A`) {
			regex = regex[2:]
		}
	}
	if strings.HasSuffix(regex, "$") {
		regex = regex[:len(regex)-1]
	} else if strings.HasSuffix(regex, `\z`) {
		regex = regex[:len(regex)-2]
	}
	if regex == "" {
		return "", fmt.Errorf("Routing requirement for \"%s\" cannot be empty.", key)
	}
	return regex, nil
}
