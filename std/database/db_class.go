package database

import (
	"github.com/php-any/origami/data"
)

func NewDBClass() *DBClass {
	return (&DBClass{}).Clone(nil).(*DBClass)
}

type DBClass struct {
	construct     data.Method
	getMethod     data.Method
	firstMethod   data.Method
	whereMethod   data.Method
	tableMethod   data.Method
	selectMethod  data.Method
	orderByMethod data.Method
	groupByMethod data.Method
	limitMethod   data.Method
	offsetMethod  data.Method
	joinMethod    data.Method
	// CRUD 方法
	insertMethod data.Method
	updateMethod data.Method
	deleteMethod data.Method
	// 原生 SQL 方法
	queryMethod data.Method
	execMethod  data.Method
}

func (d *DBClass) Clone(m map[string]data.Types) data.ClassGeneric {
	source := newDB(m)

	return &DBClass{
		construct:     &DbConstructMethod{source},
		getMethod:     &DbGetMethod{source},
		firstMethod:   &DbFirstMethod{source: source, scanner: nil},
		whereMethod:   &DbWhereMethod{source},
		tableMethod:   &DbTableMethod{source},
		selectMethod:  &DbSelectMethod{source},
		orderByMethod: &DbOrderByMethod{source},
		groupByMethod: &DbGroupByMethod{source},
		limitMethod:   &DbLimitMethod{source},
		offsetMethod:  &DbOffsetMethod{source},
		joinMethod:    &DbJoinMethod{source},
		// CRUD 方法
		insertMethod: &DbInsertMethod{source},
		updateMethod: &DbUpdateMethod{source},
		deleteMethod: &DbDeleteMethod{source},
		// 原生 SQL 方法
		queryMethod: &DbQueryMethod{source},
		execMethod:  &DbExecMethod{source},
	}
}

func (d *DBClass) GenericList() []data.Types {
	return []data.Types{
		data.Generic{Name: "M"},
	}
}

// GetValue 泛型会提前有一次 Clone
func (d *DBClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(d, ctx), nil
}

func (d *DBClass) GetFrom() data.From {
	return nil
}

func (d *DBClass) GetName() string {
	return "database\\DB"
}

func (d *DBClass) GetExtend() *string {
	return nil
}

func (d *DBClass) GetImplements() []string {
	return nil
}

func (d *DBClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (d *DBClass) GetProperties() map[string]data.Property {
	return nil
}

func (d *DBClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "get":
		return d.getMethod, true
	case "first":
		return d.firstMethod, true
	case "where":
		return d.whereMethod, true
	case "table":
		return d.tableMethod, true
	case "select":
		return d.selectMethod, true
	case "orderBy":
		return d.orderByMethod, true
	case "groupBy":
		return d.groupByMethod, true
	case "limit":
		return d.limitMethod, true
	case "offset":
		return d.offsetMethod, true
	case "join":
		return d.joinMethod, true
	// CRUD 方法
	case "insert":
		return d.insertMethod, true
	case "update":
		return d.updateMethod, true
	case "delete":
		return d.deleteMethod, true
	// 原生 SQL 方法
	case "query":
		return d.queryMethod, true
	case "exec":
		return d.execMethod, true
	}

	return nil, false
}

func (d *DBClass) GetMethods() []data.Method {
	return []data.Method{}
}

func (d *DBClass) GetConstruct() data.Method {
	return d.construct
}
