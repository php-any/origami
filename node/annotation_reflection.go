package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ClassStmtHasAnnotationTargetParameter 注解构造函数是否包含编译器注入的 target 形参。
func ClassStmtHasAnnotationTargetParameter(stmt data.ClassStmt) bool {
	if stmt == nil {
		return false
	}
	construct := stmt.GetConstruct()
	if construct == nil {
		return false
	}
	for _, param := range construct.GetParams() {
		if asAnnotationTargetParam(param) != nil {
			return true
		}
	}
	return false
}

// InstantiateAnnotationForReflection 按注解声明绑定用户实参与 target，但不执行构造副作用（无 InitAnnotation）。
func InstantiateAnnotationForReflection(a *Annotation, ctx data.Context) (data.GetValue, data.Control) {
	stmt, acl := a.resolveClass(ctx)
	if acl != nil {
		return nil, acl
	}
	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := object.(*data.ClassValue)
	if !ok || cv.Class.GetConstruct() == nil {
		return object, nil
	}
	if _, acl := bindAnnotationConstructorArgs(a, ctx, cv); acl != nil {
		return nil, acl
	}
	return cv, nil
}

// bindAnnotationConstructorArgs 将声明中的用户实参与 target 绑定到构造上下文（不调用 __construct）。
func bindAnnotationConstructorArgs(a *Annotation, ctx data.Context, object *data.ClassValue) (data.Context, data.Control) {
	method := object.Class.GetConstruct()
	varies := method.GetVariables()
	params := method.GetParams()
	fnCtx := object.CreateContext(varies)

	for index, arg := range a.Arguments {
		switch argTV := arg.(type) {
		case *NamedArgument:
			tempV, acl := argTV.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			vari, err := findVariable(varies, argTV.Name)
			if err != nil {
				return nil, data.NewErrorThrow(a.from, err)
			}
			fnCtx.SetVariableValue(vari, tempV.(data.Value))
		default:
			tempV, acl := argTV.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if index >= len(varies) {
				return nil, data.NewErrorThrow(a.from, fmt.Errorf("注解(%v)构造函数参数数量超出限制: %d", a.Name, index))
			}
			fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
		}
	}

	for index := len(a.Arguments); index < len(params); index++ {
		if index >= len(varies) {
			break
		}
		param := params[index]
		if targetParam := asAnnotationTargetParam(param); targetParam != nil {
			fnCtx.SetVariableValue(targetParam, data.NewAnyValue(unwrapAnnotationTarget(a.Target)))
			continue
		}
		if argObj, ok := param.(*Parameter); ok {
			if argObj.DefaultValue == nil {
				return nil, data.NewErrorThrow(a.from, fmt.Errorf("调用 %s 构造函数时参数 %s 缺少值和默认值", a.Name, argObj.Name))
			}
			_, acl := argObj.GetValue(fnCtx)
			if acl != nil {
				return nil, acl
			}
		}
	}

	for index, param := range params {
		if promotedParam, ok := param.(*PromotedParameter); ok {
			if index < len(varies) {
				paramValue, acl := fnCtx.GetVariableValue(varies[index])
				if acl != nil {
					if promotedParam.DefaultValue != nil {
						paramValueGet, acl := promotedParam.DefaultValue.GetValue(fnCtx)
						if acl != nil {
							return nil, acl
						}
						if paramValueGet != nil {
							paramValue = paramValueGet.(data.Value)
						}
					} else {
						continue
					}
				}
				if paramValue != nil {
					object.SetProperty(promotedParam.PropertyName, paramValue)
				}
			}
		}
	}
	return fnCtx, nil
}
