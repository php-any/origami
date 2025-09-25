package utils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/php-any/origami/data"
)

// ConvertFromIndex 从上下文索引位置转换值到指定类型
func ConvertFromIndex[S any](ctx data.Context, index int) (S, error) {
	var result S
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return result, fmt.Errorf("参数索引 %d 不存在", index)
	}

	// 直接类型转换
	if converted, err := convertValue[S](v); err == nil {
		return converted, nil
	}

	// 类型别名特殊处理
	return convertTypeAlias[S](v)
}

// convertValue 通用值转换
func convertValue[S any](v data.Value) (S, error) {
	var result S
	switch val := v.(type) {
	case data.GetSource:
		if src := val.GetSource(); src != nil {
			if converted, ok := src.(S); ok {
				return converted, nil
			}
		}
		return result, fmt.Errorf("无法从 GetSource 转换到 %T", result)

	case *data.ClassValue:
		if p, ok := val.Class.(data.GetSource); ok {
			if src := p.GetSource(); src != nil {
				if converted, ok := src.(S); ok {
					return converted, nil
				}
			}
		}
		return result, fmt.Errorf("无法从 ClassValue 转换到 %T", result)

	case *data.AnyValue:
		if converted, ok := val.Value.(S); ok {
			return converted, nil
		}
		return result, fmt.Errorf("无法从 AnyValue 转换到 %T", result)

	case *data.IntValue:
		return convertFromIntValue[S](val)

	case *data.StringValue:
		return convertFromStringValue[S](val)

	case *data.FloatValue:
		return convertFromFloatValue[S](val)

	case *data.BoolValue:
		return convertFromBoolValue[S](val)

	case *data.ArrayValue:
		return convertFromArrayValue[S](val)

	default:
		return result, fmt.Errorf("不支持的值类型: %T", v)
	}
}

// convertFromIntValue 从 IntValue 转换
func convertFromIntValue[S any](val *data.IntValue) (S, error) {
	var result S
	intVal, err := val.AsInt()
	if err != nil {
		return result, err
	}

	// 先尝试直接类型断言（性能优化）
	if s, ok := any(intVal).(S); ok {
		return s, nil
	}

	// 尝试转换到其他基本类型
	switch any(result).(type) {
	case int8:
		return any(int8(intVal)).(S), nil
	case int16:
		return any(int16(intVal)).(S), nil
	case int32:
		return any(int32(intVal)).(S), nil
	case int64:
		return any(int64(intVal)).(S), nil
	case uint:
		return any(uint(intVal)).(S), nil
	case uint8:
		return any(uint8(intVal)).(S), nil
	case uint16:
		return any(uint16(intVal)).(S), nil
	case uint32:
		return any(uint32(intVal)).(S), nil
	case uint64:
		return any(uint64(intVal)).(S), nil
	case float32:
		return any(float32(intVal)).(S), nil
	case float64:
		return any(float64(intVal)).(S), nil
	case string:
		return any(fmt.Sprintf("%d", intVal)).(S), nil
	case bool:
		return any(intVal != 0).(S), nil
	}

	// 如果直接类型断言失败，使用反射处理复杂类型
	return convertTypeAlias[S](val)
}

// convertFromStringValue 从 StringValue 转换
func convertFromStringValue[S any](val *data.StringValue) (S, error) {
	var result S
	strVal := val.AsString()

	if s, ok := any(strVal).(S); ok {
		return s, nil
	}

	// 先尝试直接类型断言（性能优化）
	switch any(result).(type) {
	case string:
		return any(strVal).(S), nil
	case bool:
		// 尝试解析为布尔值
		if boolVal, err := parseBool(strVal); err == nil {
			return any(boolVal).(S), nil
		}
		return result, fmt.Errorf("无法将字符串 '%s' 转换为布尔类型", strVal)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return result, fmt.Errorf("字符串到数值转换需要两步转换，请使用模板生成的两步转换")
	}

	// 如果直接类型断言失败，使用反射处理复杂类型
	return convertTypeAlias[S](val)
}

