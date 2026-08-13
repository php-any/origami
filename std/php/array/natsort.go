package array

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NatsortFunction 实现 natsort 函数
// natsort(array &$array): bool
// 使用"自然排序"算法对数组进行排序，保留键值对应关系
type NatsortFunction struct{}

func NewNatsortFunction() data.FuncStmt {
	return &NatsortFunction{}
}

func (f *NatsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if len(arrayRef.List) == 0 {
		return data.NewBoolValue(true), nil
	}

	// 使用自然排序，保留键值对应关系
	sort.SliceStable(arrayRef.List, func(i, j int) bool {
		return naturalCompare(arrayRef.List[i].Value.AsString(), arrayRef.List[j].Value.AsString(), false)
	})

	return data.NewBoolValue(true), nil
}

// NatcasesortFunction 实现 natcasesort 函数
// natcasesort(array &$array): bool
// 使用不区分大小写的"自然排序"算法对数组进行排序
type NatcasesortFunction struct{}

func NewNatcasesortFunction() data.FuncStmt {
	return &NatcasesortFunction{}
}

func (f *NatcasesortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if len(arrayRef.List) == 0 {
		return data.NewBoolValue(true), nil
	}

	sort.SliceStable(arrayRef.List, func(i, j int) bool {
		return naturalCompare(arrayRef.List[i].Value.AsString(), arrayRef.List[j].Value.AsString(), true)
	})

	return data.NewBoolValue(true), nil
}

// naturalCompare 实现自然排序比较（类似 PHP strnatcmp）
func naturalCompare(a, b string, caseInsensitive bool) bool {
	if caseInsensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}

	// 将字符串分割成数字和非数字片段
	re := regexp.MustCompile(`\d+|\D+`)
	aParts := re.FindAllString(a, -1)
	bParts := re.FindAllString(b, -1)

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		ap := aParts[i]
		bp := bParts[i]

		// 判断是否为数字
		aNum, aErr := strconv.ParseFloat(ap, 64)
		bNum, bErr := strconv.ParseFloat(bp, 64)

		if aErr == nil && bErr == nil {
			// 都是数字，按数值比较
			if aNum != bNum {
				return aNum < bNum
			}
		} else {
			// 至少有一个不是数字，按字符串比较
			if ap != bp {
				return ap < bp
			}
		}
	}

	// 前缀相同，较短的在前面
	return len(aParts) < len(bParts)
}

func (f *NatsortFunction) GetName() string {
	return "natsort"
}

func (f *NatsortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	}
}

func (f *NatsortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

func (f *NatcasesortFunction) GetName() string {
	return "natcasesort"
}

func (f *NatcasesortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	}
}

func (f *NatcasesortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

// ArrayIntersectUkeyFunction 实现 array_intersect_ukey 函数
// array_intersect_ukey(array $array, array ...$arrays, callable $key_compare_func): array
// 使用回调函数比较键名，返回第一个数组中与其它数组都存在的键对应的值
type ArrayIntersectUkeyFunction struct{}

func NewArrayIntersectUkeyFunction() data.FuncStmt {
	return &ArrayIntersectUkeyFunction{}
}

func (f *ArrayIntersectUkeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取所有参数（通过 Parameters 打包在 index 0）
	var arrays []*data.ArrayValue
	var compareFunc data.GetValue

	paramsVal, _ := ctx.GetIndexValue(0)
	if paramsVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	if paramsArray, ok := paramsVal.(*data.ArrayValue); ok {
		for _, item := range paramsArray.List {
			if arr, isArr := item.Value.(*data.ArrayValue); isArr {
				arrays = append(arrays, arr)
			} else if obj, isObj := item.Value.(*data.ObjectValue); isObj {
				// 将 ObjectValue 转换为 ArrayValue（按字符串键）
				objArr := objectToArray(obj)
				arrays = append(arrays, objArr)
			} else if compareFunc == nil {
				compareFunc = item.Value
			}
		}
	} else {
		// 不是打包参数模式，直接逐个读取
		for i := 0; ; i++ {
			v, ok := ctx.GetIndexValue(i)
			if !ok || v == nil {
				break
			}
			if arr, isArr := v.(*data.ArrayValue); isArr {
				arrays = append(arrays, arr)
			} else if obj, isObj := v.(*data.ObjectValue); isObj {
				objArr := objectToArray(obj)
				arrays = append(arrays, objArr)
			} else if compareFunc == nil {
				compareFunc = v
			}
		}
	}

	if len(arrays) < 2 || compareFunc == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	first := arrays[0]
	resultList := make([]*data.ZVal, 0)

	// 检查第一个数组的每个键
	for _, item := range first.List {
		key := item.Name
		foundInAll := true

		for _, arr := range arrays[1:] {
			keyExists := false
			for _, otherItem := range arr.List {
				otherKey := otherItem.Name
				// 使用回调函数比较键
				callResult, acl := callCompareFunc(ctx, compareFunc, key, otherKey)
				if acl != nil {
					return nil, acl
				}
				if callResult == 0 {
					keyExists = true
					break
				}
			}
			if !keyExists {
				foundInAll = false
				break
			}
		}

		if foundInAll {
			resultList = append(resultList, data.NewNamedZVal(key, item.Value))
		}
	}

	return &data.ArrayValue{List: resultList}, nil
}

// objectToArray 将 ObjectValue 转换为 ArrayValue（保留字符串键）
func objectToArray(obj *data.ObjectValue) *data.ArrayValue {
	list := make([]*data.ZVal, 0)
	obj.RangeProperties(func(key string, value data.Value) bool {
		list = append(list, data.NewNamedZVal(key, value))
		return true
	})
	return &data.ArrayValue{List: list}
}

// callCompareFunc 调用比较回调函数
func callCompareFunc(ctx data.Context, fn data.GetValue, a, b string) (int, data.Control) {
	// 创建两个参数的上下文
	vars := []data.Variable{
		node.NewVariable(nil, "a", 0, nil),
		node.NewVariable(nil, "b", 1, nil),
	}

	// 如果是闭包或函数值
	if closure, ok := fn.(*data.FuncValue); ok {
		callCtx := ctx.CreateContext(vars)
		callCtx.SetIndexZVal(0, data.NewZVal(data.NewStringValue(a)))
		callCtx.SetIndexZVal(1, data.NewZVal(data.NewStringValue(b)))
		v, acl := closure.Call(callCtx)
		if acl != nil {
			return 0, acl
		}
		if iv, ok := v.(data.AsInt); ok {
			val, err := iv.AsInt()
			if err == nil {
				return val, nil
			}
		}
	} else if str, ok := fn.(*data.StringValue); ok {
		funcName := str.Value
		if f, ok := ctx.GetVM().GetFunc(funcName); ok {
			callCtx := ctx.CreateContext(vars)
			callCtx.SetIndexZVal(0, data.NewZVal(data.NewStringValue(a)))
			callCtx.SetIndexZVal(1, data.NewZVal(data.NewStringValue(b)))
			v, acl := f.Call(callCtx)
			if acl != nil {
				return 0, acl
			}
			if iv, ok := v.(data.AsInt); ok {
				val, err := iv.AsInt()
				if err == nil {
					return val, nil
				}
			}
		}
	}

	return 0, nil
}

func (f *ArrayIntersectUkeyFunction) GetName() string {
	return "array_intersect_ukey"
}

func (f *ArrayIntersectUkeyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, data.Mixed{}),
	}
}

func (f *ArrayIntersectUkeyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.Mixed{}),
	}
}
