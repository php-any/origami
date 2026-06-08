package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ============================================================
// 解析阶段（NewBinaryAssign / NewPostfixIncr）识别常见整数赋值模式，
// 发出专用节点，将结构信息（变量索引、字面量值）在解析时固化，
// 运行时只需做值类型检查（*data.IntValue），无节点类型断言。
//
// 核心权衡：
//   - 使用单一具体类型 VarFastAssign（加 op 字节字段区分操作），
//     保证 for 循环体内接口调用位点单态（monomorphic），
//     CPU 分支预测可完全优化该间接跳转。
//   - 内部通过 readIdx 直接索引 ZVal 数组，无接口调用。
// ============================================================

// vfaOp 赋值操作枚举
type vfaOp byte

const (
	vfaOpCopy vfaOp = iota // $dst = $src
	vfaOpMul               // $dst = $lhs * $rhs
	vfaOpAdd               // $dst = $lhs + $rhs
)

// VarFastAssign 统一承载三种整数变量赋值模式。
// 所有变量操作数的 ZVal 索引在解析时预提取，字面量值直接存储，
// 运行时快速路径：GetIndexZVal（数组下标）+ *IntValue 断言 + AssignIntToZVal。
//
// 降级路径（非整数或复杂右侧）：调用 Slow（原始右侧节点）再走普通赋值。
type VarFastAssign struct {
	*Node  `pp:"-"`
	Dst    *VariableExpression
	DstIdx int
	LhsIdx int // ≥0: 变量索引；-1: 字面量（见 LhsLit）；-2: 复杂，需降级
	RhsIdx int
	LhsLit int
	RhsLit int
	Slow   data.GetValue // 原始右侧节点，供降级复用
	op     vfaOp
}

func (f *VarFastAssign) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	dstZv := ctx.GetIndexZVal(f.DstIdx)
	if dstZv != nil {
		switch f.op {
		case vfaOpCopy:
			// $dst = $src：LhsIdx 是 src 变量索引
			if srcZv := ctx.GetIndexZVal(f.LhsIdx); srcZv != nil {
				if iv, ok := srcZv.Value.(*data.IntValue); ok {
					data.AssignIntToZVal(dstZv, iv.Value)
					if dv, ok := dstZv.Value.(*data.IntValue); ok {
						return dv, nil
					}
				}
			}
		case vfaOpMul:
			if li, okL := readIdx(ctx, f.LhsIdx, f.LhsLit); okL {
				if ri, okR := readIdx(ctx, f.RhsIdx, f.RhsLit); okR {
					data.AssignIntToZVal(dstZv, li*ri)
					if iv, ok := dstZv.Value.(*data.IntValue); ok {
						return iv, nil
					}
				}
			}
		case vfaOpAdd:
			if li, okL := readIdx(ctx, f.LhsIdx, f.LhsLit); okL {
				if ri, okR := readIdx(ctx, f.RhsIdx, f.RhsLit); okR {
					data.AssignIntToZVal(dstZv, li+ri)
					if iv, ok := dstZv.Value.(*data.IntValue); ok {
						return iv, nil
					}
				}
			}
		}
	}

	// 降级路径
	rv, ctl := f.Slow.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if v, ok := rv.(data.Value); ok {
		return v, f.Dst.SetValue(ctx, v)
	}
	if rv == nil {
		v := data.NewNullValue()
		return v, f.Dst.SetValue(ctx, v)
	}
	return nil, data.NewErrorThrow(f.from, fmt.Errorf("VarFastAssign: unexpected result %T", rv))
}

// readIdx 读取预提取的操作数整数值：
//
//	idx ≥ 0：从变量 ZVal 读取（无节点类型断言，只做 *IntValue 值断言）
//	idx == -1：返回字面量值 lit
//	idx == -2：复杂节点，告知调用方走降级路径
func readIdx(ctx data.Context, idx, lit int) (int, bool) {
	if idx >= 0 {
		if zv := ctx.GetIndexZVal(idx); zv != nil {
			if iv, ok := zv.Value.(*data.IntValue); ok {
				return iv.Value, true
			}
		}
		return 0, false
	}
	if idx == -1 {
		return lit, true
	}
	// idx == -2：需要降级（复杂表达式作为操作数）
	return 0, false
}

