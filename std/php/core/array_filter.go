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
//   - ARRAY_FILTER_USE_KEY (1): 只传递键给回调函数
//   - ARRAY_FILTER_USE_BOTH (2): 传递值和键给回调函数
//
// 返回值: 返回过滤后的新数组，保留原数组的键（关联数组）或重新索引（索引数组）
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
	// 获取第一个参数：数组
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 转换为数组
	var sourceArray []data.Value
	var sourceMap map[string]data.Value
	isAssociative := false

	switch arr := arrayValue.(type) {
	case *data.ArrayValue:
		sourceArray = arr.ToValueList()
	case *data.ObjectValue:
		// 对象作为关联数组处理
		sourceMap = arr.GetProperties()
		isAssociative = true
		// 转换为数组格式，保留键值对信息
		sourceArray = make([]data.Value, 0, len(sourceMap))
		for _, v := range sourceMap {
			sourceArray = append(sourceArray, v)
		}
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取第二个参数：回调函数（可选）
	callbackValue, _ := ctx.GetIndexValue(1)

	// 获取第三个参数：模式（可选，默认 0）
	modeValue, _ := ctx.GetIndexValue(2)
	mode := 0
	if modeValue != nil {
		if intVal, ok := modeValue.(data.AsInt); ok {
			mode, _ = intVal.AsInt()
		}
	}

	// 如果没有回调函数（参数缺省或显式传入 null），过滤掉所有 falsy 值
	noCallback := false
	if callbackValue == nil {
		noCallback = true
	} else {
		if _, isNull := callbackValue.(*data.NullValue); isNull {
			noCallback = true
		}
	}
	if noCallback {
		var result []data.Value
		if isAssociative {
			// 关联数组：保留键
			resultObj := data.NewObjectValue()
			for k, v := range sourceMap {
				if isTruthy(v) {
					resultObj.SetProperty(k, v)
				}
			}
			return resultObj, nil
		} else {
			// 索引数组
			for _, element := range sourceArray {
				if isTruthy(element) {
					result = append(result, element)
				}
			}
			return data.NewArrayValue(result), nil
		}
	}

	// 有回调函数，需要调用回调
	// 统一使用与 preg_replace_callback / array_map 一致的回调调用约定：
	// - 将参数写入新的函数上下文（CreateContext）
	// - 保留闭包/箭头函数对外部变量的捕获行为

	// 解析回调函数
	fn, acl := f.resolveCallback(ctx, callbackValue)
	if acl != nil {
		return nil, acl
	}

	// 处理关联数组
	if isAssociative {
		resultObj := data.NewObjectValue()
		for key, element := range sourceMap {
			var args []data.Value
			switch mode {
			case 1: // ARRAY_FILTER_USE_KEY - 只传递键
				args = []data.Value{data.NewStringValue(key)}
			case 2: // ARRAY_FILTER_USE_BOTH - 传递值和键
				args = []data.Value{element, data.NewStringValue(key)}
			default: // 0 - 只传递值
				args = []data.Value{element}
			}

			ret, ctl := f.callCallback(ctx, fn, args)
			if ctl != nil {
				return nil, ctl
			}
			if ret == nil {
				continue
			}

			// 使用 PHP 的 truthy 语义决定是否保留元素
			if boolVal, ok := ret.(data.AsBool); ok {
				if isTrue, err := boolVal.AsBool(); err == nil && isTrue {
					resultObj.SetProperty(key, element)
				}
			} else if isTruthy(ret) {
				resultObj.SetProperty(key, element)
			}
		}
		return resultObj, nil
	}

	// 处理索引数组
	var result []data.Value
	for i, element := range sourceArray {
		var args []data.Value
		switch mode {
		case 1: // ARRAY_FILTER_USE_KEY - 只传递键
			args = []data.Value{data.NewIntValue(i)}
		case 2: // ARRAY_FILTER_USE_BOTH - 传递值和键
			args = []data.Value{element, data.NewIntValue(i)}
		default: // 0 - 只传递值
			args = []data.Value{element}
		}

		ret, ctl := f.callCallback(ctx, fn, args)
		if ctl != nil {
			return nil, ctl
		}
		if ret == nil {
			continue
		}

		if boolVal, ok := ret.(data.AsBool); ok {
			if isTrue, err := boolVal.AsBool(); err == nil && isTrue {
				result = append(result, element)
			}
		} else if isTruthy(ret) {
			result = append(result, element)
		}
	}

	return data.NewArrayValue(result), nil
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
	vars := fn.Value.GetVariables()
	callCtx := ctx.CreateContext(vars)
	for i := range args {
		if i >= len(vars) {
			break
		}
		callCtx.SetIndexZVal(i, data.NewZVal(args[i]))
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
		return len(arrVal.List) > 0
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

func (f *ArrayFilterFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "callback", 1, node.NewNullLiteral(nil), data.NewNullableType(data.NewBaseType("callable"))),
		node.NewParameter(nil, "mode", 2, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
	}
}

func (f *ArrayFilterFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "callback", 1, data.NewBaseType("callable")),
		node.NewVariable(nil, "mode", 2, data.NewBaseType("int")),
	}
}
