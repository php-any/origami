package protowire

import (
	"reflect"
	"testing"

	pw "google.golang.org/protobuf/encoding/protowire"
)

// ---------------------------------------------------------------------------
// Helpers: encode test data using the protowire Append* functions
// ---------------------------------------------------------------------------

func tag(num int, wtype pw.Type) []byte { return pw.AppendTag(nil, pw.Number(num), wtype) }

func varint(v uint64) []byte { return pw.AppendVarint(nil, v) }

func fixed32(v uint32) []byte { return pw.AppendFixed32(nil, v) }

func fixed64(v uint64) []byte { return pw.AppendFixed64(nil, v) }

func lengthDelimited(payload []byte) []byte { return pw.AppendBytes(nil, payload) }

// ---------------------------------------------------------------------------
// Basic types
// ---------------------------------------------------------------------------

func TestVarintField(t *testing.T) {
	data := append(tag(1, WireVarint), varint(42)...)
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	f := fields[0]
	if f.Number != 1 || f.WireType != WireVarint {
		t.Fatalf("unexpected number/wiretype: %d/%d", f.Number, f.WireType)
	}
	v, ok := f.Value.(uint64)
	if !ok || v != 42 {
		t.Fatalf("expected uint64(42), got %T(%v)", f.Value, f.Value)
	}
}

func TestFixed32Field(t *testing.T) {
	data := append(tag(2, WireFixed32), fixed32(0xDEADBEAF)...)
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	v, ok := fields[0].Value.(uint32)
	if !ok || v != 0xDEADBEAF {
		t.Fatalf("expected uint32(0xDEADBEAF), got %T(%v)", fields[0].Value, fields[0].Value)
	}
}

func TestFixed64Field(t *testing.T) {
	data := append(tag(3, WireFixed64), fixed64(0xCAFEBABEDEADBEAF)...)
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	v, ok := fields[0].Value.(uint64)
	if !ok || v != 0xCAFEBABEDEADBEAF {
		t.Fatalf("expected uint64(0xCAFEBABEDEADBEAF), got %T(%v)", fields[0].Value, fields[0].Value)
	}
}

func TestLengthDelimitedBytes(t *testing.T) {
	payload := []byte("hello protobuf")
	data := append(tag(4, WireLengthDelimited), lengthDelimited(payload)...)
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	v, ok := fields[0].Value.([]byte)
	if !ok || string(v) != "hello protobuf" {
		t.Fatalf("expected []byte(\"hello protobuf\"), got %T(%v)", fields[0].Value, fields[0].Value)
	}
}

// ---------------------------------------------------------------------------
// Multiple fields
// ---------------------------------------------------------------------------

func TestMultipleFields(t *testing.T) {
	var data []byte
	data = append(data, tag(1, WireVarint)...)
	data = append(data, varint(100)...)
	data = append(data, tag(2, WireFixed32)...)
	data = append(data, fixed32(999)...)
	data = append(data, tag(3, WireFixed64)...)
	data = append(data, fixed64(888)...)

	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[0].Value.(uint64) != 100 {
		t.Fatalf("field 1 value mismatch: %v", fields[0].Value)
	}
	if fields[1].Value.(uint32) != 999 {
		t.Fatalf("field 2 value mismatch: %v", fields[1].Value)
	}
	if fields[2].Value.(uint64) != 888 {
		t.Fatalf("field 3 value mismatch: %v", fields[2].Value)
	}
}

// ---------------------------------------------------------------------------
// Empty data
// ---------------------------------------------------------------------------

func TestEmptyData(t *testing.T) {
	fields, err := ParseRawFields(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected 0 fields, got %d", len(fields))
	}
}

// ---------------------------------------------------------------------------
// Nested messages (recursive)
// ---------------------------------------------------------------------------

