package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CallObjectDynamicProperty 表示 $obj->$name 或 $obj->{expr} 动态属性访问
// 运行时优先查找声明的属性，找不到时回退到 ArrayAccess/offsetGet 或 __get
type CallObjectDynamicProperty struct {
	*Node    `pp:"-"`
	Object   data.GetValue // 对象表达式
	NameExpr data.GetValue // 属性名表达式（求值后为字符串）
}

func NewCallObjectDynamicProperty(from data.From, object data.GetValue, nameExpr data.GetValue) *CallObjectDynamicProperty {
	return &CallObjectDynamicProperty{
		Node:     NewNode(from),
		Object:   object,
		NameExpr: nameExpr,
	}
}

func (pe *CallObjectDynamicProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 求值对象表达式
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 求值属性名表达式
	raw, acl := pe.NameExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	name := raw.(data.Value).AsString()

	switch v := o.(type) {
	case *data.ThisValue:
		// 优先查找声明的属性（包括父类）
		if prop, ok := v.GetPropertyStmt(name); ok {
			// 从 ObjectValue 动态属性存储中获取值
			if val, ctl := v.ObjectValue.GetProperty(name); ctl == nil {
				if _, isNull := val.(*data.NullValue); !isNull {
					return val, nil
				}
			}
			// 动态属性中也没有，用默认值初始化
			if def := prop.GetDefaultValue(); def != nil {
				val, ctl := def.GetValue(v)
				if ctl != nil {
					return nil, ctl
				}
				return val, nil
			}
			return data.NewNullValue(), nil
		}
		// 未找到声明属性，回退到 ArrayAccess/offsetGet
		if checkArrayAccess(ctx, v.Class) {
			return callArrayAccessOffsetGet(ctx, v.ClassValue, raw.(data.Value))
		}
		// 尝试 __get
		if magic, hasGet := v.GetMethod("__get"); hasGet {
			return pe.invokeMagicGet(v, magic, name)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("this(%s) 不存在属性(%s)", v.Class.GetName(), name))

	case *data.ClassValue:
		// 优先查找声明的属性（包括父类）
		if prop, ok := v.GetPropertyStmt(name); ok {
			if val, ctl := v.ObjectValue.GetProperty(name); ctl == nil {
				if _, isNull := val.(*data.NullValue); !isNull {
					return val, nil
				}
			}
			if def := prop.GetDefaultValue(); def != nil {
				val, ctl := def.GetValue(v)
				if ctl != nil {
					return nil, ctl
				}
				return val, nil
			}
			return data.NewNullValue(), nil
		}
		// 未找到声明属性，回退到 ArrayAccess/offsetGet
		if checkArrayAccess(ctx, v.Class) {
			return callArrayAccessOffsetGet(ctx, v, raw.(data.Value))
		}
		// 尝试 __get
		if magic, hasGet := v.GetMethod("__get"); hasGet {
			return pe.invokeMagicGet(v, magic, name)
		}
		// 动态属性
		if val, ctl := v.ObjectValue.GetProperty(name); ctl == nil {
			if _, isNull := val.(*data.NullValue); !isNull {
				return val, nil
			}
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)不存在属性(%s)", v.Class.GetName(), name))

	case data.GetProperty:
		if val, ctl := v.GetProperty(name); ctl == nil {
			return val.GetValue(ctx)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("不支持动态属性访问"))

	default:
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("值不是对象, 不能操作动态属性(%s)", name))
	}
}

// invokeMagicGet 调用 __get(string $name)
func (pe *CallObjectDynamicProperty) invokeMagicGet(object data.Context, magic data.Method, name string) (data.GetValue, data.Control) {
	varies := magic.GetVariables()
	if len(varies) < 1 {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__get 需要至少 1 个参数 (name)"))
	}
	fnCtx := object.CreateContext(varies)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(name))
	return magic.Call(fnCtx)
}
