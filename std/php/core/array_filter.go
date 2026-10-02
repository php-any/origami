package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ArrayFilterFunction 实现 array_filter 函数
//
// 使用 array_filter 过滤数组元素
//
// 语法: array_filter(array $array, ?callable $callback = null, int $mode = 0): array
//
// 参数:
//   - array: 要过滤的数组
//   - callback: 可选的回调函数，用于测试每个元素。如果为 null，则过滤掉所有 falsy 值
//   - mode: 可选参数，决定传递给回调函数的参数
//   - 0 (默认): 只传递值给回调函数
//   - ARRAY_FILTER_USE_KEY (2): 只传递键给回调函数
//   - ARRAY_FILTER_USE_BOTH (1): 传递值和键给回调函数
//
// 返回值: 返回过滤后的新数组，保留原数组的键（关联数组），包括整数键
//
// 使用示例:
//
//	// 过滤掉所有 falsy 值
//	$filtered = array_filter([0, 1, '', 'hello', null, false]);
//	// 结果: [1, 'hello']
//
//	// 使用回调函数过滤
//	$numbers = [1, 2, 3, 4, 5];
//	$evens = array_filter($numbers, fn($n) => $n % 2 == 0);
//	// 结果: [2, 4]
//
//	// 使用函数名
//	$strings = ['hello', '', 'world', null];
//	$nonEmpty = array_filter($strings, 'strlen');
//	// 结果: ['hello', 'world']
//
//	// 只使用键过滤
//	$arr = ['a' => 1, 'b' => 2, 'c' => 3];
//	$filtered = array_filter($arr, fn($key) => $key != 'b', ARRAY_FILTER_USE_KEY);
//	// 结果: ['a' => 1, 'c' => 3]
//
//	// 使用值和键
//	$arr = ['a' => 1, 'b' => 2, 'c' => 3];
//	$filtered = array_filter($arr, fn($val, $key) => $val > 1 && $key != 'c', ARRAY_FILTER_USE_BOTH);
//	// 结果: ['b' => 2]
type ArrayFilterFunction struct{}

func NewArrayFilterFunction() data.FuncStmt {
	return &ArrayFilterFunction{}
}

func (f *ArrayFilterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	callback, _ := ctx.GetIndexValue(1)
	modeValue, _ := ctx.GetIndexValue(2)
	mode := 0
	if integer, ok := modeValue.(data.AsInt); ok {
		mode, _ = integer.AsInt()
	}
	var function *data.FuncValue
	if callback != nil {
		if _, null := callback.(*data.NullValue); !null {
			var control data.Control
			function, control = f.resolveCallback(ctx, callback)
			if control != nil {
				return nil, control
			}
		}
	}
	result := data.NewArrayValue(nil).(*data.ArrayValue)
	filter := func(key, element data.Value) data.Control {
		keep := isTruthy(element)
		if function != nil {
			args := []data.Value{element}
			switch mode {
			case 2:
				args = []data.Value{key}
			case 1:
				args = []data.Value{element, key}
			}
			ret, control := f.callCallback(ctx, function, args)
			if control != nil {
				return control
			}
			keep = isTruthy(ret)
		}
		if keep {
			result.SetKey(key, element)
		}
		return nil
	}
	switch array := value.(type) {
	case *data.ArrayValue:
		for slots, position := array.View(), 0; position < slots.Len(); position++ {
			slot := slots.At(position)
			if slot != nil {
				if control := filter(slot.PHPArrayKey(position), slot.Value); control != nil {
					return nil, control
				}
			}
		}
	case *data.ObjectValue:
		var control data.Control
		array.RangeProperties(func(name string, element data.Value) bool {
			var key data.Value = data.NewStringValue(name)
			if integer, ok := data.ParseIntArrayKeyName(name); ok {
				key = data.NewIntValue(integer)
			}
			control = filter(key, element)
			return control == nil
		})
		if control != nil {
			return nil, control
		}
	}
	return result, nil
}

