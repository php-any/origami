package pdo

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"sync"

	// 如果需要 MySQL 支持，请在 go.mod 中添加 github.com/go-sql-driver/mysql
	// 如果需要 SQLite 支持，请在 go.mod 中添加 modernc.org/sqlite
	// 目前使用纯 Go 内置 database/sql，驱动需用户自行注册

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// -------------------------------------------------------------------
// PDO 内部状态
// -------------------------------------------------------------------

// pdoState：一个 PDO 实例从共享连接池借出一条连接。
// - db 指向按 DSN 共享的 *sql.DB 连接池（可服务多个 PDO）
// - conn 是本实例持有的那一条连接；close 时归还到池，不关闭池
// - 同实例多 goroutine 访问由 mu 串行化（*sql.Conn 非并发安全）
// - 多线程：每个线程/请求各自 new PDO，从池并行借出多条连接
type pdoState struct {
	mu            sync.Mutex
	db            *sql.DB // 共享池引用，close 时不 Close(db)
	conn          *sql.Conn
	tx            *sql.Tx
	driverName    string
	serverVersion string
	errMode       int // PDO::ERRMODE_*
	lastError     string
	lastSQLState  string
	lastInsertID  int64
	closed        bool
}

func (s *pdoState) getErrMode() int {
	if s == nil {
		return PDO_ERRMODE_EXCEPTION
	}
	return s.errMode
}

func (s *pdoState) ctx() context.Context {
	return context.Background()
}

func (s *pdoState) ensureOpen() error {
	if s == nil || s.closed || s.conn == nil {
		return fmt.Errorf("PDO connection is closed")
	}
	return nil
}

func (s *pdoState) queryLocked(query string, args ...any) (*sql.Rows, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	if s.tx != nil {
		return s.tx.QueryContext(s.ctx(), query, args...)
	}
	return s.conn.QueryContext(s.ctx(), query, args...)
}

func (s *pdoState) execLocked(query string, args ...any) (sql.Result, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	if s.tx != nil {
		return s.tx.ExecContext(s.ctx(), query, args...)
	}
	return s.conn.ExecContext(s.ctx(), query, args...)
}

// release 归还借出的连接到共享池。不关闭 *sql.DB。
func (s *pdoState) release() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	if s.tx != nil {
		_ = s.tx.Rollback()
		s.tx = nil
	}
	if s.conn != nil {
		_ = s.conn.Close() // database/sql：Close 把连接放回池
		s.conn = nil
	}
	s.db = nil
}

func pdoStateFinalizer(s *pdoState) {
	s.release()
}

// -------------------------------------------------------------------
// PDOClass - ClassStmt 实现（供 vm.AddClass 注册）
// -------------------------------------------------------------------

type PDOClass struct {
	node.Node
}

func (c *PDOClass) GetName() string                            { return "PDO" }
func (c *PDOClass) GetExtend() *string                         { return nil }
func (c *PDOClass) GetImplements() []string                    { return nil }
func (c *PDOClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *PDOClass) GetPropertyList() []data.Property           { return nil }

func (c *PDOClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *PDOClass) GetConstruct() data.Method {
	return &pdoConstructMethod{}
}

func (c *PDOClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return &pdoConstructMethod{}, true
	case "query":
		return &pdoQueryMethod{}, true
	case "exec":
		return &pdoExecMethod{}, true
	case "prepare":
		return &pdoPrepareMethod{}, true
	case "beginTransaction":
		return &pdoBeginTransactionMethod{}, true
	case "commit":
		return &pdoCommitMethod{}, true
	case "rollBack", "rollback":
		return &pdoRollBackMethod{}, true
	case "lastInsertId":
		return &pdoLastInsertIdMethod{}, true
	case "quote":
		return &pdoQuoteMethod{}, true
	case "setAttribute":
		return &pdoSetAttributeMethod{}, true
	case "getAttribute":
		return &pdoGetAttributeMethod{}, true
	case "errorCode":
		return &pdoErrorCodeMethod{}, true
	case "errorInfo":
		return &pdoErrorInfoMethod{}, true
	case "inTransaction":
		return &pdoInTransactionMethod{}, true
	case "close":
		return &pdoCloseMethod{}, true
	case "getAvailableDrivers":
		return &pdoGetAvailableDriversMethod{}, true
	}
	return nil, false
}

