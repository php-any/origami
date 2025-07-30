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
				stop = iTemp + 1
			}
		}

		newArr := v.Value[start:stop]
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
