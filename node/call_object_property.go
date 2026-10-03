package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CallObjectProperty 表示对象属性访问表达式
type CallObjectProperty struct {
	*Node    `pp:"-"`
	Object   data.GetValue // 对象表达式
	Property string        // 属性名
}

func (pe *CallObjectProperty) GetIndex() int {
	return -1
}

func (pe *CallObjectProperty) IsPropertyLvalue() {}

func (pe *CallObjectProperty) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	temp, acl := pe.Object.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if object, ok := temp.(*data.ThisValue); ok {
		temp = object.ClassValue
	}
	object, ok := temp.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("Cannot reference property on %T", temp))
	}
	return propertyReferenceSlot(ctx, object, pe.Property, pe.Node)
}

func propertyReferenceSlot(ctx data.Context, object *data.ClassValue, name string, from *Node) (*data.ZVal, data.Control) {
	property, declared := lookupObjectProperty(ctx, object, name)
	accessible := !declared || propertyAccessible(ctx, object, property, name)
	if accessible && declared {
		if p, ok := property.(*ClassProperty); ok && p.IsReadonly {
			return nil, data.NewErrorThrowByName(from.GetFrom(), data.NewError(from.GetFrom(), "Cannot acquire reference to readonly property", nil), "Error")
		}
		return property.GetZVal(object)
	}
	if accessible && object.ObjectValue.HasProperty(name) {
		return object.GetPropertyZVal(name)
	}
	if magic, found := object.GetMethod("__get"); found {
		proxy := &CallObjectProperty{Node: from, Property: name}
		value, ctl := proxy.invokeMagicGetReference(ctx, object, magic, name)
		if ctl != nil {
			return nil, ctl
		}
		if value, ok := value.(data.Value); ok {
			if slot, ctl := data.ReferenceSlot(value); slot != nil || ctl != nil {
				return slot, ctl
			}
		}
		return nil, data.NewErrorThrowByName(from.GetFrom(), data.NewError(from.GetFrom(), "Cannot assign by reference to overloaded object", nil), "Error")
	}
	if !accessible {
		return nil, data.NewErrorThrowByName(from.GetFrom(), fmt.Errorf("Cannot access non-public property %s", name), "Error")
	}
	if ctl := CheckDynamicPropertyCreation(object, name, from.GetFrom()); ctl != nil {
		return nil, ctl
	}
	slot, ctl := object.GetPropertyZVal(name)
	if ctl == nil {
		ctl = ReportDynamicPropertyCreation(ctx, object, name, from.GetFrom())
	}
	return slot, ctl
}

func (pe *CallObjectProperty) GetName() string {
	return pe.Property
}

func (pe *CallObjectProperty) GetType() data.Types {
	return data.NewBaseType("")
}

func (pe *CallObjectProperty) SetValue(ctx data.Context, value data.Value) data.Control {
	_, ctl := pe.AssignValue(ctx, value)
	return ctl
}

