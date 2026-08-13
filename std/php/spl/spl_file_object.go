package spl

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/stream"
	"github.com/php-any/origami/utils"
)

// SplFileObject flags（与 PHP 一致）
const (
	SFO_DROP_NEW_LINE = 1
	SFO_READ_AHEAD    = 2
	SFO_SKIP_EMPTY    = 4
	SFO_READ_CSV      = 8
)

const sfoStateKey = "__sfo_state__"

// sfoStateValue 存储 SplFileObject 运行时状�?
type sfoStateValue struct {
	stream          *stream.StreamInfo
	mode            string
	flags           int
	key             int
	current         string
	valid           bool
	eof             bool
	lines           []string
	lineIdx         int
	useBuf          bool
	reader          *bufio.Reader
	iterLinePending bool
	fgetsAtStart    bool
}

func (s *sfoStateValue) GetValue(ctx data.Context) (data.GetValue, data.Control) { return s, nil }
func (s *sfoStateValue) AsString() string                                        { return "sfoState" }
func (s *sfoStateValue) Marshal(serializer data.Serializer) ([]byte, error)      { return nil, nil }
func (s *sfoStateValue) Unmarshal(b []byte, serializer data.Serializer) error    { return nil }
func (s *sfoStateValue) ToGoValue(serializer data.Serializer) (any, error)       { return nil, nil }

// SplFileObjectClass 实现 PHP SplFileObject
type SplFileObjectClass struct {
	node.Node
}

func NewSplFileObjectClass() *SplFileObjectClass { return &SplFileObjectClass{} }

func (c *SplFileObjectClass) GetName() string { return "SplFileObject" }

func (c *SplFileObjectClass) GetExtend() *string {
	parent := "SplFileInfo"
	return &parent
}

func (c *SplFileObjectClass) GetImplements() []string {
	return []string{"RecursiveIterator", "SeekableIterator"}
}

func (c *SplFileObjectClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplFileObjectClass) GetPropertyList() []data.Property              { return nil }

func (c *SplFileObjectClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "DROP_NEW_LINE":
		return data.NewIntValue(SFO_DROP_NEW_LINE), true
	case "READ_AHEAD":
		return data.NewIntValue(SFO_READ_AHEAD), true
	case "SKIP_EMPTY":
		return data.NewIntValue(SFO_SKIP_EMPTY), true
	case "READ_CSV":
		return data.NewIntValue(SFO_READ_CSV), true
	}
	return nil, false
}

func (c *SplFileObjectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(sfiPathnameKey, data.NewStringValue(""))
	cv.SetProperty(sfoStateKey, &sfoStateValue{})
	return cv, nil
}

func (c *SplFileObjectClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &SFOConstructMethod{}, true
	case "fgets":
		return &SFOFgetsMethod{}, true
	case "fgetcsv":
		return &SFOFgetcsvMethod{}, true
	case "fputcsv":
		return &SFOFputcsvMethod{}, true
	case "fwrite":
		return &SFOFwriteMethod{}, true
	case "eof":
		return &SFOEofMethod{}, true
	case "rewind":
		return &SFORewindMethod{}, true
	case "valid":
		return &SFOValidMethod{}, true
	case "current":
		return &SFOCurrentMethod{}, true
	case "key":
		return &SFOKeyMethod{}, true
	case "next":
		return &SFONextMethod{}, true
	case "seek":
		return &SFOSeekMethod{}, true
	case "hasChildren":
		return &SFOHasChildrenMethod{}, true
	case "getChildren":
		return &SFOGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *SplFileObjectClass) GetMethods() []data.Method {
	return []data.Method{
		&SFOConstructMethod{},
		&SFOFgetsMethod{},
		&SFOFgetcsvMethod{},
		&SFOFputcsvMethod{},
		&SFOFwriteMethod{},
		&SFOEofMethod{},
		&SFORewindMethod{},
		&SFOValidMethod{},
		&SFOCurrentMethod{},
		&SFOKeyMethod{},
		&SFONextMethod{},
		&SFOSeekMethod{},
		&SFOHasChildrenMethod{},
		&SFOGetChildrenMethod{},
	}
}

func (c *SplFileObjectClass) GetConstruct() data.Method { return &SFOConstructMethod{} }

func sfoGetState(cv *data.ClassValue) *sfoStateValue {
	v, _ := cv.ObjectValue.GetProperty(sfoStateKey)
	if sv, ok := v.(*sfoStateValue); ok {
		return sv
	}
	return nil
}

