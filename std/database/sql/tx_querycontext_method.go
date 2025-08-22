package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type TxQueryContextMethod struct {
	source *sqlsrc.Tx
}

func (h *TxQueryContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	a2, ok := ctx.GetIndexValue(2)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 2"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	arg1 := a1.(*data.StringValue).AsString()
	// 警告：这是可变参数（variadic parameter）
	// 如果生成的代码有问题，请检查以下文件：
	// 1. 参数处理部分：可能需要调整 slice 展开逻辑
	// 2. GetParams 部分：可能需要使用 NewParametersReference 替代 NewParameter
	// 3. 方法调用部分：确保使用 ... 操作符展开 slice
	arg2 := make([]any, 0)
	for _, v := range a2.(*data.ArrayValue).Value {
		arg2 = append(arg2, v)
	}

	ret0, err := h.source.QueryContext(arg0, arg1, arg2...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewRowsClassFrom(ret0), ctx), nil
}

func (h *TxQueryContextMethod) GetName() string            { return "queryContext" }
func (h *TxQueryContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxQueryContextMethod) GetIsStatic() bool          { return true }
func (h *TxQueryContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "query", 1, nil, nil),
		node.NewParameters(nil, "args", 2, nil, nil),
	}
}

func (h *TxQueryContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "query", 1, nil),
		node.NewVariable(nil, "args", 2, nil),
	}
}

func (h *TxQueryContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
