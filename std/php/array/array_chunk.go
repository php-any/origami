package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayChunkFunction 实现 array_chunk
// array_chunk(array $array, int $length, bool $preserve_keys = false): array
type ArrayChunkFunction struct{}

func NewArrayChunkFunction() data.FuncStmt {
	return &ArrayChunkFunction{}
}

func (f *ArrayChunkFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	lengthVal, _ := ctx.GetIndexValue(1)
	preserveVal, _ := ctx.GetIndexValue(2)

	length := 1
	if iv, ok := lengthVal.(data.AsInt); ok {
		if v, err := iv.AsInt(); err == nil && v > 0 {
			length = v
		}
	}

	preserve := false
	if bv, ok := preserveVal.(*data.BoolValue); ok {
		preserve = bv.Value
	}

	entries := toKVEntries(arrayVal)
	if len(entries) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	chunks := make([]data.Value, 0)
	for i := 0; i < len(entries); i += length {
		end := i + length
		if end > len(entries) {
			end = len(entries)
		}
		part := entries[i:end]
		if preserve {
			chunks = append(chunks, buildResultFromEntries(part, true))
		} else {
			vals := make([]data.Value, len(part))
			for j, e := range part {
				vals[j] = e.value
			}
			chunks = append(chunks, data.NewArrayValue(vals))
		}
	}

	return data.NewArrayValue(chunks), nil
}

func (f *ArrayChunkFunction) GetName() string { return "array_chunk" }

func (f *ArrayChunkFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "length", 1, nil, nil),
		node.NewParameter(nil, "preserve_keys", 2, data.NewBoolValue(false), nil),
	}
}

func (f *ArrayChunkFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "length", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "preserve_keys", 2, data.NewBaseType("bool")),
	}
}
