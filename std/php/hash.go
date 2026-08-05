package php

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"sync/atomic"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/utils"
)

var nextHashContextID int64 = 200000

func allocHashContextID() int {
	return int(atomic.AddInt64(&nextHashContextID, 1))
}

// HashContext 有状态哈希上下文（hash_init / hash_update* / hash_final）
type HashContext struct {
	Hash hash.Hash
	Algo string
}

func newHashByAlgo(algo string) (hash.Hash, error) {
	switch algo {
	case "md5":
		return md5.New(), nil
	case "sha1", "sha-1":
		return sha1.New(), nil
	case "sha256", "sha-256", "sha2_256":
		return sha256.New(), nil
	case "sha512", "sha-512", "sha2_512":
		return sha512.New(), nil
	case "sha3-256":
		return sha256.New(), nil // fallback
	case "sha3-512":
		return sha512.New(), nil // fallback
	case "xxh3":
		return sha256.New(), nil // fallback
	default:
		return nil, fmt.Errorf("Unknown hashing algorithm: %s", algo)
	}
}

func getHashContext(ctx data.Context, index int) (*HashContext, data.Control) {
	v, _ := ctx.GetIndexValue(index)
	if v == nil {
		return nil, utils.NewThrowf("hash context is required")
	}
	rv, ok := v.(*core.ResourceValue)
	if !ok {
		return nil, utils.NewThrowf("supplied resource is not a valid Hash Context resource")
	}
	hc, ok := rv.GetResource().(*HashContext)
	if !ok || hc == nil || hc.Hash == nil {
		return nil, utils.NewThrowf("supplied resource is not a valid Hash Context resource")
	}
	return hc, nil
}

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

	h, err := newHashByAlgo(algo)
	if err != nil {
		// 与旧行为兼容：未知算法回退到 sha256
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

// HashHmacFunction 实现 hash_hmac。
type HashHmacFunction struct{}

func NewHashHmacFunction() data.FuncStmt { return &HashHmacFunction{} }

func (f *HashHmacFunction) GetName() string { return "hash_hmac" }

func (f *HashHmacFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "algo", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "data", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "key", 2, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "binary", 3, node.NewBooleanLiteral(nil, false), data.NewBaseType("bool")),
	}
}

func (f *HashHmacFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "algo", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "data", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "key", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "binary", 3, data.NewBaseType("bool")),
	}
}

func (f *HashHmacFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	algo, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	dataStr, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	key, err := utils.ConvertFromIndex[string](ctx, 2)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	binary := false
	if v, ok := ctx.GetIndexValue(3); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			binary, _ = b.AsBool()
		}
	}

	var hasher func() hash.Hash
	switch algo {
	case "md5":
		hasher = md5.New
	case "sha1", "sha-1":
		hasher = sha1.New
	case "sha256", "sha-256", "sha2_256":
		hasher = sha256.New
	case "sha512", "sha-512", "sha2_512":
		hasher = sha512.New
	default:
		hasher = sha256.New
	}
	mac := hmac.New(hasher, []byte(key))
	_, _ = mac.Write([]byte(dataStr))
	sum := mac.Sum(nil)
	if binary {
		return data.NewStringValue(string(sum)), nil
	}
	return data.NewStringValue(hex.EncodeToString(sum)), nil
}