func TestNestedMessage(t *testing.T) {
	// Inner message: field 1 (varint) = 7
	inner := append(tag(1, WireVarint), varint(7)...)
	// Outer: field 5 (length-delimited) = inner
	data := append(tag(5, WireLengthDelimited), lengthDelimited(inner)...)

	opts := &ParseOptions{
		MessageFields: map[int32]bool{5: true},
	}
	fields, err := ParseRawFields(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	innerFields, ok := fields[0].Value.([]Field)
	if !ok {
		t.Fatalf("expected []Field, got %T", fields[0].Value)
	}
	if len(innerFields) != 1 || innerFields[0].Value.(uint64) != 7 {
		t.Fatalf("inner field mismatch: %+v", innerFields)
	}
}

func TestNestedMessageDefaultBytes(t *testing.T) {
	// If a field is NOT in MessageFields, it stays as raw bytes
	inner := append(tag(1, WireVarint), varint(7)...)
	data := append(tag(5, WireLengthDelimited), lengthDelimited(inner)...)

	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := fields[0].Value.([]byte)
	if !ok {
		t.Fatalf("expected []byte (not recursively parsed), got %T", fields[0].Value)
	}
}

// ---------------------------------------------------------------------------
// Deeply nested messages (depth limit)
// ---------------------------------------------------------------------------

func TestDepthLimit(t *testing.T) {
	opts := &ParseOptions{
		MessageFields: map[int32]bool{1: true},
		MaxDepth:      5,
	}

	// Build: level5 wraps level4 wraps level3 wraps level2 wraps level1
	payload := append(tag(1, WireVarint), varint(99)...)
	for i := 0; i < 5; i++ {
		payload = append(tag(1, WireLengthDelimited), lengthDelimited(payload)...)
	}

	_, err := ParseRawFields(payload, opts)
	if err == nil {
		t.Fatal("expected depth limit error, got nil")
	}
}

// ---------------------------------------------------------------------------
// Groups
// ---------------------------------------------------------------------------

func TestGroup(t *testing.T) {
	// Group (field 10):
	//   field 1 (varint) = 123
	//   field 2 (fixed32) = 456
	// EndGroup
	var data []byte
	data = append(data, tag(10, WireStartGroup)...)
	data = append(data, tag(1, WireVarint)...)
	data = append(data, varint(123)...)
	data = append(data, tag(2, WireFixed32)...)
	data = append(data, fixed32(456)...)
	data = append(data, tag(10, WireEndGroup)...)

	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	inner, ok := fields[0].Value.([]Field)
	if !ok {
		t.Fatalf("expected []Field, got %T", fields[0].Value)
	}
	if len(inner) != 2 {
		t.Fatalf("expected 2 inner fields, got %d", len(inner))
	}
	if inner[0].Value.(uint64) != 123 || inner[1].Value.(uint32) != 456 {
		t.Fatalf("group inner fields mismatch: %+v", inner)
	}
}

func TestNestedGroup(t *testing.T) {
	// Group (field 20):
	//   Group (field 21):
	//     field 1 (varint) = 77
	//   EndGroup(21)
	// EndGroup(20)
	var data []byte
	data = append(data, tag(20, WireStartGroup)...)
	data = append(data, tag(21, WireStartGroup)...)
	data = append(data, tag(1, WireVarint)...)
	data = append(data, varint(77)...)
	data = append(data, tag(21, WireEndGroup)...)
	data = append(data, tag(20, WireEndGroup)...)

	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	outer := fields[0].Value.([]Field)
	inner := outer[0].Value.([]Field)
	if inner[0].Value.(uint64) != 77 {
		t.Fatalf("nested group inner value mismatch: %v", inner[0].Value)
	}
}

// ---------------------------------------------------------------------------
// Packed repeated fields
// ---------------------------------------------------------------------------

func TestPackedVarint(t *testing.T) {
	// Field 7, packed varint: [1, 2, 3, 4, 5]
	var payload []byte
	for _, v := range []uint64{1, 2, 3, 4, 5} {
		payload = append(payload, varint(v)...)
	}
	data := append(tag(7, WireLengthDelimited), lengthDelimited(payload)...)

	opts := &ParseOptions{
		PackedFields:      map[int32]bool{7: true},
		PackedElementType: map[int32]int32{7: WireVarint},
	}
	fields, err := ParseRawFields(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	vals, ok := fields[0].Value.([]uint64)
	if !ok {
		t.Fatalf("expected []uint64, got %T", fields[0].Value)
	}
	if !reflect.DeepEqual(vals, []uint64{1, 2, 3, 4, 5}) {
		t.Fatalf("packed values mismatch: %v", vals)
	}
}

func TestPackedFixed32(t *testing.T) {
	var payload []byte
	for _, v := range []uint32{10, 20, 30} {
		payload = append(payload, fixed32(v)...)
	}
	data := append(tag(8, WireLengthDelimited), lengthDelimited(payload)...)

	opts := &ParseOptions{
		PackedFields:      map[int32]bool{8: true},
		PackedElementType: map[int32]int32{8: WireFixed32},
	}
	fields, err := ParseRawFields(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	vals, ok := fields[0].Value.([]uint32)
	if !ok {
		t.Fatalf("expected []uint32, got %T", fields[0].Value)
	}
	if !reflect.DeepEqual(vals, []uint32{10, 20, 30}) {
		t.Fatalf("packed values mismatch: %v", vals)
	}
}

func TestPackedFixed64(t *testing.T) {
	var payload []byte
	for _, v := range []uint64{100, 200} {
		payload = append(payload, fixed64(v)...)
	}
	data := append(tag(9, WireLengthDelimited), lengthDelimited(payload)...)

	opts := &ParseOptions{
		PackedFields:      map[int32]bool{9: true},
		PackedElementType: map[int32]int32{9: WireFixed64},
	}
	fields, err := ParseRawFields(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	vals, ok := fields[0].Value.([]uint64)
	if !ok {
		t.Fatalf("expected []uint64, got %T", fields[0].Value)
	}
	if !reflect.DeepEqual(vals, []uint64{100, 200}) {
		t.Fatalf("packed values mismatch: %v", vals)
	}
}

// ---------------------------------------------------------------------------
// Repeated non-packed fields (appear as individual fields with same number)
// ---------------------------------------------------------------------------

func TestNonPackedRepeated(t *testing.T) {
	var data []byte
	for i := 0; i < 3; i++ {
		data = append(data, tag(1, WireVarint)...)
		data = append(data, varint(uint64(10+i))...)
	}
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	for i, f := range fields {
		if f.Value.(uint64) != uint64(10+i) {
			t.Fatalf("field %d value mismatch: %v", i, f.Value)
		}
	}
}

// ---------------------------------------------------------------------------
// Unknown fields (any field number, automatically handled)
// ---------------------------------------------------------------------------

func TestUnknownField(t *testing.T) {
	data := append(tag(99, WireVarint), varint(888)...)
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Number != 99 || fields[0].Value.(uint64) != 888 {
		t.Fatalf("unknown field mismatch: %+v", fields[0])
	}
}

// ---------------------------------------------------------------------------
// Malformed data
// ---------------------------------------------------------------------------

func TestMalformedTag(t *testing.T) {
	// Truncated tag (incomplete varint)
	data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}
	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for malformed tag, got nil")
	}
}

func TestMalformedVarintValue(t *testing.T) {
	// 12 continuation bytes with no terminator — exceeds max varint length
	data := append(tag(1, WireVarint), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)
	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for truncated varint value, got nil")
	}
}

func TestTruncatedLengthDelimited(t *testing.T) {
	// Valid tag and varint-length, but payload is truncated
	data := append(tag(3, WireLengthDelimited), varint(100)...) // claims 100 bytes but none follow
	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for truncated length-delimited, got nil")
	}
}

