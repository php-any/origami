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
	fnCtx := c.class.CreateContext(c.method.GetVariables())

	for i := 0; i < len(c.method.GetVariables()); i++ {
		fnCtx.SetIndexZVal(i, ctx.GetIndexZVal(i))
	}

	return c.method.Call(fnCtx)
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