// resolveCallback 解析回调函数
func (f *ArrayFilterFunction) resolveCallback(ctx data.Context, cb data.GetValue) (*data.FuncValue, data.Control) {
	switch c := cb.(type) {
	case *data.FuncValue:
		return c, nil
	case *data.ArrayValue:
		valueList := c.ToValueList()
		if len(valueList) < 2 {
			return nil, utils.NewThrow(errors.New("array_filter 回调数组长度不足"))
		}
		className := valueList[0].AsString()
		methodName := valueList[1].AsString()

		stmt, acl := ctx.GetVM().GetOrLoadClass(className)
		if acl != nil {
			return nil, acl
		}
		var method data.Method
		var ok bool
		method, ok = stmt.GetMethod(methodName)
		if !ok {
			if sm, ok2 := stmt.(data.GetStaticMethod); ok2 {
				method, ok = sm.GetStaticMethod(methodName)
			}
		}
		if !ok {
			return nil, utils.NewThrow(errors.New("array_filter 未找到方法: " + className + "::" + methodName))
		}
		fn, acl := node.NewStaticMethodFuncValue(stmt, method).GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if fv, ok := fn.(*data.FuncValue); ok {
			return fv, nil
		}
		return nil, utils.NewThrow(errors.New("array_filter 回调不是函数值"))
	default:
		// 尝试作为字符串函数名
		if str, ok := cb.(data.AsString); ok {
			funcName := str.AsString()
			fnStmt, exists := ctx.GetVM().GetFunc(funcName)
			if exists {
				fnValue := data.NewFuncValue(fnStmt)
				return fnValue, nil
			}
		}
		return nil, utils.NewThrow(errors.New("array_filter 回调不可调用"))
	}
}

// callCallback 使用给定参数调用回调，并返回其结果值
// 行为与 PregReplaceCallbackFunction.callCallback 保持一致，避免闭包/箭头函数语义分裂。
func (f *ArrayFilterFunction) callCallback(ctx data.Context, fn *data.FuncValue, args []data.Value) (data.Value, data.Control) {
	// 使用函数自身的变量表长度创建调用上下文，确保：
	// - 对于普通函数，slots 数量与形参一致
	// - 对于 LambdaExpression，slots 覆盖所有 f.vars（参数 + use 捕获变量），
	//   这样在 Lambda.Call 中通过 ctx.GetIndexZVal(i) 拷贝参数时不会越界。
	callCtx := ctx.CreateContext(fn.Value.GetVariables())
	callCtx.SetStrictTypes(false)
	if ctl := data.BindDeclaredArgs(callCtx, fn.Value, args); ctl != nil {
		return nil, ctl
	}
	ret, ctl := fn.Call(callCtx)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil {
		return nil, nil
	}
	if v, ok := ret.(data.Value); ok {
		return v, nil
	}
	return nil, nil
}

// isTruthy 检查值是否为 truthy
func isTruthy(v data.Value) bool {
	if v == nil {
		return false
	}

	// 检查 null
	if _, ok := v.(*data.NullValue); ok {
		return false
	}

	// 检查布尔值
	if boolVal, ok := v.(data.AsBool); ok {
		if b, err := boolVal.AsBool(); err == nil {
			return b
		}
	}

	// 检查整数 0
	if intVal, ok := v.(data.AsInt); ok {
		if i, err := intVal.AsInt(); err == nil {
			return i != 0
		}
	}

	// 检查浮点数 0.0
	if floatVal, ok := v.(data.AsFloat); ok {
		if f, err := floatVal.AsFloat(); err == nil {
			return f != 0.0
		}
	}

	// 检查空字符串
	if strVal, ok := v.(data.AsString); ok {
		return strVal.AsString() != ""
	}

	// 检查空数组
	if arrVal, ok := v.(*data.ArrayValue); ok {
		return arrVal.Len() > 0
	}

	// 检查空对象
	if objVal, ok := v.(*data.ObjectValue); ok {
		props := objVal.GetProperties()
		return len(props) > 0
	}

	// 其他值默认为 truthy
	return true
}

func (f *ArrayFilterFunction) GetName() string {
	return "array_filter"
}

var arrayFilterFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array")),
	node.NewParameter(nil, "callback", 1, node.NewNullLiteral(nil), data.NewNullableType(data.NewBaseType("callable"))),
	node.NewParameter(nil, "mode", 2, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
}

func (f *ArrayFilterFunction) GetParams() []data.GetValue {
	return arrayFilterFunctionGetParams
}

var arrayFilterFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	node.NewVariable(nil, "callback", 1, data.NewBaseType("callable")),
	node.NewVariable(nil, "mode", 2, data.NewBaseType("int")),
}

func (f *ArrayFilterFunction) GetVariables() []data.Variable {
	return arrayFilterFunctionGetVariables
}
