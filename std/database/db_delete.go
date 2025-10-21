package database

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

type DbDeleteMethod struct {
	source *db
}

func (d *DbDeleteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, data.NewErrorThrow(nil, errors.New("数据库连接不可用"))
	}

	// 构建删除语句
	tableName := d.source.getTableNameWithContext(ctx)

	// 构建 WHERE 子句
	whereClause := ""
	values := make([]interface{}, 0)

	if d.source.where != "" {
		whereClause = " WHERE " + d.source.where
		// 添加 WHERE 参数
		for _, arg := range d.source.whereArgs {
			values = append(values, d.convertValueToGoType(arg))
		}
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("删除操作必须指定 WHERE 条件"))
	}

	// 构建完整的 DELETE 语句
	query := fmt.Sprintf("DELETE FROM %s%s", tableName, whereClause)

	// 执行删除
	result, err := conn.Exec(query, values...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("删除失败: %w", err))
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("获取影响行数失败: %w", err))
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

// convertValueToGoType 将 data.Value 转换为 Go 原生类型
func (d *DbDeleteMethod) convertValueToGoType(val data.Value) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case *data.IntValue:
		return v.Value
	case *data.StringValue:
		return v.Value
	case *data.BoolValue:
		return v.Value
	case *data.FloatValue:
		return v.Value
	case *data.NullValue:
		return nil
	case *data.ArrayValue:
		// 对于数组，转换为 []interface{}
		result := make([]interface{}, len(v.Value))
		for i, item := range v.Value {
			result[i] = d.convertValueToGoType(item)
		}
		return result
	case *data.ObjectValue:
		// 对于对象，转换为 map[string]interface{}
		result := make(map[string]interface{})
		for k, item := range v.GetProperties() {
			result[k] = d.convertValueToGoType(item)
		}
		return result
	default:
		// 对于其他类型，尝试转换为字符串
		return v.AsString()
	}
}
