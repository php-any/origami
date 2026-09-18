package node

import (
	"github.com/php-any/origami/data"
)

// LambdaExpression 表示Lambda表达式（匿名函数）
type LambdaExpression struct {
	*FunctionStatement
	parent       map[int]int
	ctx          data.Context
	captured     map[int]data.Value // use ($x) 按值快照（创建闭包时），非 use (&$x)
	capturedRefs map[int]*data.ZVal // use (&$x) 在创建时绑定的共享 ZVal（对齐 PHP 引用捕获）
	IsStatic     bool               // static function/fn：不绑定定义处的 $this（对齐 PHP）
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
		if f.capturedRefs != nil {
			if zv, ok := f.capturedRefs[childIndex]; ok && zv != nil {
				value = zv.Value
			}
		}
		if value == nil && f.captured != nil {
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
	// PHP：use ($var) 在闭包创建时按值捕获；use (&$var) 绑定父作用域的同一 ZVal。
	// 必须在 GetValue 时固定引用槽：父函数返回后 Context 可能被还回 pool 并复用，
	// 调用时再 GetIndexZVal 会读到别的帧（Laravel FilesystemServiceProvider::serveFiles
	// 的 booted 回调里 $served[$uri] = $disk 因此失败）。
	// 闭包持有定义处 ctx（$this / self::），该帧禁止回收。
	markContextEscaped(ctx)
	captured := make(map[int]data.Value, len(f.parent))
	capturedRefs := make(map[int]*data.ZVal, len(f.parent))
	for cID, pID := range f.parent {
		if cID < 0 || cID >= len(f.vars) {
			continue
		}
		if _, isRef := f.vars[cID].(*VariableReference); isRef {
			name := ""
			if f.vars[cID] != nil {
				name = f.vars[cID].GetName()
			}
			zv := ctx.GetIndexZVal(pID)
			if zv == nil {
				zv = data.NewNamedZValSlot(name)
				ctx.SetIndexZVal(pID, zv)
			}
			capturedRefs[cID] = zv
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
		ctx:          ctx,
		parent:       f.parent,
		captured:     captured,
		capturedRefs: capturedRefs,
		IsStatic:     f.IsStatic,
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
	inner := ctx
	if bc, ok := ctx.(*data.BoundContext); ok {
		inner = bc.Context
	}
	execCtx := inner
	if defineClassCtx, ok := f.ctx.(*data.ClassMethodContext); ok && !f.IsStatic {
		execCtx = data.WrapMethodFrame(inner, defineClassCtx.ClassValue, defineClassCtx.SelfClass, defineClassCtx.StaticClass)
	} else if f.IsStatic {
		if defineClassCtx, ok := f.ctx.(*data.ClassMethodContext); ok {
			cmc := data.NewStaticMethodContext(inner, defineClassCtx.Class, defineClassCtx.StaticClass)
			if defineClassCtx.SelfClass != nil {
				cmc.SelfClass = defineClassCtx.SelfClass
			}
			if cmc.StaticClass == nil {
				cmc.StaticClass = defineClassCtx.Class
			}
			if cmc.SelfClass == nil {
				cmc.SelfClass = defineClassCtx.Class
			}
			execCtx = cmc
		}
	}
	// BoundContext 处理：
	// - ExplicitBind（Closure::bind/bindTo）：始终应用，可覆盖定义时的 $this。
	// - 调用方 CreateContext 继承的 BoundContext：仅当闭包定义时没有 $this 时才套上
	//   （否则会把 Livewire 视图里的 Login 污染到 ExtendBlade 的 EventBus 监听器）。
	if bc, ok := ctx.(*data.BoundContext); ok {
		_, hasDefThis := f.ctx.(*data.ClassMethodContext)
		if bc.ExplicitBind || !hasDefThis || f.IsStatic {
			execCtx = &data.BoundContext{
				Context:      execCtx,
				ScopeClass:   bc.ScopeClass,
				BoundThis:    bc.BoundThis,
				ExplicitBind: bc.ExplicitBind,
			}
		}
	}

	// 处理 use 捕获的外部变量
	for cID, pID := range f.parent {
		// use (&$var)：使用创建时绑定的共享 ZVal（不要在调用时再从父 ctx 取，父帧可能已回收/复用）
		if _, isRef := f.vars[cID].(*VariableReference); isRef {
			if f.capturedRefs != nil {
				if zv, ok := f.capturedRefs[cID]; ok && zv != nil {
					execCtx.SetIndexZVal(f.vars[cID].GetIndex(), zv)
					continue
				}
			}
			if f.ctx != nil {
				if parentZVal := f.ctx.GetIndexZVal(pID); parentZVal != nil {
					execCtx.SetIndexZVal(f.vars[cID].GetIndex(), parentZVal)
				}
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
		markContextEscaped(ctx)
		markContextEscaped(execCtx)
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

	persistStaticLocals(execCtx, f.vars)
	// PHP：普通闭包没有 return 时返回 null。箭头函数 fn() => expr 由解析器包成 return。
	return data.NewNullValue(), nil
}
