package uid

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type UlidClass struct {
	node.Node
	methods    map[string]data.Method
	properties []data.Property
	propIndex  map[string]data.Property
}

var ulidClassOnce = &UlidClass{}

func NewUlidClass() data.ClassStmt {
	c := ulidClassOnce
	if c.methods != nil {
		return c
	}
	c.methods = map[string]data.Method{}
	uidProp := node.NewProperty(nil, "uid", "protected", false, data.NewStringValue(""), data.NewBaseType("string"))
	c.properties = []data.Property{uidProp}
	c.propIndex = map[string]data.Property{"uid": uidProp}

	strT := data.NewBaseType("string")
	add := func(name string, static bool, params []uidParam, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = &uidMethod{name: name, static: static, params: params, fn: fn}
	}
	add("__construct", false, []uidParam{{name: "ulid", def: data.NewNullValue(), ty: strT}}, ulidConstruct)
	add("generate", true, []uidParam{{name: "time", def: data.NewNullValue()}}, ulidGenerate)
	add("fromString", true, []uidParam{{name: "ulid", ty: strT}}, ulidFromString)
	add("fromBinary", true, []uidParam{{name: "ulid", ty: strT}}, ulidFromBinary)
	add("fromRfc4122", true, []uidParam{{name: "ulid", ty: strT}}, ulidFromRfc4122)
	add("fromBase32", true, []uidParam{{name: "ulid", ty: strT}}, ulidFromBase32)
	add("isValid", true, []uidParam{{name: "ulid", ty: strT}}, ulidIsValidMethod)
	add("toString", false, nil, ulidToString)
	add("__toString", false, nil, ulidToString)
	add("toBase32", false, nil, ulidToString)
	add("toRfc4122", false, nil, ulidToRfc4122)
	add("toBinary", false, nil, ulidToBinary)
	add("getDateTime", false, nil, ulidGetDateTime)
	add("equals", false, []uidParam{{name: "other"}}, uuidEquals)
	add("compare", false, []uidParam{{name: "other"}}, uuidCompare)
	add("hash", false, nil, ulidToString)
	add("jsonSerialize", false, nil, ulidToString)
	return c
}

func (c *UlidClass) GetName() string { return ulidName }
func (c *UlidClass) GetExtend() *string {
	parent := abstractUidName
	return &parent
}
func (c *UlidClass) GetImplements() []string {
	return []string{"JsonSerializable", "Stringable", hashableName, timeBasedName}
}
func (c *UlidClass) GetProperty(name string) (data.Property, bool) {
	p, ok := c.propIndex[name]
	return p, ok
}
func (c *UlidClass) GetPropertyList() []data.Property { return c.properties }
func (c *UlidClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *UlidClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *UlidClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *UlidClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *UlidClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	if !ok || m == nil || !m.GetIsStatic() {
		return nil, false
	}
	return m, true
}
func (c *UlidClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "FORMAT_BINARY":
		return data.NewIntValue(1), true
	case "FORMAT_BASE_32":
		return data.NewIntValue(1 << 1), true
	case "FORMAT_BASE_58":
		return data.NewIntValue(1 << 2), true
	case "FORMAT_RFC_4122", "FORMAT_RFC_9562":
		return data.NewIntValue(1 << 3), true
	case "FORMAT_ALL":
		return data.NewIntValue(-1), true
	case "NIL":
		return data.NewStringValue(ulidNil), true
	case "MAX":
		return data.NewStringValue(ulidMax), true
	}
	return nil, false
}

func ulidConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := uidSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	raw, ok := argString(ctx, 0)
	if !ok {
		ms, ctl := unixMilliFromValue(ctx, nil)
		if ctl != nil {
			return nil, ctl
		}
		_ = cv.SetProperty("uid", data.NewStringValue(generateULID(ms)))
		return data.NewNullValue(), nil
	}
	if raw == ulidNil {
		_ = cv.SetProperty("uid", data.NewStringValue(ulidNil))
		return data.NewNullValue(), nil
	}
	uid := strings.ToUpper(raw)
	if uid == ulidMax {
		_ = cv.SetProperty("uid", data.NewStringValue(ulidMax))
		return data.NewNullValue(), nil
	}
	if !ulidIsValid(raw) {
		return nil, throwInvalid("Invalid ULID.")
	}
	_ = cv.SetProperty("uid", data.NewStringValue(uid))
	return data.NewNullValue(), nil
}

func ulidIsValidMethod(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	u := strings.ToUpper(raw)
	if u == ulidNil || u == ulidMax {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(ulidIsValid(raw)), nil
}

func ulidGenerate(ctx data.Context) (data.GetValue, data.Control) {
	var timeArg data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		timeArg = v
	}
	ms, ctl := unixMilliFromValue(ctx, timeArg)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(generateULID(ms)), nil
}

func ulidFromString(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok {
		return nil, throwInvalid("Invalid ULID.")
	}
	uid, errMsg := normalizeUlidString(raw)
	if errMsg != "" {
		return nil, throwInvalid(errMsg)
	}
	return newClassWithStringArg(ctx, ulidClassFor(ctx, uid), uid)
}

func ulidFromBinary(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok || len(raw) != 16 {
		return nil, throwInvalid("Invalid binary uid provided.")
	}
	return ulidFromString(ctx)
}

func ulidFromRfc4122(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok || len(raw) != 36 {
		return nil, throwInvalid("Invalid RFC4122 uid provided.")
	}
	return ulidFromString(ctx)
}

func ulidFromBase32(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok || len(raw) != 26 {
		return nil, throwInvalid("Invalid base-32 uid provided.")
	}
	return ulidFromString(ctx)
}

func ulidClassFor(ctx data.Context, uid string) data.ClassStmt {
	if stmt := lateStaticStmt(ctx); !isNativeUlid(stmt) {
		return stmt
	}
	if uid == ulidNil {
		if c := optionalClass(ctx, nilUlidName); c != nil {
			return c
		}
	}
	if uid == ulidMax {
		if c := optionalClass(ctx, maxUlidName); c != nil {
			return c
		}
	}
	return NewUlidClass()
}

func normalizeUlidString(raw string) (string, string) {
	if isRFC4122(raw) {
		b, ok := rfc4122Bytes(raw)
		if !ok {
			return "", "Invalid ULID."
		}
		return encode128ToCrockford(b), ""
	}
	if len(raw) == 16 {
		return encode128ToCrockford([]byte(raw)), ""
	}
	if raw == ulidNil || strings.ToUpper(raw) == ulidNil {
		return ulidNil, ""
	}
	u := strings.ToUpper(raw)
	if u == ulidMax {
		return ulidMax, ""
	}
	if !ulidIsValid(raw) {
		return "", "Invalid ULID."
	}
	return u, ""
}

func ulidToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(propertyUID(uidSelf(ctx))), nil
}

func ulidToRfc4122(ctx data.Context) (data.GetValue, data.Control) {
	s := propertyUID(uidSelf(ctx))
	b, ok := decodeCrockfordTo16(s)
	if !ok {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(formatRFC4122(b[:])), nil
}

func ulidToBinary(ctx data.Context) (data.GetValue, data.Control) {
	s := propertyUID(uidSelf(ctx))
	b, ok := decodeCrockfordTo16(s)
	if !ok {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(b[:])), nil
}

func ulidGetDateTime(ctx data.Context) (data.GetValue, data.Control) {
	return newDateTimeImmutable(ctx, ulidTimeMs(propertyUID(uidSelf(ctx))))
}
