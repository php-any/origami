package support

import (
	"crypto/rand"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const strClassName = "Illuminate\\Support\\Str"

type StrClass struct {
	node.Node
	methods map[string]data.Method
}

func NewStrClass() data.ClassStmt {
	c := &StrClass{methods: map[string]data.Method{}}
	c.register()
	registerMacroable(c.methods, strClassName)
	return c
}

func (c *StrClass) GetName() string                          { return strClassName }
func (c *StrClass) GetExtend() *string                       { return nil }
func (c *StrClass) GetImplements() []string                  { return nil }
func (c *StrClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *StrClass) GetPropertyList() []data.Property         { return nil }
func (c *StrClass) GetConstruct() data.Method                { return nil }
func (c *StrClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *StrClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *StrClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *StrClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *StrClass) register() {
	add := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = newStaticMethod(name, params, optionalFrom, fn, false)
	}
	add("camel", []string{"value"}, -1, strCamel)
	add("snake", []string{"value", "delimiter"}, 1, strSnake)
	add("kebab", []string{"value"}, -1, strKebab)
	add("studly", []string{"value"}, -1, strStudly)
	add("finish", []string{"value", "cap"}, -1, strFinish)
	add("start", []string{"value", "prefix"}, -1, strStart)
	add("startsWith", []string{"haystack", "needles"}, -1, strStartsWith)
	add("endsWith", []string{"haystack", "needles"}, -1, strEndsWith)
	add("lower", []string{"value"}, -1, strLower)
	add("upper", []string{"value"}, -1, strUpper)
	add("contains", []string{"haystack", "needles", "ignoreCase"}, 2, strContains)
	add("containsAll", []string{"haystack", "needles", "ignoreCase"}, 2, strContainsAll)
	add("after", []string{"subject", "search"}, -1, strAfter)
	add("afterLast", []string{"subject", "search"}, -1, strAfterLast)
	add("before", []string{"subject", "search"}, -1, strBefore)
	add("beforeLast", []string{"subject", "search"}, -1, strBeforeLast)
	add("ascii", []string{"value", "language"}, 1, strAscii)
	add("slug", []string{"title", "separator", "language", "dictionary"}, 1, strSlug)
	add("of", []string{"string"}, 1, strOf)
	add("is", []string{"pattern", "value", "ignoreCase"}, 2, strIs)
	add("title", []string{"value"}, -1, strTitle)
	add("trim", []string{"value", "charlist"}, 1, strTrim)
	add("ltrim", []string{"value", "charlist"}, 1, strLtrim)
	add("rtrim", []string{"value", "charlist"}, 1, strRtrim)
	add("replace", []string{"search", "replace", "subject", "caseSensitive"}, 3, strReplace)
	add("replaceFirst", []string{"search", "replace", "subject"}, -1, strReplaceFirst)
	add("replaceLast", []string{"search", "replace", "subject"}, -1, strReplaceLast)
	add("replaceStart", []string{"search", "replace", "subject"}, -1, strReplaceStart)
	add("replaceEnd", []string{"search", "replace", "subject"}, -1, strReplaceEnd)
	add("uuid", nil, -1, strUUID)
	add("uuid7", []string{"time"}, 1, strUUID7)
	add("orderedUuid", nil, -1, strOrderedUUID)
	add("createUuidsUsing", []string{"factory"}, 1, strCreateUuidsUsing)
	add("createUuidsNormally", nil, -1, strCreateUuidsNormally)
	add("isUuid", []string{"value", "version"}, 1, strIsUuid)
	add("ulid", []string{"time"}, 1, strUlid)
	add("random", []string{"length"}, 1, strRandom)
	add("headline", []string{"value"}, -1, strHeadline)
	add("ucfirst", []string{"string"}, -1, strUcfirst)
	add("lcfirst", []string{"string"}, -1, strLcfirst)
	add("repeat", []string{"string", "times"}, -1, strRepeat)
	add("reverse", []string{"value"}, -1, strReverse)
	add("wrap", []string{"value", "before", "after"}, 2, strWrap)
	add("isAscii", []string{"value"}, -1, strIsAscii)
	add("isJson", []string{"value"}, -1, strIsJson)
	add("pascal", []string{"value", "normalize"}, 1, strStudly)
	add("flushCache", nil, -1, strFlushCache)
	add("limit", []string{"value", "limit", "end"}, 1, strLimit)
	add("substr", []string{"string", "start", "length"}, 2, strSubstr)
	add("length", []string{"value"}, -1, strLength)
	add("parseCallback", []string{"callback", "default"}, 1, strParseCallback)
	add("substrCount", []string{"haystack", "needle", "offset", "length"}, 2, strSubstrCount)
	add("toBase64", []string{"value"}, -1, strToBase64)
	add("fromBase64", []string{"value", "strict"}, 1, strFromBase64)
	add("between", []string{"subject", "from", "to"}, -1, strBetween)
	add("betweenFirst", []string{"subject", "from", "to"}, -1, strBetweenFirst)
	add("position", []string{"haystack", "needle", "offset"}, 2, strPosition)
	add("remove", []string{"search", "subject", "caseSensitive"}, 2, strRemove)
	add("replaceArray", []string{"search", "replace", "subject"}, -1, strReplaceArray)
	add("squish", []string{"value"}, -1, strSquish)
	add("padLeft", []string{"value", "length", "pad"}, 2, strPadLeft)
	add("padRight", []string{"value", "length", "pad"}, 2, strPadRight)
	add("padBoth", []string{"value", "length", "pad"}, 2, strPadBoth)
	add("take", []string{"value", "limit"}, -1, strTake)
	add("unwrap", []string{"value", "before", "after"}, 2, strUnwrap)
	add("doesntContain", []string{"haystack", "needles", "ignoreCase"}, 2, strDoesntContain)
	add("doesntStartWith", []string{"haystack", "needles"}, -1, strDoesntStartWith)
	add("doesntEndWith", []string{"haystack", "needles"}, -1, strDoesntEndWith)
	add("isUrl", []string{"value", "protocols"}, 1, strIsUrl)
	add("isUlid", []string{"value"}, -1, strIsUlid)
	add("charAt", []string{"subject", "index"}, -1, strCharAt)
	add("numbers", []string{"value"}, -1, strNumbers)
	add("ucwords", []string{"value"}, -1, strUcwords)
	add("wordCount", []string{"value"}, -1, strWordCount)
	add("freezeUuids", []string{"callback"}, 1, strFreezeUuids)
	add("createUuidsUsingSequence", []string{"sequence", "whenMissing"}, 1, strCreateUuidsUsingSequence)
	add("plural", []string{"value", "count", "prependCount"}, 1, strPlural)
	add("singular", []string{"value"}, -1, strSingular)
	add("pluralStudly", []string{"value", "count"}, 1, strPluralStudly)
	add("pluralPascal", []string{"value", "count"}, 1, strPluralStudly)
}

