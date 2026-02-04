package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StdClass 表示 PHP 内置的 stdClass。
// 在 PHP 中，stdClass 是一个空类，通常作为通用对象使用。
type StdClass struct {
	node.Node
}

// GetValue 返回 stdClass 的一个实例。
func (s *StdClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx.CreateBaseContext()), nil
}

func (s *StdClass) GetFrom() data.From                            { return s.Node.GetFrom() }
func (s *StdClass) GetName() string                               { return "stdClass" }
func (s *StdClass) GetExtend() *string                            { return nil }
func (s *StdClass) GetImplements() []string                       { return nil }
func (s *StdClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (s *StdClass) GetPropertyList() []data.Property              { return nil }
func (s *StdClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (s *StdClass) GetMethods() []data.Method                     { return nil }
func (s *StdClass) GetConstruct() data.Method                     { return nil }

