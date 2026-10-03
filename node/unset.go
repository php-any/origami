package node

import (

	"github.com/php-any/origami/data"
)

// UnsetStatement 表示 unset 语句
type UnsetStatement struct {
	*Node `pp:"-"`
	Args  []data.GetValue // 参数表达式列表
}

// NewUnsetStatement 创建一个新的 unset 语句
func NewUnsetStatement(token *TokenFrom, args []data.GetValue) *UnsetStatement {
	return &UnsetStatement{
		Node: NewNode(token),
		Args: args,
	}
}

// GetValue 获取 unset 语句的值
func (u *UnsetStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if len(u.Args) == 0 {
		return data.NewNullValue(), nil
	}

	for _, argExpr := range u.Args {
		// 优先处理对象属性：unset($obj->prop)，需要绕过类型检查
		if callProp, ok := argExpr.(*CallObjectProperty); ok {
			objValue, acl := callProp.Object.GetValue(ctx)
			if acl != nil || objValue == nil {
				continue
			}
			switch obj := objValue.(type) {

			case *data.ClassValue:
				if ctl := unsetObjectProperty(ctx, obj, callProp.Property, callProp.Node); ctl != nil {
					return nil, ctl
				}
			case *data.ThisValue:
				if ctl := unsetObjectProperty(ctx, obj.ClassValue, callProp.Property, callProp.Node); ctl != nil {
					return nil, ctl
				}
			}
			continue
		}

		// 变量类型（普通变量 $var）
		if variable, ok := argExpr.(data.Variable); ok {
			if global, ok := variable.(interface{ SuperglobalName() string }); ok {
				slot := ctx.GetVM().EnsureGlobalZVal(global.SuperglobalName())
				slot.ReleaseRefSlot()
				*slot = *data.NewNamedZValSlot(global.SuperglobalName())
				continue
			}
			// foreach (... as &$v) 后的 unset($v)：只断开本地引用，不能清空被引用的数组元素
			if zv := ctx.GetIndexZVal(variable.GetIndex()); zv != nil && zv.RefCount() > 0 {
				zv.ReleaseRefSlot()
				ctx.SetIndexZVal(variable.GetIndex(), data.NewZVal(data.NewNullValue()))
				continue
			}
			variable.SetValue(ctx, data.NewNullValue())
			continue
		}

		// 索引表达式：unset($arr['key'])
		if indexExpr, ok := argExpr.(*IndexExpression); ok {
			if _, global := indexExpr.Array.(*GlobalsArrayVariable); global {
				slot, ctl := globalsDimensionSlot(ctx, indexExpr)
				if ctl != nil {
					return nil, ctl
				}
				slot.ReleaseRefSlot()
				*slot = *data.NewNamedZValSlot(slot.Name)
				continue
			}
			if complexDimension(indexExpr) {
				resolved, ctl := resolveDimensionMode(ctx, indexExpr, false)
				if ctl != nil {
					return nil, ctl
				}
				indexExpr = resolved
				if owner, ok := resolved.Array.(*resolvedDimensionOwner); ok && owner.overloadedClass != "" {
					emitIndirectModificationNotice(u.GetFrom(), owner.overloadedClass)
					continue
				}
			}
			arrayValue, acl := indexExpr.Array.GetValue(ctx)
			if acl != nil || arrayValue == nil {
				continue
			}
			indexValue, acl := indexExpr.Index.GetValue(ctx)
			if acl != nil || indexValue == nil {
				continue
			}
			// PHP 8.1: null is treated as empty string (no deprecation for unset)
			if _, isNull := indexValue.(*data.NullValue); isNull {
				indexValue = data.NewStringValue("")
			}
			arrayValue = cowSeparateNestedArray(ctx, indexExpr.Array, arrayValue)
			switch arr := arrayValue.(type) {
			case *data.ArrayValue:
				if iv, ok := indexValue.(data.Value); ok {
					arr.UnsetKey(iv)
				}
				writeBackArrayProperty(ctx, indexExpr.Array, arr)

			case *data.ClassValue:
				if iv, ok := indexValue.(data.Value); ok && CheckArrayAccess(ctx, arr.Class) {
					if ctl := CallArrayAccessOffsetUnset(ctx, arr, iv); ctl != nil {
						return nil, ctl
					}
					continue
				}
				if sv, ok := indexValue.(data.AsString); ok {
					arr.SetProperty(sv.AsString(), data.NewNullValue())
				}
			case *data.ThisValue:
				if arr.ClassValue != nil && CheckArrayAccess(ctx, arr.Class) {
					if iv, ok := indexValue.(data.Value); ok {
						if ctl := CallArrayAccessOffsetUnset(ctx, arr.ClassValue, iv); ctl != nil {
							return nil, ctl
						}
					}
				}
			}
		}
	}

	return data.NewNullValue(), nil
}

