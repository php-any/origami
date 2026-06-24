package database

import (
	"database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DbQueryMethod struct {
	source *db
}

func (d *DbQueryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlStr, goArgs, ctl := parseSQLCallArgs(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return d.source.runQuery(ctx, sqlStr, goArgs)
}

func (d *DbQueryMethod) GetName() string            { return "query" }
func (d *DbQueryMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbQueryMethod) GetIsStatic() bool          { return false }
func (d *DbQueryMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}
func (d *DbQueryMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}
func (d *DbQueryMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func parseSQLCallArgs(ctx data.Context) (string, []interface{}, data.Control) {
	sqlValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return "", nil, utils.NewThrow(errors.New("缺少 SQL 语句"))
	}
	sqlStrValue, ok := sqlValue.(data.AsString)
	if !ok {
		return "", nil, utils.NewThrow(errors.New("SQL 语句必须是字符串"))
	}
	return sqlStrValue.AsString(), collectBindArgs(ctx, 1), nil
}

// runQuery 执行 SELECT 类原生 SQL，返回行对象数组。
func (d *db) runQuery(ctx data.Context, sqlStr string, goArgs []interface{}) (data.GetValue, data.Control) {
	conn := d.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	rows, err := conn.Query(sqlStr, goArgs...)
	if err != nil {
		return nil, utils.NewThrowf("执行 SQL 查询失败: %v", err)
	}
	defer rows.Close()

	return scanRowsToObjects(rows)
}

func scanRowsToObjects(rows *sql.Rows) (data.GetValue, data.Control) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, utils.NewThrowf("获取列信息失败: %v", err)
	}

	scanner := NewDatabaseScanner()
	var results []data.Value
	scanValues := make([]interface{}, len(columns))
	scanPtrs := make([]interface{}, len(columns))
	for i := range scanValues {
		scanPtrs[i] = &scanValues[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanPtrs...); err != nil {
			return nil, utils.NewThrowf("扫描行失败: %v", err)
		}
		rowObj := data.NewObjectValue()
		for i, col := range columns {
			rowObj.SetProperty(col, scanner.convertToValue(scanValues[i]))
		}
		results = append(results, rowObj)
	}
	if err := rows.Err(); err != nil {
		return nil, utils.NewThrowf("遍历行时出错: %v", err)
	}
	return data.NewArrayValue(results), nil
}
