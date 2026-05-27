package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// Range 表示数组访问表达式
type Range struct {
	*Node `pp:"-"`
	Array data.GetValue // 数组表达式
	Start data.GetValue // 索引表达式
	Stop  data.GetValue // 索引表达式
}

func NewRange(token *TokenFrom, array data.GetValue, start, stop data.GetValue) *Range {
	return &Range{
		Node:  NewNode(token),
		Array: array,
		Start: start,
		Stop:  stop,
	}
}

// GetValue 获取数组访问表达式的值
func (ie *Range) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if ie.Array == nil {
		return ie.evalNumericRange(ctx)
	}
	temp, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	start := 0
	if ie.Start != nil {
		temp, acl := ie.Start.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if asIntV, ok := temp.(data.AsInt); ok {
			iTemp, err := asIntV.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			start = iTemp
		}
	}
	switch v := temp.(type) {
	case *data.ArrayValue:
		stop := len(v.List)
		if ie.Stop != nil {
			temp, acl := ie.Stop.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if asIntV, ok := temp.(data.AsInt); ok {
				iTemp, err := asIntV.AsInt()
				if err != nil {
					return nil, data.NewErrorThrow(ie.GetFrom(), err)
				}
				if iTemp >= len(v.List) {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New(fmt.Sprintf("数组索引超出范围, 索引(%v), 长度(%v)", iTemp, len(v.List))))
				}
				stop = iTemp + 1
			}
		}

		valueList := v.ToValueList()
		newArr := valueList[start:stop]
		return data.NewArrayValue(newArr), nil
	case *data.StringValue:
		stop := len(v.Value)
		if ie.Stop != nil {
			temp, acl := ie.Stop.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if asIntV, ok := temp.(data.AsInt); ok {
				iTemp, err := asIntV.AsInt()
				if err != nil {
					return nil, data.NewErrorThrow(ie.GetFrom(), err)
				}
				if iTemp >= len(v.Value) {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New(fmt.Sprintf("数组索引超出范围, 索引(%v), 长度(%v)", iTemp, len(v.Value))))
				}
				stop = iTemp
			}
		}

		return data.NewStringValue(v.Value[start:stop]), nil
	}

	return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理范围的类型值"))
}

// evalNumericRange 将 start..stop 生成为整数数组（含端点）
func (ie *Range) evalNumericRange(ctx data.Context) (data.GetValue, data.Control) {
	start, acl := ie.evalRangeEndpoint(ctx, ie.Start)
	if acl != nil {
		return nil, acl
	}
	stop, acl := ie.evalRangeEndpoint(ctx, ie.Stop)
	if acl != nil {
		return nil, acl
	}
	var list []data.Value
	if start <= stop {
		for i := start; i <= stop; i++ {
			list = append(list, data.NewIntValue(i))
		}
	} else {
		for i := start; i >= stop; i-- {
			list = append(list, data.NewIntValue(i))
		}
	}
	return data.NewArrayValue(list), nil
}

func (ie *Range) evalRangeEndpoint(ctx data.Context, expr data.GetValue) (int, data.Control) {
	if expr == nil {
		return 0, nil
	}
	temp, acl := expr.GetValue(ctx)
	if acl != nil {
		return 0, acl
	}
	if asIntV, ok := temp.(data.AsInt); ok {
		i, err := asIntV.AsInt()
		if err != nil {
			return 0, data.NewErrorThrow(ie.GetFrom(), err)
		}
		return i, nil
	}
	return 0, data.NewErrorThrow(ie.GetFrom(), errors.New("范围端点必须是整数"))
}
