package node

import "github.com/php-any/origami/data"

// ConstantName 表示运行时解析的裸常量名（如 APP_ROOT）。
// define() 在执行期注册常量，解析期尚不可见，不能烘焙成字符串字面量。
type ConstantName struct {
	*Node `pp:"-"`
	Name  string
}

func NewConstantName(from data.From, name string) *ConstantName {
	return &ConstantName{
		Node: NewNode(from),
		Name: name,
	}
}

func (c *ConstantName) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if ctx != nil {
		if vm := ctx.GetVM(); vm != nil {
			if v, ok := vm.GetConstant(c.Name); ok {
				return v, nil
			}
		}
	}
	// 与历史行为一致：未知裸标识符回落为同名字符串
	return data.NewStringValue(c.Name), nil
}

func (c *ConstantName) AsString() string {
	return c.Name
}