func (c *PDOClass) GetMethods() []data.Method { return nil }

func (c *PDOClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "connect":
		return &pdoConnectMethod{}, true
	case "getAvailableDrivers":
		return &pdoGetAvailableDriversMethod{}, true
	default:
		return nil, false
	}
}

// GetStaticProperty 实现 PDO 类常量（PHP 里用 PDO::FETCH_ASSOC 等方式访问）
func (c *PDOClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "ATTR_AUTOCOMMIT":
		return data.NewIntValue(PDO_ATTR_AUTOCOMMIT), true
	case "ATTR_ERRMODE":
		return data.NewIntValue(PDO_ATTR_ERRMODE), true
	case "ATTR_PERSISTENT":
		return data.NewIntValue(PDO_ATTR_PERSISTENT), true
	case "ATTR_DEFAULT_FETCH_MODE":
		return data.NewIntValue(PDO_ATTR_DEFAULT_FETCH_MODE), true
	case "ATTR_EMULATE_PREPARES":
		return data.NewIntValue(PDO_ATTR_EMULATE_PREPARES), true
	case "ATTR_DRIVER_NAME":
		return data.NewIntValue(PDO_ATTR_DRIVER_NAME), true
	case "ATTR_CASE":
		return data.NewIntValue(PDO_ATTR_CASE), true
	case "ATTR_ORACLE_NULLS":
		return data.NewIntValue(PDO_ATTR_ORACLE_NULLS), true
	case "ATTR_STRINGIFY_FETCHES":
		return data.NewIntValue(PDO_ATTR_STRINGIFY_FETCHES), true
	case "ATTR_STATEMENT_CLASS":
		return data.NewIntValue(PDO_ATTR_STATEMENT_CLASS), true
	case "ATTR_TIMEOUT":
		return data.NewIntValue(PDO_ATTR_TIMEOUT), true
	case "ATTR_SERVER_VERSION":
		return data.NewIntValue(PDO_ATTR_SERVER_VERSION), true
	case "ATTR_CLIENT_VERSION":
		return data.NewIntValue(PDO_ATTR_CLIENT_VERSION), true
	case "ATTR_SERVER_INFO":
		return data.NewIntValue(PDO_ATTR_SERVER_INFO), true
	case "ATTR_CONNECTION_STATUS":
		return data.NewIntValue(PDO_ATTR_CONNECTION_STATUS), true
	case "ERRMODE_SILENT":
		return data.NewIntValue(PDO_ERRMODE_SILENT), true
	case "ERRMODE_WARNING":
		return data.NewIntValue(PDO_ERRMODE_WARNING), true
	case "ERRMODE_EXCEPTION":
		return data.NewIntValue(PDO_ERRMODE_EXCEPTION), true
	case "CASE_NATURAL":
		return data.NewIntValue(PDO_CASE_NATURAL), true
	case "CASE_UPPER":
		return data.NewIntValue(PDO_CASE_UPPER), true
	case "CASE_LOWER":
		return data.NewIntValue(PDO_CASE_LOWER), true
	case "NULL_NATURAL":
		return data.NewIntValue(PDO_NULL_NATURAL), true
	case "NULL_EMPTY_STRING":
		return data.NewIntValue(PDO_NULL_EMPTY_STRING), true
	case "NULL_TO_STRING":
		return data.NewIntValue(PDO_NULL_TO_STRING), true
	case "FETCH_DEFAULT":
		return data.NewIntValue(PDO_FETCH_DEFAULT), true
	case "FETCH_LAZY":
		return data.NewIntValue(PDO_FETCH_LAZY), true
	case "FETCH_ASSOC":
		return data.NewIntValue(PDO_FETCH_ASSOC), true
	case "FETCH_NUM":
		return data.NewIntValue(PDO_FETCH_NUM), true
	case "FETCH_BOTH":
		return data.NewIntValue(PDO_FETCH_BOTH), true
	case "FETCH_OBJ":
		return data.NewIntValue(PDO_FETCH_OBJ), true
	case "FETCH_BOUND":
		return data.NewIntValue(PDO_FETCH_BOUND), true
	case "FETCH_COLUMN":
		return data.NewIntValue(PDO_FETCH_COLUMN), true
	case "FETCH_CLASS":
		return data.NewIntValue(PDO_FETCH_CLASS), true
	case "FETCH_INTO":
		return data.NewIntValue(PDO_FETCH_INTO), true
	case "FETCH_FUNC":
		return data.NewIntValue(PDO_FETCH_FUNC), true
	case "FETCH_GROUP":
		return data.NewIntValue(PDO_FETCH_GROUP), true
	case "FETCH_UNIQUE":
		return data.NewIntValue(PDO_FETCH_UNIQUE), true
	case "FETCH_KEY_PAIR":
		return data.NewIntValue(PDO_FETCH_KEY_PAIR), true
	case "FETCH_NAMED":
		return data.NewIntValue(PDO_FETCH_NAMED), true
	case "FETCH_PROPS_LATE":
		return data.NewIntValue(PDO_FETCH_PROPS_LATE), true
	case "PARAM_BOOL":
		return data.NewIntValue(PDO_PARAM_BOOL), true
	case "PARAM_NULL":
		return data.NewIntValue(PDO_PARAM_NULL), true
	case "PARAM_INT":
		return data.NewIntValue(PDO_PARAM_INT), true
	case "PARAM_STR":
		return data.NewIntValue(PDO_PARAM_STR), true
	case "PARAM_LOB":
		return data.NewIntValue(PDO_PARAM_LOB), true
	case "PARAM_INPUT_OUTPUT":
		return data.NewIntValue(PDO_PARAM_INPUT_OUTPUT), true
	case "CURSOR_FWDONLY":
		return data.NewIntValue(PDO_CURSOR_FWDONLY), true
	case "CURSOR_SCROLL":
		return data.NewIntValue(PDO_CURSOR_SCROLL), true
	// MySQL 驱动专属属性
	case "MYSQL_ATTR_USE_BUFFERED_QUERY":
		return data.NewIntValue(PDO_MYSQL_ATTR_USE_BUFFERED_QUERY), true
	case "MYSQL_ATTR_LOCAL_INFILE":
		return data.NewIntValue(PDO_MYSQL_ATTR_LOCAL_INFILE), true
	case "MYSQL_ATTR_INIT_COMMAND":
		return data.NewIntValue(PDO_MYSQL_ATTR_INIT_COMMAND), true
	case "MYSQL_ATTR_COMPRESS":
		return data.NewIntValue(PDO_MYSQL_ATTR_COMPRESS), true
	case "MYSQL_ATTR_DIRECT_QUERY":
		return data.NewIntValue(PDO_MYSQL_ATTR_DIRECT_QUERY), true
	case "MYSQL_ATTR_FOUND_ROWS":
		return data.NewIntValue(PDO_MYSQL_ATTR_FOUND_ROWS), true
	case "MYSQL_ATTR_IGNORE_SPACE":
		return data.NewIntValue(PDO_MYSQL_ATTR_IGNORE_SPACE), true
	case "MYSQL_ATTR_MAX_BUFFER_SIZE":
		return data.NewIntValue(PDO_MYSQL_ATTR_MAX_BUFFER_SIZE), true
	case "MYSQL_ATTR_READ_DEFAULT_FILE":
		return data.NewIntValue(PDO_MYSQL_ATTR_READ_DEFAULT_FILE), true
	case "MYSQL_ATTR_READ_DEFAULT_GROUP":
		return data.NewIntValue(PDO_MYSQL_ATTR_READ_DEFAULT_GROUP), true
	case "MYSQL_ATTR_CONNECT_TIMEOUT":
		return data.NewIntValue(PDO_MYSQL_ATTR_CONNECT_TIMEOUT), true
	case "MYSQL_ATTR_SERVER_PUBLIC_KEY":
		return data.NewIntValue(PDO_MYSQL_ATTR_SERVER_PUBLIC_KEY), true
	case "MYSQL_ATTR_MULTI_STATEMENTS":
		return data.NewIntValue(PDO_MYSQL_ATTR_MULTI_STATEMENTS), true
	case "MYSQL_ATTR_SSL_CA":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CA), true
	case "MYSQL_ATTR_SSL_CAPATH":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CAPATH), true
	case "MYSQL_ATTR_SSL_CERT":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CERT), true
	case "MYSQL_ATTR_SSL_KEY":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_KEY), true
	case "MYSQL_ATTR_SSL_CIPHER":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CIPHER), true
	case "MYSQL_ATTR_SSL_VERIFY_SERVER_CERT":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_VERIFY_SERVER_CERT), true
	case "MYSQL_ATTR_LOCAL_INFILE_DIRECTORY":
		return data.NewIntValue(PDO_MYSQL_ATTR_LOCAL_INFILE_DIRECTORY), true
	}
	return nil, false
}

