// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterClass 提供 PHP ReflectionParameter 类定义
// ReflectionParameter 用于获取方法参数的信息，包括参数名、位置、默认值等
type ReflectionParameterClass struct {
	node.Node
}

// GetValue 创建 ReflectionParameter 的实例
func (c *ReflectionParameterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionParameter"
func (c *ReflectionParameterClass) GetName() string { return "ReflectionParameter" }

// GetExtend 返回父类名，ReflectionParameter 没有父类
func (c *ReflectionParameterClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，ReflectionParameter 不实现任何接口
func (c *ReflectionParameterClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionParameter 没有属性
func (c *ReflectionParameterClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionParameter 没有属性
func (c *ReflectionParameterClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
func (c *ReflectionParameterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionParameterConstructMethod{}, true
	case "getName":
		return &ReflectionParameterGetNameMethod{}, true
	case "getPosition":
		return &ReflectionParameterGetPositionMethod{}, true
	case "isOptional":
		return &ReflectionParameterIsOptionalMethod{}, true
	case "isDefaultValueAvailable":
		return &ReflectionParameterIsDefaultValueAvailableMethod{}, true
	case "getDefaultValue":
		return &ReflectionParameterGetDefaultValueMethod{}, true
	case "getType":
		return &ReflectionParameterGetTypeMethod{}, true
	case "getDeclaringClass":
		return &ReflectionParameterGetDeclaringClassMethod{}, true
	case "isVariadic":
		return &ReflectionParameterIsVariadicMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionParameterClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionParameterConstructMethod{},
		&ReflectionParameterGetNameMethod{},
		&ReflectionParameterGetPositionMethod{},
		&ReflectionParameterIsOptionalMethod{},
		&ReflectionParameterIsDefaultValueAvailableMethod{},
		&ReflectionParameterGetDefaultValueMethod{},
		&ReflectionParameterGetTypeMethod{},
		&ReflectionParameterGetDeclaringClassMethod{},
		&ReflectionParameterIsVariadicMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionParameterClass) GetConstruct() data.Method {
	return &ReflectionParameterConstructMethod{}
}

// newReflectionParameter 创建一个新的 ReflectionParameter 实例（通过类名+方法名+索引）
func newReflectionParameter(ctx data.Context, className string, methodName string, paramIndex int, param data.GetValue) *data.ClassValue {
	paramClass := &ReflectionParameterClass{}
	paramValue := data.NewClassValue(paramClass, ctx.CreateBaseContext())

	paramValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	paramValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
	paramValue.ObjectValue.SetProperty("_paramIndex", data.NewIntValue(paramIndex))

	// 预先提取并存储参数信息，供 className=="" 时（Closure）直接使用
	paramValue.ObjectValue.SetProperty("_paramName", data.NewStringValue(extractParamName(param)))
	if isVariadicParam(param) {
		paramValue.ObjectValue.SetProperty("_isVariadic", data.NewBoolValue(true))
	}
	if ts := extractParamType(param); ts != "" {
		paramValue.ObjectValue.SetProperty("_paramType", data.NewStringValue(ts))
	}
	if paramHasDefault(param) {
		paramValue.ObjectValue.SetProperty("_hasDefault", data.NewBoolValue(true))
	}

	return paramValue
}

// newReflectionParameterFromVirtual 通过 virtualParam 直接创建 ReflectionParameter 实例（用于 Closure）
func newReflectionParameterFromVirtual(ctx data.Context, className string, methodName string, paramIndex int, vp *virtualParam) *data.ClassValue {
	paramClass := &ReflectionParameterClass{}
	paramValue := data.NewClassValue(paramClass, ctx.CreateBaseContext())

	paramValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	paramValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
	paramValue.ObjectValue.SetProperty("_paramIndex", data.NewIntValue(paramIndex))
	paramValue.ObjectValue.SetProperty("_paramName", data.NewStringValue(vp.name))
	if vp.variadic {
		paramValue.ObjectValue.SetProperty("_isVariadic", data.NewBoolValue(true))
	}
	if vp.typeStr != "" {
		paramValue.ObjectValue.SetProperty("_paramType", data.NewStringValue(vp.typeStr))
	}
	if vp.hasDefault {
		paramValue.ObjectValue.SetProperty("_hasDefault", data.NewBoolValue(true))
	}
	// 标记为 virtual（Closure 参数），后续 getReflectionParameterInfo 使用
	paramValue.ObjectValue.SetProperty("_isVirtual", data.NewBoolValue(true))

	return paramValue
}

// getReflectionParameterInfo 从上下文中获取 ReflectionParameter 的参数信息
// 当 _isVirtual=true 时，从预存属性构造 virtualParam 返回
func getReflectionParameterInfo(ctx data.Context) (string, string, int, data.GetValue) {
	objCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return "", "", -1, nil
	}
	if objCtx.ObjectValue == nil {
		return "", "", -1, nil
	}

	props := objCtx.ObjectValue.GetProperties()

	classNameVal, _ := props["_className"]
	methodNameVal, _ := props["_methodName"]
	paramIndexVal, hasParamIndex := props["_paramIndex"]

	if !hasParamIndex {
		return "", "", -1, nil
	}

	var className, methodName string
	var paramIndex int
	if sv, ok := classNameVal.(*data.StringValue); ok {
		className = sv.AsString()
	}
	if sv, ok := methodNameVal.(*data.StringValue); ok {
		methodName = sv.AsString()
	}
	if iv, ok := paramIndexVal.(*data.IntValue); ok {
		paramIndex, _ = iv.AsInt()
	}

	// Closure / virtual 路径：className 为空或 _isVirtual=true，从预存属性构造 virtualParam
	isVirtual := false
	if vv, ok := props["_isVirtual"]; ok {
		if bv, ok := vv.(*data.BoolValue); ok {
			isVirtual = bv.Value
		}
	}
	if className == "" || isVirtual {
		name := ""
		if nv, ok := props["_paramName"]; ok {
			if sv, ok := nv.(*data.StringValue); ok {
				name = sv.AsString()
			}
		}
		isVar := false
		if vv, ok := props["_isVariadic"]; ok {
			if bv, ok := vv.(*data.BoolValue); ok {
				isVar = bv.Value
			}
		}
		typeStr := ""
		if tv, ok := props["_paramType"]; ok {
			if sv, ok := tv.(*data.StringValue); ok {
				typeStr = sv.AsString()
			}
		}
		hasDef := false
		if dv, ok := props["_hasDefault"]; ok {
			if bv, ok := dv.(*data.BoolValue); ok {
				hasDef = bv.Value
			}
		}
		vp := &virtualParam{name: name, index: paramIndex, typeStr: typeStr, variadic: isVar, hasDefault: hasDef}
		return className, methodName, paramIndex, vp
	}

	// 常规路径：通过类名+方法名+索引查找参数
	if className != "" && methodName != "" {
		vm := ctx.GetVM()
		v, acl := vm.LoadPkg(className)
		if acl != nil {
			return "", "", -1, nil
		}
		if v != nil {
			if stmt, ok := v.(data.ClassStmt); ok {
				if method, exists := stmt.GetMethod(methodName); exists {
					params := method.GetParams()
					if paramIndex >= 0 && paramIndex < len(params) {
						return className, methodName, paramIndex, params[paramIndex]
					}
				}
			}
		}
	}

	return "", "", -1, nil
}

// virtualParam 是轻量级虚拟参数，用于 Closure 的 ReflectionParameter
type virtualParam struct {
	name       string
	index      int
	typeStr    string
	variadic   bool
	hasDefault bool
}

func (p *virtualParam) GetName() string { return p.name }
func (p *virtualParam) GetIndex() int   { return p.index }
func (p *virtualParam) GetType() data.Types {
	if p.typeStr == "" {
		return nil
	}
	return data.NewBaseType(p.typeStr)
}
func (p *virtualParam) IsVariadic() bool { return p.variadic }
func (p *virtualParam) HasDefault() bool { return p.hasDefault }
func (p *virtualParam) AsString() string { return p.name }
func (p *virtualParam) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(p.name), nil
}

// ---- 辅助函数：从 data.GetValue (参数节点) 提取信息 ----

// extractParamName 提取参数名
func extractParamName(param data.GetValue) string {
	type nameGetter interface{ GetName() string }
	if ng, ok := param.(nameGetter); ok {
		return ng.GetName()
	}
	return ""
}

// isVariadicParam 判断参数是否为可变参数
func isVariadicParam(param data.GetValue) bool {
	type variadicChecker interface{ IsVariadic() bool }
	if vc, ok := param.(variadicChecker); ok {
		return vc.IsVariadic()
	}
	return false
}

// extractParamType 提取参数类型字符串
func extractParamType(param data.GetValue) string {
	type typeGetter interface{ GetType() data.Types }
	if tg, ok := param.(typeGetter); ok {
		if t := tg.GetType(); t != nil {
			return t.String()
		}
	}
	return ""
}

// paramHasDefault 判断参数是否有默认值
func paramHasDefault(param data.GetValue) bool {
	type defaultChecker interface{ HasDefault() bool }
	if dc, ok := param.(defaultChecker); ok {
		return dc.HasDefault()
	}
	// node.Parameter 没有 HasDefault，检查是否能获取默认值（非 nil）
	type defaultValGetter interface{ GetDefaultValue() data.GetValue }
	if dg, ok := param.(defaultValGetter); ok {
		return dg.GetDefaultValue() != nil
	}
	return false
}
