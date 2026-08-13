package php

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

// EvalFunction 实现 PHP eval()
type EvalFunction struct{}

func NewEvalFunction() data.FuncStmt {
	return &EvalFunction{}
}

func (f *EvalFunction) GetName() string {
	return "eval"
}

func (f *EvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterRawAST(nil, "code", 0, data.Mixed{}),
	}
}

func (f *EvalFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "code", 0, data.Mixed{}),
	}
}

func (f *EvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	codeVal, ok := ctx.GetIndexValue(0)
	if !ok || codeVal == nil {
		return data.NewBoolValue(false), nil
	}
	if astValue, ok := codeVal.(*data.ASTValue); ok {
		val, acl := astValue.Node.GetValue(astValue.Ctx)
		if acl != nil {
			return nil, acl
		}
		codeVal = val.(data.Value)
	}
	sv, ok := codeVal.(data.AsString)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	code := sv.AsString()

	var evalFrom data.From
	for _, arg := range ctx.GetCallArgs() {
		if g, ok := arg.(node.GetFrom); ok && g.GetFrom() != nil {
			evalFrom = g.GetFrom()
			break
		}
	}

	rvm, ok := ctx.GetVM().(*runtime.VM)
	if !ok {
		return nil, data.NewErrorThrow(evalFrom, fmt.Errorf("eval 需要 runtime.VM"))
	}
	return rvm.EvalCode(code, ctx, evalFrom)
}
