package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

const (
	lsName    = "Origami Language Server"
	lsVersion = "1.0.1"
)

// 初始化日志配置
func initLogger(level int, output io.Writer) {
	// 设置输出
	if output == nil {
		output = os.Stderr
	}
	logrus.SetOutput(output)

	// 设置日志格式为文本格式，包含调用位置
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 获取文件名和行号
			filename := f.File
			if idx := strings.LastIndex(filename, "/"); idx != -1 {
				filename = filename[idx+1:]
			}
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// 启用调用位置信息
	logrus.SetReportCaller(true)

	// 设置日志级别 - 直接强转
	logrus.SetLevel(logrus.Level(level))
}

var (
	version    = flag.Bool("version", false, "Show version information")
	help       = flag.Bool("help", false, "Show help information")
	test       = flag.Bool("test", false, "Run definition jump test")
	protocol_  = flag.String("protocol", "stdio", "Protocol to use (stdio, tcp, websocket)")
	address    = flag.String("address", "localhost", "Address to bind to (for tcp/websocket)")
	port       = flag.Int("port", 8800, "Port to bind to (for tcp/websocket)")
	logLevel   = flag.Int("log-level", 5, "Log level (0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug)")
	logFile    = flag.String("log-file", "lsp.log", "Log file path (default lsp.log; empty not allowed)")
	consoleLog = flag.Bool("console-log", true, "Enable console logging in stdio mode (default: true)")
	scanDir    = flag.String("scan-dir", "", "Directory to scan for .zy files (optional)")

	// 全局 LSP 连接，用于发送通知
	globalConn *jsonrpc2.Conn
)

// 简化的 LSP 类型定义
type InitializeParams struct {
	ProcessID        *int               `json:"processId,omitempty"`
	ClientInfo       *ClientInfo        `json:"clientInfo,omitempty"`
	RootPath         *string            `json:"rootPath,omitempty"`
	RootURI          *string            `json:"rootUri,omitempty"`
	Capabilities     ClientCapabilities `json:"capabilities"`
	Trace            *string            `json:"trace,omitempty"`
	WorkspaceFolders []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type ClientCapabilities struct {
	Workspace    interface{} `json:"workspace,omitempty"`
	TextDocument interface{} `json:"textDocument,omitempty"`
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync       *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
	CompletionProvider     *CompletionOptions       `json:"completionProvider,omitempty"`
	HoverProvider          *HoverOptions            `json:"hoverProvider,omitempty"`
	DefinitionProvider     *DefinitionOptions       `json:"definitionProvider,omitempty"`
	DocumentSymbolProvider *DocumentSymbolOptions   `json:"documentSymbolProvider,omitempty"`
}

type ServerInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type TextDocumentSyncOptions struct {
	OpenClose         *bool        `json:"openClose,omitempty"`
	Change            *int         `json:"change,omitempty"`
	WillSave          *bool        `json:"willSave,omitempty"`
	WillSaveWaitUntil *bool        `json:"willSaveWaitUntil,omitempty"`
	Save              *SaveOptions `json:"save,omitempty"`
}

type SaveOptions struct {
	IncludeText *bool `json:"includeText,omitempty"`
}

type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	ResolveProvider   *bool    `json:"resolveProvider,omitempty"`
}

type HoverOptions struct {
	WorkDoneProgressOptions
}

type DefinitionOptions struct {
	WorkDoneProgressOptions
}

type DocumentSymbolOptions struct {
	WorkDoneProgressOptions
}

type WorkDoneProgressOptions struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
}

// 文档相关类型
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength *int   `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      uint32 `json:"line"`
	Character uint32 `json:"character"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// 补全相关类型
