package protowire

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ---------------------------------------------------------------------------
// PHP-callable function: protowire_parse_raw_fields
// ---------------------------------------------------------------------------

// ParseRawFieldsFunction wraps the protobuf parser as a PHP-callable function.
//
// PHP signature:
//
//	protowire_parse_raw_fields(string $data, array $options = []): array
//
// Each returned entry is an associative array:
//
//	['number' => int, 'wire_type' => int, 'value' => mixed]
type ParseRawFieldsFunction struct{}

func NewParseRawFieldsFunction() data.FuncStmt {
	return &ParseRawFieldsFunction{}
}

func (f *ParseRawFieldsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// arg 0: raw binary data string
	dataVal, ok := ctx.GetIndexValue(0)
	if !ok || dataVal == nil {
		return data.NewArrayValue(nil), nil
	}

	inputStr, ok := dataVal.(data.AsString)
	if !ok {
		return data.NewArrayValue(nil), nil
	}
	input := []byte(inputStr.AsString())

	// arg 1: optional options array
	opts := &ParseOptions{}
	if optsVal, ok := ctx.GetIndexValue(1); ok && optsVal != nil {
		if err := parseOptionsFromPHPValue(optsVal, opts); err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("protowire_parse_raw_fields: %w", err))
		}
	}

	fields, err := ParseRawFields(input, opts)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("protowire_parse_raw_fields: %w", err))
	}

	result := fieldsToPHPArray(fields)
	return result, nil
}

func (f *ParseRawFieldsFunction) GetName() string {
	return "protowire_parse_raw_fields"
}

func (f *ParseRawFieldsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, data.String{}),
		node.NewParameter(nil, "options", 1, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (f *ParseRawFieldsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.String{}),
		node.NewVariable(nil, "options", 1, data.NewBaseType("array")),
	}
}

// ---------------------------------------------------------------------------
// PHP-callable function: protowire_encode_varint
// ---------------------------------------------------------------------------

// EncodeVarintFunction encodes an integer as protobuf varint bytes.
type EncodeVarintFunction struct{}

func NewEncodeVarintFunction() data.FuncStmt {
	return &EncodeVarintFunction{}
}

func (f *EncodeVarintFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	n, err := toUint64(val)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendVarint(nil, n)
	return data.NewStringValue(string(buf)), nil
}

func (f *EncodeVarintFunction) GetName() string {
	return "protowire_encode_varint"
}

func (f *EncodeVarintFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *EncodeVarintFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

// ---------------------------------------------------------------------------
// PHP-callable function: protowire_encode_tag
// ---------------------------------------------------------------------------

type EncodeTagFunction struct{}

func NewEncodeTagFunction() data.FuncStmt {
	return &EncodeTagFunction{}
}

func (f *EncodeTagFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	numVal, _ := ctx.GetIndexValue(0)
	wtypeVal, _ := ctx.GetIndexValue(1)
	if numVal == nil || wtypeVal == nil {
		return data.NewStringValue(""), nil
	}
	num, err := toUint64(numVal)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	wtype, err := toInt32(wtypeVal)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendTag(nil, pwNumber(int(num)), pwType(int(wtype)))
	return data.NewStringValue(string(buf)), nil
}

func (f *EncodeTagFunction) GetName() string {
	return "protowire_encode_tag"
}

func (f *EncodeTagFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "number", 0, nil, nil),
		node.NewParameter(nil, "wire_type", 1, nil, nil),
	}
}

func (f *EncodeTagFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "number", 0, nil),
		node.NewVariable(nil, "wire_type", 1, nil),
	}
}

// ---------------------------------------------------------------------------
// PHP-callable function: protowire_encode_bytes
// ---------------------------------------------------------------------------

type EncodeBytesFunction struct{}

func NewEncodeBytesFunction() data.FuncStmt {
	return &EncodeBytesFunction{}
}

