package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type RowsScanMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsScanMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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

	if err := h.source.Scan(arg0...); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *RowsScanMethod) GetName() string            { return "scan" }
func (h *RowsScanMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsScanMethod) GetIsStatic() bool          { return true }
func (h *RowsScanMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParametersReference(nil, "dest", 0, nil, nil),
	}
}

func (h *RowsScanMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "dest", 0, nil),
	}
}

func (h *RowsScanMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
