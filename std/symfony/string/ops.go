package sfstring

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func indexOfUnits(haystack, needle string, offset int, ignoreCase bool, k strKind) (int, bool) {
	if needle == "" {
		return 0, false
	}
	h := units(haystack, k)
	n := units(needle, k)
	if len(n) == 0 {
		return 0, false
	}
	if offset < 0 {
		offset = len(h) + offset
	}
	if offset < 0 {
		offset = 0
	}
	if offset > len(h) {
		return 0, false
	}
	for i := offset; i+len(n) <= len(h); i++ {
		if unitsEqual(h[i:i+len(n)], n, ignoreCase) {
			return i, true
		}
	}
	return 0, false
}

func lastIndexOfUnits(haystack, needle string, offset int, ignoreCase bool, k strKind) (int, bool) {
	if needle == "" {
		return 0, false
	}
	h := units(haystack, k)
	n := units(needle, k)
	if len(n) == 0 {
		return 0, false
	}
	end := len(h)
	if offset < 0 {
		end = len(h) + offset + len(n)
		if end < 0 {
			return 0, false
		}
		if end > len(h) {
			end = len(h)
		}
		h = h[:end]
		offset = 0
	}
	found := -1
	start := offset
	if start < 0 {
		start = 0
	}
	for i := start; i+len(n) <= len(h); i++ {
		if unitsEqual(h[i:i+len(n)], n, ignoreCase) {
			found = i
		}
	}
	if found < 0 {
		return 0, false
	}
	return found, true
}

func unitsEqual(a, b []string, ignoreCase bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if ignoreCase {
			if !strings.EqualFold(a[i], b[i]) {
				return false
			}
		} else if a[i] != b[i] {
			return false
		}
	}
	return true
}

func startsWithStr(s, prefix string, ignoreCase bool, k strKind) bool {
	if prefix == "" {
		return false
	}
	if k == kindByte {
		if len(s) < len(prefix) {
			return false
		}
		if ignoreCase {
			return strings.EqualFold(s[:len(prefix)], prefix)
		}
		return strings.HasPrefix(s, prefix)
	}
	n := unitCount(prefix, k)
	head := prefixUnits(s, n, k)
	if ignoreCase {
		return strings.EqualFold(head, prefix)
	}
	return head == prefix
}

func endsWithStr(s, suffix string, ignoreCase bool, k strKind) bool {
	if suffix == "" {
		return false
	}
	if k == kindByte {
		if len(s) < len(suffix) {
			return false
		}
		if ignoreCase {
			return strings.EqualFold(s[len(s)-len(suffix):], suffix)
		}
		return strings.HasSuffix(s, suffix)
	}
	n := unitCount(suffix, k)
	u := units(s, k)
	if len(u) < n {
		return false
	}
	tail := strings.Join(u[len(u)-n:], "")
	if ignoreCase {
		return strings.EqualFold(tail, suffix)
	}
	return tail == suffix
}

func equalsStr(a, b string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.EqualFold(a, b)
	}
	return a == b
}

func replaceStr(s, from, to string, ignoreCase bool, k strKind) string {
	if from == "" {
		return s
	}
	if !ignoreCase && k == kindByte {
		return strings.ReplaceAll(s, from, to)
	}
	var b strings.Builder
	rest := s
	for {
		i, ok := indexOfUnits(rest, from, 0, ignoreCase, k)
		if !ok {
			b.WriteString(rest)
			break
		}
		head := sliceUnits(rest, 0, intPtr(i), k)
		b.WriteString(head)
		b.WriteString(to)
		fromLen := unitCount(from, k)
		rest = sliceUnits(rest, i+fromLen, nil, k)
	}
	return b.String()
}

func intPtr(n int) *int { return &n }

func splitStr(s, delim string, limit int, ignoreCase bool, k strKind) []string {
	if limit < 1 {
		limit = 1
	}
	if delim == "" {
		return []string{s}
	}
	out := make([]string, 0, 4)
	rest := s
	for len(out)+1 < limit {
		i, ok := indexOfUnits(rest, delim, 0, ignoreCase, k)
		if !ok {
			break
		}
		out = append(out, sliceUnits(rest, 0, intPtr(i), k))
		rest = sliceUnits(rest, i+unitCount(delim, k), nil, k)
	}
	out = append(out, rest)
	return out
}

func chunkStr(s string, length int, k strKind) []string {
	u := units(s, k)
	if len(u) == 0 {
		return nil
	}
	out := make([]string, 0, (len(u)+length-1)/length)
	for i := 0; i < len(u); i += length {
		end := i + length
		if end > len(u) {
			end = len(u)
		}
		out = append(out, strings.Join(u[i:end], ""))
	}
	return out
}

