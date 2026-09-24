package support

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

// 对齐 laravel/framework 的官方契约：Str::uuid() / Str::uuid7() / Str::orderedUuid()
// 返回 Ramsey\Uuid\UuidInterface 实例（调用方普遍写 Str::orderedUuid()->toString()），
// 而不是裸字符串。
//
// 实现策略（性能优先，参照 std/symfony/uid 的做法）：
//   - 用 Go 原生 ClassStmt 顶替 vendor/ramsey/uuid 的 PHP 实现，只有被 Str 家族显式创建时才会分配；
//   - 类名与接口名和官方完全一致，instanceof 走 GetName()/GetImplements() 的名字匹配；
//   - 故意不调用 vm.AddClass 注册这个 Go 类：vendor 里的 Rfc4122\UuidV4 extends Uuid、
//     UuidFactory/UuidBuilder 等 PHP 继承链因此完全不受影响，该 Go 类只作为受控实例的类型出现。
const (
	ramseyUuidName           = "Ramsey\\Uuid\\Uuid"
	ramseyUuidInterfaceName  = "Ramsey\\Uuid\\UuidInterface"
	ramseyRfc4122UuidIface   = "Ramsey\\Uuid\\Rfc4122\\UuidInterface"
	ramseyFieldsName         = "Ramsey\\Uuid\\Rfc4122\\Fields"
	ramseyRfc4122FieldsIface = "Ramsey\\Uuid\\Rfc4122\\FieldsInterface"
	ramseyFieldsIfaceName    = "Ramsey\\Uuid\\Fields\\FieldsInterface"

	// uuidValueProp 实例状态：RFC4122 规范形式（小写 36 字符）。
	// toString/__toString/jsonSerialize/比较 全部只依赖它，保证只有一个真值来源。
	uuidValueProp = "uuid"
	// fieldsValueProp Fields 实例状态：同上的规范字符串。
	fieldsValueProp = "value"
)

const (
	uuidNil = "00000000-0000-0000-0000-000000000000"
	uuidMax = "ffffffff-ffff-ffff-ffff-ffffffffffff"
)

var (
	// errInvalidUUID 对齐 ramsey InvalidUuidStringException 的报错语义。
	errInvalidUUID = errors.New("Invalid UUID string")
	// errNotUuid 对齐 compareTo() 对非 UuidInterface 参数的 TypeError。
	errNotUuid = errors.New("Argument #1 ($other) must be of type Ramsey\\Uuid\\UuidInterface")
)

// NewRamseyUuidInterface 暴露 Ramsey\Uuid\UuidInterface（含 Rfc4122 派生接口）为 Go 原生接口。
// 注册后 instanceof UuidInterface 无需再加载 vendor 里的接口文件，
// 也不依赖 ramsey/uuid 包是否存在。
func NewRamseyUuidInterface() data.InterfaceStmt {
	base := []data.Method{
		node.NewInterfaceMethod(nil, "compareTo", "public", nil, data.NewBaseType("int")),
		node.NewInterfaceMethod(nil, "equals", "public", nil, data.NewBaseType("bool")),
		node.NewInterfaceMethod(nil, "getBytes", "public", nil, data.NewBaseType("string")),
		node.NewInterfaceMethod(nil, "getHex", "public", nil, nil),
		node.NewInterfaceMethod(nil, "getInteger", "public", nil, nil),
		node.NewInterfaceMethod(nil, "getUrn", "public", nil, data.NewBaseType("string")),
		node.NewInterfaceMethod(nil, "toString", "public", nil, data.NewBaseType("string")),
		node.NewInterfaceMethod(nil, "__toString", "public", nil, data.NewBaseType("string")),
		node.NewInterfaceMethod(nil, "jsonSerialize", "public", nil, data.NewBaseType("mixed")),
	}
	return node.NewInterfaceStatement(nil, ramseyUuidInterfaceName, []string{"JsonSerializable", "Stringable"}, base)
}

// NewRamseyRfc4122UuidInterface 暴露 Ramsey\Uuid\Rfc4122\UuidInterface。
func NewRamseyRfc4122UuidInterface() data.InterfaceStmt {
	return node.NewInterfaceStatement(nil, ramseyRfc4122UuidIface, []string{ramseyUuidInterfaceName}, nil)
}

