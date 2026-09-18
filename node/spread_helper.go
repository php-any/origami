package node

import (
	"strings"

	"github.com/php-any/origami/data"
)

// spreadToValues 将 spread 展开值统一转换为平铺的值列表。
// 支持 ArrayValue、ObjectValue、Generator，以及 Iterator / IteratorAggregate
// （PHP：...$traversable 可展开；Filament Widget 配置常展开 Collection/属性数组）。
func spreadToValues(ctx data.Context, spreadVal data.GetValue) ([]data.Value, data.Control) {
	switch v := spreadVal.(type) {
	case *data.ArrayValue:
		result := make([]data.Value, 0, len(v.List))
		for _, z := range v.List {
			if z != nil {
				result = append(result, z.Value)
			}
		}
		return result, nil

	case *data.ObjectValue:
		result := make([]data.Value, 0)
		v.RangeProperties(func(_ string, val data.Value) bool {
			result = append(result, val)
			return true
		})
		return result, nil

	case *data.ClassValue:
		return iterateClassForSpread(ctx, v)
	}

	if val, ok := spreadVal.(data.Value); ok {
		return []data.Value{val}, nil
	}
	return nil, nil
}

// iterateClassForSpread 展开 Generator / Iterator / IteratorAggregate。
// 不可遍历时返回 (nil, nil)，由调用方决定是否报错。
func iterateClassForSpread(ctx data.Context, v *data.ClassValue) ([]data.Value, data.Control) {
	if v == nil || v.Class == nil {
		return nil, nil
	}
	if isGeneratorClassName(v.Class.GetName()) {
		return iterateGenerator(ctx, v)
	}

	isIterator, ctl := checkClassIs(ctx, v.Class, "Iterator")
	if ctl != nil {
		return nil, ctl
	}
	if isIterator {
		return iterateIteratorMethods(ctx, v)
	}

	isAggregate, ctl := checkClassIs(ctx, v.Class, "IteratorAggregate")
	if ctl != nil {
		return nil, ctl
	}
	if isAggregate {
		inner, ictl := callValueMethod(v, "getIterator")
		if ictl != nil {
			return nil, ictl
		}
		switch it := inner.(type) {
		case *data.ClassValue:
			return iterateClassForSpread(ctx, it)
		case *data.ThisValue:
			if it.ClassValue != nil {
				return iterateClassForSpread(ctx, it.ClassValue)
			}
		case *data.ArrayValue:
			return spreadToValues(ctx, it)
		}
	}

	// Traversable 标记接口：部分类只声明 implements Traversable
	isTrav, ctl := checkClassIs(ctx, v.Class, "Traversable")
	if ctl != nil {
		return nil, ctl
	}
	if isTrav {
		// 无 getIterator/Iterator 方法时无法展开
		return nil, nil
	}
	return nil, nil
}

func iterateIteratorMethods(ctx data.Context, obj *data.ClassValue) ([]data.Value, data.Control) {
	result := make([]data.Value, 0)
	if ctl := callVoidMethod(obj, "rewind"); ctl != nil {
		return nil, ctl
	}
	for {
		valid, ctl := callBoolMethod(obj, "valid")
		if ctl != nil {
			return nil, ctl
		}
		if !valid {
			break
		}
		val, ctl := callValueMethod(obj, "current")
		if ctl != nil {
			return nil, ctl
		}
		if val != nil {
			result = append(result, val)
		}
		if ctl := callVoidMethod(obj, "next"); ctl != nil {
			return nil, ctl
		}
	}
	return result, nil
}

// spreadToValuesForNew 将 spread 展开值转换为平铺的值列表，用于 new 构造场景。
// 与 spreadToValues 的区别：ArrayValue 的关联字符串键转换为 NamedArgument。
func spreadToValuesForNew(ctx data.Context, spreadVal data.GetValue) ([]data.GetValue, bool) {
	switch v := spreadVal.(type) {
	case *data.ArrayValue:
		result := make([]data.GetValue, 0, len(v.List))
		for _, z := range v.List {
			if z == nil {
				continue
			}
			// 关联字符串键展开为命名实参
			if z.Name != "" {
				if _, isInt := data.ParseIntArrayKeyName(z.Name); !isInt {
					result = append(result, NewNamedArgument(nil, z.Name, z.Value))
					continue
				}
			}
			result = append(result, z.Value)
		}
		return result, true

	case *data.ObjectValue:
		result := make([]data.GetValue, 0)
		v.RangeProperties(func(key string, val data.Value) bool {
			result = append(result, NewNamedArgument(nil, key, val))
			return true
		})
		return result, true

	case *data.ClassValue:
		vals, ctl := iterateClassForSpread(ctx, v)
		if ctl != nil || vals == nil {
			return nil, false
		}
		result := make([]data.GetValue, 0, len(vals))
		for _, val := range vals {
			result = append(result, val)
		}
		return result, true
	}

	if val, ok := spreadVal.(data.Value); ok {
		return []data.GetValue{val}, true
	}
	return nil, false
}

// isGeneratorClassName 判断类名是否为 Generator
func isGeneratorClassName(name string) bool {
	return name == "Generator" || strings.HasSuffix(name, "\\Generator")
}

// iterateGenerator 遍历 Generator 的所有 yield 值。
// 通过调用 valid/current/next/rewind 迭代器方法实现。
func iterateGenerator(ctx data.Context, gen *data.ClassValue) ([]data.Value, data.Control) {
	result := make([]data.Value, 0)
	genCtx := gen.CreateContext(nil)

	// rewind() 重置到第一个元素
	if m, ok := gen.GetMethod("rewind"); ok {
		if _, ctl := m.Call(genCtx); ctl != nil {
			return nil, ctl
		}
	}

	for {
		// valid() 检查是否有更多元素
		validMethod, ok := gen.GetMethod("valid")
		if !ok {
			break
		}
		validVal, ctl := validMethod.Call(genCtx)
		if ctl != nil {
			return nil, ctl
		}
		validBool, _ := validVal.(data.Value)
		if validBool == nil || validBool.AsString() != "true" {
			break
		}

		// current() 获取当前元素
		currentMethod, ok := gen.GetMethod("current")
		if !ok {
			break
		}
		currentVal, ctl := currentMethod.Call(genCtx)
		if ctl != nil {
			return nil, ctl
		}
		if val, ok := currentVal.(data.Value); ok {
			result = append(result, val)
		}

		// next() 前进到下一个元素
		nextMethod, ok := gen.GetMethod("next")
		if !ok {
			break
		}
		if _, ctl := nextMethod.Call(genCtx); ctl != nil {
			return nil, ctl
		}
	}

	return result, nil
}
