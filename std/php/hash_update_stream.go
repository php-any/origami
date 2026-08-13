package php

import (
	"io"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/stream"
	"github.com/php-any/origami/utils"
)

// HashUpdateStreamFunction 实现 PHP hash_update_stream()
type HashUpdateStreamFunction struct{}

func NewHashUpdateStreamFunction() data.FuncStmt {
	return &HashUpdateStreamFunction{}
}

func (f *HashUpdateStreamFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hc, ctrl := getHashContext(ctx, 0)
	if ctrl != nil {
		return nil, ctrl
	}

	streamVal, _ := ctx.GetIndexValue(1)
	if streamVal == nil {
		return nil, utils.NewThrowf("hash_update_stream(): Argument #2 ($stream) must be of type resource")
	}

	reader, ok := resolveStreamReader(streamVal)
	if !ok {
		return nil, utils.NewThrowf("hash_update_stream(): Argument #2 ($stream) must be of type resource")
	}

	length := -1
	if lengthVal, _ := ctx.GetIndexValue(2); lengthVal != nil {
		if _, isNull := lengthVal.(*data.NullValue); !isNull {
			if asInt, ok := lengthVal.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil {
					length = n
				}
			}
		}
	}

	var (
		n   int64
		err error
	)
	if length < 0 {
		n, err = io.Copy(hc.Hash, reader)
	} else {
		n, err = io.CopyN(hc.Hash, reader, int64(length))
		if err == io.EOF {
			err = nil
		}
	}
	if err != nil && err != io.EOF {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(n)), nil
}

func resolveStreamReader(v data.Value) (io.Reader, bool) {
	rv, ok := v.(*core.ResourceValue)
	if !ok {
		return nil, false
	}
	resource := rv.GetResource()
	if resource == nil {
		return nil, false
	}
	switch s := resource.(type) {
	case *stream.StreamInfo:
		return s, true
	case *stream.StreamInfoFromReader:
		if s == nil || s.Reader == nil {
			return nil, false
		}
		return s, true
	case io.Reader:
		return s, true
	default:
		return nil, false
	}
}

func (f *HashUpdateStreamFunction) GetName() string            { return "hash_update_stream" }
func (f *HashUpdateStreamFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashUpdateStreamFunction) GetIsStatic() bool          { return false }
func (f *HashUpdateStreamFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "context", 0, nil, nil),
		node.NewParameter(nil, "stream", 1, nil, nil),
		node.NewParameter(nil, "length", 2, node.NewNullLiteral(nil), nil),
	}
}
func (f *HashUpdateStreamFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "context", 0, nil),
		node.NewVariable(nil, "stream", 1, nil),
		node.NewVariable(nil, "length", 2, nil),
	}
}
func (f *HashUpdateStreamFunction) GetReturnType() data.Types { return data.NewBaseType("int") }
