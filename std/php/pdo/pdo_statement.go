package pdo

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// -------------------------------------------------------------------
// PDOStatement 内部状态
// -------------------------------------------------------------------

type pdoStmtState struct {
	sqlStr    string
	pdoState  *pdoState
	fetchMode int
	cols      []string // 列名缓存
	// 缓冲结果：借出的 *sql.Conn 在 Rows 未 Close 时不可复用；
	// 若不缓冲，同 PDO 上后续语句会阻塞。对齐 PHP 默认 buffered query。
	buffered []map[string]string
	rowPos   int
	// bound 按 PDO 1-based 位置参数存储（bindValue / bindParam）
	bound map[int]interface{}
}

// -------------------------------------------------------------------
// PDOStatementClass
// -------------------------------------------------------------------

type PDOStatementClass struct {
	node.Node
	state *pdoStmtState
}

func newPDOStatementClass(rows *sql.Rows, pState *pdoState) *PDOStatementClass {
	cols, buffered := bufferSQLRows(rows)
	return &PDOStatementClass{
		state: &pdoStmtState{
			pdoState:  pState,
			fetchMode: PDO_FETCH_BOTH,
			cols:      cols,
			buffered:  buffered,
			rowPos:    0,
		},
	}
}

// newPDOStatementFromSQL 延迟到 execute 时在 PDO 当前借出的连接（或事务）上执行。
func newPDOStatementFromSQL(pState *pdoState, sqlStr string) *PDOStatementClass {
	return &PDOStatementClass{
		state: &pdoStmtState{
			sqlStr:    sqlStr,
			pdoState:  pState,
			fetchMode: PDO_FETCH_BOTH,
		},
	}
}

// bufferSQLRows 读完并关闭 *sql.Rows，释放连接上的游标占用。
func bufferSQLRows(rows *sql.Rows) ([]string, []map[string]string) {
	if rows == nil {
		return nil, nil
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil
	}
	var buffered []map[string]string
	for rows.Next() {
		row, acl := scanRowToMap(rows, cols)
		if acl != nil {
			return cols, buffered
		}
		buffered = append(buffered, row)
	}
	return cols, buffered
}

func (c *PDOStatementClass) GetName() string                            { return "PDOStatement" }
func (c *PDOStatementClass) GetExtend() *string                         { return nil }
func (c *PDOStatementClass) GetImplements() []string                    { return nil }
func (c *PDOStatementClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *PDOStatementClass) GetPropertyList() []data.Property           { return nil }
func (c *PDOStatementClass) GetConstruct() data.Method                  { return nil }

func (c *PDOStatementClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *PDOStatementClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "execute":
		return &stmtExecuteMethod{state: c.state}, true
	case "fetch":
		return &stmtFetchMethod{state: c.state}, true
	case "fetchAll":
		return &stmtFetchAllMethod{state: c.state}, true
	case "fetchColumn":
		return &stmtFetchColumnMethod{state: c.state}, true
	case "rowCount":
		return &stmtRowCountMethod{state: c.state}, true
	case "columnCount":
		return &stmtColumnCountMethod{state: c.state}, true
	case "closeCursor":
		return &stmtCloseCursorMethod{state: c.state}, true
	case "setFetchMode":
		return &stmtSetFetchModeMethod{state: c.state}, true
	case "errorCode":
		return &stmtErrorCodeMethod{state: c.state}, true
	case "errorInfo":
		return &stmtErrorInfoMethod{state: c.state}, true
	case "bindParam", "bindValue":
		return &stmtBindParamMethod{state: c.state}, true
	case "getColumnMeta":
		return &stmtGetColumnMetaMethod{state: c.state}, true
	}
	return nil, false
}

func (c *PDOStatementClass) GetMethods() []data.Method { return nil }

// -------------------------------------------------------------------
// execute(?array $params=null): bool
// -------------------------------------------------------------------

type stmtExecuteMethod struct{ state *pdoStmtState }