// ramseyUuidClass Go 原生的 Ramsey\Uuid\Uuid。
type ramseyUuidClass struct {
	node.Node
	methods    map[string]data.Method
	properties []data.Property
	propIndex  map[string]data.Property
}

// 用 init() 而非包级表达式初始化，避免与 newUuidValue() 形成初始化环。
var (
	uuidValueStatement   data.ClassStmt
	fieldsValueStatement data.ClassStmt
)

func init() {
	uuidValueStatement = newRamseyUuidClass()
	fieldsValueStatement = newRamseyFieldsClass()
}

func newRamseyUuidClass() data.ClassStmt {
	c := &ramseyUuidClass{methods: map[string]data.Method{}}
	prop := node.NewProperty(nil, uuidValueProp, "protected", false, data.NewStringValue(uuidNil), data.NewBaseType("string"))
	c.properties = []data.Property{prop}
	c.propIndex = map[string]data.Property{uuidValueProp: prop}
	registerUuidMethods(c.methods)
	return c
}

func (c *ramseyUuidClass) GetName() string            { return ramseyUuidName }
func (c *ramseyUuidClass) GetExtend() *string         { return nil }
func (c *ramseyUuidClass) GetImplements() []string {
	return []string{ramseyRfc4122UuidIface, ramseyUuidInterfaceName, "JsonSerializable", "Stringable"}
}
func (c *ramseyUuidClass) GetProperty(name string) (data.Property, bool) {
	p, ok := c.propIndex[name]
	return p, ok
}
func (c *ramseyUuidClass) GetPropertyList() []data.Property { return c.properties }
func (c *ramseyUuidClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *ramseyUuidClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ramseyUuidClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *ramseyUuidClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *ramseyUuidClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	if !ok || m == nil || !m.GetIsStatic() {
		return nil, false
	}
	return m, true
}
func (c *ramseyUuidClass) GetStaticProperty(name string) (data.Value, bool) {
	// Uuid::NIL / Uuid::MAX / 命名空间常量 / UUID_TYPE_*（值对齐 ramsey/uuid 4.x）
	switch name {
	case "NIL":
		return data.NewStringValue(uuidNil), true
	case "MAX":
		return data.NewStringValue(uuidMax), true
	case "NAMESPACE_DNS":
		return data.NewStringValue("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), true
	case "NAMESPACE_URL":
		return data.NewStringValue("6ba7b811-9dad-11d1-80b4-00c04fd430c8"), true
	case "NAMESPACE_OID":
		return data.NewStringValue("6ba7b812-9dad-11d1-80b4-00c04fd430c8"), true
	case "NAMESPACE_X500":
		return data.NewStringValue("6ba7b814-9dad-11d1-80b4-00c04fd430c8"), true
	case "UUID_TYPE_TIME":
		return data.NewIntValue(1), true
	case "UUID_TYPE_DCE":
		return data.NewIntValue(2), true
	case "UUID_TYPE_HASH_MD5":
		return data.NewIntValue(3), true
	case "UUID_TYPE_RANDOM":
		return data.NewIntValue(4), true
	case "UUID_TYPE_HASH_SHA1":
		return data.NewIntValue(5), true
	case "UUID_TYPE_SORT_MAC":
		return data.NewIntValue(6), true
	case "UUID_TYPE_TIME_UNIX":
		return data.NewIntValue(7), true
	case "UUID_TYPE_CUSTOM":
		return data.NewIntValue(8), true
	}
	return nil, false
}

// ramseyFieldsClass Go 原生的 Ramsey\Uuid\Rfc4122\Fields。
type ramseyFieldsClass struct {
	node.Node
	methods    map[string]data.Method
	properties []data.Property
	propIndex  map[string]data.Property
}

func newRamseyFieldsClass() data.ClassStmt {
	c := &ramseyFieldsClass{methods: map[string]data.Method{}}
	prop := node.NewProperty(nil, fieldsValueProp, "protected", false, data.NewStringValue(uuidNil), data.NewBaseType("string"))
	c.properties = []data.Property{prop}
	c.propIndex = map[string]data.Property{fieldsValueProp: prop}
	registerFieldsMethods(c.methods)
	return c
}

func (c *ramseyFieldsClass) GetName() string        { return ramseyFieldsName }
func (c *ramseyFieldsClass) GetExtend() *string     { return nil }
func (c *ramseyFieldsClass) GetImplements() []string {
	return []string{ramseyRfc4122FieldsIface, ramseyFieldsIfaceName, "JsonSerializable"}
}
func (c *ramseyFieldsClass) GetProperty(name string) (data.Property, bool) {
	p, ok := c.propIndex[name]
	return p, ok
}
func (c *ramseyFieldsClass) GetPropertyList() []data.Property { return c.properties }
func (c *ramseyFieldsClass) GetConstruct() data.Method        { return nil }
func (c *ramseyFieldsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ramseyFieldsClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *ramseyFieldsClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}

// ------------------------------ 方法注册 ------------------------------

func registerUuidMethods(m map[string]data.Method) {
	inst := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		m[strings.ToLower(name)] = kit.InstanceMethodOpt(name, params, optionalFrom, fn)
	}
	st := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		m[strings.ToLower(name)] = kit.StaticMethod(name, params, optionalFrom, fn)
	}

	inst("__construct", []string{"uuid"}, 0, uuidConstruct)
	inst("toString", nil, -1, uuidToString)
	inst("__toString", nil, -1, uuidToString)
	inst("jsonSerialize", nil, -1, uuidToString)
	inst("getBytes", nil, -1, uuidGetBytes)
	inst("getHex", nil, -1, uuidGetHex)
	inst("getInteger", nil, -1, uuidGetInteger)
	inst("getUrn", nil, -1, uuidGetUrn)
	inst("getVersion", nil, -1, uuidGetVersion)
	inst("getVariant", nil, -1, uuidGetVariant)
	inst("isNull", nil, -1, uuidIsNull)
	inst("isNil", nil, -1, uuidIsNull)
	inst("equals", []string{"other"}, -1, uuidEquals)
	inst("compareTo", []string{"other"}, -1, uuidCompareTo)
	inst("getFields", nil, -1, uuidGetFields)
	inst("serialize", nil, -1, uuidToString)
	inst("__serialize", nil, -1, uuidSerializeArray)
	inst("unserialize", []string{"data"}, -1, uuidUnserialize)
	inst("__unserialize", []string{"data"}, -1, uuidUnserialize)
	inst("getDateTime", nil, -1, uuidGetDateTime)
	// DeprecatedUuidInterface 里的一票 *Hex 取值器：纯字符串切片，零分配压力。
	inst("getTimestampHex", nil, -1, uuidGetTimestampHex)
	inst("getTimeLowHex", nil, -1, uuidGetTimeLowHex)
	inst("getTimeMidHex", nil, -1, uuidGetTimeMidHex)
	inst("getTimeHiAndVersionHex", nil, -1, uuidGetTimeHiAndVersionHex)
	inst("getClockSeqHiAndReservedHex", nil, -1, uuidGetClockSeqHiHex)
	inst("getClockSeqLowHex", nil, -1, uuidGetClockSeqLowHex)
	inst("getClockSequenceHex", nil, -1, uuidGetClockSequenceHex)
	inst("getNodeHex", nil, -1, uuidGetNodeHex)
	inst("getLeastSignificantBitsHex", nil, -1, uuidGetLeastSignificantBitsHex)
	inst("getMostSignificantBitsHex", nil, -1, uuidGetMostSignificantBitsHex)
	inst("getFieldsHex", nil, -1, uuidGetFieldsHex)

	st("uuid4", nil, -1, uuidStaticV4)
	st("uuid7", []string{"time"}, 0, uuidStaticV7)
	st("fromString", []string{"uuid"}, -1, uuidStaticFrom)
	st("fromBytes", []string{"bytes"}, -1, uuidStaticFromBytes)
	st("isValid", []string{"uuid"}, -1, uuidStaticIsValid)
	st("isNull", []string{"uuid"}, -1, uuidStaticIsNull)
}