type CompletionParams struct {
	TextDocumentPositionParams
	Context *CompletionContext `json:"context,omitempty"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter *string               `json:"triggerCharacter,omitempty"`
}

type CompletionTriggerKind int

const (
	CompletionTriggerKindInvoked                         CompletionTriggerKind = 1
	CompletionTriggerKindTriggerCharacter                CompletionTriggerKind = 2
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label               string              `json:"label"`
	Kind                *CompletionItemKind `json:"kind,omitempty"`
	Tags                []CompletionItemTag `json:"tags,omitempty"`
	Detail              *string             `json:"detail,omitempty"`
	Documentation       *MarkupContent      `json:"documentation,omitempty"`
	Deprecated          *bool               `json:"deprecated,omitempty"`
	Preselect           *bool               `json:"preselect,omitempty"`
	SortText            *string             `json:"sortText,omitempty"`
	FilterText          *string             `json:"filterText,omitempty"`
	InsertText          *string             `json:"insertText,omitempty"`
	InsertTextFormat    *InsertTextFormat   `json:"insertTextFormat,omitempty"`
	InsertTextMode      *InsertTextMode     `json:"insertTextMode,omitempty"`
	TextEdit            *TextEdit           `json:"textEdit,omitempty"`
	TextEditText        *string             `json:"textEditText,omitempty"`
	AdditionalTextEdits []TextEdit          `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string            `json:"commitCharacters,omitempty"`
	Command             *Command            `json:"command,omitempty"`
	Data                interface{}         `json:"data,omitempty"`
}

type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

type CompletionItemTag int

const (
	CompletionItemTagDeprecated CompletionItemTag = 1
)

type InsertTextFormat int

const (
	InsertTextFormatPlainText InsertTextFormat = 1
	InsertTextFormatSnippet   InsertTextFormat = 2
)

type InsertTextMode int

const (
	InsertTextModeAsIs   InsertTextMode = 1
	InsertTextModeAdjust InsertTextMode = 2
)

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type MarkupKind string

const (
	MarkupKindPlainText MarkupKind = "plaintext"
	MarkupKindMarkdown  MarkupKind = "markdown"
)

// 悬停相关类型
type HoverParams struct {
	TextDocumentPositionParams
}

