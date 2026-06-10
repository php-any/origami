package database

import (
	"github.com/php-any/origami/data"
)

func NewDBClass() *DBClass {
	return (&DBClass{
		bindMethod:           &DbBindMethod{},
		modelMethod:          &DbModelMethod{},
		staticInsertMethod:   &DbStaticInsertMethod{},
		staticSqlMethod:      &DbStaticSqlMethod{},
		staticExecuteMethod:  &DbStaticExecuteMethod{},
		staticToEntityMethod: &DbToEntityMethod{},
	}).Clone(nil).(*DBClass)
}

type DBClass struct {
	bindMethod           data.Method
	modelMethod          data.Method
	staticInsertMethod   data.Method
	staticSqlMethod      data.Method
	staticExecuteMethod  data.Method
	staticToEntityMethod data.Method
	sqlMethod            data.Method
	construct            data.Method
	getMethod            data.Method
	firstMethod          data.Method
	whereMethod          data.Method
	tableMethod          data.Method
	selectMethod         data.Method
	orderByMethod        data.Method
	groupByMethod        data.Method
	limitMethod          data.Method
	offsetMethod         data.Method
	joinMethod           data.Method
	// CRUD 方法
	insertMethod data.Method
	updateMethod data.Method
	deleteMethod data.Method
	// 原生 SQL 方法
	queryMethod data.Method
	execMethod  data.Method
}

func (d *DBClass) Clone(m map[string]data.Types) data.ClassGeneric {
	var source *db
	if m != nil {
		source = newDB(m)
	} else {
		// 如果没有泛型参数，创建一个空的 db 对象
		source = &db{}
	}

	bindMethod := d.bindMethod
	if bindMethod == nil {
		bindMethod = &DbBindMethod{}
	}
	modelMethod := d.modelMethod
	if modelMethod == nil {
		modelMethod = &DbModelMethod{}
	}
	staticInsertMethod := d.staticInsertMethod
	if staticInsertMethod == nil {
		staticInsertMethod = &DbStaticInsertMethod{}
	}
	staticSqlMethod := d.staticSqlMethod
	if staticSqlMethod == nil {
		staticSqlMethod = &DbStaticSqlMethod{}
	}
	staticExecuteMethod := d.staticExecuteMethod
	if staticExecuteMethod == nil {
		staticExecuteMethod = &DbStaticExecuteMethod{}
	}
	staticToEntityMethod := d.staticToEntityMethod
	if staticToEntityMethod == nil {
		staticToEntityMethod = &DbToEntityMethod{}
	}

	return &DBClass{
		bindMethod:           bindMethod,
		modelMethod:          modelMethod,
		staticInsertMethod:   staticInsertMethod,
		staticSqlMethod:      staticSqlMethod,
		staticExecuteMethod:  staticExecuteMethod,
		staticToEntityMethod: staticToEntityMethod,
		sqlMethod:            &DbSqlMethod{source},
		construct:            &DbConstructMethod{source},
		getMethod:            &DbGetMethod{source},
		firstMethod:          &DbFirstMethod{source: source, scanner: nil},
		whereMethod:          &DbWhereMethod{source},
		tableMethod:          &DbTableMethod{source},
		selectMethod:         &DbSelectMethod{source},
		orderByMethod:        &DbOrderByMethod{source},
		groupByMethod:        &DbGroupByMethod{source},
		limitMethod:          &DbLimitMethod{source},
		offsetMethod:         &DbOffsetMethod{source},
		joinMethod:           &DbJoinMethod{source},
		// CRUD 方法
		insertMethod: &DbInsertMethod{source},
		updateMethod: &DbUpdateMethod{source},
		deleteMethod: &DbDeleteMethod{source},
		// 原生 SQL 方法
		queryMethod: &DbQueryMethod{source},
		execMethod:  &DbExecMethod{source},
	}
}

// CloneWithSource 使用现有的 db 对象创建新的 DBClass
func (d *DBClass) CloneWithSource(source *db) *DBClass {
	return &DBClass{
		bindMethod:           d.bindMethod,
		modelMethod:          d.modelMethod,
		staticInsertMethod:   d.staticInsertMethod,
		staticSqlMethod:      d.staticSqlMethod,
		staticExecuteMethod:  d.staticExecuteMethod,
		staticToEntityMethod: d.staticToEntityMethod,
		sqlMethod:            &DbSqlMethod{source},
		construct:            &DbConstructMethod{source},
		getMethod:            &DbGetMethod{source},
		firstMethod:          &DbFirstMethod{source: source, scanner: nil},
		whereMethod:          &DbWhereMethod{source},
		tableMethod:          &DbTableMethod{source},
		selectMethod:         &DbSelectMethod{source},
		orderByMethod:        &DbOrderByMethod{source},
		groupByMethod:        &DbGroupByMethod{source},
		limitMethod:          &DbLimitMethod{source},
		offsetMethod:         &DbOffsetMethod{source},
		joinMethod:           &DbJoinMethod{source},
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
	return "Database\\DB"
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

func (d *DBClass) GetPropertyList() []data.Property {
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
	case "sql":
		return d.sqlMethod, true
	case "exec":
		return d.execMethod, true
	}

	return nil, false
}

func (d *DBClass) GetMethods() []data.Method {
	return []data.Method{
		d.bindMethod,
		d.modelMethod,
		d.staticInsertMethod,
		d.staticSqlMethod,
		d.staticExecuteMethod,
		d.staticToEntityMethod,
		d.construct,
		d.getMethod,
		d.firstMethod,
		d.whereMethod,
		d.tableMethod,
		d.selectMethod,
		d.orderByMethod,
		d.groupByMethod,
		d.limitMethod,
		d.offsetMethod,
		d.joinMethod,
		d.insertMethod,
		d.updateMethod,
		d.deleteMethod,
		d.queryMethod,
		d.sqlMethod,
		d.execMethod,
	}
}

func (d *DBClass) GetConstruct() data.Method {
	return d.construct
}

func (d *DBClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "bind":
		return d.bindMethod, true
	case "model":
		return d.modelMethod, true
	case "insert":
		return d.staticInsertMethod, true
	case "sql":
		return d.staticSqlMethod, true
	case "toEntity":
		return d.staticToEntityMethod, true
	case "execute":
		return d.staticExecuteMethod, true
	}
	return nil, false
}
