package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DbUpdateMethod struct {
	source *db
}

func (d *DbUpdateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, data.NewErrorThrow(nil, errors.New("数据库连接不可用"))
	}

	// 获取更新数据
	dataValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少更新数据"))
	}

	// 构建更新语句
	tableName := d.source.getTableNameWithContext(ctx)

	// 处理更新数据
	var updateData map[string]interface{}
	if classValue, ok := dataValue.(*data.ClassValue); ok {
		// 从类实例获取属性，并处理注解映射
		updateData = make(map[string]interface{})
		properties := classValue.GetProperties()
		classStmt := classValue.Class

		for name, value := range properties {
			if value != nil {
				// 检查值是否为空
				if !d.isNullValue(value) {
					// 获取数据库列名
					columnName := d.getColumnName(classStmt, name)
					updateData[columnName] = d.convertValueToGoType(value)
				}
			}
		}
	} else if objectValue, ok := dataValue.(*data.ObjectValue); ok {
		// 从对象获取属性
		updateData = make(map[string]interface{})
		properties := objectValue.GetProperties()
		for name, value := range properties {
			if value != nil {
				// 检查值是否为空
				if !d.isNullValue(value) {
					updateData[name] = d.convertValueToGoType(value)
				}
			}
		}
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("更新数据必须是对象或类实例"))
	}

	// 构建 SET 子句
	setClauses := make([]string, 0, len(updateData))
	values := make([]interface{}, 0, len(updateData))

	for column, value := range updateData {
		setClauses = append(setClauses, column+" = ?")
		values = append(values, value)
	}

	// 构建 WHERE 子句
	whereClause := ""
	if d.source.where != "" {
		whereClause = " WHERE " + d.source.where
		// 添加 WHERE 参数
		for _, arg := range d.source.whereArgs {
			values = append(values, d.convertValueToGoType(arg))
		}
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("更新操作必须指定 WHERE 条件"))
	}

	// 构建完整的 UPDATE 语句
	query := fmt.Sprintf("UPDATE %s SET %s%s",
		tableName,
		strings.Join(setClauses, ", "),
		whereClause)

	// 执行更新
	result, err := conn.Exec(query, values...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("更新失败: %w", err))
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("获取影响行数失败: %w", err))
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

// convertValueToGoType 将 data.Value 转换为 Go 原生类型
func (d *DbUpdateMethod) convertValueToGoType(val data.Value) interface{} {
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

// getColumnName 获取数据库列名，支持注解映射
func (d *DbUpdateMethod) getColumnName(classStmt data.ClassStmt, propertyName string) string {
	// 获取属性定义
	properties := classStmt.GetPropertyList()
	var property data.Property
	var exists bool
	for _, prop := range properties {
		if prop.GetName() == propertyName {
			property = prop
			exists = true
			break
		}
	}
	if !exists {
		return propertyName
	}

	// 检查是否有注解
	if classProperty, ok := property.(*node.ClassProperty); ok {
		// 遍历属性的注解列表
		for _, annotation := range classProperty.Annotations {
			if annotation == nil {
				continue
			}

			// 检查是否是 Column 注解
			if annotation.Class != nil {
				className := annotation.Class.GetName()
				if className == "database\\annotation\\Column" {
					// 获取注解实例的属性
					annotationProps := annotation.GetProperties()

					// 查找 name 属性（Column 注解的第一个参数）
					if nameValue, exists := annotationProps["name"]; exists {
						if nameStr, ok := nameValue.(data.AsString); ok {
							return nameStr.AsString()
						}
					}

					// 如果 name 属性不存在，尝试其他可能的属性名
					for propName, propValue := range annotationProps {
						if propName != "name" && propValue != nil {
							if nameStr, ok := propValue.(data.AsString); ok {
								// 如果找到字符串值，可能是列名
								return nameStr.AsString()
							}
						}
					}
				}
			}
		}
	}

	// 如果没有注解，使用属性名
	return propertyName
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