// convertFromFloatValue 从 FloatValue 转换
func convertFromFloatValue[S any](val *data.FloatValue) (S, error) {
	var result S
	floatVal, err := val.AsFloat()
	if err != nil {
		return result, err
	}

	// 先尝试直接类型断言（性能优化）
	if s, ok := any(floatVal).(S); ok {
		return s, nil
	}

	// 尝试转换到其他基本类型
	switch any(result).(type) {
	case float32:
		return any(float32(floatVal)).(S), nil
	case int:
		return any(int(floatVal)).(S), nil
	case int8:
		return any(int8(floatVal)).(S), nil
	case int16:
		return any(int16(floatVal)).(S), nil
	case int32:
		return any(int32(floatVal)).(S), nil
	case int64:
		return any(int64(floatVal)).(S), nil
	case uint:
		return any(uint(floatVal)).(S), nil
	case uint8:
		return any(uint8(floatVal)).(S), nil
	case uint16:
		return any(uint16(floatVal)).(S), nil
	case uint32:
		return any(uint32(floatVal)).(S), nil
	case uint64:
		return any(uint64(floatVal)).(S), nil
	case string:
		return any(fmt.Sprintf("%g", floatVal)).(S), nil
	case bool:
		return any(floatVal != 0).(S), nil
	}

	// 如果直接类型断言失败，使用反射处理复杂类型
	return convertTypeAlias[S](val)
}

// convertFromBoolValue 从 BoolValue 转换
func convertFromBoolValue[S any](val *data.BoolValue) (S, error) {
	var result S
	boolVal, err := val.AsBool()
	if err != nil {
		return result, err
	}

	// 先尝试直接类型断言（性能优化）
	if s, ok := any(boolVal).(S); ok {
		return s, nil
	}

	// 尝试转换到其他基本类型
	switch any(result).(type) {
	case int:
		if boolVal {
			return any(1).(S), nil
		}
		return any(0).(S), nil
	case int8:
		if boolVal {
			return any(int8(1)).(S), nil
		}
		return any(int8(0)).(S), nil
	case int16:
		if boolVal {
			return any(int16(1)).(S), nil
		}
		return any(int16(0)).(S), nil
	case int32:
		if boolVal {
			return any(int32(1)).(S), nil
		}
		return any(int32(0)).(S), nil
	case int64:
		if boolVal {
			return any(int64(1)).(S), nil
		}
		return any(int64(0)).(S), nil
	case uint:
		if boolVal {
			return any(uint(1)).(S), nil
		}
		return any(uint(0)).(S), nil
	case uint8:
		if boolVal {
			return any(uint8(1)).(S), nil
		}
		return any(uint8(0)).(S), nil
	case uint16:
		if boolVal {
			return any(uint16(1)).(S), nil
		}
		return any(uint16(0)).(S), nil
	case uint32:
		if boolVal {
			return any(uint32(1)).(S), nil
		}
		return any(uint32(0)).(S), nil
	case uint64:
		if boolVal {
			return any(uint64(1)).(S), nil
		}
		return any(uint64(0)).(S), nil
	case float32:
		if boolVal {
			return any(float32(1.0)).(S), nil
		}
		return any(float32(0.0)).(S), nil
	case float64:
		if boolVal {
			return any(1.0).(S), nil
		}
		return any(0.0).(S), nil
	case string:
		if boolVal {
			return any("true").(S), nil
		}
		return any("false").(S), nil
	}

	// 如果直接类型断言失败，使用反射处理复杂类型
	return convertTypeAlias[S](val)
}

