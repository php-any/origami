package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionFunctionClass 提供 PHP ReflectionFunction 类定义
// ReflectionFunction 用于获取函数/Closure 的信息，包括参数列表等
type ReflectionFunctionClass struct {
	node.Node
}

// GetValue 创建 ReflectionFunction 的实例
func (c *ReflectionFunctionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionFunction"
func (c *ReflectionFunctionClass) GetName() string { return "ReflectionFunction" }

// GetExtend 返回父类名，ReflectionFunction 没有父类
func (c *ReflectionFunctionClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表
func (c *ReflectionFunctionClass) GetImplements() []string { return nil }

// GetProperty 获取属性
func (c *ReflectionFunctionClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表
func (c *ReflectionFunctionClass) GetPropertyList() []data.Property { return nil }

// GetMethod 根据方法名获取方法
func (c *ReflectionFunctionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionFunctionConstructMethod{}, true
	case "getParameters":
		return &ReflectionFunctionGetParametersMethod{}, true
	case "getNumberOfParameters":
		return &ReflectionFunctionGetNumberOfParametersMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionFunctionClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionFunctionConstructMethod{},
		&ReflectionFunctionGetParametersMethod{},
		&ReflectionFunctionGetNumberOfParametersMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionFunctionClass) GetConstruct() data.Method {
	return &ReflectionFunctionConstructMethod{}
}

// ---- __construct ----

type ReflectionFunctionConstructMethod struct{}

func (m *ReflectionFunctionConstructMethod) GetName() string            { return "__construct" }
func (m *ReflectionFunctionConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionFunctionConstructMethod) GetIsStatic() bool          { return false }
func (m *ReflectionFunctionConstructMethod) GetReturnType() data.Types  { return nil }

func (m *ReflectionFunctionConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "function", 0, nil, data.Mixed{}),
	}
}

func (m *ReflectionFunctionConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "function", 0, data.Mixed{}),
	}
}

// Call 执行构造函数：接受 Closure 或函数名字符串，存储到实例属性
func (m *ReflectionFunctionConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	funcVal, _ := ctx.GetIndexValue(0)
	if funcVal == nil {
		return nil, nil
	}

	objCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, nil
	}

	switch fv := funcVal.(type) {
	case *data.FuncValue:
		// 存储 Closure 的参数列表供后续 getParameters 使用
		params := fv.Value.GetParams()
		objCtx.ObjectValue.SetProperty("_isClosure", data.NewBoolValue(true))
		// 将参数数量存储为整数
		objCtx.ObjectValue.SetProperty("_paramCount", data.NewIntValue(len(params)))
		// 逐个存储参数名
		for i, p := range params {
			name := extractParamName(p)
			isVar := isVariadicParam(p)
			typeStr := extractParamType(p)
			hasDef := paramHasDefault(p)
			objCtx.ObjectValue.SetProperty("_pname_"+data.NewIntValue(i).AsString(), data.NewStringValue(name))
			if isVar {
				objCtx.ObjectValue.SetProperty("_pvar_"+data.NewIntValue(i).AsString(), data.NewBoolValue(true))
			}
			if typeStr != "" {
				objCtx.ObjectValue.SetProperty("_ptype_"+data.NewIntValue(i).AsString(), data.NewStringValue(typeStr))
			}
			if hasDef {
				objCtx.ObjectValue.SetProperty("_pdef_"+data.NewIntValue(i).AsString(), data.NewBoolValue(true))
			}
		}
	case *data.StringValue:
		objCtx.ObjectValue.SetProperty("_isClosure", data.NewBoolValue(false))
		objCtx.ObjectValue.SetProperty("_funcName", fv)
	}

	return nil, nil
}

// ---- getParameters ----

type ReflectionFunctionGetParametersMethod struct{}

func (m *ReflectionFunctionGetParametersMethod) GetName() string { return "getParameters" }
func (m *ReflectionFunctionGetParametersMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionFunctionGetParametersMethod) GetIsStatic() bool         { return false }
func (m *ReflectionFunctionGetParametersMethod) GetReturnType() data.Types { return data.Arrays{} }
func (m *ReflectionFunctionGetParametersMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (m *ReflectionFunctionGetParametersMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ReflectionFunctionGetParametersMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	objCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}
	props := objCtx.ObjectValue.GetProperties()

	countVal, hasCount := props["_paramCount"]
	if !hasCount {
		return data.NewArrayValue([]data.Value{}), nil
	}
	count := 0
	if iv, ok := countVal.(*data.IntValue); ok {
		count, _ = iv.AsInt()
	}

	result := make([]data.Value, 0, count)
	for i := 0; i < count; i++ {
		idx := data.NewIntValue(i).AsString()
		name := ""
		if nv, ok := props["_pname_"+idx]; ok {
			if sv, ok := nv.(*data.StringValue); ok {
				name = sv.AsString()
			}
		}
		isVar := false
		if vv, ok := props["_pvar_"+idx]; ok {
			if bv, ok := vv.(*data.BoolValue); ok {
				isVar = bv.Value
			}
		}
		typeStr := ""
		if tv, ok := props["_ptype_"+idx]; ok {
			if sv, ok := tv.(*data.StringValue); ok {
				typeStr = sv.AsString()
			}
		}
		hasDef := false
		if dv, ok := props["_pdef_"+idx]; ok {
			if bv, ok := dv.(*data.BoolValue); ok {
				hasDef = bv.Value
			}
		}
		vp := &virtualParam{name: name, index: i, typeStr: typeStr, variadic: isVar, hasDefault: hasDef}
		paramObj := newReflectionParameterFromVirtual(ctx, "", "", i, vp)
		result = append(result, paramObj)
	}

	return data.NewArrayValue(result), nil
}

// ---- getNumberOfParameters ----

type ReflectionFunctionGetNumberOfParametersMethod struct{}

func (m *ReflectionFunctionGetNumberOfParametersMethod) GetName() string {
	return "getNumberOfParameters"
}
func (m *ReflectionFunctionGetNumberOfParametersMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionFunctionGetNumberOfParametersMethod) GetIsStatic() bool         { return false }
func (m *ReflectionFunctionGetNumberOfParametersMethod) GetReturnType() data.Types { return data.Int{} }
func (m *ReflectionFunctionGetNumberOfParametersMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (m *ReflectionFunctionGetNumberOfParametersMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ReflectionFunctionGetNumberOfParametersMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	objCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewIntValue(0), nil
	}
	props := objCtx.ObjectValue.GetProperties()
	if cv, ok := props["_paramCount"]; ok {
		return cv, nil
	}
	return data.NewIntValue(0), nil
}
