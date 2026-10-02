package preg

import "strings"

// PCRE's horizontal/vertical whitespace classes differ from Go and .NET.
// Translate only unquoted escapes, including their use within character classes.
func translateWhitespaceClasses(pattern string, utf bool) string {
	if !strings.ContainsAny(pattern, "hHvV") || !strings.Contains(pattern, `\`) {
		return pattern
	}
	var out strings.Builder
	inClass, quoted := false, false
	changed := false
	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		if char == '\\' && i+1 < len(pattern) {
			next := pattern[i+1]
			if next == 'Q' && !quoted {
				quoted = true
			}
			if next == 'E' && quoted {
				quoted = false
			}
			if !quoted && (next == 'h' || next == 'H' || next == 'v' || next == 'V') {
				points := []rune{9, 32, 160}
				if next == 'v' || next == 'V' {
					points = []rune{10, 11, 12, 13, 133}
				}
				if utf {
					if next == 'h' || next == 'H' {
						points = []rune{9, 32, 160, 5760, 6158, 8192, 8193, 8194, 8195, 8196, 8197, 8198, 8199, 8200, 8201, 8202, 8239, 8287, 12288}
					} else {
						points = []rune{10, 11, 12, 13, 133, 8232, 8233}
					}
				}
				negative := next == 'H' || next == 'V'
				if !inClass {
					out.WriteByte('[')
					if negative {
						out.WriteByte('^')
					}
				}
				if inClass && negative {
					previous := rune(0)
					for _, point := range points {
						writeRuneRange(&out, previous, point-1)
						previous = point + 1
					}
					writeRuneRange(&out, previous, 0x10ffff)
				} else {
					for _, point := range points {
						out.WriteRune(point)
					}
				}
				if !inClass {
					out.WriteByte(']')
				}
				changed = true
				i++
				continue
			}
			out.WriteByte(char)
			out.WriteByte(next)
			i++
			continue
		}
		if !quoted {
			if char == '[' {
				inClass = true
			} else if char == ']' {
				inClass = false
			}
		}
		out.WriteByte(char)
	}
	if !changed {
		return pattern
	}
	return out.String()
}

func writeRuneRange(out *strings.Builder, start, end rune) {
	if start > end {
		return
	}
	out.WriteRune(start)
	if start != end {
		out.WriteByte('-')
		out.WriteRune(end)
	}
}