func (m *stmtExecuteMethod) GetName() string            { return "execute" }
func (m *stmtExecuteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtExecuteMethod) GetIsStatic() bool          { return false }
func (m *stmtExecuteMethod) GetReturnType() data.Types  { return nil }
func (m *stmtExecuteMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "params", 0, node.NewNullLiteral(nil), nil)}
}
func (m *stmtExecuteMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "params", 0, data.NewBaseType("array"))}
}

func (m *stmtExecuteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.state.pdoState == nil || strings.TrimSpace(m.state.sqlStr) == "" {
		return data.NewBoolValue(false), nil
	}

	// 收集参数：优先 execute($params)，否则使用 bindValue 绑定的位置参数
	var args []interface{}
	if paramVal, ok := ctx.GetIndexValue(0); ok && paramVal != nil {
		if _, isNull := paramVal.(*data.NullValue); !isNull {
			if arr, ok := paramVal.(*data.ArrayValue); ok {
				for _, v := range arr.ToValueList() {
					args = append(args, phpValueToDriver(v))
				}
			}
		}
	}
	if len(args) == 0 && len(m.state.bound) > 0 {
		max := 0
		for k := range m.state.bound {
			if k > max {
				max = k
			}
		}
		args = make([]interface{}, max)
		for i := 1; i <= max; i++ {
			args[i-1] = m.state.bound[i]
		}
	}

	sqlUpper := strings.ToUpper(strings.TrimSpace(m.state.sqlStr))
	isQuery := strings.HasPrefix(sqlUpper, "SELECT") ||
		strings.HasPrefix(sqlUpper, "WITH") ||
		strings.HasPrefix(sqlUpper, "PRAGMA") ||
		strings.HasPrefix(sqlUpper, "EXPLAIN") ||
		strings.HasPrefix(sqlUpper, "SHOW")

	state := m.state.pdoState
	state.mu.Lock()
	defer state.mu.Unlock()

	if isQuery {
		rows, err := state.queryLocked(m.state.sqlStr, args...)
		if err != nil {
			state.lastError = err.Error()
			if state.getErrMode() == PDO_ERRMODE_EXCEPTION {
				return nil, pdoException(err.Error(), ctx)
			}
			return data.NewBoolValue(false), nil
		}
		m.state.cols, m.state.buffered = bufferSQLRows(rows)
		m.state.rowPos = 0
		return data.NewBoolValue(true), nil
	}

	result, err := state.execLocked(m.state.sqlStr, args...)
	if err != nil {
		state.lastError = err.Error()
		if state.getErrMode() == PDO_ERRMODE_EXCEPTION {
			return nil, pdoException(err.Error(), ctx)
		}
		return data.NewBoolValue(false), nil
	}
	if id, err := result.LastInsertId(); err == nil {
		state.lastInsertID = id
	}
	m.state.buffered = nil
	m.state.cols = nil
	m.state.rowPos = 0
	return data.NewBoolValue(true), nil
}

// -------------------------------------------------------------------
// fetch(int $mode=PDO::FETCH_DEFAULT, ...): mixed
// -------------------------------------------------------------------

type stmtFetchMethod struct{ state *pdoStmtState }

func (m *stmtFetchMethod) GetName() string            { return "fetch" }
func (m *stmtFetchMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtFetchMethod) GetIsStatic() bool          { return false }
func (m *stmtFetchMethod) GetReturnType() data.Types  { return nil }
func (m *stmtFetchMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "mode", 0, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (m *stmtFetchMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "mode", 0, data.NewBaseType("int"))}
}

func (m *stmtFetchMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.state.rowPos >= len(m.state.buffered) {
		return data.NewBoolValue(false), nil
	}

	mode := m.state.fetchMode
	if modeVal, ok := ctx.GetIndexValue(0); ok && modeVal != nil {
		if ai, ok := modeVal.(interface{ AsInt() (int, error) }); ok {
			if v, err := ai.AsInt(); err == nil && v != 0 {
				mode = v
			}
		}
	}

	row := m.state.buffered[m.state.rowPos]
	m.state.rowPos++
	cols := m.state.cols
	if len(cols) == 0 {
		cols = make([]string, 0, len(row))
		for k := range row {
			cols = append(cols, k)
		}
	}
	return buildFetchResult(row, cols, mode), nil
}

