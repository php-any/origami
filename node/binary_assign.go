package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

type BinaryAssign struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryAssign(from data.From, left, right data.GetValue) BinaryExpression {
	switch l := left.(type) {
	case data.Variable:
		return &BinaryAssignVariable{
			Node:  NewNode(from),
			Left:  l,
			Right: right,
		}
	case *VariableList:
		return &BinaryAssignVariableList{
			Node:  NewNode(from),
			Left:  l,
			Right: right,
		}
	}

	return &BinaryAssign{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryAssign) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	if v, ok := rv.(data.Value); ok {
		switch l := b.Left.(type) {
		case *CallObjectProperty:
			temp, acl := l.Object.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			switch object := temp.(type) {
			case *data.ThisValue:
				object.SetProperty(l.Property, v)
				return v, nil
			case *data.ClassValue:
				object.SetProperty(l.Property, v)
				return v, nil
			case *data.ObjectValue:
				object.SetProperty(l.Property, v)
				return v, nil
			}
		case *IndexExpression:
			// 索引赋值 $this->where[key] = value
			idxExpr := l
			arrayVal, acl := idxExpr.Array.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if idxExpr.Index == nil {
				// TODO
			}
			indexVal, acl := idxExpr.Index.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}

			switch arr := arrayVal.(type) {
			case *data.ArrayValue:
				// 数组索引赋值
				i := 0
				if iv, ok := indexVal.(data.AsInt); ok {
					var err error
					i, err = iv.AsInt()
					if err != nil {
						return nil, data.NewErrorThrow(b.from, err)
					}
				} else {
					return nil, data.NewErrorThrow(b.from, errors.New("数组索引不是整数类型"))
				}
				if i < 0 {
					return nil, data.NewErrorThrow(b.from, errors.New("数组索引不能为负数"))
				}
				if i >= len(arr.Value) {
					// 自动扩容，填充 null
					for j := len(arr.Value); j <= i; j++ {
						arr.Value = append(arr.Value, data.NewNullValue())
					}
				}
				arr.Value[i] = v
				return v, nil
			case *data.ObjectValue:
				// 对象属性赋值
				if iv, ok := indexVal.(data.AsString); ok {
					arr.SetProperty(iv.AsString(), v)
					return v, nil
				} else {
					return nil, data.NewErrorThrow(b.from, errors.New("对象属性索引不是字符串类型"))
				}
			default:
				return nil, data.NewErrorThrow(b.from, errors.New("索引赋值仅支持数组或对象"))
			}

		default:
			return nil, data.NewErrorThrow(b.from, errors.New("TODO 赋值表达式遇到未支持的类型"))
		}
	}

	_ = rv

	return nil, data.NewErrorThrow(b.from, errors.New("TODO BinaryAssign"))
}

type BinaryAssignVariable struct {
	*Node `pp:"-"`
	Left  data.Variable
	Right data.GetValue
}

func (b *BinaryAssignVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	if v, ok := rv.(data.Value); ok {
		return v, b.Left.SetValue(ctx, v)
	}

	return nil, data.NewErrorThrow(b.from, errors.New("TODO BinaryAssign"))
}

type BinaryAssignVariableList struct {
	*Node `pp:"-"`
	Left  *VariableList
	Right data.GetValue
}

func (b *BinaryAssignVariableList) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 检查是否是ReturnControl
	if returnControl, ok := rv.(data.ReturnControl); ok {
		// 从ReturnControl中提取ArrayValue
		arrayValue := returnControl.ReturnValue()
		// 使用VariableList的SetValue方法来处理多变量赋值
		ctl := b.Left.SetValue(ctx, arrayValue)
		if ctl != nil {
			return nil, ctl
		}
		return arrayValue, nil
	}

	// 处理普通值
	if v, ok := rv.(data.Value); ok {
		// 使用VariableList的SetValue方法来处理多变量赋值
		ctl := b.Left.SetValue(ctx, v)
		if ctl != nil {
			return nil, ctl
		}
		return v, nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("多变量赋值失败"))
}
