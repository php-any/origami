package php

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewIteratorToArrayFunction 创建 iterator_to_array 函数
// PHP 语法：
// iterator_to_array(Traversable $iterator, bool $use_keys = true): array
// 将迭代器转换为数组，如果 $use_keys 为 true 则保留键名
func NewIteratorToArrayFunction() data.FuncStmt {
	return &IteratorToArrayFunction{}
}

type IteratorToArrayFunction struct{}

func (f *IteratorToArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 直接取参数，避免 ConvertFromIndex 在 ClassValue 实现 GetSource 时失败
	iterVal, _ := ctx.GetIndexValue(0)
	iterVal = unwrapThisValue(iterVal)

	// PHP：use_keys 默认 true。内部 callVMFunc 只传迭代器时，第 2 槽常是 NullValue；
	// NullValue.AsBool()==false 会把键丢掉，Collection::partition(ArrayIterator)
	// 就会把 class 变成 0= HTML 属性。
	useKeys := true
	if uk, ok := ctx.GetIndexValue(1); ok && uk != nil {
		if _, isNull := uk.(*data.NullValue); !isNull {
			if as, ok := uk.(data.AsBool); ok {
				if b, err := as.AsBool(); err == nil {
					useKeys = b
				}
			}
		}
	}

	// 优先：有 getIterator() 就先展开（IteratorAggregate）
	if classVal, ok := iterVal.(*data.ClassValue); ok {
		if method, exists := findInstanceMethod(classVal, "getIterator"); exists {
			fnCtx := classVal.CreateContext(method.GetVariables())
			result, ctl := method.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if v, ok := result.(data.Value); ok {
				iterVal = unwrapThisValue(v)
			} else {
				return nil, data.NewErrorThrowByName(nil, fmt.Errorf("getIterator 必须返回一个值"), "RuntimeException")
			}
		} else if checkInterface(ctx, "IteratorAggregate", classVal.Class) {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("IteratorAggregate 必须实现 getIterator 方法"), "RuntimeException")
		}
	}

	if classVal, ok := iterVal.(*data.ClassValue); ok {
		className := classVal.Class.GetName()
		if className == "Generator" || strings.HasSuffix(className, "\\Generator") ||
			checkInterface(ctx, "Iterator", classVal.Class) || findInstanceMethodExists(classVal, "valid") {
			return extractIteratorData(ctx, classVal, useKeys)
		}
	}

	if arrVal, ok := iterVal.(*data.ArrayValue); ok {
		return arrVal, nil
	}

	if objVal, ok := iterVal.(*data.ObjectValue); ok {
		if useKeys {
			return objVal, nil
		}
		values := make([]data.Value, 0)
		objVal.RangeProperties(func(key string, value data.Value) bool {
			values = append(values, value)
			return true
		})
		return data.NewArrayValue(values), nil
	}

	got := "nil"
	if iterVal != nil {
		got = fmt.Sprintf("%T", iterVal)
	}
	return nil, data.NewErrorThrowByName(nil, fmt.Errorf("参数必须是可迭代的对象或数组 (got %s)", got), "TypeError")
}

// unwrapThisValue 解开 fluent return 的 ThisValue，得到真实 ClassValue。
func unwrapThisValue(v data.Value) data.Value {
	if v == nil {
		return nil
	}
	if tv, ok := v.(*data.ThisValue); ok && tv != nil && tv.ClassValue != nil {
		return tv.ClassValue
	}
	return v
}

func findInstanceMethod(classVal *data.ClassValue, name string) (data.Method, bool) {
	if classVal == nil || classVal.Class == nil {
		return nil, false
	}
	if method, exists := classVal.GetMethod(name); exists {
		return method, true
	}
	current := classVal.Class
	for current != nil {
		if method, exists := current.GetMethod(name); exists {
			return method, true
		}
		ext := current.GetExtend()
		if ext == nil || *ext == "" {
			break
		}
		parent, ok := classVal.GetVM().GetClass(*ext)
		if !ok || parent == nil {
			break
		}
		current = parent
	}
	return nil, false
}

func findInstanceMethodExists(classVal *data.ClassValue, name string) bool {
	_, ok := findInstanceMethod(classVal, name)
	return ok
}

