package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayKeysFunction 实现 array_keys 函数
// 返回数组中所有的键（对普通数组为 0..n-1，对关联数组为属性名）
type ArrayKeysFunction struct{}

func NewArrayKeysFunction() data.FuncStmt {
	return &ArrayKeysFunction{}
}

func (f *ArrayKeysFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：数组
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return nil, throwMustBeArray("array_keys", nil)
	}

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		filtered := len(ctx.GetFlatCallArgs()) > 1 || len(ctx.GetCallArgs()) > 1
		search, _ := ctx.GetIndexValue(1)
		strict := false
		if flag, ok := ctx.GetIndexValue(2); ok {
			if b, ok := flag.(data.AsBool); ok {
				strict, _ = b.AsBool()
			}
		}
		eq := node.NewBinaryEq(nil, search, nil)
		identical := node.NewBinaryEqStrict(nil, search, nil)
		// 处理数组，检查是否有字符串键（Name 字段）
		length := v.Len()
		keys := make([]data.Value, 0, length)
		for i := 0; i < length; i++ {
			if filtered {
				var comparison data.GetValue
				var ctl data.Control
				if strict {
					identical.Right = v.At(i).ReadValue()
					comparison, ctl = identical.GetValue(ctx)
				} else {
					eq.Right = v.At(i).ReadValue()
					comparison, ctl = eq.GetValue(ctx)
				}
				if ctl != nil {
					return nil, ctl
				}
				if !comparison.(*data.BoolValue).Value {
					continue
				}
			}
			keys = append(keys, v.At(i).PHPArrayKey(i))
		}
		return data.NewArrayValue(keys), nil

	default:
		return nil, throwMustBeArray("array_keys", arrayValue)
	}
}

func (f *ArrayKeysFunction) GetName() string {
	return "array_keys"
}

var arrayKeysFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, nil),
	node.NewParameter(nil, "filter_value", 1, data.NewNullValue(), nil),
	node.NewParameter(nil, "strict", 2, data.NewBoolValue(false), nil),
}

func (f *ArrayKeysFunction) GetParams() []data.GetValue {
	return arrayKeysFunctionGetParams
}

var arrayKeysFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	node.NewVariable(nil, "filter_value", 1, nil),
	node.NewVariable(nil, "strict", 2, nil),
}

func (f *ArrayKeysFunction) GetVariables() []data.Variable {
	return arrayKeysFunctionGetVariables
}
