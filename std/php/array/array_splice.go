package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArraySpliceFunction() data.FuncStmt {
	return &ArraySpliceFunction{}
}

type ArraySpliceFunction struct{}

func (f *ArraySpliceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 第一个参数：array (引用)
	arrayVal, ok := ctx.GetIndexValue(0)
	if !ok || arrayVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取 array 的 ZVal 引用以支持按引用修改
	arrZVal := ctx.GetIndexZVal(0)
	if arrZVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	var source *[]*data.ZVal
	switch arr := arrZVal.Value.(type) {
	case *data.ArrayValue:
		source = &arr.List
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	start := 0
	deleteCount := len(*source)

	// 获取 offset 参数
	if startArg, ok := ctx.GetIndexValue(1); ok && startArg != nil {
		if startInt, ok := startArg.(data.AsInt); ok {
			if s, err := startInt.AsInt(); err == nil {
				start = s
			}
		}
	}

	// 获取 length 参数
	if lengthArg, ok := ctx.GetIndexValue(2); ok && lengthArg != nil {
		if lengthInt, ok := lengthArg.(data.AsInt); ok {
			if d, err := lengthInt.AsInt(); err == nil {
				deleteCount = d
			}
		}
	}

	// 处理负数索引
	if start < 0 {
		start = len(*source) + start
	}
	if start < 0 {
		start = 0
	}
	if start > len(*source) {
		start = len(*source)
	}
	if deleteCount < 0 {
		deleteCount = 0
	}
	if start+deleteCount > len(*source) {
		deleteCount = len(*source) - start
	}

	// 获取要删除的元素
	deletedElements := make([]*data.ZVal, deleteCount)
	copy(deletedElements, (*source)[start:start+deleteCount])

	// 获取要插入的元素（从第3个参数开始）
	var insertElements []*data.ZVal
	for i := 3; ; i++ {
		arg, ok := ctx.GetIndexValue(i)
		if !ok || arg == nil {
			break
		}
		insertElements = append(insertElements, data.NewZVal(arg))
	}

	// 执行 splice 操作
	newArray := make([]*data.ZVal, 0, len(*source)-deleteCount+len(insertElements))
	newArray = append(newArray, (*source)[:start]...)
	newArray = append(newArray, insertElements...)
	newArray = append(newArray, (*source)[start+deleteCount:]...)
	*source = newArray

	// 返回被删除的元素
	deletedValues := make([]data.Value, len(deletedElements))
	for i, zval := range deletedElements {
		deletedValues[i] = zval.Value
	}
	return data.NewArrayValue(deletedValues), nil
}

func (f *ArraySpliceFunction) GetName() string {
	return "array_splice"
}

func (f *ArraySpliceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.NewBaseType("array")),
		node.NewParameter(nil, "offset", 1, nil, nil),
		node.NewParameters(nil, "replacements", 2, nil, nil),
	}
}

func (f *ArraySpliceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "offset", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "replacements", 2, data.NewBaseType("array")),
	}
}