func writeBackArrayProperty(ctx data.Context, arrayExpr data.GetValue, arr data.Value) {
	switch a := arrayExpr.(type) {
	case *CallObjectProperty:
		obj, acl := a.Object.GetValue(ctx)
		if acl != nil {
			return
		}
		if _, ok := obj.(*data.ClassValue); ok {
			proxy := &CallObjectProperty{Node: a.Node, Object: &literalGetValue{v: obj}, Property: a.Property}
			_, _ = proxy.AssignValue(ctx, arr)
		} else if tv, ok := obj.(*data.ThisValue); ok && tv.ClassValue != nil {
			proxy := &CallObjectProperty{Node: a.Node, Object: tv.ClassValue, Property: a.Property}
			_, _ = proxy.AssignValue(ctx, arr)
		}
	case *IndexExpression:
		// $obj[$k][$sub][] = $v：须把子数组写回父容器对应键，不能直接用子数组覆盖父属性
		indexVal, acl := a.Index.GetValue(ctx)
		if acl != nil {
			return
		}
		parentVal, acl := a.Array.GetValue(ctx)
		if acl != nil {
			return
		}
		_ = indexSetValueOnContainer(ctx, a, parentVal, indexVal, arr)
	}
}

// cowSeparateNestedArray 在 $a[$k][$sub] 这类嵌套写入前分离内层数组。
// 赋值/传参只浅拷贝外层，内层仍共享；不分离的话 unset($snapshot['memo']['children'])
// 会改到调用方（Livewire Checksum::generate）。只在嵌套索引上克隆，不碰 Call 热路径。
func cowSeparateNestedArray(ctx data.Context, arrayExpr data.GetValue, container data.GetValue) data.GetValue {
	if resolved, ok := arrayExpr.(*resolvedDimensionOwner); ok {
		data.CowSeparateZVal(resolved.slot)
		return resolved.slot.ReadValue()
	}
	if property, ok := arrayExpr.(*CallObjectProperty); ok {
		if slot, ctl := property.GetZVal(ctx); ctl == nil && slot != nil {
			data.CowSeparateZVal(slot)
			return slot.ReadValue()
		}
	}
	if global, ok := arrayExpr.(interface{ SuperglobalName() string }); ok {
		slot := ctx.GetVM().EnsureGlobalZVal(global.SuperglobalName())
		data.CowSeparateZVal(slot)
		return slot.ReadValue()
	}

	switch variable := arrayExpr.(type) {
	case *VariableExpression:
		slot := ctx.GetIndexZVal(variable.GetIndex())
		data.CowSeparateZVal(slot)
		return slot.ReadValue()
	case *VariableReference:
		slot := ctx.GetIndexZVal(variable.GetIndex())
		data.CowSeparateZVal(slot)
		return slot.ReadValue()
	}
	if _, ok := arrayExpr.(*IndexExpression); !ok {
		return container
	}
	switch v := container.(type) {
	case *data.ArrayValue:
		return data.CloneArrayValue(v)

	default:
		return container
	}
}