func (f *EncodeBytesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	s, ok := val.(data.AsString)
	if !ok {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendBytes(nil, []byte(s.AsString()))
	return data.NewStringValue(string(buf)), nil
}

func (f *EncodeBytesFunction) GetName() string {
	return "protowire_encode_bytes"
}

func (f *EncodeBytesFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *EncodeBytesFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

// ---------------------------------------------------------------------------
// PHP-callable function: protowire_encode_fixed32
// ---------------------------------------------------------------------------

type EncodeFixed32Function struct{}

func NewEncodeFixed32Function() data.FuncStmt {
	return &EncodeFixed32Function{}
}

func (f *EncodeFixed32Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	n, err := toUint64(val)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendFixed32(nil, uint32(n))
	return data.NewStringValue(string(buf)), nil
}

func (f *EncodeFixed32Function) GetName() string {
	return "protowire_encode_fixed32"
}

func (f *EncodeFixed32Function) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *EncodeFixed32Function) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

// ---------------------------------------------------------------------------
// PHP-callable function: protowire_encode_fixed64
// ---------------------------------------------------------------------------

type EncodeFixed64Function struct{}

func NewEncodeFixed64Function() data.FuncStmt {
	return &EncodeFixed64Function{}
}

func (f *EncodeFixed64Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	n, err := toUint64(val)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendFixed64(nil, n)
	return data.NewStringValue(string(buf)), nil
}

func (f *EncodeFixed64Function) GetName() string {
	return "protowire_encode_fixed64"
}

func (f *EncodeFixed64Function) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *EncodeFixed64Function) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

// ---------------------------------------------------------------------------
// Helper: convert PHP values to Go types
// ---------------------------------------------------------------------------

func toUint64(v data.Value) (uint64, error) {
	switch tv := v.(type) {
	case data.AsInt:
		n, err := tv.AsInt()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case data.AsString:
		// try parse
		s := tv.AsString()
		var n uint64
		_, err := fmt.Sscanf(s, "%d", &n)
		return n, err
	default:
		return 0, fmt.Errorf("cannot convert %T to uint64", v)
	}
}

func toInt32(v data.Value) (int32, error) {
	switch tv := v.(type) {
	case data.AsInt:
		n, err := tv.AsInt()
		if err != nil {
			return 0, err
		}
		return int32(n), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int32", v)
	}
}

// ---------------------------------------------------------------------------
// Helper: parse PHP options array into ParseOptions
// ---------------------------------------------------------------------------

func parseOptionsFromPHPValue(val data.Value, opts *ParseOptions) error {
	// Handle ObjectValue (associative array from PHP)
	if obj, ok := val.(*data.ObjectValue); ok {
		return parseOptionsFromObject(obj, opts)
	}
	// Handle ArrayValue (could also be passed as options)
	if arr, ok := val.(*data.ArrayValue); ok {
		return parseOptionsFromArrayValue(arr, opts)
	}
	return nil
}

func parseOptionsFromObject(obj *data.ObjectValue, opts *ParseOptions) error {
	if msgFields, acl := obj.GetProperty("message_fields"); acl == nil && msgFields != nil {
		m, err := boolMapFromPHPArray(msgFields)
		if err != nil {
			return fmt.Errorf("options.message_fields: %w", err)
		}
		opts.MessageFields = m
	}
	if pFields, acl := obj.GetProperty("packed_fields"); acl == nil && pFields != nil {
		m, err := boolMapFromPHPArray(pFields)
		if err != nil {
			return fmt.Errorf("options.packed_fields: %w", err)
		}
		opts.PackedFields = m
	}
	if elemTypes, acl := obj.GetProperty("packed_element_type"); acl == nil && elemTypes != nil {
		m, err := intMapFromPHPArray(elemTypes)
		if err != nil {
			return fmt.Errorf("options.packed_element_type: %w", err)
		}
		opts.PackedElementType = m
	}
	if maxDepth, acl := obj.GetProperty("max_depth"); acl == nil && maxDepth != nil {
		if ai, ok := maxDepth.(data.AsInt); ok {
			n, err := ai.AsInt()
			if err == nil {
				opts.MaxDepth = int(n)
			}
		}
	}
	return nil
}

func parseOptionsFromArrayValue(arr *data.ArrayValue, opts *ParseOptions) error {
	for _, z := range arr.List {
		if z == nil || z.Name == "" {
			continue
		}
		switch z.Name {
		case "message_fields":
			if z.Value == nil {
				continue
			}
			m, err := boolMapFromPHPArray(z.Value)
			if err != nil {
				return fmt.Errorf("options.message_fields: %w", err)
			}
			opts.MessageFields = m
		case "packed_fields":
			if z.Value == nil {
				continue
			}
			m, err := boolMapFromPHPArray(z.Value)
			if err != nil {
				return fmt.Errorf("options.packed_fields: %w", err)
			}
			opts.PackedFields = m
		case "packed_element_type":
			if z.Value == nil {
				continue
			}
			m, err := intMapFromPHPArray(z.Value)
			if err != nil {
				return fmt.Errorf("options.packed_element_type: %w", err)
			}
			opts.PackedElementType = m
		case "max_depth":
			if ai, ok := z.Value.(data.AsInt); ok {
				n, err := ai.AsInt()
				if err == nil {
					opts.MaxDepth = int(n)
				}
			}
		}
	}
	return nil
}

