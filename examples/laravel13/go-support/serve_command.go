package gosupport

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/examples/laravel13/go-support/httpkernel"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/perfmon"
	"github.com/php-any/origami/runtime"
	illuminatehttp "github.com/php-any/origami/std/illuminate/http"
	phpcore "github.com/php-any/origami/std/php/core"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
	"github.com/php-any/origami/std/vendoraccel"
)

const (
	serveCommandClassName = "Illuminate\\Foundation\\Console\\ServeCommand"
	// PHP 网页 SAPI：max_execution_time 默认 30 秒；到期中断请求并关闭连接。
	serveMaxExecutionTime = 30 * time.Second
	serveIdleTimeout      = 60 * time.Second
)

// ServeCommandClass 在 Laravel 加载官方 ServeCommand.php 前注册同名 Go 类，
// 保留 Artisan 命令生命周期，同时用 Go net/http 替代 php -S。
type ServeCommandClass struct {
	node.Node
	properties []data.Property
}

func NewServeCommandClass() data.ClassStmt {
	return &ServeCommandClass{
		properties: []data.Property{
			node.NewProperty(nil, "name", "protected", false, node.NewStringLiteral(nil, "serve")),
			node.NewProperty(nil, "description", "protected", false, node.NewStringLiteral(nil, "Serve the application with the Origami HTTP kernel")),
			node.NewProperty(nil, "aliases", "protected", false, data.NewArrayValue(nil)),
		},
	}
}

func (c *ServeCommandClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ServeCommandClass) GetName() string { return serveCommandClassName }
func (c *ServeCommandClass) GetExtend() *string {
	parent := "Illuminate\\Console\\Command"
	return &parent
}
func (c *ServeCommandClass) GetImplements() []string          { return nil }
func (c *ServeCommandClass) GetConstruct() data.Method        { return &serveConstructMethod{} }
func (c *ServeCommandClass) GetPropertyList() []data.Property { return c.properties }

func (c *ServeCommandClass) GetProperty(name string) (data.Property, bool) {
	for _, property := range c.properties {
		if property.GetName() == name {
			return property, true
		}
	}
	return nil, false
}

func (c *ServeCommandClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "handle":
		return &serveHandleMethod{}, true
	case "getOptions":
		return &serveGetOptionsMethod{}, true
	default:
		return nil, false
	}
}

func (c *ServeCommandClass) GetMethods() []data.Method {
	return []data.Method{&serveHandleMethod{}, &serveGetOptionsMethod{}}
}

type serveConstructMethod struct{}

func (m *serveConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("serve: 无法创建命令上下文"))
	}
	parent, control := ctx.GetVM().GetOrLoadClass("Illuminate\\Console\\Command")
	if control != nil {
		return nil, control
	}
	constructor := parent.GetConstruct()
	if constructor == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("serve: Illuminate Console Command 构造函数不存在"))
	}

	callCtx := classCtx.ClassValue.CreateContext(constructor.GetVariables())
	parentCtx, ok := callCtx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("serve: 无法创建父命令上下文"))
	}
	parentCtx.SelfClass = parent
	parentCtx.StaticClass = classCtx.ClassValue.Class
	return constructor.Call(parentCtx)
}

func (m *serveConstructMethod) GetName() string               { return "__construct" }
func (m *serveConstructMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *serveConstructMethod) GetIsStatic() bool             { return false }
func (m *serveConstructMethod) GetParams() []data.GetValue    { return nil }
func (m *serveConstructMethod) GetVariables() []data.Variable { return nil }
func (m *serveConstructMethod) GetReturnType() data.Types     { return nil }

type serveHandleMethod struct{}

func (m *serveHandleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	host, port, err := serveAddress(os.Args)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	base, ok := ctx.GetVM().(*runtime.VM)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("serve: 需要常驻 runtime.VM"))
	}
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || classCtx.ClassValue == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("serve: 无法读取 Laravel Application"))
	}
	appValue, control := classCtx.ClassValue.GetProperty("laravel")
	if control != nil {
		return nil, control
	}
	app, ok := appValue.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("serve: Laravel Application 未绑定到命令"))
	}
	BootLog("ServeCommand.handle (artisan console boot finished)")
	if err := runLaravelHTTPServer(host, port, base, app); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewIntValue(0), nil
}

func (m *serveHandleMethod) GetName() string               { return "handle" }
func (m *serveHandleMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *serveHandleMethod) GetIsStatic() bool             { return false }
func (m *serveHandleMethod) GetParams() []data.GetValue    { return nil }
func (m *serveHandleMethod) GetVariables() []data.Variable { return nil }
func (m *serveHandleMethod) GetReturnType() data.Types     { return data.NewBaseType("int") }

