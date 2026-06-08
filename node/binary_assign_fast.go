package node

import "github.com/php-any/origami/data"

// intPairFromOperands 从两个操作数中同时读取整数值（用于 VarMulAssign / VarAddAssign 快速路径）。
func intPairFromOperands(ctx data.Context, left, right data.GetValue) (int, int, bool) {
	li, okL := intFromOperand(ctx, left)
	ri, okR := intFromOperand(ctx, right)
	return li, ri, okL && okR
}

// intFromOperand 从单个操作数中读取整数值：
//   - *VariableExpression：直接访问 ZVal，无接口调用
//   - *IntLiteral：直接读取预解析值
//   - 其它：降级为通用 GetValue
func intFromOperand(ctx data.Context, op data.GetValue) (int, bool) {
	switch o := op.(type) {
	case *VariableExpression:
		if zv := ctx.GetIndexZVal(o.Index); zv != nil {
			if iv, ok := zv.Value.(*data.IntValue); ok {
				return iv.Value, true
			}
		}
	case *IntLiteral:
		if iv, ok := o.V.(*data.IntValue); ok {
			return iv.Value, true
		}
	}
	v, ctl := op.GetValue(ctx)
	if ctl != nil {
		return 0, false
	}
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value, true
	}
	return 0, false
}
