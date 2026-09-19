package node

import (
	"github.com/php-any/origami/data"
)

// FuncGetArgs 表示 func_get_args 关键字表达式
type FuncGetArgs struct {
	*Node `pp:"-"`
}

// NewFuncGetArgs 创建一个新的 func_get_args 表达式
func NewFuncGetArgs(from data.From) data.GetValue {
	return &FuncGetArgs{
		Node: NewNode(from),
	}
}

// GetValue 获取 func_get_args 的值（返回所有函数参数）
// 使用 GetIndexValue 按索引读取已求值的参数，避免在错误上下文中重新求值表达式
func (f *FuncGetArgs) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	values, acl := expandCallArgsValues(ctx)
	if acl != nil {
		return nil, acl
	}
	return data.NewArrayValue(values), nil
}

// collectCallArgValues 将调用方传入的位置实参表达式列表展开为扁平值列表（用于 func_get_args）。
// 普通实参按位置取值；...$arr 展开实参按数组元素逐项取值。
func collectCallArgValues(ctx data.Context, positional []data.GetValue, fnCtx data.Context) []data.Value {
	values := make([]data.Value, 0, len(positional))
	idx := 0
	for _, arg := range positional {
		if spread, ok := arg.(*SpreadArgument); ok && spread.Expr != nil {
			spreadVal, acl := spread.GetValue(ctx)
			if acl != nil {
				continue
			}
			vals, spreadCtl := spreadToValues(ctx, spreadVal)
			if spreadCtl != nil {
				continue
			}
			values = append(values, vals...)
			idx += len(vals)
			continue
		}
		v, ok := fnCtx.GetIndexValue(idx)
		if ok && v != nil {
			values = append(values, v)
		} else {
			values = append(values, data.NewNullValue())
		}
		idx++
	}
	return values
}

// expandCallArgsValues 将本次调用的实参展开为扁平值列表。
// 优先使用调用时记录的扁平实参值（已处理 ...$arr 展开，避免递归求值 func_get_args），
// 否则退回按实参表达式逐个读取。
func expandCallArgsValues(ctx data.Context) ([]data.Value, data.Control) {
	if flat := ctx.GetFlatCallArgs(); flat != nil {
		return flat, nil
	}

	args := ctx.GetCallArgs()
	if args == nil || len(args) == 0 {
		return nil, nil
	}

	values := make([]data.Value, 0, len(args))
	idx := 0
	for _, arg := range args {
		if _, ok := arg.(*SpreadArgument); ok {
			// ...$var 的表达式槽位属于调用方。在被调帧 GetValue 会按调用方 Index
			// 读被调符号表，Go 切片越界。展开结果只能来自 FlatCallArgs。
			v, ok := ctx.GetIndexValue(idx)
			if ok && v != nil {
				values = append(values, v)
			}
			idx++
			continue
		}
		v, ok := ctx.GetIndexValue(idx)
		if ok && v != nil {
			values = append(values, v)
		} else {
			values = append(values, data.NewNullValue())
		}
		idx++
	}

	return values, nil
}
