package preg

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregReplaceCallbackFunction 实现 preg_replace_callback
//
// 完整签名：
//
//	preg_replace_callback(
//	  string|array $pattern,
//	  callable $callback,
//	  string|array $subject,
//	  int $limit = -1,
//	  int &$count = null
//	): string|array|null
//
// - 支持 pattern / subject 为字符串或一维数组
// - 支持 limit（每个 pattern / subject 独立限制）
// - 对每次匹配调用 PHP 回调，参数为 matches 数组（与 PHP 一致）
type PregReplaceCallbackFunction struct{}

func NewPregReplaceCallbackFunction() data.FuncStmt {
	return &PregReplaceCallbackFunction{}
}

func (f *PregReplaceCallbackFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternVal, _ := ctx.GetIndexValue(0)
	callbackVal, _ := ctx.GetIndexValue(1)
	subjectVal, _ := ctx.GetIndexValue(2)
	limitVal, _ := ctx.GetIndexValue(3)
	countVal, _ := ctx.GetIndexValue(4)

	if patternVal == nil || callbackVal == nil || subjectVal == nil {
		return data.NewBoolValue(false), nil
	}

	// 归一化 subject 为切片
	subjects, subjectIsArray := toStringSlice(subjectVal)
	// 归一化 pattern
	patterns, _ := toStringSlice(patternVal)

	// 解析 limit：默认 -1；<=0 视为无限制
	limit := -1
	if limitVal != nil {
		if _, isNull := limitVal.(*data.NullValue); !isNull {
			if asInt, ok := limitVal.(data.AsInt); ok {
				if v, err := asInt.AsInt(); err == nil {
					limit = v
				}
			}
		}
	}

	// 解析 callback 为可调用的 FuncValue
	fn, ctl := f.resolveCallback(ctx, callbackVal)
	if ctl != nil {
		return nil, ctl
	}

	totalCount := 0
	results := make([]data.Value, 0, len(subjects))

	for _, subj := range subjects {
		replaced := subj
		localCount := 0

		for _, pat := range patterns {
			if pat == "" {
				continue
			}
			re, err := CompileAny(pat)
			if err != nil {
				return data.NewBoolValue(false), nil
			}

			// limit <= 0: 不限制匹配次数
			if limit <= 0 {
				replaced = re.ReplaceAllStringFunc(replaced, func(match string) string {
					// 构造 matches 数组（仅包含完整匹配；暂不支持分组）
					args := []data.Value{
						data.NewArrayValue([]data.Value{data.NewStringValue(match)}),
					}
					ret, ctl := f.callCallback(ctx, fn, args)
					if ctl != nil {
						// 将控制流抛回 VM
						ctx.GetVM().ThrowControl(ctl)
						return match
					}
					localCount++
					if ret == nil {
						return ""
					}
					return ret.AsString()
				})
			} else {
				remaining := limit
				replaced = re.ReplaceAllStringFunc(replaced, func(match string) string {
					if remaining <= 0 {
						return match
					}
					args := []data.Value{
						data.NewArrayValue([]data.Value{data.NewStringValue(match)}),
					}
					ret, ctl := f.callCallback(ctx, fn, args)
					if ctl != nil {
						ctx.GetVM().ThrowControl(ctl)
						return match
					}
					remaining--
					localCount++
					if ret == nil {
						return ""
					}
					return ret.AsString()
				})
			}
		}

		totalCount += localCount
		results = append(results, data.NewStringValue(replaced))
	}

	// 更新 $count（引用参数）
	if countVal != nil {
		if _, ok := countVal.(*data.NullValue); !ok {
			if ref, ok := countVal.(*data.ReferenceValue); ok {
				ref.Ctx.SetVariableValue(ref.Val, data.NewIntValue(totalCount))
			}
		}
	}

	if subjectIsArray {
		return data.NewArrayValue(results), nil
	}
	if len(results) == 0 {
		return data.NewStringValue(""), nil
	}
	return results[0], nil
}

// resolveCallback 将 PHP callable 解析为 *data.FuncValue
func (f *PregReplaceCallbackFunction) resolveCallback(ctx data.Context, cb data.Value) (*data.FuncValue, data.Control) {
	switch c := cb.(type) {
	case *data.FuncValue:
		return c, nil
	case *data.ArrayValue:
		// 与 call_user_func 类似：['ClassName', 'method']
		valueList := c.ToValueList()
		if len(valueList) < 2 {
			return nil, data.NewErrorThrow(nil, errors.New("preg_replace_callback 回调数组长度不足"))
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
			return nil, data.NewErrorThrow(nil, errors.New("preg_replace_callback 未找到方法: "+className+"::"+methodName))
		}
		fn, acl := node.NewStaticMethodFuncValue(stmt, method).GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if fv, ok := fn.(*data.FuncValue); ok {
			return fv, nil
		}
		return nil, data.NewErrorThrow(nil, errors.New("preg_replace_callback 回调不是函数值"))
	default:
		return nil, data.NewErrorThrow(nil, errors.New("preg_replace_callback 回调不可调用"))
	}
}

// callCallback 使用给定参数调用回调，并返回其结果值
func (f *PregReplaceCallbackFunction) callCallback(ctx data.Context, fn *data.FuncValue, args []data.Value) (data.Value, data.Control) {
	callCtx := ctx.CreateContext(make([]data.Variable, len(args)))
	for i := range args {
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

func (f *PregReplaceCallbackFunction) GetName() string {
	return "preg_replace_callback"
}

func (f *PregReplaceCallbackFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "callback", 1, nil, nil),
		node.NewParameter(nil, "subject", 2, nil, nil),
		node.NewParameter(nil, "limit", 3, node.NewIntLiteral(nil, "-1"), nil),
		node.NewParameter(nil, "count", 4, node.NewNullLiteral(nil), nil),
	}
}

func (f *PregReplaceCallbackFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "callback", 1, data.NewBaseType("callable")),
		node.NewVariable(nil, "subject", 2, data.NewBaseType("mixed")),
		node.NewVariable(nil, "limit", 3, data.NewBaseType("int")),
		node.NewVariable(nil, "count", 4, data.NewBaseType("int")),
	}
}
