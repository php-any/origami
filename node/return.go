package node

import (
	"fmt"
	"github.com/php-any/origami/data"
	"os"
)

// Reference returns are selected while parsing, so ordinary returns do not
// pay for lvalue discovery. The slot survives a pooled frame and is resolved
// once, even when its index or receiver expression has side effects.
type ReferenceReturnStatement struct{ *ReturnStatement }

func NewReferenceReturnStatement(from data.From, value data.GetValue) *ReferenceReturnStatement {
	return &ReferenceReturnStatement{&ReturnStatement{Node: NewNode(from), Value: value}}
}

func (r *ReferenceReturnStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// The result may retain a pointer to a local bucket after this call exits.
	markContextEscaped(ctx)
	if r.Value == nil {
		return nil, ctx.ReturnSlot(data.NewNullValue())
	}
	switch r.Value.(type) {
	case data.Variable, *IndexExpression, interface {
		GetZVal(data.Context) (*data.ZVal, data.Control)
	}:
		slot, ctl := argumentReferenceSlot(ctx, r.Value)
		if ctl != nil {
			return nil, ctl
		}
		return nil, ctx.ReturnSlot(&data.ArraySlotRef{Slot: slot})
	}
	var value data.GetValue
	var ctl data.Control
	if call, ok := r.Value.(referenceCall); ok {
		value, ctl = call.GetReferenceValue(ctx)
	} else {
		value, ctl = r.Value.GetValue(ctx)
	}
	if ctl != nil {
		return nil, ctl
	}
	if v, ok := value.(data.Value); ok {
		if slot, ctl := data.ReferenceSlot(v); ctl != nil {
			return nil, ctl
		} else if slot != nil {
			return nil, ctx.ReturnSlot(&data.ArraySlotRef{Slot: slot})
		}
		fmt.Fprintln(os.Stderr, "Notice: Only variable references should be returned by reference")
		return nil, ctx.ReturnSlot(&data.ArraySlotRef{Slot: data.NewZVal(v)})
	}
	return nil, ctx.ReturnSlot(data.NewNullValue())
}

func (u *ReturnStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if u.Value == nil {
		return nil, ctx.ReturnSlot(data.NewNullValue())
	}
	v, ctl := u.Value.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return nil, ctx.ReturnSlot(data.NewNullValue())
	}
	return nil, ctx.ReturnSlot(v.(data.Value))
}

// ReturnStatement 表示return语句
type ReturnStatement struct {
	*Node `pp:"-"`
	Value data.GetValue
}

// NewReturnStatement 创建一个新的return语句
func NewReturnStatement(from *TokenFrom, value data.GetValue) *ReturnStatement {
	return &ReturnStatement{
		Node:  NewNode(from),
		Value: value,
	}
}

// ReturnsStatement 表示多值 return 语句

type ReturnsStatement struct {
	*Node  `pp:"-"`
	Values []data.GetValue
}

// NewReturnsStatement 创建一个新的多值 return 语句
func NewReturnsStatement(from *TokenFrom, values []data.GetValue) *ReturnsStatement {
	return &ReturnsStatement{
		Node:   NewNode(from),
		Values: values,
	}
}

func (u *ReturnsStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var result []data.Value
	for _, expr := range u.Values {
		v, ctl := expr.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if val, ok := v.(data.Value); ok {
			result = append(result, val)
		} else {
			result = append(result, data.NewNullValue())
		}
	}
	return nil, ctx.ReturnSlot(data.NewArrayValue(result))
}
