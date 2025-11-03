package database

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DbQueryMethod struct {
	source *db
}

func (d *DbQueryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 执行查询
	rows, err := conn.Query(sqlStr, args...)
	if err != nil {
		return nil, utils.NewThrowf("执行 SQL 查询失败: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return nil, utils.NewThrowf("获取列信息失败: %v", err)
	}

	// 创建结果数组
	var results []data.Value

	// 创建扫描目标
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 扫描所有行
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, utils.NewThrowf("扫描行失败: %v", err)
		}

		// 创建行对象
		rowObj := data.NewObjectValue()
		for i, col := range columns {
			rowObj.SetProperty(col, d.convertToValue(values[i]))
		}

		results = append(results, rowObj)
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		return nil, utils.NewThrowf("遍历行时出错: %v", err)
	}

	// 返回结果数组
	return data.NewArrayValue(results), nil
}

func (d *DbQueryMethod) GetName() string {
	return "query"
}

func (d *DbQueryMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbQueryMethod) GetIsStatic() bool {
	return false
}

func (d *DbQueryMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "sql", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "params", 1, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (d *DbQueryMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "sql", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "params", 1, data.NewBaseType("array")),
	}
}

func (d *DbQueryMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

// convertToValue 将数据库值转换为脚本值
func (d *DbQueryMethod) convertToValue(val interface{}) data.Value {
	if val == nil {
		return data.NewNullValue()
	}

	switch v := val.(type) {
	case int:
		return data.NewIntValue(v)
	case int8:
		return data.NewIntValue(int(v))
	case int16:
		return data.NewIntValue(int(v))
	case int32:
		return data.NewIntValue(int(v))
	case int64:
		// 检查是否超出 int 范围
		if v > int64(^uint(0)>>1) || v < int64(-1<<63) {
			// 如果超出 int 范围，转换为字符串
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case uint:
		if v > uint(^uint(0)>>1) {
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case uint8:
		return data.NewIntValue(int(v))
	case uint16:
		return data.NewIntValue(int(v))
	case uint32:
		return data.NewIntValue(int(v))
	case uint64:
		if v > uint64(^uint(0)>>1) {
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case float32:
		return data.NewFloatValue(float64(v))
	case float64:
		return data.NewFloatValue(v)
	case string:
		return data.NewStringValue(v)
	case []byte:
		return data.NewStringValue(string(v))
	case bool:
		return data.NewBoolValue(v)
	default:
		// 对于其他类型，转换为字符串
		return data.NewStringValue(fmt.Sprintf("%v", v))
	}
}
