package support

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

var strSnakeUpperRe = regexp.MustCompile(`(.)([A-Z])`)

var (
	strSnakeCache  = map[string]map[string]string{}
	strCamelCache  = map[string]string{}
	strStudlyCache = map[string]string{}
	strCacheMu     sync.RWMutex

	strRandomFactory data.Value
	strUUIDFactory   data.Value
	strULIDFactory   data.Value
)

func strIntArg(ctx data.Context, i int, def int) int {
	v := kit.Arg(ctx, i)
	if v == nil || kit.IsNull(v) {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func strBoolArg(ctx data.Context, i int, def bool) bool {
	v := kit.Arg(ctx, i)
	if v == nil || kit.IsNull(v) {
		return def
	}
	if bv, ok := v.(data.AsBool); ok {
		b, err := bv.AsBool()
		if err == nil {
			return b
		}
	}
	return kit.Truthy(v)
}

func mbLen(s string) int {
	return len([]rune(s))
}

func mbSubstr(s string, start, length int) string {
	r := []rune(s)
	if start < 0 {
		start = len(r) + start
	}
	if start < 0 {
		start = 0
	}
	if start > len(r) {
		return ""
	}
	if length < 0 {
		length = len(r) - start + length
	}
	end := start + length
	if end > len(r) {
		end = len(r)
	}
	if end < start {
		return ""
	}
	return string(r[start:end])
}

func mbStrrpos(s, sub string) int {
	r := []rune(s)
	rs := []rune(sub)
	if len(rs) == 0 {
		return -1
	}
	for i := len(r) - len(rs); i >= 0; i-- {
		ok := true
		for j := 0; j < len(rs); j++ {
			if r[i+j] != rs[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

func mbUcfirst(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func mbLcfirst(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func mbReverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func studlyCached(value string, normalize bool) string {
	if normalize {
		value = strNormalizeAcronyms(value)
	}
	strCacheMu.RLock()
	if out, ok := strStudlyCache[value]; ok {
		strCacheMu.RUnlock()
		return out
	}
	strCacheMu.RUnlock()

	words := strings.Fields(strings.NewReplacer("-", " ", "_", " ").Replace(value))
	for i, w := range words {
		words[i] = mbUcfirst(w)
	}
	out := strings.Join(words, "")

	strCacheMu.Lock()
	strStudlyCache[value] = out
	strCacheMu.Unlock()
	return out
}

func strNormalizeAcronyms(value string) string {
	var b strings.Builder
	runes := []rune(value)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if unicode.IsUpper(r) {
			start := i
			for i+1 < len(runes) && unicode.IsUpper(runes[i+1]) {
				i++
			}
			word := string(runes[start : i+1])
			if i+1 >= len(runes) || runes[i+1] == '-' || runes[i+1] == '_' || runes[i+1] == ' ' {
				word = strings.ToLower(word)
			}
			b.WriteString(word)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func snakeCached(value, delim string) string {
	strCacheMu.RLock()
	if m, ok := strSnakeCache[value]; ok {
		if out, ok2 := m[delim]; ok2 {
			strCacheMu.RUnlock()
			return out
		}
	}
	strCacheMu.RUnlock()

	out := value
	if !isASCIILower(value) {
		out = strings.ReplaceAll(out, " ", "")
		out = strSnakeUpperRe.ReplaceAllString(out, "${1}"+delim+"${2}")
		out = strings.ToLower(out)
		out = strings.ReplaceAll(out, "-", delim)
	}

	strCacheMu.Lock()
	if strSnakeCache[value] == nil {
		strSnakeCache[value] = map[string]string{}
	}
	strSnakeCache[value][delim] = out
	strCacheMu.Unlock()
	return out
}

func isASCIILower(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func randomStringLaravel(length int) string {
	var sb strings.Builder
	for sb.Len() < length {
		need := length - sb.Len()
		bytesSize := ((need + 2) / 3) * 3
		b := make([]byte, bytesSize)
		_, _ = rand.Read(b)
		chunk := strings.NewReplacer("/", "", "+", "", "=", "").Replace(base64.StdEncoding.EncodeToString(b))
		if len(chunk) > need {
			chunk = chunk[:need]
		}
		sb.WriteString(chunk)
	}
	return sb.String()
}

func formatUUIDv4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

func formatUUIDv7(ms uint64) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b[6:])
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	b[6] = (b[6] & 0x0f) | 0x70
	b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

var (
	ulidGenMu     sync.Mutex
	ulidLastMs    uint64
	ulidLastRand  [10]byte
)

func generateULIDString(ms uint64) string {
	if ms == 0 {
		ms = uint64(time.Now().UnixMilli())
	}
	ulidGenMu.Lock()
	defer ulidGenMu.Unlock()
	var randPart [10]byte
	if ms == ulidLastMs {
		randPart = ulidLastRand
		overflow := true
		for i := 9; i >= 0; i-- {
			randPart[i]++
			if randPart[i] != 0 {
				overflow = false
				break
			}
		}
		if overflow {
			ms++
		}
	} else {
		_, _ = rand.Read(randPart[:])
	}
	ulidLastMs = ms
	ulidLastRand = randPart

	var ts [6]byte
	ts[0] = byte(ms >> 40)
	ts[1] = byte(ms >> 32)
	ts[2] = byte(ms >> 24)
	ts[3] = byte(ms >> 16)
	ts[4] = byte(ms >> 8)
	ts[5] = byte(ms)
	payload := append(ts[:], randPart[:]...)
	return encodeCrockford26(payload)
}

func encodeCrockford26(b []byte) string {
	if len(b) != 16 {
		return ""
	}
	var out [26]byte
	var acc uint64
	nbits := 2
	oi := 0
	for _, by := range b {
		acc = (acc << 8) | uint64(by)
		nbits += 8
		for nbits >= 5 {
			nbits -= 5
			out[oi] = crockford[(acc>>nbits)&31]
			oi++
		}
	}
	return string(out[:])
}

func ulidValid(s string) bool {
	if len(s) != 26 {
		return false
	}
	u := strings.ToUpper(s)
	if u[0] > '7' {
		return false
	}
	for i := 0; i < 26; i++ {
		c := u[i]
		if strings.IndexByte(crockford, c) < 0 {
			return false
		}
	}
	return true
}

// asciiFallback: voku/ASCII language packs 未接入；仅做 Unicode 字母数字保留与常见拉丁变体折叠。
func asciiFallback(value, _language string) string {
	var b strings.Builder
	for _, r := range value {
		switch r {
		case 'à', 'á', 'â', 'ã', 'ä', 'å':
			b.WriteByte('a')
		case 'À', 'Á', 'Â', 'Ã', 'Ä', 'Å':
			b.WriteByte('A')
		case 'æ':
			b.WriteString("ae")
		case 'Æ':
			b.WriteString("AE")
		case 'ç':
			b.WriteByte('c')
		case 'Ç':
			b.WriteByte('C')
		case 'è', 'é', 'ê', 'ë':
			b.WriteByte('e')
		case 'È', 'É', 'Ê', 'Ë':
			b.WriteByte('E')
		case 'ì', 'í', 'î', 'ï':
			b.WriteByte('i')
		case 'ñ':
			b.WriteByte('n')
		case 'ò', 'ó', 'ô', 'õ', 'ö':
			b.WriteByte('o')
		case 'ù', 'ú', 'û', 'ü':
			b.WriteByte('u')
		case 'ß':
			b.WriteString("ss")
		default:
			if r < 128 {
				b.WriteRune(r)
			} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func transliterateFallback(s, unknown string, strict bool) string {
	var b strings.Builder
	for _, r := range s {
		if r < 128 && r >= 32 {
			b.WriteRune(r)
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		if !strict {
			b.WriteString(unknown)
		}
	}
	return b.String()
}

// englishPluralFallback: Laravel 使用 Doctrine Inflector；此处为常见规则近似。
func englishPluralFallback(word string, count int) string {
	if count == 1 {
		return englishSingularFallback(word)
	}
	lower := strings.ToLower(word)
	switch {
	case strings.HasSuffix(lower, "y") && len(word) > 1:
		return word[:len(word)-1] + "ies"
	case strings.HasSuffix(lower, "s"), strings.HasSuffix(lower, "x"), strings.HasSuffix(lower, "ch"), strings.HasSuffix(lower, "sh"):
		return word + "es"
	default:
		return word + "s"
	}
}

func englishSingularFallback(word string) string {
	lower := strings.ToLower(word)
	switch {
	case strings.HasSuffix(lower, "ies") && len(word) > 3:
		return word[:len(word)-3] + "y"
	case strings.HasSuffix(lower, "ses") && len(word) > 3:
		return word[:len(word)-2]
	case strings.HasSuffix(lower, "s") && len(word) > 1:
		return word[:len(word)-1]
	default:
		return word
	}
}

func strValueString(v data.GetValue) string {
	if v == nil {
		return ""
	}
	if val, ok := v.(data.Value); ok {
		return val.AsString()
	}
	return ""
}

func arrayAppend(av *data.ArrayValue, v data.Value) {
	if av == nil {
		return
	}
	av.List = append(av.List, data.NewZVal(v))
}

func strFlushCaches() {
	strCacheMu.Lock()
	strSnakeCache = map[string]map[string]string{}
	strCamelCache = map[string]string{}
	strStudlyCache = map[string]string{}
	strCacheMu.Unlock()
}
