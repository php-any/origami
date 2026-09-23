package support

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

var strUUIDRe = regexp.MustCompile(`^[\da-fA-F]{8}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{12}$`)

func (c *StrClass) registerMore() {
	add := func(name string, params []string, opt int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = kit.StaticMethod(name, params, opt, fn)
	}
	add("trans", []string{"key", "replace", "locale"}, 1, strTrans)
	add("afterlast", []string{"subject", "search"}, -1, strAfterLast)
	add("ascii", []string{"value", "language"}, 1, strAscii)
	add("transliterate", []string{"string", "unknown", "strict"}, 1, strTransliterate)
	add("beforelast", []string{"subject", "search"}, -1, strBeforeLast)
	add("between", []string{"subject", "from", "to"}, -1, strBetween)
	add("betweenfirst", []string{"subject", "from", "to"}, -1, strBetweenFirst)
	add("charat", []string{"subject", "index"}, -1, strCharAt)
	add("chopstart", []string{"subject", "needle"}, -1, strChopStart)
	add("chopend", []string{"subject", "needle"}, -1, strChopEnd)
	add("containsall", []string{"haystack", "needles", "ignoreCase"}, 2, strContainsAll)
	add("doesntcontain", []string{"haystack", "needles", "ignoreCase"}, 2, strDoesntContain)
	add("convertcase", []string{"string", "mode", "encoding"}, 1, strConvertCase)
	add("counted", []string{"value", "count"}, -1, strCounted)
	add("deduplicate", []string{"string", "characters"}, 1, strDeduplicate)
	add("doesntendwith", []string{"haystack", "needles"}, -1, strDoesntEndWith)
	add("excerpt", []string{"text", "phrase", "options"}, 1, strExcerpt)
	add("wrap", []string{"value", "before", "after"}, 2, strWrap)
	add("unwrap", []string{"value", "before", "after"}, 2, strUnwrap)
	add("isascii", []string{"value"}, -1, strIsAscii)
	add("isjson", []string{"value"}, -1, strIsJson)
	add("isurl", []string{"value", "protocols"}, 1, strIsUrl)
	add("isuuid", []string{"value", "version"}, 1, strIsUuid)
	add("isulid", []string{"value"}, -1, strIsUlid)
	add("words", []string{"value", "words", "end"}, 1, strWords)
	add("markdown", []string{"string", "options", "extensions"}, 1, strMarkdown)
	add("inlinemarkdown", []string{"string", "options", "extensions"}, 1, strInlineMarkdown)
	add("mask", []string{"string", "character", "index", "length", "encoding"}, 3, strMask)
	add("match", []string{"pattern", "subject"}, -1, strMatch)
	add("ismatch", []string{"pattern", "value"}, -1, strIsMatch)
	add("matchall", []string{"pattern", "subject"}, -1, strMatchAll)
	add("numbers", []string{"value"}, -1, strNumbers)
	add("padboth", []string{"value", "length", "pad"}, 2, strPadBoth)
	add("padleft", []string{"value", "length", "pad"}, 2, strPadLeft)
	add("padright", []string{"value", "length", "pad"}, 2, strPadRight)
	add("parsecallback", []string{"callback", "default"}, 1, strParseCallback)
	add("plural", []string{"value", "count", "prependCount"}, 1, strPlural)
	add("pluralstudly", []string{"value", "count"}, 1, strPluralStudly)
	add("pluralpascal", []string{"value", "count"}, 1, strPluralPascal)
	add("password", []string{"length", "letters", "numbers", "symbols", "spaces"}, 0, strPassword)
	add("position", []string{"haystack", "needle", "offset", "encoding"}, 2, strPosition)
	add("createrandomstringsusing", []string{"factory"}, -1, strCreateRandomStringsUsing)
	add("createrandomstringsusingsequence", []string{"sequence", "whenMissing"}, 1, strCreateRandomStringsUsingSequence)
	add("createrandomstringsnormally", nil, -1, strCreateRandomStringsNormally)
	add("repeat", []string{"string", "times"}, -1, strRepeat)
	add("replacearray", []string{"search", "replace", "subject"}, -1, strReplaceArray)
	add("replacefirst", []string{"search", "replace", "subject"}, -1, strReplaceFirst)
	add("replacestart", []string{"search", "replace", "subject"}, -1, strReplaceStart)
	add("replacelast", []string{"search", "replace", "subject"}, -1, strReplaceLast)
	add("replaceend", []string{"search", "replace", "subject"}, -1, strReplaceEnd)
	add("replacematches", []string{"pattern", "replace", "subject", "limit"}, 3, strReplaceMatches)
	add("remove", []string{"search", "subject", "caseSensitive"}, 2, strRemove)
	add("reverse", []string{"value"}, -1, strReverse)
	add("headline", []string{"value"}, -1, strHeadline)
	add("initials", []string{"value", "capitalize"}, 1, strInitials)
	add("apa", []string{"value"}, -1, strApa)
	add("singular", []string{"value"}, -1, strSingular)
	add("trim", []string{"value", "charlist"}, 1, strTrim)
	add("ltrim", []string{"value", "charlist"}, 1, strLtrim)
	add("rtrim", []string{"value", "charlist"}, 1, strRtrim)
	add("squish", []string{"value"}, -1, strSquish)
	add("doesntstartwith", []string{"haystack", "needles"}, -1, strDoesntStartWith)
	add("pascal", []string{"value", "normalize"}, 1, strPascal)
	add("substrcount", []string{"haystack", "needle", "offset", "length"}, 2, strSubstrCount)
	add("substrreplace", []string{"string", "replace", "offset", "length"}, 2, strSubstrReplace)
	add("swap", []string{"map", "subject"}, -1, strSwap)
	add("take", []string{"string", "limit"}, -1, strTake)
	add("tobase64", []string{"string"}, -1, strToBase64)
	add("frombase64", []string{"string", "strict"}, 1, strFromBase64)
	add("lcfirst", []string{"string"}, -1, strLcfirst)
	add("ucfirst", []string{"string"}, -1, strUcfirst)
	add("ucwords", []string{"string", "separators"}, 1, strUcwords)
	add("ucsplit", []string{"string"}, -1, strUcsplit)
	add("wordcount", []string{"string", "characters"}, 1, strWordCount)
	add("wordwrap", []string{"string", "characters", "break", "cutLongWords"}, 1, strWordWrap)
	add("uuid7", []string{"time"}, -1, strUUID7)
	add("ordereduuid", nil, -1, strOrderedUUID)
	add("createuuidsusing", []string{"factory"}, -1, strCreateUuidsUsing)
	add("createuuidsusingsequence", []string{"sequence", "whenMissing"}, 1, strCreateUuidsUsingSequence)
	add("freezeuuids", []string{"callback"}, -1, strFreezeUuids)
	add("createuuidsnormally", nil, -1, strCreateUuidsNormally)
	add("ulid", []string{"time"}, -1, strULID)
	add("createulidsnormally", nil, -1, strCreateUlidsNormally)
	add("createulidsusing", []string{"factory"}, -1, strCreateUlidsUsing)
	add("createulidsusingsequence", []string{"sequence", "whenMissing"}, 1, strCreateUlidsUsingSequence)
	add("freezeulids", []string{"callback"}, -1, strFreezeUlids)
	add("flushcache", nil, -1, strFlushCache)
	add("resetfactorystate", nil, -1, strResetFactoryState)
}