type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// 定义跳转相关类型
type DefinitionParams struct {
	TextDocumentPositionParams
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// 文档符号相关类型
type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         *string          `json:"detail,omitempty"`
	Kind           SymbolKind       `json:"kind"`
	Tags           []SymbolTag      `json:"tags,omitempty"`
	Deprecated     *bool            `json:"deprecated,omitempty"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

type SymbolKind int

const (
	SymbolKindFile          SymbolKind = 1
	SymbolKindModule        SymbolKind = 2
	SymbolKindNamespace     SymbolKind = 3
	SymbolKindPackage       SymbolKind = 4
	SymbolKindClass         SymbolKind = 5
	SymbolKindMethod        SymbolKind = 6
	SymbolKindProperty      SymbolKind = 7
	SymbolKindField         SymbolKind = 8
	SymbolKindConstructor   SymbolKind = 9
	SymbolKindEnum          SymbolKind = 10
	SymbolKindInterface     SymbolKind = 11
	SymbolKindFunction      SymbolKind = 12
	SymbolKindVariable      SymbolKind = 13
	SymbolKindConstant      SymbolKind = 14
	SymbolKindString        SymbolKind = 15
	SymbolKindNumber        SymbolKind = 16
	SymbolKindBoolean       SymbolKind = 17
	SymbolKindArray         SymbolKind = 18
	SymbolKindObject        SymbolKind = 19
	SymbolKindKey           SymbolKind = 20
	SymbolKindNull          SymbolKind = 21
	SymbolKindEnumMember    SymbolKind = 22
	SymbolKindStruct        SymbolKind = 23
	SymbolKindEvent         SymbolKind = 24
	SymbolKindOperator      SymbolKind = 25
	SymbolKindTypeParameter SymbolKind = 26
)

type SymbolTag int

const (
	SymbolTagDeprecated SymbolTag = 1
)

// 诊断相关类型
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range              Range                          `json:"range"`
	Severity           *DiagnosticSeverity            `json:"severity,omitempty"`
	Code               *interface{}                   `json:"code,omitempty"`
	Source             *string                        `json:"source,omitempty"`
	Message            string                         `json:"message"`
	Tags               []DiagnosticTag                `json:"tags,omitempty"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

type DiagnosticSeverity int

const (
	DiagnosticSeverityError   DiagnosticSeverity = 1
	DiagnosticSeverityWarning DiagnosticSeverity = 2
	DiagnosticSeverityInfo    DiagnosticSeverity = 3
	DiagnosticSeverityHint    DiagnosticSeverity = 4
)

type DiagnosticTag int

const (
	DiagnosticTagUnnecessary DiagnosticTag = 1
	DiagnosticTagDeprecated  DiagnosticTag = 2
)

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// 其他必要的类型定义
type InitializedParams struct{}

type SetTraceParams struct {
	Value string `json:"value"`
}

// 全局变量
var (
	globalLspVM *LspVM
	documents   = make(map[string]*DocumentInfo)
)

func main() {
	flag.Parse()

	// 初始化全局日志器
	var output io.Writer
	var effectiveLevel int = *logLevel

	// 强制 log-file 非空，空则重置为默认文件名
	if *logFile == "" {
		*logFile = "lsp.log"
	}

	// 转为绝对路径，方便定位
	if abs, err := filepath.Abs(*logFile); err == nil {
		*logFile = abs
	}

	// 根据协议类型决定日志输出方式
	if *protocol_ == "stdio" {
		// stdio 模式下，根据 consoleLog 选项决定日志输出方式
		var writers []io.Writer

		// 如果启用控制台日志，输出到 stderr（避免干扰 LSP 通信）
		if *consoleLog {
			writers = append(writers, os.Stderr)
		}

		// 尝试同时输出到文件（如果可能）
		if dir := filepath.Dir(*logFile); dir != "." && dir != "" {
			_ = os.MkdirAll(dir, 0755)
		}
		if file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			writers = append(writers, file)
			defer file.Close()
		}

		// 如果没有可用的输出目标，使用 Discard
		if len(writers) == 0 {
			output = io.Discard
		} else {
			// 使用 MultiWriter 同时输出到多个目标
			output = io.MultiWriter(writers...)
		}
	} else {
		// TCP/WebSocket 模式下，只输出到文件
		if dir := filepath.Dir(*logFile); dir != "." && dir != "" {
			_ = os.MkdirAll(dir, 0755)
		}
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// 无法写入日志文件时，输出到 stderr
			fmt.Fprintf(os.Stderr, "无法写入日志文件，将输出到 stderr：%s，原因：%v\n", *logFile, err)
			output = os.Stderr
		} else {
			defer file.Close()
			output = file
		}
	}

	initLogger(effectiveLevel, output)
	if effectiveLevel != int(logrus.PanicLevel) {
		if *protocol_ == "stdio" {
			if *consoleLog {
				logrus.Infof("日志输出: stderr + %s", *logFile)
			} else {
				logrus.Infof("日志文件: %s", *logFile)
			}
		} else {
			logrus.Infof("日志文件: %s", *logFile)
		}
	}

	if *version {
		fmt.Printf("%s v%s\n", lsName, lsVersion)
		return
	}

	if *help {
		showHelp()
		return
	}

	logrus.Infof("日志级别: %d", *logLevel)

	// 初始化全局 LspVM，如果指定了扫描目录则扫描该目录
	if *scanDir != "" {
		logrus.Infof("使用扫描目录模式，目录: %s", *scanDir)
		globalLspVM = NewLspVMWithScanDir(*scanDir)
	} else {
		globalLspVM = NewLspVM()
	}

	// 创建 JSON-RPC 2.0 处理器
	handler := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		return handleRequest(ctx, conn, req)
	})

	// 创建连接
	var conn *jsonrpc2.Conn
	switch *protocol_ {
	case "stdio":
		logrus.Infof("使用 stdio 协议启动 LSP 服务器")
		stream := jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{})
		conn = jsonrpc2.NewConn(context.Background(), stream, handler)
		<-conn.DisconnectNotify()
	case "tcp":
		addr := fmt.Sprintf("%s:%d", *address, *port)
		logrus.Infof("使用 TCP 协议在 %s 启动 LSP 服务器", addr)

		// 创建 TCP 监听器
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			logrus.Errorf("启动 TCP 服务器失败：%v", err)
			os.Exit(1)
		}
		defer listener.Close()

		logrus.Infof("TCP 服务器正在监听 %s", addr)

		// 接受连接
		for {
			conn, err := listener.Accept()
			if err != nil {
				logrus.Warnf("接受连接失败：%v", err)
				continue
			}

			logrus.Infof("来自 %s 的新 TCP 连接", conn.RemoteAddr())

			// 为每个连接创建一个新的 JSON-RPC 连接
			stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})
			rpcConn := jsonrpc2.NewConn(context.Background(), stream, handler)

			// 在 goroutine 中处理连接
			go func(conn net.Conn, rpcConn *jsonrpc2.Conn) {
				<-rpcConn.DisconnectNotify()
				logrus.Infof("TCP 连接已关闭：%s", conn.RemoteAddr())
				conn.Close()
			}(conn, rpcConn)
		}
	case "websocket":
		addr := fmt.Sprintf("%s:%d", *address, *port)
		logrus.Infof("使用 WebSocket 协议在 %s 启动 LSP 服务器", addr)
		logrus.Error("不支持 WebSocket 协议")
		os.Exit(1)
	default:
		logrus.Errorf("不支持的协议：%s", *protocol_)
		os.Exit(1)
	}
}

