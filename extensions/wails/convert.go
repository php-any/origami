package wails

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

// optionsFromConstructor 读取 __construct 的首个 array-like 实参（ArrayValue 或 ObjectValue）。
func optionsFromConstructor(ctx data.Context) (data.Value, bool) {
	if v, ok := ctx.GetIndexValue(0); ok && isOptionsContainer(v) {
		return v, true
	}
	for _, arg := range ctx.GetCallArgs() {
		if arg == nil {
			continue
		}
		gv, ctl := arg.GetValue(ctx)
		if ctl == nil && gv != nil {
			if val, ok := gv.(data.Value); ok && isOptionsContainer(val) {
				return val, true
			}
		}
		break
	}
	return nil, false
}

func isOptionsContainer(v data.Value) bool {
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return v != nil
	default:
		return false
	}
}

// optionGet 从 array / 关联数组（ObjectValue）中按键取值。
func optionGet(opts data.Value, key string) (data.Value, bool) {
	switch o := opts.(type) {
	case *data.ArrayValue:
		return arrayGet(o, key)
	case *data.ObjectValue:
		val, ctl := o.GetProperty(key)
		if ctl == nil && val != nil {
			if _, isNull := val.(*data.NullValue); !isNull {
				return val, true
			}
		}
	}
	return nil, false
}

func applyOptionsMap(cv *data.ClassValue, opts data.Value, keys []string) {
	for _, key := range keys {
		if v, ok := optionGet(opts, key); ok {
			cv.SetProperty(key, v)
		}
	}
}

// firstArgString 读取当前调用上下文的第一个实参字符串。
func firstArgString(ctx data.Context) string {
	return toString(argValueAt(ctx, 0))
}

// argValueAt 读取调用实参（兼容 ClassMethodContext 与 GetCallArgs）。
func argValueAt(ctx data.Context, index int) data.Value {
	if v, ok := ctx.GetIndexValue(index); ok && v != nil {
		return v
	}
	args := ctx.GetCallArgs()
	if index < len(args) && args[index] != nil {
		if gv, ctl := args[index].GetValue(ctx); ctl == nil && gv != nil {
			if val, ok := gv.(data.Value); ok {
				return val
			}
		}
	}
	return nil
}

// modifierPrefixFromValue 把修饰键数组转为 "shift+cmdorctrl+" 形式前缀。
func modifierPrefixFromValue(v data.Value) string {
	if v == nil {
		return ""
	}
	var mods []string
	collect := func(val data.Value) {
		mod := strings.ToLower(strings.TrimSpace(toString(val)))
		if mod != "" {
			mods = append(mods, mod)
		}
	}
	switch o := v.(type) {
	case *data.ArrayValue:
		for _, z := range o.List {
			if z != nil {
				collect(z.Value)
			}
		}
	case *data.ObjectValue:
		o.RangeProperties(func(_ string, val data.Value) bool {
			collect(val)
			return true
		})
	default:
		collect(v)
	}
	if len(mods) == 0 {
		return ""
	}
	return strings.Join(mods, "+") + "+"
}

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
