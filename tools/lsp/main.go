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
	"sync"

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

// 全局变量
var (
	globalLspVM    *LspVM
	documents      = make(map[string]*DocumentInfo)
	requestCancels sync.Map // map[string]context.CancelFunc
)

// CancelParams 定义取消请求参数
type CancelParams struct {
	ID jsonrpc2.ID `json:"id"`
}

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
		// TCP/WebSocket 模式下，同时输出到文件和 stdout（因为 TCP 模式下 stdout 不用于通信）
		var writers []io.Writer

		// 始终输出到 stdout
		writers = append(writers, os.Stdout)

		if dir := filepath.Dir(*logFile); dir != "." && dir != "" {
			_ = os.MkdirAll(dir, 0755)
		}
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// 无法写入日志文件时，输出到 stderr
			fmt.Fprintf(os.Stderr, "无法写入日志文件，将输出到 stderr：%s，原因：%v\n", *logFile, err)
		} else {
			defer file.Close()
			writers = append(writers, file)
		}

		output = io.MultiWriter(writers...)
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
			logrus.Infof("日志输出: stdout + %s", *logFile)
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

	// 处理取消请求
	if req.Method == "$/cancelRequest" {
		var params CancelParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			logrus.Warnf("解析取消请求参数失败：%v", err)
			return nil, nil
		}
		// 获取要取消的请求 ID
		idStr := fmt.Sprintf("%v", params.ID)
		if cancel, ok := requestCancels.Load(idStr); ok {
			logrus.Infof("正在取消请求 ID: %s", idStr)
			cancel.(context.CancelFunc)()
			requestCancels.Delete(idStr)
		} else {
			logrus.Debugf("尝试取消不存在或已完成的请求 ID: %s", idStr)
		}
		return nil, nil
	}

	// 为每个请求创建可取消的上下文
	reqCtx, cancel := context.WithCancel(ctx)
	// 如果是请求（有 ID），则记录 cancel 函数以便取消
	if !req.Notif {
		idStr := fmt.Sprintf("%v", req.ID)
		requestCancels.Store(idStr, cancel)
		defer func() {
			requestCancels.Delete(idStr)
			cancel() // 确保请求完成后释放资源
		}()
	} else {
		defer cancel()
	}

	switch req.Method {
	case "initialize":
		// 初始化语言服务器（握手）
		return handleInitialize(reqCtx, req)
	case "initialized":
		// 客户端确认初始化完成通知
		return handleInitialized(reqCtx, req)
	case "shutdown":
		// 关闭服务器请求
		return handleShutdown(reqCtx, req)
	case "textDocument/didOpen":
		// 文档打开通知
		return handleTextDocumentDidOpen(reqCtx, conn, req)
	case "textDocument/didChange":
		// 文档修改通知
		return handleTextDocumentDidChange(reqCtx, conn, req)
	case "textDocument/didClose":
		// 文档关闭通知
		return handleTextDocumentDidClose(reqCtx, conn, req)
	case "textDocument/completion":
		// 代码补全请求
		return handleTextDocumentCompletion(reqCtx, req)
	case "textDocument/hover":
		// 悬停提示请求
		return handleTextDocumentHover(reqCtx, req)
	case "textDocument/definition":
		// 定义跳转请求
		return handleTextDocumentDefinition(reqCtx, req)
	case "textDocument/documentSymbol":
		// 文档符号请求（大纲）
		return handleTextDocumentDocumentSymbol(reqCtx, req)
	case "$/setTrace":
		// 设置追踪级别
		return handleSetTrace(reqCtx, req)
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
