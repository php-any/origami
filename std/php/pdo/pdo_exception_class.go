package pdo

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/exception"
	"github.com/php-any/origami/token"
)

// PDOExceptionClass 表示 PHP 的 PDOException 类（extends RuntimeException）
type PDOExceptionClass struct {
	node.Node
	m exception.ExceptionMethods
}

func NewPDOExceptionClass() *PDOExceptionClass {
	return &PDOExceptionClass{
		m: exception.NewExceptionMethods(),
	}
}

func (c *PDOExceptionClass) GetName() string { return "PDOException" }

func (c *PDOExceptionClass) GetExtend() *string {
	s := "RuntimeException"
	return &s
}

func (c *PDOExceptionClass) GetImplements() []string { return nil }

func (c *PDOExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	newC := *c
	return data.NewClassValue(&newC, ctx), nil
}

func (c *PDOExceptionClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *PDOExceptionClass) GetPropertyList() []data.Property           { return nil }

func (c *PDOExceptionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return c.m.ConstructMethod, true
	case "error":
		return c.m.ErrorMethod, true
	case "getMessage":
		return c.m.GetMessageMethod, true
	case "getTraceAsString":
		return c.m.GetTraceMethod, true
	}
	return nil, false
}

func (c *PDOExceptionClass) GetMethods() []data.Method {
	return []data.Method{c.m.ErrorMethod, c.m.GetMessageMethod, c.m.GetTraceMethod}
}

func (c *PDOExceptionClass) GetConstruct() data.Method { return c.m.ConstructMethod }