// -------------------------------------------------------------------
// fetchAll(int $mode=PDO::FETCH_DEFAULT, ...): array
// -------------------------------------------------------------------

type stmtFetchAllMethod struct{ state *pdoStmtState }

func (m *stmtFetchAllMethod) GetName() string            { return "fetchAll" }
func (m *stmtFetchAllMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtFetchAllMethod) GetIsStatic() bool          { return false }
func (m *stmtFetchAllMethod) GetReturnType() data.Types  { return nil }
func (m *stmtFetchAllMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "mode", 0, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (m *stmtFetchAllMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "mode", 0, data.NewBaseType("int"))}
}

func (m *stmtFetchAllMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	mode := m.state.fetchMode
	if modeVal, ok := ctx.GetIndexValue(0); ok && modeVal != nil {
		if ai, ok := modeVal.(interface{ AsInt() (int, error) }); ok {
			if v, err := ai.AsInt(); err == nil && v != 0 {
				mode = v
			}
		}
	}

	cols := m.state.cols
	results := make([]data.Value, 0, len(m.state.buffered)-m.state.rowPos)
	for m.state.rowPos < len(m.state.buffered) {
		row := m.state.buffered[m.state.rowPos]
		m.state.rowPos++
		if len(cols) == 0 {
			cols = make([]string, 0, len(row))
			for k := range row {
				cols = append(cols, k)
			}
		}
		results = append(results, buildFetchResult(row, cols, mode).(data.Value))
	}
	return data.NewArrayValue(results), nil
}

// -------------------------------------------------------------------
// fetchColumn(int $column=0): mixed
// -------------------------------------------------------------------

type stmtFetchColumnMethod struct{ state *pdoStmtState }

func (m *stmtFetchColumnMethod) GetName() string            { return "fetchColumn" }
func (m *stmtFetchColumnMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtFetchColumnMethod) GetIsStatic() bool          { return false }
func (m *stmtFetchColumnMethod) GetReturnType() data.Types  { return nil }
func (m *stmtFetchColumnMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "column", 0, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (m *stmtFetchColumnMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "column", 0, data.NewBaseType("int"))}
}

func (m *stmtFetchColumnMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.state.rowPos >= len(m.state.buffered) {
		return data.NewBoolValue(false), nil
	}

	colIdx := 0
	if colVal, ok := ctx.GetIndexValue(0); ok && colVal != nil {
		if ai, ok := colVal.(interface{ AsInt() (int, error) }); ok {
			if v, err := ai.AsInt(); err == nil {
				colIdx = v
			}
		}
	}

	row := m.state.buffered[m.state.rowPos]
	m.state.rowPos++
	cols := m.state.cols
	if colIdx >= 0 && colIdx < len(cols) {
		return data.NewStringValue(row[cols[colIdx]]), nil
	}
	// 无列名缓存时按 map 迭代顺序不稳定；优先用数字键兼容
	if len(cols) == 0 && colIdx == 0 {
		for _, v := range row {
			return data.NewStringValue(v), nil
		}
	}
	return data.NewBoolValue(false), nil
}

// -------------------------------------------------------------------
// rowCount(): int
// -------------------------------------------------------------------

type stmtRowCountMethod struct{ state *pdoStmtState }

func (m *stmtRowCountMethod) GetName() string               { return "rowCount" }
func (m *stmtRowCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *stmtRowCountMethod) GetIsStatic() bool             { return false }
func (m *stmtRowCountMethod) GetReturnType() data.Types     { return nil }
func (m *stmtRowCountMethod) GetParams() []data.GetValue    { return nil }
func (m *stmtRowCountMethod) GetVariables() []data.Variable { return nil }
func (m *stmtRowCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// SELECT 语句无法直接获取 rowCount，返回 -1 (与 PHP 行为一致)
	return data.NewIntValue(-1), nil
}

// -------------------------------------------------------------------
// columnCount(): int
// -------------------------------------------------------------------

type stmtColumnCountMethod struct{ state *pdoStmtState }

