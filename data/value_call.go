package data

import "fmt"

func NewFuncValue(v FuncStmt) *FuncValue {
	return &FuncValue{
		Value: v,
	}
}

type FuncValue struct {
	Value FuncStmt
}

func (c *FuncValue) GetValue(ctx Context) (GetValue, Control) {
	return c, nil
}

func (c *FuncValue) Call(ctx Context) (GetValue, Control) {
	return c.Value.Call(ctx)
}

func (c *FuncValue) AsString() string {
	return fmt.Sprintf("%v", c.Value)
}
