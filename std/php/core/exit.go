package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ExitFunction 实现 PHP 的 exit / die。
//
// 签名近似：
//
//	exit(int|string $status = 0): never
//
// 行为：
//   - int：返回 ExitControl(code)，由 CLI/HTTP 宿主决定是否终止进程
//   - string：先写入当前输出，再返回 ExitControl(0)
//   - 无参：等价于 exit(0)
type ExitFunction struct{}

func NewExitFunction() data.FuncStmt {
	return &ExitFunction{}
}

func (f *ExitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	statusVal, _ := ctx.GetIndexValue(0)

	code := 0
	if statusVal != nil {
		if asInt, ok := statusVal.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				code = v
			}
		} else {
			s := statusVal.AsString()
			if s != "" {
				data.EmitOutput(ctx, s)
			}
		}
	}
	return data.NewNullValue(), data.NewExitControl(code)
}

func (f *ExitFunction) GetName() string {
	return "exit"
}

func (f *ExitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "status", 0, node.NewNullLiteral(nil), data.NewBaseType("mixed")),
	}
}

func (f *ExitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "status", 0, data.NewBaseType("mixed")),
	}
}
