package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/fpm"
	"github.com/php-any/origami/std/system"
)

func main() {
	host, port := "0.0.0.0", 8888
	if v := strings.TrimSpace(os.Getenv("HTTP_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			port = n
		}
	}
	if len(os.Args) > 1 {
		if n, err := strconv.Atoi(os.Args[1]); err == nil && n > 0 {
			port = n
		}
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法获取工作目录: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(filepath.Join(root, "index.php")); err != nil {
		fmt.Fprintf(os.Stderr, "当前目录缺少 index.php，请在 examples/bbs1org 下运行\n")
		os.Exit(1)
	}

	base, p := buildVM()
	handler := &bbsHandler{base: base, parser: p, root: root}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	fmt.Printf("bbs1org (Origami) 监听 http://%s\n", addr)
	fmt.Println("静态资源模式: memory (无 Range/206)")
	fmt.Println("首次访问将进入安装向导；按 Ctrl+C 停止")

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ 端口 %d 已被占用，无法启动。\n", port)
		fmt.Fprintf(os.Stderr, "   请先执行: killall -9 bbs1org\n")
		fmt.Fprintf(os.Stderr, "   或: lsof -nP -iTCP:%d -sTCP:LISTEN\n", port)
		fmt.Fprintf(os.Stderr, "   原始错误: %v\n\n", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(ln)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP 服务失败: %v\n", err)
			os.Exit(1)
		}
	case <-stop:
		fmt.Println("\n正在停止服务...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}
}

func buildVM() (*runtime.VM, *parser.Parser) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	httplib.Load(vm)
	websocket.Load(vm)
	netannotation.Load(vm)
	system.Load(vm)
	return vm, p
}

type bbsHandler struct {
	base   *runtime.VM
	parser *parser.Parser
	root   string
}

func (h *bbsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.serveStatic(w, r) {
		return
	}

	recorder := httptest.NewRecorder()
	reqVM := fpm.New(h.base, func(s string) { _, _ = io.WriteString(recorder, s) })
	reqVM.BindHTTP(r, recorder)
	_ = reqVM.SetConstant("PHP_SAPI", data.NewStringValue("cli-server"))

	_, ctrl := reqVM.LoadAndRun(filepath.Join(h.root, "index.php"))
	if ctrl != nil {
		if _, ok := ctrl.(data.ExitControl); ok {
			ctrl = nil
		}
	}
	if ctrl == nil {
		if thrown := reqVM.TakeThrow(); thrown != nil {
			if _, ok := thrown.(data.ExitControl); ok {
				thrown = nil
			} else {
				ctrl = thrown
			}
		}
	}
	if ctrl != nil {
		h.parser.ShowControl(ctrl)
		if recorder.Body.Len() == 0 {
			http.Error(recorder, "bbs1org request failed", http.StatusInternalServerError)
		}
	}
	flushRecorder(w, recorder)
}

func (h *bbsHandler) serveStatic(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	uri := path.Clean("/" + r.URL.Path)
	prefixes := []string{"/app/assets/", "/app/avatars/", "/app/upload/"}
	matched := false
	for _, prefix := range prefixes {
		if strings.HasPrefix(uri, prefix) {
			matched = true
			break
		}
	}
	if !matched {
		return false
	}
	full := filepath.Clean(filepath.Join(h.root, filepath.FromSlash(uri)))
	rel, err := filepath.Rel(h.root, full)
	if err != nil || strings.HasPrefix(rel, "..") {
		return false
	}
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		return false
	}

	// 不用 http.ServeFile：其 sendfile/Range(206) 路径在本环境下会把响应截断到 512 字节，
	// 触发浏览器 net::ERR_CONTENT_LENGTH_MISMATCH。
	body, err := os.ReadFile(full)
	if err != nil {
		http.Error(w, "file read failed", http.StatusInternalServerError)
		return true
	}
	ctype := mime.TypeByExtension(filepath.Ext(full))
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Last-Modified", info.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "no-cache")
	// 便于确认打到的是新静态逻辑（旧 ServeFile 会带 Accept-Ranges / 可能返回 206）
	w.Header().Set("X-Origami-Static", "memory")
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	}
	return true
}

func flushRecorder(w http.ResponseWriter, recorder *httptest.ResponseRecorder) {
	for name, values := range recorder.Header() {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	code := recorder.Code
	if code == 0 {
		code = http.StatusOK
	}
	w.WriteHeader(code)
	_, _ = w.Write(recorder.Body.Bytes())
}
