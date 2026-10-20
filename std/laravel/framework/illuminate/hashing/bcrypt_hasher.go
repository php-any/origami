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

func (c *BcryptHasherClass) GetName() string { return bcryptHasherClassName }
func (c *BcryptHasherClass) GetExtend() *string {
	s := "Illuminate\\Hashing\\AbstractHasher"
	return &s
}
func (c *BcryptHasherClass) GetImplements() []string { return nil }
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
	_ = cv.SetProperty("rounds", data.NewIntValue(12))
	_ = cv.SetProperty("verifyAlgorithm", data.NewBoolValue(true))
	return cv, nil
}
func (c *BcryptHasherClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *BcryptHasherClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *BcryptHasherClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

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
	_ = cv.SetProperty("rounds", data.NewIntValue(12))
	_ = cv.SetProperty("verifyAlgorithm", data.NewBoolValue(true))
	if av, ok := kit.Unwrap(kit.Arg(ctx, 0)).(*data.ArrayValue); ok && av != nil {
		if n, ok := arrayInt(av, "rounds"); ok {
			_ = cv.SetProperty("rounds", data.NewIntValue(n))
		}
		if b, ok := arrayBool(av, "verify"); ok {
			_ = cv.SetProperty("verifyAlgorithm", data.NewBoolValue(b))
		}
	}
	return cv, nil
}

func bcryptCost(cv *data.ClassValue, opts data.Value) int {
	cost := 12
	if rv, _ := cv.GetProperty("rounds"); rv != nil {
		if iv, ok := kit.Unwrap(rv).(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				cost = n
			}
		}
	}
	if av, ok := kit.Unwrap(opts).(*data.ArrayValue); ok && av != nil {
		if n, ok := arrayInt(av, "rounds"); ok {
			cost = n
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
	// PHP password_hash 固定写出 $2y$。$2a$ 会被 password_get_info 标成 unknown，Laravel 会拒绝登录。
	return data.NewStringValue(strings.Replace(string(hash), "$2a$", "$2y$", 1)), nil
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
	if !strings.HasPrefix(hash, "$2y$") {
		return data.NewBoolValue(true), nil
	}
	cost := bcryptCost(cv, kit.Arg(ctx, 1))
	return data.NewBoolValue(bcryptCostFromHash(hash) != cost), nil
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
	if strings.HasPrefix(hash, "$2y$") {
		arr.SetStringKey("algo", data.NewStringValue("2y"))
		arr.SetStringKey("algoName", data.NewStringValue("bcrypt"))
		opts := data.NewArrayValue(nil).(*data.ArrayValue)
		opts.SetStringKey("cost", data.NewIntValue(bcryptCostFromHash(hash)))
		arr.SetStringKey("options", opts)
	} else {
		arr.SetStringKey("algo", data.NewNullValue())
		arr.SetStringKey("algoName", data.NewStringValue("unknown"))
		arr.SetStringKey("options", data.NewArrayValue(nil))
	}
	return arr, nil
}

func arrayInt(av *data.ArrayValue, key string) (int, bool) {
	v := arrayAt(av, key)
	iv, ok := v.(data.AsInt)
	if !ok || iv == nil {
		return 0, false
	}
	n, err := iv.AsInt()
	if err != nil {
		return 0, false
	}
	return n, true
}

func arrayBool(av *data.ArrayValue, key string) (bool, bool) {
	v := arrayAt(av, key)
	bv, ok := v.(data.AsBool)
	if !ok || bv == nil {
		return false, false
	}
	b, err := bv.AsBool()
	if err != nil {
		return false, false
	}
	return b, true
}

func arrayAt(av *data.ArrayValue, key string) data.Value {
	if av == nil {
		return nil
	}
	slot, ok := av.LookupZValByStringKey(key)
	if !ok || slot == nil {
		return nil
	}
	return kit.Unwrap(slot.ReadValue())
}