func registerFieldsMethods(m map[string]data.Method) {
	inst := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		m[strings.ToLower(name)] = kit.InstanceMethodOpt(name, params, optionalFrom, fn)
	}
	inst("__toString", nil, -1, fieldsToString)
	inst("toString", nil, -1, fieldsToString)
	inst("jsonSerialize", nil, -1, fieldsToString)
	inst("getBytes", nil, -1, fieldsGetBytes)
	inst("getVersion", nil, -1, fieldsGetVersion)
	inst("getVariant", nil, -1, fieldsGetVariant)
	inst("isNil", nil, -1, fieldsIsNil)
	inst("getTimestamp", nil, -1, fieldsGetTimestamp)
	inst("getTimeLow", nil, -1, fieldsGetTimeLow)
	inst("getTimeMid", nil, -1, fieldsGetTimeMid)
	inst("getTimeHiAndVersion", nil, -1, fieldsGetTimeHiAndVersion)
	inst("getClockSeq", nil, -1, fieldsGetClockSeq)
	inst("getClockSeqHiAndReserved", nil, -1, fieldsGetClockSeqHi)
	inst("getClockSeqLow", nil, -1, fieldsGetClockSeqLow)
	inst("getNode", nil, -1, fieldsGetNode)
}

// ------------------------------ 构造与取值 ------------------------------