func spliceStr(s, replacement string, start int, length *int, k strKind) string {
	u := units(s, k)
	n := len(u)
	if start < 0 {
		start = n + start
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	end := n
	if length != nil {
		end = start + *length
		if *length < 0 {
			end = n + *length
		}
		if end < start {
			end = start
		}
		if end > n {
			end = n
		}
	}
	return strings.Join(u[:start], "") + replacement + strings.Join(u[end:], "")
}

func trimChars(s, chars string, left, right bool) string {
	if s == "" {
		return s
	}
	set := make(map[rune]struct{}, len(chars))
	for _, r := range chars {
		set[r] = struct{}{}
	}
	rs := []rune(s)
	start, end := 0, len(rs)
	if left {
		for start < end {
			if _, ok := set[rs[start]]; !ok {
				break
			}
			start++
		}
	}
	if right {
		for end > start {
			if _, ok := set[rs[end-1]]; !ok {
				break
			}
			end--
		}
	}
	return string(rs[start:end])
}

func trimPrefixStr(s, prefix string, ignoreCase bool) string {
	if prefix == "" || len(s) < len(prefix) {
		return s
	}
	head := s[:len(prefix)]
	if ignoreCase {
		if strings.EqualFold(head, prefix) {
			return s[len(prefix):]
		}
		return s
	}
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func trimSuffixStr(s, suffix string, ignoreCase bool) string {
	if suffix == "" || len(s) < len(suffix) {
		return s
	}
	tail := s[len(s)-len(suffix):]
	if ignoreCase {
		if strings.EqualFold(tail, suffix) {
			return s[:len(s)-len(suffix)]
		}
		return s
	}
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func reverseStr(s string, k strKind) string {
	u := units(s, k)
	for i, j := 0, len(u)-1; i < j; i, j = i+1, j-1 {
		u[i], u[j] = u[j], u[i]
	}
	return strings.Join(u, "")
}

func titleStr(s string, allWords bool, k strKind) string {
	if s == "" {
		return s
	}
	if k == kindByte {
		if allWords {
			return strings.Title(s)
		}
		return strings.ToUpper(s[:1]) + s[1:]
	}
	rs := []rune(s)
	limit := 1
	if allWords {
		limit = -1
	}
	changed := 0
	prevSpace := true
	for i, r := range rs {
		if unicode.IsSpace(r) || (!unicode.IsLetter(r) && !unicode.IsDigit(r)) {
			prevSpace = true
			continue
		}
		if prevSpace {
			if limit != -1 && changed >= limit {
				break
			}
			rs[i] = unicode.ToTitle(r)
			changed++
			prevSpace = false
		} else {
			prevSpace = false
		}
	}
	return string(rs)
}

func camelUnicode(s string) string {
	var b strings.Builder
	word := 0
	start := true
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if start {
				if word == 0 {
					b.WriteRune(unicode.ToLower(r))
				} else {
					b.WriteRune(unicode.ToTitle(r))
				}
				start = false
				word++
			} else {
				b.WriteRune(r)
			}
		} else {
			start = true
		}
	}
	return b.String()
}

func camelByte(s string) string {
	cleaned := make([]rune, 0, len(s))
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r >= 0x7f {
			cleaned = append(cleaned, r)
		} else {
			cleaned = append(cleaned, ' ')
		}
	}
	parts := strings.Fields(strings.TrimSpace(string(cleaned)))
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			if len(p) != 1 && isAllUpperASCII(p) {
				parts[i] = p
			} else {
				parts[i] = strings.ToLower(p[:1]) + p[1:]
			}
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func isAllUpperASCII(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func snakeFromCamel(s string) string {
	var b strings.Builder
	rs := []rune(s)
	for i, r := range rs {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rs[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					b.WriteByte('_')
				} else if i+1 < len(rs) && unicode.IsLower(rs[i+1]) {
					b.WriteByte('_')
				}
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func padStr(s, pad string, length int, mode int, k strKind) string {
	cur := unitCount(s, k)
	if length <= cur || pad == "" {
		return s
	}
	free := length - cur
	switch mode {
	case 0: // STR_PAD_LEFT
		return repeatPad(pad, free, k) + s
	case 2: // STR_PAD_BOTH
		left := free / 2
		right := free - left
		return repeatPad(pad, left, k) + s + repeatPad(pad, right, k)
	default: // STR_PAD_RIGHT
		return s + repeatPad(pad, free, k)
	}
}

func repeatPad(pad string, n int, k strKind) string {
	if n <= 0 || pad == "" {
		return ""
	}
	u := units(pad, k)
	if len(u) == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(u[i%len(u)])
	}
	return b.String()
}

func bytesAt(s string, offset int, k strKind) []int {
	piece := sliceUnits(s, offset, intPtr(1), k)
	if piece == "" {
		return nil
	}
	out := make([]int, 0, len(piece))
	for i := 0; i < len(piece); i++ {
		out = append(out, int(piece[i]))
	}
	return out
}

func codePointsAt(s string, offset int, k strKind) []int {
	piece := sliceUnits(s, offset, intPtr(1), k)
	if piece == "" {
		return nil
	}
	out := make([]int, 0, utf8.RuneCountInString(piece))
	for _, r := range piece {
		out = append(out, int(r))
	}
	return out
}

func joinStrings(glue string, parts []string, lastGlue *string) string {
	if lastGlue != nil && len(parts) > 1 {
		head := strings.Join(parts[:len(parts)-1], glue)
		return head + *lastGlue + parts[len(parts)-1]
	}
	return strings.Join(parts, glue)
}

func isValidUTF8(s string) bool {
	return s == "" || utf8.ValidString(s)
}
