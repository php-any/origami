package uid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type uidParam struct {
	name string
	def  data.GetValue
	ty   data.Types
}

type uidMethod struct {
	name   string
	params []uidParam
	static bool
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *uidMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *uidMethod) GetName() string                                     { return m.name }
func (m *uidMethod) GetModifier() data.Modifier                          { return data.ModifierPublic }
func (m *uidMethod) GetIsStatic() bool                                   { return m.static }
func (m *uidMethod) GetReturnType() data.Types                           { return nil }
func (m *uidMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p.name, i, p.def, p.ty)
	}
	return out
}
func (m *uidMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p.name, i, p.ty)
	}
	return out
}

func throwInvalid(msg string) data.Control {
	return data.NewErrorThrowByName(nil, errors.New(msg), invalidArgEx)
}

func uidSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func lateStaticType(ctx data.Context) int {
	var stmt data.ClassStmt
	if cm, ok := ctx.(*data.ClassMethodContext); ok {
		if cm.StaticClass != nil {
			stmt = cm.StaticClass
		} else if cm.ClassValue != nil {
			stmt = cm.ClassValue.Class
		}
	}
	if stmt == nil {
		return 0
	}
	if gsp, ok := stmt.(data.GetStaticProperty); ok {
		if v, found := gsp.GetStaticProperty("TYPE"); found && v != nil {
			if ai, ok := v.(data.AsInt); ok {
				if n, err := ai.AsInt(); err == nil {
					return n
				}
			}
		}
	}
	return 0
}

type UuidClass struct {
	node.Node
	methods    map[string]data.Method
	properties []data.Property
	propIndex  map[string]data.Property
}

var uuidClassOnce = &UuidClass{}

func NewUuidClass() data.ClassStmt {
	c := uuidClassOnce
	if c.methods != nil {
		return c
	}
	c.methods = map[string]data.Method{}
	uidProp := node.NewProperty(nil, "uid", "protected", false, data.NewStringValue(""), data.NewBaseType("string"))
	c.properties = []data.Property{uidProp}
	c.propIndex = map[string]data.Property{"uid": uidProp}

	strT := data.NewBaseType("string")
	boolT := data.NewBaseType("bool")
	add := func(name string, static bool, params []uidParam, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = &uidMethod{name: name, static: static, params: params, fn: fn}
	}
	add("__construct", false, []uidParam{
		{name: "uuid", ty: strT},
		{name: "checkVariant", def: data.NewBoolValue(false), ty: boolT},
	}, uuidConstruct)
	add("v4", true, nil, uuidV4)
	add("v7", true, []uidParam{{name: "time", def: data.NewNullValue()}}, uuidV7)
	add("fromString", true, []uidParam{{name: "uuid", ty: strT}}, uuidFromString)
	add("fromBinary", true, []uidParam{{name: "uuid", ty: strT}}, uuidFromBinary)
	add("fromRfc4122", true, []uidParam{{name: "uuid", ty: strT}}, uuidFromRfc4122)
	add("fromBase32", true, []uidParam{{name: "uuid", ty: strT}}, uuidFromBase32)
	add("isValid", true, []uidParam{
		{name: "uuid", ty: strT},
		{name: "format", def: data.NewIntValue(1 << 3)},
	}, uuidIsValid)
	add("toRfc4122", false, nil, uuidToString)
	add("toString", false, nil, uuidToString)
	add("__toString", false, nil, uuidToString)
	add("toBase32", false, nil, uuidToBase32)
	add("toBinary", false, nil, uuidToBinary)
	add("equals", false, []uidParam{{name: "other"}}, uuidEquals)
	add("compare", false, []uidParam{{name: "other"}}, uuidCompare)
	add("isNil", false, nil, uuidIsNil)
	add("hash", false, nil, uuidToString)
	add("jsonSerialize", false, nil, uuidToString)
	return c
}

