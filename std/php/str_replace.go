package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrReplaceFunction() data.FuncStmt {
	return &StrReplaceFunction{}
}

type StrReplaceFunction struct{}

func (f *StrReplaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	searchValue, _ := ctx.GetIndexValue(0)
	replaceValue, _ := ctx.GetIndexValue(1)
	subjectValue, _ := ctx.GetIndexValue(2)
	countValue, _ := ctx.GetIndexValue(3)

	if subjectValue == nil {
		return data.NewStringValue(""), nil
	}

	subject := subjectValue.AsString()

	// 处理数组情况
	if searchArray, ok := searchValue.(*data.ArrayValue); ok {
		// search 和 replace 都是数组
		var replaceArray *data.ArrayValue
		if replaceValue != nil {
			if rArr, ok := replaceValue.(*data.ArrayValue); ok {
				replaceArray = rArr
			}
		}

		result := subject
		replaceCount := 0

		for i, searchItem := range searchArray.Value {
			searchStr := searchItem.AsString()
			replaceStr := ""
			if replaceArray != nil && i < len(replaceArray.Value) {
				replaceStr = replaceArray.Value[i].AsString()
			} else if replaceValue != nil {
				replaceStr = replaceValue.AsString()
			}

			// 执行替换
			newResult := strings.ReplaceAll(result, searchStr, replaceStr)
			if newResult != result {
				replaceCount++
			}
			result = newResult
		}

		// 处理 count 参数
		if countValue != nil {
			if _, ok := countValue.(*data.NullValue); !ok {
				if countRef, ok := countValue.(*data.ReferenceValue); ok {
					// count 是引用参数，需要更新
					parentCtx := countRef.Ctx
					varRef := countRef.Val
					parentCtx.SetVariableValue(varRef, data.NewIntValue(replaceCount))
				}
			}
		}

		return data.NewStringValue(result), nil
	}

	// 单个字符串替换
	search := ""
	if searchValue != nil {
		search = searchValue.AsString()
	}

	replace := ""
	if replaceValue != nil {
		replace = replaceValue.AsString()
	}

	// 如果 search 为空字符串，不执行替换（PHP 行为）
	if search == "" {
		result := subject
		// 处理 count 参数
		if countValue != nil {
			if _, ok := countValue.(*data.NullValue); !ok {
				if countRef, ok := countValue.(*data.ReferenceValue); ok {
					// count 是引用参数，需要更新
					replaceCount := 0
					parentCtx := countRef.Ctx
					varRef := countRef.Val
					parentCtx.SetVariableValue(varRef, data.NewIntValue(replaceCount))
				}
			}
		}
		return data.NewStringValue(result), nil
	}

	// 执行替换
	result := strings.ReplaceAll(subject, search, replace)

	// 处理 count 参数
	if countValue != nil {
		if _, ok := countValue.(*data.NullValue); !ok {
			if countRef, ok := countValue.(*data.ReferenceValue); ok {
				// count 是引用参数，需要更新
				// 计算替换次数
				replaceCount := strings.Count(subject, search)
				parentCtx := countRef.Ctx
				varRef := countRef.Val
				parentCtx.SetVariableValue(varRef, data.NewIntValue(replaceCount))
			}
		}
	}

	return data.NewStringValue(result), nil
}

func (f *StrReplaceFunction) GetName() string {
	return "str_replace"
}

func (f *StrReplaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "search", 0, nil, nil),
		node.NewParameter(nil, "replace", 1, nil, nil),
		node.NewParameter(nil, "subject", 2, nil, nil),
		node.NewParameter(nil, "count", 3, node.NewNullLiteral(nil), nil),
	}
}

func (f *StrReplaceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "search", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "replace", 1, data.NewBaseType("mixed")),
		node.NewVariable(nil, "subject", 2, data.NewBaseType("mixed")),
		node.NewVariable(nil, "count", 3, data.NewBaseType("int")),
	}
}