func strTrans(ctx data.Context) (data.GetValue, data.Control) {
	key := strArg(ctx, 0)
	// __() 未接翻译器；键名原样返回 Stringable
	return newStringableValue(ctx, key)
}

func strAfterLast(ctx data.Context) (data.GetValue, data.Control) {
	subject, search := strArg(ctx, 0), strArg(ctx, 1)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	pos := mbStrrpos(subject, search)
	if pos < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(mbSubstr(subject, pos+mbLen(search), mbLen(subject))), nil
}

func strAscii(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(asciiFallback(strArg(ctx, 0), strArg(ctx, 1))), nil
}

func strTransliterate(ctx data.Context) (data.GetValue, data.Control) {
	unknown := "?"
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		unknown = v.AsString()
	}
	strict := strBoolArg(ctx, 2, false)
	return data.NewStringValue(transliterateFallback(strArg(ctx, 0), unknown, strict)), nil
}

func strBeforeLast(ctx data.Context) (data.GetValue, data.Control) {
	subject, search := strArg(ctx, 0), strArg(ctx, 1)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	pos := mbStrrpos(subject, search)
	if pos < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(mbSubstr(subject, 0, pos)), nil
}

func strBetween(ctx data.Context) (data.GetValue, data.Control) {
	subject, from, to := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if from == "" || to == "" {
		return data.NewStringValue(subject), nil
	}
	rest := strAfterImpl(subject, from)
	return data.NewStringValue(strBeforeLastImpl(rest, to)), nil
}

func strBetweenFirstImpl(subject, from, to string) string {
	return strBeforeImpl(strAfterImpl(subject, from), to)
}

func strAfterImpl(subject, search string) string {
	if search == "" {
		return subject
	}
	parts := strings.SplitN(subject, search, 2)
	if len(parts) < 2 {
		return subject
	}
	return parts[1]
}

func strBeforeImpl(subject, search string) string {
	if search == "" {
		return subject
	}
	i := strings.Index(subject, search)
	if i < 0 {
		return subject
	}
	return subject[:i]
}