// extractIteratorData 从 Iterator 对象中提取数据
func extractIteratorData(ctx data.Context, classVal *data.ClassValue, useKeys bool) (data.GetValue, data.Control) {
	if method, exists := classVal.GetMethod("rewind"); exists {
		fnCtx := classVal.CreateContext(method.GetVariables())
		_, ctl := method.Call(fnCtx)
		if ctl != nil {
			return nil, ctl
		}
	}

	result := make([]*data.ZVal, 0)
	keyPositions := make(map[string]int)
	index := 0

	for {
		validMethod, validExists := classVal.GetMethod("valid")
		if !validExists {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Iterator 必须实现 valid 方法"), "RuntimeException")
		}

		validCtx := classVal.CreateContext(validMethod.GetVariables())
		validResult, ctl := validMethod.Call(validCtx)
		if ctl != nil {
			return nil, ctl
		}

		isValid := false
		if boolVal, ok := validResult.(data.AsBool); ok {
			if val, err := boolVal.AsBool(); err == nil {
				isValid = val
			}
		}

		if !isValid {
			break
		}

		keyMethod, keyExists := classVal.GetMethod("key")
		if !keyExists {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Iterator 必须实现 key 方法"), "RuntimeException")
		}

		keyCtx := classVal.CreateContext(keyMethod.GetVariables())
		keyResult, ctl := keyMethod.Call(keyCtx)
		if ctl != nil {
			return nil, ctl
		}

		var key string
		if useKeys {
			if strVal, ok := keyResult.(data.AsString); ok {
				key = strVal.AsString()
			} else if intVal, ok := keyResult.(data.AsInt); ok {
				if i, err := intVal.AsInt(); err == nil {
					key = fmt.Sprintf("%d", i)
				} else {
					key = fmt.Sprintf("%d", index)
				}
			} else {
				key = fmt.Sprintf("%d", index)
			}
		} else {
			key = fmt.Sprintf("%d", index)
		}

		currentMethod, currentExists := classVal.GetMethod("current")
		if !currentExists {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Iterator 必须实现 current 方法"), "RuntimeException")
		}

		currentCtx := classVal.CreateContext(currentMethod.GetVariables())
		currentResult, ctl := currentMethod.Call(currentCtx)
		if ctl != nil {
			return nil, ctl
		}

		value, ok := currentResult.(data.Value)
		if !ok {
			value = data.NewNullValue()
		}
		zv := data.NewZVal(value)
		if useKeys {
			zv.Name = key
			if position, exists := keyPositions[key]; exists {
				result[position] = zv
			} else {
				keyPositions[key] = len(result)
				result = append(result, zv)
			}
		} else {
			result = append(result, zv)
		}

		nextMethod, nextExists := classVal.GetMethod("next")
		if !nextExists {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Iterator 必须实现 next 方法"), "RuntimeException")
		}

		nextCtx := classVal.CreateContext(nextMethod.GetVariables())
		_, ctl = nextMethod.Call(nextCtx)
		if ctl != nil {
			return nil, ctl
		}

		index++
	}

	return &data.ArrayValue{List: result}, nil
}

// checkInterface 检查类是否实现了指定接口
func checkInterface(ctx data.Context, interfaceName string, classStmt data.ClassStmt) bool {
	normalize := func(name string) string {
		for len(name) > 0 && name[0] == '\\' {
			name = name[1:]
		}
		return name
	}
	want := normalize(interfaceName)

	matches := func(impl string) bool {
		n := normalize(impl)
		if n == want {
			return true
		}
		if i := strings.LastIndex(n, "\\"); i >= 0 && n[i+1:] == want {
			return true
		}
		return false
	}

	if implements := classStmt.GetImplements(); implements != nil {
		for _, impl := range implements {
			if matches(impl) {
				return true
			}
		}
	}

	vm := ctx.GetVM()
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok {
			break
		}

		if implements := next.GetImplements(); implements != nil {
			for _, impl := range implements {
				if matches(impl) {
					return true
				}
			}
		}

		last = next
	}

	return false
}

func (f *IteratorToArrayFunction) GetName() string {
	return "iterator_to_array"
}

var iteratorToArrayFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "iterator", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "use_keys", 1, data.NewBoolValue(true), data.Bool{}),
}

func (f *IteratorToArrayFunction) GetParams() []data.GetValue {
	return iteratorToArrayFunctionGetParams
}

var iteratorToArrayFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "iterator", 0, data.Mixed{}),
	node.NewVariable(nil, "use_keys", 1, data.Bool{}),
}

func (f *IteratorToArrayFunction) GetVariables() []data.Variable {
	return iteratorToArrayFunctionGetVariables
}
