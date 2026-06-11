package database

import (
	"github.com/php-any/origami/data"
)

// collectBindArgs 从上下文收集 SQL 绑定参数，支持两种写法：
//   - DB::sql($sql, $a, $b)          与 where 一致的可变参数
//   - DB::sql($sql, [$a, $b])        兼容旧版数组传参
func collectBindArgs(ctx data.Context, bindIndex int) []interface{} {
	paramValue, ok := ctx.GetIndexValue(bindIndex)
	if !ok {
		return nil
	}
	if _, isNull := paramValue.(*data.NullValue); isNull {
		return nil
	}

	if paramArray, ok := paramValue.(*data.ArrayValue); ok {
		valueList := flattenBindArgList(paramArray.ToValueList())
		args := make([]interface{}, len(valueList))
		for i, param := range valueList {
			args[i] = ConvertValueToGoType(param)
		}
		return args
	}
	if val, ok := paramValue.(data.Value); ok {
		return []interface{}{ConvertValueToGoType(val)}
	}
	return nil
}

// flattenBindArgList 若仅传入一个数组实参，则展开为多个绑定值。
func flattenBindArgList(list []data.Value) []data.Value {
	if len(list) == 1 {
		if nested, ok := list[0].(*data.ArrayValue); ok {
			return nested.ToValueList()
		}
	}
	return list
}
