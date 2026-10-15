package node

import (
	"testing"

	"github.com/php-any/origami/data"
)

type issetThrowExpression struct {
	control data.Control
}

func (e *issetThrowExpression) GetValue(data.Context) (data.GetValue, data.Control) {
	return nil, e.control
}

func TestIssetPropagatesPHPException(t *testing.T) {
	thrown := &data.ThrowValue{Name: "Exception", Object: &data.ClassValue{Class: &ClassStatement{Name: "Exception"}}}
	expr := NewIssetStatement(nil, []data.GetValue{&issetThrowExpression{control: thrown}})
	_, ctl := expr.GetValue(nil)
	if ctl != thrown {
		t.Fatalf("isset swallowed PHP exception: %v", ctl)
	}
}
