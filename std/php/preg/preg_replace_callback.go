package preg

import (
	"errors"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

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

	subjects, subjectIsArray := toStringSlice(subjectVal)
	patterns, _ := toStringSlice(patternVal)

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

			allMatches := re.FindAllStringSubmatchIndex(replaced, -1)
			if len(allMatches) == 0 {
				continue
			}

			maxMatches := len(allMatches)
			if limit > 0 && maxMatches > limit {
				maxMatches = limit
			}

			var sb strings.Builder
			pos := 0
			for mi := 0; mi < maxMatches; mi++ {
				loc := allMatches[mi]
				if len(loc) < 2 {
					continue
				}
				start := loc[0]
				end := loc[1]

				matchValues := make([]data.Value, 0, len(loc)/2)
				for g := 0; g < len(loc); g += 2 {
					if loc[g] >= 0 && loc[g+1] >= 0 && loc[g] < len(replaced) && loc[g+1] <= len(replaced) {
						matchValues = append(matchValues, data.NewStringValue(replaced[loc[g]:loc[g+1]]))
					} else {
						matchValues = append(matchValues, data.NewNullValue())
					}
				}

				sb.WriteString(replaced[pos:start])

				ret, ctl := f.callWithSubmatches(ctx, fn, matchValues)
				if ctl != nil {
					ctx.GetVM().ThrowControl(ctl)
					return matchValues[0], nil
				}
				localCount++
				if ret == nil || ret.AsString() == "" {
					sb.WriteString("")
				} else {
					sb.WriteString(ret.AsString())
				}

				pos = end
			}
			sb.WriteString(replaced[pos:])
			replaced = sb.String()
		}

		totalCount += localCount
		results = append(results, data.NewStringValue(replaced))
	}

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

func (f *PregReplaceCallbackFunction) resolveCallback(ctx data.Context, cb data.Value) (*data.FuncValue, data.Control) {
	switch c := cb.(type) {
	case *data.FuncValue:
		return c, nil
	case *data.ArrayValue:
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

// callWithSubmatches 使用子匹配数组调用回调
func (f *PregReplaceCallbackFunction) callWithSubmatches(ctx data.Context, fn *data.FuncValue, matchValues []data.Value) (data.Value, data.Control) {
	args := []data.Value{
		data.NewArrayValue(matchValues),
	}
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

func (f *PregReplaceCallbackFunction) GetName() string            { return "preg_replace_callback" }
func (f *PregReplaceCallbackFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *PregReplaceCallbackFunction) GetIsStatic() bool          { return false }
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
func (f *PregReplaceCallbackFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }
