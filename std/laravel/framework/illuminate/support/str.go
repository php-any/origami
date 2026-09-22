package support

import (
	"crypto/rand"
	"encoding/hex"
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

func strSnake(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	delim := "_"
	if d, ok := ctx.GetIndexValue(1); ok && d != nil && !kit.IsNull(d) {
		delim = d.AsString()
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
	out = strings.ReplaceAll(out, "-", delim)
	return data.NewStringValue(strings.ToLower(out)), nil
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
		for _, e := range kit.Entries(av) {
			out = append(out, e.Value.AsString())
		}
		return out
	}
	return []string{v.AsString()}
}

func strContains(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needles, _ := ctx.GetIndexValue(1)
	for _, n := range strNeedles(needles) {
		if n != "" && strings.Contains(haystack, n) {
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
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !kit.IsNull(v) {
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

func strSlug(ctx data.Context) (data.GetValue, data.Control) {
	title := strings.ToLower(strings.TrimSpace(strArg(ctx, 0)))
	sep := "-"
	if s := strArg(ctx, 1); s != "" {
		sep = s
	}
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
	// Stringable 可后做；先返回原始字符串
	return data.NewStringValue(strArg(ctx, 0)), nil
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
	s := strArg(ctx, 0)
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(unicode.ToLower(r))
		} else if r == '_' || r == ' ' {
			b.WriteByte('-')
		} else {
			b.WriteRune(r)
		}
	}
	return data.NewStringValue(strings.ToLower(b.String())), nil
}
