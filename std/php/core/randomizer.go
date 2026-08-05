package core

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RandomizerClass 实现 PHP 8.2 Random\Randomizer 中 Laravel/Symfony 使用的 API。
type RandomizerClass struct {
	node.Node
}

func (c *RandomizerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *RandomizerClass) GetName() string                          { return `Random\Randomizer` }
func (c *RandomizerClass) GetExtend() *string                       { return nil }
func (c *RandomizerClass) GetImplements() []string                  { return nil }
func (c *RandomizerClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *RandomizerClass) GetPropertyList() []data.Property         { return nil }
func (c *RandomizerClass) GetConstruct() data.Method                { return nil }

func (c *RandomizerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "pickArrayKeys":
		return &randomizerPickArrayKeysMethod{}, true
	case "shuffleArray":
		return &randomizerShuffleArrayMethod{}, true
	case "getBytesFromString":
		return &randomizerGetBytesFromStringMethod{}, true
	default:
		return nil, false
	}
}

func (c *RandomizerClass) GetMethods() []data.Method {
	return []data.Method{
		&randomizerPickArrayKeysMethod{},
		&randomizerShuffleArrayMethod{},
		&randomizerGetBytesFromStringMethod{},
	}
}

func secureRandomIndex(limit int) (int, error) {
	if limit <= 0 {
		return 0, fmt.Errorf("随机范围必须大于 0")
	}
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(limit)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func randomizerIntArg(ctx data.Context, index int) (int, bool) {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return 0, false
	}
	asInt, ok := v.(data.AsInt)
	if !ok {
		return 0, false
	}
	n, err := asInt.AsInt()
	return n, err == nil
}

type randomizerArrayEntry struct {
	key   data.Value
	value data.Value
}

func randomizerArrayEntries(value data.Value) ([]randomizerArrayEntry, bool) {
	var entries []randomizerArrayEntry
	switch array := value.(type) {
	case *data.ArrayValue:
		entries = make([]randomizerArrayEntry, 0, len(array.List))
		for index, zv := range array.List {
			if zv == nil {
				continue
			}
			var key data.Value = data.NewIntValue(index)
			if zv.Name != "" {
				if numericKey, ok := data.ParseIntArrayKeyName(zv.Name); ok {
					key = data.NewIntValue(numericKey)
				} else {
					key = data.NewStringValue(zv.Name)
				}
			}
			entries = append(entries, randomizerArrayEntry{key: key, value: zv.Value})
		}
	case *data.ObjectValue:
		array.RangeProperties(func(name string, value data.Value) bool {
			var key data.Value = data.NewStringValue(name)
			if numericKey, ok := data.ParseIntArrayKeyName(name); ok {
				key = data.NewIntValue(numericKey)
			}
			entries = append(entries, randomizerArrayEntry{key: key, value: value})
			return true
		})
	default:
		return nil, false
	}
	return entries, true
}

type randomizerPickArrayKeysMethod struct{}

func (m *randomizerPickArrayKeysMethod) GetName() string            { return "pickArrayKeys" }
func (m *randomizerPickArrayKeysMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *randomizerPickArrayKeysMethod) GetIsStatic() bool          { return false }
func (m *randomizerPickArrayKeysMethod) GetReturnType() data.Types  { return nil }
func (m *randomizerPickArrayKeysMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
		node.NewParameter(nil, "num", 1, nil, data.Int{}),
	}
}
func (m *randomizerPickArrayKeysMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Arrays{}),
		node.NewVariable(nil, "num", 1, data.Int{}),
	}
}
func (m *randomizerPickArrayKeysMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::pickArrayKeys 缺少 array 参数"))
	}
	entries, ok := randomizerArrayEntries(value)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::pickArrayKeys 需要 array 参数"))
	}
	count, ok := randomizerIntArg(ctx, 1)
	if !ok || count < 1 || count > len(entries) {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::pickArrayKeys 的 num 超出范围"))
	}

	indices := make([]int, len(entries))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < count; i++ {
		offset, err := secureRandomIndex(len(indices) - i)
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		j := i + offset
		indices[i], indices[j] = indices[j], indices[i]
	}

	keys := make([]data.Value, 0, count)
	for _, index := range indices[:count] {
		keys = append(keys, entries[index].key)
	}
	return data.NewArrayValue(keys), nil
}

type randomizerShuffleArrayMethod struct{}

func (m *randomizerShuffleArrayMethod) GetName() string            { return "shuffleArray" }
func (m *randomizerShuffleArrayMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *randomizerShuffleArrayMethod) GetIsStatic() bool          { return false }
func (m *randomizerShuffleArrayMethod) GetReturnType() data.Types  { return nil }
func (m *randomizerShuffleArrayMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, nil, data.Arrays{})}
}
func (m *randomizerShuffleArrayMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.Arrays{})}
}
func (m *randomizerShuffleArrayMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::shuffleArray 缺少 array 参数"))
	}
	entries, ok := randomizerArrayEntries(value)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::shuffleArray 需要 array 参数"))
	}
	values := make([]data.Value, 0, len(entries))
	for _, entry := range entries {
		values = append(values, entry.value)
	}
	for i := len(values) - 1; i > 0; i-- {
		j, err := secureRandomIndex(i + 1)
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		values[i], values[j] = values[j], values[i]
	}
	return data.NewArrayValue(values), nil
}

type randomizerGetBytesFromStringMethod struct{}

func (m *randomizerGetBytesFromStringMethod) GetName() string            { return "getBytesFromString" }
func (m *randomizerGetBytesFromStringMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *randomizerGetBytesFromStringMethod) GetIsStatic() bool          { return false }
func (m *randomizerGetBytesFromStringMethod) GetReturnType() data.Types  { return data.String{} }
func (m *randomizerGetBytesFromStringMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "source", 0, nil, data.String{}),
		node.NewParameter(nil, "length", 1, nil, data.Int{}),
	}
}
func (m *randomizerGetBytesFromStringMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "source", 0, data.String{}),
		node.NewVariable(nil, "length", 1, data.Int{}),
	}
}
func (m *randomizerGetBytesFromStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::getBytesFromString 缺少 source 参数"))
	}
	source := []byte(value.AsString())
	length, ok := randomizerIntArg(ctx, 1)
	if !ok || length < 1 || len(source) == 0 {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Randomizer::getBytesFromString 参数无效"))
	}
	result := make([]byte, length)
	for i := range result {
		index, err := secureRandomIndex(len(source))
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		result[i] = source[index]
	}
	return data.NewStringValue(string(result)), nil
}