// newUuidValue 由规范字符串构造 Ramsey\Uuid\Uuid 实例（唯一的实例化入口）。
func newUuidValue(ctx data.Context, canonical string) *data.ClassValue {
	cv := data.NewClassValue(uuidValueStatement, ctx.CreateBaseContext())
	_ = cv.SetProperty(uuidValueProp, data.NewStringValue(canonical))
	return cv
}

func uuidSelf(ctx data.Context) *data.ClassValue {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		return cv.ClassValue
	}
	return nil
}

// uuidValue 读实例上的规范字符串；ctx 不是实例上下文时退回实例自身（静态调用场景）。
// Uuid 的状态属性名是 "uuid"，Fields 的是 "value"，两者共用同一套实例方法实现，
// 所以这里按 "uuid" -> "value" 顺序兜底（Uuid 上 "uuid" 必然存在，不会走到兜底分支）。
func uuidValue(ctx data.Context) string {
	if v := uuidPropValue(ctx, uuidValueProp); v != "" {
		return v
	}
	return uuidPropValue(ctx, fieldsValueProp)
}

func fieldsValue(ctx data.Context) string {
	return uuidValue(ctx)
}

func uuidPropValue(ctx data.Context, prop string) string {
	cv := uuidSelf(ctx)
	if cv == nil {
		if v, ok := ctx.(*data.ClassValue); ok {
			cv = v
		}
	}
	if cv == nil {
		return ""
	}
	p, _ := cv.GetProperty(prop)
	if p == nil {
		return ""
	}
	return p.AsString()
}

func uuidConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := uuidSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	raw := kit.Arg(ctx, 0)
	if raw == nil || kit.IsNull(raw) {
		return data.NewNullValue(), nil
	}
	canonical, ok := uuidCanonical(raw.AsString())
	if !ok {
		return nil, data.NewErrorThrow(nil, errInvalidUUID)
	}
	_ = cv.SetProperty(uuidValueProp, data.NewStringValue(canonical))
	return data.NewNullValue(), nil
}

// uuidCanonical 统一成小写 36 字符规范形式：接受 RFC4122 文本或 16 字节二进制。
func uuidCanonical(raw string) (string, bool) {
	switch {
	case len(raw) == 36:
		if !strUUIDRe.MatchString(raw) {
			return "", false
		}
		return strings.ToLower(raw), true
	case len(raw) == 16:
		return uuidFormatBytes([]byte(raw)), true
	default:
		return "", false
	}
}

// uuidFormatBytes 把 16 字节按 8-4-4-4-12 输出为小写规范形式。
func uuidFormatBytes(b []byte) string {
	s := hex.EncodeToString(b)
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

func uuidPlainHex(v string) string {
	return strings.ReplaceAll(v, "-", "")
}

func uuidBytes(v string) []byte {
	b, err := hex.DecodeString(uuidPlainHex(v))
	if err != nil || len(b) != 16 {
		return make([]byte, 16)
	}
	return b
}

// ------------------------------ 生成 ------------------------------

// uuidFormatV4 随机 v4。
func uuidFormatV4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return uuidFormatBytes(b)
}

