package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DbWhereMethod struct {
	source *db
}

func (d *DbWhereMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 sql 参数"))
	}
	sql := a1.(data.AsString).AsString()

	var args []data.Value
	if a2, ok := ctx.GetIndexValue(1); ok {
		if arr, ok := a2.(*data.ArrayValue); ok {
			args = flattenBindArgList(arr.ToValueList())
		} else if val, ok := a2.(data.Value); ok {
			args = []data.Value{val}
		}
	}

	// 创建新的 db 对象实例
	newDB := d.source.clone()
	newDB.where = sql
	newDB.whereArgs = args

	// 创建新的 DBClass 实例
	newDBClass := (&DBClass{}).CloneWithSource(newDB)

	// 返回新的 DB 类实例
	return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
}

func (d *DbWhereMethod) GetName() string {
	return "where"
}

func (d *DbWhereMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbWhereMethod) GetIsStatic() bool {
	return false
}

func (d *DbWhereMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (d *DbWhereMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}

// GetReturnType 返回方法返回类型
func (d *DbWhereMethod) GetReturnType() data.Types {
	return data.Generic{}
}