// -------------------------------------------------------------------
// __construct(string $dsn, string $username='', string $password='', array $options=[])
// -------------------------------------------------------------------

type pdoConstructMethod struct{}

func (m *pdoConstructMethod) GetName() string            { return token.ConstructName }
func (m *pdoConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoConstructMethod) GetIsStatic() bool          { return false }
func (m *pdoConstructMethod) GetReturnType() data.Types  { return nil }

func (m *pdoConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "dsn", 0, nil, nil),
		node.NewParameter(nil, "username", 1, node.NewStringLiteral(nil, ""), nil),
		node.NewParameter(nil, "password", 2, node.NewStringLiteral(nil, ""), nil),
		node.NewParameter(nil, "options", 3, node.NewNullLiteral(nil), nil),
	}
}

func (m *pdoConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "dsn", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "username", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "password", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 3, data.NewBaseType("array")),
	}
}

func (m *pdoConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	dsnVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("PDO::__construct() expects at least 1 argument"))
	}
	dsn := dsnVal.AsString()

	username := ""
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			username = v.AsString()
		}
	}
	password := ""
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			password = v.AsString()
		}
	}

	// 解析 DSN：driver:...
	colonIdx := strings.Index(dsn, ":")
	if colonIdx < 0 {
		return nil, pdoException(fmt.Sprintf("invalid data source name: %s", dsn), ctx)
	}
	driver := dsn[:colonIdx]
	rest := dsn[colonIdx+1:]

	// 将 PHP DSN 转换为 Go driver DSN
	goDSN, err := convertDSN(driver, rest, username, password)
	if err != nil {
		return nil, pdoException(err.Error(), ctx)
	}

	goDriver := phpDriverToGo(driver)
	db, err2 := getSharedDB(goDriver, goDSN)
	if err2 != nil {
		return nil, pdoException(err2.Error(), ctx)
	}
	// 从共享池借出一条连接；本 PDO 生命周期内所有 SQL 都走这条连接。
	conn, connErr := db.Conn(context.Background())
	if connErr != nil {
		return nil, pdoException(connErr.Error(), ctx)
	}

	cv := pdoGetClassValue(ctx)
	if cv == nil {
		_ = conn.Close() // 归还到池
		return nil, data.NewErrorThrow(nil, fmt.Errorf("PDO::__construct() missing instance context"))
	}
	state := &pdoState{
		db:            db,
		conn:          conn,
		driverName:    driver,
		serverVersion: queryServerVersionConn(conn, driver),
		errMode:       PDO_ERRMODE_EXCEPTION,
	}
	runtime.SetFinalizer(state, pdoStateFinalizer)
	cv.ObjectValue.SetProperty("__pdo_state__", &pdoStateValue{state: state})

	return nil, nil
}

