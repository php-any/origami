package hashing

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const bcryptHasherClassName = "Illuminate\\Hashing\\BcryptHasher"

type BcryptHasherClass struct {
	node.Node
	methods map[string]data.Method
}

func NewBcryptHasherClass() data.ClassStmt {
	c := &BcryptHasherClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *BcryptHasherClass) GetName() string                          { return bcryptHasherClassName }
func (c *BcryptHasherClass) GetExtend() *string                       { s := "Illuminate\\Hashing\\AbstractHasher"; return &s }
func (c *BcryptHasherClass) GetImplements() []string                  { return nil }
func (c *BcryptHasherClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "rounds", "verifyAlgorithm", "limit":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *BcryptHasherClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 3)
	for _, n := range []string{"rounds", "verifyAlgorithm", "limit"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *BcryptHasherClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *BcryptHasherClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	_ = cv.SetProperty("rounds", data.NewIntValue(bcrypt.DefaultCost))
	_ = cv.SetProperty("verifyAlgorithm", data.NewBoolValue(true))
	return cv, nil
}
func (c *BcryptHasherClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *BcryptHasherClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *BcryptHasherClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *BcryptHasherClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"options"}, 0, bcryptConstruct)
	c.methods["make"] = kit.InstanceMethodOpt("make", []string{"value", "options"}, 1, bcryptMake)
	c.methods["check"] = kit.InstanceMethodOpt("check", []string{"value", "hashedValue", "options"}, 2, bcryptCheck)
	c.methods["needsrehash"] = kit.InstanceMethodOpt("needsRehash", []string{"hashedValue", "options"}, 1, bcryptNeedsRehash)
	c.methods["info"] = kit.InstanceMethod("info", []string{"hashedValue"}, bcryptInfo)
}

func bcryptRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("BcryptHasher missing $this"))
}

func bcryptConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bcryptRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("rounds", data.NewIntValue(bcrypt.DefaultCost))
	_ = cv.SetProperty("verifyAlgorithm", data.NewBoolValue(true))
	opts := kit.Arg(ctx, 0)
	if av, ok := kit.Unwrap(opts).(*data.ArrayValue); ok && av != nil {
		if r, ctl := av.GetProperty("rounds"); ctl == nil && r != nil {
			if iv, ok := r.(data.AsInt); ok {
				if n, err := iv.AsInt(); err == nil {
					_ = cv.SetProperty("rounds", data.NewIntValue(n))
				}
			}
		}
		if v, ctl := av.GetProperty("verify"); ctl == nil && v != nil {
			if bv, ok := v.(data.AsBool); ok {
				b, _ := bv.AsBool()
				_ = cv.SetProperty("verifyAlgorithm", data.NewBoolValue(b))
			}
		}
	}
	return cv, nil
}

func bcryptCost(cv *data.ClassValue, opts data.Value) int {
	cost := bcrypt.DefaultCost
	if rv, _ := cv.GetProperty("rounds"); rv != nil {
		if iv, ok := kit.Unwrap(rv).(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				cost = n
			}
		}
	}
	if av, ok := kit.Unwrap(opts).(*data.ArrayValue); ok && av != nil {
		if r, ctl := av.GetProperty("rounds"); ctl == nil && r != nil {
			if iv, ok := r.(data.AsInt); ok {
				if n, err := iv.AsInt(); err == nil {
					cost = n
				}
			}
		}
	}
	if cost < bcrypt.MinCost {
		cost = bcrypt.MinCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return cost
}

func bcryptMake(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bcryptRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	plain := ""
	if v := kit.Arg(ctx, 0); v != nil {
		plain = v.AsString()
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost(cv, kit.Arg(ctx, 1)))
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Bcrypt hashing not supported."))
	}
	return data.NewStringValue(string(hash)), nil
}

func bcryptCheck(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bcryptRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	plain := ""
	hash := ""
	if v := kit.Arg(ctx, 0); v != nil {
		plain = v.AsString()
	}
	if v := kit.Arg(ctx, 1); v != nil {
		hash = v.AsString()
	}
	if hash == "" || kit.IsNull(kit.Arg(ctx, 1)) {
		return data.NewBoolValue(false), nil
	}
	verifyV, _ := cv.GetProperty("verifyAlgorithm")
	if kit.Truthy(verifyV) && !strings.HasPrefix(hash, "$2y$") && !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("This password does not use the Bcrypt algorithm."))
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return data.NewBoolValue(err == nil), nil
}

func bcryptNeedsRehash(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bcryptRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	hash := ""
	if v := kit.Arg(ctx, 0); v != nil {
		hash = v.AsString()
	}
	cost := bcryptCost(cv, kit.Arg(ctx, 1))
	cur := bcryptCostFromHash(hash)
	return data.NewBoolValue(cur != cost), nil
}

func bcryptCostFromHash(hash string) int {
	parts := strings.Split(hash, "$")
	if len(parts) < 3 {
		return bcrypt.DefaultCost
	}
	var cost int
	fmt.Sscanf(parts[2], "%d", &cost)
	if cost < bcrypt.MinCost {
		return bcrypt.DefaultCost
	}
	return cost
}

func bcryptInfo(ctx data.Context) (data.GetValue, data.Control) {
	hash := ""
	if v := kit.Arg(ctx, 0); v != nil {
		hash = v.AsString()
	}
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	arr.SetStringKey("algoName", data.NewStringValue("bcrypt"))
	opts := data.NewArrayValue(nil).(*data.ArrayValue)
	opts.SetStringKey("cost", data.NewIntValue(bcryptCostFromHash(hash)))
	arr.SetStringKey("options", opts)
	return arr, nil
}