func (pe *CallObjectProperty) AssignValue(ctx data.Context, value data.Value) (data.Value, data.Control) {
	temp, acl := pe.Object.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if this, ok := temp.(*data.ThisValue); ok {
		temp = this.ClassValue
	}
	switch object := temp.(type) {
	case *data.ClassValue:
		property, ok := lookupObjectProperty(ctx, object, pe.Property)
		if ok && !propertyAccessible(ctx, object, property, pe.Property) {
			if magic, found := object.GetMethod("__set"); found {
				return value, pe.invokeMagicSet(ctx, object, magic, pe.Property, value)
			}
			return nil, data.NewErrorThrowByName(pe.GetFrom(), fmt.Errorf("Cannot access non-public property %s::$%s", object.Class.GetName(), pe.Property), "Error")
		}
		if ok {
			var cloneState *readonlyCloneState
			if readonly, yes := property.(*ClassProperty); yes && readonly.IsReadonly {
				if object.ObjectValue.HasProperty(data.PropertyStorageName(property)) {
					cloneState = readonlyClonePermission(ctx, object.ObjectValue, data.PropertyStorageName(property))
					if cloneState == nil {
						return nil, data.NewErrorThrowByName(pe.GetFrom(), fmt.Errorf("Cannot modify readonly property %s::$%s", readonly.DeclaringClass, pe.Property), "Error")
					}
				}
				owner, _ := ctx.GetVM().GetClass(readonly.DeclaringClass)
				if !memberAccessible(ctx, data.ModifierProtected, owner) {
					return nil, data.NewErrorThrowByName(pe.GetFrom(), fmt.Errorf("Cannot initialize readonly property %s::$%s from global scope", readonly.DeclaringClass, pe.Property), "Error")
				}
			}
			prepared, accepted, conversion := data.PreparePropertyValue(propertyObject(object), data.PropertyStorageName(property), property.GetType(), value, ctx)

			if conversion != nil {
				return nil, conversion
			}
			if !accepted {
				return nil, data.NewTypeError(pe.GetFrom(), fmt.Errorf("%s 属性 %s 因为类型不一致无法赋值", TryGetCallClassName(object), pe.Property))
			}
			result, ctl := assignPropertyValue(ctx, propertyObject(object), data.PropertyStorageName(property), prepared)
			if ctl == nil && cloneState != nil {
				cloneState.written[data.PropertyStorageName(property)] = true
			}
			return result, ctl
		}
		// 无声明属性时尝试 __set(string $name, mixed $value)
		if object.ObjectValue.HasProperty(pe.Property) {
			return value, object.SetProperty(pe.Property, value)
		}
		if magic, hasSet := object.GetMethod("__set"); hasSet {
			return value, pe.invokeMagicSet(ctx, object, magic, pe.Property, value)
		}
		if ctl := CheckDynamicPropertyCreation(object, pe.Property, pe.GetFrom()); ctl != nil {
			return nil, ctl
		}
		if ctl := object.SetProperty(pe.Property, value); ctl != nil {
			return nil, ctl
		}
		return value, ReportDynamicPropertyCreation(ctx, object, pe.Property, pe.GetFrom())
	case data.SetProperty:
		return value, object.SetProperty(pe.Property, value)
	default:
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("Cannot assign property %s on %T", pe.Property, temp))
	}
}

// NewObjectProperty 创建一个新的对象属性访问表达式
func NewObjectProperty(token *TokenFrom, object data.GetValue, property string) *CallObjectProperty {
	return &CallObjectProperty{
		Node:     NewNode(token),
		Object:   object,
		Property: property,
	}
}

// invokeMagicIsset 调用 __isset(string $name)，用于 isset/?? 判断魔术属性是否存在
func (pe *CallObjectProperty) invokeMagicIsset(ctx, object data.Context, magic data.Method, name string) (bool, data.Control) {
	varies := magic.GetVariables()
	if len(varies) < 1 {
		return false, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__isset 需要至少 1 个参数 (name)"))
	}
	fnCtx := implicitReceiverFrame(ctx, object, magic)
	defer tryReleaseCallContext(magic, fnCtx)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(name))
	rv, acl := magic.Call(fnCtx)
	if acl != nil {
		return false, acl
	}
	if bv, ok := rv.(data.AsBool); ok {
		if b, err := bv.AsBool(); err == nil {
			return b, nil
		}
	}
	return false, nil
}

// invokeMagicGet 调用 __get(string $name)，用于读取不存在或不可见属性时的魔法分发
func (pe *CallObjectProperty) invokeMagicGet(ctx, object data.Context, magic data.Method, name string) (data.GetValue, data.Control) {
	return callValue(pe.invokeMagicGetReference(ctx, object, magic, name))
}

func (pe *CallObjectProperty) invokeMagicGetReference(ctx, object data.Context, magic data.Method, name string) (data.GetValue, data.Control) {
	varies := magic.GetVariables()
	if len(varies) < 1 {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__get 需要至少 1 个参数 (name)"))
	}
	fnCtx := implicitReceiverFrame(ctx, object, magic)
	defer tryReleaseCallContext(magic, fnCtx)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(name))
	return magic.Call(fnCtx)
}

