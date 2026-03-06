package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// IndexExpression 表示数组访问表达式
type IndexExpression struct {
	*Node `pp:"-"`
	Array data.GetValue // 数组表达式
	Index data.GetValue // 索引表达式
}

// NewIndexExpression 创建一个新的数组访问表达式
func NewIndexExpression(token *TokenFrom, array data.GetValue, index data.GetValue) *IndexExpression {
	return &IndexExpression{
		Node:  NewNode(token),
		Array: array,
		Index: index,
	}
}

func (ie *IndexExpression) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	temp, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	switch v := temp.(type) {
	case *data.ArrayValue:
		i := 0
		switch iv := index.(type) {
		case *data.IntValue:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		case *data.StringValue:
			if len(v.List) == 0 {
				return data.NewZVal(data.NewNullValue()), nil
			}

			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("未实现自动转化为对象的能力"))
		case data.AsInt:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		case *data.BoolValue:
			if iv.Value {
				i = 1
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		default:
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
		}

		return v.List[i], nil
	case *data.ObjectValue:
		switch iv := index.(type) {
		case data.AsString:
			ov, acl := v.GetZVal(iv.AsString())
			if acl != nil {
				return nil, acl
			}
			if ov == nil {
				return data.NewZVal(data.NewNullValue()), nil
			}
			return ov, nil
		}
	case *data.NullValue:
		return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("NULL值不支持数组索引操作"), "UndefinedIndexExpression")
	}
	return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("无法处理索引的类型值"), "UndefinedIndexExpression")
}

func (ie *IndexExpression) SetValue(ctx data.Context, value data.Value) data.Control {
	// 获取数组值
	arrayVal, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		// 在赋值语境下，未定义索引（UndefinedIndexExpression）应视为 null，
		// 以便支持 PHP 风格的链式自动创建：
		// $namespace['commands'][1] = 'foo';
		if tv, ok := acl.(*data.ThrowValue); ok && tv.Name == "UndefinedIndexExpression" {
			arrayVal = data.NewNullValue()
		} else {
			return acl
		}
	}

	// 如果数组当前为 null，根据左侧表达式类型决定如何自动创建容器
	if _, ok := arrayVal.(*data.NullValue); ok {
		switch base := ie.Array.(type) {
		case data.Variable:
			// 场景：$namespace['commands'] = ... 或 $namespace['commands'][1] = ...
			obj := data.NewObjectValue()
			if ctl := base.SetValue(ctx, obj); ctl != nil {
				return ctl
			}
			// 变量赋值时会对 ObjectValue 做 clone，这里需要重新从上下文读取，
			// 确保后续写入操作作用在真实存储的容器上。
			var errCtl data.Control
			arrayVal, errCtl = ie.Array.GetValue(ctx)
			if errCtl != nil {
				return errCtl
			}
		case *IndexExpression:
			// 多级访问，自动创建空对象并挂到上一层：
			// $config['db']['host'] = '...';
			obj := data.NewObjectValue()
			if _, errCtl := NewBinaryAssign(ie.GetFrom(), base, obj).GetValue(ctx); errCtl != nil {
				return errCtl
			}
			// 同样需要重新读取当前数组值，拿到实际挂在上一层上的容器（可能经历了 clone）
			var errCtl2 data.Control
			arrayVal, errCtl2 = ie.Array.GetValue(ctx)
			if errCtl2 != nil {
				return errCtl2
			}
		default:
			// 其它情况保持 null，后续类型分支会给出统一错误信息
		}
	}

	// 获取索引值
	indexVal, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return acl
	}

	switch arr := arrayVal.(type) {
	case *data.ArrayValue:
		// 数组索引赋值
		i := 0
		if iv, ok := indexVal.(data.AsInt); ok {
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return data.NewErrorThrow(ie.GetFrom(), err)
			}
		} else if iv, ok := indexVal.(data.AsString); ok {
			// 字符串索引：将数组转换为对象
			objectVal := data.NewObjectValue()
			valueList := arr.ToValueList()
			for i2, val := range valueList {
				objectVal.SetProperty(fmt.Sprintf("%d", i2), val)
			}
			objectVal.SetProperty(iv.AsString(), value)
			// 重新赋值
			_, acl = NewBinaryAssign(ie.GetFrom(), ie.Array, objectVal).GetValue(ctx)
			if acl != nil {
				return acl
			}
			return nil
		} else {
			return data.NewErrorThrow(ie.GetFrom(), errors.New("数组索引不是整数类型"))
		}

		if i < 0 {
			return data.NewErrorThrow(ie.GetFrom(), errors.New("数组索引不能为负数"))
		}

		// 如果索引超出范围，自动扩容
		if i >= len(arr.List) {
			for j := len(arr.List); j <= i; j++ {
				arr.List = append(arr.List, data.NewZVal(data.NewNullValue()))
			}
		}

		// 设置值
		arr.List[i] = data.NewZVal(value)
		return nil

	case data.SetProperty:
		// 对象属性赋值
		if iv, ok := indexVal.(data.AsString); ok {
			arr.SetProperty(iv.AsString(), value)
			return nil
		} else if iv, ok := indexVal.(data.AsInt); ok {
			// 整数索引转换为字符串
			if i, err := iv.AsInt(); err == nil {
				arr.SetProperty(fmt.Sprintf("%d", i), value)
				return nil
			}
		}
		return data.NewErrorThrow(ie.GetFrom(), errors.New("对象索引必须是字符串或整数"))

	case *data.ClassValue:
		// 类实例属性赋值
		if iv, ok := indexVal.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := arr.GetPropertyStmt(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				// 使用 SetProperty 方法设置属性值
				return arr.SetProperty(name, value)
			}
			return data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return data.NewErrorThrow(ie.GetFrom(), errors.New("ClassValue索引必须是字符串"))

	case *data.ThisValue:
		// $this[$name] 动态设置当前对象属性
		if iv, ok := indexVal.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := arr.Class.GetProperty(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				// ThisValue 包含 ClassValue，直接使用 SetProperty
				return arr.ClassValue.SetProperty(name, value)
			}
			return data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return data.NewErrorThrow(ie.GetFrom(), errors.New("ThisValue索引必须是字符串"))

	default:
		return data.NewErrorThrow(ie.GetFrom(), errors.New("无法设置索引表达式的值"))
	}
}

