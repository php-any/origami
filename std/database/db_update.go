package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbUpdateMethod struct {
	source *db
}

func (d *DbUpdateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	// 获取更新数据
	dataValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少更新数据"))
	}

	// 构建更新语句
	tableName, ctl := d.source.getTableNameWithContext(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 处理更新数据
	var updateData map[string]interface{}
	if classValue, ok := dataValue.(*data.ClassValue); ok {
		// 从类实例获取属性，并处理注解映射
		updateData = make(map[string]interface{})
		properties := classValue.GetProperties()
		classStmt := classValue.Class

		for name, value := range properties {
			if value != nil {
				// 对于 UPDATE 操作，只跳过真正的 null 值，不跳过 0、false、空字符串等
				// 因为用户可能想要将字段更新为这些值
				if _, ok := value.(*data.NullValue); !ok {
					// 获取数据库列名
					columnName := d.getColumnName(classStmt, name)
					updateData[columnName] = ConvertValueToGoType(value)
				}
			}
		}
	} else if objectValue, ok := dataValue.(*data.ObjectValue); ok {
		// 从对象获取属性
		updateData = make(map[string]interface{})
		properties := objectValue.GetProperties()
		for name, value := range properties {
			if value != nil {
				// 对于 UPDATE 操作，只跳过真正的 null 值，不跳过 0、false、空字符串等
				// 因为用户可能想要将字段更新为这些值
				if _, ok := value.(*data.NullValue); !ok {
					updateData[name] = ConvertValueToGoType(value)
				}
			}
		}
	} else {
		return nil, utils.NewThrow(errors.New("更新数据必须是对象或类实例"))
	}

	// 构建 SET 子句
	setClauses := make([]string, 0, len(updateData))
	values := make([]interface{}, 0, len(updateData))

	for column, value := range updateData {
		setClauses = append(setClauses, column+" = ?")
		values = append(values, value)
	}

	// 检查是否有要更新的字段
	if len(setClauses) == 0 {
		return nil, utils.NewThrow(errors.New("没有要更新的字段。UPDATE 语句必须包含 SET 子句。如果需要执行 SQL 表达式（如自增、NOW() 等），请使用 exec() 方法执行原生 SQL"))
	}

	// 构建 WHERE 子句
	whereClause := ""
	if d.source.where != "" {
		whereClause = " WHERE " + d.source.where
		// 添加 WHERE 参数
		for _, arg := range d.source.whereArgs {
			values = append(values, ConvertValueToGoType(arg))
		}
	}

	// 构建完整的 UPDATE 语句
	query := fmt.Sprintf("UPDATE %s SET %s%s",
		tableName,
		strings.Join(setClauses, ", "),
		whereClause)

	// 执行更新
	result, err := conn.Exec(query, values...)
	if err != nil {
		return nil, utils.NewThrowf("更新失败: %v; sql(%v", err, query)
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.NewThrowf("获取影响行数失败: %v", err)
	}

	// 返回更新结果
	resultObj := data.NewObjectValue()
	resultObj.SetProperty("rowsAffected", data.NewIntValue(int(rowsAffected)))
	resultObj.SetProperty("success", data.NewBoolValue(true))

	return resultObj, nil
}

func (d *DbUpdateMethod) GetName() string {
	return "update"
}

func (d *DbUpdateMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbUpdateMethod) GetIsStatic() bool {
	return false
}

func (d *DbUpdateMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("data", 0),
	}
}

func (d *DbUpdateMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("data", 0, data.NewBaseType("object")),
	}
}

func (d *DbUpdateMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// getColumnName 获取数据库列名，支持注解映射
func (d *DbUpdateMethod) getColumnName(classStmt data.ClassStmt, propertyName string) string {
	return getColumnName(classStmt, propertyName)
}

// isNullValue 检查值是否为空
func (d *DbUpdateMethod) isNullValue(val data.Value) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case *data.NullValue:
		return true
	case *data.StringValue:
		return v.Value == ""
	case *data.IntValue:
		return v.Value == 0
	case *data.FloatValue:
		return v.Value == 0.0
	case *data.BoolValue:
		return !v.Value
	default:
		return false
	}
}