type serveGetOptionsMethod struct{}

func (m *serveGetOptionsMethod) Call(data.Context) (data.GetValue, data.Control) {
	option := func(values ...data.Value) data.Value {
		return data.NewArrayValue(values)
	}
	null := data.NewNullValue()
	return data.NewArrayValue([]data.Value{
		option(data.NewStringValue("host"), null, data.NewIntValue(4), data.NewStringValue("The host address to serve the application on"), data.NewStringValue("127.0.0.1")),
		option(data.NewStringValue("port"), null, data.NewIntValue(4), data.NewStringValue("The port to serve the application on"), null),
		option(data.NewStringValue("tries"), null, data.NewIntValue(4), data.NewStringValue("The max number of ports to attempt to serve from"), data.NewIntValue(10)),
		option(data.NewStringValue("no-reload"), null, data.NewIntValue(1), data.NewStringValue("Do not reload the development server on .env file changes")),
	}), nil
}

func (m *serveGetOptionsMethod) GetName() string               { return "getOptions" }
func (m *serveGetOptionsMethod) GetModifier() data.Modifier    { return data.ModifierProtected }
func (m *serveGetOptionsMethod) GetIsStatic() bool             { return false }
func (m *serveGetOptionsMethod) GetParams() []data.GetValue    { return nil }
func (m *serveGetOptionsMethod) GetVariables() []data.Variable { return nil }
func (m *serveGetOptionsMethod) GetReturnType() data.Types     { return data.NewBaseType("array") }

type laravelHTTPKernel struct {
	base       *runtime.VM
	parser     *parser.Parser
	app        *data.ClassValue
	kernel     *data.ClassValue
	initError  data.Control
	initOnce   sync.Once
	connTimers sync.Map // net.Conn -> *time.Timer，到期关闭连接
}

func newLaravelHTTPKernel(base *runtime.VM, app *data.ClassValue) *laravelHTTPKernel {
	return &laravelHTTPKernel{
		base: base,
		app:  app,
	}
}

func (k *laravelHTTPKernel) trackConn(c net.Conn, s http.ConnState) {
	switch s {
	case http.StateActive:
		k.armConnTimeout(c)
	case http.StateIdle, http.StateClosed, http.StateHijacked:
		k.disarmConnTimeout(c)
	}
}

func (k *laravelHTTPKernel) armConnTimeout(c net.Conn) {
	k.disarmConnTimeout(c)
	t := time.AfterFunc(serveMaxExecutionTime, func() {
		_ = c.Close()
	})
	k.connTimers.Store(c, t)
}

func (k *laravelHTTPKernel) disarmConnTimeout(c net.Conn) {
	if v, ok := k.connTimers.LoadAndDelete(c); ok {
		if t, ok := v.(*time.Timer); ok {
			t.Stop()
		}
	}
}

// beginPHPWebTimeLimit 按 PHP 网页 SAPI 给当前请求套上 max_execution_time=30。
// 截止时间写在当前 goroutine 的 CallState 上，禁止改进程级槽。
func beginPHPWebTimeLimit() func() {
	sec := phpcore.DefaultHTTPMaxExecutionTime
	phpcore.SetExecutionDeadline(sec)
	return func() {
		phpcore.SetExecutionDeadline(0)
	}
}

func closeTimedOutRequest(w http.ResponseWriter) {
	w.Header().Set("Connection", "close")
	http.Error(w, "Fatal error: Maximum execution time of 30 seconds exceeded", http.StatusInternalServerError)
}

func (k *laravelHTTPKernel) ensureBase() data.Control {
	k.initOnce.Do(func() {
		if k.base == nil || k.app == nil {
			k.initError = data.NewErrorThrow(nil, fmt.Errorf("Laravel Application 未初始化"))
			return
		}
		k.parser = k.base.GetParser()
		var control data.Control
		k.kernel, control = httpkernel.Resolve(k.app)
		if control != nil {
			k.initError = control
			return
		}
		if control = httpkernel.Bootstrap(k.kernel); control != nil {
			k.initError = control
		}
	})
	return k.initError
}

