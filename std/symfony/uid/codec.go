package uid

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/php-any/origami/data"
)

const (
	abstractUidName = "Symfony\\Component\\Uid\\AbstractUid"
	uuidName        = "Symfony\\Component\\Uid\\Uuid"
	ulidName        = "Symfony\\Component\\Uid\\Ulid"
	uuidV4Name      = "Symfony\\Component\\Uid\\UuidV4"
	uuidV7Name      = "Symfony\\Component\\Uid\\UuidV7"
	nilUuidName     = "Symfony\\Component\\Uid\\NilUuid"
	maxUuidName     = "Symfony\\Component\\Uid\\MaxUuid"
	nilUlidName     = "Symfony\\Component\\Uid\\NilUlid"
	maxUlidName     = "Symfony\\Component\\Uid\\MaxUlid"
	hashableName    = "Symfony\\Component\\Uid\\HashableInterface"
	timeBasedName   = "Symfony\\Component\\Uid\\TimeBasedUidInterface"
	invalidArgEx    = "Symfony\\Component\\Uid\\Exception\\InvalidArgumentException"

	uuidNil = "00000000-0000-0000-0000-000000000000"
	uuidMax = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	ulidNil = "00000000000000000000000000"
	ulidMax = "7ZZZZZZZZZZZZZZZZZZZZZZZZZ"

	crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
)

var crockfordDecode [256]int

func init() {
	for i := range crockfordDecode {
		crockfordDecode[i] = -1
	}
	for i, c := range []byte(crockfordAlphabet) {
		crockfordDecode[c] = i
		if c >= 'A' && c <= 'Z' {
			crockfordDecode[c+('a'-'A')] = i
		}
	}
}

func formatRFC4122(b []byte) string {
	h := hex.EncodeToString(b)
	if len(h) != 32 {
		return ""
	}
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

func isRFC4122(s string) bool {
	if len(s) != 36 {
		return false
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return false
	}
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		c := s[i]
		if c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' {
			continue
		}
		return false
	}
	return true
}

func rfc4122Bytes(s string) ([]byte, bool) {
	if !isRFC4122(s) {
		return nil, false
	}
	b, err := hex.DecodeString(s[0:8] + s[9:13] + s[14:18] + s[19:23] + s[24:36])
	if err != nil || len(b) != 16 {
		return nil, false
	}
	return b, true
}

func generateUUIDv4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return formatRFC4122(b)
}

func generateUUIDv7(ms uint64) string {
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
	return formatRFC4122(b)
}

func encode128ToCrockford(b []byte) string {
	if len(b) != 16 {
		return ""
	}
	var out [26]byte
	var acc uint64
	nbits := 2 // 130 = 26*5，前置 2 bit 填 0，使首字符落在 0-7
	oi := 0
	for _, by := range b {
		acc = (acc << 8) | uint64(by)
		nbits += 8
		for nbits >= 5 {
			nbits -= 5
			out[oi] = crockfordAlphabet[(acc>>nbits)&31]
			oi++
		}
	}
	return string(out[:])
}

func decodeCrockfordTo16(s string) ([16]byte, bool) {
	var out [16]byte
	if len(s) != 26 {
		return out, false
	}
	var acc uint64
	nbits := 0
	skip := 2
	oi := 0
	for i := 0; i < 26; i++ {
		v := crockfordDecode[s[i]]
		if v < 0 {
			return out, false
		}
		acc = (acc << 5) | uint64(v)
		nbits += 5
		if skip > 0 {
			if nbits < skip {
				continue
			}
			nbits -= skip
			if nbits == 0 {
				acc = 0
			} else {
				acc &= (uint64(1) << nbits) - 1
			}
			skip = 0
		}
		for nbits >= 8 && oi < 16 {
			nbits -= 8
			out[oi] = byte(acc >> nbits)
			acc &= (uint64(1) << nbits) - 1
			oi++
		}
	}
	return out, oi == 16
}

