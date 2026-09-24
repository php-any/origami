package php

import (
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	passwordDefault = 1
	passwordBcrypt  = 1
)

// PasswordHashFunction 实现 password_hash。
type PasswordHashFunction struct{}

func NewPasswordHashFunction() data.FuncStmt { return &PasswordHashFunction{} }

func (f *PasswordHashFunction) GetName() string { return "password_hash" }

var passwordHashFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "password", 0, nil, data.String{}),
	node.NewParameter(nil, "algo", 1, nil, data.Mixed{}),
	node.NewParameter(nil, "options", 2, node.NewNullLiteral(nil), data.Mixed{}),
}

func (f *PasswordHashFunction) GetParams() []data.GetValue {
	return passwordHashFunctionGetParams
}

var passwordHashFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "password", 0, data.String{}),
	node.NewVariable(nil, "algo", 1, data.Mixed{}),
	node.NewVariable(nil, "options", 2, data.Mixed{}),
}

func (f *PasswordHashFunction) GetVariables() []data.Variable {
	return passwordHashFunctionGetVariables
}

func (f *PasswordHashFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pwdVal, _ := ctx.GetIndexValue(0)
	pwd := ""
	if pwdVal != nil {
		pwd = pwdVal.AsString()
	}
	cost := bcrypt.DefaultCost
	if optVal, ok := ctx.GetIndexValue(2); ok && optVal != nil {
		var costVal data.Value
		switch o := optVal.(type) {
		case *data.ArrayValue:
			if c, ctl := o.GetProperty("cost"); ctl == nil {
				costVal = c
			}
		case *data.ObjectValue:
			if c, ctl := o.GetProperty("cost"); ctl == nil {
				costVal = c
			}
		}
		if costVal != nil {
			if asInt, ok := costVal.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil && n >= bcrypt.MinCost && n <= bcrypt.MaxCost {
					cost = n
				}
			}
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(string(hash)), nil
}

// PasswordVerifyFunction 实现 password_verify。
type PasswordVerifyFunction struct{}

func NewPasswordVerifyFunction() data.FuncStmt { return &PasswordVerifyFunction{} }

func (f *PasswordVerifyFunction) GetName() string { return "password_verify" }

var passwordVerifyFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "password", 0, nil, data.String{}),
	node.NewParameter(nil, "hash", 1, nil, data.String{}),
}

func (f *PasswordVerifyFunction) GetParams() []data.GetValue {
	return passwordVerifyFunctionGetParams
}

var passwordVerifyFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "password", 0, data.String{}),
	node.NewVariable(nil, "hash", 1, data.String{}),
}

func (f *PasswordVerifyFunction) GetVariables() []data.Variable {
	return passwordVerifyFunctionGetVariables
}

func (f *PasswordVerifyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pwdVal, _ := ctx.GetIndexValue(0)
	hashVal, _ := ctx.GetIndexValue(1)
	pwd, hash := "", ""
	if pwdVal != nil {
		pwd = pwdVal.AsString()
	}
	if hashVal != nil {
		hash = hashVal.AsString()
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return data.NewBoolValue(err == nil), nil
}

func loadPasswordFunctions(vm data.VM) {
	vm.SetConstant("PASSWORD_DEFAULT", data.NewIntValue(passwordDefault))
	vm.SetConstant("PASSWORD_BCRYPT", data.NewIntValue(passwordBcrypt))
	_ = vm.AddFunc(NewPasswordHashFunction())
	_ = vm.AddFunc(NewPasswordVerifyFunction())
	_ = vm.AddFunc(NewPasswordGetInfoFunction())
	_ = vm.AddFunc(NewPasswordNeedsRehashFunction())
}

func passwordBcryptCost(hash string) int {
	parts := strings.Split(hash, "$")
	if len(parts) >= 3 {
		n, _ := strconv.Atoi(parts[2])
		return n
	}
	return 0
}

func isBcryptHash(hash string) bool {
	return strings.HasPrefix(hash, "$2y$") || strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$")
}

// PasswordGetInfoFunction 实现 password_get_info。
type PasswordGetInfoFunction struct{}

func NewPasswordGetInfoFunction() data.FuncStmt { return &PasswordGetInfoFunction{} }

func (f *PasswordGetInfoFunction) GetName() string { return "password_get_info" }

var passwordGetInfoFunctionGetParams = []data.GetValue{node.NewParameter(nil, "hash", 0, nil, data.String{})}

func (f *PasswordGetInfoFunction) GetParams() []data.GetValue {
	return passwordGetInfoFunctionGetParams
}

var passwordGetInfoFunctionGetVariables = []data.Variable{node.NewVariable(nil, "hash", 0, data.String{})}

func (f *PasswordGetInfoFunction) GetVariables() []data.Variable {
	return passwordGetInfoFunctionGetVariables
}

func (f *PasswordGetInfoFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hash := ""
	if v, _ := ctx.GetIndexValue(0); v != nil {
		hash = v.AsString()
	}
	algo := 0
	algoName := "unknown"
	cost := 0
	if isBcryptHash(hash) {
		algo = passwordBcrypt
		algoName = "bcrypt"
		cost = passwordBcryptCost(hash)
	}
	options := []*data.ZVal{}
	if cost > 0 {
		options = []*data.ZVal{data.NewNamedZVal("cost", data.NewIntValue(cost))}
	}
	return &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("algo", data.NewIntValue(algo)),
			data.NewNamedZVal("algoName", data.NewStringValue(algoName)),
			data.NewNamedZVal("options", &data.ArrayValue{List: options}),
		},
	}, nil
}

// PasswordNeedsRehashFunction 实现 password_needs_rehash。
type PasswordNeedsRehashFunction struct{}

func NewPasswordNeedsRehashFunction() data.FuncStmt { return &PasswordNeedsRehashFunction{} }

func (f *PasswordNeedsRehashFunction) GetName() string { return "password_needs_rehash" }

var passwordNeedsRehashFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "hash", 0, nil, data.String{}),
	node.NewParameter(nil, "algo", 1, nil, data.Mixed{}),
	node.NewParameter(nil, "options", 2, node.NewNullLiteral(nil), data.Mixed{}),
}

func (f *PasswordNeedsRehashFunction) GetParams() []data.GetValue {
	return passwordNeedsRehashFunctionGetParams
}

var passwordNeedsRehashFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "hash", 0, data.String{}),
	node.NewVariable(nil, "algo", 1, data.Mixed{}),
	node.NewVariable(nil, "options", 2, data.Mixed{}),
}

func (f *PasswordNeedsRehashFunction) GetVariables() []data.Variable {
	return passwordNeedsRehashFunctionGetVariables
}

func (f *PasswordNeedsRehashFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hash := ""
	if v, _ := ctx.GetIndexValue(0); v != nil {
		hash = v.AsString()
	}
	if !isBcryptHash(hash) {
		return data.NewBoolValue(true), nil
	}
	wantCost := bcrypt.DefaultCost
	if optVal, ok := ctx.GetIndexValue(2); ok && optVal != nil {
		var costVal data.Value
		switch o := optVal.(type) {
		case *data.ArrayValue:
			if c, ctl := o.GetProperty("cost"); ctl == nil {
				costVal = c
			}
		case *data.ObjectValue:
			if c, ctl := o.GetProperty("cost"); ctl == nil {
				costVal = c
			}
		}
		if costVal != nil {
			if asInt, ok := costVal.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil && n > 0 {
					wantCost = n
				}
			}
		}
	}
	return data.NewBoolValue(passwordBcryptCost(hash) != wantCost), nil
}
