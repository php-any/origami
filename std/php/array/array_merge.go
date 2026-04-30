package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayMergeFunction() data.FuncStmt {
	return &ArrayMergeFunction{}
}

type ArrayMergeFunction struct{}

func (f *ArrayMergeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 Parameters 参数（包含所有传入的数组）
	paramsValue, _ := ctx.GetIndexValue(0)
	if paramsValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// Parameters 返回的是 ArrayValue，包含所有参数
	paramsArray, ok := paramsValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	allListValues := make([]data.Value, 0)
	resultAssoc := (*data.ObjectValue)(nil)

	paramsList := paramsArray.ToValueList()
	for _, paramValue := range paramsList {
		switch v := paramValue.(type) {
		case *data.ArrayValue:
			// 列表数组：依次取值
			values := v.ToValueList()
			if resultAssoc != nil {
				// 结果已经是关联数组：这些值追加成新的 int 键
				for _, val := range values {
					// 使用当前属性数量作为新键（转成字符串）
					key := len(resultAssoc.GetProperties())
					resultAssoc.SetProperty(fmt.Sprintf("%d", key), val)
				}
			} else {
				allListValues = append(allListValues, values...)
			}

		case *data.ObjectValue:
			// 关联数组：需要按键合并
			if resultAssoc == nil {
				resultAssoc = data.NewObjectValue()
				// 把之前累积的列表值先转成 0..n-1 的 int 键
				for i, val := range allListValues {
					resultAssoc.SetProperty(fmt.Sprintf("%d", i), val)
				}
			}
			for key, val := range v.GetProperties() {
				// 这里 key 是 string，保持为关联键；若有同名键，直接覆盖
				resultAssoc.SetProperty(key, val)
			}

		default:
			// 非数组参数，PHP 中会发 warning，我们这里简单按值附加
			if resultAssoc != nil {
				key := len(resultAssoc.GetProperties())
				resultAssoc.SetProperty(fmt.Sprintf("%d", key), v)
			} else {
				allListValues = append(allListValues, v)
			}
		}
	}

	if resultAssoc != nil {
		return resultAssoc, nil
	}
	return data.NewArrayValue(allListValues), nil
}

func (f *ArrayMergeFunction) GetName() string {
	return "array_merge"
}

func (f *ArrayMergeFunction) GetParams() []data.GetValue {
	// 使用可变参数
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayMergeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