func sfoSetState(cv *data.ClassValue, st *sfoStateValue) {
	cv.ObjectValue.SetProperty(sfoStateKey, st)
}

func sfoDropNewline(s string, flags int) string {
	if flags&SFO_DROP_NEW_LINE != 0 {
		s = strings.TrimRight(s, "\r\n")
	}
	return s
}

func sfoIsEmptyLine(s string, flags int) bool {
	if flags&SFO_SKIP_EMPTY == 0 {
		return false
	}
	return strings.TrimSpace(s) == ""
}

func sfoReadRawLine(st *sfoStateValue) (string, error) {
	if st.stream == nil || st.stream.IsClosed() {
		return "", io.EOF
	}
	if st.reader == nil {
		st.reader = bufio.NewReader(st.stream.File)
	}
	line, err := st.reader.ReadString('\n')
	if err == io.EOF && line == "" {
		st.eof = true
		return "", io.EOF
	}
	if err != nil && err != io.EOF {
		return "", err
	}
	if err == io.EOF {
		st.eof = true
	}
	return line, nil
}

func sfoLoadAllLines(st *sfoStateValue) error {
	st.lines = nil
	st.reader = nil
	if st.stream != nil {
		st.stream.Seek(0, io.SeekStart)
	}
	for {
		line, err := sfoReadRawLine(st)
		if err == io.EOF && line == "" {
			break
		}
		st.lines = append(st.lines, line)
		if err == io.EOF {
			break
		}
	}
	st.lineIdx = 0
	st.eof = len(st.lines) == 0
	return nil
}

func sfoReadNext(st *sfoStateValue) {
	if st == nil {
		return
	}
	if st.flags&SFO_READ_AHEAD != 0 {
		for {
			if st.lineIdx >= len(st.lines) {
				st.valid = false
				st.current = ""
				st.eof = true
				return
			}
			line := st.lines[st.lineIdx]
			if sfoIsEmptyLine(line, st.flags) {
				st.lineIdx++
				continue
			}
			st.current = sfoDropNewline(line, st.flags)
			st.key = st.lineIdx
			st.valid = true
			st.eof = false
			return
		}
	}
	for {
		line, err := sfoReadRawLine(st)
		if err == io.EOF && line == "" {
			st.valid = false
			st.current = ""
			st.eof = true
			return
		}
		line = sfoDropNewline(line, st.flags)
		if sfoIsEmptyLine(line, st.flags) {
			continue
		}
		st.current = line
		st.valid = true
		st.eof = false
		return
	}
}

func sfoAdvance(st *sfoStateValue) {
	if st == nil {
		return
	}
	if st.flags&SFO_READ_AHEAD != 0 {
		st.lineIdx++
		sfoReadNext(st)
		return
	}
	st.key++
	sfoReadNext(st)
}

func sfoRewind(st *sfoStateValue) error {
	if st == nil || st.stream == nil {
		return errors.New("SplFileObject not initialized")
	}
	st.key = 0
	st.eof = false
	st.valid = false
	st.current = ""
	st.fgetsAtStart = true
	if st.flags&SFO_READ_AHEAD != 0 {
		st.lineIdx = 0
		sfoReadNext(st)
		return nil
	}
	st.reader = nil
	if _, err := st.stream.Seek(0, io.SeekStart); err != nil {
		return err
	}
	sfoReadNext(st)
	return nil
}

func sfoSeek(st *sfoStateValue, offset int) error {
	if st == nil {
		return errors.New("SplFileObject not initialized")
	}
	if offset < 0 {
		offset = 0
	}
	if st.flags&SFO_READ_AHEAD != 0 {
		st.lineIdx = offset
		sfoReadNext(st)
		return nil
	}
	if err := sfoRewind(st); err != nil {
		return err
	}
	for i := 0; i < offset; i++ {
		if !st.valid {
			break
		}
		st.key++
		sfoReadNext(st)
	}
	return nil
}

func sfoOpenFile(cv *data.ClassValue, filename, mode string, flags int) error {
	if mode == "" {
		mode = "r"
	}
	si, err := stream.OpenFile(filename, mode)
	if err != nil {
		return err
	}
	st := &sfoStateValue{
		stream: si,
		mode:   mode,
		flags:  flags,
	}
	if flags&SFO_READ_AHEAD != 0 {
		if err := sfoLoadAllLines(st); err != nil {
			si.Close()
			return err
		}
	}
	sfiSetPathname(cv, filename)
	sfoSetState(cv, st)
	return sfoRewind(st)
}

