package reflection

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const reflectorClassName = "Illuminate\\Support\\Reflector"

type ReflectorClass struct {
	node.Node
	methods map[string]data.Method
}

func NewReflectorClass() data.ClassStmt {
	c := &ReflectorClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *ReflectorClass) GetName() string                          { return reflectorClassName }
func (c *ReflectorClass) GetExtend() *string                       { return nil }
func (c *ReflectorClass) GetImplements() []string                  { return nil }
func (c *ReflectorClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *ReflectorClass) GetPropertyList() []data.Property         { return nil }
func (c *ReflectorClass) GetConstruct() data.Method                { return nil }
func (c *ReflectorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ReflectorClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *ReflectorClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *ReflectorClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *ReflectorClass) register() {
	c.methods["iscallable"] = kit.StaticMethod("isCallable", []string{"var", "syntaxOnly"}, 1, reflectorIsCallable)
	c.methods["getparameterclassname"] = kit.StaticMethod("getParameterClassName", []string{"parameter"}, -1, reflectorGetParameterClassName)
	c.methods["getparameterclassnames"] = kit.StaticMethod("getParameterClassNames", []string{"parameter"}, -1, reflectorGetParameterClassNames)
	c.methods["isparametersubclassof"] = kit.StaticMethod("isParameterSubclassOf", []string{"parameter", "className"}, -1, reflectorIsParameterSubclassOf)
	c.methods["isparameterbackedenumwithstringbackingtype"] = kit.StaticMethod("isParameterBackedEnumWithStringBackingType", []string{"parameter"}, -1, reflectorIsParameterBackedEnum)
	// getClassAttribute / getClassAttributes 依赖 ReflectionAttribute 全链，暂走 __call 回退不够；单独实现轻量版
	c.methods["getclassattribute"] = kit.StaticMethod("getClassAttribute", []string{"objectOrClass", "attribute", "ascend"}, 2, reflectorGetClassAttribute)
	c.methods["getclassattributes"] = kit.StaticMethod("getClassAttributes", []string{"objectOrClass", "attribute", "includeParents"}, 2, reflectorGetClassAttributes)
}

