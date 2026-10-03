package node

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

func NewClassClosure(class *data.ClassValue, methodName string) (*data.FuncValue, data.Control) {
	method, ok := class.GetMethod(methodName)
	if !ok {
		return nil, utils.NewThrow(errors.New("method not found"))
	}

	return data.NewFuncValue(&ClassClosure{
		class:  class,
		method: method,
	}), nil
}

// ClassClosure 基于对象的闭包
type ClassClosure struct {
	class  *data.ClassValue
	method data.Method
}

func (c *ClassClosure) Call(ctx data.Context) (data.GetValue, data.Control) {
	self := findDeclaringClassForMethod(ctx.GetVM(), c.class.Class, c.method.GetName())
	fnCtx := data.WrapMethodFrame(ctx, c.class, self, c.class.Class)
	return c.method.Call(fnCtx)
}

func (c *ClassClosure) RequestScopeObjects() []*data.ClassValue { return []*data.ClassValue{c.class} }
func (c *ClassClosure) BindRequestScope(ctx data.Context, objects map[*data.ObjectValue]*data.ClassValue) data.FuncStmt {
	clone := *c
	if target := objects[c.class.ObjectValue]; target != nil {
		clone.class = target
	}
	return &clone
}
func (c *ClassClosure) BindRequestCapture(scope *data.RequestCaptureScope) data.FuncStmt {
	clone := *c
	if scope.ScopeObject != nil {
		clone.class = scope.ScopeObject(c.class)
	} else if target := scope.Objects[c.class.ObjectValue]; target != nil {
		clone.class = target
	}
	return &clone
}

func (c *ClassClosure) GetName() string {
	return c.method.GetName()
}

func (c *ClassClosure) GetParams() []data.GetValue {
	return c.method.GetParams()
}

func (c *ClassClosure) GetVariables() []data.Variable {
	return c.method.GetVariables()
}
