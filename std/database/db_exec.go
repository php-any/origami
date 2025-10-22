package database

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DbExecMethod struct {
	source *db
}

func (d *DbExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, data.NewErrorThrow(nil, errors.New("数据库连接不可用"))
	}

	// 获取 SQL 语句
	sqlValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少 SQL 语句"))
	}

	var sqlStr string
	if sqlStrValue, ok := sqlValue.(data.AsString); ok {
		sqlStr = sqlStrValue.AsString()
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("SQL 语句必须是字符串"))
	}

	// 获取参数
	var args []interface{}
	if paramValue, ok := ctx.GetIndexValue(1); ok {
		if paramArray, ok := paramValue.(*data.ArrayValue); ok {
			args = make([]interface{}, len(paramArray.Value))
			for i, param := range paramArray.Value {
				args[i] = ConvertValueToGoType(param)
			}
		} else {
			// 单个参数
			args = []interface{}{ConvertValueToGoType(paramValue)}
		}
	}

	// 执行 SQL 语句
	result, err := conn.Exec(sqlStr, args...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("执行 SQL 语句失败: %w", err))
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("获取影响行数失败: %w", err))
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
		node.NewParameter(nil, "sql", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "params", 1, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (d *DbExecMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "sql", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "params", 1, data.NewBaseType("array")),
	}
}

func (d *DbExecMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
