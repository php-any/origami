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
	for cID, pID := range f.parent {
		v, ok := f.ctx.GetIndexValue(pID)
		if !ok {
			continue
		}

		// 如果子变量是 VariableReference，说明是 use (&$var) 按引用捕获
		if _, isRef := f.vars[cID].(*VariableReference); isRef {
			// 创建一个“父级变量”描述，索引指向父作用域中的变量
			parentVar := NewVariable(
				nil,
				f.vars[cID].GetName(),
				pID,
				f.vars[cID].GetType(),
			)
			refValue := data.NewReferenceValue(parentVar, f.ctx)
			ctx.SetVariableValue(f.vars[cID], refValue)
		} else {
			// 普通按值捕获
			ctx.SetVariableValue(f.vars[cID], v)
		}
	}

	var v data.GetValue
	var ctl data.Control
	for _, statement := range f.Body {
		v, ctl = statement.GetValue(ctx)
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
