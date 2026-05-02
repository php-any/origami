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
	case *data.ObjectValue:
		// ObjectValue -> 转为 ArrayValue 再操作
		tmp, _ := data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
		arr.RangeProperties(func(key string, v data.Value) bool {
			tmp.List = append(tmp.List, data.NewZVal(v))
			return true
		})
		arrZVal.Value = tmp
		source = &tmp.List
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

	// 获取要插入的元素（replacement 是第4个参数，一个数组）
	var insertElements []*data.ZVal
	if replacementArg, ok := ctx.GetIndexValue(3); ok && replacementArg != nil {
		// replacement 是一个数组，展开其元素插入
		if replacementArray, isArr := replacementArg.(*data.ArrayValue); isArr {
			for _, z := range replacementArray.List {
				insertElements = append(insertElements, z)
			}
		} else if replacementObj, isObj := replacementArg.(*data.ObjectValue); isObj {
			replacementObj.RangeProperties(func(key string, v data.Value) bool {
				insertElements = append(insertElements, data.NewZVal(v))
				return true
			})
		} else {
			// 单个值
			insertElements = append(insertElements, data.NewZVal(replacementArg))
		}
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
		node.NewParameter(nil, "offset", 1, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "length", 2, node.NewNullLiteral(nil), data.NewBaseType("int")),
		node.NewParameter(nil, "replacement", 3, node.NewNullLiteral(nil), data.NewBaseType("array")),
	}
}

func (f *ArraySpliceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "offset", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "replacement", 3, data.NewBaseType("array")),
	}
}
