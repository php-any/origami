// Package protowire implements a generic protobuf binary parser that works
// without .proto files or code generation. It relies entirely on the
// google.golang.org/protobuf/encoding/protowire package for low-level wire
// format operations.
package protowire

import (
	"errors"
	"fmt"

	pw "google.golang.org/protobuf/encoding/protowire"
)

// Encoding helpers re-exported from protowire for use by functions.go.
// These are thin wrappers to keep the encoding API available without
// exposing the external dependency directly to PHP-callable code.

func pwAppendTag(b []byte, num pw.Number, typ pw.Type) []byte { return pw.AppendTag(b, num, typ) }
func pwAppendVarint(b []byte, v uint64) []byte                { return pw.AppendVarint(b, v) }
func pwAppendFixed32(b []byte, v uint32) []byte               { return pw.AppendFixed32(b, v) }
func pwAppendFixed64(b []byte, v uint64) []byte               { return pw.AppendFixed64(b, v) }
func pwAppendBytes(b []byte, v []byte) []byte                 { return pw.AppendBytes(b, v) }
func pwNumber(n int) pw.Number                                { return pw.Number(n) }
func pwType(t int) pw.Type                                    { return pw.Type(t) }
func pwConsumeVarint(b []byte) (uint64, int)                  { return pw.ConsumeVarint(b) }
func pwConsumeFixed32(b []byte) (uint32, int)                 { return pw.ConsumeFixed32(b) }
func pwConsumeFixed64(b []byte) (uint64, int)                 { return pw.ConsumeFixed64(b) }
func pwConsumeBytes(b []byte) ([]byte, int)                   { return pw.ConsumeBytes(b) }
func pwConsumeTag(b []byte) (pw.Number, pw.Type, int)         { return pw.ConsumeTag(b) }

// Wire type constants matching the protobuf wire format specification.
const (
	WireVarint          = 0
	WireFixed64         = 1
	WireLengthDelimited = 2
	WireStartGroup      = 3
	WireEndGroup        = 4
	WireFixed32         = 5
)

// Errors returned by the parser.
var (
	ErrMaxDepth           = errors.New("maximum recursion depth exceeded")
	ErrInvalidVarint      = errors.New("invalid varint encoding")
	ErrInvalidFixed64     = errors.New("invalid fixed64 encoding")
	ErrInvalidFixed32     = errors.New("invalid fixed32 encoding")
	ErrInvalidTag         = errors.New("invalid tag")
	ErrInvalidLength      = errors.New("invalid length-delimited encoding")
	ErrUnexpectedEnd      = errors.New("unexpected end of data")
	ErrUnexpectedEndGroup = errors.New("unexpected end group (not inside a group)")
)

// Field represents a single parsed protobuf field.
//
// The Value field's concrete type depends on the WireType and options:
//   - WireVarint  -> uint64
//   - WireFixed32 -> uint32
//   - WireFixed64 -> uint64
//   - WireLengthDelimited -> []byte (default), []Field (if configured as message),
//     or []uint64/[]uint32 (if configured as packed)
//   - WireStartGroup -> []Field
type Field struct {
	Number   int32
	WireType int32
	Value    interface{}
}

// ParseOptions configures the behavior of ParseRawFields.
type ParseOptions struct {
	// MessageFields specifies which field numbers should be treated as
	// embedded message types and recursively parsed.
	MessageFields map[int32]bool

	// PackedFields specifies which field numbers use packed repeated encoding.
	PackedFields map[int32]bool

	// PackedElementType specifies the wire type of elements inside packed
	// fields. Required for each field listed in PackedFields.
	// Valid values: WireVarint, WireFixed32, WireFixed64.
	PackedElementType map[int32]int32

	// MaxDepth limits the maximum recursion depth for nested messages
	// and groups. Defaults to 64 if zero or negative.
	MaxDepth int
}

// ParseRawFields parses raw protobuf binary data into a slice of Field values.
// It operates solely on the wire format and does not require .proto schema
// information, though ParseOptions can provide optional hints for recursive
// message parsing and packed field decoding.
func ParseRawFields(data []byte, opts *ParseOptions) ([]Field, error) {
	if opts == nil {
		opts = &ParseOptions{}
	}
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = 64
	}
	return parseFields(data, opts, 0)
}

// parseFields recursively parses a sequence of protobuf fields from data.
func parseFields(data []byte, opts *ParseOptions, depth int) ([]Field, error) {
	if depth >= opts.MaxDepth {
		return nil, ErrMaxDepth
	}

	var fields []Field
	for len(data) > 0 {
		num, wtype, n := pw.ConsumeTag(data)
		if n <= 0 {
			return nil, fmt.Errorf("field tag: %w at remaining %d bytes", ErrInvalidTag, len(data))
		}
		data = data[n:]

		if wtype == WireEndGroup {
			// EndGroup at the top level is unexpected; only consumeGroup
			// should encounter it. Break so the caller (consumeGroup) can
			// handle it.
			break
		}

		field := Field{Number: int32(num), WireType: int32(wtype)}
		remaining, err := consumeFieldValue(data, &field, opts, depth)
		if err != nil {
			return nil, fmt.Errorf("field %d (wire type %d): %w", num, wtype, err)
		}
		data = remaining
		fields = append(fields, field)
	}
	return fields, nil
}