func (c *UuidClass) GetName() string { return uuidName }
func (c *UuidClass) GetExtend() *string {
	parent := abstractUidName
	return &parent
}
func (c *UuidClass) GetImplements() []string {
	return []string{"JsonSerializable", "Stringable", hashableName}
}
func (c *UuidClass) GetProperty(name string) (data.Property, bool) {
	p, ok := c.propIndex[name]
	return p, ok
}
func (c *UuidClass) GetPropertyList() []data.Property { return c.properties }
func (c *UuidClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *UuidClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *UuidClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *UuidClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *UuidClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	if !ok || m == nil || !m.GetIsStatic() {
		return nil, false
	}
	return m, true
}
func (c *UuidClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "NAMESPACE_DNS":
		return data.NewStringValue("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), true
	case "NAMESPACE_URL":
		return data.NewStringValue("6ba7b811-9dad-11d1-80b4-00c04fd430c8"), true
	case "NAMESPACE_OID":
		return data.NewStringValue("6ba7b812-9dad-11d1-80b4-00c04fd430c8"), true
	case "NAMESPACE_X500":
		return data.NewStringValue("6ba7b814-9dad-11d1-80b4-00c04fd430c8"), true
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
	case "TYPE":
		return data.NewIntValue(0), true
	case "NIL":
		return data.NewStringValue(uuidNil), true
	case "MAX":
		return data.NewStringValue(uuidMax), true
	}
	return nil, false
}

func uuidConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := uidSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	raw, ok := argString(ctx, 0)
	if !ok {
		return nil, throwInvalid("Invalid UUID.")
	}
	if !isRFC4122(raw) {
		typ := lateStaticType(ctx)
		return nil, throwInvalid(invalidUUIDMsg(typ))
	}
	uid := strings.ToLower(raw)
	typ := lateStaticType(ctx)
	if typ > 0 {
		ver := int(uid[14] - '0')
		if uid[14] < '0' || uid[14] > '9' || ver != typ {
			return nil, throwInvalid(invalidUUIDMsg(typ))
		}
	}
	if argBool(ctx, 1, false) {
		switch uid[19] {
		case '8', '9', 'a', 'b':
		default:
			return nil, throwInvalid(invalidUUIDMsg(typ))
		}
	}
	_ = cv.SetProperty("uid", data.NewStringValue(uid))
	return data.NewNullValue(), nil
}

func invalidUUIDMsg(typ int) string {
	if typ > 0 {
		return fmt.Sprintf("Invalid UUIDv%d.", typ)
	}
	return "Invalid UUID."
}

func uuidV4(ctx data.Context) (data.GetValue, data.Control) {
	if cls := optionalClass(ctx, uuidV4Name); cls != nil {
		return newClassNoArg(ctx, cls)
	}
	return newClassWithStringArg(ctx, NewUuidClass(), generateUUIDv4())
}

func uuidV7(ctx data.Context) (data.GetValue, data.Control) {
	var timeArg data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		timeArg = v
	}
	if timeArg == nil || isNullVal(timeArg) {
		if cls := optionalClass(ctx, uuidV7Name); cls != nil {
			return newClassNoArg(ctx, cls)
		}
		ms, _ := unixMilliFromValue(ctx, nil)
		return newClassWithStringArg(ctx, NewUuidClass(), generateUUIDv7(ms))
	}
	ms, ctl := unixMilliFromValue(ctx, timeArg)
	if ctl != nil {
		return nil, ctl
	}
	s := generateUUIDv7(ms)
	if cls := optionalClass(ctx, uuidV7Name); cls != nil {
		return newClassWithStringArg(ctx, cls, s)
	}
	return newClassWithStringArg(ctx, NewUuidClass(), s)
}

func uuidFromString(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok {
		return nil, throwInvalid("Invalid UUID.")
	}
	rfc, ok := transformToRFC4122(raw)
	if !ok {
		return nil, throwInvalid("Invalid UUID.")
	}
	if stmt := lateStaticStmt(ctx); !isNativeUuid(stmt) {
		return newClassWithStringArg(ctx, stmt, rfc)
	}
	return newClassWithStringArg(ctx, uuidClassForRFC4122(ctx, rfc), rfc)
}

func uuidFromBinary(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok || len(raw) != 16 {
		return nil, throwInvalid("Invalid binary uid provided.")
	}
	return uuidFromString(ctx)
}

func uuidFromRfc4122(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok || len(raw) != 36 {
		return nil, throwInvalid("Invalid RFC4122 uid provided.")
	}
	return uuidFromString(ctx)
}

