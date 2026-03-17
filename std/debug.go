package std

import (
	"github.com/php-any/origami/node"
)
import "github.com/php-any/origami/data"

func NewDebugFunction() data.FuncStmt {
	return &DebugFunction{}
}

// 用于便捷测试，会在入参的ast节点处打断，返回入参的ast节点的值
type DebugFunction struct{}

func (f *DebugFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arg, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}

	return arg.(*data.ASTValue).Node.GetValue(ctx)
}
func (f *DebugFunction) GetName() string {
	return "gg"
}
func (f *DebugFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterRawAST(nil, "ast", 0, nil),
	}
}
func (f *DebugFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ast", 0, nil),
	}
}