func (k *laravelHTTPKernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if control := k.ensureBase(); control != nil {
		if k.parser != nil {
			k.parser.ShowControl(control)
		}
		http.Error(w, "Laravel bootstrap failed", http.StatusInternalServerError)
		return
	}

	// 对齐 artisan serve 的 server.php：public 下存在的静态文件直接返回。
	if k.servePublicFile(w, r) {
		return
	}
	if k.serveLivewireDist(w, r) {
		return
	}
	if k.serveMissingStaticAsset(w, r) {
		return
	}

	if r.Context().Err() != nil {
		return
	}

	// 每个请求并行 Handle：输出/调用栈/Container::$instance 按 goroutine 隔离，
	// Application/Router 用 clone 沙箱，互不排队。
	handleCtx, handleCancel := context.WithTimeout(context.Background(), serveMaxExecutionTime)
	defer handleCancel()

	outcome := "ok"
	span := perfmon.BeginRequest(r.Method, r.URL.Path)
	defer func() { span.End(outcome) }()

	// 对齐 Octane Worker：每请求独立 ob_start，残留 echo 并入 HTTP body，不写进程 stdout。
	defer runtime.BeginRequestOutput()()
	runtime.StartRequestOutputBuffer()
	runtime.MuteRequestStdout()
	defer runtime.BeginRequestDeadline(handleCtx)()
	defer beginPHPWebTimeLimit()()

	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		if data.IsRequestCanceled(rec) {
			outcome = "timeout"
			if r.Context().Err() != nil {
				return
			}
			closeTimedOutRequest(w)
			return
		}
		outcome = "panic"
		fmt.Fprintf(os.Stderr, "origami ServeHTTP panic: %v\n%s\n", rec, debug.Stack())
		if r.Context().Err() != nil {
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}()

	requestCtx := runtime.NewTempVM(k.base).CreateContext(nil)
	tSandbox := perfmon.Now()
	kernel := httpkernel.Sandbox(requestCtx, k.kernel)
	perfmon.NoteFileRun("sandbox", perfmon.Since(tSandbox))
	defer httpkernel.ResetViewEngines(requestCtx, kernel)
	request := illuminatehttp.NewIlluminateRequestValue(requestCtx, r)
	var response data.GetValue
	tHandle := time.Now()
	response, control := httpkernel.Handle(requestCtx, kernel, request)
	perfmon.NoteFileRun("handle", time.Since(tHandle))
	if firstHTTPLogged.CompareAndSwap(false, true) {
		BootLog("first-request Handle " + time.Since(tHandle).Truncate(time.Millisecond).String() + " path=" + r.URL.Path)
	}
	leftover := runtime.TakeRequestOutput()
	sentOK := false
	if control == nil {
		var sent *data.ClassValue
		sent, control = httpfoundation.SendResponseTo(w, response, leftover)
		if control == nil {
			sentOK = true
			if term := httpkernel.Terminate(requestCtx, kernel, request, sent); term != nil {
				if k.parser != nil {
					k.parser.ShowControl(term)
				}
				// 响应已写出：terminate（Telescope 等）失败不能再 WriteHeader。
				return
			}
		}
	}
	// PHP exit/die 在请求生命周期内是正常结束，不当作错误。
	if exit, ok := control.(data.ExitControl); ok && exit.IsExit() {
		control = nil
	}
	if control != nil {
		if outcome == "ok" {
			outcome = "error"
		}
		if httpfoundation.IsClientAbortControl(control) {
			return
		}
		if k.parser != nil {
			k.parser.ShowControl(control)
		}
		if !sentOK {
			http.Error(w, "Laravel request failed", http.StatusInternalServerError)
		}
	}
}

// servePublicFile 模拟 Laravel resources/server.php：若 public{$uri} 存在则直出。
func (k *laravelHTTPKernel) servePublicFile(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	uri := path.Clean("/" + r.URL.Path)
	if uri == "/" {
		return false
	}
	root, err := os.Getwd()
	if err != nil {
		return false
	}
	publicRoot := filepath.Clean(filepath.Join(root, "public"))
	full := filepath.Clean(filepath.Join(publicRoot, filepath.FromSlash(uri)))
	rel, err := filepath.Rel(publicRoot, full)
	if err != nil || strings.HasPrefix(rel, "..") {
		return false
	}
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		return false
	}
	http.ServeFile(w, r, full)
	return true
}

// serveMissingStaticAsset 对 public 下不存在的静态扩展名直接 404。
// 浏览器会狂刷 /favicon.ico；若走完整 Kernel，串行锁会把后续 /admin 等请求堵住。
func (k *laravelHTTPKernel) serveMissingStaticAsset(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	uri := path.Clean("/" + r.URL.Path)
	if strings.HasPrefix(strings.ToLower(uri), "/.well-known/") {
		http.NotFound(w, r)
		return true
	}
	switch strings.ToLower(path.Base(uri)) {
	case "favicon.ico", "robots.txt", "apple-touch-icon.png", "apple-touch-icon-precomposed.png":
	default:
		return false
	}
	http.NotFound(w, r)
	return true
}