// ------------------------------------------------------------------
// VarPostIncr: $var++（简单变量后自增，表达式语境）
// 预提取 VarIdx，运行时只做 *IntValue 值断言，无节点类型断言。
// 整数路径：复用旧指针作为返回值（原始值），只分配一个新 *IntValue。
// 由 NewPostfixIncr 发出；在 for 循环增量中由 NewForStatement 换为
// VarStmtIncr（0 次分配）。
// ------------------------------------------------------------------

type VarPostIncr struct {
	*Node    `pp:"-"`
	VarIdx   int
	Var      *VariableExpression
	Fallback *PostfixIncr
}

func (f *VarPostIncr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if zv := ctx.GetIndexZVal(f.VarIdx); zv != nil {
		if iv, ok := zv.Value.(*data.IntValue); ok {
			old := iv                                      // 旧指针直接作为返回值（原始值），无额外分配
			zv.Value = &data.IntValue{Value: iv.Value + 1} // 仅一次分配
			return old, nil
		}
	}
	return f.Fallback.GetValue(ctx)
}

// ------------------------------------------------------------------
// VarStmtIncr: for 循环增量专用自增（返回值始终被丢弃）
// 由 NewForStatement 在 AST 构建时将 VarPostIncr 替换为此节点。
//
// 关键优化：因为 for 循环永远丢弃增量表达式的返回值，
// 可以直接原地改写 *IntValue.Value（0 次分配）而不需要保存旧值。
// 安全保障：SetVariableValue 的 copy-on-assign 确保变量槽持有唯一
// 所有权的 *IntValue，不会影响其他变量。
// ------------------------------------------------------------------

type VarStmtIncr struct {
	*Node    `pp:"-"`
	VarIdx   int
	Fallback *PostfixIncr
}

func (f *VarStmtIncr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if zv := ctx.GetIndexZVal(f.VarIdx); zv != nil {
		if iv, ok := zv.Value.(*data.IntValue); ok {
			iv.Value++ // 原地改写，0 次分配
			return iv, nil
		}
	}
	return f.Fallback.GetValue(ctx)
}

// ------------------------------------------------------------------
// VarIntLe: $var <= IntLiteral（变量与整数字面量的 ≤ 比较）
// 由 NewBinaryLe 在解析阶段发出。
//
// 实现 BoolTest 接口，允许 ForStatement 直接获取 bool 值，
// 避免每次循环分配 *BoolValue（消除 1M 次堆分配）。
// ------------------------------------------------------------------

// BoolTest 是一个可选接口，由能直接返回 bool 的条件节点实现。
// ForStatement 在检测到此接口时调用 testBool，绕过 BoolValue 分配。
type BoolTest interface {
	testBool(ctx data.Context) (bool, data.Control)
}

// VarIntLe 针对 $var <= IntLiteral 模式，预提取变量索引和字面量值。
type VarIntLe struct {
	*Node  `pp:"-"`
	VarIdx int
	Lit    int
	Le     *BinaryLe // 供降级路径复用
}

// testBool 直接返回 bool，不分配 BoolValue。
func (f *VarIntLe) testBool(ctx data.Context) (bool, data.Control) {
	if zv := ctx.GetIndexZVal(f.VarIdx); zv != nil {
		if iv, ok := zv.Value.(*data.IntValue); ok {
			return iv.Value <= f.Lit, nil
		}
	}
	// 降级：调用原始节点再解包
	v, ctl := f.Le.GetValue(ctx)
	if ctl != nil {
		return false, ctl
	}
	if bv, ok := v.(*data.BoolValue); ok {
		return bv.Value, nil
	}
	if ab, ok := v.(data.AsBool); ok {
		b, err := ab.AsBool()
		if err != nil {
			return false, data.NewErrorThrow(f.from, err)
		}
		return b, nil
	}
	return v != nil, nil
}

// GetValue 保持兼容，在非 ForStatement 场景下仍返回 BoolValue。
func (f *VarIntLe) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	b, ctl := f.testBool(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(b), nil
}

// ------------------------------------------------------------------
// 解析时辅助：提取操作数的预编译信息
// ------------------------------------------------------------------

// preExtract 解析时分析操作数类型，返回 (idx, lit)：
//
//	VariableExpression → (varIndex, 0)
//	IntLiteral         → (-1, literalValue)
//	其他               → (-2, 0)   触发运行时降级
func preExtract(op data.GetValue) (idx, lit int) {
	switch o := op.(type) {
	case *VariableExpression:
		return o.Index, 0
	case *IntLiteral:
		if iv, ok := o.V.(*data.IntValue); ok {
			return -1, iv.Value
		}
	}
	return -2, 0
}
