package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

const (
	lsName    = "Origami Language Server"
	lsVersion = "1.0.0"
)

var (
	version   = flag.Bool("version", false, "Show version information")
	help      = flag.Bool("help", false, "Show help information")
	test      = flag.Bool("test", false, "Run definition jump test")
	protocol_ = flag.String("protocol", "stdio", "Protocol to use (stdio, tcp, websocket)")
	address   = flag.String("address", "localhost", "Address to bind to (for tcp/websocket)")
	port      = flag.Int("port", 8080, "Port to bind to (for tcp/websocket)")
	logLevel  = flag.Int("log-level", 1, "Log level (0=off, 1=error, 2=warn, 3=info, 4=debug)")
	enableLog = flag.Bool("enable-log", true, "Enable logging")
	logFile   = flag.String("log-file", "", "Log file path (empty for stderr)")
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

type DocumentInfo struct {
	Content string
	Version int32
	AST     interface{}
	Parser  interface{}
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s\n", lsName, lsVersion)
		return
	}

	if *help {
		showHelp()
		return
	}

	if *logLevel > 0 {
		fmt.Fprintf(os.Stderr, "日志级别: %d\n", *logLevel)
	}

	if *test {
		fmt.Println("=== 运行定义跳转测试 ===")
		testDefinitionJumpFeature()
		fmt.Println("=== 测试完成 ===")
		return
	}

	// 初始化全局 LspVM
	globalLspVM = NewLspVM()

	// 创建 JSON-RPC 2.0 处理器
	handler := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		return handleRequest(ctx, conn, req)
	})

	// 创建连接
	var conn *jsonrpc2.Conn
	switch *protocol_ {
	case "stdio":
		if *logLevel > 0 {
			fmt.Fprintf(os.Stderr, "[INFO] Starting LSP server with stdio protocol\n")
		}
		stream := jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{})
		conn = jsonrpc2.NewConn(context.Background(), stream, handler)
		<-conn.DisconnectNotify()
	case "tcp":
		addr := fmt.Sprintf("%s:%d", *address, *port)
		if *logLevel > 0 {
			fmt.Fprintf(os.Stderr, "[INFO] Starting LSP server with TCP protocol on %s\n", addr)
		}

		// 创建 TCP 监听器
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Failed to start TCP server: %v\n", err)
			os.Exit(1)
		}
		defer listener.Close()

		if *logLevel > 0 {
			fmt.Fprintf(os.Stderr, "[INFO] TCP server listening on %s\n", addr)
		}

		// 接受连接
		for {
			conn, err := listener.Accept()
			if err != nil {
				if *logLevel > 1 {
					fmt.Fprintf(os.Stderr, "[WARNING] Failed to accept connection: %v\n", err)
				}
				continue
			}

			if *logLevel > 2 {
				fmt.Fprintf(os.Stderr, "[INFO] New TCP connection from %s\n", conn.RemoteAddr())
			}

			// 为每个连接创建一个新的 JSON-RPC 连接
			stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})
			rpcConn := jsonrpc2.NewConn(context.Background(), stream, handler)

			// 在 goroutine 中处理连接
			go func(conn net.Conn, rpcConn *jsonrpc2.Conn) {
				<-rpcConn.DisconnectNotify()
				if *logLevel > 2 {
					fmt.Fprintf(os.Stderr, "[INFO] TCP connection closed: %s\n", conn.RemoteAddr())
				}
				conn.Close()
			}(conn, rpcConn)
		}
	case "websocket":
		addr := fmt.Sprintf("%s:%d", *address, *port)
		if *logLevel > 0 {
			fmt.Fprintf(os.Stderr, "[INFO] Starting LSP server with WebSocket protocol on %s\n", addr)
		}
		fmt.Fprintf(os.Stderr, "WebSocket protocol not supported\n")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported protocol: %s\n", *protocol_)
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
	// 打印完整的 JSON 参数
	if *logLevel > 3 {
		paramsJSON, _ := json.MarshalIndent(req.Params, "", "  ")
		fmt.Fprintf(os.Stderr, "[DEBUG] LSP Request: %s\nParams: %s\n", req.Method, string(paramsJSON))
	}

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
		if *logLevel > 1 {
			fmt.Fprintf(os.Stderr, "[WARNING] Unknown method: %s\n", req.Method)
		}
		return nil, nil
	}
}

// 中间件：统一的 LSP 协议日志处理
func logLSPCommunication(method string, isRequest bool, params interface{}) {
	if *logLevel > 2 {
		msgType := "notification"
		if isRequest {
			msgType = "request"
		}
		fmt.Fprintf(os.Stderr, "[INFO] LSP %s: %s\n", msgType, method)
	}
	if *logLevel > 2 && params != nil {
		paramsJSON, _ := json.MarshalIndent(params, "", "  ")
		fmt.Fprintf(os.Stderr, "[DEBUG] Params: %s\n", string(paramsJSON))
	}
}

func logLSPResponse(method string, result interface{}, err error) {
	if err != nil {
		if *logLevel > 1 {
			fmt.Fprintf(os.Stderr, "[ERROR] LSP %s failed: %v\n", method, err)
		}
	} else {
		if *logLevel > 2 {
			fmt.Fprintf(os.Stderr, "[INFO] LSP %s completed\n", method)
		}
		if *logLevel > 3 && result != nil {
			resultJSON, _ := json.MarshalIndent(result, "", "  ")
			fmt.Fprintf(os.Stderr, "[DEBUG] Result: %s\n", string(resultJSON))
		}
	}
}

func showHelp() {
	fmt.Printf("%s v%s\n\n", lsName, lsVersion)
	fmt.Println("A Language Server Protocol implementation for Origami language.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Start with stdio (default)")
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Start with TCP on port 8080")
	fmt.Printf("  %s -protocol tcp -port 8080\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Start with WebSocket and custom logging")
	fmt.Printf("  %s -protocol websocket -port 9000 -log-level 4 -log-file lsp.log\n", os.Args[0])
}

func testDefinitionJumpFeature() {
	// 测试定义跳转功能
	fmt.Println("测试定义跳转功能...")
	// 这里可以添加具体的测试逻辑
}
