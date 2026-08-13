package protowire

import (
	"fmt"
	"math"

	pw "google.golang.org/protobuf/encoding/protowire"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SerializeMethod implements Protowire::serialize()
// Serialize a PHP object to protobuf binary using @Field annotations on class properties.
type SerializeMethod struct{}

func NewSerializeMethod() data.Method {
	return &SerializeMethod{}
}

func (m *SerializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	instanceVal, ok := ctx.GetIndexValue(0)
	if !ok || instanceVal == nil {
		return data.NewStringValue(""), nil
	}

	classVal, ok := instanceVal.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::serialize: parameter must be an object"))
	}

	fieldPlans := readFieldPlans(classVal)

	buf, ctl := encodeFieldPlans(fieldPlans)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(string(buf)), nil
}

func (m *SerializeMethod) GetName() string {
	return "serialize"
}

func (m *SerializeMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *SerializeMethod) GetIsStatic() bool {
	return true
}

func (m *SerializeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "instance", 0, nil, nil),
	}
}

func (m *SerializeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "instance", 0, nil),
	}
}

func (m *SerializeMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

// ---------------------------------------------------------------------------
// fieldPlan holds a fully resolved encode plan for one field.
// ---------------------------------------------------------------------------

type fieldPlan struct {
	num      int32
	wtype    int32
	value    data.Value
	encoding string // "", "float", "double", "zigzag"
}

// readFieldPlans reads @Field annotations from a ClassValue and builds encode plans.
func readFieldPlans(classVal *data.ClassValue) []fieldPlan {
	var plans []fieldPlan
	for _, prop := range classVal.Class.GetPropertyList() {
		cp, ok := prop.(*node.ClassProperty)
		if !ok {
			continue
		}
		num, wtype, encoding := getFieldAnnotationFull(cp)
		if num < 0 {
			continue
		}
		propVal, acl := classVal.GetProperty(cp.GetName())
		if acl != nil || propVal == nil {
			continue
		}
		plans = append(plans, fieldPlan{num: num, wtype: wtype, value: propVal, encoding: encoding})
	}
	return plans
}

// encodeFieldPlans encodes a slice of fieldPlan into protobuf binary.
func encodeFieldPlans(plans []fieldPlan) ([]byte, data.Control) {
	var buf []byte
	for _, p := range plans {
		buf = pwAppendTag(buf, pwNumber(int(p.num)), pwType(int(p.wtype)))
		var err error
		buf, err = encodeFieldValue(buf, p)
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("Protowire::serialize: field %d: %w", p.num, err))
		}
	}
	return buf, nil
}

// encodeFieldValue appends the encoded value for a single fieldPlan.
func encodeFieldValue(buf []byte, p fieldPlan) ([]byte, error) {
	switch p.wtype {
	case WireVarint:
		n, err := toUint64(p.value)
		if err != nil {
			return nil, err
		}
		if p.encoding == "zigzag" {
			n = uint64(pw.EncodeZigZag(int64(n)))
		}
		return pwAppendVarint(buf, n), nil

	case WireFixed64:
		if p.encoding == "double" {
			f, err := toFloat64(p.value)
			if err != nil {
				return nil, err
			}
			return pwAppendFixed64(buf, math.Float64bits(f)), nil
		}
		n, err := toUint64(p.value)
		if err != nil {
			return nil, err
		}
		return pwAppendFixed64(buf, n), nil

	case WireFixed32:
		if p.encoding == "float" {
			f, err := toFloat64(p.value) // PHP float is float64; cast to float32
			if err != nil {
				return nil, err
			}
			return pwAppendFixed32(buf, math.Float32bits(float32(f))), nil
		}
		n, err := toUint64(p.value)
		if err != nil {
			return nil, err
		}
		return pwAppendFixed32(buf, uint32(n)), nil

	case WireLengthDelimited:
		return encodeLengthDelimited(buf, p)

	case WireStartGroup:
		// Recursively serialize nested sub-object or array of fields
		inner, err := encodeGroupContent(p)
		if err != nil {
			return nil, err
		}
		buf = append(buf, inner...)
		buf = pwAppendTag(buf, pwNumber(int(p.num)), pwType(WireEndGroup))
		return buf, nil

	default:
		return nil, fmt.Errorf("unsupported wire type %d", p.wtype)
	}
}

