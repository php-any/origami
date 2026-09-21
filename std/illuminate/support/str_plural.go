package support

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/data"
)

var englishIrregularPlural = map[string]string{
	"child": "children", "person": "people", "man": "men", "woman": "women",
	"mouse": "mice", "louse": "lice", "goose": "geese", "foot": "feet",
	"tooth": "teeth", "ox": "oxen", "die": "dice", "index": "indices",
	"matrix": "matrices", "vertex": "vertices", "axis": "axes",
	"analysis": "analyses", "basis": "bases", "crisis": "crises",
	"thesis": "theses", "quiz": "quizzes", "wife": "wives", "knife": "knives",
	"life": "lives", "leaf": "leaves", "half": "halves", "self": "selves",
	"elf": "elves", "loaf": "loaves", "potato": "potatoes", "tomato": "tomatoes",
	"hero": "heroes", "echo": "echoes", "embargo": "embargoes",
}

var englishUncountable = map[string]struct{}{
	"sheep": {}, "fish": {}, "deer": {}, "series": {}, "species": {},
	"money": {}, "information": {}, "equipment": {}, "rice": {}, "news": {},
	"recommended": {}, "related": {},
}

func englishPlural(word string) string {
	if word == "" {
		return word
	}
	lower := strings.ToLower(word)
	if _, ok := englishUncountable[lower]; ok {
		return word
	}
	if p, ok := englishIrregularPlural[lower]; ok {
		return matchEnglishCase(p, word)
	}
	runes := []rune(lower)
	n := len(runes)
	plural := lower + "s"
	if n >= 1 {
		last := runes[n-1]
		switch {
		case last == 'y' && n >= 2 && !isEnglishVowel(runes[n-2]):
			plural = string(runes[:n-1]) + "ies"
		case last == 's' || last == 'x' || last == 'z' || last == 'o':
			plural = lower + "es"
		case n >= 2 && ((runes[n-2] == 'c' && last == 'h') || (runes[n-2] == 's' && last == 'h')):
			plural = lower + "es"
		case last == 'f':
			plural = string(runes[:n-1]) + "ves"
		case n >= 2 && runes[n-2] == 'f' && last == 'e':
			plural = string(runes[:n-2]) + "ves"
		}
	}
	return matchEnglishCase(plural, word)
}

func englishSingular(word string) string {
	if word == "" {
		return word
	}
	lower := strings.ToLower(word)
	if _, ok := englishUncountable[lower]; ok {
		return word
	}
	for s, p := range englishIrregularPlural {
		if p == lower {
			return matchEnglishCase(s, word)
		}
	}
	switch {
	case strings.HasSuffix(lower, "ies") && len(lower) > 3:
		return matchEnglishCase(lower[:len(lower)-3]+"y", word)
	case strings.HasSuffix(lower, "ves") && len(lower) > 3:
		return matchEnglishCase(lower[:len(lower)-3]+"f", word)
	case strings.HasSuffix(lower, "ses") || strings.HasSuffix(lower, "xes") || strings.HasSuffix(lower, "zes") ||
		strings.HasSuffix(lower, "ches") || strings.HasSuffix(lower, "shes"):
		return matchEnglishCase(lower[:len(lower)-2], word)
	case strings.HasSuffix(lower, "oes") && len(lower) > 3:
		return matchEnglishCase(lower[:len(lower)-2], word)
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") && len(lower) > 1:
		return matchEnglishCase(lower[:len(lower)-1], word)
	}
	return word
}

func isEnglishVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func matchEnglishCase(out, sample string) string {
	if sample == "" {
		return out
	}
	r, _ := utf8.DecodeRuneInString(sample)
	if unicode.IsUpper(r) && strings.ToUpper(sample) == sample {
		return strings.ToUpper(out)
	}
	if unicode.IsUpper(r) {
		or, size := utf8.DecodeRuneInString(out)
		return string(unicode.ToUpper(or)) + out[size:]
	}
	return out
}

func strCountable(ctx data.Context, i int, def int) int {
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil || isNull(v) {
		return def
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return len(toEntries(av))
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func strPlural(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	count := strCountable(ctx, 1, 2)
	prepend := false
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			prepend, _ = b.AsBool()
		}
	}
	out := value
	if absInt(count) != 1 {
		out = englishPlural(value)
	}
	if prepend {
		out = strconv.Itoa(count) + " " + out
	}
	return data.NewStringValue(out), nil
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func strSingular(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(englishSingular(strArg(ctx, 0))), nil
}

func strPluralStudly(ctx data.Context) (data.GetValue, data.Control) {
	value := strArg(ctx, 0)
	count := strCountable(ctx, 1, 2)
	if value == "" {
		return data.NewStringValue(""), nil
	}
	i := len(value)
	for i > 0 {
		r, size := utf8.DecodeLastRuneInString(value[:i])
		if i < len(value) && unicode.IsUpper(r) {
			break
		}
		i -= size
	}
	head, last := value[:i], value[i:]
	if last == "" {
		last = value
		head = ""
	}
	out := last
	if absInt(count) != 1 {
		out = englishPlural(last)
	}
	return data.NewStringValue(head + out), nil
}
