package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CompactStatement evaluates names in the caller's symbol table.
type CompactStatement struct {
	*Node    `pp:"-"`
	VarNames []data.GetValue
}

func NewCompactStatement(token *TokenFrom, varNames []data.GetValue) *CompactStatement {
	return &CompactStatement{Node: NewNode(token), VarNames: varNames}
}

func (c *CompactStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	result := data.NewArrayValue(nil).(*data.ArrayValue)
	var collect func(data.Value) data.Control
	collect = func(value data.Value) data.Control {
		switch names := value.(type) {
		case *data.StringValue:
			name := names.Value
			if ctx.HasVariableByName(name) {
				if value, ok := ctx.GetVariableByName(name); ok {
					result.SetStringKey(name, value)
				}
			} else {
				return data.EmitPHPError(ctx, 2, fmt.Sprintf("compact(): Undefined variable $%s", name), c.GetFrom())
			}
		case *data.ArrayValue:
			for slots, i := names.View(), 0; i < slots.Len(); i++ {
				if ctl := collect(slots.At(i).ReadValue()); ctl != nil {
					return ctl
				}
			}
		default:
			return data.EmitPHPError(ctx, 2, "compact(): Argument must be string or array of strings", c.GetFrom())
		}
		return nil
	}
	for _, expression := range c.VarNames {
		value, ctl := expression.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if value, ok := value.(data.Value); ok {
			if ctl := collect(value); ctl != nil {
				return nil, ctl
			}
		}
	}
	return result, nil
}