func boolMapFromPHPArray(val data.Value) (map[int32]bool, error) {
	result := make(map[int32]bool)
	if obj, ok := val.(*data.ObjectValue); ok {
		obj.RangeProperties(func(key string, v data.Value) bool {
			var n int
			if _, err := fmt.Sscanf(key, "%d", &n); err == nil {
				if bv, ok := v.(*data.BoolValue); ok {
					result[int32(n)] = bv.Value
				} else {
					// non-false value = true
					result[int32(n)] = true
				}
			}
			return true
		})
		return result, nil
	}
	if arr, ok := val.(*data.ArrayValue); ok {
		for _, z := range arr.List {
			if z == nil {
				continue
			}
			if z.Name != "" {
				var n int
				if _, err := fmt.Sscanf(z.Name, "%d", &n); err == nil {
					if bv, ok := z.Value.(*data.BoolValue); ok {
						result[int32(n)] = bv.Value
					} else {
						result[int32(n)] = true
					}
				}
			}
		}
		return result, nil
	}
	return result, nil
}

func intMapFromPHPArray(val data.Value) (map[int32]int32, error) {
	result := make(map[int32]int32)
	if obj, ok := val.(*data.ObjectValue); ok {
		obj.RangeProperties(func(key string, v data.Value) bool {
			var n int
			if _, err := fmt.Sscanf(key, "%d", &n); err == nil {
				if ai, ok := v.(data.AsInt); ok {
					if iv, err := ai.AsInt(); err == nil {
						result[int32(n)] = int32(iv)
					}
				}
			}
			return true
		})
		return result, nil
	}
	if arr, ok := val.(*data.ArrayValue); ok {
		for _, z := range arr.List {
			if z == nil || z.Name == "" {
				continue
			}
			var n int
			if _, err := fmt.Sscanf(z.Name, "%d", &n); err == nil {
				if ai, ok := z.Value.(data.AsInt); ok {
					if iv, err := ai.AsInt(); err == nil {
						result[int32(n)] = int32(iv)
					}
				}
			}
		}
		return result, nil
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Helper: convert []Field to PHP array
// ---------------------------------------------------------------------------

func fieldsToPHPArray(fields []Field) data.Value {
	entries := make([]data.Value, len(fields))
	for i, f := range fields {
		entries[i] = fieldToPHPValue(f)
	}
	return data.NewArrayValue(entries)
}

func fieldToPHPValue(f Field) data.Value {
	// Build an associative array entry
	props := make(map[string]data.Value)
	props["number"] = data.NewIntValue(int(f.Number))
	props["wire_type"] = data.NewIntValue(int(f.WireType))
	props["value"] = goValueToPHPValue(f.Value)
	return objectValueFromMap(props)
}

func goValueToPHPValue(v interface{}) data.Value {
	switch tv := v.(type) {
	case uint64:
		return data.NewIntValue(int(tv))
	case uint32:
		return data.NewIntValue(int(tv))
	case int32:
		return data.NewIntValue(int(tv))
	case int:
		return data.NewIntValue(tv)
	case []byte:
		return data.NewStringValue(string(tv))
	case string:
		return data.NewStringValue(tv)
	case []Field:
		return fieldsToPHPArray(tv)
	case []uint64:
		vals := make([]data.Value, len(tv))
		for i, n := range tv {
			vals[i] = data.NewIntValue(int(n))
		}
		return data.NewArrayValue(vals)
	case []uint32:
		vals := make([]data.Value, len(tv))
		for i, n := range tv {
			vals[i] = data.NewIntValue(int(n))
		}
		return data.NewArrayValue(vals)
	default:
		return data.NewNullValue()
	}
}

func objectValueFromMap(m map[string]data.Value) data.Value {
	obj := data.NewObjectValue()
	for k, v := range m {
		obj.SetProperty(k, v)
	}
	return obj
}