// invokeMagicSet 调用 __set(string $name, mixed $value)，用于写入不存在或不可见属性时的魔法分发
func (pe *CallObjectProperty) invokeMagicSet(ctx, object data.Context, magic data.Method, name string, value data.Value) data.Control {
	varies := magic.GetVariables()
	if len(varies) < 2 {
		return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("__set 需要至少 2 个参数 (name, value)"))
	}
	fnCtx := implicitReceiverFrame(ctx, object, magic)
	defer tryReleaseCallContext(magic, fnCtx)
	fnCtx.SetVariableValue(varies[0], data.NewStringValue(name))
	fnCtx.SetVariableValue(varies[1], value)
	_, acl := magic.Call(fnCtx)
	return acl
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallObjectProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		if ctl, ok := ctl.(data.AddStack); ok {
			ctl.AddStackWithInfo(pe.from, TryGetCallClassName(pe.Object), pe.Property)
		}
		return nil, ctl
	}
	if this, ok := o.(*data.ThisValue); ok {
		o = this.ClassValue
	}
	switch v := o.(type) {
	case *data.NullValue:
		return data.NewNullValue(), nil
	case *data.ClassValue:
		property, ok := lookupObjectProperty(ctx, v, pe.Property)
		if ok {
			if !propertyAccessible(ctx, v, property, pe.Property) {
				if magic, found := v.GetMethod("__get"); found {
					return pe.invokeMagicGet(ctx, v, magic, pe.Property)
				}
				return nil, data.NewErrorThrowByName(pe.GetFrom(), fmt.Errorf("Cannot access non-public property %s::$%s", v.Class.GetName(), pe.Property), "Error")
			}
			return property.GetValue(v)
		}
		// 无声明属性时尝试 __get(string $name)
		if magic, hasGet := v.GetMethod("__get"); hasGet {
			return pe.invokeMagicGet(ctx, v, magic, pe.Property)
		}
		dynVal, acl := v.ObjectValue.GetProperty(pe.Property)
		if acl != nil {
			return nil, acl
		}
		return dynVal, nil
	case data.GetProperty:
		ov, acl := v.GetProperty(pe.Property)
		if acl != nil {
			return nil, acl
		}
		return ov.GetValue(ctx)
	default:
		return nil, data.NewErrorThrow(pe.from, fmt.Errorf("值(%s)不是对象, 不能操作属性(%s)", TryGetCallClassName(pe.Object), pe.Property))
	}
}

// Public properties do not need caller inspection or an inheritance walk.
func propertyAccessible(ctx data.Context, object *data.ClassValue, property data.Property, name string) bool {
	if property.GetModifier() == data.ModifierPublic {
		return true
	}
	if property, ok := property.(*ClassProperty); ok && property.DeclaringClass != "" {
		class, _ := ctx.GetVM().GetClass(property.DeclaringClass)
		return memberAccessible(ctx, property.GetModifier(), class)
	}
	class := object.Class
	for class != nil {
		if _, found := class.GetProperty(name); found {
			return memberAccessible(ctx, property.GetModifier(), class)
		}
		if class.GetExtend() == nil {
			break
		}
		class, _ = ctx.GetVM().GetClass(*class.GetExtend())
	}
	return false
}

func lookupObjectProperty(ctx data.Context, object *data.ClassValue, name string) (data.Property, bool) {
	var lexical data.ClassStmt
	switch owner := ctx.(type) {
	case *data.ClassMethodContext:
		lexical = owner.SelfClass
	case *data.BoundContext:
		if owner.ScopeClass != "" {
			lexical, _ = ctx.GetVM().GetClass(owner.ScopeClass)
		}
	}
	if lexical != nil && lexical != object.Class && data.NominalIsA(object.Class, lexical.GetName(), ctx.GetVM()) {
		if property, found := lexical.GetProperty(name); found && property.GetModifier() == data.ModifierPrivate {
			return property, true
		}
	}
	property, found := object.GetPropertyStmt(name)
	if found && property.GetModifier() == data.ModifierPrivate {
		if declared, ok := property.(*ClassProperty); ok && declared.DeclaringClass != object.Class.GetName() && !propertyAccessible(ctx, object, property, name) {
			return nil, false
		}
	}
	return property, found
}

func unsetObjectProperty(ctx data.Context, object *data.ClassValue, name string, from *Node) data.Control {
	property, declared := lookupObjectProperty(ctx, object, name)
	if declared && propertyAccessible(ctx, object, property, name) {
		if readonly, ok := property.(*ClassProperty); ok && readonly.IsReadonly {
			return data.NewErrorThrowByName(from.GetFrom(), fmt.Errorf("Cannot unset readonly property %s::$%s", object.Class.GetName(), name), "Error")
		}
		object.ObjectValue.UnsetProperty(data.PropertyStorageName(property))
		return nil
	}
	if magic, found := object.GetMethod("__unset"); found {
		proxy := &CallObjectProperty{Node: from}
		_, ctl := proxy.invokeMagicGet(ctx, object, magic, name)
		return ctl
	}
	if declared {
		return data.NewErrorThrowByName(from.GetFrom(), fmt.Errorf("Cannot unset non-public property %s::$%s", object.Class.GetName(), name), "Error")
	}
	object.ObjectValue.UnsetProperty(name)
	return nil
}
