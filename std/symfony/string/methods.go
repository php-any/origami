package sfstring

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/preg"
)

func methodConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	s := argString(ctx, 0, "")
	if strKindOf(cv) != kindByte {
		if ctl := mustUTF8(s); ctl != nil {
			return nil, ctl
		}
	}
	_ = cv.SetProperty("string", data.NewStringValue(s))
	_ = cv.SetProperty("ignoreCase", data.NewBoolValue(false))
	return data.NewNullValue(), nil
}

func methodClone(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	_ = cv.SetProperty("ignoreCase", data.NewBoolValue(false))
	return data.NewNullValue(), nil
}

func methodToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(getString(strSelf(ctx))), nil
}

func methodJsonSerialize(ctx data.Context) (data.GetValue, data.Control) {
	return methodToString(ctx)
}

func methodLength(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	return data.NewIntValue(unitCount(getString(cv), strKindOf(cv))), nil
}

func methodIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(getString(strSelf(ctx)) == ""), nil
}

func methodIgnoreCase(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	out := newOf(cv, getString(cv))
	_ = out.SetProperty("ignoreCase", data.NewBoolValue(true))
	return out, nil
}

func methodAppend(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	s := getString(cv) + strings.Join(variadicStrings(ctx, 0), "")
	if strKindOf(cv) != kindByte {
		if ctl := mustUTF8(s); ctl != nil {
			return nil, ctl
		}
	}
	return newOf(cv, s), nil
}

func methodPrepend(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	s := strings.Join(variadicStrings(ctx, 0), "") + getString(cv)
	if strKindOf(cv) != kindByte {
		if ctl := mustUTF8(s); ctl != nil {
			return nil, ctl
		}
	}
	return newOf(cv, s), nil
}

func methodSlice(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	start := argInt(ctx, 0, 0)
	length := argOptionalInt(ctx, 1)
	return newOf(cv, sliceUnits(getString(cv), start, length, strKindOf(cv))), nil
}

func methodSplice(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	repl := argString(ctx, 0, "")
	if strKindOf(cv) != kindByte {
		if ctl := mustUTF8(repl); ctl != nil {
			return nil, ctl
		}
	}
	start := argInt(ctx, 1, 0)
	length := argOptionalInt(ctx, 2)
	return newOf(cv, spliceStr(getString(cv), repl, start, length, strKindOf(cv))), nil
}

func methodChunk(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	n := argInt(ctx, 0, 1)
	if n < 1 {
		return nil, throwInvalid("The chunk length must be greater than zero.")
	}
	parts := chunkStr(getString(cv), n, strKindOf(cv))
	out := make([]*data.ClassValue, len(parts))
	for i, p := range parts {
		out[i] = newOf(cv, p)
	}
	return toArray(out), nil
}

func methodSplit(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	delim := argString(ctx, 0, "")
	if delim == "" {
		return nil, throwInvalid("Split delimiter is empty.")
	}
	limit := 2147483647
	if v := argValue(ctx, 1); v != nil {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				limit = n
			}
		}
	}
	if limit < 1 {
		return nil, throwInvalid("Split limit must be a positive integer.")
	}
	if flags := argValue(ctx, 2); flags != nil {
		// flags 路径交给父类 PHP preg_split；这里做 unicode 标记
		_ = flags
	}
	k := strKindOf(cv)
	if k != kindByte {
		if ctl := mustUTF8(delim); ctl != nil {
			return nil, throwInvalid("Split delimiter is not a valid UTF-8 string.")
		}
	}
	parts := splitStr(getString(cv), delim, limit, ignoreCaseOf(cv), k)
	out := make([]*data.ClassValue, len(parts))
	for i, p := range parts {
		out[i] = newOf(cv, p)
	}
	return toArray(out), nil
}