// consumeFieldValue reads the value for a single field from data, advancing
// the buffer past the consumed bytes. It returns the remaining data.
func consumeFieldValue(data []byte, field *Field, opts *ParseOptions, depth int) ([]byte, error) {
	num := pw.Number(field.Number)
	wtype := field.WireType

	switch wtype {
	case WireVarint:
		v, n := pw.ConsumeVarint(data)
		if n <= 0 {
			return nil, ErrInvalidVarint
		}
		field.Value = v
		return data[n:], nil

	case WireFixed64:
		v, n := pw.ConsumeFixed64(data)
		if n <= 0 {
			return nil, fmt.Errorf("fixed64: %w", ErrInvalidFixed64)
		}
		field.Value = v
		return data[n:], nil

	case WireLengthDelimited:
		payload, n := pw.ConsumeBytes(data)
		if n <= 0 {
			return nil, fmt.Errorf("length-delimited: %w", ErrInvalidLength)
		}
		data = data[n:]

		// Packed repeated field?
		if opts.PackedFields[int32(num)] {
			elemType, ok := opts.PackedElementType[int32(num)]
			if !ok {
				return nil, fmt.Errorf("packed field %d: PackedElementType not configured", int32(num))
			}
			packed, err := unpackPacked(payload, elemType)
			if err != nil {
				return nil, fmt.Errorf("packed field %d: %w", int32(num), err)
			}
			field.Value = packed
			return data, nil
		}

		// Nested message field?
		if opts.MessageFields[int32(num)] {
			inner, err := parseFields(payload, opts, depth+1)
			if err != nil {
				return nil, fmt.Errorf("nested message: %w", err)
			}
			field.Value = inner
			return data, nil
		}

		// Default: raw bytes
		field.Value = payload
		return data, nil

	case WireStartGroup:
		inner, remaining, err := consumeGroup(data, num, opts, depth)
		if err != nil {
			return nil, fmt.Errorf("group: %w", err)
		}
		field.Value = inner
		return remaining, nil

	case WireEndGroup:
		return nil, ErrUnexpectedEndGroup

	case WireFixed32:
		v, n := pw.ConsumeFixed32(data)
		if n <= 0 {
			return nil, fmt.Errorf("fixed32: %w", ErrInvalidFixed32)
		}
		field.Value = v
		return data[n:], nil

	default:
		return nil, fmt.Errorf("unsupported wire type %d", wtype)
	}
}

// consumeGroup reads fields inside a StartGroup/EndGroup pair. It returns
// the parsed fields and the remaining data after the matching EndGroup tag.
func consumeGroup(data []byte, groupNum pw.Number, opts *ParseOptions, depth int) ([]Field, []byte, error) {
	if depth >= opts.MaxDepth {
		return nil, nil, ErrMaxDepth
	}

	var fields []Field
	for len(data) > 0 {
		num, wtype, n := pw.ConsumeTag(data)
		if n <= 0 {
			return nil, nil, fmt.Errorf("group %d: %w", groupNum, ErrInvalidTag)
		}
		data = data[n:]

		if wtype == WireEndGroup {
			if num != groupNum {
				return nil, nil, fmt.Errorf("group %d: mismatched end group tag number %d", groupNum, num)
			}
			return fields, data, nil
		}

		field := Field{Number: int32(num), WireType: int32(wtype)}
		remaining, err := consumeFieldValue(data, &field, opts, depth+1)
		if err != nil {
			return nil, nil, fmt.Errorf("group %d field %d: %w", groupNum, num, err)
		}
		data = remaining
		fields = append(fields, field)
	}
	return nil, nil, fmt.Errorf("group %d: %w", groupNum, ErrUnexpectedEnd)
}

// unpackPacked decodes a packed repeated field payload into a slice of the
// appropriate Go type based on the element wire type.
func unpackPacked(data []byte, elemWireType int32) (interface{}, error) {
	switch elemWireType {
	case WireVarint:
		var vals []uint64
		for len(data) > 0 {
			v, n := pw.ConsumeVarint(data)
			if n <= 0 {
				return nil, ErrInvalidVarint
			}
			vals = append(vals, v)
			data = data[n:]
		}
		return vals, nil

	case WireFixed32:
		var vals []uint32
		for len(data) > 0 {
			v, n := pw.ConsumeFixed32(data)
			if n <= 0 {
				return nil, ErrInvalidFixed32
			}
			vals = append(vals, v)
			data = data[n:]
		}
		return vals, nil

	case WireFixed64:
		var vals []uint64
		for len(data) > 0 {
			v, n := pw.ConsumeFixed64(data)
			if n <= 0 {
				return nil, ErrInvalidFixed64
			}
			vals = append(vals, v)
			data = data[n:]
		}
		return vals, nil

	default:
		return nil, fmt.Errorf("unsupported packed element wire type %d", elemWireType)
	}
}
