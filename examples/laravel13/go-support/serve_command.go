package gosupport

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/examples/laravel13/go-support/httpfoundation"
	"github.com/php-any/origami/examples/laravel13/go-support/httpkernel"
	"github.com/php-any/origami/examples/laravel13/go-support/requestvm"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

const serveCommandClassName = "Illuminate\\Foundation\\Console\\ServeCommand"

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
	base      *runtime.VM
	parser    *parser.Parser
	app       *data.ClassValue
	kernel    *data.ClassValue
	initError data.Control
	initOnce  sync.Once
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

	recorder := httptest.NewRecorder()
	reqVM := requestvm.New(k.base, func(s string) { _, _ = io.WriteString(recorder, s) })
	reqVM.BindHTTP(r, recorder)

	requestCtx := reqVM.CreateContext(nil)
	request := httpfoundation.NewIlluminateRequestValue(requestCtx, r)
	var response data.GetValue
	response, control := httpkernel.Handle(requestCtx, k.kernel, request)
	if control == nil {
		var sent *data.ClassValue
		sent, control = httpfoundation.SendResponse(response)
		if control == nil {
			control = httpkernel.Terminate(requestCtx, k.kernel, request, sent)
		}
	}
	// PHP exit/die 在请求生命周期内是正常结束，不当作错误。
	if exit, ok := control.(data.ExitControl); ok && exit.IsExit() {
		control = nil
	}
	if control != nil {
		k.parser.ShowControl(control)
		http.Error(recorder, "Laravel request failed", http.StatusInternalServerError)
	} else if thrown := reqVM.TakeThrow(); thrown != nil {
		if exit, ok := thrown.(data.ExitControl); ok && exit.IsExit() {
			thrown = nil
		}
		if thrown != nil {
			k.parser.ShowControl(thrown)
			http.Error(recorder, "Laravel request failed", http.StatusInternalServerError)
		}
	}
	k.flushRecorder(w, recorder)
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

func (k *laravelHTTPKernel) flushRecorder(w http.ResponseWriter, recorder *httptest.ResponseRecorder) {
	for name, values := range recorder.Header() {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(recorder.Code)
	_, _ = w.Write(recorder.Body.Bytes())
}

func runLaravelHTTPServer(host string, port int, base *runtime.VM, app *data.ClassValue) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	fmt.Printf("   INFO  Server running on [http://%s].\n\n", addr)
	fmt.Println("  Press Ctrl+C to stop the server")

	server := &http.Server{
		Addr: addr,
		Handler: &laravelHTTPKernel{
			base: base,
			app:  app,
		},
		ReadHeaderTimeout: 10 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
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
