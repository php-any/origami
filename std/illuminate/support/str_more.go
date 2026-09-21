package support

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/data"
)

func strContainsAll(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needles, _ := ctx.GetIndexValue(1)
	ignore := false
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if b, err := v.(data.AsBool).AsBool(); err == nil {
			ignore = b
		}
	}
	h := haystack
	if ignore {
		h = strings.ToLower(h)
	}
	for _, n := range strNeedles(needles) {
		if n == "" {
			continue
		}
		nn := n
		if ignore {
			nn = strings.ToLower(nn)
		}
		if !strings.Contains(h, nn) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func strAfterLast(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	search := strArg(ctx, 1)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	i := strings.LastIndex(subject, search)
	if i < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(subject[i+len(search):]), nil
}

func strBeforeLast(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	search := strArg(ctx, 1)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	i := strings.LastIndex(subject, search)
	if i < 0 {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(subject[:i]), nil
}

func strAscii(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(asciiFold(strArg(ctx, 0))), nil
}

func asciiFold(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r < utf8.RuneSelf {
			b.WriteByte(byte(r))
			continue
		}
		switch r {
		case 'à', 'á', 'â', 'ã', 'ä', 'å', 'ā', 'ă', 'ą':
			b.WriteByte('a')
		case 'À', 'Á', 'Â', 'Ã', 'Ä', 'Å', 'Ā', 'Ă', 'Ą':
			b.WriteByte('A')
		case 'è', 'é', 'ê', 'ë', 'ē', 'ĕ', 'ė', 'ę':
			b.WriteByte('e')
		case 'È', 'É', 'Ê', 'Ë', 'Ē', 'Ĕ', 'Ė', 'Ę':
			b.WriteByte('E')
		case 'ì', 'í', 'î', 'ï', 'ī':
			b.WriteByte('i')
		case 'Ì', 'Í', 'Î', 'Ï', 'Ī':
			b.WriteByte('I')
		case 'ò', 'ó', 'ô', 'õ', 'ö', 'ō', 'ø':
			b.WriteByte('o')
		case 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', 'Ō', 'Ø':
			b.WriteByte('O')
		case 'ù', 'ú', 'û', 'ü', 'ū':
			b.WriteByte('u')
		case 'Ù', 'Ú', 'Û', 'Ü', 'Ū':
			b.WriteByte('U')
		case 'ç', 'ć', 'č':
			b.WriteByte('c')
		case 'ñ':
			b.WriteByte('n')
		case 'ß':
			b.WriteString("ss")
		default:
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func strSlug(ctx data.Context) (data.GetValue, data.Control) {
	title := strArg(ctx, 0)
	sep := "-"
	if s := strArg(ctx, 1); s != "" {
		sep = s
	}
	lang := "en"
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		lang = v.AsString()
	}
	if lang != "" {
		title = asciiFold(title)
	}
	flip := "-"
	if sep == "-" {
		flip = "_"
	}
	title = strings.ReplaceAll(title, flip, sep)
	dict := map[string]string{"@": "at"}
	if raw, ok := ctx.GetIndexValue(3); ok && raw != nil && !isNull(raw) {
		dict = map[string]string{}
		for _, e := range toEntries(raw) {
			dict[e.keyStr] = e.value.AsString()
		}
	}
	for k, v := range dict {
		title = strings.ReplaceAll(title, k, sep+v+sep)
	}
	title = strings.ToLower(title)
	var b strings.Builder
	lastSep := false
	for _, r := range title {
		ok := unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || strings.ContainsRune(sep, r)
		if !ok {
			continue
		}
		if unicode.IsSpace(r) || strings.ContainsRune(sep, r) {
			if !lastSep && b.Len() > 0 {
				b.WriteString(sep)
				lastSep = true
			}
			continue
		}
		b.WriteRune(r)
		lastSep = false
	}
	return data.NewStringValue(strings.Trim(b.String(), sep)), nil
}

func strTrim(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if v, ok := ctx.GetIndexValue(1); ok && v != nil && !isNull(v) {
		return data.NewStringValue(strings.Trim(s, v.AsString())), nil
	}
	return data.NewStringValue(strings.TrimSpace(s)), nil
}

func strLtrim(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if v, ok := ctx.GetIndexValue(1); ok && v != nil && !isNull(v) {
		return data.NewStringValue(strings.TrimLeft(s, v.AsString())), nil
	}
	return data.NewStringValue(strings.TrimLeft(s, " \t\n\r")), nil
}

func strRtrim(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if v, ok := ctx.GetIndexValue(1); ok && v != nil && !isNull(v) {
		return data.NewStringValue(strings.TrimRight(s, v.AsString())), nil
	}
	return data.NewStringValue(strings.TrimRight(s, " \t\n\r")), nil
}

func strReplaceFirst(ctx data.Context) (data.GetValue, data.Control) {
	search, replace, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if search == "" {
		return data.NewStringValue(subject), nil
	}
	return data.NewStringValue(strings.Replace(subject, search, replace, 1)), nil
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

func strReplaceStart(ctx data.Context) (data.GetValue, data.Control) {
	search, replace, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if search != "" && strings.HasPrefix(subject, search) {
		return data.NewStringValue(replace + subject[len(search):]), nil
	}
	return data.NewStringValue(subject), nil
}

func strReplaceEnd(ctx data.Context) (data.GetValue, data.Control) {
	search, replace, subject := strArg(ctx, 0), strArg(ctx, 1), strArg(ctx, 2)
	if search != "" && strings.HasSuffix(subject, search) {
		return data.NewStringValue(subject[:len(subject)-len(search)] + replace), nil
	}
	return data.NewStringValue(subject), nil
}

func strHeadline(ctx data.Context) (data.GetValue, data.Control) {
	s := strings.ReplaceAll(strArg(ctx, 0), "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	parts := strings.Fields(s)
	for i, p := range parts {
		r := []rune(strings.ToLower(p))
		if len(r) == 0 {
			continue
		}
		r[0] = unicode.ToUpper(r[0])
		parts[i] = string(r)
	}
	return data.NewStringValue(strings.Join(parts, " ")), nil
}

func strUcfirst(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if s == "" {
		return data.NewStringValue(""), nil
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return data.NewStringValue(string(r)), nil
}

func strLcfirst(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if s == "" {
		return data.NewStringValue(""), nil
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return data.NewStringValue(string(r)), nil
}

func strRepeat(ctx data.Context) (data.GetValue, data.Control) {
	n := 0
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			n, _ = iv.AsInt()
		}
	}
	if n < 0 {
		n = 0
	}
	return data.NewStringValue(strings.Repeat(strArg(ctx, 0), n)), nil
}

func strReverse(ctx data.Context) (data.GetValue, data.Control) {
	r := []rune(strArg(ctx, 0))
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return data.NewStringValue(string(r)), nil
}

func strWrap(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	before := strArg(ctx, 1)
	after := before
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		after = v.AsString()
	}
	return data.NewStringValue(before + value + after), nil
}

func strIsAscii(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func strIsJson(ctx data.Context) (data.GetValue, data.Control) {
	s := strings.TrimSpace(strArg(ctx, 0))
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	var js any
	err := json.Unmarshal([]byte(s), &js)
	return data.NewBoolValue(err == nil), nil
}

func strFlushCache(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func strOf(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, strArg(ctx, 0))
}

func strToBase64(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(base64.StdEncoding.EncodeToString([]byte(strArg(ctx, 0)))), nil
}

func strFromBase64(ctx data.Context) (data.GetValue, data.Control) {
	b, err := base64.StdEncoding.DecodeString(strArg(ctx, 0))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(string(b)), nil
}