func TestTruncatedFixed32(t *testing.T) {
	// Valid tag but only 3 of 4 bytes
	data := append(tag(4, WireFixed32), 0x01, 0x02, 0x03)
	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for truncated fixed32, got nil")
	}
}

func TestTruncatedFixed64(t *testing.T) {
	// Valid tag but only 4 of 8 bytes
	data := append(tag(5, WireFixed64), 0x01, 0x02, 0x03, 0x04)
	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for truncated fixed64, got nil")
	}
}

// ---------------------------------------------------------------------------
// Mixed: combination of all wire types
// ---------------------------------------------------------------------------

func TestMixedFields(t *testing.T) {
	var data []byte
	data = append(data, tag(1, WireVarint)...)
	data = append(data, varint(1)...)
	data = append(data, tag(2, WireFixed64)...)
	data = append(data, fixed64(2)...)
	data = append(data, tag(3, WireLengthDelimited)...)
	data = append(data, lengthDelimited([]byte("three"))...)
	data = append(data, tag(4, WireFixed32)...)
	data = append(data, fixed32(4)...)

	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(fields))
	}
	if fields[0].Value.(uint64) != 1 {
		t.Fatalf("field 1: expected 1, got %v", fields[0].Value)
	}
	if fields[1].Value.(uint64) != 2 {
		t.Fatalf("field 2: expected 2, got %v", fields[1].Value)
	}
	if string(fields[2].Value.([]byte)) != "three" {
		t.Fatalf("field 3: expected \"three\", got %v", fields[2].Value)
	}
	if fields[3].Value.(uint32) != 4 {
		t.Fatalf("field 4: expected 4, got %v", fields[3].Value)
	}
}

// ---------------------------------------------------------------------------
// Group with mismatched end tag
// ---------------------------------------------------------------------------

func TestGroupMismatchedEndTag(t *testing.T) {
	var data []byte
	data = append(data, tag(10, WireStartGroup)...)
	data = append(data, tag(1, WireVarint)...)
	data = append(data, varint(1)...)
	data = append(data, tag(99, WireEndGroup)...) // wrong field number

	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for mismatched end group tag, got nil")
	}
}