// encodeLengthDelimited handles string, bytes, nested message, and packed repeated.
func encodeLengthDelimited(buf []byte, p fieldPlan) ([]byte, error) {
	switch p.encoding {
	case "packed":
		return encodePacked(buf, p)
	case "message":
		if cv, ok := p.value.(*data.ClassValue); ok {
			innerPlans := readFieldPlans(cv)
			inner, ctl := encodeFieldPlans(innerPlans)
			if ctl != nil {
				return nil, fmt.Errorf("nested message: %s", ctl.AsString())
			}
			return pwAppendBytes(buf, inner), nil
		}
		// Fall through to string encoding
	}
	// Default: string / bytes
	if s, ok := p.value.(data.AsString); ok {
		return pwAppendBytes(buf, []byte(s.AsString())), nil
	}
	return pwAppendBytes(buf, []byte(p.value.AsString())), nil
}

// encodePacked encodes a PHP array as packed repeated wire format.
func encodePacked(buf []byte, p fieldPlan) ([]byte, error) {
	arr, ok := p.value.(*data.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("packed field requires array value, got %T", p.value)
	}
	var payload []byte
	for _, z := range arr.List {
		if z == nil {
			continue
		}
		n, err := toUint64(z.Value)
		if err != nil {
			return nil, fmt.Errorf("packed element: %w", err)
		}
		payload = pwAppendVarint(payload, n)
	}
	return pwAppendBytes(buf, payload), nil
}

// encodeGroupContent serializes a group's inner fields.
// Value must be a ClassValue (object with @Field annotations).
func encodeGroupContent(p fieldPlan) ([]byte, error) {
	cv, ok := p.value.(*data.ClassValue)
	if !ok {
		return nil, fmt.Errorf("group value must be an object, got %T", p.value)
	}
	innerPlans := readFieldPlans(cv)
	inner, ctl := encodeFieldPlans(innerPlans)
	if ctl != nil {
		return nil, fmt.Errorf("group content: %s", ctl.AsString())
	}
	return inner, nil
}

// ---------------------------------------------------------------------------
// Float conversion helpers
// ---------------------------------------------------------------------------

func toFloat64(v data.Value) (float64, error) {
	switch tv := v.(type) {
	case data.AsInt:
		n, err := tv.AsInt()
		if err != nil {
			return 0, err
		}
		return float64(n), nil
	case data.AsString:
		s := tv.AsString()
		var f float64
		_, err := fmt.Sscanf(s, "%f", &f)
		return f, err
	default:
		// Fallback: try AsFloat if available
		if af, ok := v.(interface{ AsFloat() (float64, error) }); ok {
			return af.AsFloat()
		}
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// ---------------------------------------------------------------------------
// Full annotation extraction with encoding hint
// ---------------------------------------------------------------------------

// getFieldAnnotationFull extracts field number, wire type and optional encoding hint
// from a ClassProperty's @Field annotation.
// Returns (-1, 0, "") if no @Field annotation found.
func getFieldAnnotationFull(cp *node.ClassProperty) (int32, int32, string) {
	for _, ann := range cp.Annotations {
		if ann == nil || ann.Class == nil {
			continue
		}
		if ann.Class.GetName() != "Protowire\\Annotation\\Field" {
			continue
		}
		props := ann.GetProperties()
		var num int32 = -1
		var wtype int32
		var encoding string
		if nv, ok := props["number"]; ok && nv != nil {
			if ai, ok := nv.(data.AsInt); ok {
				if n, err := ai.AsInt(); err == nil {
					num = int32(n)
				}
			}
		}
		if tv, ok := props["type"]; ok && tv != nil {
			if ai, ok := tv.(data.AsInt); ok {
				if t, err := ai.AsInt(); err == nil {
					wtype = int32(t)
				}
			}
		}
		if ev, ok := props["encoding"]; ok && ev != nil {
			if s, ok := ev.(data.AsString); ok {
				encoding = s.AsString()
			}
		}
		if num >= 0 {
			return num, wtype, encoding
		}
	}
	return -1, 0, ""
}
