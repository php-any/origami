package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

// LambdaExpression 表示Lambda表达式（匿名函数）
type LambdaExpression struct {
	*FunctionStatement
	parent         map[int]int
	ctx            data.Context
	captured       map[int]data.Value // use ($x) 按值快照（创建闭包时），非 use (&$x)
	capturedRefs   map[int]*data.ZVal // use (&$x) 在创建时绑定的共享 ZVal（对齐 PHP 引用捕获）
	IsStatic       bool               // static function/fn：不绑定定义处的 $this（对齐 PHP）
	parameterScope bool
}

// RequestScopeObjects walks retained captures at the request setup boundary.
// Executing the closure still uses its fixed PHP object handles directly.
func (f *LambdaExpression) RequestScopeObjects() []*data.ClassValue {
	var objects []*data.ClassValue
	seenValues := map[data.Value]bool{}
	seenClosures := map[*LambdaExpression]bool{}
	var visit func(data.Value)
	var closure func(*LambdaExpression)
	closure = func(lambda *LambdaExpression) {
		if lambda == nil || seenClosures[lambda] {
			return
		}
		seenClosures[lambda] = true
		if owner, ok := lambda.ctx.(*data.ClassMethodContext); ok && !lambda.IsStatic && owner.ObjectValue != nil {
			visit(owner.ClassValue)
		}
		for _, value := range lambda.captured {
			visit(value)
		}
		for _, slot := range lambda.capturedRefs {
			if slot != nil {
				visit(slot.ReadValue())
			}
		}
	}
	visit = func(value data.Value) {
		if value == nil || seenValues[value] {
			return
		}
		seenValues[value] = true
		switch v := value.(type) {
		case *data.ClassValue:
			if v.ObjectValue != nil {
				objects = append(objects, v)
			}
		case *data.ThisValue:
			visit(v.ClassValue)
		case *data.ArrayValue:
			for _, slot := range v.Range() {
				if slot != nil {
					visit(slot.ReadValue())
				}
			}
		case *data.FuncValue:
			if lambda, ok := v.Value.(*LambdaExpression); ok {
				closure(lambda)
			}
		case *data.BoundFuncValue:
			visit(v.BoundObject)
			if lambda, ok := v.Value.(*LambdaExpression); ok {
				closure(lambda)
			}
		case *data.ArraySlotRef:
			if v.Slot != nil {
				visit(v.Slot.ReadValue())
			}
		}
	}
	closure(f)
	return objects
}

func (f *LambdaExpression) BindRequestScope(ctx data.Context, objects map[*data.ObjectValue]*data.ClassValue) data.FuncStmt {
	return f.BindRequestCapture(data.NewRequestCaptureScope(ctx, objects))
}

func (f *LambdaExpression) BindRequestCapture(scope *data.RequestCaptureScope) data.FuncStmt {
	if clone := scope.Closure(f); clone != nil {
		return clone
	}
	clone := *f
	scope.RememberClosure(f, &clone)
	if owner, ok := f.ctx.(*data.ClassMethodContext); ok && owner.ObjectValue != nil {
		scoped := scope.Objects[owner.ObjectValue]
		if scoped == nil && scope.ScopeObject != nil && !f.IsStatic {
			scoped = scope.ScopeObject(owner.ClassValue)
		}
		if scoped != nil {
			clone.ctx = data.WrapMethodFrame(scope.Context.CreateBaseContext(), scoped, owner.SelfClass, owner.StaticClass)
		}
	}
	if len(f.captured) > 0 {
		clone.captured = make(map[int]data.Value, len(f.captured))
		for index, value := range f.captured {
			clone.captured[index] = scope.Bind(value)
		}
	}
	if len(f.capturedRefs) > 0 {
		clone.capturedRefs = make(map[int]*data.ZVal, len(f.capturedRefs))
		for index, slot := range f.capturedRefs {
			clone.capturedRefs[index] = scope.BindSlot(slot)
		}
	}
	return &clone
}

// NewLambdaExpression 创建一个新的Lambda表达式
func NewLambdaExpression(from data.From, params []data.GetValue, body []data.GetValue, vars []data.Variable, parent map[int]int, strict ...bool) *LambdaExpression {
	return &LambdaExpression{
		FunctionStatement: &FunctionStatement{
			Node:        NewNode(from),
			Params:      params,
			Body:        body,
			vars:        vars,
			IsGenerator: containsYield(body),
			StrictTypes: len(strict) != 0 && strict[0],
		},
		parent:         parent,
		parameterScope: parametersRequireClassScope(params),
	}
}