// stdrwc 实现 io.ReadWriteCloser 接口
type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}

// 处理 LSP 请求
func handleRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	// 设置全局连接，用于发送通知
	globalConn = conn

	// 打印完整的 JSON 参数
	paramsJSON, _ := json.Marshal(req.Params)
	logrus.Debugf("LSP 请求：%s\n参数：%s", req.Method, string(paramsJSON))

	switch req.Method {
	case "initialize":
		return handleInitialize(req)
	case "initialized":
		return handleInitialized(req)
	case "shutdown":
		return handleShutdown(req)
	case "textDocument/didOpen":
		return handleTextDocumentDidOpen(conn, req)
	case "textDocument/didChange":
		return handleTextDocumentDidChange(conn, req)
	case "textDocument/didClose":
		return handleTextDocumentDidClose(conn, req)
	case "textDocument/completion":
		return handleTextDocumentCompletion(req)
	case "textDocument/hover":
		return handleTextDocumentHover(req)
	case "textDocument/definition":
		return handleTextDocumentDefinition(req)
	case "textDocument/documentSymbol":
		return handleTextDocumentDocumentSymbol(req)
	case "$/setTrace":
		return handleSetTrace(req)
	default:
		logrus.Warnf("未知方法：%s", req.Method)
		return nil, nil
	}
}

// 中间件：统一的 LSP 协议日志处理
func logLSPCommunication(method string, isRequest bool, params interface{}) {
	msgType := "通知"
	if isRequest {
		msgType = "请求"
	}
	logrus.Infof("LSP %s：%s", msgType, method)
	if params != nil {
		paramsJSON, _ := json.Marshal(params)
		logrus.Debugf("参数：%s", string(paramsJSON))
	}
}

func logLSPResponse(method string, result interface{}, err error) {
	if err != nil {
		logrus.Errorf("LSP %s 失败：%v", method, err)
	} else {
		if result != nil {
			resultJSON, _ := json.Marshal(result)
			logrus.Infof("结果：%s", string(resultJSON))
		}
	}
}

func showHelp() {
	fmt.Printf("%s v%s\n", lsName, lsVersion)
	fmt.Printf("Origami 语言的语言服务器协议实现。\n")
	fmt.Printf("\n")
	fmt.Printf("用法：\n")
	fmt.Printf("  %s [选项]\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("选项：\n")
	fmt.Printf("  -version        显示版本信息\n")
	fmt.Printf("  -help           显示帮助信息\n")
	fmt.Printf("  -test           运行定义跳转测试\n")
	fmt.Printf("  -protocol       协议类型 (stdio, tcp, websocket) [默认: stdio]\n")
	fmt.Printf("  -address        绑定地址 (tcp/websocket) [默认: localhost]\n")
	fmt.Printf("  -port           绑定端口 (tcp/websocket) [默认: 8800]\n")
	fmt.Printf("  -log-level      日志级别 (0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug) [默认: 5]\n")
	fmt.Printf("  -log-file       日志文件路径 [默认: lsp.log]\n")
	fmt.Printf("  -console-log    在 stdio 模式下启用控制台日志 [默认: true]\n")
	fmt.Printf("  -scan-dir       扫描指定目录中的所有 .zy 文件 [可选]\n")
	fmt.Printf("\n")
	fmt.Printf("示例：\n")
	fmt.Printf("  # 使用 stdio 启动（默认，控制台日志启用）\n")
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("  # 使用 stdio 启动，禁用控制台日志\n")
	fmt.Printf("  %s -console-log=false\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("  # 扫描指定目录中的所有 .zy 文件\n")
	fmt.Printf("  %s -scan-dir /path/to/project\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("  # 使用 TCP 在端口 8080 启动\n")
	fmt.Printf("  %s -protocol tcp -port 8080\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("  # 使用 WebSocket 和自定义日志启动\n")
	fmt.Printf("  %s -protocol websocket -port 9000 -log-level 4 -log-file lsp.log\n", os.Args[0])
}

func testDefinitionJumpFeature() {
	// 测试定义跳转功能
	logrus.Infof("测试定义跳转功能...")
	// 这里可以添加具体的测试逻辑
}
