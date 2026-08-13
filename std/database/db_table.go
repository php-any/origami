package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbTableMethod struct {
	source *db
}

func (d *DbTableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取表名参数
	if tableName, ok := ctx.GetIndexValue(0); ok {
		if tableNameStr, ok := tableName.(data.AsString); ok {
			// 创建新的 db 对象实例
			newDB := d.source.clone()
			newDB.setTableName(tableNameStr.AsString())

			// 创建新的 DBClass 实例
			newDBClass := (&DBClass{}).CloneWithSource(newDB)

			// 返回新的 DB 类实例
			return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
		}
	}
	return nil, utils.NewThrow(errors.New("表名必须是字符串"))
}

func (d *DbTableMethod) GetName() string {
	return "table"
}

func (d *DbTableMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbTableMethod) GetIsStatic() bool {
	return false
}

func (d *DbTableMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("tableName", 0),
	}
}

func (d *DbTableMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("tableName", 0, data.NewBaseType("string")),
	}
}

func (d *DbTableMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