type pdoConnectMethod struct{}

func (m *pdoConnectMethod) GetName() string            { return "connect" }
func (m *pdoConnectMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoConnectMethod) GetIsStatic() bool          { return true }
func (m *pdoConnectMethod) GetReturnType() data.Types  { return nil }
func (m *pdoConnectMethod) GetParams() []data.GetValue {
	return (&pdoConstructMethod{}).GetParams()
}
func (m *pdoConnectMethod) GetVariables() []data.Variable {
	return (&pdoConstructMethod{}).GetVariables()
}
func (m *pdoConnectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if _, ctl := (&pdoConstructMethod{}).Call(ctx); ctl != nil {
		return nil, ctl
	}
	cv := pdoGetClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("PDO::connect() missing class context"))
	}
	return cv, nil
}

func queryServerVersionConn(conn *sql.Conn, driver string) string {
	var query string
	switch strings.ToLower(driver) {
	case "mysql":
		query = "SELECT VERSION()"
	case "sqlite":
		query = "SELECT sqlite_version()"
	case "pgsql", "postgres", "postgresql":
		query = "SHOW server_version"
	default:
		return ""
	}

	var version string
	if err := conn.QueryRowContext(context.Background(), query).Scan(&version); err != nil {
		return ""
	}
	return version
}