// uuidOrderedComb 复刻 Laravel Str::orderedUuid() 的官方算法：
// Ramsey CombGenerator（末 6 字节放 10µs 精度的 Unix 时间戳）+ TimestampFirstCombCodec
// （编码时把首 6 字节与末 6 字节对调）。
// 结果仍是 v4（版本位 4、变体位 10xx），但字符串前 48 位就是时间戳，按字符串排序即时序，
// 这就是官方语义下的「有序 UUID」（早期 RFC4122 COMB，Laravel 至今沿用）。
func uuidOrderedComb() string {
	now := time.Now()
	var b [16]byte
	_, _ = rand.Read(b[:10])
	// microtime(false) 对齐：秒 * 1e5 + 微秒 / 10（CombGenerator 的 0.00001s 精度）
	ts := uint64(now.Unix())*100000 + uint64(now.Nanosecond()/10000)
	b[10] = byte(ts >> 40)
	b[11] = byte(ts >> 32)
	b[12] = byte(ts >> 24)
	b[13] = byte(ts >> 16)
	b[14] = byte(ts >> 8)
	b[15] = byte(ts)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	var out [16]byte
	copy(out[0:6], b[10:16])
	copy(out[6:10], b[6:10])
	copy(out[10:16], b[0:6])
	return uuidFormatBytes(out[:])
}

// uuidFormatV7 由 Unix 毫秒生成 v7（时间前缀，官方 Str::uuid7() 语义）。
func uuidFormatV7(ms uint64) string {
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
	return uuidFormatBytes(b)
}

// ------------------------------ 实例方法 ------------------------------

func uuidToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(uuidValue(ctx)), nil
}

func uuidGetBytes(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(string(uuidBytes(uuidValue(ctx)))), nil
}

func uuidGetHex(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(uuidPlainHex(uuidValue(ctx))), nil
}

// uuidGetInteger 官方返回 Ramsey\Uuid\Type\Integer；这里返回其十进制文本。
func uuidGetInteger(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	// 128 位无符号十进制，超出 int64 用大数风格手工拼
	return data.NewStringValue(uuidDecimal(b)), nil
}

func uuidDecimal(b []byte) string {
	digits := []byte{'0'}
	for _, by := range b {
		carry := int(by)
		for i := len(digits) - 1; i >= 0; i-- {
			cur := int(digits[i]-'0')*256 + carry
			digits[i] = byte('0' + cur%10)
			carry = cur / 10
		}
		for carry > 0 {
			digits = append([]byte{byte('0' + carry%10)}, digits...)
			carry /= 10
		}
	}
	return string(digits)
}

func uuidGetUrn(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue("urn:uuid:" + uuidValue(ctx)), nil
}

func uuidGetVersion(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(uuidVersion(uuidValue(ctx))), nil
}

// uuidVersion 取第 15 个字符（版本位）；RFC4122 v7 为 7，COMB v4 为 4。
func uuidVersion(v string) int {
	if len(v) != 36 {
		return 0
	}
	c := v[14]
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	if c >= 'a' && c <= 'f' {
		return int(c-'a') + 10
	}
	return 0
}

// uuidVariant 取变体位（RFC4122：0b10xx → 返回 0b100/0b110 等前缀右移结果，对齐 ramsey 的取法）。
func uuidVariant(v string) int {
	if len(v) != 36 {
		return 0
	}
	b := uuidBytes(v)
	return int(b[8] >> 5)
}

func uuidIsNull(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(uuidValue(ctx) == uuidNil), nil
}

func uuidEquals(ctx data.Context) (data.GetValue, data.Control) {
	self := uuidValue(ctx)
	other, ok := uuidValueOf(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(self == other), nil
}

func uuidCompareTo(ctx data.Context) (data.GetValue, data.Control) {
	self := uuidValue(ctx)
	other, ok := uuidValueOf(ctx, 0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errNotUuid)
	}
	switch {
	case self < other:
		return data.NewIntValue(-1), nil
	case self > other:
		return data.NewIntValue(1), nil
	default:
		return data.NewIntValue(0), nil
	}
}

func uuidGetFields(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(fieldsValueStatement, ctx.CreateBaseContext())
	_ = cv.SetProperty(fieldsValueProp, data.NewStringValue(uuidValue(ctx)))
	return cv, nil
}

