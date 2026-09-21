package support

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// 对齐 Laravel Str::uuid / orderedUuid：返回带 __toString 的对象，不抢注 Ramsey\Uuid\Uuid，
// 以免挡住 vendor ramsey 的其余类。
const strUUIDObjectName = "Illuminate\\Support\\NativeUuid"

var (
	strUUIDFactoryMu sync.RWMutex
	strUUIDFactory   data.Value
	uuidRFC4122Re    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type strUUIDClass struct {
	node.Node
	methods map[string]data.Method
}

var strUUIDClassOnce = newStrUUIDClass()

func newStrUUIDClass() *strUUIDClass {
	c := &strUUIDClass{methods: map[string]data.Method{}}
	add := func(name string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = newInstanceMethod(name, nil, fn, false)
	}
	add("__toString", strUUIDToString)
	add("toString", strUUIDToString)
	add("toRfc4122", strUUIDToString)
	add("jsonSerialize", strUUIDToString)
	return c
}

func (c *strUUIDClass) GetName() string { return strUUIDObjectName }
func (c *strUUIDClass) GetExtend() *string {
	return nil
}
func (c *strUUIDClass) GetImplements() []string {
	return []string{"Ramsey\\Uuid\\UuidInterface", "JsonSerializable", "Stringable"}
}
func (c *strUUIDClass) GetProperty(name string) (data.Property, bool) {
	if name == "uid" {
		return node.NewProperty(nil, "uid", "protected", false, data.NewStringValue("")), true
	}
	return nil, false
}
func (c *strUUIDClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "uid", "protected", false, data.NewStringValue(""))}
}
func (c *strUUIDClass) GetConstruct() data.Method { return nil }
func (c *strUUIDClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *strUUIDClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *strUUIDClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *strUUIDClass) GetStaticMethod(string) (data.Method, bool) { return nil, false }

func wrapUUID(ctx data.Context, s string) *data.ClassValue {
	base := ctx
	if ctx != nil {
		base = ctx.CreateBaseContext()
	}
	cv := data.NewClassValue(strUUIDClassOnce, base)
	_ = cv.SetProperty("uid", data.NewStringValue(s))
	return cv
}

func strUUIDToString(ctx data.Context) (data.GetValue, data.Control) {
	cv := overlayReceiver(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	v, ctl := cv.GetProperty("uid")
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(v.AsString()), nil
}

func formatRFC4122(b []byte) string {
	s := hex.EncodeToString(b)
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

func generateUUIDv4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return formatRFC4122(b)
}

func generateUUIDv7() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	ms := uint64(time.Now().UnixMilli())
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

func strUUIDFromFactory(ctx data.Context) (data.GetValue, data.Control, bool) {
	strUUIDFactoryMu.RLock()
	f := strUUIDFactory
	strUUIDFactoryMu.RUnlock()
	if f == nil || isNull(f) {
		return nil, nil, false
	}
	v, ctl := callValue(ctx, f)
	return v, ctl, true
}

func strUUID(ctx data.Context) (data.GetValue, data.Control) {
	if v, ctl, ok := strUUIDFromFactory(ctx); ok {
		return v, ctl
	}
	return wrapUUID(ctx, generateUUIDv4()), nil
}

func strUUID7(ctx data.Context) (data.GetValue, data.Control) {
	if v, ctl, ok := strUUIDFromFactory(ctx); ok {
		return v, ctl
	}
	return wrapUUID(ctx, generateUUIDv7()), nil
}

func strOrderedUUID(ctx data.Context) (data.GetValue, data.Control) {
	if v, ctl, ok := strUUIDFromFactory(ctx); ok {
		return v, ctl
	}
	return wrapUUID(ctx, generateUUIDv7()), nil
}

func strCreateUuidsUsing(ctx data.Context) (data.GetValue, data.Control) {
	strUUIDFactoryMu.Lock()
	defer strUUIDFactoryMu.Unlock()
	if v, ok := ctx.GetIndexValue(0); ok && v != nil && !isNull(v) {
		strUUIDFactory = v
	} else {
		strUUIDFactory = nil
	}
	return data.NewNullValue(), nil
}

func strCreateUuidsNormally(ctx data.Context) (data.GetValue, data.Control) {
	strUUIDFactoryMu.Lock()
	strUUIDFactory = nil
	strUUIDFactoryMu.Unlock()
	return data.NewNullValue(), nil
}

func strIsUuid(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok || v == nil {
		return data.NewBoolValue(false), nil
	}
	if _, isStr := v.(*data.StringValue); !isStr {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(uuidRFC4122Re.MatchString(v.AsString())), nil
}

func strUlid(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewStringValue(""), nil
	}
	c, ok := vm.GetClass("Symfony\\Component\\Uid\\Ulid")
	if !ok || c == nil {
		return data.NewStringValue(""), nil
	}
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	ctor := c.GetConstruct()
	if ctor == nil {
		return cv, nil
	}
	vars := ctor.GetVariables()
	nctx := cv.CreateContext(vars)
	if _, ctl := ctor.Call(nctx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}
