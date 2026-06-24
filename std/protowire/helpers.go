package protowire

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ---------------------------------------------------------------------------
// PHP value → Go type conversions
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
// PHP options → ParseOptions conversion
// ---------------------------------------------------------------------------

func parseOptionsFromPHPValue(val data.Value, opts *ParseOptions) error {
	if obj, ok := val.(*data.ObjectValue); ok {
		return parseOptionsFromObject(obj, opts)
	}
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
// []Field → PHP data.Value conversion
// ---------------------------------------------------------------------------

func fieldsToPHPArray(fields []Field) data.Value {
	entries := make([]data.Value, len(fields))
	for i, f := range fields {
		entries[i] = fieldToPHPValue(f)
	}
	return data.NewArrayValue(entries)
}

func fieldToPHPValue(f Field) data.Value {
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