func strArg(ctx data.Context, i int) string {
	v, _ := ctx.GetIndexValue(i)
	if v == nil {
		return ""
	}
	return v.AsString()
}

func strCamel(ctx data.Context) (data.GetValue, data.Control) {
	s := studly(strArg(ctx, 0))
	if s == "" {
		return data.NewStringValue(""), nil
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return data.NewStringValue(string(r)), nil
}

func strStudly(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(studly(strArg(ctx, 0))), nil
}

func studly(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	parts := strings.Fields(s)
	for i, p := range parts {
		if p == "" {
			continue
		}
		r := []rune(p)
		r[0] = unicode.ToUpper(r[0])
		parts[i] = string(r)
	}
	return strings.Join(parts, "")
}

func strKebab(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(kebabString(strArg(ctx, 0))), nil
}

func kebabString(s string) string {
	return snakeString(s, "-")
}

func strAfter(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	search := strArg(ctx, 1)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	i := strings.Index(subject, search)
	if i < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(subject[i+len(search):]), nil
}

func strBefore(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	search := strArg(ctx, 1)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	i := strings.Index(subject, search)
	if i < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(subject[:i]), nil
}

func snakeString(s, delim string) string {
	if delim == "" {
		delim = "_"
	}
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteString(delim)
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	out := strings.ReplaceAll(b.String(), " ", delim)
	out = strings.ReplaceAll(out, "_", delim)
	if delim != "-" {
		out = strings.ReplaceAll(out, "-", delim)
	}
	return strings.ToLower(out)
}

func strSnake(ctx data.Context) (data.GetValue, data.Control) {
	delim := "_"
	if d, ok := ctx.GetIndexValue(1); ok && d != nil && !isNull(d) {
		delim = d.AsString()
	}
	return data.NewStringValue(snakeString(strArg(ctx, 0), delim)), nil
}

func strFinish(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	cap := strArg(ctx, 1)
	if cap == "" {
		return data.NewStringValue(value), nil
	}
	for strings.HasSuffix(value, cap) {
		value = strings.TrimSuffix(value, cap)
	}
	return data.NewStringValue(value + cap), nil
}

func strStart(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	prefix := strArg(ctx, 1)
	if prefix != "" && !strings.HasPrefix(value, prefix) {
		value = prefix + value
	}
	return data.NewStringValue(value), nil
}

func strIs(ctx data.Context) (data.GetValue, data.Control) {
	pattern := strArg(ctx, 0)
	value := strArg(ctx, 1)
	if pattern == value {
		return data.NewBoolValue(true), nil
	}
	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			return data.NewBoolValue(strings.HasPrefix(value, parts[0]) && strings.HasSuffix(value, parts[1])), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func strNeedles(v data.Value) []string {
	if v == nil {
		return nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		out := make([]string, 0)
		for _, e := range toEntries(av) {
			out = append(out, e.value.AsString())
		}
		return out
	}
	return []string{v.AsString()}
}

func strContains(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needles, _ := ctx.GetIndexValue(1)
	ignore := false
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			ignore, _ = b.AsBool()
		}
	}
	if ignore {
		haystack = strings.ToLower(haystack)
	}
	for _, n := range strNeedles(needles) {
		if n == "" {
			continue
		}
		if ignore {
			n = strings.ToLower(n)
		}
		if strings.Contains(haystack, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func strStartsWith(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needles, _ := ctx.GetIndexValue(1)
	for _, n := range strNeedles(needles) {
		if n != "" && strings.HasPrefix(haystack, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func strEndsWith(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needles, _ := ctx.GetIndexValue(1)
	for _, n := range strNeedles(needles) {
		if n != "" && strings.HasSuffix(haystack, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func strLower(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(strings.ToLower(strArg(ctx, 0))), nil
}

func strUpper(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(strings.ToUpper(strArg(ctx, 0))), nil
}

func strTitle(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(strings.Title(strings.ToLower(strArg(ctx, 0)))), nil
}

func strRandom(ctx data.Context) (data.GetValue, data.Control) {
	length := 16
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil && n > 0 {
				length = n
			}
		}
	}
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rb := make([]byte, length)
	_, _ = rand.Read(rb)
	for i := 0; i < length; i++ {
		b[i] = alphabet[int(rb[i])%len(alphabet)]
	}
	return data.NewStringValue(string(b)), nil
}

func strLimit(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	limit := 100
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				limit = n
			}
		}
	}
	end := "..."
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		end = v.AsString()
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return data.NewStringValue(value), nil
	}
	return data.NewStringValue(string(runes[:limit]) + end), nil
}

func strReplace(ctx data.Context) (data.GetValue, data.Control) {
	search := strArg(ctx, 0)
	replace := strArg(ctx, 1)
	subject := strArg(ctx, 2)
	return data.NewStringValue(strings.ReplaceAll(subject, search, replace)), nil
}

func strSubstr(ctx data.Context) (data.GetValue, data.Control) {
	s := []rune(strArg(ctx, 0))
	start := 0
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				start = n
			}
		}
	}
	if start < 0 {
		start = len(s) + start
	}
	if start < 0 {
		start = 0
	}
	if start > len(s) {
		return data.NewStringValue(""), nil
	}
	length := len(s) - start
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				length = n
			}
		}
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	if end < start {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(s[start:end])), nil
}

func strLength(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(len([]rune(strArg(ctx, 0)))), nil
}