func sfoCtxInt(ctx data.Context, idx int, def int) int {
	v, _ := ctx.GetIndexValue(idx)
	if v == nil {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func sfoCtxString(ctx data.Context, idx int, def string) string {
	v, _ := ctx.GetIndexValue(idx)
	if v == nil {
		return def
	}
	if s, ok := v.(data.AsString); ok {
		return s.AsString()
	}
	return v.AsString()
}

func sfoSyncForLineRead(st *sfoStateValue) error {
	if st == nil {
		return errors.New("SplFileObject not initialized")
	}
	if st.fgetsAtStart {
		st.reader = nil
		if _, err := st.stream.Seek(0, io.SeekStart); err != nil {
			return err
		}
		st.fgetsAtStart = false
	}
	return nil
}

// ---- __construct ----

type SFOConstructMethod struct{}

func (m *SFOConstructMethod) GetName() string            { return "__construct" }
func (m *SFOConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SFOConstructMethod) GetIsStatic() bool          { return false }
func (m *SFOConstructMethod) GetReturnType() data.Types  { return nil }
func (m *SFOConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "mode", 1, data.NewStringValue("r"), data.NewBaseType("string")),
		node.NewParameter(nil, "flags", 2, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *SFOConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "mode", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "flags", 2, data.NewBaseType("int")),
	}
}
func (m *SFOConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	filename := sfoCtxString(ctx, 0, "")
	if filename == "" {
		return nil, utils.NewThrow(errors.New("SplFileObject::__construct(): Argument #1 ($filename) cannot be empty"))
	}
	mode := sfoCtxString(ctx, 1, "r")
	flags := sfoCtxInt(ctx, 2, 0)
	if err := sfoOpenFile(cv, filename, mode, flags); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

// ---- fgets ----

type SFOFgetsMethod struct{}

func (m *SFOFgetsMethod) GetName() string               { return "fgets" }
func (m *SFOFgetsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOFgetsMethod) GetIsStatic() bool             { return false }
func (m *SFOFgetsMethod) GetReturnType() data.Types     { return nil }
func (m *SFOFgetsMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOFgetsMethod) GetVariables() []data.Variable { return nil }
func (m *SFOFgetsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil || st.stream == nil {
		return data.NewBoolValue(false), nil
	}
	if err := sfoSyncForLineRead(st); err != nil {
		return data.NewBoolValue(false), nil
	}
	line, err := sfoReadRawLine(st)
	if err == io.EOF && line == "" {
		return data.NewBoolValue(false), nil
	}
	line = sfoDropNewline(line, st.flags)
	st.current = line
	st.valid = line != "" || !st.eof
	return data.NewStringValue(line), nil
}

// ---- fgetcsv ----

type SFOFgetcsvMethod struct{}

func (m *SFOFgetcsvMethod) GetName() string            { return "fgetcsv" }
func (m *SFOFgetcsvMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SFOFgetcsvMethod) GetIsStatic() bool          { return false }
func (m *SFOFgetcsvMethod) GetReturnType() data.Types  { return nil }
func (m *SFOFgetcsvMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "separator", 0, data.NewStringValue(","), data.NewBaseType("string")),
		node.NewParameter(nil, "enclosure", 1, data.NewStringValue("\""), data.NewBaseType("string")),
		node.NewParameter(nil, "escape", 2, data.NewStringValue("\\"), data.NewBaseType("string")),
	}
}
func (m *SFOFgetcsvMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "separator", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "enclosure", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "escape", 2, data.NewBaseType("string")),
	}
}
func (m *SFOFgetcsvMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil || st.stream == nil {
		return data.NewBoolValue(false), nil
	}
	if err := sfoSyncForLineRead(st); err != nil {
		return data.NewBoolValue(false), nil
	}
	sep := sfoCtxString(ctx, 0, ",")
	enc := sfoCtxString(ctx, 1, "\"")
	esc := sfoCtxString(ctx, 2, "\\")
	line, err := sfoReadRawLine(st)
	if err == io.EOF && line == "" {
		return data.NewBoolValue(false), nil
	}
	line = sfoDropNewline(line, st.flags)
	r := csv.NewReader(strings.NewReader(line))
	r.Comma = rune(sep[0])
	if len(enc) > 0 {
		r.LazyQuotes = true
	}
	_ = esc
	fields, err := r.Read()
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	vals := make([]data.Value, len(fields))
	for i, f := range fields {
		vals[i] = data.NewStringValue(f)
	}
	return data.NewArrayValue(vals), nil
}

