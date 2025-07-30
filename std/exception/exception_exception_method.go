package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ExceptionExceptionMethod struct {
	source *Exception
}

func (h *ExceptionExceptionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	msg, acl := ctx.GetVariableValue(node.NewVariable(nil, "msg", 0, nil))
	if acl != nil {
		return nil, acl
	}

	h.source.Exception(msg.(*data.StringValue).AsString())
	return nil, nil
}

func (h *ExceptionExceptionMethod) GetName() string {
	return "exception"
}

func (h *ExceptionExceptionMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ExceptionExceptionMethod) GetIsStatic() bool {
	return false
}

func (h *ExceptionExceptionMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "msg", 0, nil, nil),
	}
}

func (h *ExceptionExceptionMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "msg", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *ExceptionExceptionMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