// -------------------------------------------------------------------
// DSN 转换辅助
// -------------------------------------------------------------------

// phpDriverToGo 把 PHP DSN 驱动名转换为 Go sql driver 名
func phpDriverToGo(phpDriver string) string {
	switch strings.ToLower(phpDriver) {
	case "mysql":
		return "mysql"
	case "sqlite":
		return "sqlite"
	case "pgsql", "postgresql":
		return "postgres"
	default:
		return phpDriver
	}
}

// convertDSN 将 PHP PDO DSN 右侧部分 + username/password 转为 Go DSN
func convertDSN(driver, rest, username, password string) (string, error) {
	switch strings.ToLower(driver) {
	case "mysql":
		// PHP: host=127.0.0.1;port=3306;dbname=testdb;charset=utf8
		// Go:  user:pass@tcp(host:port)/dbname?charset=utf8
		params := parseSemicolonParams(rest)
		host := params["host"]
		if host == "" {
			host = "127.0.0.1"
		}
		port := params["port"]
		if port == "" {
			port = "3306"
		}
		dbname := params["dbname"]
		charset := params["charset"]
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
		if charset != "" {
			dsn += "?charset=" + charset + "&parseTime=True&loc=Local"
		} else {
			dsn += "?parseTime=True&loc=Local"
		}
		return dsn, nil

	case "sqlite":
		// PHP: /path/to/db.sqlite  or  :memory:
		if rest == "" || rest == ":memory:" {
			return ":memory:", nil
		}
		return rest, nil

	case "pgsql", "postgresql":
		// PHP: host=localhost;port=5432;dbname=testdb
		params := parseSemicolonParams(rest)
		host := params["host"]
		if host == "" {
			host = "localhost"
		}
		port := params["port"]
		if port == "" {
			port = "5432"
		}
		dbname := params["dbname"]
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, dbname), nil

	default:
		// 透传
		return rest, nil
	}
}

func parseSemicolonParams(s string) map[string]string {
	result := make(map[string]string)
	for _, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "="); idx > 0 {
			result[strings.TrimSpace(part[:idx])] = strings.TrimSpace(part[idx+1:])
		}
	}
	return result
}

// -------------------------------------------------------------------
// query(string $sql, ?int $fetchMode=null, ...$fetchModeArgs): PDOStatement|false
// -------------------------------------------------------------------

type pdoQueryMethod struct{}

func (m *pdoQueryMethod) GetName() string            { return "query" }
func (m *pdoQueryMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoQueryMethod) GetIsStatic() bool          { return false }
func (m *pdoQueryMethod) GetReturnType() data.Types  { return nil }
func (m *pdoQueryMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "sql", 0, nil, nil)}
}
func (m *pdoQueryMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "sql", 0, data.NewBaseType("string"))}
}

func (m *pdoQueryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	sqlVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	query := sqlVal.AsString()

	state.mu.Lock()
	rows, err := state.queryLocked(query)
	if err != nil {
		state.lastError = err.Error()
		state.mu.Unlock()
		if state.getErrMode() == PDO_ERRMODE_EXCEPTION {
			return nil, pdoException(err.Error(), ctx)
		}
		return data.NewBoolValue(false), nil
	}
	// 在持锁期间缓冲并关闭 Rows，释放本 PDO 借出连接上的游标
	stmtClass := newPDOStatementClass(rows, state)
	state.mu.Unlock()
	return data.NewClassValue(stmtClass, ctx), nil
}

// -------------------------------------------------------------------
// exec(string $sql): int|false
// -------------------------------------------------------------------

type pdoExecMethod struct{}

func (m *pdoExecMethod) GetName() string            { return "exec" }
func (m *pdoExecMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoExecMethod) GetIsStatic() bool          { return false }
func (m *pdoExecMethod) GetReturnType() data.Types  { return nil }
func (m *pdoExecMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "sql", 0, nil, nil)}
}
func (m *pdoExecMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "sql", 0, data.NewBaseType("string"))}
}