// GetValue 获取数组访问表达式的值
func (ie *IndexExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	temp, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	switch v := temp.(type) {
	case *data.ArrayValue:
		i := 0
		switch iv := index.(type) {
		case *data.IntValue:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		case *data.StringValue:
			if len(v.List) == 0 {
				return data.NewNullValue(), nil
			}

			return nil, data.NewErrorThrow(ie.GetFrom(), fmt.Errorf("array[%s] 未实现自动转化为对象的能力", iv.AsString()))
		case data.AsInt:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		case *data.BoolValue:
			if iv.Value {
				i = 1
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		default:
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
		}

		return v.List[i].Value, nil
	case *data.ObjectValue:
		// 支持整数索引（转换为字符串）和字符串索引
		var key string
		if iv, ok := index.(data.AsString); ok {
			key = iv.AsString()
		} else if iv, ok := index.(data.AsInt); ok {
			// 将整数索引转换为字符串
			if i, err := iv.AsInt(); err == nil {
				key = fmt.Sprintf("%d", i)
			} else {
				return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
			}
		} else {
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ObjectValue无法处理索引的类型值"))
		}
		ov, acl := v.GetProperty(key)
		if acl != nil {
			return nil, acl
		}
		return ov, nil
	case *data.ClassValue:
		// 支持对类实例通过字符串索引访问公开属性：
		// $obj[$name]，在动态属性语法 $obj->$name 降级为索引访问后会走到这里
		if iv, ok := index.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := v.GetPropertyStmt(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				return prop.GetValue(v)
			}
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ClassValue无法处理索引的类型值"))
	case *data.ThisValue:
		// $this[$name] 动态访问当前对象属性
		if iv, ok := index.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := v.Class.GetProperty(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				return prop.GetValue(ctx)
			}
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ThisValue无法处理索引的类型值"))
	case *data.StringValue:
		// 获取字符串指定位置符号
		if iv, ok := index.(data.AsInt); ok {
			var err error
			i, err := iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.Value) {
				return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("字符串索引超出范围"))
			}
			return data.NewStringValue(string(v.Value[i])), nil
		} else {
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("字符串无法处理非int值"))
		}
	case *data.NullValue:
		return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("null 无法处理索引的类型值"), "UndefinedIndexExpression")
	case *data.BoolValue:
		return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("bool 无法处理索引的类型值"), "UndefinedIndexExpression")
	}
	return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("无法处理索引的类型值"), "UndefinedIndexExpression")
}
