package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StmtQueryRowMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtQueryRowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	// 警告：这是可变参数（variadic parameter）
	// 如果生成的代码有问题，请检查以下文件：
	// 1. 参数处理部分：可能需要调整 slice 展开逻辑
	// 2. GetParams 部分：可能需要使用 NewParametersReference 替代 NewParameter
	// 3. 方法调用部分：确保使用 ... 操作符展开 slice
	arg0 := make([]any, 0)
	for _, v := range a0.(*data.ArrayValue).Value {
		arg0 = append(arg0, v)
	}

	ret0 := h.source.QueryRow(arg0...)
	return data.NewClassValue(NewRowClassFrom(ret0), ctx), nil
}

func (h *StmtQueryRowMethod) GetName() string            { return "queryRow" }
func (h *StmtQueryRowMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtQueryRowMethod) GetIsStatic() bool          { return true }
func (h *StmtQueryRowMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

func (h *StmtQueryRowMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}

func (h *StmtQueryRowMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