func (m *pdoExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	sqlVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	state.mu.Lock()
	result, err := state.execLocked(sqlVal.AsString())
	if err != nil {
		state.lastError = err.Error()
		state.mu.Unlock()
		if state.getErrMode() == PDO_ERRMODE_EXCEPTION {
			return nil, pdoException(err.Error(), ctx)
		}
		return data.NewBoolValue(false), nil
	}
	affected, _ := result.RowsAffected()
	if lastID, err2 := result.LastInsertId(); err2 == nil {
		state.lastInsertID = lastID
	}
	state.mu.Unlock()
	return data.NewIntValue(int(affected)), nil
}

// -------------------------------------------------------------------
// prepare(string $sql, array $options=[]): PDOStatement|false
// -------------------------------------------------------------------

type pdoPrepareMethod struct{}

func (m *pdoPrepareMethod) GetName() string            { return "prepare" }
func (m *pdoPrepareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoPrepareMethod) GetIsStatic() bool          { return false }
func (m *pdoPrepareMethod) GetReturnType() data.Types  { return nil }
func (m *pdoPrepareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "sql", 0, nil, nil),
		node.NewParameter(nil, "options", 1, node.NewNullLiteral(nil), nil),
	}
}
func (m *pdoPrepareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "sql", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 1, data.NewBaseType("array")),
	}
}

func (m *pdoPrepareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	sqlVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 只保存 SQL：execute 时在当前 session（conn 或 tx）上执行，
	// 保证 beginTransaction / exec("BEGIN") 后仍落在同一条独占连接上。
	stmtClass := newPDOStatementFromSQL(state, sqlVal.AsString())
	return data.NewClassValue(stmtClass, ctx), nil
}

// -------------------------------------------------------------------
// beginTransaction() / commit() / rollBack()
// -------------------------------------------------------------------

type pdoBeginTransactionMethod struct{}

func (m *pdoBeginTransactionMethod) GetName() string               { return "beginTransaction" }
func (m *pdoBeginTransactionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoBeginTransactionMethod) GetIsStatic() bool             { return false }
func (m *pdoBeginTransactionMethod) GetReturnType() data.Types     { return nil }
func (m *pdoBeginTransactionMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoBeginTransactionMethod) GetVariables() []data.Variable { return nil }
func (m *pdoBeginTransactionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.tx != nil {
		if state.getErrMode() == PDO_ERRMODE_EXCEPTION {
			return nil, pdoException("There is already an active transaction", ctx)
		}
		return data.NewBoolValue(false), nil
	}
	tx, err := state.conn.BeginTx(state.ctx(), nil)
	if err != nil {
		if state.getErrMode() == PDO_ERRMODE_EXCEPTION {
			return nil, pdoException(err.Error(), ctx)
		}
		return data.NewBoolValue(false), nil
	}
	state.tx = tx
	return data.NewBoolValue(true), nil
}

type pdoCommitMethod struct{}

func (m *pdoCommitMethod) GetName() string               { return "commit" }
func (m *pdoCommitMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoCommitMethod) GetIsStatic() bool             { return false }
func (m *pdoCommitMethod) GetReturnType() data.Types     { return nil }
func (m *pdoCommitMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoCommitMethod) GetVariables() []data.Variable { return nil }
func (m *pdoCommitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.tx == nil {
		return data.NewBoolValue(false), nil
	}
	if err := state.tx.Commit(); err != nil {
		return data.NewBoolValue(false), nil
	}
	state.tx = nil
	return data.NewBoolValue(true), nil
}

type pdoRollBackMethod struct{}

func (m *pdoRollBackMethod) GetName() string               { return "rollBack" }
func (m *pdoRollBackMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoRollBackMethod) GetIsStatic() bool             { return false }
func (m *pdoRollBackMethod) GetReturnType() data.Types     { return nil }
func (m *pdoRollBackMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoRollBackMethod) GetVariables() []data.Variable { return nil }
func (m *pdoRollBackMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.tx == nil {
		return data.NewBoolValue(false), nil
	}
	if err := state.tx.Rollback(); err != nil {
		return data.NewBoolValue(false), nil
	}
	state.tx = nil
	return data.NewBoolValue(true), nil
}

// -------------------------------------------------------------------
// lastInsertId(?string $name=null): string|false
// -------------------------------------------------------------------

type pdoLastInsertIdMethod struct{}

