package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbOrderByMethod struct {
	source *db
}

func (d *DbOrderByMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取排序参数
	if orderBy, ok := ctx.GetIndexValue(0); ok {
		if orderByStr, ok := orderBy.(data.AsString); ok {
			// 创建新的 db 对象实例
			newDB := d.source.clone()
			newDB.setOrderBy(orderByStr.AsString())

			// 创建新的 DBClass 实例
			newDBClass := (&DBClass{}).CloneWithSource(newDB)

			// 返回新的 DB 类实例
			return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
		}
	}
	return nil, data.NewErrorThrow(nil, errors.New("排序字段必须是字符串"))
}

func (d *DbOrderByMethod) GetName() string {
	return "orderBy"
}

func (d *DbOrderByMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbOrderByMethod) GetIsStatic() bool {
	return false
}

func (d *DbOrderByMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("orderBy", 0),
	}
}

func (d *DbOrderByMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("orderBy", 0, data.NewBaseType("string")),
	}
}

func (d *DbOrderByMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
