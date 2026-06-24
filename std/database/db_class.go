package database

import (
	"github.com/php-any/origami/data"
)

func NewDBClass() *DBClass {
	return (&DBClass{
		modelMethod:            &DbModelMethod{},
		staticConnectionMethod: &DbStaticConnectionMethod{},
		staticInsertMethod:     &DbStaticInsertMethod{},
		staticQueryMethod:      &DbStaticQueryMethod{},
		staticExecuteMethod:    &DbStaticExecuteMethod{},
		staticToEntityMethod:   &DbToEntityMethod{},
	}).Clone(nil).(*DBClass)
}

type DBClass struct {
	modelMethod            data.Method
	staticConnectionMethod data.Method
	staticInsertMethod     data.Method
	staticQueryMethod      data.Method
	staticExecuteMethod    data.Method
	staticToEntityMethod   data.Method
	construct              data.Method
	getMethod              data.Method
	firstMethod            data.Method
	whereMethod            data.Method
	tableMethod            data.Method
	connectionMethod       data.Method
	selectMethod           data.Method
	orderByMethod          data.Method
	groupByMethod          data.Method
	limitMethod            data.Method
	offsetMethod           data.Method
	joinMethod             data.Method
	insertMethod           data.Method
	updateMethod           data.Method
	deleteMethod           data.Method
	queryMethod            data.Method
	executeMethod          data.Method
}

func (d *DBClass) Clone(m map[string]data.Types) data.ClassGeneric {
	var source *db
	if m != nil {
		source = newDB(m)
	} else {
		source = &db{}
	}

	modelMethod := d.modelMethod
	if modelMethod == nil {
		modelMethod = &DbModelMethod{}
	}
	staticConnectionMethod := d.staticConnectionMethod
	if staticConnectionMethod == nil {
		staticConnectionMethod = &DbStaticConnectionMethod{}
	}
	staticInsertMethod := d.staticInsertMethod
	if staticInsertMethod == nil {
		staticInsertMethod = &DbStaticInsertMethod{}
	}
	staticQueryMethod := d.staticQueryMethod
	if staticQueryMethod == nil {
		staticQueryMethod = &DbStaticQueryMethod{}
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
		modelMethod:            modelMethod,
		staticConnectionMethod: staticConnectionMethod,
		staticInsertMethod:     staticInsertMethod,
		staticQueryMethod:      staticQueryMethod,
		staticExecuteMethod:    staticExecuteMethod,
		staticToEntityMethod:   staticToEntityMethod,
		construct:              &DbConstructMethod{source},
		getMethod:              &DbGetMethod{source},
		firstMethod:            &DbFirstMethod{source: source, scanner: nil},
		whereMethod:            &DbWhereMethod{source},
		tableMethod:            &DbTableMethod{source},
		connectionMethod:       &DbConnectionMethod{source},
		selectMethod:           &DbSelectMethod{source},
		orderByMethod:          &DbOrderByMethod{source},
		groupByMethod:          &DbGroupByMethod{source},
		limitMethod:            &DbLimitMethod{source},
		offsetMethod:           &DbOffsetMethod{source},
		joinMethod:             &DbJoinMethod{source},
		insertMethod:           &DbInsertMethod{source},
		updateMethod:           &DbUpdateMethod{source},
		deleteMethod:           &DbDeleteMethod{source},
		queryMethod:            &DbQueryMethod{source},
		executeMethod:          &DbExecuteMethod{source},
	}
}

// CloneWithSource 使用现有的 db 对象创建新的 DBClass
func (d *DBClass) CloneWithSource(source *db) *DBClass {
	return &DBClass{
		modelMethod:            d.modelMethod,
		staticConnectionMethod: d.staticConnectionMethod,
		staticInsertMethod:     d.staticInsertMethod,
		staticQueryMethod:      d.staticQueryMethod,
		staticExecuteMethod:    d.staticExecuteMethod,
		staticToEntityMethod:   d.staticToEntityMethod,
		construct:              &DbConstructMethod{source},
		getMethod:              &DbGetMethod{source},
		firstMethod:            &DbFirstMethod{source: source, scanner: nil},
		whereMethod:            &DbWhereMethod{source},
		tableMethod:            &DbTableMethod{source},
		connectionMethod:       &DbConnectionMethod{source},
		selectMethod:           &DbSelectMethod{source},
		orderByMethod:          &DbOrderByMethod{source},
		groupByMethod:          &DbGroupByMethod{source},
		limitMethod:            &DbLimitMethod{source},
		offsetMethod:           &DbOffsetMethod{source},
		joinMethod:             &DbJoinMethod{source},
		insertMethod:           &DbInsertMethod{source},
		updateMethod:           &DbUpdateMethod{source},
		deleteMethod:           &DbDeleteMethod{source},
		queryMethod:            &DbQueryMethod{source},
		executeMethod:          &DbExecuteMethod{source},
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
	case "connection":
		return d.connectionMethod, true
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
	case "insert":
		return d.insertMethod, true
	case "update":
		return d.updateMethod, true
	case "delete":
		return d.deleteMethod, true
	case "query":
		return d.queryMethod, true
	case "execute":
		return d.executeMethod, true
	}

	return nil, false
}

func (d *DBClass) GetMethods() []data.Method {
	return []data.Method{
		d.modelMethod,
		d.staticConnectionMethod,
		d.staticInsertMethod,
		d.staticQueryMethod,
		d.staticExecuteMethod,
		d.staticToEntityMethod,
		d.construct,
		d.getMethod,
		d.firstMethod,
		d.whereMethod,
		d.tableMethod,
		d.connectionMethod,
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
		d.executeMethod,
	}
}

func (d *DBClass) GetConstruct() data.Method {
	return d.construct
}

func (d *DBClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "model":
		return d.modelMethod, true
	case "connection":
		return d.staticConnectionMethod, true
	case "insert":
		return d.staticInsertMethod, true
	case "query":
		return d.staticQueryMethod, true
	case "toEntity":
		return d.staticToEntityMethod, true
	case "execute":
		return d.staticExecuteMethod, true
	}
	return nil, false
}