func (m *pdoLastInsertIdMethod) GetName() string            { return "lastInsertId" }
func (m *pdoLastInsertIdMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoLastInsertIdMethod) GetIsStatic() bool          { return false }
func (m *pdoLastInsertIdMethod) GetReturnType() data.Types  { return nil }
func (m *pdoLastInsertIdMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "name", 0, node.NewNullLiteral(nil), nil)}
}
func (m *pdoLastInsertIdMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "name", 0, data.NewBaseType("string"))}
}
func (m *pdoLastInsertIdMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	// lastInsertId 需要通过 exec 结果获取，这里从 state 缓存中读取
	return data.NewStringValue(fmt.Sprintf("%d", state.lastInsertID)), nil
}

// -------------------------------------------------------------------
// quote(string $value, int $type=PDO::PARAM_STR): string|false
// -------------------------------------------------------------------

type pdoQuoteMethod struct{}

func (m *pdoQuoteMethod) GetName() string            { return "quote" }
func (m *pdoQuoteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoQuoteMethod) GetIsStatic() bool          { return false }
func (m *pdoQuoteMethod) GetReturnType() data.Types  { return nil }
func (m *pdoQuoteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "type", 1, node.NewIntLiteral(nil, "2"), nil),
	}
}
func (m *pdoQuoteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "type", 1, data.NewBaseType("int")),
	}
}
func (m *pdoQuoteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	// 简单转义单引号
	s = strings.ReplaceAll(s, "'", "\\'")
	return data.NewStringValue("'" + s + "'"), nil
}

// -------------------------------------------------------------------
// setAttribute(int $attr, mixed $value): bool
// -------------------------------------------------------------------

type pdoSetAttributeMethod struct{}

func (m *pdoSetAttributeMethod) GetName() string            { return "setAttribute" }
func (m *pdoSetAttributeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoSetAttributeMethod) GetIsStatic() bool          { return false }
func (m *pdoSetAttributeMethod) GetReturnType() data.Types  { return nil }
func (m *pdoSetAttributeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "attr", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (m *pdoSetAttributeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "attr", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "value", 1, data.NewBaseType("mixed")),
	}
}
func (m *pdoSetAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	attrVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	valVal, ok2 := ctx.GetIndexValue(1)
	if !ok2 {
		return data.NewBoolValue(false), nil
	}

	attr := 0
	if ai, ok := attrVal.(interface{ AsInt() (int, error) }); ok {
		v, err := ai.AsInt()
		if err == nil {
			attr = v
		}
	}

	switch attr {
	case PDO_ATTR_ERRMODE:
		if vi, ok := valVal.(interface{ AsInt() (int, error) }); ok {
			v, err := vi.AsInt()
			if err == nil {
				state.errMode = v
			}
		}
	}
	return data.NewBoolValue(true), nil
}

// -------------------------------------------------------------------
// getAttribute(int $attr): mixed
// -------------------------------------------------------------------

type pdoGetAttributeMethod struct{}

func (m *pdoGetAttributeMethod) GetName() string            { return "getAttribute" }
func (m *pdoGetAttributeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *pdoGetAttributeMethod) GetIsStatic() bool          { return false }
func (m *pdoGetAttributeMethod) GetReturnType() data.Types  { return nil }
func (m *pdoGetAttributeMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "attr", 0, nil, nil)}
}
func (m *pdoGetAttributeMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "attr", 0, data.NewBaseType("int"))}
}
func (m *pdoGetAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	attrVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewNullValue(), nil
	}
	attr := 0
	if ai, ok := attrVal.(interface{ AsInt() (int, error) }); ok {
		v, err := ai.AsInt()
		if err == nil {
			attr = v
		}
	}
	switch attr {
	case PDO_ATTR_ERRMODE:
		return data.NewIntValue(state.errMode), nil
	case PDO_ATTR_DRIVER_NAME:
		return data.NewStringValue(state.driverName), nil
	case PDO_ATTR_SERVER_VERSION:
		return data.NewStringValue(state.serverVersion), nil
	}
	return data.NewNullValue(), nil
}

// -------------------------------------------------------------------
// errorCode(): string|null
// -------------------------------------------------------------------

type pdoErrorCodeMethod struct{}