func reflectorIsCallable(ctx data.Context) (data.GetValue, data.Control) {
	varV := kit.Arg(ctx, 0)
	syntaxOnly := false
	if v := kit.Arg(ctx, 1); v != nil {
		if bv, ok := v.(data.AsBool); ok {
			syntaxOnly, _ = bv.AsBool()
		}
	}
	if av, ok := kit.Unwrap(varV).(*data.ArrayValue); ok && av != nil {
		list := av.ToValueList()
		if len(list) < 2 {
			return data.NewBoolValue(false), nil
		}
		methodV := list[1]
		methodName := ""
		if sv, ok := methodV.(*data.StringValue); ok {
			methodName = sv.AsString()
		} else if methodV != nil {
			methodName = methodV.AsString()
		} else {
			return data.NewBoolValue(false), nil
		}
		target := list[0]
		if syntaxOnly {
			if _, ok := target.(*data.StringValue); ok {
				return data.NewBoolValue(methodName != ""), nil
			}
			if _, ok := kit.Unwrap(target).(*data.ClassValue); ok {
				return data.NewBoolValue(methodName != ""), nil
			}
		}
		className := ""
		if sv, ok := target.(*data.StringValue); ok {
			className = sv.AsString()
		} else if cv, ok := kit.Unwrap(target).(*data.ClassValue); ok && cv != nil {
			className = cv.Class.GetName()
			if m, ok := cv.GetMethod(strings.ToLower(methodName)); ok && m != nil {
				return data.NewBoolValue(true), nil
			}
			if _, ok := cv.GetMethod("__call"); ok {
				return data.NewBoolValue(true), nil
			}
			return data.NewBoolValue(false), nil
		}
		if className != "" {
			stmt, ctl := ctx.GetVM().GetOrLoadClass(className)
			if ctl != nil || stmt == nil {
				return data.NewBoolValue(false), nil
			}
			if _, ok := stmt.GetMethod(strings.ToLower(methodName)); ok {
				return data.NewBoolValue(true), nil
			}
			if _, ok := stmt.GetMethod("__callstatic"); ok {
				return data.NewBoolValue(true), nil
			}
		}
		return data.NewBoolValue(false), nil
	}
	if fn, ok := ctx.GetVM().GetFunc("is_callable"); ok {
		fctx := ctx.CreateContext(fn.GetVariables())
		vars := fn.GetVariables()
		if len(vars) > 0 {
			_ = fctx.SetVariableValue(vars[0], varV)
		}
		if len(vars) > 1 {
			_ = fctx.SetVariableValue(vars[1], data.NewBoolValue(syntaxOnly))
		}
		ret, ctl := fn.Call(fctx)
		if ctl != nil {
			return nil, ctl
		}
		if bv, ok := ret.(data.AsBool); ok {
			b, _ := bv.AsBool()
			return data.NewBoolValue(b), nil
		}
	}
	switch kit.Unwrap(varV).(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return data.NewBoolValue(true), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func callParamMethod(ctx data.Context, param data.Value, method string, args ...data.Value) (data.GetValue, data.Control) {
	cv, ok := kit.Unwrap(param).(*data.ClassValue)
	if !ok || cv == nil {
		return data.NewNullValue(), nil
	}
	// Reflection* 的 GetMethod 多为大小写敏感 switch；先原名再小写。
	m, ok := cv.Class.GetMethod(method)
	if !ok {
		m, ok = cv.GetMethod(method)
	}
	if !ok {
		m, ok = cv.GetMethod(strings.ToLower(method))
	}
	if !ok || m == nil {
		return data.NewNullValue(), nil
	}
	vars := m.GetVariables()
	inner := cv.CreateContext(vars)
	fnCtx := data.WrapMethodFrame(inner, cv, cv.Class, cv.Class)
	for i, v := range vars {
		if i < len(args) && args[i] != nil {
			_ = fnCtx.SetVariableValue(v, args[i])
		}
	}
	fnCtx.SetFlatCallArgs(args)
	return m.Call(fnCtx)
}

func callObjMethod(ctx data.Context, cv *data.ClassValue, method string, args ...data.Value) (data.GetValue, data.Control) {
	if cv == nil {
		return data.NewNullValue(), nil
	}
	m, ok := cv.Class.GetMethod(method)
	if !ok {
		m, ok = cv.GetMethod(method)
	}
	if !ok {
		m, ok = cv.GetMethod(strings.ToLower(method))
	}
	if !ok || m == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s::%s does not exist.", cv.Class.GetName(), method))
	}
	vars := m.GetVariables()
	inner := cv.CreateContext(vars)
	fnCtx := data.WrapMethodFrame(inner, cv, cv.Class, cv.Class)
	for i, v := range vars {
		if i < len(args) && args[i] != nil {
			_ = fnCtx.SetVariableValue(v, args[i])
		}
	}
	fnCtx.SetFlatCallArgs(args)
	return m.Call(fnCtx)
}

func reflectorGetParameterClassName(ctx data.Context) (data.GetValue, data.Control) {
	param := kit.Arg(ctx, 0)
	typeRet, ctl := callParamMethod(ctx, param, "getType")
	if ctl != nil {
		return nil, ctl
	}
	typeCV, ok := typeRet.(*data.ClassValue)
	if !ok || typeCV == nil {
		return data.NewNullValue(), nil
	}
	// 仅 ReflectionNamedType
	name := typeCV.Class.GetName()
	if !strings.Contains(strings.ToLower(name), "namedtype") && name != "ReflectionNamedType" {
		return data.NewNullValue(), nil
	}
	builtinRet, ctl := callObjMethod(ctx, typeCV, "isBuiltin")
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := builtinRet.(data.AsBool); ok {
		if b, _ := bv.AsBool(); b {
			return data.NewNullValue(), nil
		}
	}
	return reflectorGetTypeName(ctx, param, typeCV)
}

func reflectorGetParameterClassNames(ctx data.Context) (data.GetValue, data.Control) {
	param := kit.Arg(ctx, 0)
	typeRet, ctl := callParamMethod(ctx, param, "getType")
	if ctl != nil {
		return nil, ctl
	}
	typeCV, ok := typeRet.(*data.ClassValue)
	if !ok || typeCV == nil {
		return data.NewArrayValue(nil), nil
	}
	className := typeCV.Class.GetName()
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if strings.Contains(strings.ToLower(className), "union") {
		typesRet, ctl := callObjMethod(ctx, typeCV, "getTypes")
		if ctl != nil {
			return nil, ctl
		}
		if av, ok := typesRet.(*data.ArrayValue); ok {
			i := 0
			for _, e := range av.List {
				if e == nil || e.Value == nil {
					continue
				}
				listed, ok := kit.Unwrap(e.Value).(*data.ClassValue)
				if !ok || listed == nil {
					continue
				}
				if !strings.Contains(strings.ToLower(listed.Class.GetName()), "named") {
					continue
				}
				builtinRet, ctl := callObjMethod(ctx, listed, "isBuiltin")
				if ctl != nil {
					continue
				}
				if bv, ok := builtinRet.(data.AsBool); ok {
					if b, _ := bv.AsBool(); b {
						continue
					}
				}
				tn, ctl := reflectorGetTypeName(ctx, param, listed)
				if ctl != nil || tn == nil || kit.IsNull(tn.(data.Value)) {
					continue
				}
				out.SetIntKey(i, tn.(data.Value))
				i++
			}
		}
		return out, nil
	}
	single, ctl := reflectorGetParameterClassName(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if single != nil && !kit.IsNull(single.(data.Value)) {
		out.SetIntKey(0, single.(data.Value))
	}
	return out, nil
}

func reflectorGetTypeName(ctx data.Context, param data.Value, namedType *data.ClassValue) (data.GetValue, data.Control) {
	nameRet, ctl := callObjMethod(ctx, namedType, "getName")
	if ctl != nil {
		return nil, ctl
	}
	typeName := ""
	if nameRet != nil {
		if v, ok := nameRet.(data.Value); ok {
			typeName = v.AsString()
		}
	}
	if typeName == "self" || typeName == "parent" {
		declRet, ctl := callParamMethod(ctx, param, "getDeclaringClass")
		if ctl != nil {
			return nil, ctl
		}
		if decl, ok := declRet.(*data.ClassValue); ok && decl != nil {
			if typeName == "self" {
				n, _ := callObjMethod(ctx, decl, "getName")
				if n != nil {
					return n, nil
				}
			}
			if typeName == "parent" {
				parentRet, _ := callObjMethod(ctx, decl, "getParentClass")
				if parent, ok := parentRet.(*data.ClassValue); ok && parent != nil {
					n, _ := callObjMethod(ctx, parent, "getName")
					if n != nil {
						return n, nil
					}
				}
			}
		}
	}
	return data.NewStringValue(typeName), nil
}

func reflectorIsParameterSubclassOf(ctx data.Context) (data.GetValue, data.Control) {
	paramClass, ctl := reflectorGetParameterClassName(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if paramClass == nil || kit.IsNull(paramClass.(data.Value)) {
		return data.NewBoolValue(false), nil
	}
	paramClassName := paramClass.(data.Value).AsString()
	className := ""
	if v := kit.Arg(ctx, 1); v != nil {
		className = v.AsString()
	}
	if paramClassName == "" || className == "" {
		return data.NewBoolValue(false), nil
	}
	stmt, ctl := ctx.GetVM().GetOrLoadClass(paramClassName)
	if ctl != nil || stmt == nil {
		return data.NewBoolValue(false), nil
	}
	// ReflectionClass::isSubclassOf
	rcStmt, ctl := ctx.GetVM().GetOrLoadClass("ReflectionClass")
	if ctl != nil || rcStmt == nil {
		return data.NewBoolValue(false), nil
	}
	rc := data.NewClassValue(rcStmt, ctx.CreateBaseContext())
	cons := rcStmt.GetConstruct()
	if cons != nil {
		cctx := rc.CreateContext(cons.GetVariables())
		data.BindDeclaredArgs(cctx, cons, []data.Value{data.NewStringValue(paramClassName)})
		if _, ctl := cons.Call(cctx); ctl != nil {
			return data.NewBoolValue(false), nil
		}
	}
	ret, ctl := callObjMethod(ctx, rc, "isSubclassOf", data.NewStringValue(className))
	if ctl != nil {
		return data.NewBoolValue(false), nil
	}
	if bv, ok := ret.(data.AsBool); ok {
		b, _ := bv.AsBool()
		return data.NewBoolValue(b), nil
	}
	return data.NewBoolValue(false), nil
}

func reflectorIsParameterBackedEnum(ctx data.Context) (data.GetValue, data.Control) {
	param := kit.Arg(ctx, 0)
	typeRet, ctl := callParamMethod(ctx, param, "getType")
	if ctl != nil {
		return nil, ctl
	}
	typeCV, ok := typeRet.(*data.ClassValue)
	if !ok || typeCV == nil {
		return data.NewBoolValue(false), nil
	}
	if !strings.Contains(strings.ToLower(typeCV.Class.GetName()), "named") {
		return data.NewBoolValue(false), nil
	}
	nameRet, ctl := callObjMethod(ctx, typeCV, "getName")
	if ctl != nil || nameRet == nil {
		return data.NewBoolValue(false), nil
	}
	enumClass := nameRet.(data.Value).AsString()
	if enumClass == "" {
		return data.NewBoolValue(false), nil
	}
	if fn, ok := ctx.GetVM().GetFunc("enum_exists"); ok {
		fctx := ctx.CreateContext(fn.GetVariables())
		data.BindDeclaredArgs(fctx, fn, []data.Value{data.NewStringValue(enumClass)})
		ret, ctl := fn.Call(fctx)
		if ctl != nil {
			return data.NewBoolValue(false), nil
		}
		if bv, ok := ret.(data.AsBool); ok {
			if b, _ := bv.AsBool(); !b {
				return data.NewBoolValue(false), nil
			}
		}
	}
	reStmt, ctl := ctx.GetVM().GetOrLoadClass("ReflectionEnum")
	if ctl != nil || reStmt == nil {
		return data.NewBoolValue(false), nil
	}
	re := data.NewClassValue(reStmt, ctx.CreateBaseContext())
	if cons := reStmt.GetConstruct(); cons != nil {
		cctx := re.CreateContext(cons.GetVariables())
		data.BindDeclaredArgs(cctx, cons, []data.Value{data.NewStringValue(enumClass)})
		if _, ctl := cons.Call(cctx); ctl != nil {
			return data.NewBoolValue(false), nil
		}
	}
	backed, ctl := callObjMethod(ctx, re, "isBacked")
	if ctl != nil {
		return data.NewBoolValue(false), nil
	}
	if bv, ok := backed.(data.AsBool); ok {
		if b, _ := bv.AsBool(); !b {
			return data.NewBoolValue(false), nil
		}
	}
	bt, ctl := callObjMethod(ctx, re, "getBackingType")
	if ctl != nil {
		return data.NewBoolValue(false), nil
	}
	if btcv, ok := bt.(*data.ClassValue); ok && btcv != nil {
		n, _ := callObjMethod(ctx, btcv, "getName")
		if n != nil && n.(data.Value).AsString() == "string" {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func reflectorGetClassAttribute(ctx data.Context) (data.GetValue, data.Control) {
	attrs, ctl := reflectorGetClassAttributes(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if cv, ok := attrs.(*data.ClassValue); ok && cv != nil {
		flat, ctl := callObjMethod(ctx, cv, "flatten")
		if ctl != nil {
			return nil, ctl
		}
		if fcv, ok := flat.(*data.ClassValue); ok && fcv != nil {
			return callObjMethod(ctx, fcv, "first")
		}
	}
	return data.NewNullValue(), nil
}

func reflectorGetClassAttributes(ctx data.Context) (data.GetValue, data.Control) {
	objectOrClass := kit.Arg(ctx, 0)
	attribute := kit.Arg(ctx, 1)
	includeParents := false
	if v := kit.Arg(ctx, 2); v != nil {
		if bv, ok := v.(data.AsBool); ok {
			includeParents, _ = bv.AsBool()
		}
	}
	rcStmt, ctl := ctx.GetVM().GetOrLoadClass("ReflectionClass")
	if ctl != nil || rcStmt == nil {
		return data.NewNullValue(), nil
	}
	rc := data.NewClassValue(rcStmt, ctx.CreateBaseContext())
	if cons := rcStmt.GetConstruct(); cons != nil {
		cctx := rc.CreateContext(cons.GetVariables())
		arg := objectOrClass
		if arg == nil {
			arg = data.NewNullValue()
		}
		data.BindDeclaredArgs(cctx, cons, []data.Value{arg})
		if _, ctl := cons.Call(cctx); ctl != nil {
			return nil, ctl
		}
	}
	attrName := ""
	if attribute != nil {
		attrName = attribute.AsString()
	}
	// 简化：仅当前类 attributes，不爬父类（includeParents 时仍返回一层 Collection map）
	attrsRet, ctl := callObjMethod(ctx, rc, "getAttributes", data.NewStringValue(attrName))
	if ctl != nil {
		return nil, ctl
	}
	instances := data.NewArrayValue(nil).(*data.ArrayValue)
	i := 0
	if av, ok := attrsRet.(*data.ArrayValue); ok {
		for _, e := range av.List {
			if e == nil || e.Value == nil {
				continue
			}
			if acv, ok := kit.Unwrap(e.Value).(*data.ClassValue); ok {
				inst, ctl := callObjMethod(ctx, acv, "newInstance")
				if ctl != nil {
					continue
				}
				if v, ok := inst.(data.Value); ok {
					instances.SetIntKey(i, v)
					i++
				}
			}
		}
	}
	collStmt, ctl := ctx.GetVM().GetOrLoadClass("Illuminate\\Support\\Collection")
	if ctl != nil || collStmt == nil {
		return instances, nil
	}
	coll := data.NewClassValue(collStmt, ctx.CreateBaseContext())
	if cons := collStmt.GetConstruct(); cons != nil {
		cctx := coll.CreateContext(cons.GetVariables())
		data.BindDeclaredArgs(cctx, cons, []data.Value{instances})
		if _, ctl := cons.Call(cctx); ctl != nil {
			return instances, nil
		}
	}
	if !includeParents {
		return coll, nil
	}
	// includeParents: Collection of className => Collection
	classNameRet, _ := callObjMethod(ctx, rc, "getName")
	className := ""
	if classNameRet != nil {
		className = classNameRet.(data.Value).AsString()
	}
	outer := data.NewArrayValue(nil).(*data.ArrayValue)
	outer.SetStringKey(className, coll)
	outerColl := data.NewClassValue(collStmt, ctx.CreateBaseContext())
	if cons := collStmt.GetConstruct(); cons != nil {
		cctx := outerColl.CreateContext(cons.GetVariables())
		data.BindDeclaredArgs(cctx, cons, []data.Value{outer})
		if _, ctl := cons.Call(cctx); ctl != nil {
			return coll, nil
		}
	}
	return outerColl, nil
}