func methodIndexOf(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	needle := argValue(ctx, 0)
	offset := argInt(ctx, 1, 0)
	best := -1
	for _, n := range needles(ctx, needle) {
		if i, ok := indexOfUnits(getString(cv), n, offset, ignoreCaseOf(cv), strKindOf(cv)); ok {
			if best < 0 || i < best {
				best = i
			}
		}
	}
	if best < 0 {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(best), nil
}

func methodIndexOfLast(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	needle := argValue(ctx, 0)
	offset := argInt(ctx, 1, 0)
	best := -1
	found := false
	for _, n := range needles(ctx, needle) {
		if i, ok := lastIndexOfUnits(getString(cv), n, offset, ignoreCaseOf(cv), strKindOf(cv)); ok {
			if !found || i >= best {
				best = i
				found = true
				offset = i
			}
		}
	}
	if !found {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(best), nil
}

func methodStartsWith(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	for _, n := range needles(ctx, argValue(ctx, 0)) {
		if startsWithStr(getString(cv), n, ignoreCaseOf(cv), strKindOf(cv)) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func methodEndsWith(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	for _, n := range needles(ctx, argValue(ctx, 0)) {
		if endsWithStr(getString(cv), n, ignoreCaseOf(cv), strKindOf(cv)) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func methodEqualsTo(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	for _, n := range needles(ctx, argValue(ctx, 0)) {
		if equalsStr(getString(cv), n, ignoreCaseOf(cv)) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func methodContainsAny(ctx data.Context) (data.GetValue, data.Control) {
	ret, ctl := methodIndexOf(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if _, isNull := ret.(*data.NullValue); isNull {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func methodReplace(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	from := argString(ctx, 0, "")
	to := argString(ctx, 1, "")
	k := strKindOf(cv)
	if k != kindByte {
		if from != "" && !isValidUTF8(from) {
			return newOf(cv, getString(cv)), nil
		}
		if to != "" {
			if ctl := mustUTF8(to); ctl != nil {
				return nil, ctl
			}
		}
	}
	return newOf(cv, replaceStr(getString(cv), from, to, ignoreCaseOf(cv), k)), nil
}

func methodReplaceMatches(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	pattern := argString(ctx, 0, "")
	to := argValue(ctx, 1)
	if ignoreCaseOf(cv) && !strings.Contains(pattern, "i") {
		pattern = withRegexFlag(pattern, 'i')
	}
	k := strKindOf(cv)
	if k != kindByte && !strings.Contains(pattern, "u") {
		pattern = withRegexFlag(pattern, 'u')
	}
	if to != nil {
		if _, isFn := to.(*data.FuncValue); isFn {
			ret, ctl := callPHPFunc(ctx, "preg_replace_callback", data.NewStringValue(pattern), to, data.NewStringValue(getString(cv)))
			if ctl != nil {
				return nil, ctl
			}
			if ret == nil {
				return nil, throwInvalid("Matching failed with unknown error code.")
			}
			if bv, ok := ret.(*data.BoolValue); ok && !truthyBool(bv) {
				return nil, throwInvalid("Matching failed with unknown error code.")
			}
			return newOf(cv, valueString(ctx, ret.(data.Value))), nil
		}
		if _, isBound := to.(*data.BoundFuncValue); isBound {
			ret, ctl := callPHPFunc(ctx, "preg_replace_callback", data.NewStringValue(pattern), to, data.NewStringValue(getString(cv)))
			if ctl != nil {
				return nil, ctl
			}
			return newOf(cv, valueString(ctx, ret.(data.Value))), nil
		}
	}
	repl := ""
	if to != nil {
		repl = valueString(ctx, to)
		if k != kindByte {
			if ctl := mustUTF8(repl); ctl != nil {
				return nil, ctl
			}
		}
	}
	re, err := preg.CompileAny(pattern)
	if err != nil {
		return nil, throwInvalid(err.Error())
	}
	out, _ := preg.ReplaceAllPhp(re, getString(cv), repl, -1)
	return newOf(cv, out), nil
}

func withRegexFlag(pattern string, flag byte) string {
	if pattern == "" {
		return pattern
	}
	return pattern + string(flag)
}

func truthyBool(b *data.BoolValue) bool {
	ok, err := b.AsBool()
	return err == nil && ok
}

func methodLower(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	s := getString(cv)
	if strKindOf(cv) == kindByte {
		return newOf(cv, strings.ToLower(s)), nil
	}
	s = strings.ReplaceAll(s, "İ", "i̇")
	return newOf(cv, strings.ToLower(s)), nil
}

func methodUpper(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	return newOf(cv, strings.ToUpper(getString(cv))), nil
}

func methodFolded(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	s := getString(cv)
	if strKindOf(cv) == kindByte {
		return newOf(cv, strings.ToLower(s)), nil
	}
	s = strings.ToLower(replacePairs(s, foldFrom, foldTo))
	return newOf(cv, s), nil
}

func methodCamel(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	s := getString(cv)
	if strKindOf(cv) == kindByte {
		return newOf(cv, camelByte(s)), nil
	}
	return newOf(cv, camelUnicode(s)), nil
}

func methodSnake(ctx data.Context) (data.GetValue, data.Control) {
	cam, ctl := methodCamel(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cv := cam.(*data.ClassValue)
	return newOf(cv, snakeFromCamel(getString(cv))), nil
}

func methodKebab(ctx data.Context) (data.GetValue, data.Control) {
	sn, ctl := methodSnake(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cv := sn.(*data.ClassValue)
	return newOf(cv, strings.ReplaceAll(getString(cv), "_", "-")), nil
}

func methodPascal(ctx data.Context) (data.GetValue, data.Control) {
	cam, ctl := methodCamel(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cv := cam.(*data.ClassValue)
	return newOf(cv, titleStr(getString(cv), false, strKindOf(cv))), nil
}

func methodTitle(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	return newOf(cv, titleStr(getString(cv), argBool(ctx, 0, false), strKindOf(cv))), nil
}

func methodTrim(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	chars := argString(ctx, 0, defaultTrim(cv))
	if strKindOf(cv) != kindByte && !isValidUTF8(chars) {
		return nil, throwInvalid("Invalid UTF-8 chars.")
	}
	return newOf(cv, trimChars(getString(cv), chars, true, true)), nil
}

func methodTrimStart(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	chars := argString(ctx, 0, defaultTrim(cv))
	if strKindOf(cv) != kindByte && !isValidUTF8(chars) {
		return nil, throwInvalid("Invalid UTF-8 chars.")
	}
	return newOf(cv, trimChars(getString(cv), chars, true, false)), nil
}

func methodTrimEnd(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	chars := argString(ctx, 0, defaultTrim(cv))
	if strKindOf(cv) != kindByte && !isValidUTF8(chars) {
		return nil, throwInvalid("Invalid UTF-8 chars.")
	}
	return newOf(cv, trimChars(getString(cv), chars, false, true)), nil
}

func defaultTrim(cv *data.ClassValue) string {
	if strKindOf(cv) == kindByte {
		return byteTrimDefault
	}
	return unicodeTrimDefault
}

func methodTrimPrefix(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	v := argValue(ctx, 0)
	if av, ok := v.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z == nil || z.Value == nil {
				continue
			}
			t := trimPrefixStr(getString(cv), valueString(ctx, z.Value), ignoreCaseOf(cv))
			if t != getString(cv) {
				return newOf(cv, t), nil
			}
		}
		return newOf(cv, getString(cv)), nil
	}
	prefix := valueString(ctx, v)
	return newOf(cv, trimPrefixStr(getString(cv), prefix, ignoreCaseOf(cv))), nil
}

func methodTrimSuffix(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	v := argValue(ctx, 0)
	if av, ok := v.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z == nil || z.Value == nil {
				continue
			}
			t := trimSuffixStr(getString(cv), valueString(ctx, z.Value), ignoreCaseOf(cv))
			if t != getString(cv) {
				return newOf(cv, t), nil
			}
		}
		return newOf(cv, getString(cv)), nil
	}
	suffix := valueString(ctx, v)
	return newOf(cv, trimSuffixStr(getString(cv), suffix, ignoreCaseOf(cv))), nil
}

func methodReverse(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	return newOf(cv, reverseStr(getString(cv), strKindOf(cv))), nil
}

func methodRepeat(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	n := argInt(ctx, 0, 0)
	if n < 0 {
		return nil, throwInvalid("Multiplier must be positive, " + itoa(n) + " given.")
	}
	return newOf(cv, strings.Repeat(getString(cv), n)), nil
}

func methodPadBoth(ctx data.Context) (data.GetValue, data.Control) {
	return methodPad(ctx, 2)
}
func methodPadEnd(ctx data.Context) (data.GetValue, data.Control) {
	return methodPad(ctx, 1)
}
func methodPadStart(ctx data.Context) (data.GetValue, data.Control) {
	return methodPad(ctx, 0)
}

func methodPad(ctx data.Context, mode int) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	length := argInt(ctx, 0, 0)
	pad := argString(ctx, 1, " ")
	if pad == "" {
		return nil, throwInvalid("Invalid UTF-8 string.")
	}
	if strKindOf(cv) != kindByte {
		if ctl := mustUTF8(pad); ctl != nil {
			return nil, ctl
		}
	}
	return newOf(cv, padStr(getString(cv), pad, length, mode, strKindOf(cv))), nil
}

func methodJoin(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	arr := argValue(ctx, 0)
	var parts []string
	if av, ok := arr.(*data.ArrayValue); ok {
		for _, z := range av.List {
			if z != nil && z.Value != nil {
				parts = append(parts, valueString(ctx, z.Value))
			}
		}
	}
	var last *string
	if v := argValue(ctx, 1); v != nil {
		s := v.AsString()
		last = &s
	}
	out := joinStrings(getString(cv), parts, last)
	if strKindOf(cv) != kindByte {
		if ctl := mustUTF8(out); ctl != nil {
			return nil, ctl
		}
	}
	return newOf(cv, out), nil
}

func methodBytesAt(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	bs := bytesAt(getString(cv), argInt(ctx, 0, 0), strKindOf(cv))
	vals := make([]data.Value, len(bs))
	for i, b := range bs {
		vals[i] = data.NewIntValue(b)
	}
	return data.NewArrayValue(vals), nil
}

func methodCodePointsAt(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	ps := codePointsAt(getString(cv), argInt(ctx, 0, 0), strKindOf(cv))
	vals := make([]data.Value, len(ps))
	for i, p := range ps {
		vals[i] = data.NewIntValue(p)
	}
	return data.NewArrayValue(vals), nil
}

func methodToByteString(ctx data.Context) (data.GetValue, data.Control) {
	return newTyped(ctx, ensureByteClass(), getString(strSelf(ctx))), nil
}

func methodToCodePointString(ctx data.Context) (data.GetValue, data.Control) {
	s := getString(strSelf(ctx))
	if ctl := mustUTF8(s); ctl != nil {
		return nil, ctl
	}
	return newTyped(ctx, ensureCodePointClass(), s), nil
}

func methodToUnicodeString(ctx data.Context) (data.GetValue, data.Control) {
	s := getString(strSelf(ctx))
	if ctl := mustUTF8(s); ctl != nil {
		return nil, ctl
	}
	return newTyped(ctx, ensureUnicodeClass(), s), nil
}

func methodToStringAlias(ctx data.Context) (data.GetValue, data.Control) {
	return methodToString(ctx)
}

func methodIsUtf8(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(isValidUTF8(getString(strSelf(ctx)))), nil
}

func methodNormalize(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	form := argInt(ctx, 0, 4) // NFC
	switch form {
	case 1, 2, 4, 5: // NFD NFKD NFC NFKC
		return newOf(cv, getString(cv)), nil
	default:
		return nil, throwInvalid("Unsupported normalization form.")
	}
}

func methodWidth(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	ignoreAnsi := argBool(ctx, 0, true)
	s := getString(cv)
	if strKindOf(cv) == kindByte && !isValidUTF8(s) {
		var b strings.Builder
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				b.WriteByte('?')
			} else {
				b.WriteByte(s[i])
			}
		}
		s = b.String()
	}
	return data.NewIntValue(stringWidth(s, ignoreAnsi)), nil
}

func methodAscii(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	return newOf(cv, asciiTransliterate(getString(cv))), nil
}

func methodEnsureEnd(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	suffix := argString(ctx, 0, "")
	if endsWithStr(getString(cv), suffix, ignoreCaseOf(cv), strKindOf(cv)) {
		return newOf(cv, getString(cv)), nil
	}
	return newOf(cv, getString(cv)+suffix), nil
}

func methodEnsureStart(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	prefix := argString(ctx, 0, "")
	if startsWithStr(getString(cv), prefix, ignoreCaseOf(cv), strKindOf(cv)) {
		return newOf(cv, getString(cv)), nil
	}
	return newOf(cv, prefix+getString(cv)), nil
}

func methodAfter(ctx data.Context) (data.GetValue, data.Control) {
	return afterBefore(ctx, true, false)
}
func methodAfterLast(ctx data.Context) (data.GetValue, data.Control) {
	return afterBefore(ctx, true, true)
}
func methodBefore(ctx data.Context) (data.GetValue, data.Control) {
	return afterBefore(ctx, false, false)
}
func methodBeforeLast(ctx data.Context) (data.GetValue, data.Control) {
	return afterBefore(ctx, false, true)
}

func afterBefore(ctx data.Context, after, last bool) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	k := strKindOf(cv)
	include := argBool(ctx, 1, false)
	offset := argInt(ctx, 2, 0)
	s := getString(cv)
	best := -1
	bestNeedle := ""
	for _, n := range needles(ctx, argValue(ctx, 0)) {
		var i int
		var ok bool
		if last {
			i, ok = lastIndexOfUnits(s, n, offset, ignoreCaseOf(cv), k)
		} else {
			i, ok = indexOfUnits(s, n, offset, ignoreCaseOf(cv), k)
		}
		if !ok {
			continue
		}
		if last {
			if i >= best {
				best, bestNeedle = i, n
			}
		} else if best < 0 || i < best {
			best, bestNeedle = i, n
		}
	}
	if best < 0 {
		return newOf(cv, s), nil
	}
	if after {
		if !include {
			best += unitCount(bestNeedle, k)
		}
		return newOf(cv, sliceUnits(s, best, nil, k)), nil
	}
	if include {
		best += unitCount(bestNeedle, k)
	}
	return newOf(cv, sliceUnits(s, 0, intPtr(best), k)), nil
}

func methodTruncate(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	length := argInt(ctx, 0, 0)
	ellipsis := argString(ctx, 1, "")
	k := strKindOf(cv)
	s := getString(cv)
	if unitCount(s, k) <= length {
		return newOf(cv, s), nil
	}
	ellLen := 0
	if ellipsis != "" {
		ellLen = unitCount(ellipsis, k)
		if length < ellLen {
			ellLen = 0
		}
	}
	cut := sliceUnits(s, 0, intPtr(length-ellLen), k)
	if ellLen > 0 {
		cut = strings.TrimRight(cut, " \t\n\r") + ellipsis
	}
	return newOf(cv, cut), nil
}

func methodWordwrap(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	width := argInt(ctx, 0, 75)
	brk := argString(ctx, 1, "\n")
	cut := argBool(ctx, 2, false)
	k := strKindOf(cv)
	s := getString(cv)
	if width < 1 {
		return newOf(cv, s), nil
	}
	lines := []string{s}
	if brk != "" {
		lines = splitStr(s, brk, 2147483647, false, k)
	}
	if len(lines) == 1 && lines[0] == "" {
		return newOf(cv, ""), nil
	}
	var chars []string
	var mask strings.Builder
	for i, line := range lines {
		if i > 0 {
			chars = append(chars, brk)
			mask.WriteByte('#')
		}
		for _, ch := range chunkStr(line, 1, k) {
			chars = append(chars, ch)
			if ch == " " {
				mask.WriteByte(' ')
			} else {
				mask.WriteByte('?')
			}
		}
	}
	wrapped, ctl := callPHPFunc(ctx, "wordwrap", data.NewStringValue(mask.String()), data.NewIntValue(width), data.NewStringValue("#"), data.NewBoolValue(cut))
	if ctl != nil {
		return nil, ctl
	}
	maskOut := valueString(ctx, wrapped.(data.Value))
	var b strings.Builder
	j := 0
	i := -1
	pos := 0
	for {
		idx := strings.IndexByte(maskOut[pos:], '#')
		if idx < 0 {
			break
		}
		bpos := pos + idx
		for i = i + 1; i < bpos; i++ {
			if j < len(chars) {
				b.WriteString(chars[j])
				j++
			}
		}
		if j < len(chars) && (chars[j] == brk || chars[j] == " ") {
			j++
		}
		b.WriteString(brk)
		pos = bpos + 1
	}
	for ; j < len(chars); j++ {
		b.WriteString(chars[j])
	}
	return newOf(cv, b.String()), nil
}

func methodCollapseWhitespace(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	s := getString(cv)
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' || r == '\f' {
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return newOf(cv, strings.Trim(b.String(), " \n\r\t\f")), nil
}

func methodMatch(ctx data.Context) (data.GetValue, data.Control) {
	cv := strSelf(ctx)
	pattern := argString(ctx, 0, "")
	flags := argInt(ctx, 1, 0)
	offset := argInt(ctx, 2, 0)
	if ignoreCaseOf(cv) {
		pattern = withRegexFlag(pattern, 'i')
	}
	if strKindOf(cv) != kindByte {
		pattern = withRegexFlag(pattern, 'u')
	}
	re, err := preg.CompileAny(pattern)
	if err != nil {
		return nil, throwInvalid(err.Error())
	}
	subject := getString(cv)
	if flags&1 != 0 || flags&2 != 0 {
		all := re.FindAllStringSubmatchIndex(subject, -1)
		vals := make([]data.Value, 0, len(all))
		for _, loc := range all {
			if len(loc) < 2 || loc[0] < 0 {
				continue
			}
			vals = append(vals, data.NewStringValue(subject[loc[0]:loc[1]]))
		}
		return data.NewArrayValue(vals), nil
	}
	caps := preg.FindCaptures(re, subject, offset, preg.HasModifier(pattern, 'A'))
	if caps == nil {
		return data.NewArrayValue(nil), nil
	}
	return preg.BuildMatchArray(caps, flags|512), nil
}

func callPHPFunc(ctx data.Context, name string, args ...data.Value) (data.GetValue, data.Control) {
	fn, ok := ctx.GetVM().GetFunc(name)
	if !ok {
		return nil, data.NewErrorThrow(nil, errFunc(name))
	}
	callCtx := ctx.CreateContext(fn.GetVariables())
	data.BindDeclaredArgs(callCtx, fn, args)
	return fn.Call(callCtx)
}

func errFunc(name string) error { return &simpleErr{s: "function " + name + " not found"} }

type simpleErr struct{ s string }

func (e *simpleErr) Error() string { return e.s }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func replacePairs(s string, from, to []string) string {
	if len(from) == 0 {
		return s
	}
	repl := make([]string, 0, len(from)*2)
	n := len(from)
	if len(to) < n {
		n = len(to)
	}
	for i := 0; i < n; i++ {
		repl = append(repl, from[i], to[i])
	}
	return strings.NewReplacer(repl...).Replace(s)
}
