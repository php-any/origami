package php

import (
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

func (f *PasswordHashFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "password", 0, nil, data.String{}),
		node.NewParameter(nil, "algo", 1, nil, data.Mixed{}),
		node.NewParameter(nil, "options", 2, node.NewNullLiteral(nil), data.Mixed{}),
	}
}

func (f *PasswordHashFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "password", 0, data.String{}),
		node.NewVariable(nil, "algo", 1, data.Mixed{}),
		node.NewVariable(nil, "options", 2, data.Mixed{}),
	}
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

func (f *PasswordVerifyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "password", 0, nil, data.String{}),
		node.NewParameter(nil, "hash", 1, nil, data.String{}),
	}
}

func (f *PasswordVerifyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "password", 0, data.String{}),
		node.NewVariable(nil, "hash", 1, data.String{}),
	}
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
}