func uuidFromBase32(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok || len(raw) != 26 {
		return nil, throwInvalid("Invalid base-32 uid provided.")
	}
	return uuidFromString(ctx)
}

func uuidClassForRFC4122(ctx data.Context, rfc string) data.ClassStmt {
	if rfc == uuidNil {
		if c := optionalClass(ctx, nilUuidName); c != nil {
			return c
		}
	}
	if rfc == uuidMax {
		if c := optionalClass(ctx, maxUuidName); c != nil {
			return c
		}
	}
	if len(rfc) >= 15 {
		switch rfc[14] {
		case '1', '3', '4', '5', '6', '7', '8':
			name := "Symfony\\Component\\Uid\\UuidV" + string(rfc[14])
			if c := optionalClass(ctx, name); c != nil {
				return c
			}
		}
	}
	return NewUuidClass()
}

func uuidToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(propertyUID(uidSelf(ctx))), nil
}

func uuidToBase32(ctx data.Context) (data.GetValue, data.Control) {
	s := propertyUID(uidSelf(ctx))
	b, ok := rfc4122Bytes(s)
	if !ok {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(encode128ToCrockford(b)), nil
}

func uuidToBinary(ctx data.Context) (data.GetValue, data.Control) {
	s := propertyUID(uidSelf(ctx))
	b, ok := rfc4122Bytes(s)
	if !ok {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(b)), nil
}

func uuidEquals(ctx data.Context) (data.GetValue, data.Control) {
	self := propertyUID(uidSelf(ctx))
	other, ok := ctx.GetIndexValue(0)
	if !ok || other == nil {
		return data.NewBoolValue(false), nil
	}
	got, isObj := uidFromValue(other)
	if !isObj {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(self == got), nil
}

func uuidCompare(ctx data.Context) (data.GetValue, data.Control) {
	self := propertyUID(uidSelf(ctx))
	other, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrowByName(nil, errors.New("Argument #1 ($other) must be of type Symfony\\Component\\Uid\\AbstractUid"), "TypeError")
	}
	got, isObj := uidFromValue(other)
	if !isObj {
		return nil, data.NewErrorThrowByName(nil, errors.New("Argument #1 ($other) must be of type Symfony\\Component\\Uid\\AbstractUid"), "TypeError")
	}
	if len(self) != len(got) {
		if len(self) < len(got) {
			return data.NewIntValue(-1), nil
		}
		return data.NewIntValue(1), nil
	}
	return data.NewIntValue(strings.Compare(self, got)), nil
}

func uuidIsNil(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(propertyUID(uidSelf(ctx)) == uuidNil), nil
}

func uuidIsValid(ctx data.Context) (data.GetValue, data.Control) {
	raw, ok := argString(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	format := 1 << 3 // FORMAT_RFC_9562
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if ai, ok := v.(data.AsInt); ok {
			if n, err := ai.AsInt(); err == nil {
				format = n
			}
		}
	}
	rfc, ok := uuidToRFC4122ForFormat(raw, format)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	if rfc == uuidNil || rfc == uuidMax {
		return data.NewBoolValue(true), nil
	}
	if len(rfc) != 36 {
		return data.NewBoolValue(false), nil
	}
	switch rfc[19] {
	case '8', '9', 'a', 'b':
	default:
		return data.NewBoolValue(false), nil
	}
	typ := lateStaticType(ctx)
	if typ > 0 {
		if rfc[14] < '0' || rfc[14] > '9' || int(rfc[14]-'0') != typ {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func uuidToRFC4122ForFormat(raw string, format int) (string, bool) {
	all := format == -1
	if len(raw) == 36 {
		if !all && format&(1<<3) == 0 {
			return "", false
		}
		if !isRFC4122(raw) {
			return "", false
		}
		return strings.ToLower(raw), true
	}
	if len(raw) == 16 && (all || format&1 != 0) {
		return formatRFC4122([]byte(raw)), true
	}
	if ulidIsValid(raw) && (all || format&(1<<1) != 0) {
		b, ok := decodeCrockfordTo16(strings.ToUpper(raw))
		if !ok {
			return "", false
		}
		return formatRFC4122(b[:]), true
	}
	return "", false
}