// uuidSerializeArray 对应 __serialize(): array —— PHP 序列化只需带出规范字符串。
func uuidSerializeArray(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewArrayValue([]data.Value{data.NewStringValue(uuidValue(ctx))}), nil
}

func uuidUnserialize(ctx data.Context) (data.GetValue, data.Control) {
	cv := uuidSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	raw := ""
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		raw = v.AsString()
	}
	if canonical, ok := uuidCanonical(raw); ok {
		_ = cv.SetProperty(uuidValueProp, data.NewStringValue(canonical))
	}
	return data.NewNullValue(), nil
}

// uuidGetDateTime 官方对时间型 UUID 返回 DateTimeImmutable；随机型（v4/COMB）无时间语义，
// 官方实现此处会抛异常，这里按 null 返回（见汇报中的已知缺口）。
func uuidGetDateTime(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func uuidGetTimestampHex(ctx data.Context) (data.GetValue, data.Control) {
	v := uuidValue(ctx)
	if len(v) != 36 {
		return data.NewStringValue(""), nil
	}
	b := uuidBytes(v)
	return data.NewStringValue(hex.EncodeToString(b[0:4]) + hex.EncodeToString(b[4:6]) + uuidTimeHiHex(b)), nil
}

func uuidTimeHiHex(b []byte) string {
	s := hex.EncodeToString(b[6:8])
	// 抹掉版本半字节，与 ramsey Fields::getTimeHiAndVersion() 一致
	return "0" + strings.ToLower(s[1:])
}

func uuidGetTimeLowHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[0:4])), nil
}

func uuidGetTimeMidHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[4:6])), nil
}

func uuidGetTimeHiAndVersionHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(uuidTimeHiHex(b)), nil
}

func uuidGetClockSeqHiHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[8:9])), nil
}

func uuidGetClockSeqLowHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[9:10])), nil
}

func uuidGetClockSequenceHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[8:10])), nil
}

func uuidGetNodeHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[10:16])), nil
}

func uuidGetLeastSignificantBitsHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[8:16])), nil
}

func uuidGetMostSignificantBitsHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	return data.NewStringValue(hex.EncodeToString(b[0:8])), nil
}

func uuidGetFieldsHex(ctx data.Context) (data.GetValue, data.Control) {
	b := uuidBytes(uuidValue(ctx))
	arr := &data.ArrayValue{}
	put := func(k, v string) {
		arr.List = append(arr.List, data.NewNamedZVal(k, data.NewStringValue(v)))
	}
	put("time_low", hex.EncodeToString(b[0:4]))
	put("time_mid", hex.EncodeToString(b[4:6]))
	put("time_hi_and_version", uuidTimeHiHex(b))
	put("clock_seq_hi_and_reserved", hex.EncodeToString(b[8:9]))
	put("clock_seq_low", hex.EncodeToString(b[9:10]))
	put("node", hex.EncodeToString(b[10:16]))
	return arr, nil
}

func uuidGetVariant(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(uuidVariant(uuidValue(ctx))), nil
}

// ------------------------------ 静态方法 ------------------------------

func uuidStaticV4(ctx data.Context) (data.GetValue, data.Control) {
	return newUuidValue(ctx, uuidFormatV4()), nil
}

func uuidStaticV7(ctx data.Context) (data.GetValue, data.Control) {
	return newUuidValue(ctx, uuidFormatV7(uuidTimeMilli(ctx, 0))), nil
}

func uuidStaticFrom(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 0)
	if v == nil || kit.IsNull(v) {
		return nil, data.NewErrorThrow(nil, errInvalidUUID)
	}
	canonical, ok := uuidCanonical(v.AsString())
	if !ok {
		return nil, data.NewErrorThrow(nil, errInvalidUUID)
	}
	return newUuidValue(ctx, canonical), nil
}

func uuidStaticFromBytes(ctx data.Context) (data.GetValue, data.Control) {
	raw := ""
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		raw = v.AsString()
	}
	if len(raw) != 16 {
		return nil, data.NewErrorThrow(nil, errInvalidUUID)
	}
	return newUuidValue(ctx, uuidFormatBytes([]byte(raw))), nil
}

