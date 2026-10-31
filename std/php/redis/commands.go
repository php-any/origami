package redis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

type commandSpec struct {
	min, max     int
	keys, result string
}
type prepared struct {
	args       []string
	result     string
	fields     []string
	serialized bool
}

var commands = map[string]commandSpec{
	"ping": {0, 1, "", ""}, "echo": {1, 1, "", ""}, "auth": {1, 1, "", "bool"}, "select": {1, 1, "", "bool"},
	"get": {1, 1, "first", ""}, "getdel": {1, 1, "first", ""}, "getset": {2, 2, "first", ""}, "set": {2, 3, "first", "bool"}, "setex": {3, 3, "first", "bool"}, "psetex": {3, 3, "first", "bool"}, "setnx": {2, 2, "first", "bool"},
	"mget": {1, 1, "all", ""}, "mset": {1, 1, "pairs", "bool"}, "msetnx": {1, 1, "pairs", "bool"},
	"del": {1, -1, "all", ""}, "unlink": {1, -1, "all", ""}, "exists": {1, -1, "all", ""}, "touch": {1, -1, "all", ""},
	"expire": {2, 3, "first", "bool"}, "pexpire": {2, 3, "first", "bool"}, "expireat": {2, 3, "first", "bool"}, "pexpireat": {2, 3, "first", "bool"}, "ttl": {1, 1, "first", ""}, "pttl": {1, 1, "first", ""}, "persist": {1, 1, "first", "bool"}, "type": {1, 1, "first", "type"},
	"incr": {1, 2, "first", ""}, "decr": {1, 2, "first", ""}, "incrby": {2, 2, "first", ""}, "decrby": {2, 2, "first", ""}, "incrbyfloat": {2, 2, "first", "float"}, "append": {2, 2, "first", ""}, "strlen": {1, 1, "first", ""}, "getrange": {3, 3, "first", ""}, "setrange": {3, 3, "first", ""},
	"rename": {2, 2, "all", "bool"}, "renamenx": {2, 2, "all", "bool"}, "keys": {1, 1, "first", ""}, "dbsize": {0, 0, "", ""}, "time": {0, 0, "", ""},
	"hget": {2, 2, "first", ""}, "hset": {2, -1, "first", ""}, "hsetnx": {3, 3, "first", "bool"}, "hmget": {2, 2, "first", "hmget"}, "hmset": {2, 2, "first", "bool"}, "hgetall": {1, 1, "first", "pairs"}, "hdel": {2, -1, "first", ""}, "hexists": {2, 2, "first", "bool"}, "hlen": {1, 1, "first", ""}, "hkeys": {1, 1, "first", ""}, "hvals": {1, 1, "first", ""}, "hincrby": {3, 3, "first", ""}, "hincrbyfloat": {3, 3, "first", "float"},
	"lpush": {2, -1, "first", ""}, "rpush": {2, -1, "first", ""}, "lpushx": {2, -1, "first", ""}, "rpushx": {2, -1, "first", ""}, "lpop": {1, 2, "first", ""}, "rpop": {1, 2, "first", ""}, "llen": {1, 1, "first", ""}, "lrange": {3, 3, "first", ""}, "lindex": {2, 2, "first", ""}, "lset": {3, 3, "first", "bool"}, "ltrim": {3, 3, "first", "bool"}, "lrem": {2, 3, "first", ""},
	"sadd": {2, -1, "first", ""}, "srem": {2, -1, "first", ""}, "sismember": {2, 2, "first", "bool"}, "smembers": {1, 1, "first", ""}, "scard": {1, 1, "first", ""}, "spop": {1, 2, "first", ""}, "srandmember": {1, 2, "first", ""}, "sunion": {1, -1, "all", ""}, "sinter": {1, -1, "all", ""}, "sdiff": {1, -1, "all", ""},
	"zadd": {3, -1, "first", ""}, "zrem": {2, -1, "first", ""}, "zcard": {1, 1, "first", ""}, "zscore": {2, 2, "first", "float"}, "zrank": {2, 2, "first", ""}, "zrevrank": {2, 2, "first", ""}, "zincrby": {3, 3, "first", "float"}, "zrange": {3, 4, "first", ""}, "zrevrange": {3, 4, "first", ""}, "zcount": {3, 3, "first", ""}, "zremrangebyscore": {3, 3, "first", ""}, "zremrangebyrank": {3, 3, "first", ""},
	"eval": {1, 3, "eval", ""}, "evalsha": {1, 3, "eval", ""}, "script": {1, -1, "", ""}, "rawcommand": {1, -1, "raw", "raw"},
	"watch": {1, -1, "all", "bool"}, "unwatch": {0, 0, "", "bool"}, "flushdb": {0, 1, "", "bool"}, "flushall": {0, 1, "", "bool"}, "info": {0, 1, "", "info"}, "client": {1, -1, "", "raw"},
}

