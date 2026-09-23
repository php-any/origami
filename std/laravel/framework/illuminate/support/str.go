package support

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const strClassName = "Illuminate\\Support\\Str"

type StrClass struct {
	node.Node
	methods map[string]data.Method
}

func NewStrClass() data.ClassStmt {
	c := &StrClass{methods: map[string]data.Method{}}
	c.register()
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
	add := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = kit.StaticMethod(name, params, -1, fn)
	}
	add("camel", []string{"value"}, strCamel)
	add("snake", []string{"value", "delimiter"}, strSnake)
	add("studly", []string{"value"}, strStudly)
	add("finish", []string{"value", "cap"}, strFinish)
	add("start", []string{"value", "prefix"}, strStart)
	add("is", []string{"pattern", "value"}, strIs)
	add("contains", []string{"haystack", "needles"}, strContains)
	add("startsWith", []string{"haystack", "needles"}, strStartsWith)
	add("endsWith", []string{"haystack", "needles"}, strEndsWith)
	add("lower", []string{"value"}, strLower)
	add("upper", []string{"value"}, strUpper)
	add("title", []string{"value"}, strTitle)
	add("uuid", nil, strUUID)
	add("random", []string{"length"}, strRandom)
	add("limit", []string{"value", "limit", "end"}, strLimit)
	add("replace", []string{"search", "replace", "subject"}, strReplace)
	add("substr", []string{"string", "start", "length"}, strSubstr)
	add("length", []string{"value"}, strLength)
	add("slug", []string{"title", "separator", "language", "dictionary"}, strSlug)
	add("of", []string{"string"}, strOf)
	add("after", []string{"subject", "search"}, strAfter)
	add("before", []string{"subject", "search"}, strBefore)
	add("kebab", []string{"value"}, strKebab)
	c.methods["explode"] = kit.StaticMethod("explode", []string{"value", "delimiter", "limit"}, 2, strExplode)
	c.registerMore()
	kit.RegisterMacroable(c.methods, strClassName)
}

func strArg(ctx data.Context, i int) string {
	v, _ := ctx.GetIndexValue(i)
	if v == nil {
		return ""
	}
	return v.AsString()
}

func strCamel(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	strCacheMu.RLock()
	if out, ok := strCamelCache[value]; ok {
		strCacheMu.RUnlock()
		return data.NewStringValue(out), nil
	}
	strCacheMu.RUnlock()
	s := studlyCached(value, false)
	if s == "" {
		return data.NewStringValue(""), nil
	}
	out := mbLcfirst(s)
	strCacheMu.Lock()
	strCamelCache[value] = out
	strCacheMu.Unlock()
	return data.NewStringValue(out), nil
}

func strStudly(ctx data.Context) (data.GetValue, data.Control) {
	norm := strBoolArg(ctx, 1, false)
	return data.NewStringValue(studlyCached(strArg(ctx, 0), norm)), nil
}

func strSnake(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	delim := "_"
	if d, ok := ctx.GetIndexValue(1); ok && d != nil && !kit.IsNull(d) {
		delim = d.AsString()
	}
	return data.NewStringValue(snakeCached(s, delim)), nil
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
	value := strArg(ctx, 1)
	ignoreCase := strBoolArg(ctx, 2, false)
	for _, pattern := range strNeedles(kit.Arg(ctx, 0)) {
		if pattern == "*" || pattern == value {
			return data.NewBoolValue(true), nil
		}
		if ignoreCase && strings.EqualFold(pattern, value) {
			return data.NewBoolValue(true), nil
		}
		quoted := regexp.QuoteMeta(pattern)
		quoted = strings.ReplaceAll(quoted, `\*`, ".*")
		flags := ""
		if ignoreCase {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + "^" + quoted + "$")
		if err == nil && re.MatchString(value) {
			return data.NewBoolValue(true), nil
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
		for _, e := range kit.Entries(av) {
			out = append(out, e.Value.AsString())
		}
		return out
	}
	return []string{v.AsString()}
}

func strContains(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	if haystack == "" {
		return data.NewBoolValue(false), nil
	}
	needles, _ := ctx.GetIndexValue(1)
	ignore := strBoolArg(ctx, 2, false)
	for _, n := range strNeedles(needles) {
		h, nd := haystack, n
		if ignore {
			h, nd = strings.ToLower(h), strings.ToLower(nd)
		}
		if nd != "" && strings.Contains(h, nd) {
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

func strUUID(ctx data.Context) (data.GetValue, data.Control) {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return data.NewStringValue(s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]), nil
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
	if strRandomFactory != nil {
		ret, ctl := kit.Call(ctx, strRandomFactory, data.NewIntValue(length))
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil {
			return ret, nil
		}
	}
	return data.NewStringValue(randomStringLaravel(length)), nil
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
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !kit.IsNull(v) {
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
	if strBoolArg(ctx, 3, true) {
		return data.NewStringValue(strings.ReplaceAll(subject, search, replace)), nil
	}
	return data.NewStringValue(strReplaceFold(subject, search, replace)), nil
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
	// 默认取到末尾
	end := len(s)
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !kit.IsNull(v) {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				if n < 0 {
					// 负 length：从末尾回退 |n| 个字符（对齐 PHP substr）
					end = len(s) + n
				} else {
					end = start + n
				}
			}
		}
	}
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

func strSlug(ctx data.Context) (data.GetValue, data.Control) {
	title := strArg(ctx, 0)
	lang := strArg(ctx, 2)
	if lang != "" {
		title = asciiFallback(title, lang)
	}
	title = strings.ToLower(strings.TrimSpace(title))
	sep := "-"
	if s := strArg(ctx, 1); s != "" {
		sep = s
	}
	flip := "_"
	if sep == "_" {
		flip = "-"
	}
	title = strings.ReplaceAll(title, flip, sep)
	var b strings.Builder
	lastSep := false
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastSep = false
			continue
		}
		if !lastSep && b.Len() > 0 {
			b.WriteString(sep)
			lastSep = true
		}
	}
	out := strings.Trim(b.String(), sep)
	return data.NewStringValue(out), nil
}

func strOf(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, strArg(ctx, 0))
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

func strKebab(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(snakeCached(strArg(ctx, 0), "-")), nil
}

func strExplode(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	delimiter := strArg(ctx, 1)
	limit := int(^uint(0) >> 1) // PHP_INT_MAX-ish
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		if iv, ok := v.(*data.IntValue); ok {
			limit = iv.Value
		} else if b, ok := v.(data.AsInt); ok {
			if n, err := b.AsInt(); err == nil {
				limit = n
			}
		}
	}
	parts := explodeLimited(value, delimiter, limit)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for i, p := range parts {
		out.SetIntKey(i, data.NewStringValue(p))
	}
	return out, nil
}

func explodeLimited(value, delimiter string, limit int) []string {
	if delimiter == "" {
		return []string{value}
	}
	if limit == 0 {
		return []string{}
	}
	if limit < 0 {
		// PHP negative limit: return all but last |limit| elements
		all := strings.Split(value, delimiter)
		n := len(all) + limit
		if n < 0 {
			n = 0
		}
		if n > len(all) {
			n = len(all)
		}
		return all[:n]
	}
	if limit == 1 {
		return []string{value}
	}
	parts := strings.SplitN(value, delimiter, limit)
	return parts
}
