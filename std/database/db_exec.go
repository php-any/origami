package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DbExecuteMethod struct {
	source *db
}

func (d *DbExecuteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlStr, goArgs, ctl := parseSQLCallArgs(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return d.source.runExecute(ctx, sqlStr, goArgs)
}

func (d *DbExecuteMethod) GetName() string            { return "execute" }
func (d *DbExecuteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbExecuteMethod) GetIsStatic() bool          { return false }
func (d *DbExecuteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}
func (d *DbExecuteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}
func (d *DbExecuteMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// runExecute 执行 INSERT/UPDATE/DELETE/DDL，返回影响行数等信息。
func (d *db) runExecute(ctx data.Context, sqlStr string, goArgs []interface{}) (data.GetValue, data.Control) {
	conn := d.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	result, err := conn.Exec(sqlStr, goArgs...)
	if err != nil {
		return nil, utils.NewThrowf("执行 SQL 语句失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.NewThrowf("获取影响行数失败: %v", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		lastInsertId = 0
	}

	resultObj := data.NewObjectValue()
	resultObj.SetProperty("rowsAffected", data.NewIntValue(int(rowsAffected)))
	resultObj.SetProperty("lastInsertId", data.NewIntValue(int(lastInsertId)))
	resultObj.SetProperty("success", data.NewBoolValue(true))

	return resultObj, nil
}