func parametersRequireClassScope(params []data.GetValue) bool {
	for _, parameter := range params {
		if typed, ok := parameter.(interface{ GetType() data.Types }); ok {
			if ref, ok := typed.GetType().(data.TypeRef); ok && ref.RequiresClassScope() {
				return true
			}
		}
	}
	return false
}

func (f *LambdaExpression) ParameterTypeContext(ctx data.Context) data.Context {
	if !f.parameterScope {
		return ctx
	}
	if owner, ok := f.ctx.(*data.ClassMethodContext); ok {
		if f.IsStatic {
			return data.NewStaticMethodContext(ctx, owner.SelfClass, owner.StaticClass)
		}
		return data.WrapMethodFrame(ctx, owner.ClassValue, owner.SelfClass, owner.StaticClass)
	}
	return ctx
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
				value = zv.ReadValue()
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
	// Retain lexical/object identity without retaining the caller's locals.
	// Captured values and independent reference buckets own their lifetimes.
	definition := ctx.CreateBaseContext()
	if owner, ok := ctx.(*data.ClassMethodContext); ok {
		if owner.ObjectValue == nil {
			definition = data.NewStaticMethodContext(definition, owner.SelfClass, owner.StaticClass)
		} else {
			definition = data.WrapMethodFrame(definition, owner.ClassValue, owner.SelfClass, owner.StaticClass)
		}
	}
	markContextEscaped(definition)
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
			zv.AddRefSlot()
			capturedRefs[cID] = data.CopyReferenceBucket(zv)
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
			Node:             f.Node,
			Params:           f.Params,
			Body:             f.Body,
			vars:             f.vars,
			IsGenerator:      f.IsGenerator,
			Name:             f.Name,
			Ret:              f.Ret,
			StrictTypes:      f.StrictTypes,
			ReturnsReference: f.ReturnsReference,
		},
		ctx:            definition,
		parent:         f.parent,
		captured:       captured,
		capturedRefs:   capturedRefs,
		IsStatic:       f.IsStatic,
		parameterScope: f.parameterScope,
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
	ctx.SetStrictTypes(f.StrictTypes)
	inner := ctx
	if bc, ok := ctx.(*data.BoundContext); ok {
		inner = bc.Context
	}
	execCtx := inner
	var borrowed *data.ClassMethodContext
	if defineClassCtx, ok := f.ctx.(*data.ClassMethodContext); ok && !f.IsStatic {
		borrowed = data.WrapMethodFrame(inner, defineClassCtx.ClassValue, defineClassCtx.SelfClass, defineClassCtx.StaticClass)
		execCtx = borrowed
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
			borrowed = cmc
		}
	}
	if borrowed != nil {
		defer borrowed.ReleaseBorrowedFrame()
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
					data.BindContextReference(execCtx, f.vars[cID].GetIndex(), zv)
					continue
				}
			}
			if f.ctx != nil {
				if parentZVal := f.ctx.GetIndexZVal(pID); parentZVal != nil {
					data.BindContextReference(execCtx, f.vars[cID].GetIndex(), parentZVal)
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

	frame := data.CallFrame{Function: f.Name}
	if frame.Function == "" {
		frame.Function = "{closure}"
	}
	if from := f.GetFrom(); from != nil {
		frame.File = from.GetSource()
		line, _ := from.GetStartPosition()
		frame.Line = line + 1
	}
	phpCallEnter(ctx, frame)
	defer phpCallLeave(ctx)

	var ctl data.Control
	for bodyIndex := 0; bodyIndex < len(f.Body); bodyIndex++ {
		statement := f.Body[bodyIndex]
		_, ctl = statement.GetValue(execCtx)
		if ctl != nil {
			switch rv := ctl.(type) {
			case data.ExitControl:
				return nil, ctl
			case data.ReturnControl:
				ret := rv.ReturnValue()
				if f.Ret == data.TypeVoid {
					return data.NewNullValue(), nil
				}
				if f.Ret != data.TypeInvalid {
					prepared, ok, conversion := data.PrepareDeclaredValueInContext(f.Ret, ret, execCtx)
					if conversion != nil {
						return nil, data.ReturnTypeError(f.from, fmt.Errorf("closure return type must be %s", f.Ret.String()), conversion)
					}
					if !ok {
						return nil, data.NewTypeError(f.from, fmt.Errorf("closure return type must be %s", f.Ret.String()))
					}
					return prepared, nil
				}
				return ret, nil
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
	if !(f.Ret == data.TypeInvalid || f.Ret == data.TypeVoid) {
		return nil, data.NewTypeError(f.from, fmt.Errorf("closure must return a value"))
	}
	// PHP：普通闭包没有 return 时返回 null。箭头函数 fn() => expr 由解析器包成 return。
	return data.NewNullValue(), nil
}
