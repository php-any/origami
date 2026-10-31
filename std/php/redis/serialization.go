package redis

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php"
)

type phpControlError struct{ control data.Control }

func (e *phpControlError) Error() string { return "PHP serialization hook failed" }

// Redis serializes values, never key names, hash fields, Lua arguments, numeric
// increments or rawCommand arguments. Positions below follow prepare's mapping.
func serializedArgument(name string, i, zStart int) bool {
	switch name {
	case "set", "setnx", "getset":
		return i == 1
	case "setex", "psetex", "lset", "lrem":
		return i == 2
	case "mset", "msetnx":
		return i%2 == 1
	case "hset", "hmset", "hsetnx":
		return i >= 2 && i%2 == 0
	case "lpush", "rpush", "lpushx", "rpushx", "sadd", "srem", "sismember", "zrem", "zscore", "zrank", "zrevrank":
		return i > 0
	case "zincrby":
		return i == 2
	case "zadd":
		return i >= zStart && (i-zStart)%2 == 0
	}
	return false
}
func (s *clientState) prepare(ctx data.Context, name string, a []data.Value) (prepared, error) {
	if s.serializer != 1 {
		return prepare(name, a, s.prefix)
	}
	p, err := prepare(name, a, s.prefix, func(v data.Value) (string, error) { return s.pack(ctx, v) })
	p.serialized = true
	return p, err
}
func (s *clientState) pack(ctx data.Context, v data.Value) (string, error) {
	fn := php.NewSerializeFunction()
	frame := ctx.CreateContext(fn.GetVariables())
	frame.SetIndexZVal(0, data.NewZVal(v))
	// __serialize may re-enter this same Redis object. Never hold the connection
	// lock while executing user PHP code. No command has been sent at this point.
	s.mu.Unlock()
	defer s.mu.Lock()
	result, ctl := fn.Call(frame)
	if ctl != nil {
		return "", &phpControlError{ctl}
	}
	str, ok := result.(*data.StringValue)
	if !ok {
		return "", fmt.Errorf("PHP serialization failed")
	}
	return str.Value, nil
}
func (s *clientState) unpack(ctx data.Context, value data.Value) (data.Value, error) {
	str, ok := value.(*data.StringValue)
	if !ok {
		return value, nil
	}
	s.mu.Unlock()
	defer s.mu.Lock()
	v, valid, ctl := php.DecodeSerializedValue(ctx, str.Value)
	if ctl != nil {
		return nil, &phpControlError{ctl}
	}
	if !valid {
		return value, nil
	}
	return v, nil
}
func (s *clientState) convertValue(ctx data.Context, p prepared, reply any) (data.Value, error) {
	v, err := s.convert(p, reply)
	if err != nil || !p.serialized || len(p.args) == 0 || p.result == "raw" {
		return v, err
	}
	switch p.args[0] {
	case "GET", "GETDEL", "GETSET", "HGET", "LINDEX":
		return s.unpack(ctx, v)
	case "SET":
		if p.result != "bool" {
			return s.unpack(ctx, v)
		}
	case "MGET", "HMGET", "HGETALL", "HVALS", "LRANGE", "SMEMBERS", "SUNION", "SINTER", "SDIFF", "LPOP", "RPOP", "SPOP", "SRANDMEMBER", "ZRANGE", "ZREVRANGE":
		if p.result == "scores" {
			return v, nil
		} // Members remain string keys in score maps.
		if array, ok := v.(*data.ArrayValue); ok {
			out := data.NewArrayValue(nil).(*data.ArrayValue)
			for i, slot := range array.Range() {
				decoded, err := s.unpack(ctx, slot.ReadValue())
				if err != nil {
					return nil, err
				}
				out.SetKey(slot.PHPArrayKey(i), decoded)
			}
			return out, nil
		}
		return s.unpack(ctx, v)
	}
	return v, nil
}
