// Package redis provides an opt-in subset of the phpredis API. Hosts must call
// Load explicitly; the PHP standard library does not import this module.
package redis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/exception"
)

const moduleVersion = "0.1.0-origami"
const (
	optSerializer       = 1
	optPrefix           = 2
	optReadTimeout      = 3
	optCompression      = 7
	optScan             = 4
	optMaxRetries       = 11
	optBackoffAlgorithm = 12
	optBackoffBase      = 13
	optBackoffCap       = 14
)
const (
	atomicMode   = 0
	multiMode    = 1
	pipelineMode = 2
)

type redisClass struct {
	node.Node
	methods map[string]*redisMethod
}

func (c *redisClass) GetName() string                          { return "Redis" }
func (c *redisClass) GetExtend() *string                       { return nil }
func (c *redisClass) GetImplements() []string                  { return nil }
func (c *redisClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *redisClass) GetPropertyList() []data.Property         { return nil }
func (c *redisClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	v := data.NewClassValue(c, ctx)
	v.InstanceSource = newState()
	return v, nil
}
func (c *redisClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *redisClass) GetMethod(n string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(n)]
	if !ok {
		return nil, false
	}
	return m, true
}
func (c *redisClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *redisClass) GetStaticProperty(n string) (data.Value, bool) {
	v, ok := constants[n]
	if !ok {
		return nil, false
	}
	return data.NewIntValue(v), true
}

var constants = map[string]int{
	"ATOMIC": 0, "MULTI": 1, "PIPELINE": 2, "OPT_SERIALIZER": 1, "OPT_PREFIX": 2, "OPT_READ_TIMEOUT": 3, "OPT_COMPRESSION": 7,
	"SERIALIZER_NONE": 0, "COMPRESSION_NONE": 0,
	"SERIALIZER_PHP": 1, "OPT_SCAN": optScan,
	"SCAN_NORETRY": 0, "SCAN_RETRY": 1, "SCAN_PREFIX": 2, "SCAN_NOPREFIX": 3,
	"OPT_MAX_RETRIES": optMaxRetries, "OPT_BACKOFF_ALGORITHM": optBackoffAlgorithm,
	"OPT_BACKOFF_BASE": optBackoffBase, "OPT_BACKOFF_CAP": optBackoffCap,
	"BACKOFF_ALGORITHM_DEFAULT": 0, "BACKOFF_ALGORITHM_DECORRELATED_JITTER": 1,
	"BACKOFF_ALGORITHM_FULL_JITTER": 2, "BACKOFF_ALGORITHM_EQUAL_JITTER": 3,
	"BACKOFF_ALGORITHM_EXPONENTIAL": 4, "BACKOFF_ALGORITHM_UNIFORM": 5, "BACKOFF_ALGORITHM_CONSTANT": 6,
	"REDIS_NOT_FOUND": 0, "REDIS_STRING": 1, "REDIS_SET": 2, "REDIS_LIST": 3, "REDIS_ZSET": 4, "REDIS_HASH": 5, "REDIS_STREAM": 6,
}

type redisMethod struct {
	name string
	call func(data.Context, *clientState, []data.Value) (data.GetValue, data.Control)
}

func (m *redisMethod) GetName() string            { return m.name }
func (m *redisMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *redisMethod) GetIsStatic() bool          { return false }
func (m *redisMethod) GetReturnType() data.Types  { return nil }
func (m *redisMethod) GetParams() []data.GetValue {
	if m.name == "scan" {
		return scanParams
	}
	return []data.GetValue{node.NewParameters(nil, "args", 0, nil, nil)}
}
func (m *redisMethod) GetVariables() []data.Variable {
	if m.name == "scan" {
		return scanVariables
	}
	return []data.Variable{node.NewVariable(nil, "args", 0, nil)}
}
func (m *redisMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	o := instance(ctx)
	if o == nil {
		return failure(fmt.Errorf("Redis method requires an instance"))
	}
	s, ok := o.InstanceSource.(*clientState)
	if !ok {
		return failure(fmt.Errorf("uninitialized Redis instance"))
	}
	var args []data.Value
	if m.name == "scan" {
		for i := range scanParams {
			v, _ := ctx.GetIndexValue(i)
			args = append(args, v)
		}
	} else if v, ok := ctx.GetIndexValue(0); ok {
		if a, ok := v.(*data.ArrayValue); ok {
			args = a.ToValueList()
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return m.call(ctx, s, args)
}
func instance(ctx data.Context) *data.ClassValue {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ClassValue.InstanceIdentity()
	case *data.ClassValue:
		return c.InstanceIdentity()
	}
	return nil
}
func failure(err error) (data.GetValue, data.Control) {
	var e *phpControlError
	if errors.As(err, &e) {
		return nil, e.control
	}
	return nil, data.NewErrorThrowByName(nil, err, "RedisException")
}
func bad(format string, args ...any) (data.GetValue, data.Control) {
	return failure(fmt.Errorf(format, args...))
}
func null(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}
func scalar(v data.Value) (string, error) {
	switch v.(type) {
	case *data.StringValue, *data.IntValue, *data.FloatValue, *data.BoolValue, *data.NullValue:
		return v.AsString(), nil
	}
	return "", fmt.Errorf("Redis values must be scalar when serialization is disabled")
}
func integer(v data.Value) (int, error) {
	n, ok := v.(*data.IntValue)
	if !ok {
		return 0, fmt.Errorf("expected integer")
	}
	return n.Value, nil
}
func arity(name string, a []data.Value, min, max int) error {
	if len(a) < min || (max >= 0 && len(a) > max) {
		return fmt.Errorf("Redis::%s received %d arguments; expected %d..%d", name, len(a), min, max)
	}
	return nil
}

type redisException struct{ *exception.ExceptionClass }

func (c *redisException) GetName() string    { return "RedisException" }
func (c *redisException) GetExtend() *string { n := "RuntimeException"; return &n }
func (c *redisException) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

// Load only registers declarations and constants. It never opens a connection.
func Load(vm data.VM) {
	c := &redisClass{methods: map[string]*redisMethod{}}
	for n, f := range localMethods() {
		c.methods[n] = &redisMethod{name: n, call: f}
	}
	for name := range commands {
		n := name
		c.methods[n] = &redisMethod{name: n, call: func(ctx data.Context, s *clientState, args []data.Value) (data.GetValue, data.Control) {
			p, err := s.prepare(ctx, n, args)
			if err != nil {
				return failure(err)
			}
			return s.execute(ctx, p)
		}}
	}
	vm.AddClass(&redisException{exception.NewExceptionClass()})
	vm.AddClass(c)
	vm.SetConstant("REDIS_VERSION", data.NewStringValue(moduleVersion))
}
