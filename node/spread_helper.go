package node

import (
	"strings"

	"github.com/php-any/origami/data"
)

// spreadToValues 将 spread 展开值统一转换为平铺的值列表。
// 支持 ArrayValue、ObjectValue 和 Generator（生成器），
// 与 PHP 语义一致：...$generator 会遍历生成器所有 yield 值。
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
		// Generator 展开：遍历生成器的所有 yield 值
		if v.Class != nil && isGeneratorClassName(v.Class.GetName()) {
			return iterateGenerator(ctx, v)
		}
	}

	if val, ok := spreadVal.(data.Value); ok {
		return []data.Value{val}, nil
	}
	return nil, nil
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
		// Generator 展开：遍历生成器的所有 yield 值
		if v.Class != nil && isGeneratorClassName(v.Class.GetName()) {
			vals, ctl := iterateGenerator(ctx, v)
			if ctl != nil {
				return nil, false
			}
			result := make([]data.GetValue, 0, len(vals))
			for _, val := range vals {
				result = append(result, val)
			}
			return result, true
		}
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
