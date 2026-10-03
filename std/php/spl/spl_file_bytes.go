package spl

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"io"
)

type SFOByteMethod struct{ name string }

func (m *SFOByteMethod) GetName() string          { return m.name }
func (*SFOByteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (*SFOByteMethod) GetIsStatic() bool          { return false }
func (m *SFOByteMethod) GetReturnType() data.Types {
	if m.name == "fread" {
		return data.NewUnionType([]data.Types{data.String{}, data.TypeFalse})
	}
	if m.name == "ftell" {
		return data.NewUnionType([]data.Types{data.Int{}, data.TypeFalse})
	}
	return data.Int{}
}

var sfoReadParams = []data.GetValue{node.NewParameter(nil, "length", 0, nil, data.Int{})}
var sfoReadVariables = []data.Variable{node.NewVariable(nil, "length", 0, data.Int{})}
var sfoSeekParams = []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Int{}), node.NewParameter(nil, "whence", 1, data.NewIntValue(io.SeekStart), data.Int{})}
var sfoSeekVariables = []data.Variable{node.NewVariable(nil, "offset", 0, data.Int{}), node.NewVariable(nil, "whence", 1, data.Int{})}

func (m *SFOByteMethod) GetParams() []data.GetValue {
	if m.name == "fread" {
		return sfoReadParams
	}
	if m.name == "fseek" {
		return sfoSeekParams
	}
	return nil
}
func (m *SFOByteMethod) GetVariables() []data.Variable {
	if m.name == "fread" {
		return sfoReadVariables
	}
	if m.name == "fseek" {
		return sfoSeekVariables
	}
	return nil
}
func (m *SFOByteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state := sfoGetState(sfiGetCV(ctx))
	if state == nil || state.stream == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("SplFileObject is not initialized"), "Error")
	}
	data.CheckRequest(ctx.GoContext())
	defer data.CheckRequest(ctx.GoContext())
	// Iterator lookahead must not move the byte cursor before the first read.
	if state.fgetsAtStart {
		if _, err := state.stream.Seek(0, io.SeekStart); err != nil {
			return data.NewBoolValue(false), nil
		}
		state.fgetsAtStart = false
		state.reader = nil
		state.eof = false
	}
	switch m.name {
	case "fread":
		length := sfoCtxInt(ctx, 0, 0)
		if length <= 0 {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("SplFileObject::fread(): Argument #1 ($length) must be greater than 0"), "ValueError")
		}
		reader := io.Reader(state.stream)
		if state.reader != nil {
			reader = state.reader
		}
		// A limited reader avoids allocating the requested size for tiny files.
		bytes, err := io.ReadAll(io.LimitReader(reader, int64(length)))
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		state.eof = len(bytes) < length
		state.valid = false
		state.current = ""
		return data.NewStringValue(string(bytes)), nil
	case "fseek":
		offset, whence := sfoCtxInt(ctx, 0, 0), sfoCtxInt(ctx, 1, io.SeekStart)
		if whence == io.SeekCurrent && state.reader != nil {
			offset -= state.reader.Buffered()
		}
		if _, err := state.stream.Seek(int64(offset), whence); err != nil {
			return data.NewIntValue(-1), nil
		}
		state.reader = nil
		state.eof = false
		state.valid = false
		state.current = ""
		state.key = 0
		return data.NewIntValue(0), nil
	case "ftell":
		position, err := state.stream.Seek(0, io.SeekCurrent)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		if state.reader != nil {
			position -= int64(state.reader.Buffered())
		}
		return data.NewIntValue(int(position)), nil
	}
	return data.NewBoolValue(false), nil
}
