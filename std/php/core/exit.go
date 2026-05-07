package core

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ExitFunction 实现 PHP 的 exit / die 语言结构（作为普通函数使用）
//
// 签名近似：
//
//	exit(int|string $status = 0): void
//
// 当前实现：
//   - 如果传入 int，则使用该值作为进程退出码
//   - 如果传入 string，则打印到 stdout，退出码为 0
//   - 如果未传参，则等价于 exit(0)
//
// 注意：这里直接调用 os.Exit() 终止当前进程，与 PHP 行为一致。
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
				_, _ = os.Stdout.WriteString(s)
			}
		}
	}
	os.Exit(code)
	return data.NewNullValue(), nil
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