func strBeforeLastImpl(subject, search string) string {
	if search == "" {
		return subject
	}
	pos := mbStrrpos(subject, search)
	if pos < 0 {
		return subject
	}
	return mbSubstr(subject, 0, pos)
}

func strBetweenFirst(ctx data.Context) (data.GetValue, data.Control) {
	subject, from, to := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if from == "" || to == "" {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(strBetweenFirstImpl(subject, from, to)), nil
}

func strCharAt(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	idx := strIntArg(ctx, 1, 0)
	length := mbLen(subject)
	if idx < 0 {
		if idx < -length {
			return data.NewBoolValue(false), nil
		}
	} else if idx > length-1 {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(mbSubstr(subject, idx, 1)), nil
}

func strChopStart(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	for _, n := range strNeedles(kit.Arg(ctx, 1)) {
		if n != "" && strings.HasPrefix(subject, n) {
			return data.NewStringValue(mbSubstr(subject, mbLen(n), mbLen(subject))), nil
		}
	}
	return data.NewStringValue(subject), nil
}

func strChopEnd(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	for _, n := range strNeedles(kit.Arg(ctx, 1)) {
		if n != "" && strings.HasSuffix(subject, n) {
			return data.NewStringValue(mbSubstr(subject, 0, mbLen(subject)-mbLen(n))), nil
		}
	}
	return data.NewStringValue(subject), nil
}

func strContainsAll(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needles := kit.Arg(ctx, 1)
	ignore := strBoolArg(ctx, 2, false)
	any := false
	for _, n := range strNeedles(needles) {
		any = true
		h, nd := haystack, n
		if ignore {
			h, nd = strings.ToLower(h), strings.ToLower(nd)
		}
		if nd == "" || !strings.Contains(h, nd) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(any), nil
}

func strDoesntContain(ctx data.Context) (data.GetValue, data.Control) {
	b, ctl := strContains(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := b.(data.AsBool).AsBool()
	return data.NewBoolValue(!v), nil
}

func strConvertCase(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	mode := strIntArg(ctx, 1, 3) // MB_CASE_FOLD default in PHP 8.3+
	switch mode {
	case 0, 4:
		return data.NewStringValue(strings.ToUpper(s)), nil
	case 1, 5:
		return data.NewStringValue(strings.ToLower(s)), nil
	case 2, 6:
		return data.NewStringValue(strings.Title(strings.ToLower(s))), nil
	default:
		return data.NewStringValue(strings.ToLower(s)), nil
	}
}

func strCounted(ctx data.Context) (data.GetValue, data.Control) {
	word := strArg(ctx, 0)
	count := strIntArg(ctx, 1, 2)
	out := englishPluralFallback(word, count)
	return data.NewStringValue(fmt.Sprintf("%s %s", strconv.Itoa(count), out)), nil
}

func strDeduplicate(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	ch := kit.Arg(ctx, 1)
	chars := []string{" "}
	if ch != nil && !kit.IsNull(ch) {
		if av, ok := ch.(*data.ArrayValue); ok {
			for _, e := range kit.Entries(av) {
				chars = append(chars, e.Value.AsString())
			}
		} else {
			chars = []string{ch.AsString()}
		}
	}
	out := s
	for _, c := range chars {
		if c == "" {
			continue
		}
		re := regexp.MustCompile(regexp.QuoteMeta(c) + "+")
		out = re.ReplaceAllString(out, c)
	}
	return data.NewStringValue(out), nil
}

func strDoesntEndWith(ctx data.Context) (data.GetValue, data.Control) {
	b, ctl := strEndsWith(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := b.(data.AsBool).AsBool()
	return data.NewBoolValue(!v), nil
}

func strExcerpt(ctx data.Context) (data.GetValue, data.Control) {
	text, phrase := strArg(ctx, 0), strArg(ctx, 1)
	if phrase == "" {
		return data.NewNullValue(), nil
	}
	re := regexp.MustCompile("(?is)^(.*?)" + regexp.QuoteMeta(phrase) + "(.*)$")
	m := re.FindStringSubmatch(text)
	if len(m) < 3 {
		return data.NewNullValue(), nil
	}
	radius := 100
	if opt := kit.Arg(ctx, 2); opt != nil {
		if av, ok := opt.(*data.ArrayValue); ok {
			for _, e := range kit.Entries(av) {
				if e.KeyStr == "radius" {
					radius = strIntArgFromValue(e.Value, 100)
				}
			}
		}
	}
	omission := "..."
	start := m[1]
	if len([]rune(start)) > radius {
		start = omission + mbSubstr(start, len([]rune(start))-radius, radius)
	}
	end := m[2]
	if len([]rune(end)) > radius {
		end = mbSubstr(end, 0, radius) + omission
	}
	return data.NewStringValue(strings.TrimSpace(start) + phrase + strings.TrimSpace(end)), nil
}

func strIntArgFromValue(v data.Value, def int) int {
	if v == nil {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func strWrap(ctx data.Context) (data.GetValue, data.Control) {
	value, before := strArg(ctx, 0), strArg(ctx, 1)
	after := before
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		after = v.AsString()
	}
	return data.NewStringValue(before + value + after), nil
}

func strUnwrap(ctx data.Context) (data.GetValue, data.Control) {
	value, before := strArg(ctx, 0), strArg(ctx, 1)
	after := before
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		after = v.AsString()
	}
	if strings.HasPrefix(value, before) {
		value = mbSubstr(value, mbLen(before), mbLen(value))
	}
	if strings.HasSuffix(value, after) {
		value = mbSubstr(value, 0, mbLen(value)-mbLen(after))
	}
	return data.NewStringValue(value), nil
}

func strIsAscii(ctx data.Context) (data.GetValue, data.Control) {
	for _, r := range strArg(ctx, 0) {
		if r > unicode.MaxASCII {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func strIsJson(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	var js any
	return data.NewBoolValue(json.Unmarshal([]byte(s), &js) == nil), nil
}

func strIsUrl(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
			u, err = url.Parse(s)
			return data.NewBoolValue(err == nil && u.Host != ""), nil
		}
		return data.NewBoolValue(false), nil
	}
	protos := kit.Arg(ctx, 1)
	if av, ok := protos.(*data.ArrayValue); ok && len(kit.Entries(av)) > 0 {
		for _, e := range kit.Entries(av) {
			if strings.EqualFold(u.Scheme, e.Value.AsString()) {
				return data.NewBoolValue(true), nil
			}
		}
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func strIsUuid(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if !strUUIDRe.MatchString(s) {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func strIsUlid(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(ulidValid(strArg(ctx, 0))), nil
}

func strWords(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	maxWords := strIntArg(ctx, 1, 100)
	end := "..."
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		end = v.AsString()
	}
	parts := strings.Fields(value)
	if len(parts) <= maxWords {
		return data.NewStringValue(value), nil
	}
	return data.NewStringValue(strings.Join(parts[:maxWords], " ") + end), nil
}

func strMarkdown(ctx data.Context) (data.GetValue, data.Control) {
	// League CommonMark 未接入；返回转义后的纯文本段落近似。
	s := strArg(ctx, 0)
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return data.NewStringValue("<p>" + strings.ReplaceAll(s, "\n", "<br>\n") + "</p>"), nil
}

func strInlineMarkdown(ctx data.Context) (data.GetValue, data.Control) {
	return strMarkdown(ctx)
}

func strMask(ctx data.Context) (data.GetValue, data.Control) {
	s, char := strArg(ctx, 0), strArg(ctx, 1)
	if char == "" {
		return data.NewStringValue(s), nil
	}
	idx := strIntArg(ctx, 2, 0)
	length := strIntArg(ctx, 3, mbLen(mbSubstr(s, idx, mbLen(s))))
	if idx < 0 {
		idx = mbLen(s) + idx
		if idx < 0 {
			idx = 0
		}
	}
	seg := mbSubstr(s, idx, length)
	if seg == "" {
		return data.NewStringValue(s), nil
	}
	mask := strings.Repeat(string([]rune(char)[0]), mbLen(seg))
	return data.NewStringValue(mbSubstr(s, 0, idx) + mask + mbSubstr(s, idx+mbLen(seg), mbLen(s))), nil
}

func strMatch(ctx data.Context) (data.GetValue, data.Control) {
	pattern, subject := strArg(ctx, 0), strArg(ctx, 1)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	m := re.FindStringSubmatch(subject)
	if len(m) == 0 {
		return data.NewStringValue(""), nil
	}
	if len(m) > 1 {
		return data.NewStringValue(m[1]), nil
	}
	return data.NewStringValue(m[0]), nil
}

func strIsMatch(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 1)
	for _, p := range strNeedles(kit.Arg(ctx, 0)) {
		re, err := regexp.Compile(p)
		if err == nil && re.MatchString(value) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func strMatchAll(ctx data.Context) (data.GetValue, data.Control) {
	pattern, subject := strArg(ctx, 0), strArg(ctx, 1)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return strEmptyCollection(ctx)
	}
	m := re.FindAllStringSubmatch(subject, -1)
	if len(m) == 0 {
		return strEmptyCollection(ctx)
	}
	list := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, row := range m {
		if len(row) > 1 {
			arrayAppend(list, data.NewStringValue(row[1]))
		} else if len(row) > 0 {
			arrayAppend(list, data.NewStringValue(row[0]))
		}
	}
	return strNewCollection(ctx, list)
}

func strEmptyCollection(ctx data.Context) (data.GetValue, data.Control) {
	return strNewCollection(ctx, data.NewArrayValue(nil))
}

func strNewCollection(ctx data.Context, items data.Value) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	cls, ok := vm.GetClass("Illuminate\\Support\\Collection")
	if !ok || cls == nil {
		var ctl data.Control
		cls, ctl = vm.GetOrLoadClass("Illuminate\\Support\\Collection")
		if ctl != nil {
			return data.NewArrayValue(nil), nil
		}
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = cv.SetProperty("items", items)
	return cv, nil
}

func strNumbers(ctx data.Context) (data.GetValue, data.Control) {
	re := regexp.MustCompile(`[^0-9]`)
	return data.NewStringValue(re.ReplaceAllString(strArg(ctx, 0), "")), nil
}

func strPadBoth(ctx data.Context) (data.GetValue, data.Control) {
	return strPad(ctx, 0)
}
func strPadLeft(ctx data.Context) (data.GetValue, data.Control) {
	return strPad(ctx, 1)
}
func strPadRight(ctx data.Context) (data.GetValue, data.Control) {
	return strPad(ctx, 2)
}

func strPad(ctx data.Context, mode int) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	length := strIntArg(ctx, 1, 0)
	pad := " "
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		pad = v.AsString()
	}
	if pad == "" {
		pad = " "
	}
	runes := []rune(value)
	padRunes := []rune(pad)
	for len(runes) < length {
		switch mode {
		case 1:
			runes = append([]rune{padRunes[0]}, runes...)
		case 2:
			runes = append(runes, padRunes[0])
		default:
			runes = append([]rune{padRunes[0]}, runes...)
			if len(runes) < length {
				runes = append(runes, padRunes[0])
			}
		}
	}
	if len(runes) > length {
		runes = runes[:length]
	}
	return data.NewStringValue(string(runes)), nil
}

func strParseCallback(ctx data.Context) (data.GetValue, data.Control) {
	cb := strArg(ctx, 0)
	def := kit.Arg(ctx, 1)
	defStr := ""
	if def != nil && !kit.IsNull(def) {
		defStr = def.AsString()
	}
	if strings.Contains(cb, "@") {
		parts := strings.SplitN(cb, "@", 2)
		arr := data.NewArrayValue(nil).(*data.ArrayValue)
		arrayAppend(arr, data.NewStringValue(parts[0]))
		arrayAppend(arr, data.NewStringValue(parts[1]))
		return arr, nil
	}
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	arrayAppend(arr, data.NewStringValue(cb))
	arrayAppend(arr, data.NewStringValue(defStr))
	return arr, nil
}

func strPluralWithPrepend(ctx data.Context, prepend bool) (data.GetValue, data.Control) {
	word := strArg(ctx, 0)
	count := strIntArg(ctx, 1, 2)
	if av, ok := kit.Arg(ctx, 1).(*data.ArrayValue); ok {
		count = len(kit.Entries(av))
	}
	out := englishPluralFallback(word, count)
	if prepend && strBoolArg(ctx, 2, prepend) {
		out = strings.TrimSpace(strings.Join([]string{data.NewIntValue(count).AsString(), out}, " "))
	}
	return data.NewStringValue(out), nil
}

func strPlural(ctx data.Context) (data.GetValue, data.Control) {
	return strPluralWithPrepend(ctx, strBoolArg(ctx, 2, false))
}

func strPluralStudly(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	count := strIntArg(ctx, 1, 2)
	parts := regexp.MustCompile(`(.)`).Split(value, -1)
	if len(parts) == 0 {
		return data.NewStringValue(value), nil
	}
	last := parts[len(parts)-1]
	parts[len(parts)-1] = englishPluralFallback(last, count)
	return data.NewStringValue(strings.Join(parts, "")), nil
}

func strPluralPascal(ctx data.Context) (data.GetValue, data.Control) {
	return strPluralStudly(ctx)
}

func strPassword(ctx data.Context) (data.GetValue, data.Control) {
	length := strIntArg(ctx, 0, 32)
	letters := strBoolArg(ctx, 1, true)
	numbers := strBoolArg(ctx, 2, true)
	symbols := strBoolArg(ctx, 3, true)
	spaces := strBoolArg(ctx, 4, false)
	var pool []rune
	if letters {
		pool = append(pool, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")...)
	}
	if numbers {
		pool = append(pool, []rune("0123456789")...)
	}
	if symbols {
		pool = append(pool, []rune("~!#$%^&*()-_.,<>?/\\{}[]|:;")...)
	}
	if spaces {
		pool = append(pool, ' ')
	}
	if len(pool) == 0 {
		pool = []rune("abcdefghijklmnopqrstuvwxyz")
	}
	b := make([]byte, length)
	_, _ = rand.Read(b)
	out := make([]rune, length)
	for i := range out {
		out[i] = pool[int(b[i])%len(pool)]
	}
	return data.NewStringValue(string(out)), nil
}

func strPosition(ctx data.Context) (data.GetValue, data.Control) {
	haystack, needle := strArg(ctx, 0), strArg(ctx, 1)
	offset := strIntArg(ctx, 2, 0)
	if needle == "" {
		return data.NewBoolValue(false), nil
	}
	sub := haystack
	if offset > 0 {
		sub = mbSubstr(haystack, offset, mbLen(haystack))
	} else if offset < 0 {
		sub = mbSubstr(haystack, offset, mbLen(haystack))
	}
	i := strings.Index(sub, needle)
	if i < 0 {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(i + max(0, offset)), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func strCreateRandomStringsUsing(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 0)
	if v == nil || kit.IsNull(v) {
		strRandomFactory = nil
	} else {
		strRandomFactory = v
	}
	return data.NewNullValue(), nil
}

func strCreateRandomStringsUsingSequence(ctx data.Context) (data.GetValue, data.Control) {
	strCreateRandomStringsNormally(ctx)
	return data.NewNullValue(), nil
}

func strCreateRandomStringsNormally(ctx data.Context) (data.GetValue, data.Control) {
	strRandomFactory = nil
	return data.NewNullValue(), nil
}

func strRepeat(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(strings.Repeat(strArg(ctx, 0), strIntArg(ctx, 1, 0))), nil
}

func strReplaceArray(ctx data.Context) (data.GetValue, data.Control) {
	search := strArg(ctx, 0)
	subject := strArg(ctx, 2)
	segments := strings.Split(subject, search)
	repl := kit.Arg(ctx, 1)
	var replace []string
	if av, ok := repl.(*data.ArrayValue); ok {
		for _, e := range kit.Entries(av) {
			replace = append(replace, e.Value.AsString())
		}
	}
	if len(segments) == 0 {
		return data.NewStringValue(subject), nil
	}
	out := segments[0]
	for i := 1; i < len(segments); i++ {
		r := search
		if len(replace) > 0 {
			r = replace[0]
			replace = replace[1:]
		}
		out += r + segments[i]
	}
	return data.NewStringValue(out), nil
}

func strReplaceFirst(ctx data.Context) (data.GetValue, data.Control) {
	search, replace, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	i := strings.Index(subject, search)
	if i < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(subject[:i] + replace + subject[i+len(search):]), nil
}

func strReplaceStart(ctx data.Context) (data.GetValue, data.Control) {
	search, _, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if strings.HasPrefix(subject, search) {
		return strReplaceFirst(ctx)
	}
	return data.NewStringValue(subject), nil
}

func strReplaceLast(ctx data.Context) (data.GetValue, data.Control) {
	search, replace, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	i := strings.LastIndex(subject, search)
	if i < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(subject[:i] + replace + subject[i+len(search):]), nil
}

func strReplaceEnd(ctx data.Context) (data.GetValue, data.Control) {
	search, _, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if strings.HasSuffix(subject, search) {
		return strReplaceLast(ctx)
	}
	return data.NewStringValue(subject), nil
}

func strReplaceMatches(ctx data.Context) (data.GetValue, data.Control) {
	pattern := strArg(ctx, 0)
	replace := kit.Arg(ctx, 1)
	subject := strArg(ctx, 2)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return data.NewStringValue(subject), nil
	}
	if fv, ok := replace.(*data.FuncValue); ok && fv != nil {
		out := re.ReplaceAllStringFunc(subject, func(m string) string {
			ret, ctl := kit.Call(ctx, fv, data.NewStringValue(m))
			if ctl != nil {
				return m
			}
			if ret != nil {
				return strValueString(ret)
			}
			return m
		})
		return data.NewStringValue(out), nil
	}
	rep := replace.AsString()
	return data.NewStringValue(re.ReplaceAllString(subject, rep)), nil
}

func strRemove(ctx data.Context) (data.GetValue, data.Control) {
	search := kit.Arg(ctx, 0)
	subject := strArg(ctx, 1)
	caseSensitive := strBoolArg(ctx, 2, true)
	needles := strNeedles(search)
	out := subject
	for _, n := range needles {
		if caseSensitive {
			out = strings.ReplaceAll(out, n, "")
		} else {
			out = strReplaceFold(out, n, "")
		}
	}
	return data.NewStringValue(out), nil
}

func strReplaceFold(haystack, needle, repl string) string {
	if needle == "" {
		return haystack
	}
	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(needle))
	return re.ReplaceAllString(haystack, repl)
}

func strReverse(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(mbReverse(strArg(ctx, 0))), nil
}

func strHeadline(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	parts := strings.Fields(value)
	if len(parts) > 1 {
		for i, p := range parts {
			parts[i] = strings.Title(strings.ToLower(p))
		}
		return data.NewStringValue(strings.Join(parts, " ")), nil
	}
	words := strUcsplitImpl(value)
	for i, w := range words {
		words[i] = strings.Title(strings.ToLower(w))
	}
	collapsed := strings.ReplaceAll(strings.Join(words, "_"), "-", "_")
	collapsed = strings.ReplaceAll(collapsed, " ", "_")
	return data.NewStringValue(strings.ReplaceAll(collapsed, "_", " ")), nil
}

func strInitials(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	cap := strBoolArg(ctx, 1, false)
	parts := strings.Fields(value)
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(mbSubstr(p, 0, 1))
	}
	out := b.String()
	if cap {
		out = strings.ToUpper(out)
	}
	return data.NewStringValue(out), nil
}

func strApa(ctx data.Context) (data.GetValue, data.Control) {
	value := strings.TrimSpace(strArg(ctx, 0))
	if value == "" {
		return data.NewStringValue(value), nil
	}
	parts := strings.Fields(value)
	for i, w := range parts {
		if i == 0 {
			parts[i] = mbUcfirst(strings.ToLower(w))
		} else {
			parts[i] = strings.ToLower(w)
		}
	}
	return data.NewStringValue(strings.Join(parts, " ")), nil
}

func strSingular(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(englishSingularFallback(strArg(ctx, 0))), nil
}

func strTrim(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		return data.NewStringValue(strings.Trim(value, v.AsString())), nil
	}
	return data.NewStringValue(strings.TrimSpace(value)), nil
}

func strLtrim(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		return data.NewStringValue(strings.TrimLeft(value, v.AsString())), nil
	}
	return data.NewStringValue(strings.TrimLeft(value, " \t\n\r\v\f\x00")), nil
}

func strRtrim(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		return data.NewStringValue(strings.TrimRight(value, v.AsString())), nil
	}
	return data.NewStringValue(strings.TrimRight(value, " \t\n\r\v\f\x00")), nil
}

func strSquish(ctx data.Context) (data.GetValue, data.Control) {
	value, ctl := strTrim(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s := strValueString(value)
	re := regexp.MustCompile(`\s+`)
	return data.NewStringValue(re.ReplaceAllString(s, " ")), nil
}

func strDoesntStartWith(ctx data.Context) (data.GetValue, data.Control) {
	b, ctl := strStartsWith(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := b.(data.AsBool).AsBool()
	return data.NewBoolValue(!v), nil
}

func strPascal(ctx data.Context) (data.GetValue, data.Control) {
	norm := strBoolArg(ctx, 1, false)
	return data.NewStringValue(studlyCached(strArg(ctx, 0), norm)), nil
}

func strSubstrCount(ctx data.Context) (data.GetValue, data.Control) {
	haystack, needle := strArg(ctx, 0), strArg(ctx, 1)
	offset := strIntArg(ctx, 2, 0)
	if needle == "" {
		return data.NewIntValue(0), nil
	}
	sub := haystack
	if offset != 0 {
		if offset < 0 || offset > len(haystack) {
			sub = haystack
		} else {
			sub = haystack[offset:]
		}
	}
	return data.NewIntValue(strings.Count(sub, needle)), nil
}

func strSubstrReplace(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	replace := strArg(ctx, 1)
	offset := strIntArg(ctx, 2, 0)
	length := strIntArg(ctx, 3, mbLen(s))
	return data.NewStringValue(mbSubstr(s, 0, offset) + replace + mbSubstr(s, offset+length, mbLen(s))), nil
}

func strSwap(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 1)
	m := kit.Arg(ctx, 0)
	if av, ok := m.(*data.ArrayValue); ok {
		repl := make([]string, 0, len(kit.Entries(av))*2)
		for _, e := range kit.Entries(av) {
			repl = append(repl, e.KeyStr, e.Value.AsString())
		}
		return data.NewStringValue(strings.NewReplacer(repl...).Replace(subject)), nil
	}
	return data.NewStringValue(subject), nil
}

func strTake(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	limit := strIntArg(ctx, 1, 0)
	if limit < 0 {
		return data.NewStringValue(mbSubstr(s, limit, mbLen(s))), nil
	}
	return data.NewStringValue(mbSubstr(s, 0, limit)), nil
}

func strToBase64(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(base64.StdEncoding.EncodeToString([]byte(strArg(ctx, 0)))), nil
}

func strFromBase64(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	strict := strBoolArg(ctx, 1, false)
	dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil && strict {
		return data.NewBoolValue(false), nil
	}
	if err != nil {
		dec, _ = base64.RawStdEncoding.DecodeString(s)
	}
	return data.NewStringValue(string(dec)), nil
}

func strLcfirst(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(mbLcfirst(strArg(ctx, 0))), nil
}

func strUcfirst(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(mbUcfirst(strArg(ctx, 0))), nil
}

func strUcwords(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	seps := " \t\r\n\f\v"
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		seps = v.AsString()
	}
	re := regexp.MustCompile(`(^|[` + regexp.QuoteMeta(seps) + `])([\p{Ll}])`)
	out := re.ReplaceAllStringFunc(s, func(m string) string {
		r := []rune(m)
		if len(r) >= 2 {
			r[len(r)-1] = unicode.ToUpper(r[len(r)-1])
		}
		return string(r)
	})
	return data.NewStringValue(out), nil
}

func strUcsplit(ctx data.Context) (data.GetValue, data.Control) {
	parts := strUcsplitImpl(strArg(ctx, 0))
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, p := range parts {
		arrayAppend(arr, data.NewStringValue(p))
	}
	return arr, nil
}

func strUcsplitImpl(s string) []string {
	re := regexp.MustCompile(`(?=\p{Lu})`)
	return re.Split(s, -1)
}

func strWordCount(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	return data.NewIntValue(len(strings.Fields(s))), nil
}

func strWordWrap(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	width := strIntArg(ctx, 1, 75)
	br := "\n"
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		br = v.AsString()
	}
	cut := strBoolArg(ctx, 3, false)
	words := strings.Fields(s)
	if len(words) == 0 {
		return data.NewStringValue(s), nil
	}
	var lines []string
	var line strings.Builder
	for _, w := range words {
		if line.Len() == 0 {
			line.WriteString(w)
			continue
		}
		if line.Len()+1+len(w) > width && !cut {
			lines = append(lines, line.String())
			line.Reset()
			line.WriteString(w)
		} else {
			line.WriteString(" ")
			line.WriteString(w)
		}
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}
	return data.NewStringValue(strings.Join(lines, br)), nil
}

func strUUID7(ctx data.Context) (data.GetValue, data.Control) {
	ms := uint64(time.Now().UnixMilli())
	return data.NewStringValue(formatUUIDv7(ms)), nil
}

func strOrderedUUID(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(formatUUIDv4()), nil
}

func strCreateUuidsUsing(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 0)
	if v == nil || kit.IsNull(v) {
		strUUIDFactory = nil
	} else {
		strUUIDFactory = v
	}
	return data.NewNullValue(), nil
}

func strCreateUuidsUsingSequence(ctx data.Context) (data.GetValue, data.Control) {
	strCreateUuidsNormally(ctx)
	return data.NewNullValue(), nil
}

func strFreezeUuids(ctx data.Context) (data.GetValue, data.Control) {
	u := formatUUIDv4()
	strUUIDFactory = data.NewStringValue(u)
	cb := kit.Arg(ctx, 0)
	if cb != nil && !kit.IsNull(cb) {
		_, _ = kit.Call(ctx, cb, data.NewStringValue(u))
	}
	strUUIDFactory = nil
	return data.NewStringValue(u), nil
}

func strCreateUuidsNormally(ctx data.Context) (data.GetValue, data.Control) {
	strUUIDFactory = nil
	return data.NewNullValue(), nil
}

func strULID(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(generateULIDString(0)), nil
}

func strCreateUlidsNormally(ctx data.Context) (data.GetValue, data.Control) {
	strULIDFactory = nil
	return data.NewNullValue(), nil
}

func strCreateUlidsUsing(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 0)
	if v == nil || kit.IsNull(v) {
		strULIDFactory = nil
	} else {
		strULIDFactory = v
	}
	return data.NewNullValue(), nil
}

func strCreateUlidsUsingSequence(ctx data.Context) (data.GetValue, data.Control) {
	strCreateUlidsNormally(ctx)
	return data.NewNullValue(), nil
}

func strFreezeUlids(ctx data.Context) (data.GetValue, data.Control) {
	u := generateULIDString(0)
	cb := kit.Arg(ctx, 0)
	if cb != nil && !kit.IsNull(cb) {
		_, _ = kit.Call(ctx, cb, data.NewStringValue(u))
	}
	return data.NewStringValue(u), nil
}

func strFlushCache(ctx data.Context) (data.GetValue, data.Control) {
	strFlushCaches()
	return data.NewNullValue(), nil
}

func strResetFactoryState(ctx data.Context) (data.GetValue, data.Control) {
	strCreateRandomStringsNormally(ctx)
	strCreateUlidsNormally(ctx)
	strCreateUuidsNormally(ctx)
	return data.NewNullValue(), nil
}