// ---------------------------------------------------------------------------
// Group with no matching end tag (truncated)
// ---------------------------------------------------------------------------

func TestGroupTruncated(t *testing.T) {
	data := append(tag(10, WireStartGroup), tag(1, WireVarint)...)
	data = append(data, varint(1)...)
	// missing EndGroup
	_, err := ParseRawFields(data, nil)
	if err == nil {
		t.Fatal("expected error for truncated group, got nil")
	}
}

// ---------------------------------------------------------------------------
// Nil opts (should use defaults)
// ---------------------------------------------------------------------------

func TestNilOpts(t *testing.T) {
	data := append(tag(1, WireVarint), varint(42)...)
	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
}

// ---------------------------------------------------------------------------
// Depth limit for groups
// ---------------------------------------------------------------------------

func TestGroupDepthLimit(t *testing.T) {
	opts := &ParseOptions{MaxDepth: 2}

	// Nested groups: group(10) -> group(11) -> should fail at depth 3
	var data []byte
	data = append(data, tag(10, WireStartGroup)...)
	data = append(data, tag(11, WireStartGroup)...)
	data = append(data, tag(1, WireVarint)...)
	data = append(data, varint(1)...)
	data = append(data, tag(11, WireEndGroup)...)
	// omit EndGroup(10) since we expect failure before getting there

	_, err := ParseRawFields(data, opts)
	if err == nil {
		t.Fatal("expected depth limit error, got nil")
	}
}

// ---------------------------------------------------------------------------
// Message + group + packed all together
// ---------------------------------------------------------------------------

func TestComplexNested(t *testing.T) {
	opts := &ParseOptions{
		MessageFields:     map[int32]bool{1: true},
		PackedFields:      map[int32]bool{3: true},
		PackedElementType: map[int32]int32{3: WireVarint},
		MaxDepth:          32,
	}

	// Inner message: field 2 (group[10]: varint=9), field 3 (packed: [1,2,3])
	var inner []byte
	// group
	inner = append(inner, tag(10, WireStartGroup)...)
	inner = append(inner, tag(1, WireVarint)...)
	inner = append(inner, varint(9)...)
	inner = append(inner, tag(10, WireEndGroup)...)
	// packed
	var packedPayload []byte
	for _, v := range []uint64{1, 2, 3} {
		packedPayload = append(packedPayload, varint(v)...)
	}
	inner = append(inner, tag(3, WireLengthDelimited)...)
	inner = append(inner, lengthDelimited(packedPayload)...)

	// Outer: field 1 (message) = inner
	data := append(tag(1, WireLengthDelimited), lengthDelimited(inner)...)

	fields, err := ParseRawFields(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	innerFields := fields[0].Value.([]Field)
	if len(innerFields) != 2 {
		t.Fatalf("expected 2 inner fields, got %d", len(innerFields))
	}
	// Check group
	groupFields := innerFields[0].Value.([]Field)
	if len(groupFields) != 1 || groupFields[0].Value.(uint64) != 9 {
		t.Fatalf("group inner mismatch: %+v", groupFields)
	}
	// Check packed
	packedVals := innerFields[1].Value.([]uint64)
	if !reflect.DeepEqual(packedVals, []uint64{1, 2, 3}) {
		t.Fatalf("packed values mismatch: %v", packedVals)
	}
}

// ---------------------------------------------------------------------------
// Empty packed field
// ---------------------------------------------------------------------------

func TestEmptyPacked(t *testing.T) {
	// Empty packed field: length-delimited with zero-length payload
	data := append(tag(5, WireLengthDelimited), lengthDelimited(nil)...)

	opts := &ParseOptions{
		PackedFields:      map[int32]bool{5: true},
		PackedElementType: map[int32]int32{5: WireVarint},
	}
	fields, err := ParseRawFields(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	vals, ok := fields[0].Value.([]uint64)
	if !ok || len(vals) != 0 {
		t.Fatalf("expected empty []uint64, got %T(%v)", fields[0].Value, fields[0].Value)
	}
}

// ---------------------------------------------------------------------------
// Large field numbers
// ---------------------------------------------------------------------------

func TestLargeFieldNumber(t *testing.T) {
	// Field number 1<<29 - 4 (max valid field number), wire type varint, value 1
	num := pw.Number(1<<29 - 4)
	data := pw.AppendTag(nil, num, WireVarint)
	data = append(data, varint(1)...)

	fields, err := ParseRawFields(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || int32(num) != fields[0].Number {
		t.Fatalf("expected field number %d, got %d", num, fields[0].Number)
	}
}
