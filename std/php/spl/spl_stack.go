package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplStackClass 实现 PHP SPL �?SplStack（LIFO），继承 SplDoublyLinkedList
type SplStackClass struct {
	node.Node
}

func NewSplStackClass() *SplStackClass {
	return &SplStackClass{}
}

func (c *SplStackClass) GetName() string { return "SplStack" }
func (c *SplStackClass) GetExtend() *string {
	parent := "SplDoublyLinkedList"
	return &parent
}
func (c *SplStackClass) GetImplements() []string {
	return []string{"Iterator", "Countable"}
}
func (c *SplStackClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplStackClass) GetPropertyList() []data.Property              { return nil }
func (c *SplStackClass) GetConstruct() data.Method                     { return nil }
func (c *SplStackClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	splListInitCV(cv)
	cv.SetProperty(splListModeKey, data.NewIntValue(SplITModeLIFO))
	return cv, nil
}

func (c *SplStackClass) GetMethod(name string) (data.Method, bool) {
	return splExtendGetMethod(c, name, nil)
}

func (c *SplStackClass) GetMethods() []data.Method {
	return splExtendGetMethods(c, nil)
}