// serveLivewireDist 把 /livewire-{hash}/livewire(.min).js 映射到 vendor 发行文件。
// Livewire 4 默认走 hashed 路由 + BinaryFileResponse；直出 dist 可避免 file response 路径 500。
func (k *laravelHTTPKernel) serveLivewireDist(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	uri := path.Clean("/" + r.URL.Path)
	dir, base := path.Dir(uri), path.Base(uri)
	if !strings.HasPrefix(path.Base(dir), "livewire-") {
		return false
	}
	switch base {
	case "livewire.js", "livewire.min.js", "livewire.js.map", "livewire.min.js.map",
		"livewire.csp.js", "livewire.csp.min.js", "livewire.csp.min.js.map",
		"livewire.esm.js", "livewire.csp.esm.js":
	default:
		return false
	}
	root, err := os.Getwd()
	if err != nil {
		return false
	}
	full := filepath.Join(root, "vendor", "livewire", "livewire", "dist", base)
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		return false
	}
	http.ServeFile(w, r, full)
	return true
}

func runLaravelHTTPServer(host string, port int, base *runtime.VM, app *data.ClassValue) error {
	// 默认 VM 遇到 throw 会 os.Exit(1)。常驻 serve 下单请求失败不能把整个进程打死。
	base.SetThrowControl(func(acl data.Control) {
		if p := base.GetParser(); p != nil {
			p.ShowControl(acl)
		}
	})

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	ln, err := listenTCP(addr)
	if err != nil {
		if isAddrInUse(err) {
			return fmt.Errorf("serve: 端口 %s 已被占用，拒绝启动", addr)
		}
		return fmt.Errorf("serve: 无法监听 %s: %w", addr, err)
	}

	k := newLaravelHTTPKernel(base, app)
	seed, err := http.NewRequest(http.MethodGet, "http://"+addr+"/", nil)
	if err != nil {
		_ = ln.Close()
		return fmt.Errorf("serve: 无法构造启动 Request: %w", err)
	}
	seed.Host = addr
	if ctl := httpkernel.BindRequest(app, illuminatehttp.NewIlluminateRequestValue(app, seed)); ctl != nil {
		_ = ln.Close()
		return fmt.Errorf("serve: 绑定启动 Request 失败: %s", ctl.AsString())
	}
	tBoot := time.Now()
	if ctl := k.ensureBase(); ctl != nil {
		_ = ln.Close()
		return fmt.Errorf("serve: HTTP Kernel bootstrap 失败: %s", ctl.AsString())
	}
	BootLog("HTTP Kernel resolve+bootstrap " + time.Since(tBoot).Truncate(time.Millisecond).String())

	fmt.Printf("   INFO  Server running on [http://%s].\n\n", addr)
	fmt.Println("  Press Ctrl+C to stop the server")
	BootLog("listen+print (HTTP accept loop starting)")

	// 预热不阻塞监听：完整 classmap 要数分钟，且会误加载依赖 PHPUnit 的 Testing 类。
	if vendoraccel.ShouldWarmup() {
		go func() {
			if root, err := os.Getwd(); err == nil {
				vendoraccel.WarmupVendorClassmap(base, root)
			}
		}()
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           k,
		ConnState:         k.trackConn,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       serveIdleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Serve(ln)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("serve: HTTP 服务失败: %w", err)
		}
		return nil
	case <-stop:
		fmt.Println("\n   INFO  Stopping server...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("serve: 停止 HTTP 服务失败: %w", err)
	}
	if err := <-serverErr; err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serve: HTTP 服务失败: %w", err)
	}
	return nil
}

func serveAddress(args []string) (string, int, error) {
	host, port := "127.0.0.1", 8000
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case strings.HasPrefix(arg, "--host="):
			host = strings.TrimPrefix(arg, "--host=")
		case arg == "--host" && i+1 < len(args):
			i++
			host = args[i]
		case strings.HasPrefix(arg, "--port="):
			value, err := strconv.Atoi(strings.TrimPrefix(arg, "--port="))
			if err != nil {
				return "", 0, fmt.Errorf("serve: 无效端口 %q", arg)
			}
			port = value
		case arg == "--port" && i+1 < len(args):
			i++
			value, err := strconv.Atoi(args[i])
			if err != nil {
				return "", 0, fmt.Errorf("serve: 无效端口 %q", args[i])
			}
			port = value
		}
	}
	if parsedHost, parsedPort, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
		if parsed, parseErr := strconv.Atoi(parsedPort); parseErr == nil {
			port = parsed
		}
	}
	if port < 0 || port > 65535 {
		return "", 0, fmt.Errorf("serve: port 必须在 0 到 65535 之间")
	}
	return host, port, nil
}
