package sfstring

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func asciiTransliterate(s string) string {
	s = replacePairs(s, translitFrom, translitTo)
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r < 0x80 {
			b.WriteRune(r)
			continue
		}
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if a, ok := latinASCII[r]; ok {
			b.WriteString(a)
			continue
		}
		// 预组合拉丁字母粗略去音调
		if mapped := stripLatinAccent(r); mapped != 0 {
			b.WriteRune(mapped)
			continue
		}
		b.WriteByte('?')
	}
	return b.String()
}

func stripLatinAccent(r rune) rune {
	switch {
	case r >= 'À' && r <= 'Å':
		return 'A'
	case r == 'Ç':
		return 'C'
	case r >= 'È' && r <= 'Ë':
		return 'E'
	case r >= 'Ì' && r <= 'Ï':
		return 'I'
	case r == 'Ñ':
		return 'N'
	case r >= 'Ò' && r <= 'Ö':
		return 'O'
	case r >= 'Ù' && r <= 'Ü':
		return 'U'
	case r == 'Ý':
		return 'Y'
	case r >= 'à' && r <= 'å':
		return 'a'
	case r == 'ç':
		return 'c'
	case r >= 'è' && r <= 'ë':
		return 'e'
	case r >= 'ì' && r <= 'ï':
		return 'i'
	case r == 'ñ':
		return 'n'
	case r >= 'ò' && r <= 'ö':
		return 'o'
	case r >= 'ù' && r <= 'ü':
		return 'u'
	case r == 'ý' || r == 'ÿ':
		return 'y'
	}
	return 0
}

var latinASCII = map[rune]string{
	'Æ': "AE", 'æ': "ae", 'Œ': "OE", 'œ': "oe",
	'ß': "ss", 'Ø': "O", 'ø': "o", 'Đ': "D", 'đ': "d",
	'Ł': "L", 'ł': "l", 'Þ': "TH", 'þ': "th",
}

func stringWidth(s string, ignoreAnsi bool) int {
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.ReplaceAll(s, "\x05", "")
	s = strings.ReplaceAll(s, "\x07", "")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	width := 0
	for _, line := range strings.Split(s, "\n") {
		if ignoreAnsi {
			line = stripANSI(line)
		}
		w := wcswidth(line)
		if w > width {
			width = w
		}
	}
	return width
}

func stripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == 0x1b {
			i++
			if i < len(s) && s[i] == '[' {
				i++
				for i < len(s) {
					c := s[i]
					i++
					if c >= 0x40 && c <= 0x7e {
						break
					}
				}
				continue
			}
			if i < len(s) && s[i] >= 0x40 && s[i] <= 0x5f {
				i++
				continue
			}
			continue
		}
		if s[i] < 0x20 || s[i] == 0x7f {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			i++
			continue
		}
		b.WriteString(s[i : i+size])
		i += size
	}
	return b.String()
}

func wcswidth(s string) int {
	width := 0
	for _, r := range s {
		if r == 0 || r == 0x034F || (r >= 0x200B && r <= 0x200F) || r == 0x2028 || r == 0x2029 ||
			(r >= 0x202A && r <= 0x202E) || (r >= 0x2060 && r <= 0x2063) {
			continue
		}
		if r < 32 || (r >= 0x7F && r < 0xA0) {
			return -1
		}
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) {
			continue
		}
		if eastAsianWide(r) {
			width += 2
			continue
		}
		width++
	}
	return width
}

func eastAsianWide(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x115F:
		return true
	case r >= 0x231A && r <= 0x231B:
		return true
	case r >= 0x2329 && r <= 0x232A:
		return true
	case r >= 0x2E80 && r <= 0xA4CF && r != 0x303F:
		return true
	case r >= 0xAC00 && r <= 0xD7A3:
		return true
	case r >= 0xF900 && r <= 0xFAFF:
		return true
	case r >= 0xFE10 && r <= 0xFE19:
		return true
	case r >= 0xFE30 && r <= 0xFE6F:
		return true
	case r >= 0xFF01 && r <= 0xFF60:
		return true
	case r >= 0xFFE0 && r <= 0xFFE6:
		return true
	case r >= 0x1F300 && r <= 0x1F64F:
		return true
	case r >= 0x1F900 && r <= 0x1F9FF:
		return true
	case r >= 0x20000 && r <= 0x3FFFD:
		return true
	}
	return false
}
