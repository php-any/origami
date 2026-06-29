package wails

import (
	"strconv"

	"github.com/php-any/origami/data"
)

// getThis 从 ClassMethodContext 中获取当前实例的 *data.ClassValue
func getThis(ctx data.Context) *data.ClassValue {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		return cv.ClassValue
	}
	return nil
}

// arrayGet 从 ArrayValue 中按字符串键查找值
func arrayGet(av *data.ArrayValue, key string) (data.Value, bool) {
	if av == nil {
		return nil, false
	}
	for _, z := range av.List {
		if z != nil && z.Name == key {
			return z.Value, true
		}
	}
	return nil, false
}

// toInt 从 Value 中提取整数
func toInt(v data.Value) int {
	if v == nil {
		return 0
	}
	if iv, ok := v.(data.AsInt); ok {
		n, _ := iv.AsInt()
		return n
	}
	return 0
}

// toFloat 从 Value 中提取浮点数
func toFloat(v data.Value) float64 {
	if v == nil {
		return 0
	}
	if fv, ok := v.(data.AsFloat); ok {
		f, _ := fv.AsFloat()
		return f
	}
	if iv, ok := v.(data.AsInt); ok {
		n, _ := iv.AsInt()
		return float64(n)
	}
	return 0
}

// toBool 从 Value 中提取布尔值
func toBool(v data.Value) bool {
	if v == nil {
		return false
	}
	if bv, ok := v.(data.AsBool); ok {
		b, _ := bv.AsBool()
		return b
	}
	return false
}

// toString 从 Value 中提取字符串
func toString(v data.Value) string {
	if v == nil {
		return ""
	}
	if sv, ok := v.(data.AsString); ok {
		return sv.AsString()
	}
	return ""
}

// toIntOrDefault 从 Value 提取整数，带默认值
func toIntOrDefault(v data.Value, defaultVal int) int {
	if v == nil {
		return defaultVal
	}
	if iv, ok := v.(data.AsInt); ok {
		n, _ := iv.AsInt()
		return n
	}
	return defaultVal
}

// toStringOrDefault 从 Value 提取字符串，带默认值
func toStringOrDefault(v data.Value, defaultVal string) string {
	if v == nil {
		return defaultVal
	}
	if sv, ok := v.(data.AsString); ok {
		return sv.AsString()
	}
	return defaultVal
}

// toBoolOrDefault 从 Value 提取布尔值，带默认值
func toBoolOrDefault(v data.Value, defaultVal bool) bool {
	if v == nil {
		return defaultVal
	}
	if bv, ok := v.(data.AsBool); ok {
		b, _ := bv.AsBool()
		return b
	}
	return defaultVal
}

// setDefaultIntProperty 在 ClassValue 上设置默认整数属性
func setDefaultIntProperty(cv *data.ClassValue, name string, val int) {
	if cv != nil {
		cv.SetProperty(name, data.NewIntValue(val))
	}
}

// setDefaultStringProperty 在 ClassValue 上设置默认字符串属性
func setDefaultStringProperty(cv *data.ClassValue, name string, val string) {
	if cv != nil {
		cv.SetProperty(name, data.NewStringValue(val))
	}
}

// setDefaultBoolProperty 在 ClassValue 上设置默认布尔属性
func setDefaultBoolProperty(cv *data.ClassValue, name string, val bool) {
	if cv != nil {
		cv.SetProperty(name, data.NewBoolValue(val))
	}
}

// setDefaultFloatProperty 在 ClassValue 上设置默认浮点属性
func setDefaultFloatProperty(cv *data.ClassValue, name string, val float64) {
	if cv != nil {
		cv.SetProperty(name, data.NewFloatValue(val))
	}
}

// applyArrayToClassValue 将 ArrayValue 中的键值对应用到 ClassValue 的属性
func applyArrayToClassValue(cv *data.ClassValue, av *data.ArrayValue, keys []string) {
	if cv == nil || av == nil {
		return
	}
	for _, key := range keys {
		if v, ok := arrayGet(av, key); ok {
			cv.SetProperty(key, v)
		}
	}
}

// intArrayKeyName 将整数转换为数组键名字符串
func intArrayKeyName(i int) string {
	return strconv.Itoa(i)
}