// convertFromArrayValue 从 ArrayValue 转换
func convertFromArrayValue[S any](val *data.ArrayValue) (S, error) {
	var result S

	// 先尝试直接类型断言（性能优化）
	if s, ok := any(val.Value).(S); ok {
		return s, nil
	}

	// 尝试转换到基本类型切片
	switch any(result).(type) {
	case []int:
		slice := make([]int, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[int](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []int8:
		slice := make([]int8, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[int8](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []int16:
		slice := make([]int16, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[int16](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []int32:
		slice := make([]int32, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[int32](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []int64:
		slice := make([]int64, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[int64](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []uint:
		slice := make([]uint, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[uint](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []uint8:
		slice := make([]uint8, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[uint8](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []uint16:
		slice := make([]uint16, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[uint16](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []uint32:
		slice := make([]uint32, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[uint32](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []uint64:
		slice := make([]uint64, 0, len(val.Value))
		for _, item := range val.Value {
			if intVal, err := convertValue[uint64](item); err == nil {
				slice = append(slice, intVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []float32:
		slice := make([]float32, 0, len(val.Value))
		for _, item := range val.Value {
			if floatVal, err := convertValue[float32](item); err == nil {
				slice = append(slice, floatVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []float64:
		slice := make([]float64, 0, len(val.Value))
		for _, item := range val.Value {
			if floatVal, err := convertValue[float64](item); err == nil {
				slice = append(slice, floatVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []string:
		slice := make([]string, 0, len(val.Value))
		for _, item := range val.Value {
			if strVal, err := convertValue[string](item); err == nil {
				slice = append(slice, strVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	case []bool:
		slice := make([]bool, 0, len(val.Value))
		for _, item := range val.Value {
			if boolVal, err := convertValue[bool](item); err == nil {
				slice = append(slice, boolVal)
			} else {
				return result, fmt.Errorf("转换数组元素失败: %w", err)
			}
		}
		return any(slice).(S), nil
	}

	// 如果直接类型断言失败，使用反射处理复杂类型
	targetType := reflect.TypeOf((*S)(nil)).Elem()

	// 只支持切片类型
	if targetType.Kind() != reflect.Slice {
		return result, fmt.Errorf("无法将数组转换为非切片类型 %T", result)
	}

	// 创建目标切片
	slice := reflect.MakeSlice(targetType, len(val.Value), len(val.Value))

	for i, item := range val.Value {
		// 递归转换每个元素
		if converted, err := convertValue[any](item); err == nil {
			slice.Index(i).Set(reflect.ValueOf(converted))
		} else {
			return result, fmt.Errorf("无法转换数组元素 %d: %v", i, err)
		}
	}

	return slice.Interface().(S), nil
}

// convertTypeAlias 处理类型别名转换
func convertTypeAlias[S any](v data.Value) (S, error) {
	var result S

	// 先尝试直接类型断言（性能优化）
	if s, ok := any(v).(S); ok {
		return s, nil
	}

	// 特殊处理 time.Duration
	if _, ok := any(result).(time.Duration); ok {
		if intVal, ok := v.(data.AsInt); ok {
			if duration, err := intVal.AsInt(); err == nil {
				return any(time.Duration(duration)).(S), nil
			}
		}
		return result, fmt.Errorf("无法将 %T 转换为 time.Duration", v)
	}

	// 如果直接类型断言失败，使用反射处理复杂类型
	targetType := reflect.TypeOf((*S)(nil)).Elem()

	// 检查是否为具名类型
	if targetType.PkgPath() != "" && targetType.Name() != "" {
		// 特殊处理 time.Duration
		if targetType.PkgPath() == "time" && targetType.Name() == "Duration" {
			if intVal, ok := v.(data.AsInt); ok {
				if duration, err := intVal.AsInt(); err == nil {
					return any(time.Duration(duration)).(S), nil
				}
			}
			return result, fmt.Errorf("无法将 %T 转换为 time.Duration", v)
		}

		// 其他具名类型尝试直接转换
		if converted, ok := v.(S); ok {
			return converted, nil
		}
	}

	return result, fmt.Errorf("无法转换类型 %T 到 %T", v, result)
}

// 辅助函数
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%g", &result)
	return result, err
}

func parseBool(s string) (bool, error) {
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("无效的布尔值: %s", s)
	}
}

// Convert 将 data.Value 转换为指定类型
func Convert[S any](v data.Value) (S, error) {
	// 直接使用 convertValue 函数，与 ConvertFromIndex 保持一致
	return convertValue[S](v)
}