func ulidIsValid(s string) bool {
	if len(s) != 26 {
		return false
	}
	u := strings.ToUpper(s)
	if u[0] > '7' {
		return false
	}
	for i := 0; i < 26; i++ {
		c := u[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		if crockfordDecode[c] < 0 {
			return false
		}
		// 严格字母表：不含 I/L/O/U
		switch c {
		case 'I', 'L', 'O', 'U':
			return false
		}
	}
	return true
}

var (
	ulidMu      sync.Mutex
	ulidLastMs  = ^uint64(0)
	ulidLastEnt [10]byte
)

func generateULID(ms uint64) string {
	ulidMu.Lock()
	defer ulidMu.Unlock()
	var ent [10]byte
	if ms == ulidLastMs {
		ent = ulidLastEnt
		overflow := true
		for i := 9; i >= 0; i-- {
			ent[i]++
			if ent[i] != 0 {
				overflow = false
				break
			}
		}
		if overflow {
			ms++
			_, _ = rand.Read(ent[:])
		}
	} else {
		_, _ = rand.Read(ent[:])
	}
	ulidLastMs = ms
	ulidLastEnt = ent
	var b [16]byte
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	copy(b[6:], ent[:])
	return encode128ToCrockford(b[:])
}

func optionalClass(ctx data.Context, name string) data.ClassStmt {
	vm := ctx.GetVM()
	if vm == nil {
		return nil
	}
	if c, ok := vm.GetClass(name); ok && c != nil {
		return c
	}
	if _, ctl := vm.LoadPkg(name); ctl != nil {
		return nil
	}
	if c, ok := vm.GetClass(name); ok && c != nil {
		return c
	}
	return nil
}

func newClassWithStringArg(ctx data.Context, stmt data.ClassStmt, arg string) (*data.ClassValue, data.Control) {
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	ctor := stmt.GetConstruct()
	if ctor == nil {
		_ = cv.SetProperty("uid", data.NewStringValue(arg))
		return cv, nil
	}
	vars := ctor.GetVariables()
	fnCtx := cv.CreateContext(vars)
	if len(vars) > 0 {
		_ = fnCtx.SetVariableValue(vars[0], data.NewStringValue(arg))
	}
	if _, ctl := ctor.Call(fnCtx); ctl != nil {
		return nil, ctl
	}
	if v, _ := cv.GetProperty("uid"); v == nil || v.AsString() == "" {
		_ = cv.SetProperty("uid", data.NewStringValue(arg))
	}
	return cv, nil
}

// newClassNoArg 调用无参（或默认参）构造，对应 PHP `new UuidV4()` / `new UuidV7()`。
func newClassNoArg(ctx data.Context, stmt data.ClassStmt) (*data.ClassValue, data.Control) {
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	ctor := stmt.GetConstruct()
	if ctor == nil {
		return cv, nil
	}
	vars := ctor.GetVariables()
	fnCtx := cv.CreateContext(vars)
	if len(vars) > 0 {
		_ = fnCtx.SetVariableValue(vars[0], data.NewNullValue())
	}
	if _, ctl := ctor.Call(fnCtx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

func lateStaticStmt(ctx data.Context) data.ClassStmt {
	if cm, ok := ctx.(*data.ClassMethodContext); ok {
		if cm.StaticClass != nil {
			return cm.StaticClass
		}
		if cm.ClassValue != nil {
			return cm.ClassValue.Class
		}
	}
	return nil
}

func isNativeUuid(stmt data.ClassStmt) bool {
	return stmt == nil || stmt.GetName() == uuidName
}

func isNativeUlid(stmt data.ClassStmt) bool {
	return stmt == nil || stmt.GetName() == ulidName
}

func ulidTimeMs(s string) uint64 {
	b, ok := decodeCrockfordTo16(s)
	if !ok {
		return 0
	}
	return uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5])
}

func newDateTimeImmutable(ctx data.Context, ms uint64) (data.GetValue, data.Control) {
	t := time.UnixMilli(int64(ms)).UTC()
	stamp := t.Format("2006-01-02 15:04:05")
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewStringValue(stamp), nil
	}
	cls, ok := vm.GetClass("DateTimeImmutable")
	if !ok || cls == nil {
		if _, ctl := vm.LoadPkg("DateTimeImmutable"); ctl != nil {
			return nil, ctl
		}
		cls, ok = vm.GetClass("DateTimeImmutable")
	}
	if !ok || cls == nil {
		return data.NewStringValue(stamp), nil
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	ctor := cls.GetConstruct()
	if ctor == nil {
		return cv, nil
	}
	vars := ctor.GetVariables()
	fnCtx := cv.CreateContext(vars)
	if len(vars) > 0 {
		_ = fnCtx.SetVariableValue(vars[0], data.NewStringValue(stamp))
	}
	if _, ctl := ctor.Call(fnCtx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

func unixMilliFromValue(ctx data.Context, v data.Value) (uint64, data.Control) {
	if v == nil || isNullVal(v) {
		return uint64(time.Now().UnixMilli()), nil
	}
	if cv, ok := v.(*data.ClassValue); ok && cv != nil {
		if m, ok := cv.GetMethod("format"); ok && m != nil {
			vars := m.GetVariables()
			fnCtx := cv.CreateContext(vars)
			if len(vars) > 0 {
				_ = fnCtx.SetVariableValue(vars[0], data.NewStringValue("Uv"))
			}
			ret, ctl := m.Call(fnCtx)
			if ctl != nil {
				return 0, ctl
			}
			if ret != nil {
				if val, ok := ret.(data.Value); ok {
					s := val.AsString()
					if strings.HasPrefix(s, "-") {
						return 0, throwInvalid("The timestamp must be positive.")
					}
					if n, err := strconv.ParseUint(s, 10, 64); err == nil {
						return n, nil
					}
				}
			}
		}
		if m, ok := cv.GetMethod("getTimestamp"); ok && m != nil {
			ret, ctl := m.Call(cv.CreateContext(m.GetVariables()))
			if ctl != nil {
				return 0, ctl
			}
			if ret != nil {
				if val, ok := ret.(data.Value); ok {
					if ai, ok := val.(data.AsInt); ok {
						if n, err := ai.AsInt(); err == nil {
							if n < 0 {
								return 0, throwInvalid("The timestamp must be positive.")
							}
							return uint64(n) * 1000, nil
						}
					}
				}
			}
		}
	}
	if ai, ok := v.(data.AsInt); ok {
		if n, err := ai.AsInt(); err == nil {
			if n < 0 {
				return 0, throwInvalid("The timestamp must be positive.")
			}
			u := uint64(n)
			if u < 1e11 {
				return u * 1000, nil
			}
			return u, nil
		}
	}
	return uint64(time.Now().UnixMilli()), nil
}

func isNullVal(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func argString(ctx data.Context, i int) (string, bool) {
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil || isNullVal(v) {
		return "", false
	}
	return v.AsString(), true
}

func argBool(ctx data.Context, i int, def bool) bool {
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil || isNullVal(v) {
		return def
	}
	if ab, ok := v.(data.AsBool); ok {
		b, err := ab.AsBool()
		if err == nil {
			return b
		}
	}
	return def
}

func propertyUID(cv *data.ClassValue) string {
	if cv == nil {
		return ""
	}
	v, _ := cv.GetProperty("uid")
	if v == nil {
		return ""
	}
	return v.AsString()
}

func uidFromValue(v data.Value) (string, bool) {
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil || cv.Class == nil {
		return "", false
	}
	name := cv.Class.GetName()
	if name != abstractUidName && name != uuidName && name != ulidName && !strings.Contains(name, "\\Uid\\") {
		return "", false
	}
	return propertyUID(cv), true
}

func transformToRFC4122(raw string) (string, bool) {
	if isRFC4122(raw) {
		return strings.ToLower(raw), true
	}
	if len(raw) == 16 {
		return formatRFC4122([]byte(raw)), true
	}
	if ulidIsValid(raw) {
		b, ok := decodeCrockfordTo16(strings.ToUpper(raw))
		if !ok {
			return "", false
		}
		return formatRFC4122(b[:]), true
	}
	return "", false
}