func values(v data.Value) ([]data.Value, error) {
	a, ok := v.(*data.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("expected array")
	}
	return a.ToValueList(), nil
}
func pairs(v data.Value) ([]data.Value, error) {
	a, ok := v.(*data.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("expected associative array")
	}
	out := make([]data.Value, 0, a.Len()*2)
	for i, s := range a.Range() {
		out = append(out, s.PHPArrayKey(i), s.ReadValue())
	}
	return out, nil
}
func prepare(name string, input []data.Value, prefix string, encoders ...func(data.Value) (string, error)) (prepared, error) {
	spec := commands[name]
	p := prepared{result: spec.result}
	if e := arity(name, input, spec.min, spec.max); e != nil {
		return p, e
	}
	a := append([]data.Value(nil), input...)
	cmd := strings.ToUpper(name)
	zMemberStart := 2
	switch name {
	case "auth":
		if _, ok := a[0].(*data.ArrayValue); ok {
			var e error
			a, e = values(a[0])
			if e != nil {
				return p, e
			}
			if len(a) != 1 && len(a) != 2 {
				return p, fmt.Errorf("auth expects a password or [username, password]")
			}
		}
	case "mget":
		var e error
		a, e = values(a[0])
		if e != nil {
			return p, e
		}
	case "mset", "msetnx":
		var e error
		a, e = pairs(a[0])
		if e != nil {
			return p, e
		}
	case "hmset":
		v, e := pairs(a[1])
		if e != nil {
			return p, e
		}
		a = append(a[:1], v...)
	case "hset":
		if len(a) == 2 {
			v, e := pairs(a[1])
			if e != nil {
				return p, e
			}
			a = append(a[:1], v...)
		}
		if len(a) < 3 || len(a)%2 != 1 {
			return p, fmt.Errorf("hset requires field/value pairs")
		}
	case "hmget":
		v, e := values(a[1])
		if e != nil {
			return p, e
		}
		if len(v) == 0 {
			return p, fmt.Errorf("hmget requires at least one field")
		}
		a = append(a[:1], v...)
		for _, v := range v {
			s, e := scalar(v)
			if e != nil {
				return p, e
			}
			p.fields = append(p.fields, s)
		}
	case "del", "unlink", "exists", "touch", "watch", "sunion", "sinter", "sdiff":
		if len(a) == 1 {
			if _, ok := a[0].(*data.ArrayValue); ok {
				a, _ = values(a[0])
			}
		}
	case "incr", "decr":
		if len(a) == 2 {
			cmd += "BY"
		}
	case "lrem":
		count := data.NewIntValue(0)
		if len(a) == 3 {
			count = a[2]
		}
		a = []data.Value{a[0], count, a[1]}
	case "set":
		if len(a) == 3 {
			opt := a[2]
			a = a[:2]
			if !null(opt) {
				if n, ok := opt.(*data.IntValue); ok {
					if n.Value <= 0 {
						return p, fmt.Errorf("SET expiry must be positive")
					}
					a = append(a, data.NewStringValue("EX"), n)
				} else {
					o, ok := opt.(*data.ArrayValue)
					if !ok {
						return p, fmt.Errorf("SET options must be array, integer or null")
					}
					seen := map[string]bool{}
					for i, s := range o.Range() {
						k := s.PHPArrayKey(i)
						v := s.ReadValue()
						var flag string
						if _, numeric := k.(*data.IntValue); numeric {
							if null(v) {
								continue
							}
							str, e := scalar(v)
							if e != nil {
								return p, e
							}
							flag = strings.ToUpper(str)
							switch flag {
							case "NX", "XX", "KEEPTTL", "GET":
							default:
								return p, fmt.Errorf("unsupported SET flag %s", flag)
							}
							a = append(a, data.NewStringValue(flag))
						} else {
							flag = strings.ToUpper(k.AsString())
							switch flag {
							case "EX", "PX", "EXAT", "PXAT":
								n, e := integer(v)
								if e != nil || n <= 0 {
									return p, fmt.Errorf("%s requires a positive integer", flag)
								}
								a = append(a, data.NewStringValue(flag), v)
							default:
								return p, fmt.Errorf("unsupported SET option %s", flag)
							}
						}
						if seen[flag] {
							return p, fmt.Errorf("duplicate SET option %s", flag)
						}
						seen[flag] = true
					}
					if seen["NX"] && seen["XX"] {
						return p, fmt.Errorf("SET NX and XX are mutually exclusive")
					}
					if seen["GET"] {
						p.result = ""
					}
				}
			}
		}
	case "zrange", "zrevrange":
		if len(a) == 4 {
			b, ok := a[3].(*data.BoolValue)
			if !ok {
				return p, fmt.Errorf("range options currently require a withScores boolean")
			}
			a = a[:3]
			if b.Value {
				a = append(a, data.NewStringValue("WITHSCORES"))
				p.result = "scores"
			}
		}
	case "zadd":
		start := 1
		if opt, ok := a[1].(*data.ArrayValue); ok {
			flags := []data.Value{}
			for _, s := range opt.Range() {
				f, e := scalar(s.ReadValue())
				if e != nil {
					return p, e
				}
				f = strings.ToUpper(f)
				switch f {
				case "NX", "XX", "GT", "LT", "CH", "INCR":
				default:
					return p, fmt.Errorf("unsupported ZADD flag %s", f)
				}
				if f == "INCR" {
					p.result = "float"
				}
				flags = append(flags, data.NewStringValue(f))
			}
			rest := append([]data.Value(nil), a[2:]...)
			a = append(a[:1], flags...)
			start = len(a)
			a = append(a, rest...)
		}
		if len(a)-start < 2 || (len(a)-start)%2 != 0 {
			return p, fmt.Errorf("ZADD requires score/member pairs")
		}
		zMemberStart = start + 1
	case "eval", "evalsha":
		var rest []data.Value
		var e error
		n := 0
		if len(a) > 1 {
			rest, e = values(a[1])
			if e != nil {
				return p, e
			}
		}
		if len(a) > 2 {
			n, e = integer(a[2])
			if e != nil {
				return p, e
			}
		}
		if n < 0 || n > len(rest) {
			return p, fmt.Errorf("numKeys must be between 0 and the argument count")
		}
		for i := 0; i < n; i++ {
			s, e := scalar(rest[i])
			if e != nil {
				return p, e
			}
			rest[i] = data.NewStringValue(prefix + s)
		}
		a = append([]data.Value{a[0], data.NewIntValue(n)}, rest...)
	case "ping":
		if len(a) == 1 && null(a[0]) {
			a = nil
		}
	case "flushdb", "flushall":
		if len(a) == 1 {
			b, ok := a[0].(*data.BoolValue)
			if !ok {
				return p, fmt.Errorf("async must be boolean")
			}
			a = nil
			if b.Value {
				a = []data.Value{data.NewStringValue("ASYNC")}
			}
		}
	}
	if len(a) == 0 && spec.min > 0 {
		return p, fmt.Errorf("%s requires nonempty arguments", name)
	}
	p.args = append(p.args, cmd)
	for i, v := range a {
		encode := scalar
		if len(encoders) > 0 && serializedArgument(name, i, zMemberStart) {
			encode = encoders[0]
		}
		s, e := encode(v)
		if e != nil {
			return p, e
		}
		if spec.keys == "first" && i == 0 || spec.keys == "all" || spec.keys == "pairs" && i%2 == 0 {
			s = prefix + s
		}
		p.args = append(p.args, s)
	}
	if name == "rawcommand" {
		p.args = p.args[1:]
		p.args[0] = strings.ToUpper(p.args[0])
		switch p.args[0] {
		case "HELLO", "SUBSCRIBE", "PSUBSCRIBE", "SSUBSCRIBE", "UNSUBSCRIBE", "PUNSUBSCRIBE", "SUNSUBSCRIBE", "MONITOR", "SYNC", "PSYNC", "MULTI", "EXEC", "DISCARD", "WATCH", "UNWATCH", "QUIT", "RESET", "SELECT", "AUTH", "CLIENT":
			return p, fmt.Errorf("rawCommand %s changes connection state; use a supported dedicated method", p.args[0])
		}
	}
	if name == "client" {
		sub := strings.ToUpper(p.args[1])
		switch sub {
		case "GETNAME", "SETNAME", "ID", "INFO", "LIST":
		default:
			return p, fmt.Errorf("unsupported CLIENT subcommand %s", sub)
		}
	}
	return p, nil
}

