package support

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
)

func phpList(vals ...data.Value) *data.ArrayValue {
	av := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, v := range vals {
		av.AppendValue(v)
	}
	return av
}

func strParseCallback(ctx data.Context) (data.GetValue, data.Control) {
	callback := strArg(ctx, 0)
	var def data.Value = data.NewNullValue()
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		def = v
	}
	anon := "@anonymous\x00"
	if strings.Contains(callback, anon) {
		if strings.Count(callback, "@") > 1 {
			i := strings.LastIndex(callback, "@")
			return phpList(data.NewStringValue(callback[:i]), data.NewStringValue(callback[i+1:])), nil
		}
		return phpList(data.NewStringValue(callback), def), nil
	}
	if strings.Contains(callback, "@") {
		parts := strings.SplitN(callback, "@", 2)
		second := ""
		if len(parts) > 1 {
			second = parts[1]
		}
		return phpList(data.NewStringValue(parts[0]), data.NewStringValue(second)), nil
	}
	return phpList(data.NewStringValue(callback), def), nil
}

func strSubstrCount(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needle := strArg(ctx, 1)
	if needle == "" {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(strings.Count(haystack, needle)), nil
}

func strBetween(ctx data.Context) (data.GetValue, data.Control) {
	subject := strArg(ctx, 0)
	from := strArg(ctx, 1)
	to := strArg(ctx, 2)
	if from == "" && to == "" {
		return data.NewStringValue(subject), nil
	}
	s := subject
	if from != "" {
		if i := strings.Index(s, from); i >= 0 {
			s = s[i+len(from):]
		}
	}
	if to != "" {
		if i := strings.Index(s, to); i >= 0 {
			s = s[:i]
		}
	}
	return data.NewStringValue(s), nil
}

func strBetweenFirst(ctx data.Context) (data.GetValue, data.Control) {
	return strBetween(ctx)
}

func strPosition(ctx data.Context) (data.GetValue, data.Control) {
	haystack := strArg(ctx, 0)
	needle := strArg(ctx, 1)
	offset := 0
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				offset = n
			}
		}
	}
	runes := []rune(haystack)
	if offset < 0 {
		offset = 0
	}
	if offset > len(runes) {
		return data.NewBoolValue(false), nil
	}
	rest := string(runes[offset:])
	i := strings.Index(rest, needle)
	if i < 0 {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(offset + len([]rune(rest[:i]))), nil
}

func strRemove(ctx data.Context) (data.GetValue, data.Control) {
	search, _ := ctx.GetIndexValue(0)
	subject := strArg(ctx, 1)
	for _, n := range strNeedles(search) {
		if n == "" {
			continue
		}
		subject = strings.ReplaceAll(subject, n, "")
	}
	return data.NewStringValue(subject), nil
}

func strReplaceArray(ctx data.Context) (data.GetValue, data.Control) {
	search := strArg(ctx, 0)
	replace, _ := ctx.GetIndexValue(1)
	subject := strArg(ctx, 2)
	reps := strNeedles(replace)
	i := 0
	for {
		j := strings.Index(subject, search)
		if j < 0 || i >= len(reps) {
			break
		}
		subject = subject[:j] + reps[i] + subject[j+len(search):]
		i++
	}
	return data.NewStringValue(subject), nil
}

func strSquish(ctx data.Context) (data.GetValue, data.Control) {
	s := strings.TrimSpace(strArg(ctx, 0))
	space := false
	var b strings.Builder
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !space {
				b.WriteByte(' ')
				space = true
			}
			continue
		}
		space = false
		b.WriteRune(r)
	}
	return data.NewStringValue(b.String()), nil
}

func strPad(s string, length int, pad string, dir int) string {
	if pad == "" {
		pad = " "
	}
	n := length - len([]rune(s))
	if n <= 0 {
		return s
	}
	var extra strings.Builder
	for extra.Len() < n {
		extra.WriteString(pad)
	}
	p := []rune(extra.String())[:n]
	ps := string(p)
	switch dir {
	case -1:
		return ps + s
	case 1:
		return s + ps
	default:
		left := n / 2
		return string(p[:left]) + s + string(p[left:])
	}
}

func strPadInt(ctx data.Context, i int, def int) int {
	if v, ok := ctx.GetIndexValue(i); ok && v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				return n
			}
		}
	}
	return def
}

func strPadLeft(ctx data.Context) (data.GetValue, data.Control) {
	pad := " "
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		pad = v.AsString()
	}
	return data.NewStringValue(strPad(strArg(ctx, 0), strPadInt(ctx, 1, 0), pad, -1)), nil
}

func strPadRight(ctx data.Context) (data.GetValue, data.Control) {
	pad := " "
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		pad = v.AsString()
	}
	return data.NewStringValue(strPad(strArg(ctx, 0), strPadInt(ctx, 1, 0), pad, 1)), nil
}

func strPadBoth(ctx data.Context) (data.GetValue, data.Control) {
	pad := " "
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		pad = v.AsString()
	}
	return data.NewStringValue(strPad(strArg(ctx, 0), strPadInt(ctx, 1, 0), pad, 0)), nil
}

