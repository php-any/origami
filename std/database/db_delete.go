package database

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbDeleteMethod struct {
	source *db
}

func (d *DbDeleteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	// 构建删除语句
	tableName, ctl := d.source.getTableNameWithContext(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 构建 WHERE 子句
	whereClause := ""
	values := make([]interface{}, 0)

	if d.source.where != "" {
		whereClause = " WHERE " + d.source.where
		// 添加 WHERE 参数
		for _, arg := range d.source.whereArgs {
			values = append(values, ConvertValueToGoType(arg))
		}
	}

	// 构建完整的 DELETE 语句
	query := fmt.Sprintf("DELETE FROM %s%s", tableName, whereClause)

	// 执行删除
	result, err := conn.Exec(query, values...)
	if err != nil {
		return nil, utils.NewThrowf("删除失败: %v", err)
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.NewThrowf("获取影响行数失败: %v", err)
	}

	// 返回删除结果
	resultObj := data.NewObjectValue()
	resultObj.SetProperty("rowsAffected", data.NewIntValue(int(rowsAffected)))
	resultObj.SetProperty("success", data.NewBoolValue(true))

	return resultObj, nil
}

func (d *DbDeleteMethod) GetName() string {
	return "delete"
}

func (d *DbDeleteMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbDeleteMethod) GetIsStatic() bool {
	return false
}

func (d *DbDeleteMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (d *DbDeleteMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (d *DbDeleteMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
