package pdo

import (
	"database/sql"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// -------------------------------------------------------------------
// PDOStatement 内部状态
// -------------------------------------------------------------------

type pdoStmtState struct {
	rows        *sql.Rows // query 返回的结果集
	stmt        *sql.Stmt // prepared statement
	sqlStr      string
	pdoState    *pdoState
	fetchMode   int
	cols        []string // 列名缓存
	colsFetched bool
}

// -------------------------------------------------------------------
// PDOStatementClass
// -------------------------------------------------------------------

type PDOStatementClass struct {
	node.Node
	state *pdoStmtState
}

func newPDOStatementClass(rows *sql.Rows, pState *pdoState) *PDOStatementClass {
	return &PDOStatementClass{
		state: &pdoStmtState{
			rows:      rows,
			pdoState:  pState,
			fetchMode: PDO_FETCH_BOTH,
		},
	}
}

func newPDOStatementClassFromPrepared(stmt *sql.Stmt, pState *pdoState, sqlStr string) *PDOStatementClass {
	return &PDOStatementClass{
		state: &pdoStmtState{
			stmt:      stmt,
			sqlStr:    sqlStr,
			pdoState:  pState,
			fetchMode: PDO_FETCH_BOTH,
		},
	}
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
	if m.state.stmt == nil {
		return data.NewBoolValue(false), nil
	}

	// 收集参数
	var args []interface{}
	if paramVal, ok := ctx.GetIndexValue(0); ok && paramVal != nil {
		if _, isNull := paramVal.(*data.NullValue); !isNull {
			if arr, ok := paramVal.(*data.ArrayValue); ok {
				for _, v := range arr.ToValueList() {
					args = append(args, v.AsString())
				}
			}
		}
	}

	rows, err := m.state.stmt.Query(args...)
	if err != nil {
		if m.state.pdoState != nil && m.state.pdoState.getErrMode() == PDO_ERRMODE_EXCEPTION {
			return nil, pdoException(err.Error(), ctx)
		}
		return data.NewBoolValue(false), nil
	}
	m.state.rows = rows
	m.state.cols = nil
	m.state.colsFetched = false
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
	if m.state.rows == nil {
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

	if !m.state.rows.Next() {
		return data.NewBoolValue(false), nil
	}

	cols, err := m.state.rows.Columns()
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	row, acl := scanRowToMap(m.state.rows, cols)
	if acl != nil {
		return nil, acl
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
	if m.state.rows == nil {
		return data.NewArrayValue(nil), nil
	}

	mode := m.state.fetchMode
	if modeVal, ok := ctx.GetIndexValue(0); ok && modeVal != nil {
		if ai, ok := modeVal.(interface{ AsInt() (int, error) }); ok {
			if v, err := ai.AsInt(); err == nil && v != 0 {
				mode = v
			}
		}
	}

	cols, err := m.state.rows.Columns()
	if err != nil {
		return data.NewArrayValue(nil), nil
	}

	var results []data.Value
	for m.state.rows.Next() {
		row, acl := scanRowToMap(m.state.rows, cols)
		if acl != nil {
			return nil, acl
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
	if m.state.rows == nil || !m.state.rows.Next() {
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

	cols, _ := m.state.rows.Columns()
	row, acl := scanRowToMap(m.state.rows, cols)
	if acl != nil {
		return nil, acl
	}

	if colIdx < len(cols) {
		return data.NewStringValue(row[cols[colIdx]]), nil
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
	if m.state.rows == nil {
		return data.NewIntValue(0), nil
	}
	cols, err := m.state.rows.Columns()
	if err != nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(cols)), nil
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
	if m.state.rows != nil {
		m.state.rows.Close()
		m.state.rows = nil
	}
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
	// 简化实现：暂不支持按名绑定，execute 时直接传参
	return data.NewBoolValue(true), nil
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
		if vals[i] == nil {
			row[col] = ""
		} else {
			row[col] = fmt.Sprintf("%v", vals[i])
		}
	}
	return row, nil
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
