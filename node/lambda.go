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
			Node:        NewNode(from),
			Params:      params,
			Body:        body,
			vars:        vars,
			IsGenerator: containsYield(body),
		},
		parent: parent,
	}
}

// GetParentBindings 返回 use 捕获的父作用域变量索引映射（编译期代码生成使用）
func (f *LambdaExpression) GetParentBindings() map[int]int {
	return f.parent
}

// GetStaticVariables 返回闭包 use 捕获的变量，对齐 ReflectionFunction::getStaticVariables()。
func (f *LambdaExpression) GetStaticVariables() map[string]data.Value {
	result := make(map[string]data.Value, len(f.parent))
	if f.ctx == nil {
		return result
	}
	for childIndex, parentIndex := range f.parent {
		if childIndex < 0 || childIndex >= len(f.vars) {
			continue
		}
		value, ok := f.ctx.GetIndexValue(parentIndex)
		if !ok || value == nil {
			continue
		}
		result[f.vars[childIndex].GetName()] = value
	}
	return result
}

func (f *LambdaExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewFuncValue(&LambdaExpression{
		FunctionStatement: &FunctionStatement{
			Node:        f.Node,
			Params:      f.Params,
			Body:        f.Body,
			vars:        f.vars,
			IsGenerator: f.IsGenerator,
			Name:        f.Name,
			Ret:         f.Ret,
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
		// VM 属于本次调用而不是闭包定义作用域。常驻对象中的闭包可能跨请求复用，
		// 但输出、HTTP 和调用栈必须落到当前调用方 VM。
		execCtx.SetVM(ctx.GetVM())
	} else {
		// 普通场景：基于当前 ctx 再创建一层函数上下文，隔离变量写入。
		execCtx = ctx.CreateContext(f.vars)
	}
	// 保留 BoundContext（来自 Closure::bind/bindTo）以便闭包内可以访问 $this 与私有成员
	if bc := getBoundContext(ctx); bc != nil {
		execCtx = &data.BoundContext{Context: execCtx, ScopeClass: bc.ScopeClass, BoundThis: bc.BoundThis}
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

	// PHP：含 yield 的闭包调用时立即返回 Generator（Symfony TableRows 依赖）
	if f.IsGenerator {
		generator := NewFuncYieldStackState(execCtx, f, f.Body, 0, nil, nil)
		generatorClass := NewGeneratorClass(generator)
		return generatorClass.GetValue(execCtx)
	}

	var v data.GetValue
	var ctl data.Control
	for bodyIndex := 0; bodyIndex < len(f.Body); bodyIndex++ {
		statement := f.Body[bodyIndex]
		v, ctl = statement.GetValue(execCtx)
		if ctl != nil {
			switch rv := ctl.(type) {
			case data.ExitControl:
				return nil, ctl
			case data.ReturnControl:
				return rv.ReturnValue(), nil
			case data.GotoControl:
				offset, acl := resolveGotoBodyIndex(f.from, f.Body, rv)
				if acl != nil {
					return nil, acl
				}
				bodyIndex = offset - 1
				continue
			case LabelControl:
				continue
			case data.YieldControl:
				generator := rv.CreateStackState(execCtx, f, f.Body, bodyIndex)
				generatorClass := NewGeneratorClass(generator)
				return generatorClass.GetValue(execCtx)
			case data.YieldValueControl:
				generator := NewFuncYieldStackState(execCtx, f, f.Body, bodyIndex+1, rv.GetYieldKey(), rv.GetYieldValue())
				generatorClass := NewGeneratorClass(generator)
				return generatorClass.GetValue(execCtx)
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

// getBoundContext 从上下文链中查找 BoundContext（来自 Closure::bind/bindTo）
func getBoundContext(ctx data.Context) *data.BoundContext {
	return data.FindBoundContext(ctx)
}
