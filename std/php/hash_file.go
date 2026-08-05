package php

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// HashFileFunction 实现 PHP hash_file()
type HashFileFunction struct{}

func NewHashFileFunction() data.FuncStmt {
	return &HashFileFunction{}
}

func (f *HashFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	algo, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	filename, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	h, err := newHashByAlgo(algo)
	if err != nil {
		setLastError(2 /* E_WARNING */, err.Error(), "", 0)
		return data.NewBoolValue(false), nil
	}

	file, err := os.Open(filename)
	if err != nil {
		setLastError(2 /* E_WARNING */, "hash_file("+filename+"): Failed to open stream: "+err.Error(), "", 0)
		return data.NewBoolValue(false), nil
	}
	defer file.Close()

	if _, err := io.Copy(h, file); err != nil {
		setLastError(2 /* E_WARNING */, err.Error(), "", 0)
		return data.NewBoolValue(false), nil
	}

	rawOutput := false
	if rawVal, _ := ctx.GetIndexValue(2); rawVal != nil {
		if _, isNull := rawVal.(*data.NullValue); !isNull {
			if asBool, ok := rawVal.(data.AsBool); ok {
				if b, err := asBool.AsBool(); err == nil {
					rawOutput = b
				}
			}
		}
	}

	sum := h.Sum(nil)
	if rawOutput {
		return data.NewStringValue(string(sum)), nil
	}
	return data.NewStringValue(hex.EncodeToString(sum)), nil
}

func (f *HashFileFunction) GetName() string            { return "hash_file" }
func (f *HashFileFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashFileFunction) GetIsStatic() bool          { return false }
func (f *HashFileFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "algo", 0, nil, nil),
		node.NewParameter(nil, "filename", 1, nil, nil),
		node.NewParameter(nil, "binary", 2, node.NewNullLiteral(nil), nil),
	}
}
func (f *HashFileFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "algo", 0, nil),
		node.NewVariable(nil, "filename", 1, nil),
		node.NewVariable(nil, "binary", 2, nil),
	}
}
func (f *HashFileFunction) GetReturnType() data.Types {
	return data.NewUnionType([]data.Types{data.NewBaseType("string"), data.NewBaseType("bool")})
}