func (m *stmtColumnCountMethod) GetName() string               { return "columnCount" }
func (m *stmtColumnCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *stmtColumnCountMethod) GetIsStatic() bool             { return false }
func (m *stmtColumnCountMethod) GetReturnType() data.Types     { return nil }
func (m *stmtColumnCountMethod) GetParams() []data.GetValue    { return nil }
func (m *stmtColumnCountMethod) GetVariables() []data.Variable { return nil }
func (m *stmtColumnCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(len(m.state.cols)), nil
}

// -------------------------------------------------------------------
// closeCursor(): bool
// -------------------------------------------------------------------

type stmtCloseCursorMethod struct{ state *pdoStmtState }

func (m *stmtCloseCursorMethod) GetName() string               { return "closeCursor" }
func (m *stmtCloseCursorMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *stmtCloseCursorMethod) GetIsStatic() bool             { return false }
func (m *stmtCloseCursorMethod) GetReturnType() data.Types     { return nil }
func (m *stmtCloseCursorMethod) GetParams() []data.GetValue    { return nil }
func (m *stmtCloseCursorMethod) GetVariables() []data.Variable { return nil }
func (m *stmtCloseCursorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	m.state.buffered = nil
	m.state.cols = nil
	m.state.rowPos = 0
	return data.NewBoolValue(true), nil
}

// -------------------------------------------------------------------
// setFetchMode(int $mode, ...): bool
// -------------------------------------------------------------------

type stmtSetFetchModeMethod struct{ state *pdoStmtState }

func (m *stmtSetFetchModeMethod) GetName() string            { return "setFetchMode" }
func (m *stmtSetFetchModeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtSetFetchModeMethod) GetIsStatic() bool          { return false }
func (m *stmtSetFetchModeMethod) GetReturnType() data.Types  { return nil }
func (m *stmtSetFetchModeMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "mode", 0, nil, nil)}
}
func (m *stmtSetFetchModeMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "mode", 0, data.NewBaseType("int"))}
}
func (m *stmtSetFetchModeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if modeVal, ok := ctx.GetIndexValue(0); ok && modeVal != nil {
		if ai, ok := modeVal.(interface{ AsInt() (int, error) }); ok {
			if v, err := ai.AsInt(); err == nil {
				m.state.fetchMode = v
			}
		}
	}
	return data.NewBoolValue(true), nil
}

// -------------------------------------------------------------------
// errorCode / errorInfo
// -------------------------------------------------------------------

type stmtErrorCodeMethod struct{ state *pdoStmtState }

func (m *stmtErrorCodeMethod) GetName() string               { return "errorCode" }
func (m *stmtErrorCodeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *stmtErrorCodeMethod) GetIsStatic() bool             { return false }
func (m *stmtErrorCodeMethod) GetReturnType() data.Types     { return nil }
func (m *stmtErrorCodeMethod) GetParams() []data.GetValue    { return nil }
func (m *stmtErrorCodeMethod) GetVariables() []data.Variable { return nil }
func (m *stmtErrorCodeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.state.pdoState == nil || m.state.pdoState.lastSQLState == "" {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(m.state.pdoState.lastSQLState), nil
}

type stmtErrorInfoMethod struct{ state *pdoStmtState }

func (m *stmtErrorInfoMethod) GetName() string               { return "errorInfo" }
func (m *stmtErrorInfoMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *stmtErrorInfoMethod) GetIsStatic() bool             { return false }
func (m *stmtErrorInfoMethod) GetReturnType() data.Types     { return nil }
func (m *stmtErrorInfoMethod) GetParams() []data.GetValue    { return nil }
func (m *stmtErrorInfoMethod) GetVariables() []data.Variable { return nil }
func (m *stmtErrorInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	state := m.state.pdoState
	sqlState := ""
	errMsg := ""
	if state != nil {
		sqlState = state.lastSQLState
		errMsg = state.lastError
	}
	return data.NewArrayValue([]data.Value{
		data.NewStringValue(sqlState),
		data.NewNullValue(),
		data.NewStringValue(errMsg),
	}), nil
}

// -------------------------------------------------------------------
// bindParam / bindValue (stub)
// -------------------------------------------------------------------

type stmtBindParamMethod struct{ state *pdoStmtState }

func (m *stmtBindParamMethod) GetName() string            { return "bindParam" }
func (m *stmtBindParamMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtBindParamMethod) GetIsStatic() bool          { return false }
func (m *stmtBindParamMethod) GetReturnType() data.Types  { return nil }
func (m *stmtBindParamMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param", 0, nil, nil),
		node.NewParameter(nil, "var", 1, nil, nil),
		node.NewParameter(nil, "type", 2, node.NewIntLiteral(nil, "2"), nil),
	}
}
func (m *stmtBindParamMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "var", 1, data.NewBaseType("mixed")),
		node.NewVariable(nil, "type", 2, data.NewBaseType("int")),
	}
}
func (m *stmtBindParamMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramVal, ok := ctx.GetIndexValue(0)
	if !ok || paramVal == nil {
		return data.NewBoolValue(false), nil
	}
	varVal, ok := ctx.GetIndexValue(1)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if m.state.bound == nil {
		m.state.bound = make(map[int]interface{})
	}

	idx := 0
	if ai, ok := paramVal.(interface{ AsInt() (int, error) }); ok {
		if v, err := ai.AsInt(); err == nil {
			idx = v
		}
	}
	if idx <= 0 {
		// 非正整数位置：若是数字字符串则解析，否则按调用顺序追加
		if n, err := strconv.Atoi(strings.TrimSpace(paramVal.AsString())); err == nil && n > 0 {
			idx = n
		} else {
			idx = len(m.state.bound) + 1
		}
	}
	driverVal := phpValueToDriver(varVal)
	m.state.bound[idx] = driverVal
	return data.NewBoolValue(true), nil
}