// ---- fwrite ----

type SFOFwriteMethod struct{}

func (m *SFOFwriteMethod) GetName() string            { return "fwrite" }
func (m *SFOFwriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SFOFwriteMethod) GetIsStatic() bool          { return false }
func (m *SFOFwriteMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SFOFwriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "length", 1, data.NewNullValue(), nil),
	}
}
func (m *SFOFwriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "length", 1, nil),
	}
}
func (m *SFOFwriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil || st.stream == nil {
		return data.NewBoolValue(false), nil
	}
	dataStr := sfoCtxString(ctx, 0, "")
	if dataStr == "" {
		return data.NewIntValue(0), nil
	}
	length := len(dataStr)
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			length = sfoCtxInt(ctx, 1, length)
			if length < 0 {
				length = 0
			}
			if length > len(dataStr) {
				length = len(dataStr)
			}
		}
	}
	n, err := st.stream.Write([]byte(dataStr[:length]))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	st.reader = nil
	return data.NewIntValue(n), nil
}

// ---- fputcsv ----

type SFOFputcsvMethod struct{}

func (m *SFOFputcsvMethod) GetName() string            { return "fputcsv" }
func (m *SFOFputcsvMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SFOFputcsvMethod) GetIsStatic() bool          { return false }
func (m *SFOFputcsvMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SFOFputcsvMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "fields", 0, nil, nil),
		node.NewParameter(nil, "separator", 1, data.NewStringValue(","), data.NewBaseType("string")),
		node.NewParameter(nil, "enclosure", 2, data.NewStringValue("\""), data.NewBaseType("string")),
		node.NewParameter(nil, "escape", 3, data.NewStringValue("\\"), data.NewBaseType("string")),
	}
}
func (m *SFOFputcsvMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "fields", 0, nil),
		node.NewVariable(nil, "separator", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "enclosure", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "escape", 3, data.NewBaseType("string")),
	}
}
func (m *SFOFputcsvMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil || st.stream == nil {
		return data.NewBoolValue(false), nil
	}
	fieldsVal, _ := ctx.GetIndexValue(0)
	arr, ok := fieldsVal.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	sep := sfoCtxString(ctx, 1, ",")
	enc := sfoCtxString(ctx, 2, "\"")
	if enc == "" {
		enc = "\""
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = rune(sep[0])
	strFields := make([]string, len(arr.List))
	for i, z := range arr.List {
		if z != nil {
			strFields[i] = z.Value.AsString()
		}
	}
	if err := w.Write(strFields); err != nil {
		return data.NewBoolValue(false), nil
	}
	w.Flush()
	out := buf.String()
	n, err := st.stream.Write([]byte(out))
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	_ = st.stream.Flush()
	st.reader = nil
	return data.NewIntValue(n), nil
}

// ---- eof ----

type SFOEofMethod struct{}

func (m *SFOEofMethod) GetName() string               { return "eof" }
func (m *SFOEofMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOEofMethod) GetIsStatic() bool             { return false }
func (m *SFOEofMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SFOEofMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOEofMethod) GetVariables() []data.Variable { return nil }
func (m *SFOEofMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(st.eof), nil
}

// ---- Iterator methods ----

type SFORewindMethod struct{}

func (m *SFORewindMethod) GetName() string               { return "rewind" }
func (m *SFORewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFORewindMethod) GetIsStatic() bool             { return false }
func (m *SFORewindMethod) GetReturnType() data.Types     { return nil }
func (m *SFORewindMethod) GetParams() []data.GetValue    { return nil }
func (m *SFORewindMethod) GetVariables() []data.Variable { return nil }
func (m *SFORewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if err := sfoRewind(st); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

type SFOValidMethod struct{}

func (m *SFOValidMethod) GetName() string               { return "valid" }
func (m *SFOValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOValidMethod) GetIsStatic() bool             { return false }
func (m *SFOValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SFOValidMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOValidMethod) GetVariables() []data.Variable { return nil }
func (m *SFOValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil {
		return data.NewBoolValue(false), nil
	}
	if !st.valid && !st.eof {
		sfoReadNext(st)
	}
	return data.NewBoolValue(st.valid), nil
}

type SFOCurrentMethod struct{}

func (m *SFOCurrentMethod) GetName() string               { return "current" }
func (m *SFOCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOCurrentMethod) GetIsStatic() bool             { return false }
func (m *SFOCurrentMethod) GetReturnType() data.Types     { return data.String{} }
func (m *SFOCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *SFOCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil {
		return data.NewStringValue(""), nil
	}
	if !st.valid && !st.eof {
		sfoReadNext(st)
	}
	if !st.valid {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(st.current), nil
}

type SFOKeyMethod struct{}

func (m *SFOKeyMethod) GetName() string               { return "key" }
func (m *SFOKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOKeyMethod) GetIsStatic() bool             { return false }
func (m *SFOKeyMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SFOKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOKeyMethod) GetVariables() []data.Variable { return nil }
func (m *SFOKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(st.key), nil
}

type SFONextMethod struct{}

func (m *SFONextMethod) GetName() string               { return "next" }
func (m *SFONextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFONextMethod) GetIsStatic() bool             { return false }
func (m *SFONextMethod) GetReturnType() data.Types     { return nil }
func (m *SFONextMethod) GetParams() []data.GetValue    { return nil }
func (m *SFONextMethod) GetVariables() []data.Variable { return nil }
func (m *SFONextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	if st == nil {
		return nil, utils.NewThrow(errors.New("SplFileObject not initialized"))
	}
	sfoAdvance(st)
	return nil, nil
}

type SFOSeekMethod struct{}

func (m *SFOSeekMethod) GetName() string            { return "seek" }
func (m *SFOSeekMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SFOSeekMethod) GetIsStatic() bool          { return false }
func (m *SFOSeekMethod) GetReturnType() data.Types  { return nil }
func (m *SFOSeekMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "offset", 0, data.NewIntValue(0), data.Int{}),
	}
}
func (m *SFOSeekMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Int{})}
}
func (m *SFOSeekMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	st := sfoGetState(cv)
	offset := sfoCtxInt(ctx, 0, 0)
	if err := sfoSeek(st, offset); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

type SFOHasChildrenMethod struct{}

func (m *SFOHasChildrenMethod) GetName() string               { return "hasChildren" }
func (m *SFOHasChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOHasChildrenMethod) GetIsStatic() bool             { return false }
func (m *SFOHasChildrenMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SFOHasChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOHasChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *SFOHasChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

type SFOGetChildrenMethod struct{}

func (m *SFOGetChildrenMethod) GetName() string               { return "getChildren" }
func (m *SFOGetChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SFOGetChildrenMethod) GetIsStatic() bool             { return false }
func (m *SFOGetChildrenMethod) GetReturnType() data.Types     { return nil }
func (m *SFOGetChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *SFOGetChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *SFOGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, utils.NewThrow(errors.New("SplFileObject::getChildren(): Cannot get children, not a directory"))
}

// sfoCreateTempFile 创建临时文件并返回路径（�?SplTempFileObject 使用�?
func sfoCreateTempFile(prefix string) (string, error) {
	if prefix == "" {
		prefix = "splt"
	}
	f, err := os.CreateTemp("", prefix+"_*")
	if err != nil {
		return "", err
	}
	name := f.Name()
	f.Close()
	return name, nil
}

// sfoResolveTempFilename 解析 SplTempFileObject 构造参数中的文件名
func sfoResolveTempFilename(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || name == "php://temp" || name == "php://memory" {
		return sfoCreateTempFile("splt")
	}
	return name, nil
}

// sfoOpenFileForTemp �?SplTempFileObject 复用的打开逻辑
func sfoOpenFileForTemp(cv *data.ClassValue, filename, mode string, flags int) error {
	resolved, err := sfoResolveTempFilename(filename)
	if err != nil {
		return err
	}
	return sfoOpenFile(cv, resolved, mode, flags)
}

// SFOConstantNames 返回 SplFileObject 常量名列表（�?load.go 注册�?
func SFOConstantNames() map[string]int {
	return map[string]int{
		"DROP_NEW_LINE": SFO_DROP_NEW_LINE,
		"READ_AHEAD":    SFO_READ_AHEAD,
		"SKIP_EMPTY":    SFO_SKIP_EMPTY,
		"READ_CSV":      SFO_READ_CSV,
	}
}

// SFOConstantValue 解析 SplFileObject::CONST 访问
func SFOConstantValue(name string) (data.Value, bool) {
	if v, ok := SFOConstantNames()[name]; ok {
		return data.NewIntValue(v), true
	}
	return nil, false
}
