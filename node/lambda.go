package node

import (
	"github.com/php-any/origami/data"
)

// LambdaExpression 表示Lambda表达式（匿名函数）
type LambdaExpression struct {
	*FunctionStatement
	parent   map[int]int
	ctx      data.Context
	captured map[int]data.Value // use ($x) 按值快照（创建闭包时），非 use (&$x)
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
	for childIndex, parentIndex := range f.parent {
		if childIndex < 0 || childIndex >= len(f.vars) {
			continue
		}
		var value data.Value
		if f.captured != nil {
			if v, ok := f.captured[childIndex]; ok {
				value = v
			}
		}
		if value == nil && f.ctx != nil {
			if v, ok := f.ctx.GetIndexValue(parentIndex); ok {
				value = v
			}
		}
		if value == nil {
			continue
		}
		result[f.vars[childIndex].GetName()] = value
	}
	return result
}

func (f *LambdaExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// PHP：use ($var) 在闭包创建时按值捕获；use (&$var) 在调用时共享 ZVal。
	// 必须在 GetValue 时快照，否则 foreach 里创建的多个闭包会全部看到循环变量终值
	//（Livewire ComponentHookRegistry 因此只注册到最后一个 Feature）。
	captured := make(map[int]data.Value, len(f.parent))
	for cID, pID := range f.parent {
		if cID < 0 || cID >= len(f.vars) {
			continue
		}
		if _, isRef := f.vars[cID].(*VariableReference); isRef {
			continue
		}
		v, ok := ctx.GetIndexValue(pID)
		if !ok || v == nil {
			continue
		}
		captured[cID] = snapshotUseValue(v)
	}
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
		ctx:      ctx,
		parent:   f.parent,
		captured: captured,
	}), nil
}

// snapshotUseValue 按值捕获：标量/数组拷贝；对象保留同一句柄（与 PHP 一致）。
func snapshotUseValue(v data.Value) data.Value {
	switch t := v.(type) {
	case *data.ArrayValue:
		return data.CloneArrayValue(t)
	case *data.IntValue:
		return data.NewIntValue(t.Value)
	case *data.FloatValue:
		return data.NewFloatValue(t.Value)
	case *data.StringValue:
		return data.NewStringValue(t.Value)
	case *data.BoolValue:
		return data.NewBoolValue(t.Value)
	case *data.NullValue:
		return data.NewNullValue()
	default:
		return v
	}
}

func (f *LambdaExpression) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 为 lambda 创建独立的执行上下文，避免直接复用调用方 ctx 而污染上层环境。
	var execCtx data.Context
	if defineClassCtx, ok := f.ctx.(*data.ClassMethodContext); ok {
		// 在类方法中定义的 lambda：使用定义时对象创建新的 ClassMethodContext 作为执行上下文，
		// 以保证 this / static:: 语义正确。
		execCtx = defineClassCtx.ClassValue.CreateContext(f.vars)
		if cmc, ok := execCtx.(*data.ClassMethodContext); ok {
			cmc.StaticClass = defineClassCtx.StaticClass
			cmc.SelfClass = defineClassCtx.SelfClass
			if cmc.StaticClass == nil {
				cmc.StaticClass = defineClassCtx.Class
			}
			if cmc.SelfClass == nil {
				cmc.SelfClass = defineClassCtx.Class
			}
		}
		// VM 属于本次调用而不是闭包定义作用域。常驻对象中的闭包可能跨请求复用，
		// 但输出、HTTP 和调用栈必须落到当前调用方 VM。
		execCtx.SetVM(ctx.GetVM())
	} else {
		// 普通场景：基于当前 ctx 再创建一层函数上下文，隔离变量写入。
		execCtx = ctx.CreateContext(f.vars)
	}
	// 仅当本次调用本身由 BoundFuncValue 注入 BoundContext 时保留绑定。
	// 不可沿父上下文链向上找：HTTP 请求里 Request macro 等会在祖先留下 BoundContext，
	// 误套到无关闭包上会导致 static::/self:: 类型断言失败（Livewire Utils 回调）。
	if bc, ok := ctx.(*data.BoundContext); ok {
		execCtx = &data.BoundContext{Context: execCtx, ScopeClass: bc.ScopeClass, BoundThis: bc.BoundThis}
	}
	// 将调用方 ctx 中已经绑定好的参数 ZVal 复制到新的执行上下文中
	for i := range f.vars {
		zv := ctx.GetIndexZVal(i)
		if zv != nil {
			execCtx.SetIndexZVal(i, zv)
		}
	}

	// 处理 use 捕获的外部变量
	for cID, pID := range f.parent {
		// use (&$var)：调用时共享父作用域 ZVal
		if _, isRef := f.vars[cID].(*VariableReference); isRef {
			parentZVal := f.ctx.GetIndexZVal(pID)
			if parentZVal != nil {
				execCtx.SetIndexZVal(f.vars[cID].GetIndex(), parentZVal)
			}
			continue
		}

		// use ($var)：使用创建闭包时的按值快照
		if f.captured != nil {
			if v, ok := f.captured[cID]; ok && v != nil {
				execCtx.SetVariableValue(f.vars[cID], v)
				continue
			}
		}
		// 兼容未快照的旧路径
		if v, ok := f.ctx.GetIndexValue(pID); ok {
			execCtx.SetVariableValue(f.vars[cID], v)
		}
	}

	// PHP：含 yield 的闭包调用时立即返回 Generator（Symfony TableRows 依赖）
	if f.IsGenerator {
		generator := NewFuncYieldStackState(execCtx, f, f.Body, 0, nil, nil)
		generatorClass := NewGeneratorClass(generator)
		return generatorClass.GetValue(execCtx)
	}

	var ctl data.Control
	for bodyIndex := 0; bodyIndex < len(f.Body); bodyIndex++ {
		statement := f.Body[bodyIndex]
		_, ctl = statement.GetValue(execCtx)
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

	// PHP：普通闭包没有 return 时返回 null。箭头函数 fn() => expr 由解析器包成 return。
	return data.NewNullValue(), nil
}
