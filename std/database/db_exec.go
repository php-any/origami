package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DbExecMethod struct {
	source *db
}

func (d *DbExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	// 获取 SQL 语句
	sqlValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 SQL 语句"))
	}

	var sqlStr string
	if sqlStrValue, ok := sqlValue.(data.AsString); ok {
		sqlStr = sqlStrValue.AsString()
	} else {
		return nil, utils.NewThrow(errors.New("SQL 语句必须是字符串"))
	}

	args := collectBindArgs(ctx, 1)

	// 执行 SQL 语句
	result, err := conn.Exec(sqlStr, args...)
	if err != nil {
		return nil, utils.NewThrowf("执行 SQL 语句失败: %v", err)
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.NewThrowf("获取影响行数失败: %v", err)
	}

	// 获取最后插入的 ID（如果是 INSERT 语句）
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		// 如果不是 INSERT 语句，lastInsertId 可能为 0，这是正常的
		lastInsertId = 0
	}

	// 返回执行结果
	resultObj := data.NewObjectValue()
	resultObj.SetProperty("rowsAffected", data.NewIntValue(int(rowsAffected)))
	resultObj.SetProperty("lastInsertId", data.NewIntValue(int(lastInsertId)))
	resultObj.SetProperty("success", data.NewBoolValue(true))

	return resultObj, nil
}

func (d *DbExecMethod) GetName() string {
	return "exec"
}

func (d *DbExecMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbExecMethod) GetIsStatic() bool {
	return false
}

func (d *DbExecMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (d *DbExecMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}

func (d *DbExecMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
