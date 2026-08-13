package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbGroupByMethod struct {
	source *db
}

func (d *DbGroupByMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取分组参数
	if groupBy, ok := ctx.GetIndexValue(0); ok {
		if groupByStr, ok := groupBy.(data.AsString); ok {
			// 创建新的 db 对象实例
			newDB := d.source.clone()
			newDB.setGroupBy(groupByStr.AsString())

			// 创建新的 DBClass 实例
			newDBClass := (&DBClass{}).CloneWithSource(newDB)

			// 返回新的 DB 类实例
			return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
		}
	}
	return nil, utils.NewThrow(errors.New("分组字段必须是字符串"))
}

func (d *DbGroupByMethod) GetName() string {
	return "groupBy"
}

func (d *DbGroupByMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbGroupByMethod) GetIsStatic() bool {
	return false
}

func (d *DbGroupByMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("groupBy", 0),
	}
}

func (d *DbGroupByMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("groupBy", 0, data.NewBaseType("string")),
	}
}

func (d *DbGroupByMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