func uuidStaticIsValid(ctx data.Context) (data.GetValue, data.Control) {
	canonical, ok := uuidValueOf(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(uuidIsRFC4122(canonical)), nil
}

func uuidStaticIsNull(ctx data.Context) (data.GetValue, data.Control) {
	canonical, ok := uuidValueOf(ctx, 0)
	return data.NewBoolValue(ok && canonical == uuidNil), nil
}

// uuidIsRFC4122 校验 RFC4122 文本 + 变体位（nil/max 特例与官方一致）。
func uuidIsRFC4122(v string) bool {
	if !strUUIDRe.MatchString(v) {
		return false
	}
	if v == uuidNil || v == uuidMax {
		return true
	}
	switch v[19] {
	case '8', '9', 'a', 'b', 'A', 'B':
		return true
	}
	return false
}

// uuidTimeMilli 解析 uuid7($time) 的时间参数：null → 当前；
// DateTimeInterface（带 getTimestamp()）→ 秒；数字 → 按量级区分秒/毫秒。
func uuidTimeMilli(ctx data.Context, i int) uint64 {
	v := kit.Arg(ctx, i)
	if v == nil || kit.IsNull(v) {
		return uint64(time.Now().UnixMilli())
	}
	if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok && cv != nil {
		if m, ok := cv.GetMethod("getTimestamp"); ok && m != nil {
			fnCtx := cv.CreateContext(m.GetVariables())
			if ret, ctl := m.Call(fnCtx); ctl == nil && ret != nil {
				if iv, ok := ret.(data.AsInt); ok {
					if n, err := iv.AsInt(); err == nil && n > 0 {
						return uint64(n) * 1000
					}
				}
			}
		}
	}
	if ai, ok := v.(data.AsInt); ok {
		if n, err := ai.AsInt(); err == nil && n > 0 {
			if n < 100000000000 { // < 1e11 视为 Unix 秒
				return uint64(n) * 1000
			}
			return uint64(n)
		}
	}
	return uint64(time.Now().UnixMilli())
}

// uuidValueOf 从第 i 个实参解析 uuid：实例走实例属性/__toString，字符串按规范形式解析。
func uuidValueOf(ctx data.Context, i int) (string, bool) {
	v := kit.Arg(ctx, i)
	if v == nil || kit.IsNull(v) {
		return "", false
	}
	if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok && cv != nil {
		if p, _ := cv.GetProperty(uuidValueProp); p != nil {
			if canonical, ok := uuidCanonical(p.AsString()); ok {
				return canonical, true
			}
		}
		if m, ok := cv.GetMethod("__toString"); ok && m != nil {
			if ret, ctl := m.Call(cv.CreateContext(m.GetVariables())); ctl == nil && ret != nil {
				if sv, ok := ret.(data.Value); ok {
					if canonical, ok := uuidCanonical(sv.AsString()); ok {
						return canonical, true
					}
				}
			}
		}
		return "", false
	}
	return uuidCanonical(v.AsString())
}

// ------------------------------ Fields 方法 ------------------------------

func fieldsToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(fieldsValue(ctx)), nil
}

func fieldsGetBytes(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(string(uuidBytes(fieldsValue(ctx)))), nil
}

func fieldsGetVersion(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(uuidVersion(fieldsValue(ctx))), nil
}

func fieldsGetVariant(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(uuidVariant(fieldsValue(ctx))), nil
}

func fieldsIsNil(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(fieldsValue(ctx) == uuidNil), nil
}

func fieldsGetTimestamp(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetTimestampHex(ctx)
}

func fieldsGetTimeLow(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetTimeLowHex(ctx)
}

func fieldsGetTimeMid(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetTimeMidHex(ctx)
}

func fieldsGetTimeHiAndVersion(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetTimeHiAndVersionHex(ctx)
}

func fieldsGetClockSeq(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetClockSequenceHex(ctx)
}

func fieldsGetClockSeqHi(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetClockSeqHiHex(ctx)
}

func fieldsGetClockSeqLow(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetClockSeqLowHex(ctx)
}

func fieldsGetNode(ctx data.Context) (data.GetValue, data.Control) {
	return uuidGetNodeHex(ctx)
}