func strTake(ctx data.Context) (data.GetValue, data.Control) {
	s := []rune(strArg(ctx, 0))
	n := strPadInt(ctx, 1, 0)
	if n < 0 {
		if -n >= len(s) {
			return data.NewStringValue(string(s)), nil
		}
		return data.NewStringValue(string(s[len(s)+n:])), nil
	}
	if n >= len(s) {
		return data.NewStringValue(string(s)), nil
	}
	return data.NewStringValue(string(s[:n])), nil
}

func strUnwrap(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	before := strArg(ctx, 1)
	after := before
	if v, ok := ctx.GetIndexValue(2); ok && v != nil && !isNull(v) {
		after = v.AsString()
	}
	if strings.HasPrefix(value, before) {
		value = value[len(before):]
	}
	if strings.HasSuffix(value, after) {
		value = value[:len(value)-len(after)]
	}
	return data.NewStringValue(value), nil
}

func strDoesntContain(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := strContains(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := v.(*data.BoolValue); ok {
		return data.NewBoolValue(!bv.Value), nil
	}
	return data.NewBoolValue(true), nil
}

func strDoesntStartWith(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := strStartsWith(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := v.(*data.BoolValue); ok {
		return data.NewBoolValue(!bv.Value), nil
	}
	return data.NewBoolValue(true), nil
}

func strDoesntEndWith(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := strEndsWith(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := v.(*data.BoolValue); ok {
		return data.NewBoolValue(!bv.Value), nil
	}
	return data.NewBoolValue(true), nil
}

var urlLike = regexp.MustCompile(`(?i)^https?://`)

func strIsUrl(ctx data.Context) (data.GetValue, data.Control) {
	s := strArg(ctx, 0)
	if !urlLike.MatchString(s) {
		return data.NewBoolValue(false), nil
	}
	u, err := url.Parse(s)
	return data.NewBoolValue(err == nil && u.Host != ""), nil
}

func strIsUlid(ctx data.Context) (data.GetValue, data.Control) {
	s := strings.ToUpper(strArg(ctx, 0))
	if len(s) != 26 {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z')) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func strCharAt(ctx data.Context) (data.GetValue, data.Control) {
	s := []rune(strArg(ctx, 0))
	i := strPadInt(ctx, 1, 0)
	if i < 0 {
		i = len(s) + i
	}
	if i < 0 || i >= len(s) {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(string(s[i])), nil
}

func strNumbers(ctx data.Context) (data.GetValue, data.Control) {
	var b strings.Builder
	for _, r := range strArg(ctx, 0) {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return data.NewStringValue(b.String()), nil
}

func strUcwords(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(strings.Title(strings.ToLower(strArg(ctx, 0)))), nil
}

func strWordCount(ctx data.Context) (data.GetValue, data.Control) {
	s := strings.TrimSpace(strArg(ctx, 0))
	if s == "" {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(strings.Fields(s))), nil
}

func strFreezeUuids(ctx data.Context) (data.GetValue, data.Control) {
	u, ctl := strUUID(ctx)
	if ctl != nil {
		return nil, ctl
	}
	uv, _ := u.(data.Value)
	strUUIDFactoryMu.Lock()
	strUUIDFactory = data.NewFuncValue(&frozenUUIDFunc{u: uv})
	strUUIDFactoryMu.Unlock()
	if cb, ok := ctx.GetIndexValue(0); ok && cb != nil && !isNull(cb) {
		v, ctl := callValue(ctx, cb, uv)
		strCreateUuidsNormally(ctx)
		return v, ctl
	}
	return u, nil
}

type frozenUUIDFunc struct{ u data.Value }

func (f *frozenUUIDFunc) Call(ctx data.Context) (data.GetValue, data.Control) {
	return f.u, nil
}
func (f *frozenUUIDFunc) GetName() string            { return "frozenUuid" }
func (f *frozenUUIDFunc) GetParams() []data.GetValue { return nil }
func (f *frozenUUIDFunc) GetVariables() []data.Variable {
	return nil
}

func strCreateUuidsUsingSequence(ctx data.Context) (data.GetValue, data.Control) {
	seq, _ := ctx.GetIndexValue(0)
	items := strNeedles(seq)
	i := 0
	fn := &seqUUIDFunc{items: items, i: &i, ctxsrc: ctx}
	strUUIDFactoryMu.Lock()
	strUUIDFactory = data.NewFuncValue(fn)
	strUUIDFactoryMu.Unlock()
	return data.NewNullValue(), nil
}

type seqUUIDFunc struct {
	items  []string
	i      *int
	ctxsrc data.Context
}

func (f *seqUUIDFunc) Call(ctx data.Context) (data.GetValue, data.Control) {
	if *f.i < len(f.items) {
		s := f.items[*f.i]
		*f.i++
		return wrapUUID(ctx, s), nil
	}
	return wrapUUID(ctx, generateUUIDv4()), nil
}
func (f *seqUUIDFunc) GetName() string            { return "seqUuid" }
func (f *seqUUIDFunc) GetParams() []data.GetValue { return nil }
func (f *seqUUIDFunc) GetVariables() []data.Variable {
	return nil
}
