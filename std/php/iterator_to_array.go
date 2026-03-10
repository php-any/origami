package php

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
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
	// 读取第一个参数 $iterator
	iterVal, _ := utils.ConvertFromIndex[data.Value](ctx, 0)

	// 读取第二个参数 $use_keys (默认 true)
	useKeys, _ := utils.ConvertFromIndex[bool](ctx, 1)

	// 检查是否实现了 IteratorAggregate 接口
	if classVal, ok := iterVal.(*data.ClassValue); ok {
		if checkInterface(ctx, "IteratorAggregate", classVal.Class) {
			// 调用 getIterator() 方法获取迭代器
			method, exists := classVal.GetMethod("getIterator")
			if !exists {
				return nil, data.NewErrorThrowByName(nil, fmt.Errorf("IteratorAggregate 必须实现 getIterator 方法"), "RuntimeException")
			}

			// 创建调用上下文
			fnCtx := classVal.CreateContext(method.GetVariables())
			result, ctl := method.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}

			// 更新 iterVal 为 getIterator() 的返回值
			if v, ok := result.(data.Value); ok {
				iterVal = v
			} else {
				return nil, data.NewErrorThrowByName(nil, fmt.Errorf("getIterator 必须返回一个值"), "RuntimeException")
			}
		}
	}

	// 检查是否实现了 Iterator 接口（包括 Generator 类）
	if classVal, ok := iterVal.(*data.ClassValue); ok {
		if classVal.Class.GetName() == "Generator" || checkInterface(ctx, "Iterator", classVal.Class) {
			// 从 Iterator 对象中提取数据
			return extractIteratorData(ctx, classVal, useKeys)
		}
	}

	// 如果已经是数组，直接返回
	if arrVal, ok := iterVal.(*data.ArrayValue); ok {
		return arrVal, nil
	}

	// 如果是 ObjectValue（关联数组），也返回
	if objVal, ok := iterVal.(*data.ObjectValue); ok {
		if useKeys {
			return objVal, nil
		}
		// 不使用键名，只返回值
		values := make([]data.Value, 0)
		objVal.RangeProperties(func(key string, value data.Value) bool {
			values = append(values, value)
			return true
		})
		return data.NewArrayValue(values), nil
	}

	return nil, data.NewErrorThrowByName(nil, fmt.Errorf("参数必须是可迭代的对象或数组"), "TypeError")
}

// extractIteratorData 从 Iterator 对象中提取数据
func extractIteratorData(ctx data.Context, classVal *data.ClassValue, useKeys bool) (data.GetValue, data.Control) {
	// 调用 rewind() 方法
	if method, exists := classVal.GetMethod("rewind"); exists {
		fnCtx := classVal.CreateContext(method.GetVariables())
		_, ctl := method.Call(fnCtx)
		if ctl != nil {
			return nil, ctl
		}
	}

	result := make(map[string]data.Value)
	index := 0

	// 循环遍历迭代器
	for {
		// 调用 valid() 检查当前位置是否有效
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

		// 调用 key() 获取当前键
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

		// 调用 current() 获取当前值
		currentMethod, currentExists := classVal.GetMethod("current")
		if !currentExists {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Iterator 必须实现 current 方法"), "RuntimeException")
		}

		currentCtx := classVal.CreateContext(currentMethod.GetVariables())
		currentResult, ctl := currentMethod.Call(currentCtx)
		if ctl != nil {
			return nil, ctl
		}

		if val, ok := currentResult.(data.Value); ok {
			result[key] = val
		} else {
			result[key] = data.NewNullValue()
		}

		// 调用 next() 移动到下一个位置
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

	// 构建结果数组
	// PHP 语义：iterator_to_array 总是返回 array 类型
	// 即使 use_keys=true，也返回关联数组（ArrayValue），而不是 ObjectValue
	if useKeys {
		// 使用键名，创建关联数组
		objVal := data.NewObjectValue()
		for k, v := range result {
			objVal.SetProperty(k, v)
		}
		// 将 ObjectValue 包装成 ArrayValue 的第一个元素？不对！
		// 应该直接将 ObjectValue 作为数组返回，但需要是 ArrayValue 类型
		// 实际上 PHP 的 array 可以有关联键，我们用 ObjectValue 来表示关联数组
		return objVal, nil
	} else {
		// 不使用键名，返回索引数组
		values := make([]data.Value, 0, len(result))
		for _, v := range result {
			values = append(values, v)
		}
		return data.NewArrayValue(values), nil
	}
}

// checkInterface 检查类是否实现了指定接口
func checkInterface(ctx data.Context, interfaceName string, classStmt data.ClassStmt) bool {
	// 检查直接实现的接口
	if implements := classStmt.GetImplements(); implements != nil {
		for _, impl := range implements {
			if impl == interfaceName {
				return true
			}
		}
	}

	// 检查继承链
	vm := ctx.GetVM()
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok {
			break
		}

		// 检查父类实现的接口
		if implements := next.GetImplements(); implements != nil {
			for _, impl := range implements {
				if impl == interfaceName {
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

func (f *IteratorToArrayFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "use_keys", 1, data.NewBoolValue(true), data.Bool{}),
	}
}

func (f *IteratorToArrayFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.Mixed{}),
		node.NewVariable(nil, "use_keys", 1, data.Bool{}),
	}
}
