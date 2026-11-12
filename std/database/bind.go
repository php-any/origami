package database

import (
	"database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	sqlpkg "github.com/php-any/origami/std/database/sql"
	"github.com/php-any/origami/utils"
)

// NewBindExecutorFunction 让 Database\DB 使用 Database\Sql 返回的 *sql.DB 或 *sql.Tx。
func NewBindExecutorFunction() data.FuncStmt { return &BindExecutorFunction{} }

type BindExecutorFunction struct{}

func (f *BindExecutorFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	builderValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少查询器"))
	}
	targetValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少数据库执行器"))
	}
	builder, ok := scriptObject(builderValue)
	target, targetOK := scriptObject(targetValue)
	if !ok || !targetOK {
		return nil, utils.NewThrow(errors.New("查询器和执行器都必须是对象"))
	}
	state, ok := builderState(builder)
	if !ok {
		return nil, utils.NewThrow(errors.New("执行器只能绑定到 Database\\DB"))
	}
	switch class := target.Class.(type) {
	case *sqlpkg.TxClass:
		tx, _ := class.GetSource().(*sql.Tx)
		if tx == nil {
			return nil, utils.NewThrow(errors.New("事务尚未开始"))
		}
		state.exec = tx
	case *sqlpkg.DBClass:
		pool, _ := class.GetSource().(*sql.DB)
		if pool == nil {
			return nil, utils.NewThrow(errors.New("连接池尚未打开"))
		}
		state.exec = nil
		state.conn = pool
	default:
		return nil, utils.NewThrow(errors.New("执行器必须是 Database\\Sql\\DB 或 Database\\Sql\\Tx"))
	}
	return builder, nil
}

func (f *BindExecutorFunction) GetName() string { return "Database\\bind" }
func (f *BindExecutorFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "builder", 0, nil, data.NewBaseType("object")),
		node.NewParameter(nil, "executor", 1, nil, data.NewBaseType("object")),
	}
}
func (f *BindExecutorFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "builder", 0, data.NewBaseType("object")),
		node.NewVariable(nil, "executor", 1, data.NewBaseType("object")),
	}
}

func scriptObject(value data.GetValue) (*data.ClassValue, bool) {
	switch current := value.(type) {
	case *data.ClassValue:
		return current, true
	case *data.ThisValue:
		if current == nil || current.ClassValue == nil {
			return nil, false
		}
		return current.ClassValue, true
	case *data.ZValValue:
		if current.ZVal == nil {
			return nil, false
		}
		return scriptObject(current.ZVal.ReadValue())
	default:
		return nil, false
	}
}

func builderState(object *data.ClassValue) (*db, bool) {
	class, ok := object.Class.(*DBClass)
	if !ok || class == nil || class.queryMethod == nil {
		return nil, false
	}
	query, ok := class.queryMethod.(*DbQueryMethod)
	if !ok || query.source == nil {
		return nil, false
	}
	return query.source, true
}