func phpValueToDriver(v data.Value) interface{} {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case *data.NullValue:
		return nil
	case *data.BoolValue:
		return t.Value
	case *data.IntValue:
		if i, err := t.AsInt(); err == nil {
			return i
		}
	case *data.FloatValue:
		if f, err := t.AsFloat(); err == nil {
			return f
		}
	}
	return v.AsString()
}

// -------------------------------------------------------------------
// getColumnMeta(int $column): array|false
// -------------------------------------------------------------------

type stmtGetColumnMetaMethod struct{ state *pdoStmtState }

func (m *stmtGetColumnMetaMethod) GetName() string            { return "getColumnMeta" }
func (m *stmtGetColumnMetaMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *stmtGetColumnMetaMethod) GetIsStatic() bool          { return false }
func (m *stmtGetColumnMetaMethod) GetReturnType() data.Types  { return nil }
func (m *stmtGetColumnMetaMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "column", 0, nil, nil)}
}
func (m *stmtGetColumnMetaMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "column", 0, data.NewBaseType("int"))}
}
func (m *stmtGetColumnMetaMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

// -------------------------------------------------------------------
// 行扫描辅助
// -------------------------------------------------------------------

func scanRowToMap(rows *sql.Rows, cols []string) (map[string]string, data.Control) {
	vals := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("PDO scan error: %v", err))
	}
	row := make(map[string]string, len(cols))
	for i, col := range cols {
		row[col] = driverValueString(vals[i])
	}
	return row, nil
}

func driverValueString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case []byte:
		return string(v)
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

func buildFetchResult(row map[string]string, cols []string, mode int) data.GetValue {
	switch mode {
	case PDO_FETCH_ASSOC:
		obj := data.NewObjectValue()
		for k, v := range row {
			obj.SetProperty(k, data.NewStringValue(v))
		}
		return obj

	case PDO_FETCH_NUM:
		vals := make([]data.Value, len(cols))
		for i, col := range cols {
			vals[i] = data.NewStringValue(row[col])
		}
		return data.NewArrayValue(vals)

	case PDO_FETCH_OBJ:
		obj := data.NewObjectValue()
		for k, v := range row {
			obj.SetProperty(k, data.NewStringValue(v))
		}
		return obj

	default: // PDO_FETCH_BOTH
		obj := data.NewObjectValue()
		for i, col := range cols {
			obj.SetProperty(col, data.NewStringValue(row[col]))
			obj.SetProperty(fmt.Sprintf("%d", i), data.NewStringValue(row[col]))
		}
		return obj
	}
}
