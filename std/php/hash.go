package php

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// HashFunction 实现 PHP hash() 函数
type HashFunction struct{}

func NewHashFunction() data.FuncStmt {
	return &HashFunction{}
}

func (f *HashFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	algo, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	dataStr, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	var h hash.Hash
	switch algo {
	case "md5":
		h = md5.New()
	case "sha1", "sha-1":
		h = sha1.New()
	case "sha256", "sha-256", "sha2_256":
		h = sha256.New()
	case "sha512", "sha-512", "sha2_512":
		h = sha512.New()
	case "sha3-256":
		h = sha256.New() // fallback: use sha256
	case "sha3-512":
		h = sha512.New() // fallback: use sha512
	case "xxh3":
		// xxh3 is not available in std lib, fall back to sha256
		h = sha256.New()
	default:
		// default fallback
		h = sha256.New()
	}

	h.Write([]byte(dataStr))
	return data.NewStringValue(hex.EncodeToString(h.Sum(nil))), nil
}

func (f *HashFunction) GetName() string            { return "hash" }
func (f *HashFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashFunction) GetIsStatic() bool          { return false }
func (f *HashFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "algo", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, nil),
	}
}
func (f *HashFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "algo", 0, nil),
		node.NewVariable(nil, "data", 1, nil),
	}
}
func (f *HashFunction) GetReturnType() data.Types { return data.NewBaseType("string") }
