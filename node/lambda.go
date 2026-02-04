package node

import (
	"github.com/php-any/origami/data"
)

// LambdaExpression 表示Lambda表达式（匿名函数）
type LambdaExpression struct {
	*FunctionStatement
	parent map[int]int
	ctx    data.Context
}

// NewLambdaExpression 创建一个新的Lambda表达式
func NewLambdaExpression(from data.From, params []data.GetValue, body []data.GetValue, vars []data.Variable, parent map[int]int) *LambdaExpression {
	return &LambdaExpression{
		FunctionStatement: &FunctionStatement{
			Node:   NewNode(from),
			Params: params,
			Body:   body,
			vars:   vars,
		},
		parent: parent,
	}
}

func (f *LambdaExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewFuncValue(&LambdaExpression{
		FunctionStatement: &FunctionStatement{
			Node:   f.Node,
			Params: f.Params,
			Body:   f.Body,
			vars:   f.vars,
		},
		ctx:    ctx,
		parent: f.parent,
	}), nil
}

func (f *LambdaExpression) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 为 lambda 创建独立的执行上下文，避免直接复用调用方 ctx 而污染上层环境。
	var execCtx data.Context
	if defineClassCtx, ok := f.ctx.(*data.ClassMethodContext); ok {
		// 在类方法中定义的 lambda：使用定义时对象创建新的 ClassMethodContext 作为执行上下文，
		// 以保证 this 语义正确。
		execCtx = defineClassCtx.ClassValue.CreateContext(f.vars)
	} else {
		// 普通场景：基于当前 ctx 再创建一层函数上下文，隔离变量写入。
		execCtx = ctx.CreateContext(f.vars)
	}
	// 将调用方 ctx 中已经绑定好的参数 ZVal 复制到新的执行上下文中
	for i := range f.vars {
		zv := ctx.GetIndexZVal(i)
		if zv != nil {
			execCtx.SetIndexZVal(i, zv)
		}
	}

	// 处理 use 捕获的外部变量：从定义时上下文 f.ctx 读取，写入 execCtx
	for cID, pID := range f.parent {
		v, ok := f.ctx.GetIndexValue(pID)
		if !ok {
			continue
		}

		// 如果子变量是 VariableReference，说明是 use (&$var) 按引用捕获：
		// 直接让子变量槽引用父作用域同一个 ZVal，实现引用语义，
		// 避免使用已标记 deprecated 的 data.NewReferenceValue。
		if _, isRef := f.vars[cID].(*VariableReference); isRef {
			parentZVal := f.ctx.GetIndexZVal(pID)
			if parentZVal != nil {
				execCtx.SetIndexZVal(f.vars[cID].GetIndex(), parentZVal)
			}
		} else {
			// 普通按值捕获
			execCtx.SetVariableValue(f.vars[cID], v)
		}
	}

	var v data.GetValue
	var ctl data.Control
	for _, statement := range f.Body {
		v, ctl = statement.GetValue(execCtx)
		if ctl != nil {
			switch rv := ctl.(type) {
			case data.ReturnControl:
				return rv.ReturnValue(), nil
			case data.AddStack:
				switch call := statement.(type) {
				case *CallExpression:
					rv.AddStackWithInfo(call.from, "", call.FunName)
				case *CallObjectMethod:
					rv.AddStackWithInfo(call.from, "->", call.Method)
				}
				rv.AddStackWithInfo(f.from, "lambda", f.Name)
			}
			return nil, ctl
		}
	}

	return v, nil
}
