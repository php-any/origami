package pdo

import (
	_ "github.com/go-sql-driver/mysql" // 自动注册 MySQL 驱动
	"github.com/php-any/origami/data"
)

// Load 注册 PDO 相关类和常量到 VM
func Load(vm data.VM) {
	// 注册 PDO 主类
	vm.AddClass(&PDOClass{})

	// 注册 PDOException 类（extends RuntimeException）
	vm.AddClass(NewPDOExceptionClass())

	// 注册 PDOStatement 类（空壳，实际由 PDO::query/prepare 返回）
	vm.AddClass(&PDOStatementClass{})

	// 注册 PDO 类常量（以 PDO::ATTR_xxx 形式）
	vm.SetConstant("PDO::ATTR_AUTOCOMMIT", data.NewIntValue(PDO_ATTR_AUTOCOMMIT))
	vm.SetConstant("PDO::ATTR_ERRMODE", data.NewIntValue(PDO_ATTR_ERRMODE))
	vm.SetConstant("PDO::ATTR_PERSISTENT", data.NewIntValue(PDO_ATTR_PERSISTENT))
	vm.SetConstant("PDO::ATTR_DEFAULT_FETCH_MODE", data.NewIntValue(PDO_ATTR_DEFAULT_FETCH_MODE))
	vm.SetConstant("PDO::ATTR_EMULATE_PREPARES", data.NewIntValue(PDO_ATTR_EMULATE_PREPARES))
	vm.SetConstant("PDO::ATTR_DRIVER_NAME", data.NewIntValue(PDO_ATTR_DRIVER_NAME))
	vm.SetConstant("PDO::ATTR_CASE", data.NewIntValue(PDO_ATTR_CASE))
	vm.SetConstant("PDO::ATTR_ORACLE_NULLS", data.NewIntValue(PDO_ATTR_ORACLE_NULLS))
	vm.SetConstant("PDO::ATTR_STRINGIFY_FETCHES", data.NewIntValue(PDO_ATTR_STRINGIFY_FETCHES))
	vm.SetConstant("PDO::ATTR_STATEMENT_CLASS", data.NewIntValue(PDO_ATTR_STATEMENT_CLASS))
	vm.SetConstant("PDO::ATTR_TIMEOUT", data.NewIntValue(PDO_ATTR_TIMEOUT))
	vm.SetConstant("PDO::ATTR_SERVER_VERSION", data.NewIntValue(PDO_ATTR_SERVER_VERSION))
	vm.SetConstant("PDO::ATTR_CLIENT_VERSION", data.NewIntValue(PDO_ATTR_CLIENT_VERSION))
	vm.SetConstant("PDO::ATTR_SERVER_INFO", data.NewIntValue(PDO_ATTR_SERVER_INFO))
	vm.SetConstant("PDO::ATTR_CONNECTION_STATUS", data.NewIntValue(PDO_ATTR_CONNECTION_STATUS))

	vm.SetConstant("PDO::ERRMODE_SILENT", data.NewIntValue(PDO_ERRMODE_SILENT))
	vm.SetConstant("PDO::ERRMODE_WARNING", data.NewIntValue(PDO_ERRMODE_WARNING))
	vm.SetConstant("PDO::ERRMODE_EXCEPTION", data.NewIntValue(PDO_ERRMODE_EXCEPTION))

	vm.SetConstant("PDO::CASE_NATURAL", data.NewIntValue(PDO_CASE_NATURAL))
	vm.SetConstant("PDO::CASE_UPPER", data.NewIntValue(PDO_CASE_UPPER))
	vm.SetConstant("PDO::CASE_LOWER", data.NewIntValue(PDO_CASE_LOWER))

	vm.SetConstant("PDO::NULL_NATURAL", data.NewIntValue(PDO_NULL_NATURAL))
	vm.SetConstant("PDO::NULL_EMPTY_STRING", data.NewIntValue(PDO_NULL_EMPTY_STRING))
	vm.SetConstant("PDO::NULL_TO_STRING", data.NewIntValue(PDO_NULL_TO_STRING))

	vm.SetConstant("PDO::FETCH_DEFAULT", data.NewIntValue(PDO_FETCH_DEFAULT))
	vm.SetConstant("PDO::FETCH_LAZY", data.NewIntValue(PDO_FETCH_LAZY))
	vm.SetConstant("PDO::FETCH_ASSOC", data.NewIntValue(PDO_FETCH_ASSOC))
	vm.SetConstant("PDO::FETCH_NUM", data.NewIntValue(PDO_FETCH_NUM))
	vm.SetConstant("PDO::FETCH_BOTH", data.NewIntValue(PDO_FETCH_BOTH))
	vm.SetConstant("PDO::FETCH_OBJ", data.NewIntValue(PDO_FETCH_OBJ))
	vm.SetConstant("PDO::FETCH_BOUND", data.NewIntValue(PDO_FETCH_BOUND))
	vm.SetConstant("PDO::FETCH_COLUMN", data.NewIntValue(PDO_FETCH_COLUMN))
	vm.SetConstant("PDO::FETCH_CLASS", data.NewIntValue(PDO_FETCH_CLASS))
	vm.SetConstant("PDO::FETCH_INTO", data.NewIntValue(PDO_FETCH_INTO))
	vm.SetConstant("PDO::FETCH_FUNC", data.NewIntValue(PDO_FETCH_FUNC))
	vm.SetConstant("PDO::FETCH_GROUP", data.NewIntValue(PDO_FETCH_GROUP))
	vm.SetConstant("PDO::FETCH_UNIQUE", data.NewIntValue(PDO_FETCH_UNIQUE))
	vm.SetConstant("PDO::FETCH_KEY_PAIR", data.NewIntValue(PDO_FETCH_KEY_PAIR))
	vm.SetConstant("PDO::FETCH_NAMED", data.NewIntValue(PDO_FETCH_NAMED))
	vm.SetConstant("PDO::FETCH_PROPS_LATE", data.NewIntValue(PDO_FETCH_PROPS_LATE))

	vm.SetConstant("PDO::PARAM_BOOL", data.NewIntValue(PDO_PARAM_BOOL))
	vm.SetConstant("PDO::PARAM_NULL", data.NewIntValue(PDO_PARAM_NULL))
	vm.SetConstant("PDO::PARAM_INT", data.NewIntValue(PDO_PARAM_INT))
	vm.SetConstant("PDO::PARAM_STR", data.NewIntValue(PDO_PARAM_STR))
	vm.SetConstant("PDO::PARAM_LOB", data.NewIntValue(PDO_PARAM_LOB))
	vm.SetConstant("PDO::PARAM_INPUT_OUTPUT", data.NewIntValue(PDO_PARAM_INPUT_OUTPUT))

	vm.SetConstant("PDO::CURSOR_FWDONLY", data.NewIntValue(PDO_CURSOR_FWDONLY))
	vm.SetConstant("PDO::CURSOR_SCROLL", data.NewIntValue(PDO_CURSOR_SCROLL))

	// MySQL 驱动专属属性常量
	vm.SetConstant("PDO::MYSQL_ATTR_USE_BUFFERED_QUERY", data.NewIntValue(PDO_MYSQL_ATTR_USE_BUFFERED_QUERY))
	vm.SetConstant("PDO::MYSQL_ATTR_LOCAL_INFILE", data.NewIntValue(PDO_MYSQL_ATTR_LOCAL_INFILE))
	vm.SetConstant("PDO::MYSQL_ATTR_INIT_COMMAND", data.NewIntValue(PDO_MYSQL_ATTR_INIT_COMMAND))
	vm.SetConstant("PDO::MYSQL_ATTR_COMPRESS", data.NewIntValue(PDO_MYSQL_ATTR_COMPRESS))
	vm.SetConstant("PDO::MYSQL_ATTR_DIRECT_QUERY", data.NewIntValue(PDO_MYSQL_ATTR_DIRECT_QUERY))
	vm.SetConstant("PDO::MYSQL_ATTR_FOUND_ROWS", data.NewIntValue(PDO_MYSQL_ATTR_FOUND_ROWS))
	vm.SetConstant("PDO::MYSQL_ATTR_IGNORE_SPACE", data.NewIntValue(PDO_MYSQL_ATTR_IGNORE_SPACE))
	vm.SetConstant("PDO::MYSQL_ATTR_MAX_BUFFER_SIZE", data.NewIntValue(PDO_MYSQL_ATTR_MAX_BUFFER_SIZE))
	vm.SetConstant("PDO::MYSQL_ATTR_READ_DEFAULT_FILE", data.NewIntValue(PDO_MYSQL_ATTR_READ_DEFAULT_FILE))
	vm.SetConstant("PDO::MYSQL_ATTR_READ_DEFAULT_GROUP", data.NewIntValue(PDO_MYSQL_ATTR_READ_DEFAULT_GROUP))
	vm.SetConstant("PDO::MYSQL_ATTR_CONNECT_TIMEOUT", data.NewIntValue(PDO_MYSQL_ATTR_CONNECT_TIMEOUT))
	vm.SetConstant("PDO::MYSQL_ATTR_SERVER_PUBLIC_KEY", data.NewIntValue(PDO_MYSQL_ATTR_SERVER_PUBLIC_KEY))
	vm.SetConstant("PDO::MYSQL_ATTR_MULTI_STATEMENTS", data.NewIntValue(PDO_MYSQL_ATTR_MULTI_STATEMENTS))
	vm.SetConstant("PDO::MYSQL_ATTR_SSL_CA", data.NewIntValue(PDO_MYSQL_ATTR_SSL_CA))
	vm.SetConstant("PDO::MYSQL_ATTR_SSL_CAPATH", data.NewIntValue(PDO_MYSQL_ATTR_SSL_CAPATH))
	vm.SetConstant("PDO::MYSQL_ATTR_SSL_CERT", data.NewIntValue(PDO_MYSQL_ATTR_SSL_CERT))
	vm.SetConstant("PDO::MYSQL_ATTR_SSL_KEY", data.NewIntValue(PDO_MYSQL_ATTR_SSL_KEY))
	vm.SetConstant("PDO::MYSQL_ATTR_SSL_CIPHER", data.NewIntValue(PDO_MYSQL_ATTR_SSL_CIPHER))
	vm.SetConstant("PDO::MYSQL_ATTR_SSL_VERIFY_SERVER_CERT", data.NewIntValue(PDO_MYSQL_ATTR_SSL_VERIFY_SERVER_CERT))
}