func (m *pdoErrorCodeMethod) GetName() string               { return "errorCode" }
func (m *pdoErrorCodeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoErrorCodeMethod) GetIsStatic() bool             { return false }
func (m *pdoErrorCodeMethod) GetReturnType() data.Types     { return nil }
func (m *pdoErrorCodeMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoErrorCodeMethod) GetVariables() []data.Variable { return nil }
func (m *pdoErrorCodeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	if state.lastSQLState == "" {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(state.lastSQLState), nil
}

// -------------------------------------------------------------------
// errorInfo(): array
// -------------------------------------------------------------------

type pdoErrorInfoMethod struct{}

func (m *pdoErrorInfoMethod) GetName() string               { return "errorInfo" }
func (m *pdoErrorInfoMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoErrorInfoMethod) GetIsStatic() bool             { return false }
func (m *pdoErrorInfoMethod) GetReturnType() data.Types     { return nil }
func (m *pdoErrorInfoMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoErrorInfoMethod) GetVariables() []data.Variable { return nil }
func (m *pdoErrorInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	return data.NewArrayValue([]data.Value{
		data.NewStringValue(state.lastSQLState),
		data.NewNullValue(),
		data.NewStringValue(state.lastError),
	}), nil
}

// -------------------------------------------------------------------
// inTransaction(): bool
// -------------------------------------------------------------------

type pdoInTransactionMethod struct{}

func (m *pdoInTransactionMethod) GetName() string               { return "inTransaction" }
func (m *pdoInTransactionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoInTransactionMethod) GetIsStatic() bool             { return false }
func (m *pdoInTransactionMethod) GetReturnType() data.Types     { return nil }
func (m *pdoInTransactionMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoInTransactionMethod) GetVariables() []data.Variable { return nil }
func (m *pdoInTransactionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return data.NewBoolValue(state.tx != nil), nil
}

// -------------------------------------------------------------------
// close(): void — 归还连接到共享池
// -------------------------------------------------------------------

type pdoCloseMethod struct{}

func (m *pdoCloseMethod) GetName() string               { return "close" }
func (m *pdoCloseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoCloseMethod) GetIsStatic() bool             { return false }
func (m *pdoCloseMethod) GetReturnType() data.Types     { return nil }
func (m *pdoCloseMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoCloseMethod) GetVariables() []data.Variable { return nil }
func (m *pdoCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state, acl := getPDOState(ctx)
	if acl != nil {
		return nil, acl
	}
	runtime.SetFinalizer(state, nil)
	state.release()
	return nil, nil
}

// -------------------------------------------------------------------
// getAvailableDrivers(): array
// -------------------------------------------------------------------

type pdoGetAvailableDriversMethod struct{}

func (m *pdoGetAvailableDriversMethod) GetName() string               { return "getAvailableDrivers" }
func (m *pdoGetAvailableDriversMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *pdoGetAvailableDriversMethod) GetIsStatic() bool             { return true }
func (m *pdoGetAvailableDriversMethod) GetReturnType() data.Types     { return nil }
func (m *pdoGetAvailableDriversMethod) GetParams() []data.GetValue    { return nil }
func (m *pdoGetAvailableDriversMethod) GetVariables() []data.Variable { return nil }
func (m *pdoGetAvailableDriversMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	drivers := sql.Drivers()
	vals := make([]data.Value, len(drivers))
	for i, d := range drivers {
		vals[i] = data.NewStringValue(d)
	}
	return data.NewArrayValue(vals), nil
}

// -------------------------------------------------------------------
// pdoStateValue — 存储 Go 对象的 Value 包装
// -------------------------------------------------------------------

type pdoStateValue struct {
	state *pdoState
}

func (v *pdoStateValue) GetValue(_ data.Context) (data.GetValue, data.Control) { return v, nil }
func (v *pdoStateValue) AsString() string                                      { return "[PDO]" }
func (v *pdoStateValue) SetValue(_ data.Value)                                 {}

// -------------------------------------------------------------------
// 辅助函数
// -------------------------------------------------------------------

func pdoGetClassValue(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

func getPDOState(ctx data.Context) (*pdoState, data.Control) {
	cv := pdoGetClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("PDO object is not initialized"))
	}
	v, _ := cv.ObjectValue.GetProperty("__pdo_state__")
	if sv, ok := v.(*pdoStateValue); ok {
		return sv.state, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("PDO object is not initialized"))
}

func pdoException(msg string, ctx data.Context) data.Control {
	exClass := NewPDOExceptionClass()
	exClass.m.SetMessage(msg)
	cv := data.NewClassValue(exClass, ctx)
	return data.NewErrorThrowFromClassValue(nil, cv)
}