// Ordinary Redis command failures return false and remain available through
// getLastError. Connection/auth/loading errors throw, matching phpredis 6.1.
func shouldThrow(s string) bool {
	for _, prefix := range []string{"NOSCRIPT", "NOQUORUM", "NOGOODSLAVE", "WRONGTYPE", "BUSYGROUP", "NOGROUP"} {
		if strings.HasPrefix(s, prefix) {
			return false
		}
	}
	return !strings.HasPrefix(s, "ERR") || strings.HasPrefix(s, "ERR AUTH")
}
func (s *clientState) decode(v any, raw bool) (data.Value, error) {
	switch v := v.(type) {
	case nil:
		return data.NewBoolValue(false), nil
	case serverError:
		s.lastError = string(v)
		if shouldThrow(string(v)) {
			return nil, fmt.Errorf("%s", v)
		}
		return data.NewBoolValue(false), nil
	case status:
		if raw {
			return data.NewStringValue(string(v)), nil
		}
		return data.NewBoolValue(true), nil
	case string:
		return data.NewStringValue(v), nil
	case int64:
		return data.NewIntValue(int(v)), nil
	case []any:
		a := make([]data.Value, 0, len(v))
		for _, x := range v {
			d, e := s.decode(x, raw)
			if e != nil {
				return nil, e
			}
			a = append(a, d)
		}
		return data.NewArrayValue(a), nil
	default:
		return nil, fmt.Errorf("invalid Redis response %T", v)
	}
}
func (s *clientState) convert(p prepared, v any) (data.Value, error) {
	if _, ok := v.(serverError); ok {
		return s.decode(v, false)
	}
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	switch p.result {
	case "bool":
		switch x := v.(type) {
		case status:
			return data.NewBoolValue(x == "OK"), nil
		case int64:
			return data.NewBoolValue(x != 0), nil
		}
		return nil, fmt.Errorf("expected Redis boolean result")
	case "float":
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expected Redis numeric string")
		}
		n, e := strconv.ParseFloat(str, 64)
		if e != nil {
			return nil, e
		}
		return data.NewFloatValue(n), nil
	case "type":
		var str string
		switch x := v.(type) {
		case status:
			str = string(x)
		case string:
			str = x
		default:
			return nil, fmt.Errorf("invalid TYPE response")
		}
		n, ok := map[string]int{"none": 0, "string": 1, "set": 2, "list": 3, "zset": 4, "hash": 5, "stream": 6}[str]
		if !ok {
			return nil, fmt.Errorf("unknown Redis type %q", str)
		}
		return data.NewIntValue(n), nil
	case "pairs", "scores", "hmget":
		a, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("expected Redis array result")
		}
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		if p.result == "hmget" {
			if len(a) != len(p.fields) {
				return nil, fmt.Errorf("HMGET reply length mismatch")
			}
			for i, x := range a {
				d, e := s.decode(x, false)
				if e != nil {
					return nil, e
				}
				out.SetKey(data.NewStringValue(p.fields[i]), d)
			}
			return out, nil
		}
		if len(a)%2 != 0 {
			return nil, fmt.Errorf("odd key/value reply length")
		}
		for i := 0; i < len(a); i += 2 {
			k, ok := a[i].(string)
			if !ok {
				return nil, fmt.Errorf("invalid reply key")
			}
			var d data.Value
			var e error
			if p.result == "scores" {
				d, e = s.convert(prepared{result: "float"}, a[i+1])
			} else {
				d, e = s.decode(a[i+1], false)
			}
			if e != nil {
				return nil, e
			}
			out.SetKey(data.NewStringValue(k), d)
		}
		return out, nil
	case "info":
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid INFO reply")
		}
		a := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, line := range strings.Split(str, "\r\n") {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			k, x, ok := strings.Cut(line, ":")
			if ok {
				var d data.Value = data.NewStringValue(x)
				if n, e := strconv.ParseInt(x, 10, 64); e == nil {
					d = data.NewIntValue(int(n))
				}
				a.SetStringKey(k, d)
			}
		}
		return a, nil
	}
	return s.decode(v, p.result == "raw")
}
