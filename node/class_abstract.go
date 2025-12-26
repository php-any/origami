package node

import (
	"github.com/php-any/origami/data"
)

// AbstractClassStatement 表示抽象类定义语句
type AbstractClassStatement struct {
	*ClassStatement
}

// NewAbstractClassStatement 创建一个新的抽象类定义语句
func NewAbstractClassStatement(class *ClassStatement) *AbstractClassStatement {
	return &AbstractClassStatement{
		ClassStatement: class,
	}
}

// GetValue 获取抽象类定义语句的值
// 抽象类不能被实例化
func (c *AbstractClassStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return c.ClassStatement.GetValue(ctx)
}

// AddAnnotations 添加注解（实现 AddAnnotations 接口）
func (c *AbstractClassStatement) AddAnnotations(a *data.ClassValue) {
	c.ClassStatement.AddAnnotations(a)
}
